#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内存memory实现函数
* 1.函数realloc
* 2.函数free
* 3.函数memset
* 4.C库函数strcpy
* @author liweic
*/

CONTRACT MemoryFunction_2 : public platon::Contract{
	public:
    ACTION void init(){}

	CONST std::string getrealloc(){
	   std::string strTmp = "WasmTest2";
       char* buf = (char*)malloc(strTmp.size() + 1);
       char* newbuf = (char*)realloc(buf, strTmp.size() + 1);
       memset(newbuf, 0, strTmp.size()+1);
       strcpy(newbuf, strTmp.c_str());
       std::string str_temp (newbuf);
       free(newbuf);
       return str_temp;
    }
};

PLATON_DISPATCH(MemoryFunction_2, (init)(getrealloc))