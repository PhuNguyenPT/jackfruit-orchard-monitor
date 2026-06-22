#include <Arduino.h>
#include <PubSubClient.h>
#include <WiFi.h>
#include <WiFiClientSecure.h>

#include "Logger.h"
#include "MQTTManager.h"
#include "SHT40Poller.h"
#include "SoilPoller.h"
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
const uint32_t kSoilPollIntervalMs = 10000U;
const char* TAG = "Main";
}  // namespace

// ---------------------------------------------------------------------------
// Global objects — main.cpp is the composition root
// ---------------------------------------------------------------------------
WiFiClientSecure espClient;
PubSubClient client(espClient);

static uint32_t lastSHT40 = 0U;
static uint32_t last_MKE_S13_Poll = 0U;

// ---------------------------------------------------------------------------
// Wi-Fi
// ---------------------------------------------------------------------------
void setupWiFi() {
    ESP_LOGI(TAG, "Initializing Wi-Fi interface...");
    ESP_LOGI(TAG, "Connecting to SSID: %s", WIFI_SSID);

    WiFi.disconnect(true);
    delay(kWifiInitDelayMs);
    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    while (WiFi.status() != WL_CONNECTED) {
        delay(kWifiReconnectDelayMs);
        Serial.print('.');
    }
    Serial.println();
    ESP_LOGI(TAG, "Wi-Fi Connected! IP Assigned: %s", WiFi.localIP().toString().c_str());
}

// ---------------------------------------------------------------------------
// Arduino entry points
// ---------------------------------------------------------------------------
void setup() {
    Serial.begin(1000000);
    delay(kSerialInitDelayMs);
    Logger::setup();

    SHT40Poller::init(SHT40Poller::RxPin{XY485_RX}, SHT40Poller::TxPin{XY485_TX});

    SoilPoller::init();

    setupWiFi();
    TimeSync::setup();
    MQTTManager::setup(espClient, client, MQTT_SERVER, MQTT_PORT, ROOT_CA);
    MQTTManager::connect(client, MQTT_USER, MQTT_PASS);  // blocking — fine at boot

    lastSHT40 = millis();
    ESP_LOGI(TAG, "System Pipeline Initialized. Commencing telemetry loops.");
}

void loop() {
    if (WiFi.status() != WL_CONNECTED) {
        ESP_LOGW(TAG, "Link state dropped! Re-asserting Wi-Fi stack...");
        setupWiFi();
    }

    TimeSync::maybeResync();

    MQTTManager::maybeReconnect(client, MQTT_USER, MQTT_PASS);
    client.loop();

    if (millis() - lastSHT40 >= POLL_INTERVAL_MS) {
        lastSHT40 = millis();
        ESP_LOGI(TAG, "Executing scheduled Modbus scan...");

        for (uint8_t i = 0; i < NUM_SENSORS; i++) {
            SHT40Poller::poll(SLAVE_ADDRS.at(i), client);
            if (i < NUM_SENSORS - 1U) {
                delay(INTER_SLAVE_MS);
            }
        }
    }

    if (millis() - last_MKE_S13_Poll >= kSoilPollIntervalMs) {
        last_MKE_S13_Poll = millis();
        ESP_LOGI(TAG, "Executing scheduled soil moisture scan...");

        SoilPoller::poll(client);
    }
}