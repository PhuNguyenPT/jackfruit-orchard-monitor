#ifndef SHT40_H
#define SHT40_H
#include <array>
#include <cstdint>

inline constexpr uint8_t NUM_SENSORS = 2U;
inline constexpr std::array<uint8_t, NUM_SENSORS> SLAVE_ADDRS = {1U, 2U};

#endif