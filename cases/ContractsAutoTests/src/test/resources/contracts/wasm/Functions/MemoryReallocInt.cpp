#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内存memory实现函数
* @author liweic
*/

CONTRACT MemoryReallocInt : public platon::Contract{
	public:
    ACTION void init(){}

	CONST int getrealloc(){
       int *p = (int *)malloc(sizeof(int));
       *p = 10;
       free(p);
       int *p1 = (int *)malloc(2*sizeof(int));
       *p1 = 100;
       free(p1);

       int *p_new=(int *)realloc(p, 10*sizeof(int));
       *p_new = 100;
       int temp = *p_new;
       free(p_new);
       return temp;
    }
};

PLATON_DISPATCH(MemoryReallocInt, (init)(getrealloc))