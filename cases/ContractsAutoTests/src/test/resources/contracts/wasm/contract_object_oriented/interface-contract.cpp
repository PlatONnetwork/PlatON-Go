//
// Created by 许文 on 2020-02-06.
//

#include <platon/platon.hpp>
using namespace platon;

class Count {
public:
    virtual uint64_t getCount()=0;
    virtual void setCount(uint64_t count)=0;
private:
    uint64_t count;
};

CONTRACT migrate: public platon::Contract,public Count {
public:
    ACTION void init(){
        this->count = 0;
    };
    CONST uint64_t getCount() {
        return this->count;
    };
    ACTION void setCount(uint64_t count){
        this->count = count;
    };
private:
    uint64_t count;
};

PLATON_DISPATCH(migrate, (init)(getCount)(setCount));
