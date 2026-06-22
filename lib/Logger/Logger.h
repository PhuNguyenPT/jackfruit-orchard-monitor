#ifndef LOGGER_H
#define LOGGER_H
#include <Arduino.h>
#include <cstdarg>
#include <cstdio>
#include <ctime>
#include "esp_log.h"

namespace Logger {

inline int vprintf_handler(const char* format, va_list args) {
    struct tm timeinfo{};
    char timebuf[32];

    if (getLocalTime(&timeinfo)) {
        strftime(timebuf, sizeof(timebuf), "[%Y-%m-%d %H:%M:%S] ", &timeinfo);
    } else {
        snprintf(timebuf, sizeof(timebuf), "[%10lu s] ", millis() / 1000UL);
    }
    Serial.print(timebuf);

    char buf[256];
    vsnprintf(buf, sizeof(buf), format, args);
    Serial.print(buf);
    return 0;
}

inline void setup() {
    esp_log_set_vprintf(vprintf_handler);
}

}  // namespace Logger
#endif