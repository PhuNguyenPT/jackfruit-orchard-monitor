#include "SoilPoller.h"
#include <Arduino.h>
#include <array>
#include "Logger.h"

namespace SoilPoller {

namespace {
constexpr uint8_t  kMaxSensors    = 4U;
constexpr uint8_t  kNumSamples    = 30U;
constexpr uint32_t kSampleDelayMs = 10U;
uint8_t            gPins[kMaxSensors]{};
uint8_t            gCount = 0U;
}  // namespace

static uint16_t readAvg(uint8_t pin) {
    uint32_t sum = 0U;
    for (uint8_t s = 0U; s < kNumSamples; s++) {
        sum += static_cast<uint32_t>(analogRead(pin));
        delay(kSampleDelayMs);
    }
    return static_cast<uint16_t>(sum / kNumSamples);
}

void init(const uint8_t* pins, uint8_t count) {
    gCount = (count > kMaxSensors) ? kMaxSensors : count;
    for (uint8_t i = 0U; i < gCount; i++) {
        gPins[i] = pins[i];
        pinMode(gPins[i], INPUT);
    }
}

void poll(PubSubClient& mqttClient) {
    for (uint8_t i = 0U; i < gCount; i++) {
        const uint16_t raw     = readAvg(gPins[i]);
        const float    percent = toPercent(raw);

        Logger::log(Logger::Level::SUCCESS,
                    "Soil Sensor %d (GPIO%d): raw=%d -> %.1f %%",
                    i, gPins[i], raw, percent);

        if (!mqttClient.connected()) continue;

        std::array<char, kTopicBufSize>   topic{};
        std::array<char, kPayloadBufSize> payload{};

        // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
        snprintf(topic.data(), topic.size(), kTopicTemplate, i);
        // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
        snprintf(payload.data(), payload.size(),
                 R"({"moisture":%.1f, "raw":%d})", percent, raw);

        if (!mqttClient.publish(topic.data(), payload.data())) {
            Logger::log(Logger::Level::ERROR,
                        "MQTT Frame dropped. Publish failed for soil sensor %d.", i);
            continue;
        }
        Logger::log(Logger::Level::INFO, "MQTT Outbound -> [%s] Payload: %s",
                    topic.data(), payload.data());
    }
}

}  // namespace SoilPoller