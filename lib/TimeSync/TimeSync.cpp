#include "TimeSync.h"
#include <Arduino.h>
#include <ctime>
#include "Logger.h"

namespace TimeSync {

namespace {
const uint32_t kSyncTimeoutMs   = 30000U;  // give up after 30000 ms = 30s
const uint32_t kPollDelayMs     = 500U;
// Floor chosen as "before this project existed" — any clock reading
// earlier than this is the ESP32's un-synced power-on default, not a
// real NTP result. Exact date has no other significance; bump it
// forward over time if you want, it just needs to predate "now."
const time_t kMinPlausibleTs = 1735689600L;  // 2025-01-01T00:00:00Z
const uint32_t kResyncIntervalMs = 6UL * 60UL * 60UL * 1000UL;  // 6 hours
uint32_t       lastResyncMs      = 0U;
bool           synced           = false;
}  // namespace

void setup() {
    Logger::log(Logger::Level::INFO, "Querying NTP pools for network time sync...");
    configTime(0, 0, "pool.ntp.org", "time.nist.gov");  // UTC, no offset — see note below

    time_t now = time(nullptr);
    const uint32_t startMs = millis();

    while (now < kMinPlausibleTs) {
        if (millis() - startMs >= kSyncTimeoutMs) {
            Logger::log(Logger::Level::WARN,
                        "NTP sync timed out after %lu ms — proceeding without synced time.",
                        static_cast<unsigned long>(kSyncTimeoutMs));
            synced = false;
            return;
        }
        delay(kPollDelayMs);
        now = time(nullptr);
    }

    synced = true;
    Logger::log(Logger::Level::SUCCESS, "NTP Time synchronized perfectly.");
}

void maybeResync() {
    if (millis() - lastResyncMs < kResyncIntervalMs) return;
    lastResyncMs = millis();

    Logger::log(Logger::Level::INFO, "Performing periodic NTP re-sync...");
    configTime(0, 0, "pool.ntp.org", "time.nist.gov");

    // Non-blocking check — don't stall loop() if NTP is briefly unreachable;
    // just try again next interval.
    time_t now = time(nullptr);
    if (now >= kMinPlausibleTs) {
        synced = true;
        Logger::log(Logger::Level::SUCCESS, "NTP re-sync successful.");
    } else {
        Logger::log(Logger::Level::WARN, "NTP re-sync did not return valid time yet.");
    }
}

bool isSynced() {
    return synced;
}

time_t now() {
    return time(nullptr);
}

}  // namespace TimeSync