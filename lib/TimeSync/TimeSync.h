#ifndef TIME_SYNC_H
#define TIME_SYNC_H
#pragma once
#include <ctime>

namespace TimeSync {
    void setup();
    void maybeResync();
    auto isSynced() -> bool;
    auto now() -> time_t;
}
#endif