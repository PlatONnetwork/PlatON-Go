#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内存memory实现函数
* @author liweic
*/

CONTRACT MemoryFunction_3 : public platon::Contract{
	public:
    ACTION void init(){}

	CONST int getcalloc(){
        int *p = (int*)calloc(5, sizeof(int));
        free(p);
        int *p2 = (int*)calloc(10, 5*sizeof(int));
        free(p2);
        int *p3 = (int*)calloc(20, 10*sizeof(int));
        free(p3);

        int *p4 = (int*)calloc(50, 50*sizeof(int));
        *p4 = 10;
        int *temp = p4;
        free(p4);
        return *temp;
    }
};

PLATON_DISPATCH(MemoryFunction_3, (init)(getcalloc))