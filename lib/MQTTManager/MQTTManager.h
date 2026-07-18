#ifndef MQTT_MANAGER_H
#define MQTT_MANAGER_H
#include <mqtt_client.h>

#include <cstdint>

namespace MQTTManager {

// Caller fills in host/port/username/password (and transport/cert_pem if
// MQTT_SECURE is defined) before calling setup(). Everything else in cfg
// is left at whatever the caller set — MQTTManager doesn't own defaults
// beyond what it explicitly sets internally (keepalive).
void setup(const esp_mqtt_client_config_t& cfg);

// Blocking — use in setup() only. Returns false on timeout (proceed anyway;
// the client keeps retrying in the background).
auto connect(uint32_t timeoutMs) -> bool;

auto isConnected() -> bool;

// QoS1 publish. Returns the broker-assigned msg_id (>=0) once queued, or
// a negative value if it couldn't even be queued. A non-negative return
// does NOT mean the broker got it — call waitForAcks() before sleeping.
auto publish(const char* topic, const char* payload) -> int;

// Blocks until every QoS1 publish since setup() has been PUBACKed, or
// timeoutMs elapses. Returns false on timeout — caller should sleep anyway
// rather than hold the radio open indefinitely.
auto waitForAcks(uint32_t timeoutMs) -> bool;

void disconnect();

}  // namespace MQTTManager
#endif