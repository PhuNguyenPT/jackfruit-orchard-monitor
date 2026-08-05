#include <unity.h>
#include <Arduino.h>
#include <HardwareSerial.h>
#include <ModbusMaster.h>
#include <cstdint>
#include <cstdio>
#include "PowerRail.h"
#include "SHT40Common.h"
#include "SHT40Poller.h"   // add — for kMinWarmupMs
#include "config.h"
#include "gpio.h"
#include "sht40.h"
#include <driver/gpio.h>

// cppcheck-suppress unusedFunction
void setUp(void) {}
// cppcheck-suppress unusedFunction
void tearDown(void) {}

namespace {

constexpr uint32_t kBaud = 4800U;

HardwareSerial modbusSerial(1);
ModbusMaster node;

struct PollResult {
    bool ok;
    uint8_t exceptionCode;
    bool zeroed;
};

auto pollOnce(uint8_t slaveAddr) -> PollResult {
    node.begin(slaveAddr, modbusSerial);
    const uint8_t result = node.readHoldingRegisters(0x0000, 2);
    if (result != ModbusMaster::ku8MBSuccess) {
        return PollResult{false, result, false};
    }
    const bool zeroed =
        SHT40Poller::isZeroedReading(node.getResponseBuffer(0), node.getResponseBuffer(1));
    return PollResult{true, 0U, zeroed};
}

auto measureSettleMs(uint8_t slaveAddr, uint32_t pollIntervalMs, uint32_t maxWaitMs) -> int32_t {
    PowerRail::parkForSleep();
    delay(2000U);  // confirmed sufficient for full rail discharge

    modbusSerial.begin(kBaud, SERIAL_8N1, static_cast<int>(XY485_RX), static_cast<int>(XY485_TX));

    const uint32_t startMs = millis();
    gpio_hold_dis(static_cast<gpio_num_t>(15U));
    pinMode(15U, OUTPUT);
    digitalWrite(15U, HIGH);

    while (millis() - startMs < maxWaitMs) {
        const PollResult r = pollOnce(slaveAddr);
        if (r.ok && !r.zeroed) {
            return static_cast<int32_t>(millis() - startMs);
        }
        delay(pollIntervalMs);
    }
    return -1;
}

}  // namespace

// ---------------------------------------------------------------------------
// Regression guard — confirms real settle time stays safely under the
// production kMinWarmupMs floor SHT40Poller::poll() waits out before its
// first read attempt. This is a hardware fact (2s internal refresh cycle
// per datasheet), not a retry-budget question — see SHT40Poller.h.
//
// Two things must hold:
//   1. settleMs must land at or below kMinWarmupMs — if it doesn't,
//      production's first read after the warm-up wait would still see a
//      zeroed register, reproducing the original bug.
//   2. settleMs must land above a floor (kSanityFloorMs) — if a trial
//      "settles" near-instantly, that's a sign the sensor isn't actually
//      power-cycling (e.g. the discharge delay above stopped being
//      sufficient), not a sign things got better.
// ---------------------------------------------------------------------------
void test_hw_settle_within_warmup_floor(void) {
    const uint32_t kPollIntervalMs = 50U;
    const uint32_t kSanityFloorMs = 1900U;  // real refresh cycle is 2000ms per datasheet
    const uint8_t kTrials = 5U;

    for (uint8_t trial = 0U; trial < kTrials; trial++) {
        for (size_t i = 0; i < NUM_SENSORS; i++) {
            const uint8_t addr = SLAVE_ADDRS.at(i);
            const int32_t settleMs =
                measureSettleMs(addr, kPollIntervalMs, SHT40Poller::kMinWarmupMs + 500U);

            char msg[112];
            snprintf(msg, sizeof(msg),
                     "trial=%d slave=%d settled at %ld ms (production waits %lu ms)", trial, addr,
                     static_cast<long>(settleMs),
                     static_cast<unsigned long>(SHT40Poller::kMinWarmupMs));
            TEST_ASSERT_TRUE_MESSAGE(
                settleMs >= static_cast<int32_t>(kSanityFloorMs) &&
                    settleMs <= static_cast<int32_t>(SHT40Poller::kMinWarmupMs),
                msg);
        }
    }
}

// ---------------------------------------------------------------------------
// Diagnostic sweep — unchanged, still useful for eyeballing the real curve.
// ---------------------------------------------------------------------------
void test_hw_settle_sweep(void) {
    const uint32_t kPollIntervalMs = 10U;
    const uint32_t kMaxWaitMs = 5000U;
    const uint8_t kTrialsPerSlave = 5U;

    for (size_t s = 0; s < NUM_SENSORS; s++) {
        const uint8_t addr = SLAVE_ADDRS.at(s);

        char header[48];
        snprintf(header, sizeof(header), "--- slave=%d ---", addr);
        Serial.println(header);

        for (uint8_t trial = 0U; trial < kTrialsPerSlave; trial++) {
            const int32_t settleMs = measureSettleMs(addr, kPollIntervalMs, kMaxWaitMs);

            char line[64];
            if (settleMs >= 0) {
                snprintf(line, sizeof(line), "  trial=%d settled at %ld ms", trial,
                         static_cast<long>(settleMs));
            } else {
                snprintf(line, sizeof(line), "  trial=%d NEVER settled within %lu ms", trial,
                         static_cast<unsigned long>(kMaxWaitMs));
            }
            Serial.println(line);
        }
    }

    TEST_ASSERT_TRUE(true);
}

// ---------------------------------------------------------------------------
void setup() {
    Serial.begin(1000000);
    delay(10000);
    UNITY_BEGIN();

    RUN_TEST(test_hw_settle_sweep);
    RUN_TEST(test_hw_settle_within_warmup_floor);

    UNITY_END();
}

void loop() {}