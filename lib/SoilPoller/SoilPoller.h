#ifndef SOIL_POLLER_H
#define SOIL_POLLER_H

#include "SoilCommon.h"

namespace SoilPoller {

void init();
void poll();
void parkForSleep();
}  // namespace SoilPoller
#endif