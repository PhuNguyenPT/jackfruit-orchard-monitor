#include <Arduino.h>
#include <HardwareSerial.h>
#include <ModbusMaster.h>
#include <PubSubClient.h>
#include <WiFi.h>
#include <WiFiClientSecure.h>

#include "broker_config.h"
#include "config.h"
#include "gpio.h"
#include "sht40.h"

// ---------------------------------------------------------------------------
// Compile-time constants (replaces magic numbers throughout)
// ---------------------------------------------------------------------------
namespace {

constexpr uint32_t WIFI_INIT_DELAY_MS = 100U;
constexpr uint32_t WIFI_RECONNECT_DELAY_MS = 500U;
constexpr uint32_t MQTT_RETRY_DELAY_MS = 5000U;
constexpr uint32_t SERIAL_INIT_DELAY_MS = 10U;

constexpr float SENSOR_SCALE = 10.0F;

constexpr size_t TOPIC_BUF_SIZE = 50U;
constexpr size_t PAYLOAD_BUF_SIZE = 100U;

}  // namespace

// ---------------------------------------------------------------------------
// Global objects
// Arduino's single-TU architecture requires file-scope objects
// ---------------------------------------------------------------------------
HardwareSerial modbusSerial(1);
ModbusMaster node;
uint32_t lastPoll = 0;

WiFiClientSecure espClient;
PubSubClient client(espClient);

// ---------------------------------------------------------------------------
// Wi-Fi
// ---------------------------------------------------------------------------
void setupWiFi() {
    Serial.println("\n--- Initializing Wi-Fi ---");
    Serial.print("Connecting to: ");
    Serial.println(WIFI_SSID);

    WiFi.disconnect(true);
    delay(WIFI_INIT_DELAY_MS);

    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    while (WiFi.status() != WL_CONNECTED) {
        delay(WIFI_RECONNECT_DELAY_MS);
        Serial.print(".");
    }

    Serial.println("\n[SUCCESS] Wi-Fi Connected!");
    Serial.print("IP Address: ");
    Serial.println(WiFi.localIP());
}

// ---------------------------------------------------------------------------
// MQTT
// ---------------------------------------------------------------------------
void connectMQTT() {
    while (!client.connected() && WiFi.status() == WL_CONNECTED) {
        Serial.print("Attempting Secure MQTT connection... ");

        String clientId = "ESP32-Gateway-" + String(random(0xffff), HEX);

        if (client.connect(clientId.c_str(), MQTT_USER, MQTT_PASS)) {
            Serial.println("[SUCCESS] Connected to Secure Broker!");
        } else {
            Serial.print("[FAILED] rc=");
            Serial.print(client.state());
            Serial.println(" -> Trying again in 5 seconds");
            delay(MQTT_RETRY_DELAY_MS);
        }
    }
}

// ---------------------------------------------------------------------------
// Modbus polling & MQTT publish
// ---------------------------------------------------------------------------
void pollSensor(uint8_t idx) {
    node.begin(SLAVE_ADDRS[idx], modbusSerial);
    uint8_t result = node.readHoldingRegisters(0x0000, 2);

    if (result == node.ku8MBSuccess) {
        float hum = static_cast<float>(node.getResponseBuffer(0)) / SENSOR_SCALE;
        float temp =
            static_cast<float>(static_cast<int16_t>(node.getResponseBuffer(1))) / SENSOR_SCALE;

        Serial.printf("Sensor%d  %.1f %%RH  %.1f C\n", SLAVE_ADDRS[idx], hum, temp);

        if (client.connected()) {
            std::array<char, TOPIC_BUF_SIZE> topic{};
            std::array<char, PAYLOAD_BUF_SIZE> payload{};

            snprintf(topic.data(), topic.size(), MQTT_TOPIC_TEMPLATE, SLAVE_ADDRS[idx]);

            snprintf(payload.data(), payload.size(), R"({"temperature":%.1f, "humidity":%.1f})",
                     temp, hum);

            const bool success = client.publish(topic.data(), payload.data());

            if (!success) {
                Serial.println("  --> MQTT Publish FAILED");
                return;
            }
            Serial.printf("  --> MQTT Published to [%s]: %s\n", topic.data(), payload.data());
        }
    } else {
        Serial.printf("Sensor%d ERROR 0x%02X\n", SLAVE_ADDRS[idx], result);
    }
}

// ---------------------------------------------------------------------------
// Arduino entry points
// ---------------------------------------------------------------------------
void setup() {
    Serial.begin(1000000);
    delay(SERIAL_INIT_DELAY_MS);

    modbusSerial.begin(4800, SERIAL_8N1, XY485_RX, XY485_TX);

    setupWiFi();
    espClient.setCACert(ROOT_CA);
    client.setServer(MQTT_SERVER, MQTT_PORT);
    connectMQTT();

    lastPoll = millis();
    Serial.println("\n=== Gateway Initialized and Running ===");
}

void loop() {
    // 1. Maintain Wi-Fi health
    if (WiFi.status() != WL_CONNECTED) {
        Serial.println("[WARNING] Wi-Fi lost! Halting system and reconnecting...");
        setupWiFi();
    }

    // 2. Maintain MQTT health
    if (!client.connected()) {
        connectMQTT();
    }
    client.loop();

    // 3. Modbus polling cadence
    if (millis() - lastPoll < POLL_INTERVAL_MS) {
        return;
    }
    lastPoll = millis();

    for (uint8_t i = 0; i < NUM_SENSORS; i++) {
        pollSensor(i);
        if (i < NUM_SENSORS - 1U) {
            delay(INTER_SLAVE_MS);
        }
    }
}