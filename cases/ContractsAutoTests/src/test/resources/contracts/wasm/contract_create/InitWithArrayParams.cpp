#define TESTNET
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

      //array size
      CONST uint8_t get_array_size(){
          return strarray.self().size();
      }

      //判断数组中是否有指定元素的数据
      CONST bool get_array_contain_element(std::string &value){
          bool flg = false;
          for(auto iter = strarray.self().begin(); iter != strarray.self().end(); iter++) {
            DEBUG("InitWithArrayParams", "get_array_contain_element", *iter);
            if( *iter == value ){
                flg = true;
                break;
            }
          }
          return flg;
      }

   private:
      platon::StorageType<"strarray"_n, std::array<std::string,10>> strarray;
};

PLATON_DISPATCH(InitWithArrayParams, (init)(set_array)(get_array)(get_array_size)(get_array_contain_element))
