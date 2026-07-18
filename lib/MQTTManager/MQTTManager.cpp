#include "MQTTManager.h"

#include <Arduino.h>
#include <esp_log.h>
#include <mqtt_client.h>

#include <atomic>

namespace MQTTManager {

namespace {
const char* TAG = "MQTT";
const uint32_t kPollMs = 20U;

esp_mqtt_client_handle_t client = nullptr;
std::atomic<bool> connected{false};
std::atomic<int> pendingAcks{0};

void eventHandler(void* /*handlerArgs*/, esp_event_base_t /*base*/, int32_t eventId,
                  void* eventData) {
    auto* event = static_cast<esp_mqtt_event_handle_t>(eventData);
    switch (static_cast<esp_mqtt_event_id_t>(eventId)) {
        case MQTT_EVENT_CONNECTED:
#ifdef MQTT_SECURE
            ESP_LOGI(TAG, "TLS session established. Connected to broker.");
#else
            ESP_LOGI(TAG, "Connected to broker (plain TCP, dev/test mode).");
#endif
            connected.store(true, std::memory_order_release);
            break;
        case MQTT_EVENT_DISCONNECTED:
            ESP_LOGW(TAG, "MQTT disconnected.");
            connected.store(false, std::memory_order_release);
            break;
        case MQTT_EVENT_PUBLISHED:
            ESP_LOGI(TAG, "PUBACK received for msg_id=%d.", event->msg_id);
            pendingAcks.fetch_sub(1, std::memory_order_acq_rel);
            break;
        case MQTT_EVENT_ERROR:
            ESP_LOGE(TAG, "MQTT error, type=%d.",
                     static_cast<int>(event->error_handle->error_type));
            break;
        default:
            break;
    }
}
}  // namespace

void setup(const esp_mqtt_client_config_t& cfg) {
    esp_mqtt_client_config_t localCfg = cfg;
    localCfg.keepalive = 30;

    client = esp_mqtt_client_init(&localCfg);
    esp_mqtt_client_register_event(client, static_cast<esp_mqtt_event_id_t>(ESP_EVENT_ANY_ID),
                                   eventHandler, nullptr);
    esp_mqtt_client_start(client);
}

auto connect(uint32_t timeoutMs) -> bool {
    const uint32_t start = millis();
    while (!connected.load(std::memory_order_acquire)) {
        if (millis() - start >= timeoutMs) {
            ESP_LOGW(TAG, "MQTT connect timed out after %lu ms.",
                     static_cast<unsigned long>(timeoutMs));
            return false;
        }
        delay(kPollMs);
    }
    return true;
}

auto isConnected() -> bool { return connected.load(std::memory_order_acquire); }

auto publish(const char* topic, const char* payload) -> int {
    const int msgId = esp_mqtt_client_publish(client, topic, payload, 0, /*qos=*/1, /*retain=*/0);
    if (msgId < 0) {
        ESP_LOGE(TAG, "Failed to queue publish on topic %s.", topic);
        return msgId;
    }
    pendingAcks.fetch_add(1, std::memory_order_acq_rel);
    return msgId;
}

auto waitForAcks(uint32_t timeoutMs) -> bool {
    const uint32_t start = millis();
    while (pendingAcks.load(std::memory_order_acquire) > 0) {
        if (millis() - start >= timeoutMs) {
            ESP_LOGW(TAG, "Timed out with %d PUBACK(s) still outstanding.",
                     pendingAcks.load(std::memory_order_acquire));
            return false;
        }
        delay(kPollMs);
    }
    return true;
}

void disconnect() {
    if (client != nullptr) {
        esp_mqtt_client_stop(client);
    }
}

}  // namespace MQTTManager