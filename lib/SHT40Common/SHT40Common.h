#ifndef SHT40_COMMON_H
#define SHT40_COMMON_H
#include <cstddef>
#include <cstdint>

namespace SHT40Poller {

static const float kSensorScale = 10.0F;
static const size_t kTopicBufSize = 50U;
static const size_t kPayloadBufSize = 100U;
static const char kTopicTemplate[] = "sht40/%d/data";

auto scaleHumidity(uint16_t raw) -> float;
auto scaleTemperature(uint16_t raw) -> float;

}  // namespace SHT40Poller
#endif