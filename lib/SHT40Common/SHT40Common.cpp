#include "SHT40Common.h"

namespace SHT40Poller {

auto scaleHumidity(uint16_t raw) -> float { return static_cast<float>(raw) / kSensorScale; }

auto scaleTemperature(uint16_t raw) -> float {
    return static_cast<float>(static_cast<int16_t>(raw)) / kSensorScale;
}
auto isZeroedReading(uint16_t rawHumidity, uint16_t rawTemperature) -> bool {
    return rawHumidity == 0U && rawTemperature == 0U;
}
}  // namespace SHT40Poller