#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内存memory实现函数
* 1.函数malloc
* 2.函数free
* 3.函数memset
* 4.C库函数strcpy
* @author liweic
*/

CONTRACT MemoryFunction_1 : public platon::Contract{
	public:
    ACTION void init(){}

	CONST std::string getmalloc(){
	   std::string strTmp = "WasmTest";
       char* buf = (char*)malloc(strTmp.size() + 1);
       memset(buf, 0, strTmp.size()+1);
       strcpy(buf, strTmp.c_str());
       std::string str_temp (buf);
       free(buf);
       return str_temp;
    }
};

PLATON_DISPATCH(MemoryFunction_1, (init)(getmalloc))