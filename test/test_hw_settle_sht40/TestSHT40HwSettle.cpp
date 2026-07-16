#include <unity.h>
#include <Arduino.h>
#include <HardwareSerial.h>
#include <ModbusMaster.h>
#include <cstdint>
#include <cstdio>
#include "config.h"
#include "gpio.h"
#include "sht40.h"

// cppcheck-suppress unusedFunction
void setUp(void) {}
// cppcheck-suppress unusedFunction
void tearDown(void) {}

namespace {

// Match SHT40Poller::init()'s default baud — bump this alongside the real
// code if/when you move to a higher baud rate, so the sweep reflects the
// actual bus speed you're testing against.
constexpr uint32_t kBaud = 4800U;

HardwareSerial modbusSerial(1);
ModbusMaster node;

// Result of a single poll attempt against one slave.
struct PollResult {
    bool ok;
    uint8_t exceptionCode;  // valid only when !ok
};

auto pollOnce(uint8_t slaveAddr) -> PollResult {
    node.begin(slaveAddr, modbusSerial);
    const uint8_t result = node.readHoldingRegisters(0x0000, 2);
    return PollResult{result == ModbusMaster::ku8MBSuccess, result};
}

}  // namespace

// ---------------------------------------------------------------------------
// Regression guard — confirms the CURRENT INTER_SLAVE_MS is sufficient for
// every configured slave, back-to-back, over several cycles.
// ---------------------------------------------------------------------------
void test_inter_slave_delay_is_sufficient(void) {
    modbusSerial.begin(kBaud, SERIAL_8N1, static_cast<int>(XY485_RX),
                        static_cast<int>(XY485_TX));

    const uint8_t kCycles = 10U;
    for (uint8_t cycle = 0U; cycle < kCycles; cycle++) {
        for (size_t i = 0; i < NUM_SENSORS; i++) {
            const uint8_t addr = SLAVE_ADDRS.at(i);
            const PollResult r = pollOnce(addr);

            char msg[80];
            snprintf(msg, sizeof(msg),
                     "cycle=%d slave=%d failed at INTER_SLAVE_MS=%d (code=0x%02X)", cycle, addr,
                     INTER_SLAVE_MS, r.exceptionCode);
            TEST_ASSERT_TRUE_MESSAGE(r.ok, msg);

            if (i + 1 < NUM_SENSORS) {
                delay(INTER_SLAVE_MS);
            }
        }
        delay(INTER_SLAVE_MS);  // gap before next cycle, same as steady-state polling
    }
}

// ---------------------------------------------------------------------------
// Diagnostic sweep — NOT a pass/fail check. For each candidate inter-slave
// delay, hammers every configured slave several times back-to-back and prints
// the success rate over Serial. Read the printed values to pick the real
// minimum delay by hand, the same way the mux channel-settle sweep was read.
// ---------------------------------------------------------------------------
void test_inter_slave_delay_sweep(void) {
    modbusSerial.begin(kBaud, SERIAL_8N1, static_cast<int>(XY485_RX),
                        static_cast<int>(XY485_TX));

    // Checkpoints in milliseconds — covers well below and well above the
    // current 200ms setting so you can see the full curve.
    const uint32_t checkpointsMs[] = {200, 100, 50, 20, 10, 5, 2, 0};
    const size_t numCheckpoints = sizeof(checkpointsMs) / sizeof(checkpointsMs[0]);
    const uint8_t kTrialsPerCheckpoint = 20U;

    for (size_t i = 0; i < numCheckpoints; i++) {
        const uint32_t delayMs = checkpointsMs[i];

        char header[48];
        snprintf(header, sizeof(header), "--- INTER_SLAVE_MS=%lu ---",
                 static_cast<unsigned long>(delayMs));
        Serial.println(header);

        for (size_t s = 0; s < NUM_SENSORS; s++) {
            const uint8_t addr = SLAVE_ADDRS.at(s);
            uint8_t successes = 0U;
            uint8_t lastFailCode = 0U;

            for (uint8_t trial = 0U; trial < kTrialsPerCheckpoint; trial++) {
                const PollResult r = pollOnce(addr);
                if (r.ok) {
                    successes++;
                } else {
                    lastFailCode = r.exceptionCode;
                }
                if (delayMs > 0) {
                    delay(delayMs);
                }
            }

            char line[80];
            snprintf(line, sizeof(line), "  slave=%d  %d/%d ok  (last fail code=0x%02X)", addr,
                     successes, kTrialsPerCheckpoint, lastFailCode);
            Serial.println(line);
        }
    }

    // Always "passes" — this test exists to produce log output, not to gate
    // anything. Read the printed success rates to decide the real delay floor.
    TEST_ASSERT_TRUE(true);
}

// ---------------------------------------------------------------------------
void setup() {
    Serial.begin(1000000);
    delay(10000);
    UNITY_BEGIN();

    RUN_TEST(test_inter_slave_delay_sweep);
    RUN_TEST(test_inter_slave_delay_is_sufficient);

    UNITY_END();
}

void loop() {}