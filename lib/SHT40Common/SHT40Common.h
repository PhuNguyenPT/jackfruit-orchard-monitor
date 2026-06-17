#ifndef SHT40_COMMON_H
#define SHT40_COMMON_H
#include <cstddef>
#include <cstdint>

namespace SHT40Poller {

static const float kSensorScale = 10.0F;
static const size_t kTopicBufSize   = 50U;
static const size_t kPayloadBufSize = 100U;
static const char kTopicTemplate[] = "sht40/%d/data";

inline float scaleHumidity(uint16_t raw) {
    return static_cast<float>(raw) / kSensorScale;
}

inline float scaleTemperature(uint16_t raw) {
    return static_cast<float>(static_cast<int16_t>(raw)) / kSensorScale;
}

}  // namespace SHT40Poller
#endif