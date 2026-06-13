#include <unity.h>

#include <array>
#include <cstdint>
#include <cstdio>
#include <cstring>
#include "sht40.h"

// ---------------------------------------------------------------------------
// Constants not defined in sht40.h (mirror main.cpp anonymous namespace)
// ---------------------------------------------------------------------------
namespace {
constexpr float  SENSOR_SCALE     = 10.0F;
constexpr size_t TOPIC_BUF_SIZE   = 50U;
constexpr size_t PAYLOAD_BUF_SIZE = 100U;
}  // namespace
// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
static float scaleHumidity(uint16_t raw) {
    return static_cast<float>(raw) / SENSOR_SCALE;
}

static float scaleTemperature(uint16_t raw) {
    return static_cast<float>(static_cast<int16_t>(raw)) / SENSOR_SCALE;
}

// ---------------------------------------------------------------------------
// setUp / tearDown (required by Unity)
// ---------------------------------------------------------------------------
// cppcheck-suppress unusedFunction
void setUp(void) {}
// cppcheck-suppress unusedFunction
void tearDown(void) {}

// ---------------------------------------------------------------------------
// Scaling tests
// ---------------------------------------------------------------------------
void test_humidity_scaling_normal(void) {
    // 650 raw -> 65.0 %RH
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 65.0F, scaleHumidity(650));
}

void test_humidity_scaling_zero(void) {
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 0.0F, scaleHumidity(0));
}

void test_humidity_scaling_max(void) {
    // 1000 raw -> 100.0 %RH
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 100.0F, scaleHumidity(1000));
}

void test_temperature_scaling_positive(void) {
    // 301 raw -> 30.1 °C
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 30.1F, scaleTemperature(301));
}

void test_temperature_scaling_negative(void) {
    // -50 as int16 -> -5.0 °C
    // Cast to uint16_t as ModbusMaster returns uint16_t
    uint16_t raw = static_cast<uint16_t>(static_cast<int16_t>(-50));
    TEST_ASSERT_FLOAT_WITHIN(0.01F, -5.0F, scaleTemperature(raw));
}

void test_temperature_scaling_zero(void) {
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 0.0F, scaleTemperature(0));
}

// ---------------------------------------------------------------------------
// MQTT payload formatting tests
// ---------------------------------------------------------------------------
void test_payload_format_positive_temp(void) {
    std::array<char, PAYLOAD_BUF_SIZE> payload{};
    snprintf(payload.data(), payload.size(),
             R"({"temperature":%.1f, "humidity":%.1f})",
             30.1F, 65.0F);
    TEST_ASSERT_EQUAL_STRING(
        R"({"temperature":30.1, "humidity":65.0})",
        payload.data());
}

void test_payload_format_negative_temp(void) {
    std::array<char, PAYLOAD_BUF_SIZE> payload{};
    snprintf(payload.data(), payload.size(),
             R"({"temperature":%.1f, "humidity":%.1f})",
             -5.0F, 80.0F);
    TEST_ASSERT_EQUAL_STRING(
        R"({"temperature":-5.0, "humidity":80.0})",
        payload.data());
}

void test_payload_does_not_overflow(void) {
    std::array<char, PAYLOAD_BUF_SIZE> payload{};
    int written = snprintf(payload.data(), payload.size(),
                           R"({"temperature":%.1f, "humidity":%.1f})",
                           -99.9F, 100.0F);
    TEST_ASSERT_TRUE(written > 0);
    TEST_ASSERT_TRUE(static_cast<size_t>(written) < PAYLOAD_BUF_SIZE);
}

// ---------------------------------------------------------------------------
// MQTT topic formatting tests
// ---------------------------------------------------------------------------
void test_topic_format_addr_1(void) {
    std::array<char, TOPIC_BUF_SIZE> topic{};
    snprintf(topic.data(), topic.size(), MQTT_TOPIC_TEMPLATE, 1);
    TEST_ASSERT_EQUAL_STRING("sht40/sensor1/data", topic.data());
}

void test_topic_format_addr_2(void) {
    std::array<char, TOPIC_BUF_SIZE> topic{};
    snprintf(topic.data(), topic.size(), MQTT_TOPIC_TEMPLATE, 2);
    TEST_ASSERT_EQUAL_STRING("sht40/sensor2/data", topic.data());
}

void test_topic_does_not_overflow(void) {
    std::array<char, TOPIC_BUF_SIZE> topic{};
    int written = snprintf(topic.data(), topic.size(), MQTT_TOPIC_TEMPLATE, 255);
    TEST_ASSERT_TRUE(written > 0);
    TEST_ASSERT_TRUE(static_cast<size_t>(written) < TOPIC_BUF_SIZE);
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------
#ifdef ARDUINO
#include <Arduino.h>

void setup() {
    delay(2000);  // Give the board time to settle before tests run
    UNITY_BEGIN();

    RUN_TEST(test_humidity_scaling_normal);
    RUN_TEST(test_humidity_scaling_zero);
    RUN_TEST(test_humidity_scaling_max);
    RUN_TEST(test_temperature_scaling_positive);
    RUN_TEST(test_temperature_scaling_negative);
    RUN_TEST(test_temperature_scaling_zero);
    RUN_TEST(test_payload_format_positive_temp);
    RUN_TEST(test_payload_format_negative_temp);
    RUN_TEST(test_payload_does_not_overflow);
    RUN_TEST(test_topic_format_addr_1);
    RUN_TEST(test_topic_format_addr_2);
    RUN_TEST(test_topic_does_not_overflow);

    UNITY_END();
}

void loop() {}

#else  // native

int main(int argc, char** argv) {
    UNITY_BEGIN();

    RUN_TEST(test_humidity_scaling_normal);
    RUN_TEST(test_humidity_scaling_zero);
    RUN_TEST(test_humidity_scaling_max);
    RUN_TEST(test_temperature_scaling_positive);
    RUN_TEST(test_temperature_scaling_negative);
    RUN_TEST(test_temperature_scaling_zero);
    RUN_TEST(test_payload_format_positive_temp);
    RUN_TEST(test_payload_format_negative_temp);
    RUN_TEST(test_payload_does_not_overflow);
    RUN_TEST(test_topic_format_addr_1);
    RUN_TEST(test_topic_format_addr_2);
    RUN_TEST(test_topic_does_not_overflow);

    return UNITY_END();
}

#endif