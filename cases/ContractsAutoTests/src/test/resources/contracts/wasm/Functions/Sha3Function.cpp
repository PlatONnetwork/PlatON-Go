#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 函数platon_sha3
*/

// extern char const contract_ower[] = "contract_ower";
// extern const uint8_t result[32] = {};

CONTRACT Sha3Function : public platon::Contract{
	public:
    ACTION void init(){}

    CONST uint32_t Sha3Result(){
        uint8_t src[5] = {0x01, 0x07, 0x10, 0x20, 0x30};
        uint8_t result[32] = {};
        platon_sha3(src, 5, result, 32);
        uint32_t a;
        memcpy(&a, result, sizeof(result));
        return a;
    }

};

PLATON_DISPATCH(Sha3Function, (init)(Sha3Result))