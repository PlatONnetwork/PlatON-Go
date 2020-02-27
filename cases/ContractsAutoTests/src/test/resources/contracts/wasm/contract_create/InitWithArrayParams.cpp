#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT InitWithArrayParams : public platon::Contract{
   public:
      ACTION void init(const std::array<std::string,10>  &inArray){
         strarray.self() = inArray;
      }

      ACTION void set_array(const std::array<std::string,10>  &inArray){
          strarray.self() = inArray;
      }
      CONST std::array<std::string,10> get_array(){
          return strarray.self();
      }

   private:
      platon::StorageType<"strarray"_n, std::array<std::string,10>> strarray;
};

PLATON_DISPATCH(InitWithArrayParams, (init)(set_array)(get_array))
