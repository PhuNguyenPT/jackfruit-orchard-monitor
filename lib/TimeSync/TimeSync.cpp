#include "TimeSync.h"
#include <Arduino.h>
#include <ctime>
#include "Logger.h"

namespace TimeSync {

void setup() {
    Logger::log(Logger::Level::INFO, "Querying NTP pools for network time sync...");
    configTime(7 * 3600, 0, "pool.ntp.org", "time.nist.gov");

    time_t now = time(nullptr);
    while (now < 1600000000L) {
        delay(500);
        now = time(nullptr);
    }

    struct tm timeinfo{};
    if (gmtime_r(&now, &timeinfo) == nullptr) {
        Logger::log(Logger::Level::WARN, "gmtime_r returned null, timestamp may be inaccurate");
    }
    Logger::log(Logger::Level::SUCCESS, "NTP Time synchronized perfectly.");
}

}  // namespace TimeSync