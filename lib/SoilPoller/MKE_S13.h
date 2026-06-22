#ifndef MKE_S13_H
#define MKE_S13_H
#include <array>
#include <cstdint>

namespace SoilConfig {

inline constexpr uint8_t kMuxS0 = 25U;
inline constexpr uint8_t kMuxS1 = 26U;
inline constexpr uint8_t kMuxS2 = 27U;
inline constexpr uint8_t kMuxS3 = 14U;

struct MuxBoard {
    uint8_t sigPin;
    uint8_t enPin;
    uint8_t numCh;
};

inline constexpr std::array<MuxBoard, 2> kBoards = {{
    {32U, 18U, 2U},  // MUX1: SIG=GPIO32, EN=GPIO18, 2 sensors on CH0-CH1
    {33U, 19U, 2U},  // MUX2: SIG=GPIO33, EN=GPIO19, 2 sensors on CH0-CH1
}};

inline constexpr uint8_t kNumBoards = static_cast<uint8_t>(kBoards.size());

inline constexpr uint8_t kNumSensors = []() constexpr {
    uint8_t n = 0U;
    for (const auto& board : kBoards) { n += board.numCh; }
    return n;
}();

}  // namespace SoilConfig
#endif