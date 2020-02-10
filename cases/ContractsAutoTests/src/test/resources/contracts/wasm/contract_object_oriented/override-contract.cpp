//
// Created by 许文 on 2020-02-04.
//

#include <platon/platon.hpp>
using namespace platon;

class operate {
public:
    int add(int x, int y) {
        return x + y;
    }
    int minus(int x, int y) {
        return x - y;
    }
};

CONTRACT arithmetic: public platon::Contract, public operate {
public:
    ACTION void init() {};
    CONST int add(int x, int y) {
        return x+2*y;
    };
    CONST int minus(int x, int y) {
        return x-2*y;
    };
};

PLATON_DISPATCH(arithmetic, (init)(add)(minus))
