#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 函数platon_caller
*/
CONTRACT CallerFunction : public platon::Contract{
	public:

    ACTION void init() {}

    CONST std::string get_platon_caller() {
        bytes addr;
      	addr.resize(20);
        ::platon_caller(addr.data());
        Address address(addr);
        //DEBUG("caller", "addr", address.toString());
        return Address(addr).toString();
    }
};

PLATON_DISPATCH(CallerFunction, (init)(get_platon_caller))