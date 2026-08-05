// PowerRail.cpp
#include "PowerRail.h"

#include <Arduino.h>
#include <driver/gpio.h>

namespace PowerRail {

namespace {
const uint8_t kSensorPowerPin = 15U;
}  // namespace

void init() {
    gpio_hold_dis(static_cast<gpio_num_t>(kSensorPowerPin));
    pinMode(kSensorPowerPin, OUTPUT);
    digitalWrite(kSensorPowerPin, HIGH);  // both MKE-M06 modules ON -> 3.3V and 5V rails live
    delay(50U);  // let sensors/RS485 transceiver settle before init() calls
}

void parkForSleep() {
    digitalWrite(kSensorPowerPin, LOW);  // both modules OFF -> 3.3V and 5V rails cut
    gpio_hold_en(static_cast<gpio_num_t>(kSensorPowerPin));
}

}  // namespace PowerRail