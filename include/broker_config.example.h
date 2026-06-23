// ==========================================
// MQTT CONFIGURATION
// ==========================================
// Copy this file to broker_config.h and fill in your values.
// Build with -DMQTT_SECURE (esp32prod env) for production TLS.
// Default (esp32dev env) uses setInsecure() + plain TCP on port 1883.
// ==========================================
#ifndef BROKER_CONFIG_H
#define BROKER_CONFIG_H

#ifdef MQTT_SECURE
// ---------------------------------------------------------------------------
// Production — TLS verified, port 8883
// ---------------------------------------------------------------------------
inline constexpr char MQTT_SERVER[] = "mqtt.your-domain.com";
inline constexpr int MQTT_PORT = 8883;
inline constexpr char MQTT_USER[] = "your-mqtt-user";
inline constexpr char MQTT_PASS[] = "your-mqtt-password";

inline constexpr char ROOT_CA[] = R"EOF(
-----BEGIN CERTIFICATE-----
<paste your CA certificate here>
-----END CERTIFICATE-----
)EOF";

#else
// ---------------------------------------------------------------------------
// Dev / test — plain TCP, no CA needed, port 1883
// ---------------------------------------------------------------------------
inline constexpr char MQTT_SERVER[] = "192.168.x.x";
inline constexpr int MQTT_PORT = 1883;
inline constexpr char MQTT_USER[] = "your-mqtt-user";
inline constexpr char MQTT_PASS[] = "your-mqtt-password";

inline constexpr char ROOT_CA[] = "";  // unused — setInsecure() active

#endif

#endif  // BROKER_CONFIG_H