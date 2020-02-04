#include <platon/platon.hpp>
#include <string>
using namespace platon;

extern "C"{
    uint64_t platon_block_number();
    int64_t platon_timestamp();
}
CONTRACT PlatONSpecialFunctionsA : public platon::Contract{
	public:
    ACTION void init(){}

	CONST uint64_t getBlockNumber(){
        return platon_block_number();
    }

    CONST int64_t getTimestamp(){
        return platon_timestamp();
    }

};

PLATON_DISPATCH(PlatONSpecialFunctionsA, (init)(getBlockNumber)(getTimestamp))