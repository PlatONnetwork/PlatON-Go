#define TESTNET
//
// Created by 许文 on 2020-02-06.
//

#include <platon/platon.hpp>
#include <string>
using namespace platon;

CONTRACT operate: public platon::Contract {
public:
    ACTION void init() {};
    CONST uint64_t add(uint64_t x, uint64_t y) {
        return x+y;
    };

    CONST std::string add(std::string x, std::string y) {
        return x + y;
    };
    uint64_t max(uint64_t x, uint64_t y) {
        if(x > y){
            return x;
        } else {
            return y;
        }
    };
    std::string max(std::string s1, std::string s2) {
        if(s1 > s2){
            return s1;
        } else {
            return s2;
        }
    };
};

PLATON_DISPATCH(operate, (init)(add))
