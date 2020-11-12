#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT InitWithSetParams : public platon::Contract{
   public:
      ACTION void init(const std::set<std::string>  &inSet){
         strSet.self() = inSet;
      }

      ACTION void add_set(const std::set<std::string>  &inSet){
          strSet.self() = inSet;
      }
      CONST std::set<std::string> get_set(){
          return strSet.self();
      }

   private:
      platon::StorageType<"strset"_n, std::set<std::string>> strSet;
};

PLATON_DISPATCH(InitWithSetParams, (init)(add_set)(get_set))
