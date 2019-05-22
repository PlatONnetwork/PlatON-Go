/* ===-- ashriti3.c - Implement __ashriti3 ------------------------------===
 * ===-------------------------------------------------------------------===
 */

#include <stdint.h>

__int128 ___ashriti3(uint64_t low, uint64_t high, uint32_t shift) {
    __int128 ret = high;
    ret <<= 64;
    ret |= low;
    ret >>= shift;
    return ret;
}
