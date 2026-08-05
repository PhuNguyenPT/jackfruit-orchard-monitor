#include "SHT40Poller.h"

#include <Arduino.h>
#include <HardwareSerial.h>
#include <ModbusMaster.h>

#include <array>

#include "MQTTManager.h"
#include "TimeSync.h"

namespace SHT40Poller {

namespace {
HardwareSerial modbusSerial(1);
ModbusMaster node;
const char* TAG = "SHT40";
const uint8_t kMaxAttempts = 3U;
const uint32_t kRetryDelayMs = 50U;
uint32_t g_railOnMs = 0U;
}  // namespace

void init(RxPin rxPin, TxPin txPin, uint32_t baud) {
    modbusSerial.begin(baud, SERIAL_8N1, static_cast<int>(rxPin), static_cast<int>(txPin));
    g_railOnMs = millis();
}

void poll(uint8_t slaveAddr) {
    const uint32_t elapsed = millis() - g_railOnMs;
    if (elapsed < kMinWarmupMs) {
        delay(kMinWarmupMs - elapsed);
    }

    node.begin(slaveAddr, modbusSerial);
    uint8_t result = node.readHoldingRegisters(0x0000, 2);

    uint8_t attempt = 1U;
    while (result == ModbusMaster::ku8MBSuccess &&
           isZeroedReading(node.getResponseBuffer(0), node.getResponseBuffer(1)) &&
           attempt < kMaxAttempts) {
        ESP_LOGW(TAG, "Sensor %d returned zeroed registers (attempt %d/%d), retrying...", slaveAddr,
                 attempt, kMaxAttempts);
        delay(kRetryDelayMs);
        result = node.readHoldingRegisters(0x0000, 2);
        attempt++;
    }

    if (result == ModbusMaster::ku8MBSuccess &&
        isZeroedReading(node.getResponseBuffer(0), node.getResponseBuffer(1))) {
        ESP_LOGE(TAG, "Sensor %d still zeroed after %d attempts — skipping publish.", slaveAddr,
                 kMaxAttempts);
        return;
    }

    if (result == ModbusMaster::ku8MBSuccess) {
        const float hum = scaleHumidity(node.getResponseBuffer(0));
        const float temp = scaleTemperature(node.getResponseBuffer(1));

        ESP_LOGI(TAG, "Sensor %d Readout: %.1f %%RH | %.1f C", slaveAddr, hum, temp);

        if (MQTTManager::isConnected()) {
            std::array<char, kTopicBufSize> topic{};
            std::array<char, kPayloadBufSize> payload{};

            // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
            snprintf(topic.data(), topic.size(), kTopicTemplate, slaveAddr);
            if (TimeSync::isSynced()) {
                // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
                snprintf(payload.data(), payload.size(),
                         R"({"temperature": %.1f, "humidity": %.1f, "ts": %ld})", temp, hum,
                         static_cast<long>(TimeSync::now()));
            } else {
                // NOLINTNEXTLINE(cppcoreguidelines-pro-type-vararg)
                snprintf(payload.data(), payload.size(),
                         R"({"temperature": %.1f, "humidity": %.1f})", temp, hum);
            }

            if (MQTTManager::publish(topic.data(), payload.data()) < 0) {
                ESP_LOGE(TAG, "MQTT publish failed to queue for topic %s.", topic.data());
                return;
            }
            ESP_LOGI(TAG, "MQTT Outbound -> [%s] Payload: %s", topic.data(), payload.data());
        }
    } else {
        ESP_LOGE(TAG, "Modbus fault on Sensor %d. Exception Code: 0x%02X", slaveAddr, result);
    }
}

}  // namespace SHT40Poller