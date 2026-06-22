#include "MQTTManager.h"
#include <Arduino.h>
#include <WiFi.h>
#include "Logger.h"

namespace MQTTManager {

namespace {
const uint32_t kRetryDelayMs = 5000U;
const size_t   kClientIdSize = 24U;
uint32_t           lastAttemptMs = 0U;
const char* TAG = "MQTT";


// Single connection attempt — shared by connect() and maybeReconnect()
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

void setup(WiFiClientSecure& espClient, PubSubClient& client,
           const char* server, uint16_t port, const char* caCert) {
    espClient.setCACert(caCert);
    client.setServer(server, port);
}

void connect(PubSubClient& client, const char* user, const char* pass) {
    while (!client.connected() && WiFi.status() == WL_CONNECTED) {
        if (attempt(client, user, pass)) { return; }
        delay(kRetryDelayMs);
    }
}

void maybeReconnect(PubSubClient& client, const char* user, const char* pass) {
    if (client.connected() || WiFi.status() != WL_CONNECTED) { return; }

    const uint32_t now = millis();
    if (now - lastAttemptMs < kRetryDelayMs) { return; }
    lastAttemptMs = now;

    attempt(client, user, pass);
}

}  // namespace MQTTManager