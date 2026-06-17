#ifndef LOGGER_H
#define LOGGER_H
#include <Arduino.h>
#include <array>
#include <ctime>

namespace Logger {

enum class Level : unsigned char { INFO, SUCCESS, WARN, ERROR };

namespace {
static const size_t kTimeBufSize = 32U;
static const size_t kLogBufSize  = 128U;
static const std::array<const char*, 4> kLevelTags = {
    "[INFO] ", "[SUCC] ", "[WARN] ", "[ERRO] "
};
}  // namespace

template <typename... Args>
void log(Level level, const char* format, Args... args) {
    std::array<char, kTimeBufSize> timebuf{};
    struct tm timeinfo{};

    if (getLocalTime(&timeinfo)) {
        strftime(timebuf.data(), timebuf.size(), "[%Y-%m-%d %H:%M:%S] ", &timeinfo);
    } else {
        // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
        snprintf(timebuf.data(), timebuf.size(), "[%10lu s] ", millis() / 1000UL);
    }
    Serial.print(timebuf.data());
    Serial.print(kLevelTags.at(static_cast<size_t>(level)));

    std::array<char, kLogBufSize> buffer{};
    // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
    snprintf(buffer.data(), buffer.size(), format, args...);
    Serial.println(buffer.data());
}

}  // namespace Logger
#endif