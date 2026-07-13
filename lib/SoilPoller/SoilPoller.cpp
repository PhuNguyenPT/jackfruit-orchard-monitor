#include "SoilPoller.h"

#include <Arduino.h>
#include <driver/gpio.h>

#include <array>

#include "MKE_S13.h"
#include "MQTTManager.h"
#include "TimeSync.h"

namespace SoilPoller {

namespace {
const uint8_t kNumSamples = 30U;
const uint32_t kSampleDelayMs = 10U;
const char* TAG = "Soil";

enum class BoardIdx : uint8_t {};
enum class ChannelIdx : uint8_t {};

void selectChannel(uint8_t channel) {
    digitalWrite(SoilConfig::kMuxS0, (channel >> 0U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS1, (channel >> 1U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS2, (channel >> 2U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS3, (channel >> 3U) & 0x01U);
    delayMicroseconds(10U);
}

void enableBoard(uint8_t enPin, bool enable) {
    digitalWrite(enPin, enable ? LOW : HIGH);
    delayMicroseconds(5U);
}

void disableAllBoards() {
    for (const auto& board : SoilConfig::kBoards) {
        enableBoard(board.enPin, false);
    }
}

auto readAvg(uint8_t sigPin) -> uint16_t {
    uint32_t sum = 0U;
    for (uint8_t sampleIdx = 0U; sampleIdx < kNumSamples; sampleIdx++) {
        sum += static_cast<uint32_t>(analogRead(sigPin));
        delay(kSampleDelayMs);
    }
    return static_cast<uint16_t>(sum / kNumSamples);
}

auto readSensor(BoardIdx boardIdx, ChannelIdx chanIdx) -> uint16_t {
    const SoilConfig::MuxBoard& board = SoilConfig::kBoards.at(static_cast<uint8_t>(boardIdx));
    disableAllBoards();
    enableBoard(board.enPin, true);
    selectChannel(static_cast<uint8_t>(chanIdx));
    delay(5U);
    const uint16_t raw = readAvg(board.sigPin);
    enableBoard(board.enPin, false);
    return raw;
}
}  // namespace

void init() {
    // Release any hold left over from the previous deep-sleep cycle before
    // touching pinMode() — held pins ignore reconfiguration until released.
    gpio_hold_dis(static_cast<gpio_num_t>(SoilConfig::kMuxS0));
    gpio_hold_dis(static_cast<gpio_num_t>(SoilConfig::kMuxS1));
    gpio_hold_dis(static_cast<gpio_num_t>(SoilConfig::kMuxS2));
    gpio_hold_dis(static_cast<gpio_num_t>(SoilConfig::kMuxS3));
    for (const auto& board : SoilConfig::kBoards) {
        gpio_hold_dis(static_cast<gpio_num_t>(board.enPin));
    }

    pinMode(SoilConfig::kMuxS0, OUTPUT);
    pinMode(SoilConfig::kMuxS1, OUTPUT);
    pinMode(SoilConfig::kMuxS2, OUTPUT);
    pinMode(SoilConfig::kMuxS3, OUTPUT);
    selectChannel(0U);

    for (const auto& board : SoilConfig::kBoards) {
        pinMode(board.enPin, OUTPUT);
        pinMode(board.sigPin, INPUT);
    }
    disableAllBoards();
}

void poll(PubSubClient& mqttClient) {
    uint8_t sensorId = 0U;

    for (uint8_t boardIdx = 0U; boardIdx < SoilConfig::kNumBoards; boardIdx++) {
        for (uint8_t chanIdx = 0U; chanIdx < SoilConfig::kBoards.at(boardIdx).numCh; chanIdx++) {
            const uint16_t raw = readSensor(BoardIdx{boardIdx}, ChannelIdx{chanIdx});
            const float percent = toPercent(raw);

            ESP_LOGI(TAG, "Soil Sensor %d (MUX%d CH%d): raw=%d -> %.1f %%", sensorId, boardIdx + 1U,
                     chanIdx, raw, percent);

            if (!MQTTManager::isConnected(mqttClient)) {
                sensorId++;
                continue;
            }

            std::array<char, kTopicBufSize> topic{};
            std::array<char, kPayloadBufSize> payload{};

            // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
            snprintf(topic.data(), topic.size(), kTopicTemplate, sensorId);
            if (TimeSync::isSynced()) {
                // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
                snprintf(payload.data(), payload.size(),
                         R"({"moisture": %.1f, "raw": %d, "ts": %ld})", percent, raw,
                         static_cast<long>(TimeSync::now()));
            } else {
                // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
                snprintf(payload.data(), payload.size(), R"({"moisture": %.1f, "raw": %d})",
                         percent, raw);
            }
            if (!MQTTManager::publish(mqttClient, topic.data(), payload.data())) {
                ESP_LOGE(TAG, "MQTT Frame dropped. Publish failed for soil sensor %d.", sensorId);
            } else {
                ESP_LOGI(TAG, "MQTT Outbound -> [%s] Payload: %s", topic.data(), payload.data());
            }

            sensorId++;
        }
    }
}

void parkForSleep() {
    disableAllBoards();  // EN = HIGH on both boards
    selectChannel(0U);   // deterministic select-line state

    // All four are RTC-capable — hold is meaningful and will survive deep sleep.
    gpio_hold_en(static_cast<gpio_num_t>(SoilConfig::kMuxS0));
    gpio_hold_en(static_cast<gpio_num_t>(SoilConfig::kMuxS1));
    gpio_hold_en(static_cast<gpio_num_t>(SoilConfig::kMuxS2));
    gpio_hold_en(static_cast<gpio_num_t>(SoilConfig::kMuxS3));

    // GPIO13/4 are RTC-capable now (post GPIO18/19 rewire), so this hold
    // actually persists through esp_deep_sleep_start() — unlike the old
    // GPIO18/19 assignment, where this call was a no-op during deep sleep.
    for (const auto& board : SoilConfig::kBoards) {
        gpio_hold_en(static_cast<gpio_num_t>(board.enPin));
    }
}

}  // namespace SoilPoller