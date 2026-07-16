#include <unity.h>
#include <Arduino.h>
#include <cstdint>
#include <cstdio>
#include "MKE_S13.h"

// cppcheck-suppress unusedFunction
void setUp(void) {}
// cppcheck-suppress unusedFunction
void tearDown(void) {}

namespace {

void configurePins() {
    for (const auto& board : SoilConfig::kBoards) {
        pinMode(board.enPin, OUTPUT);
        pinMode(board.sigPin, INPUT);
    }
    pinMode(SoilConfig::kMuxS0, OUTPUT);
    pinMode(SoilConfig::kMuxS1, OUTPUT);
    pinMode(SoilConfig::kMuxS2, OUTPUT);
    pinMode(SoilConfig::kMuxS3, OUTPUT);
}

void selectAndEnable(uint8_t enPin, uint8_t channel) {
    digitalWrite(SoilConfig::kMuxS0, (channel >> 0U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS1, (channel >> 1U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS2, (channel >> 2U) & 0x01U);
    digitalWrite(SoilConfig::kMuxS3, (channel >> 3U) & 0x01U);
    digitalWrite(enPin, LOW);
}

void disableAll() {
    for (const auto& board : SoilConfig::kBoards) {
        digitalWrite(board.enPin, HIGH);
    }
}

}  // namespace

// ---------------------------------------------------------------------------
// Regression guard (unchanged) — confirms the CURRENT settle delay is enough.
// ---------------------------------------------------------------------------
void test_channel_settle_is_sufficient(void) {
    configurePins();

    for (const auto& board : SoilConfig::kBoards) {
        for (uint8_t ch = 0; ch < board.numCh; ch++) {
            disableAll();
            selectAndEnable(board.enPin, ch);

            delayMicroseconds(100U);
            uint16_t a = analogRead(board.sigPin);
            delay(2);
            uint16_t b = analogRead(board.sigPin);

            disableAll();

            uint16_t diff = (a > b) ? (a - b) : (b - a);
            char msg[64];
            snprintf(msg, sizeof(msg), "board EN=%d ch=%d still drifting (a=%d b=%d)",
                     board.enPin, ch, a, b);
            TEST_ASSERT_TRUE_MESSAGE(diff < 30, msg);
        }
    }
}

// ---------------------------------------------------------------------------
// Diagnostic sweep — NOT a pass/fail check. Prints raw ADC values over Serial
// at increasing settle times after a mux switch so you can read the curve off
// the serial monitor and pick the real minimum delay by hand.
// ---------------------------------------------------------------------------
void test_channel_settle_sweep(void) {
    configurePins();

    // Checkpoints in microseconds — covers well below and well above the
    // current 1ms (1000us) setting so you can see the full curve.
    const uint32_t checkpointsUs[] = {0, 50, 100, 200, 500, 1000, 2000, 5000};
    const size_t numCheckpoints = sizeof(checkpointsUs) / sizeof(checkpointsUs[0]);

    for (const auto& board : SoilConfig::kBoards) {
        for (uint8_t ch = 0; ch < board.numCh; ch++) {
            char header[64];
            snprintf(header, sizeof(header), "--- board EN=%d ch=%d ---", board.enPin, ch);
            Serial.println(header);

            for (size_t i = 0; i < numCheckpoints; i++) {
                // Fresh switch each time so every checkpoint measures settle
                // time from a cold channel change, not from a warmed-up line.
                disableAll();
                delay(5);  // let the previous channel fully release first
                selectAndEnable(board.enPin, ch);

                if (checkpointsUs[i] > 0) {
                    delayMicroseconds(checkpointsUs[i]);
                }
                uint16_t val = analogRead(board.sigPin);

                disableAll();

                char line[64];
                snprintf(line, sizeof(line), "  t=%5luus  raw=%d",
                         static_cast<unsigned long>(checkpointsUs[i]), val);
                Serial.println(line);
            }
        }
    }

    // Always "passes" — this test exists to produce log output, not to gate
    // anything. Read the printed values to decide the real settle time.
    TEST_ASSERT_TRUE(true);
}

// ---------------------------------------------------------------------------
void setup() {
    Serial.begin(1000000);
    delay(10000);
    UNITY_BEGIN();

    RUN_TEST(test_channel_settle_sweep);
    RUN_TEST(test_channel_settle_is_sufficient);

    UNITY_END();
}

void loop() {}