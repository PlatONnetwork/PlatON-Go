#define TESTNET
//
// Created by 许文 on 2020-02-06.
//

#include <platon/platon.hpp>
using namespace platon;

class base {
public:
    uint64_t get_count(){
        return this->count;
    };
    void set_count(uint64_t count) {
        this->count = count;
    };
private:
    uint64_t count;
};

PLATON_DISPATCH(base, (init)(getCount)(setCount));

CONTRACT inherit: public platon::Contract, public base {
public:
    ACTION void init(){
        base::set_count(10);
    };
    CONST uint64_t getCount() {
        return base::get_count();
    }
};

PLATON_DISPATCH(inherit, (init)(getCount));