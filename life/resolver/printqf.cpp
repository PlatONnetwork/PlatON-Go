#include "printqf.h"

#include <string.h>
#include <sstream>
#include <limits>

extern "C" {
#include "./softfloat/source/include/softfloat.h"
}

const char* __printqf(uint64_t low, uint64_t high) {
    float128_t val{{low, high}};
    extFloat80_t val_approx;

    f128M_to_extF80M(&val, &val_approx);

    std::ostringstream oss;
    auto old_prec = oss.precision();
    oss.precision(std::numeric_limits<long double>::digits10);
    oss << (*(long double*)&val_approx);
    oss.precision(old_prec);
    static char buf[128+1];
    ::memset(buf, 0, sizeof(buf));
    buf[sizeof(buf) - 1] = '\0';
    for (size_t i = 0; i < oss.str().size(); i++) {
        buf[i] = oss.str()[i];
    }
    return &buf[0];
}