#include "MQTTManager.h"

#include <Arduino.h>
#include <WiFi.h>
#include <freertos/FreeRTOS.h>
#include <freertos/semphr.h>

#include <array>

namespace MQTTManager {

namespace {
const uint32_t kRetryDelayMs = 5000U;
const size_t kClientIdSize = 24U;
uint32_t lastAttemptMs = 0U;
const char* TAG = "MQTT";

// Guards every PubSubClient socket op once multiple tasks share one client.
// Constructed at global-init time (before setup() runs), so there's no
// ordering dependency on which task touches publish()/isConnected() first.
SemaphoreHandle_t publishMutex = xSemaphoreCreateMutex();

auto attempt(PubSubClient& client, const char* user, const char* pass) -> bool {
    std::array<char, kClientIdSize> clientId{};
    // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
    snprintf(clientId.data(), clientId.size(), "ESP32-Gateway-%04X",
             static_cast<unsigned int>(random(0xffff)));

    ESP_LOGI(TAG, "Attempting TLS encrypted MQTT handshake...");
    if (client.connect(clientId.data(), user, pass)) {
        ESP_LOGI(TAG, "TLS Session established. Connected to broker.");
        return true;
    }
    ESP_LOGE(TAG, "MQTT connection failure, rc=%d.", client.state());
    return false;
}
}  // namespace

void setup(WiFiClientSecure& espClient, PubSubClient& client, const char* server, uint16_t port,
           const char* caCert) {
#ifdef MQTT_SECURE
    ESP_LOGI(TAG, "TLS: verifying broker against CA cert.");
    espClient.setCACert(caCert);
#else
    ESP_LOGW(TAG, "TLS verification DISABLED — dev/test mode.");
    espClient.setInsecure();
    (void)caCert;
#endif
    client.setServer(server, port);
}

void connect(PubSubClient& client, const char* user, const char* pass) {
    while (!client.connected() && WiFi.status() == WL_CONNECTED) {
        if (attempt(client, user, pass)) {
            return;
        }
        delay(kRetryDelayMs);
    }
}

void maybeReconnect(PubSubClient& client, const char* user, const char* pass) {
    if (client.connected() || WiFi.status() != WL_CONNECTED) {
        return;
    }

    const uint32_t now = millis();
    if (now - lastAttemptMs < kRetryDelayMs) {
        return;
    }
    lastAttemptMs = now;

    attempt(client, user, pass);
}

auto publish(PubSubClient& client, const char* topic, const char* payload) -> bool {
    xSemaphoreTake(publishMutex, portMAX_DELAY);
    const bool published = client.publish(topic, payload);
    xSemaphoreGive(publishMutex);
    return published;
}

auto isConnected(PubSubClient& client) -> bool {
    xSemaphoreTake(publishMutex, portMAX_DELAY);
    const bool connected = client.connected();
    xSemaphoreGive(publishMutex);
    return connected;
}

}  // namespace MQTTManager