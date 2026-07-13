#include "TimeSync.h"

#include <Arduino.h>

#include <atomic>
#include <ctime>

namespace TimeSync {

namespace {
const uint32_t kSyncTimeoutMs = 30000U;  // give up after 30000 ms = 30s
const uint32_t kPollDelayMs = 500U;
// Floor chosen as "before this project existed" — any clock reading
// earlier than this is the ESP32's un-synced power-on default, not a
// real NTP result. Exact date has no other significance; bump it
// forward over time if you want, it just needs to predate "now."
const time_t kMinPlausibleTs = 1735689600L;  // 2025-01-01T00:00:00Z
std::atomic<bool> synced{false};
const char* TAG = "TimeSync";
}  // namespace

void setup() {
    ESP_LOGI(TAG, "Querying NTP pools for network time sync...");

    configTime(0, 0, "pool.ntp.org", "time.nist.gov");

    time_t currentTime = time(nullptr);
    const uint32_t startMs = millis();

    while (currentTime < kMinPlausibleTs) {
        if (millis() - startMs >= kSyncTimeoutMs) {
            ESP_LOGW(TAG, "NTP sync timed out after %lu ms — proceeding without synced time.",
                     static_cast<unsigned long>(kSyncTimeoutMs));
            synced.store(false, std::memory_order_release);
            return;
        }
        delay(kPollDelayMs);
        currentTime = time(nullptr);
    }
    synced.store(true, std::memory_order_release);
    ESP_LOGI(TAG, "NTP Time synchronized perfectly.");
}

auto isSynced() -> bool { return synced.load(std::memory_order_acquire); }

auto now() -> time_t { return time(nullptr); }
}  // namespace TimeSync