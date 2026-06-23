#ifndef SHT40_COMMON_H
#define SHT40_COMMON_H
#include <cstddef>
#include <cstdint>

namespace SHT40Poller {

inline constexpr float kSensorScale = 10.0F;
inline constexpr size_t kTopicBufSize = 50U;
inline constexpr size_t kPayloadBufSize = 100U;
inline constexpr char kTopicTemplate[] = "sht40/%d/data";

auto scaleHumidity(uint16_t raw) -> float;
auto scaleTemperature(uint16_t raw) -> float;

}  // namespace SHT40Poller
#endif