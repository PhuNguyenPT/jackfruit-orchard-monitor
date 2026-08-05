#include <Arduino.h>
#include <WiFi.h>
#include <driver/gpio.h>
#include <esp_sleep.h>
#include <freertos/FreeRTOS.h>
#include <freertos/event_groups.h>
#include <freertos/task.h>
#include <mqtt_client.h>

#include "MQTTManager.h"
#include "PowerRail.h"
#include "SHT40Poller.h"
#include "SoilPoller.h"
#include "TimeSync.h"
#include "broker_config.h"
#include "config.h"
#include "gpio.h"
#include "sht40.h"
#include "wifi.h"

namespace {
const uint32_t kWifiInitDelayMs = 100U;
const uint32_t kWifiReconnectDelayMs = 500U;
const uint32_t kSerialInitDelayMs = 10U;
const uint32_t kMqttConnectTimeoutMs = 15000U;
const uint32_t kMqttAckTimeoutMs = 2000U;
#ifndef DEEP_SLEEP_SEC
#define DEEP_SLEEP_SEC 3600ULL
#endif

const uint64_t kDeepSleepUs = static_cast<uint64_t>(DEEP_SLEEP_SEC) * 1000000ULL;
const char* TAG = "Main";

RTC_DATA_ATTR uint32_t bootCount = 0U;

// --- sync primitives ---------------------------------------------------
EventGroupHandle_t g_syncEvents = nullptr;

constexpr EventBits_t kNetworkReadyBit = BIT0;
constexpr EventBits_t kHardwareReadyBit = BIT1;
constexpr EventBits_t kSht40DoneBit = BIT2;
constexpr EventBits_t kSoilDoneBit = BIT3;

constexpr uint32_t kNetTaskStack = 12288U;  // TLS handshake is stack-hungry
constexpr uint32_t kHwInitStack = 3072U;
constexpr uint32_t kSht40PollStack = 8192U;
constexpr uint32_t kSoilPollStack = 8192U;
}  // namespace

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
// Phase A: network chain (core 0)
// ---------------------------------------------------------------------------
void networkSetupTask(void* /*pvParameters*/) {
    setupWiFi();
    TimeSync::setup();

    esp_mqtt_client_config_t mqttCfg = {};
    mqttCfg.host = MQTT_SERVER;
    mqttCfg.port = static_cast<uint16_t>(MQTT_PORT);
    mqttCfg.username = MQTT_USER;
    mqttCfg.password = MQTT_PASS;
#ifdef MQTT_SECURE
    mqttCfg.transport = MQTT_TRANSPORT_OVER_SSL;
    mqttCfg.cert_pem = ROOT_CA;
#else
    mqttCfg.transport = MQTT_TRANSPORT_OVER_TCP;
#endif
    MQTTManager::setup(mqttCfg);
    MQTTManager::connect(kMqttConnectTimeoutMs);  // bounded — don't block forever pre-sleep

    xEventGroupSetBits(g_syncEvents, kNetworkReadyBit);
    vTaskDelete(nullptr);
}

// ---------------------------------------------------------------------------
// Phase B: hardware init (core 1) — zero overlap with WiFi radio
// ---------------------------------------------------------------------------
void hardwareInitTask(void* /*pvParameters*/) {
    PowerRail::init();
    SHT40Poller::init(SHT40Poller::RxPin{XY485_RX}, SHT40Poller::TxPin{XY485_TX});
    SoilPoller::init();

    xEventGroupSetBits(g_syncEvents, kHardwareReadyBit);
    vTaskDelete(nullptr);
}

// ---------------------------------------------------------------------------
// Phase C: SHT40 poll (core 0)
// ---------------------------------------------------------------------------
void sht40PollTask(void* /*pvParameters*/) {
    ESP_LOGI(TAG, "Executing scheduled Modbus scan...");
    for (uint8_t i = 0; i < NUM_SENSORS; i++) {
        SHT40Poller::poll(SLAVE_ADDRS.at(i));
        if (i < NUM_SENSORS - 1U) {
            delay(INTER_SLAVE_MS);
        }
    }
    xEventGroupSetBits(g_syncEvents, kSht40DoneBit);
    vTaskDelete(nullptr);
}

// ---------------------------------------------------------------------------
// Phase D: Soil poll (core 1)
// ---------------------------------------------------------------------------
void soilPollTask(void* /*pvParameters*/) {
    ESP_LOGI(TAG, "Executing scheduled soil moisture scan...");
    SoilPoller::poll();
    xEventGroupSetBits(g_syncEvents, kSoilDoneBit);
    vTaskDelete(nullptr);
}

// ---------------------------------------------------------------------------
void goToSleep() {
    ESP_LOGI(TAG,
             "Cycle complete. Shutting down radios and entering deep sleep for %llu seconds...",
             static_cast<unsigned long long>(DEEP_SLEEP_SEC));

    SoilPoller::parkForSleep();
    PowerRail::parkForSleep();
    gpio_deep_sleep_hold_en();  // global switch — required for per-pin holds to survive deep sleep

    MQTTManager::waitForAcks(kMqttAckTimeoutMs);  // let PUBACKs land before killing the radio
    MQTTManager::disconnect();
    WiFi.disconnect(true);
    WiFi.mode(WIFI_OFF);

    esp_sleep_enable_timer_wakeup(kDeepSleepUs);
    esp_deep_sleep_start();
}

void setup() {
    Serial.begin(1000000);
    delay(kSerialInitDelayMs);

    bootCount++;
    ESP_LOGI(TAG, "Boot #%lu (wake cause: %d)", static_cast<unsigned long>(bootCount),
             static_cast<int>(esp_sleep_get_wakeup_cause()));

    g_syncEvents = xEventGroupCreate();

    // --- Phase A + B in parallel ---
    xTaskCreatePinnedToCore(networkSetupTask, "NetSetup", kNetTaskStack, nullptr, 1, nullptr, 1);
    xTaskCreatePinnedToCore(hardwareInitTask, "HwInit", kHwInitStack, nullptr, 1, nullptr, 0);

    xEventGroupWaitBits(g_syncEvents, kNetworkReadyBit | kHardwareReadyBit,
                        pdTRUE /*clear on exit*/, pdTRUE /*wait for ALL*/, portMAX_DELAY);

    ESP_LOGI(TAG, "System Pipeline Initialized. Running telemetry cycle.");

    // --- Phase C + D in parallel ---
    xTaskCreatePinnedToCore(sht40PollTask, "Sht40Poll", kSht40PollStack, nullptr, 1, nullptr, 0);
    xTaskCreatePinnedToCore(soilPollTask, "SoilPoll", kSoilPollStack, nullptr, 1, nullptr, 1);

    xEventGroupWaitBits(g_syncEvents, kSht40DoneBit | kSoilDoneBit, pdTRUE, pdTRUE, portMAX_DELAY);

    vEventGroupDelete(g_syncEvents);

    goToSleep();
}

void loop() {
    // Unreachable — esp_deep_sleep_start() resets the chip on wake.
}