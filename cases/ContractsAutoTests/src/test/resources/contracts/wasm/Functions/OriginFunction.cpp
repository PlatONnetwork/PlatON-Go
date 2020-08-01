#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 函数platon_origin
*/


CONTRACT OriginFunction : public platon::Contract{
	public:

    ACTION void init() {}

    CONST std::string get_platon_origin() {
        bytes addr;
      	addr.resize(20);
        ::platon_origin(addr.data());
        return Address(addr).toString();
    }
};

PLATON_DISPATCH(OriginFunction, (init)(get_platon_origin))