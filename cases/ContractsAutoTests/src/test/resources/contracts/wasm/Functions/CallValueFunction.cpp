#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 函数platon_call_value
*/

CONTRACT CallValueFunction : public platon::Contract{
	public:

    ACTION void init() {}

    CONST uint8_t get_platon_call_value() {
      uint8_t result;
      uint8_t val[32];
      result = platon_call_value(val);
      return result;
    }
};

PLATON_DISPATCH(CallValueFunction, (init)(get_platon_call_value))