#ifndef MKE_S13_H
#define MKE_S13_H
#include <cstdint>

namespace SoilConfig {

// ---------------------------------------------------------------------------
// MUX control pins (shared across both 74HC4067 boards)
// ---------------------------------------------------------------------------
static const uint8_t kMuxS0 = 25U;
static const uint8_t kMuxS1 = 26U;
static const uint8_t kMuxS2 = 27U;
static const uint8_t kMuxS3 = 14U;

// ---------------------------------------------------------------------------
// Per-board SIG and EN pins
// ---------------------------------------------------------------------------
struct MuxBoard {
    uint8_t sigPin;  // ADC1 pin for this board's SIG line
    uint8_t enPin;   // EN pin (active LOW)
    uint8_t numCh;   // number of populated sensor channels
};

static const MuxBoard kBoards[] = {
    {32U, 18U, 2U},  // MUX1: SIG=GPIO32, EN=GPIO18, 2 sensors on CH0-CH1
    {33U, 19U, 2U},  // MUX2: SIG=GPIO33, EN=GPIO19, 2 sensors on CH0-CH1
};

static const uint8_t kNumBoards =
    static_cast<uint8_t>(sizeof(kBoards) / sizeof(kBoards[0]));

// Total sensor count derived from board descriptors
static const uint8_t kNumSensors = []() {
    uint8_t n = 0U;
    for (uint8_t i = 0U; i < kNumBoards; i++) n += kBoards[i].numCh;
    return n;
}();

}  // namespace SoilConfig
#endif