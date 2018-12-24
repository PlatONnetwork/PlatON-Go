#include "printqf.h"

#include <sstream>

extern "C" {
#include "../softfloat/include/softfloat.h"
}

const char* __printqf(uint64_t low, uint64_t high) {
    float128_t val{{low, high}};
    extFloat80_t val_approx;

    f128M_to_extF80M(&val, &val_approx);

    std::ostringstream oss;
    oss << (*(long double*)&val_approx);
    static char buf[128+1];
    buf[sizeof(buf) - 1] = '\0';
    for (size_t i = 0; i < oss.str().size(); i++) {
        buf[i] = oss.str()[i];
    }
    return &buf[0];
}