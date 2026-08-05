// PowerRail.h
#ifndef POWER_RAIL_H
#define POWER_RAIL_H

namespace PowerRail {
void init();          // release any prior hold, drive rail ON
void parkForSleep();  // drive rail OFF, hold through deep sleep
}  // namespace PowerRail
#endif