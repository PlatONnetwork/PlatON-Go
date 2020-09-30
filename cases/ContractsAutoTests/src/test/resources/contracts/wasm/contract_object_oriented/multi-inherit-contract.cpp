#define TESTNET
//
// Created by 许文 on 2020-02-06.
//

#include <platon/platon.hpp>
#include <string>
using namespace platon;

class animal {
public:
    void setName(std::string name) {
        this->name = name;
    };
    std::string getName() {
        return this->name;
    }
private:
    std::string name;
};

class voice {
public:
    std::string utterance() {
        return this->shut;
    };
    void setShut(std::string shut) {
        this->shut = shut;
    };
private:
    std::string shut;
};

CONTRACT dog: public platon::Contract,public animal, public voice {
public:
    ACTION void init(){
        voice::setShut("汪汪汪");
    };
    ACTION void setName(std::string name) {
        animal::setName(name);
    };
    CONST std::string getName(){
        return animal::getName();
    }
    CONST std::string shut() {
        return voice::utterance();
    }
};

PLATON_DISPATCH(dog, (init)(setName)(getName)(shut));


