#include <Arduino.h>
#include <PubSubClient.h>
#include <WiFi.h>
#include <WiFiClientSecure.h>

#include "Logger.h"
#include "MQTTManager.h"
#include "SHT40Poller.h"
#include "TimeSync.h"
#include "broker_config.h"
#include "config.h"
#include "gpio.h"
#include "sht40.h"
#include "wifi.h"

// ---------------------------------------------------------------------------
// Compile-time constants
// ---------------------------------------------------------------------------
namespace {
const uint32_t kWifiInitDelayMs = 100U;
const uint32_t kWifiReconnectDelayMs = 500U;
const uint32_t kSerialInitDelayMs = 10U;
}  // namespace

// ---------------------------------------------------------------------------
// Global objects — main.cpp is the composition root
// ---------------------------------------------------------------------------
WiFiClientSecure espClient;
PubSubClient client(espClient);

static uint32_t lastPoll = 0;

// ---------------------------------------------------------------------------
// Wi-Fi
// ---------------------------------------------------------------------------
void setupWiFi() {
    Logger::log(Logger::Level::INFO, "Initializing Wi-Fi interface...");
    Logger::log(Logger::Level::INFO, "Connecting to SSID: %s", WIFI_SSID);

    WiFi.disconnect(true);
    delay(kWifiInitDelayMs);
    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    while (WiFi.status() != WL_CONNECTED) {
        delay(kWifiReconnectDelayMs);
        Serial.print('.');
    }
    Serial.println();
    Logger::log(Logger::Level::SUCCESS, "Wi-Fi Connected! IP Assigned: %s",
                WiFi.localIP().toString().c_str());
}

// ---------------------------------------------------------------------------
// Arduino entry points
// ---------------------------------------------------------------------------
void setup() {
    Serial.begin(1000000);
    delay(kSerialInitDelayMs);

    SHT40Poller::init(XY485_RX, XY485_TX);
    setupWiFi();
    TimeSync::setup();
    MQTTManager::setup(espClient, client, MQTT_SERVER, MQTT_PORT, ROOT_CA);
    MQTTManager::connect(client, MQTT_USER, MQTT_PASS);  // blocking — fine at boot

    lastPoll = millis();
    Logger::log(Logger::Level::INFO, "System Pipeline Initialized. Commencing telemetry loops.");
}

void loop() {
    if (WiFi.status() != WL_CONNECTED) {
        Logger::log(Logger::Level::WARN, "Link state dropped! Re-asserting Wi-Fi stack...");
        setupWiFi();
    }

    MQTTManager::maybeReconnect(client, MQTT_USER, MQTT_PASS);
    client.loop();  // always runs, even during reconnect cooldown

    if (millis() - lastPoll < POLL_INTERVAL_MS) {
        return;
    }
    lastPoll = millis();

    Logger::log(Logger::Level::INFO, "Executing scheduled Modbus scan...");
    for (uint8_t i = 0; i < NUM_SENSORS; i++) {
        SHT40Poller::poll(SLAVE_ADDRS.at(i), client);
        if (i < NUM_SENSORS - 1U) {
            delay(INTER_SLAVE_MS);
        }
    }
}