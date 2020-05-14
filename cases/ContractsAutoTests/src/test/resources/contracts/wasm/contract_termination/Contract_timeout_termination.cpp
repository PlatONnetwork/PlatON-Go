#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";
/**
 * 当value值大于100将会限入死循环
 * 将会终止合约
 */
CONTRACT ContractTimeoutTermination : public platon::Contract{
   public:
      ACTION void init(){
      }
 
      ACTION void forfunction(std::string name, uint64_t value){
            while(value>100){
               value=value+1;
            }
            stringstorage.self() = name;
      }

      CONST std::string get_string_storage(){
          return stringstorage.self();
      }

   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractTimeoutTermination, (init)(forfunction)(get_string_storage))
