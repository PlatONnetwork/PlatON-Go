#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;


CONTRACT InitWithMap : public platon::Contract{
    public:
    ACTION void init(std::string &key, std::string &value){
        maps.self()[key] = value;
    }

    ACTION void set_map(std::string &key, std::string &value){
        maps.self()[key] = value;
    }


    CONST std::string get_map(std::string &key){
        return maps.self()[key];
    }   

    private:
    platon::StorageType<"initmap"_n, std::map<std::string,std::string>> maps;
};

PLATON_DISPATCH(InitWithMap, (init)(set_map)(get_map))
