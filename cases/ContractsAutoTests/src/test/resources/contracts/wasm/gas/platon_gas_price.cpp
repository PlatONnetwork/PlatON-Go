#define TESTNET
#include <platon/platon.hpp>
#include "Gas.h"
using namespace platon;

CONTRACT platon_gas_price : public platon::Contract{
    public:
    ACTION void init() {}

    ACTION void test() {
        //Gas gas("platon_call");
        platon_gas();
        //gas.Reset("platon_ecrecover");
        //platon_ecrecover();
        //gas.Reset("platon_sha3");
        //platon_sha3();
    }
};

PLATON_DISPATCH(platon_gas_price, (init)(test))



