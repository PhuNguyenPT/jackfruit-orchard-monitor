#include <unity.h>
#include <array>
#include <cstdint>
#include <cstdio>
#include "SoilCommon.h"

// ---------------------------------------------------------------------------
// setUp / tearDown (required by Unity on embedded targets)
// ---------------------------------------------------------------------------
// cppcheck-suppress unusedFunction
void setUp(void) {}
// cppcheck-suppress unusedFunction
void tearDown(void) {}

// ---------------------------------------------------------------------------
// toPercent() tests
// ---------------------------------------------------------------------------
void test_toPercent_dry_boundary(void) {
    // At kDryValue -> 0%
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 0.0F, SoilPoller::toPercent(3340));
}

void test_toPercent_above_dry_clamps_to_zero(void) {
    // Above kDryValue -> clamp to 0%
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 0.0F, SoilPoller::toPercent(4095));
}

void test_toPercent_wet_boundary(void) {
    // At kWetValue -> 100%
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 100.0F, SoilPoller::toPercent(1805));
}

void test_toPercent_below_wet_clamps_to_hundred(void) {
    // Below kWetValue -> clamp to 100%
    TEST_ASSERT_FLOAT_WITHIN(0.01F, 100.0F, SoilPoller::toPercent(0));
}

void test_toPercent_midpoint(void) {
    // Midpoint between kWetValue and kDryValue -> ~50%
    const uint16_t mid = static_cast<uint16_t>((3340U + 1805U) / 2U);
    TEST_ASSERT_FLOAT_WITHIN(1.0F, 50.0F, SoilPoller::toPercent(mid));
}

void test_toPercent_quarter_point(void) {
    // 75% of range from wet -> ~75% moisture
    const uint16_t quarter = static_cast<uint16_t>(1805U + (3340U - 1805U) / 4U);
    TEST_ASSERT_FLOAT_WITHIN(1.0F, 75.0F, SoilPoller::toPercent(quarter));
}

// ---------------------------------------------------------------------------
// Topic formatting tests
// ---------------------------------------------------------------------------
void test_topic_format_index_0(void) {
    std::array<char, SoilPoller::kTopicBufSize> topic{};
    snprintf(topic.data(), topic.size(), SoilPoller::kTopicTemplate, 0);
    TEST_ASSERT_EQUAL_STRING("mke-s13/0/data", topic.data());
}

void test_topic_format_index_3(void) {
    std::array<char, SoilPoller::kTopicBufSize> topic{};
    snprintf(topic.data(), topic.size(), SoilPoller::kTopicTemplate, 3);
    TEST_ASSERT_EQUAL_STRING("mke-s13/3/data", topic.data());
}

void test_topic_does_not_overflow(void) {
    std::array<char, SoilPoller::kTopicBufSize> topic{};
    int written = snprintf(topic.data(), topic.size(),
                           SoilPoller::kTopicTemplate, 255);
    TEST_ASSERT_TRUE(written > 0);
    TEST_ASSERT_TRUE(static_cast<size_t>(written) < SoilPoller::kTopicBufSize);
}

// ---------------------------------------------------------------------------
// Payload formatting tests
// ---------------------------------------------------------------------------
void test_payload_format_normal(void) {
    std::array<char, SoilPoller::kPayloadBufSize> payload{};
    snprintf(payload.data(), payload.size(),
             R"({"moisture":%.1f, "raw":%d})", 65.0F, 2500);
    TEST_ASSERT_EQUAL_STRING(
        R"({"moisture":65.0, "raw":2500})",
        payload.data());
}

void test_payload_format_fully_wet(void) {
    std::array<char, SoilPoller::kPayloadBufSize> payload{};
    snprintf(payload.data(), payload.size(),
             R"({"moisture":%.1f, "raw":%d})", 100.0F, 1805);
    TEST_ASSERT_EQUAL_STRING(
        R"({"moisture":100.0, "raw":1805})",
        payload.data());
}

void test_payload_does_not_overflow(void) {
    std::array<char, SoilPoller::kPayloadBufSize> payload{};
    int written = snprintf(payload.data(), payload.size(),
                           R"({"moisture":%.1f, "raw":%d})", 100.0F, 4095);
    TEST_ASSERT_TRUE(written > 0);
    TEST_ASSERT_TRUE(static_cast<size_t>(written) < SoilPoller::kPayloadBufSize);
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------
#ifdef ARDUINO
#include <Arduino.h>

void setup() {
    delay(10000);
    UNITY_BEGIN();

    RUN_TEST(test_toPercent_dry_boundary);
    RUN_TEST(test_toPercent_above_dry_clamps_to_zero);
    RUN_TEST(test_toPercent_wet_boundary);
    RUN_TEST(test_toPercent_below_wet_clamps_to_hundred);
    RUN_TEST(test_toPercent_midpoint);
    RUN_TEST(test_toPercent_quarter_point);
    RUN_TEST(test_topic_format_index_0);
    RUN_TEST(test_topic_format_index_3);
    RUN_TEST(test_topic_does_not_overflow);
    RUN_TEST(test_payload_format_normal);
    RUN_TEST(test_payload_format_fully_wet);
    RUN_TEST(test_payload_does_not_overflow);

    UNITY_END();
}

void loop() {}

#else  // native

int main(int argc, char** argv) {
    UNITY_BEGIN();

    RUN_TEST(test_toPercent_dry_boundary);
    RUN_TEST(test_toPercent_above_dry_clamps_to_zero);
    RUN_TEST(test_toPercent_wet_boundary);
    RUN_TEST(test_toPercent_below_wet_clamps_to_hundred);
    RUN_TEST(test_toPercent_midpoint);
    RUN_TEST(test_toPercent_quarter_point);
    RUN_TEST(test_topic_format_index_0);
    RUN_TEST(test_topic_format_index_3);
    RUN_TEST(test_topic_does_not_overflow);
    RUN_TEST(test_payload_format_normal);
    RUN_TEST(test_payload_format_fully_wet);
    RUN_TEST(test_payload_does_not_overflow);

    return UNITY_END();
}

#endif