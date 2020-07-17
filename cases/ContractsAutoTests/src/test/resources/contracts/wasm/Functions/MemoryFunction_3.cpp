#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内存memory实现函数
* 1.函数calloc
* 2.函数free
* 3.函数memcpy
* @author liweic
*/

CONTRACT MemoryFunction_3 : public platon::Contract{
	public:
    ACTION void init(){}

	CONST std::string getcalloc(){
	   std::string strTmp = "WasmTest3";
       char* buf = (char*)calloc(strTmp.size() + 1, strTmp.size() + 1);
       memset(buf, 0, strTmp.size()+1);
       memcpy(buf, strTmp.c_str(), strTmp.size() + 1);
       std::string str_temp (buf);
       free(buf);
       return str_temp;
    }
};

PLATON_DISPATCH(MemoryFunction_3, (init)(getcalloc))