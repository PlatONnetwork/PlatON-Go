//
// Created by 许文 on 2020-02-06.
//

#include <platon/platon.hpp>
using namespace platon;

CONTRACT abstract: public platon::Contract {
public:
    ACTION void init();
    CONST uint64_t getCount();
    ACTION void setCount(uint64_t count);
};

PLATON_DISPATCH(abstract, (init)(getCount)(setCount));

