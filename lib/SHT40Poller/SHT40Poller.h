#ifndef SHT40_POLLER_H
#define SHT40_POLLER_H

#include <cstdint>

#include "SHT40Common.h"

namespace SHT40Poller {
enum class RxPin : int {};
enum class TxPin : int {};

// Per datasheet: SHT sensor's internal temperature/humidity register
// refresh cycle is 2s. Registers read all-zero until the first cycle
// completes post power-on. poll() enforces this internally; exposed here
// so hardware tests can validate real settle time against the same
// constant production uses, rather than duplicating the number.
inline constexpr uint32_t kMinWarmupMs = 2500U;

void init(RxPin rxPin, TxPin txPin, uint32_t baud = 4800U);
void poll(uint8_t slaveAddr);

}  // namespace SHT40Poller
#endif