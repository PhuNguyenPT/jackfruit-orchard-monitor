#include "MQTTManager.h"
#include <Arduino.h>
#include <WiFi.h>
#include "Logger.h"

namespace MQTTManager {

namespace {
constexpr uint32_t kRetryDelayMs = 5000U;
constexpr size_t   kClientIdSize = 24U;
uint32_t           lastAttemptMs = 0U;
}

// Single connection attempt — shared by connect() and maybeReconnect()
bool attempt(PubSubClient& client, const char* user, const char* pass) {
    char clientId[kClientIdSize];
    // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
    snprintf(clientId, sizeof(clientId), "ESP32-Gateway-%04X",
             static_cast<unsigned int>(random(0xffff)));

    Logger::log(Logger::Level::INFO, "Attempting TLS encrypted MQTT handshake...");
    if (client.connect(clientId, user, pass)) {
        Logger::log(Logger::Level::SUCCESS, "TLS Session established. Connected to broker.");
        return true;
    }
    Logger::log(Logger::Level::ERROR,
                "MQTT connection failure, rc=%d.", client.state());
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
        if (attempt(client, user, pass)) return;
        delay(kRetryDelayMs);
    }
}

void maybeReconnect(PubSubClient& client, const char* user, const char* pass) {
    if (client.connected() || WiFi.status() != WL_CONNECTED) return;

    const uint32_t now = millis();
    if (now - lastAttemptMs < kRetryDelayMs) return;
    lastAttemptMs = now;

    attempt(client, user, pass);
}

}  // namespace MQTTManager