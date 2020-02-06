#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
* 验证内置的一些与链交互的函数
* 函数platon_caller
*/

extern const char contract_ower[] = "ower";

CONTRACT CallerFunction : public platon::Contract{
	public:
    ACTION void init(){}

    ACTION std::string init(const std::string address = ""){
        if (address.empty()) {
            Address platon_address;
            platon_caller(platon_address);
            contract_ower.self() = platon_address;
        } else {
            contract_ower.self() = Address(address);
        }
        return contract_ower.self().toString();
    }

    CONST std::string getPlatONCaller(){
        return contract_ower.self().toString();
    }

    private:
        StorageType<contract_ower, Address> contract_ower;

};

PLATON_DISPATCH(CallerFunction, (init)(getPlatONCaller))