#ifndef __RESOLVER_PRINT128_H
#define __RESOLVER_PRINT128_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

const char* printi128(uint64_t lo, uint64_t hi);
const char* printui128(uint64_t lo, uint64_t hi);

#ifdef __cplusplus
}
#endif

#endif

