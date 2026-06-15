#ifndef SOIL_POLLER_H
#define SOIL_POLLER_H
#include <PubSubClient.h>
#include <cstdint>
#include "SoilCommon.h"

namespace SoilPoller {

void init(const uint8_t* pins, uint8_t count);
void poll(PubSubClient& mqttClient);

}  // namespace SoilPoller
#endif