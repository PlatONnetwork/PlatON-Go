#include <platon/platon.hpp>
#include <string>
using namespace platon;

extern "C"{
    uint64_t platon_gas();
    uint64_t platon_gas_limit();
    uint64_t platon_gas_price();
}
CONTRACT PlatONSpecialFunctionsB : public platon::Contract{
	public:
    ACTION void init(){}

    CONST uint64_t getPlatONGas(){
        return platon_gas();
    }

    CONST uint64_t getPlatONGasLimit(){
        return platon_gas_limit();
    }

    CONST uint64_t getPlatONGasPrice(){
        return platon_gas_price();
    }

};

PLATON_DISPATCH(PlatONSpecialFunctionsB, (init)(getPlatONGas)(getPlatONGasLimit)(getPlatONGasPrice))