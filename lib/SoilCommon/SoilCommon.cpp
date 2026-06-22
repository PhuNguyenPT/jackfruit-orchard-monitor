#include "SoilCommon.h"

namespace SoilPoller {

auto toPercent(uint16_t raw) -> float {
    if (raw >= kDryValue) {
        return 0.0F;
    }
    if (raw <= kWetValue) {
        return 100.0F;
    }
    return (1.0F -
            static_cast<float>(raw - kWetValue) / static_cast<float>(kDryValue - kWetValue)) *
           100.0F;
}

}  // namespace SoilPoller