#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT InitWithParams : public platon::Contract{
   public:
      ACTION void init(const std::map<std::string,std::string>  &inMap){
         strmap.self() = inMap;
      }

      ACTION void add_map(const std::map<std::string,std::string>  &inMap){
          strmap.self() = inMap;
      }
      CONST std::map<std::string,std::string> get_map(){
          return strmap.self();
      }

   private:
      platon::StorageType<"strmap"_n, std::map<std::string,std::string>> strmap;
};

PLATON_DISPATCH(InitWithParams, (init)(add_map)(get_map))
