#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内存memory实现函数
* @author liweic
*/

CONTRACT MemoryCallocInt : public platon::Contract{
	public:
    ACTION void init(){}

	CONST int getcalloc(){
        int *p1 = (int*)calloc(5, sizeof(int));
        *p1 = 10;
        free(p1);
        int *p2 = (int*)calloc(10, 5*sizeof(int));
        *p2 = 50;
        free(p2);
        int *p3 = (int*)calloc(20, 10*sizeof(int));
        *p3 = 200;
        free(p3);

        int *p4 = (int*)calloc(50, 50*sizeof(int));
        *p4 = 2500;
        int temp = *p4;
        free(p4);
        return temp;
    }
};

PLATON_DISPATCH(MemoryCallocInt, (init)(getcalloc))