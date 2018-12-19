/* ===-- modti3.c - Implement __modti3 ----------------------------------===
 * ===-------------------------------------------------------------------===
 */

#include <stdint.h>
#include <assert.h>

__int128 ___modti3(uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb) {
    __int128 lhs = ha;
    __int128 rhs = hb;

    lhs <<= 64;
    lhs |=  la;

    rhs <<= 64;
    rhs |=  lb;

    assert(rhs != 0);

    lhs %= rhs;
    return lhs;
}
