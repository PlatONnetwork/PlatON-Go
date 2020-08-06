#define TESTNET
// Created by 许文 on 2020-02-06.
//

#include <platon/platon.hpp>
#include <string>
using namespace platon;

class Count {
public:
    virtual uint64_t getCount()=0;
    virtual void setCount(const uint64_t &count)=0;
};

CONTRACT migrate: public platon::Contract,public Count {
public:
    ACTION void init(){
        this->count.self() = 0;
    };
    CONST uint64_t getCount() {
        return this->count.self();
    };
    ACTION void setCount(const uint64_t &count){
        this->count.self() = count;
    };
private:
    platon::StorageType<"count"_n, uint64_t> count;
};

PLATON_DISPATCH(migrate, (init)(getCount)(setCount));