#ifndef SOIL_POLLER_H
#define SOIL_POLLER_H
#include <PubSubClient.h>

#include "SoilCommon.h"

namespace SoilPoller {

void init();
void poll(PubSubClient& mqttClient);

}  // namespace SoilPoller
#endif