#include <Arduino.h>

#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include <ModbusMaster.h>
#include <HardwareSerial.h>
#include "config.h"
#include "broker_config.h"
#include "gpio.h"
#include "sht40.h"

// --- Global Objects ---
HardwareSerial modbusSerial(1);
ModbusMaster   node;
uint32_t       lastPoll = 0;

WiFiClientSecure espClient;
PubSubClient client(espClient);

// --- Wi-Fi Connection Function ---
void setupWiFi() {
  Serial.println("\n--- Initializing Wi-Fi ---");
  Serial.print("Connecting to: ");
  Serial.println(WIFI_SSID);

  WiFi.disconnect(true);
  delay(100);

  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("\n[SUCCESS] Wi-Fi Connected!");
  Serial.print("IP Address: ");
  Serial.println(WiFi.localIP());
}

// --- MQTT Connection Function ---
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
      delay(5000);
    }
  }
}

// --- Modbus Polling & MQTT Publishing ---
void pollSensor(uint8_t idx) {
  node.begin(SLAVE_ADDRS[idx], modbusSerial);
  uint8_t result = node.readHoldingRegisters(0x0000, 2);

  if (result == node.ku8MBSuccess) {
    float hum  =          node.getResponseBuffer(0) / 10.0f;
    float temp = (int16_t)node.getResponseBuffer(1) / 10.0f;

    Serial.printf("Sensor%d  %.1f %%RH  %.1f C\n", SLAVE_ADDRS[idx], hum, temp);

    // If MQTT is connected, publish the data
    if (client.connected()) {
      char topic[50];
      char payload[100];

      // Construct the topic dynamically using your defined template
      snprintf(topic, sizeof(topic), MQTT_TOPIC_TEMPLATE, SLAVE_ADDRS[idx]);

      // Format data as JSON
      snprintf(payload, sizeof(payload), "{\"temperature\":%.1f, \"humidity\":%.1f}", temp, hum);

      if (client.publish(topic, payload)) {
        Serial.printf("  --> MQTT Published to [%s]: %s\n", topic, payload);
      } else {
        Serial.println("  --> MQTT Publish FAILED");
      }
    }
  } else {
    Serial.printf("Sensor%d ERROR 0x%02X\n", SLAVE_ADDRS[idx], result);
  }
}

void setup() {
  Serial.begin(1000000);
  delay(10);

  // Initialize Modbus Serial
  modbusSerial.begin(4800, SERIAL_8N1, XY485_RX, XY485_TX);

  // Establish Network and Cloud Layers
  setupWiFi();
  espClient.setCACert(ROOT_CA);
  client.setServer(MQTT_SERVER, MQTT_PORT);
  connectMQTT();

  lastPoll = millis();
  Serial.println("\n=== Gateway Initialized and Running ===");
}

void loop() {
  // 1. Maintain Wi-Fi Health
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("[WARNING] Wi-Fi lost! Halting system and reconnecting...");
    setupWiFi();
  }

  // 2. Maintain MQTT Health
  if (!client.connected()) {
    connectMQTT();
  }
  client.loop();

  // 3. Modbus Routine
  if (millis() - lastPoll < POLL_INTERVAL_MS) return;
  lastPoll = millis();

  for (int i = 0; i < NUM_SENSORS; i++) {
    pollSensor(i);
    if (i < NUM_SENSORS - 1) delay(INTER_SLAVE_MS);
  }
}