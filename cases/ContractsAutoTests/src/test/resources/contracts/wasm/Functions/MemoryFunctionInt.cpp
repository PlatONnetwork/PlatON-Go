#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内存memory实现函数-整型存储
* 1.函数realloc
* 2.函数free
* 3.函数memset
* @author liweic
*/

CONTRACT MemoryFunctionInt : public platon::Contract{
	public:
    ACTION void init(){}

	CONST int getmallocint(){
       int intTmp = -1;
       int* buf = (int*)malloc(sizeof(intTmp));
       memcpy(buf, &intTmp, sizeof(intTmp));
       int int_temp = *buf;
       free(buf);
       return int_temp;
    }
};

PLATON_DISPATCH(MemoryFunctionInt, (init)(getmallocint))