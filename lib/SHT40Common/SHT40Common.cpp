#include "SHT40Common.h"

namespace SHT40Poller {

auto scaleHumidity(uint16_t raw) -> float { return static_cast<float>(raw) / kSensorScale; }

auto scaleTemperature(uint16_t raw) -> float {
    return static_cast<float>(static_cast<int16_t>(raw)) / kSensorScale;
}

}  // namespace SHT40Poller