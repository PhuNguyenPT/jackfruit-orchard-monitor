#ifndef MQTT_MANAGER_H
#define MQTT_MANAGER_H
#include <PubSubClient.h>
#include <WiFiClientSecure.h>

#include <cstdint>

namespace MQTTManager {

void setup(WiFiClientSecure& espClient, PubSubClient& client, const char* server, uint16_t port,
           const char* caCert);

// Blocking — use in setup() only; retries until connected or WiFi lost
void connect(PubSubClient& client, const char* user, const char* pass);

// Non-blocking — use in loop(); one attempt per kRetryDelayMs, returns immediately
void maybeReconnect(PubSubClient& client, const char* user, const char* pass);

// Thread-safe wrappers — safe to call concurrently from multiple tasks that
// share the same PubSubClient. Only the actual socket op is locked; callers
// do their sensor reads/formatting outside the lock.
auto publish(PubSubClient& client, const char* topic, const char* payload) -> bool;
auto isConnected(PubSubClient& client) -> bool;

}  // namespace MQTTManager
#endif