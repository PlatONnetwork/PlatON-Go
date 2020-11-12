#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";

CONTRACT ContractDistory : public platon::Contract{
   public:
      ACTION void init(){
      }

      ACTION int32_t distory_contract(){
           Address platon_address = platon_origin();
           return platon_destroy(platon_address);
      }

      ACTION void set_string(std::string &name){
            stringstorage.self() = name;
      }

      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractDistory, (init)(distory_contract)(set_string)(get_string))
