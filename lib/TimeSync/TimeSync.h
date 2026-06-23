#ifndef TIME_SYNC_H
#define TIME_SYNC_H
#include <ctime>

namespace TimeSync {
void setup();
void maybeResync();
auto isSynced() -> bool;
auto now() -> time_t;
}  // namespace TimeSync
#endif