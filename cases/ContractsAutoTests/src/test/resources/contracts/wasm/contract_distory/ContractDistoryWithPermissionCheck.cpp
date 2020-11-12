#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";

CONTRACT ContractDistory : public platon::Contract{
   public:
      ACTION void init(){
        contract_ower.self() = platon_origin().toString();
      }

      ACTION int32_t distory_contract(){
        Address platon_address = platon_origin();
        if (contract_ower.self() != platon_address.toString()){
            return -1;
        }
        return platon_destroy(platon_address);
      }

      ACTION void set_string(std::string &name){
            contract_ower.self() = name;
      }

      CONST std::string get_string(){
          return contract_ower.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> contract_ower;
};

PLATON_DISPATCH(ContractDistory, (init)(distory_contract)(set_string)(get_string))
