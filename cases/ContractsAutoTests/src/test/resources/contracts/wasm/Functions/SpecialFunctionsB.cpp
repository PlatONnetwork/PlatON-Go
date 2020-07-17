#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 1.函数platon_gas
* 2.函数platon_gas_limit
* 3.函数platon_gas_price
*/

//extern "C"{
//    uint64_t platon_gas();
//    uint64_t platon_gas_limit();
//    uint64_t platon_gas_price();
//}

CONTRACT SpecialFunctionsB : public platon::Contract{
	public:
    ACTION void init(){}

    CONST uint64_t getPlatONGas(){
        return platon_gas();
    }

    CONST uint64_t getPlatONGasLimit(){
        return platon_gas_limit();
    }

    CONST std::string getPlatONGasPrice(){
    	u128 gasprice = platon_gas_price();
    	return std::to_string(gasprice);
    }

};

PLATON_DISPATCH(SpecialFunctionsB, (init)(getPlatONGas)(getPlatONGasLimit)(getPlatONGasPrice))