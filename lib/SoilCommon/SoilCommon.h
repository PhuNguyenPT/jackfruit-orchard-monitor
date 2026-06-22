#ifndef SOIL_COMMON_H
#define SOIL_COMMON_H
#include <cstddef>
#include <cstdint>

namespace SoilPoller {

static const size_t kTopicBufSize = 50U;
static const size_t kPayloadBufSize = 100U;
static const char kTopicTemplate[] = "mke-s13/%d/data";

// Calibration (MKE-S13 at 5V)
static const uint16_t kDryValue = 3500U;
static const uint16_t kWetValue = 1760U;

inline float toPercent(uint16_t raw) {
    if (raw >= kDryValue) return 0.0F;
    if (raw <= kWetValue) return 100.0F;
    return (1.0F -
            static_cast<float>(raw - kWetValue) / static_cast<float>(kDryValue - kWetValue)) *
           100.0F;
}

}  // namespace SoilPoller
#endif