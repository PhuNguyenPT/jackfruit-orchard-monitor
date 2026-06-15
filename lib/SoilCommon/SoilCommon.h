#ifndef SOIL_COMMON_H
#define SOIL_COMMON_H
#include <cstddef>
#include <cstdint>

namespace SoilPoller {

constexpr size_t   kTopicBufSize   = 50U;
constexpr size_t   kPayloadBufSize = 100U;
constexpr char     kTopicTemplate[] = "mke-s13/%d/data";

// Calibration (MKE-S13 at 5V)
constexpr uint16_t kDryValue = 3340U;
constexpr uint16_t kWetValue = 1805U;

inline float toPercent(uint16_t raw) {
    if (raw >= kDryValue) return 0.0F;
    if (raw <= kWetValue) return 100.0F;
    return (1.0F - static_cast<float>(raw - kWetValue) /
                   static_cast<float>(kDryValue - kWetValue)) * 100.0F;
}

}  // namespace SoilPoller
#endif