#ifndef SHT40_POLLER_H
#define SHT40_POLLER_H

#include <cstdint>

#include "SHT40Common.h"

namespace SHT40Poller {
enum class RxPin : int {};
enum class TxPin : int {};

void init(RxPin rxPin, TxPin txPin, uint32_t baud = 4800U);
void poll(uint8_t slaveAddr);

}  // namespace SHT40Poller
#endif