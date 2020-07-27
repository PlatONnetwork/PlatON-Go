#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";
/**
 * platon_assert
 */
CONTRACT ContractTermination : public platon::Contract{
   public:
      ACTION void init(){
      }
 
      ACTION bool transfer_assert(std::string name, uint64_t value){
            platon_assert(value >= 100, "bad value");
            stringstorage.self() = name;
            return false;
      }

      CONST std::string get_string_storage(){
          return stringstorage.self();
      }

   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractTermination, (init)(transfer_assert)(get_string_storage))
