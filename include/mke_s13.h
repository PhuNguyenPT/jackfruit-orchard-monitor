#ifndef MKE_S13_H
#define MKE_S13_H
#include <cstdint>

namespace SoilConfig {

constexpr uint8_t kSoilPins[] = {32, 33, 34, 35};
constexpr uint8_t kNumSensors = static_cast<uint8_t>(sizeof(kSoilPins) / sizeof(kSoilPins[0]));

}  // namespace SoilConfig
#endif