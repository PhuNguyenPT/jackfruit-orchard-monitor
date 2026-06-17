#include "SoilPoller.h"
#include <Arduino.h>
#include <array>
#include "Logger.h"
#include "MKE_S13.h"

namespace SoilPoller {

namespace {
static const uint8_t  kNumSamples    = 30U;
static const uint32_t kSampleDelayMs = 10U;

// ---------------------------------------------------------------------------
// MUX helpers
// ---------------------------------------------------------------------------
void selectChannel(uint8_t ch) {
    digitalWrite(SoilConfig::kMuxS0, (ch >> 0U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS1, (ch >> 1U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS2, (ch >> 2U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS3, (ch >> 3U) & 0x01U);
    delayMicroseconds(10U);
}

void enableBoard(uint8_t enPin, bool enable) {
    digitalWrite(enPin, enable ? LOW : HIGH);  // EN is active LOW
    delayMicroseconds(5U);
}

void disableAllBoards() {
    for (uint8_t b = 0U; b < SoilConfig::kNumBoards; b++) {
        enableBoard(SoilConfig::kBoards[b].enPin, false);
    }
}

// ---------------------------------------------------------------------------
// ADC averaging
// ---------------------------------------------------------------------------
uint16_t readAvg(uint8_t sigPin) {
    uint32_t sum = 0U;
    for (uint8_t s = 0U; s < kNumSamples; s++) {
        sum += static_cast<uint32_t>(analogRead(sigPin));
        delay(kSampleDelayMs);
    }
    return static_cast<uint16_t>(sum / kNumSamples);
}

// ---------------------------------------------------------------------------
// Single sensor read: isolate board, select channel, sample
// ---------------------------------------------------------------------------
uint16_t readSensor(uint8_t boardIdx, uint8_t channel) {
    const SoilConfig::MuxBoard& board = SoilConfig::kBoards[boardIdx];
    disableAllBoards();
    enableBoard(board.enPin, true);
    selectChannel(channel);
    delay(5U);
    const uint16_t raw = readAvg(board.sigPin);
    enableBoard(board.enPin, false);
    return raw;
}
}  // namespace

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------
void init() {
    pinMode(SoilConfig::kMuxS0, OUTPUT);
    pinMode(SoilConfig::kMuxS1, OUTPUT);
    pinMode(SoilConfig::kMuxS2, OUTPUT);
    pinMode(SoilConfig::kMuxS3, OUTPUT);
    selectChannel(0U);

    for (uint8_t b = 0U; b < SoilConfig::kNumBoards; b++) {
        const SoilConfig::MuxBoard& board = SoilConfig::kBoards[b];
        pinMode(board.enPin,  OUTPUT);
        pinMode(board.sigPin, INPUT);
    }
    disableAllBoards();
}

void poll(PubSubClient& mqttClient) {
    uint8_t sensorId = 0U;

    for (uint8_t b = 0U; b < SoilConfig::kNumBoards; b++) {
        for (uint8_t ch = 0U; ch < SoilConfig::kBoards[b].numCh; ch++) {
            const uint16_t raw     = readSensor(b, ch);
            const float    percent = toPercent(raw);

            Logger::log(Logger::Level::SUCCESS,
                        "Soil Sensor %d (MUX%d CH%d): raw=%d -> %.1f %%",
                        sensorId, b + 1U, ch, raw, percent);

            if (!mqttClient.connected()) {
                sensorId++;
                continue;
            }

            std::array<char, kTopicBufSize>   topic{};
            std::array<char, kPayloadBufSize> payload{};

            // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
            snprintf(topic.data(), topic.size(), kTopicTemplate, sensorId);
            // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
            snprintf(payload.data(), payload.size(),
                     R"({"moisture":%.1f, "raw":%d})", percent, raw);

            if (!mqttClient.publish(topic.data(), payload.data())) {
                Logger::log(Logger::Level::ERROR,
                            "MQTT Frame dropped. Publish failed for soil sensor %d.", sensorId);
            } else {
                Logger::log(Logger::Level::INFO, "MQTT Outbound -> [%s] Payload: %s",
                            topic.data(), payload.data());
            }

            sensorId++;
        }
    }
}

}  // namespace SoilPoller