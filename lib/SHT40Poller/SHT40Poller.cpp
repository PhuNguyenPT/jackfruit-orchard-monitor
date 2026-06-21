#include "SHT40Poller.h"
#include <Arduino.h>
#include <HardwareSerial.h>
#include <ModbusMaster.h>
#include <array>
#include "Logger.h"
#include "TimeSync.h"

namespace SHT40Poller {

namespace {
HardwareSerial modbusSerial(1);  // UART1
ModbusMaster   node;
}  // namespace

void init(int rxPin, int txPin, uint32_t baud) {
    modbusSerial.begin(baud, SERIAL_8N1, rxPin, txPin);
}

void poll(uint8_t slaveAddr, PubSubClient& mqttClient) {
    node.begin(slaveAddr, modbusSerial);
    const uint8_t result = node.readHoldingRegisters(0x0000, 2);

    if (result == node.ku8MBSuccess) {
        const float hum  = scaleHumidity(node.getResponseBuffer(0));
        const float temp = scaleTemperature(node.getResponseBuffer(1));

        Logger::log(Logger::Level::SUCCESS,
                    "Sensor %d Readout: %.1f %%RH | %.1f C", slaveAddr, hum, temp);

        if (mqttClient.connected()) {
            std::array<char, kTopicBufSize>   topic{};
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

            if (!mqttClient.publish(topic.data(), payload.data())) {
                Logger::log(Logger::Level::ERROR, "MQTT Frame dropped. Publish failed.");
                return;
            }
            Logger::log(Logger::Level::INFO, "MQTT Outbound -> [%s] Payload: %s",
                        topic.data(), payload.data());
        }
    } else {
        Logger::log(Logger::Level::ERROR,
                    "Modbus fault on Sensor %d. Exception Code: 0x%02X", slaveAddr, result);
    }
}

}  // namespace SHT40Poller