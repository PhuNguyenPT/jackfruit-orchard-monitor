#ifndef SHT40_H
#define SHT40_H
static const uint8_t NUM_SENSORS = 2;
constexpr std::array<uint8_t, NUM_SENSORS> SLAVE_ADDRS = {1, 2};

// --- MQTT Topic Definitions ---
static const char MQTT_TOPIC_TEMPLATE[] = "sht40/%d/data";

#endif