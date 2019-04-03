#pragma once
#include <stdint.h>
//#include <softfloat.hpp>

//extern "C" {
__int128 ___fixdfti(uint64_t);
__int128 ___fixsfti(uint32_t);
__int128 ___fixtfti( float128_t);
unsigned __int128 ___fixunsdfti(uint64_t);
unsigned __int128 ___fixunssfti(uint32_t);
unsigned __int128 ___fixunstfti(float128_t);
//double ___floattidf(__int128);
double ___floattidf(uint64_t l, uint64_t h);
//double ___floatuntidf(unsigned __int128);
double ___floatuntidf(uint64_t l, uint64_t h);

__int128 ___ashriti3(uint64_t, uint64_t, uint32_t);
__int128 ___divti3(uint64_t, uint64_t, uint64_t, uint64_t);
unsigned __int128 ___udivti3(uint64_t, uint64_t, uint64_t, uint64_t);
__int128 ___modti3(uint64_t, uint64_t, uint64_t, uint64_t);
unsigned __int128 ___umodti3(uint64_t, uint64_t, uint64_t, uint64_t);
__int128 ___multi3(uint64_t, uint64_t, uint64_t, uint64_t);

//}
