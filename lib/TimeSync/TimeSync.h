#ifndef TIME_SYNC_H
#define TIME_SYNC_H
#pragma once
#include <ctime>

namespace TimeSync {
    void setup();
    void maybeResync();
    bool isSynced();
    time_t now();
}
#endif