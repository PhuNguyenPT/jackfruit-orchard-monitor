#ifndef SHT40_H
#define SHT40_H
#define NUM_SENSORS 2
static const uint8_t SLAVE_ADDRS[NUM_SENSORS] = {1, 2};

// --- MQTT Topic Definitions ---
#define MQTT_TOPIC_TEMPLATE "sht40/sensor%d/data"

#endif