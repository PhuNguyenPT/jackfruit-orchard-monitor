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

auto toPercent(uint16_t raw) -> float;
}  // namespace SoilPoller
#endif