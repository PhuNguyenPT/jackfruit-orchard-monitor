#include <unity.h>
#include <array>
#include <cstdint>
#include <cstdio>
#include <cstring>
#include "SHT40Common.h"

// ---------------------------------------------------------------------------
// setUp / tearDown (required by Unity on embedded targets)
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
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 65.0F, SHT40Poller::scaleHumidity(650));
}

void test_humidity_scaling_zero(void) {
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 0.0F, SHT40Poller::scaleHumidity(0));
}

void test_humidity_scaling_max(void) {
    // 1000 raw -> 100.0 %RH
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 100.0F, SHT40Poller::scaleHumidity(1000));
}

void test_temperature_scaling_positive(void) {
    // 301 raw -> 30.1 °C
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 30.1F, SHT40Poller::scaleTemperature(301));
}

void test_temperature_scaling_negative(void) {
    // -50 as int16 -> -5.0 °C
    // Cast to uint16_t as ModbusMaster returns uint16_t
    uint16_t raw = static_cast<uint16_t>(static_cast<int16_t>(-50));
    TEST_ASSERT_FLOAT_WITHIN(0.01F, -5.0F, SHT40Poller::scaleTemperature(raw));
}

void test_temperature_scaling_zero(void) {
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 0.0F, SHT40Poller::scaleTemperature(0));
}

// ---------------------------------------------------------------------------
// MQTT payload formatting tests
// ---------------------------------------------------------------------------
void test_payload_format_positive_temp(void) {
    std::array<char, SHT40Poller::kPayloadBufSize> payload{};
    snprintf(payload.data(), payload.size(),
             R"({"temperature":%.1f, "humidity":%.1f})",
             30.1F, 65.0F);
    TEST_ASSERT_EQUAL_STRING(
        R"({"temperature":30.1, "humidity":65.0})",
        payload.data());
}

void test_payload_format_negative_temp(void) {
    std::array<char, SHT40Poller::kPayloadBufSize> payload{};
    snprintf(payload.data(), payload.size(),
             R"({"temperature":%.1f, "humidity":%.1f})",
             -5.0F, 80.0F);
    TEST_ASSERT_EQUAL_STRING(
        R"({"temperature":-5.0, "humidity":80.0})",
        payload.data());
}

void test_payload_does_not_overflow(void) {
    std::array<char, SHT40Poller::kPayloadBufSize> payload{};
    int written = snprintf(payload.data(), payload.size(),
                           R"({"temperature":%.1f, "humidity":%.1f})",
                           -99.9F, 100.0F);
    TEST_ASSERT_TRUE(written > 0);
    TEST_ASSERT_TRUE(static_cast<size_t>(written) < SHT40Poller::kPayloadBufSize);
}

// ---------------------------------------------------------------------------
// MQTT topic formatting tests
// ---------------------------------------------------------------------------
void test_topic_format_addr_1(void) {
    std::array<char, SHT40Poller::kTopicBufSize> topic{};
    snprintf(topic.data(), topic.size(), SHT40Poller::kTopicTemplate, 1);
    TEST_ASSERT_EQUAL_STRING("sht40/1/data", topic.data());
}

void test_topic_format_addr_2(void) {
    std::array<char, SHT40Poller::kTopicBufSize> topic{};
    snprintf(topic.data(), topic.size(), SHT40Poller::kTopicTemplate, 2);
    TEST_ASSERT_EQUAL_STRING("sht40/2/data", topic.data());
}

void test_topic_does_not_overflow(void) {
    std::array<char, SHT40Poller::kTopicBufSize> topic{};
    int written = snprintf(topic.data(), topic.size(),
                           SHT40Poller::kTopicTemplate, 255);
    TEST_ASSERT_TRUE(written > 0);
    TEST_ASSERT_TRUE(static_cast<size_t>(written) < SHT40Poller::kTopicBufSize);
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------
#ifdef ARDUINO
#include <Arduino.h>

void setup() {
    delay(10000);
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