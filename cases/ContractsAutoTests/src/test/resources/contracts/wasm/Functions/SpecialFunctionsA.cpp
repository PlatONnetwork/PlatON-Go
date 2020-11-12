#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 1.函数platon_block_number
* 2.函数platon_timestamp
*/

//extern "C"{
//    uint64_t platon_block_number();
//    int64_t platon_timestamp();
//}

CONTRACT SpecialFunctionsA : public platon::Contract{
	public:
    ACTION void init(){}

	CONST uint64_t getBlockNumber(){
        return platon_block_number();
    }

    CONST int64_t getTimestamp(){
        return platon_timestamp();
    }

};

PLATON_DISPATCH(SpecialFunctionsA, (init)(getBlockNumber)(getTimestamp))