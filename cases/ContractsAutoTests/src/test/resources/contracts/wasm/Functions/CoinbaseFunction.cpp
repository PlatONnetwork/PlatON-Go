#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 函数platon_origin
*/


CONTRACT CoinbaseFunction : public platon::Contract{
	public:

    ACTION void init() {}

    CONST std::string get_platon_coinbase() {
        bytes addr;
      	addr.resize(20);
        ::platon_coinbase(addr.data());
        return Address(addr).toString();
    }
};

PLATON_DISPATCH(CoinbaseFunction, (init)(get_platon_coinbase))