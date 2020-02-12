#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些转账和查询余额函数
* 函数platon_balance
* 函数platon_transfer
*/


CONTRACT TransferFunction : public platon::Contract{
	public:
    ACTION void init(){}

    CONST uint8_t get_Balance(bytes addr){
        addr.resize(20);
        uint8_t balance[32];
        uint8_t result;
        result = platon_balance(addr.data(), balance);
        return result;
    }

    ACTION int32_t get_platon_transfer(bytes addr){
        addr.resize(20);
        int32_t ptransfer;
        uint8_t *amount = new uint8_t(100);
        ptransfer = ::platon_transfer(addr.data(), amount, 18);
        return ptransfer;
    }

};

PLATON_DISPATCH(TransferFunction, (init)(get_Balance)(get_platon_transfer))