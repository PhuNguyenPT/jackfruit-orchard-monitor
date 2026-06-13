#include <Arduino.h>
#include <HardwareSerial.h>
#include <ModbusMaster.h>
#include <PubSubClient.h>
#include <WiFi.h>
#include <WiFiClientSecure.h>

#include <array>
#include <cstdint>
#include <ctime>

#include "broker_config.h"
#include "config.h"
#include "gpio.h"
#include "sht40.h"
#include "wifi.h"

// ---------------------------------------------------------------------------
// Compile-time constants
// ---------------------------------------------------------------------------
namespace {

constexpr uint32_t WIFI_INIT_DELAY_MS = 100U;
constexpr uint32_t WIFI_RECONNECT_DELAY_MS = 500U;
constexpr uint32_t MQTT_RETRY_DELAY_MS = 5000U;
constexpr uint32_t SERIAL_INIT_DELAY_MS = 10U;

constexpr float SENSOR_SCALE = 10.0F;

constexpr size_t TOPIC_BUF_SIZE = 50U;
constexpr size_t PAYLOAD_BUF_SIZE = 100U;
constexpr size_t LOG_BUF_SIZE = 128U;
constexpr size_t TIME_BUF_SIZE = 32U;

}  // namespace

// ---------------------------------------------------------------------------
// Logger
// ---------------------------------------------------------------------------
namespace Logger {

enum class Level : std::uint8_t { INFO, SUCCESS, WARN, ERROR };

static constexpr std::array<const char*, 4> kLevelTags = {"[INFO] ", "[SUCC] ", "[WARN] ",
                                                          "[ERRO] "};

template <typename... Args>
void log(Level level, const char* format, Args... args) {
    std::array<char, TIME_BUF_SIZE> timebuf{};
    struct tm timeinfo {};

    if (getLocalTime(&timeinfo)) {
        strftime(timebuf.data(), timebuf.size(), "[%Y-%m-%d %H:%M:%S] ", &timeinfo);
    } else {
        // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
        snprintf(timebuf.data(), timebuf.size(), "[%10lu s] ", millis() / 1000UL);
    }
    Serial.print(timebuf.data());
    Serial.print(kLevelTags.at(static_cast<uint8_t>(level)));

    std::array<char, LOG_BUF_SIZE> buffer{};
    // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
    snprintf(buffer.data(), buffer.size(), format, args...);
    Serial.println(buffer.data());
}

}  // namespace Logger

// ---------------------------------------------------------------------------
// Global objects
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
    Logger::log(Logger::Level::INFO, "Initializing Wi-Fi interface...");
    Logger::log(Logger::Level::INFO, "Connecting to SSID: %s", WIFI_SSID);

    WiFi.disconnect(true);
    delay(WIFI_INIT_DELAY_MS);

    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    while (WiFi.status() != WL_CONNECTED) {
        delay(WIFI_RECONNECT_DELAY_MS);
        Serial.print(".");  // Kept raw for visual connection tracking
    }
    Serial.println();

    Logger::log(Logger::Level::SUCCESS, "Wi-Fi Connected! IP Assigned: %s",
                WiFi.localIP().toString().c_str());
}

// ---------------------------------------------------------------------------
// Time Synchronization
// ---------------------------------------------------------------------------
void setupTime() {
    Logger::log(Logger::Level::INFO, "Querying NTP pools for network time sync...");
    configTime(7 * 3600, 0, "pool.ntp.org", "time.nist.gov");

    time_t now = time(nullptr);
    while (now < 1600000000) {
        delay(500);
        now = time(nullptr);
    }

    struct tm timeinfo {};
    if (gmtime_r(&now, &timeinfo) == nullptr) {
        Logger::log(Logger::Level::WARN, "gmtime_r returned null, timestamp may be inaccurate");
    }
    Logger::log(Logger::Level::SUCCESS, "NTP Time synchronized perfectly.");
}

// ---------------------------------------------------------------------------
// MQTT
// ---------------------------------------------------------------------------
void connectMQTT() {
    while (!client.connected() && WiFi.status() == WL_CONNECTED) {
        Logger::log(Logger::Level::INFO, "Attempting TLS encrypted MQTT handshake...");

        String clientId = "ESP32-Gateway-" + String(random(0xffff), HEX);

        if (client.connect(clientId.c_str(), MQTT_USER, MQTT_PASS)) {
            Logger::log(Logger::Level::SUCCESS, "TLS Session established. Connected to broker.");
        } else {
            Logger::log(Logger::Level::ERROR, "MQTT connection failure, rc=%d. Retrying in 5s...",
                        client.state());
            delay(MQTT_RETRY_DELAY_MS);
        }
    }
}

// ---------------------------------------------------------------------------
// Modbus polling & MQTT publish
// ---------------------------------------------------------------------------
void pollSensor(uint8_t idx) {
    const uint8_t addr = SLAVE_ADDRS.at(idx);
    node.begin(addr, modbusSerial);
    uint8_t result = node.readHoldingRegisters(0x0000, 2);

    if (result == node.ku8MBSuccess) {
        float hum = static_cast<float>(node.getResponseBuffer(0)) / SENSOR_SCALE;
        float temp =
            static_cast<float>(static_cast<int16_t>(node.getResponseBuffer(1))) / SENSOR_SCALE;

        Logger::log(Logger::Level::SUCCESS, "Sensor %d Readout: %.1f %%RH | %.1f C", addr, hum,
                    temp);

        if (client.connected()) {
            std::array<char, TOPIC_BUF_SIZE> topic{};
            std::array<char, PAYLOAD_BUF_SIZE> payload{};

            // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
            snprintf(topic.data(), topic.size(), MQTT_TOPIC_TEMPLATE, addr);
            // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
            snprintf(payload.data(), payload.size(), R"({"temperature":%.1f, "humidity":%.1f})",
                     temp, hum);

            if (!client.publish(topic.data(), payload.data())) {
                Logger::log(Logger::Level::ERROR, "MQTT Frame dropped. Publish failed.");
                return;
            }
            Logger::log(Logger::Level::INFO, "MQTT Outbound -> [%s] Payload: %s", topic.data(),
                        payload.data());
        }
    } else {
        Logger::log(Logger::Level::ERROR, "Modbus fault on Sensor %d. Exception Code: 0x%02X", addr,
                    result);
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
    setupTime();

    espClient.setCACert(ROOT_CA);
    client.setServer(MQTT_SERVER, MQTT_PORT);
    connectMQTT();

    lastPoll = millis();
    Logger::log(Logger::Level::INFO, "System Pipeline Initialized. Commencing telemetry loops.");
}

void loop() {
    if (WiFi.status() != WL_CONNECTED) {
        Logger::log(Logger::Level::WARN, "Link state dropped! Re-asserting Wi-Fi stack...");
        setupWiFi();
    }

    if (!client.connected()) {
        connectMQTT();
    }
    client.loop();

    if (millis() - lastPoll < POLL_INTERVAL_MS) {
        return;
    }
    lastPoll = millis();

    Logger::log(Logger::Level::INFO, "Executing scheduled Modbus scan...");
    for (uint8_t i = 0; i < NUM_SENSORS; i++) {
        pollSensor(i);
        if (i < NUM_SENSORS - 1U) {
            delay(INTER_SLAVE_MS);
        }
    }
}