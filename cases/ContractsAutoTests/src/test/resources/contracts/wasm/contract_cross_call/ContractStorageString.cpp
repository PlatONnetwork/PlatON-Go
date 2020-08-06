#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT storge_str : public platon::Contract{
   public:

     ACTION void init(){}
      
      ACTION void set_string(const std::string &name){
          DEBUG("storge_str set_string", "name:", name);
          str.self() = name;
      }
      
      CONST const std::string get_string(){
          DEBUG("storge_str get_string", "name:", str.self());
          return str.self();
      }

      
   private:
      platon::StorageType<"str"_n, std::string> str;
};

PLATON_DISPATCH(storge_str, (init)(set_string)(get_string))