#include "print128.h"

#include <stdint.h>
#include <iterator>

const char* printi128(uint64_t lo, uint64_t hi) {
    __int128 u = hi;
    u <<= 64;
    u |= lo;

    unsigned __int128 tmp = u < 0 ? -u : u;
    static char buffer[128+1];
    buffer[sizeof(buffer)-1] = '\0';
    char* d = std::end(buffer)-1;
    do
    {
        --d;
        *d = "0123456789"[ tmp % 10 ];
        tmp /= 10;
    } while ( tmp != 0 );
    if ( u < 0 ) {
        --d;
        *d = '-';
    }

    return d;
}

const char* printui128(uint64_t lo, uint64_t hi) {
    unsigned __int128 u = hi;
    u <<= 64;
    u |= lo;

    unsigned __int128 tmp = u;
    static char buffer[128+1];
    buffer[sizeof(buffer)-1] = '\0';
    char* d = std::end(buffer)-1;
    do
    {
        --d;
        *d = "0123456789"[ tmp % 10 ];
        tmp /= 10;
    } while ( tmp != 0 );
    return d;
}

