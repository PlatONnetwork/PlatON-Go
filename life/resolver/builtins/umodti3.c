/* ===-- umodti3.c - Implement __umodti3 --------------------------------===
 * ===-------------------------------------------------------------------===
 */

#include <stdint.h>
#include <assert.h>

unsigned __int128 ___umodti3(uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb) {
    unsigned __int128 lhs = ha;
    unsigned __int128 rhs = hb;

    lhs <<= 64;
    lhs |=  la;

    rhs <<= 64;
    rhs |=  lb;

    assert(rhs != 0);

    lhs %= rhs;
    return lhs;
}
