#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";
//extern char const string_storage1[] = "stringstorage1";
/**
 * platon_panic 函数退出合约，此时会把用户的全部 gas 用完
 */
CONTRACT ContractPanic : public platon::Contract{
   public:
      ACTION void init(){
      }

      /**
       * 合约终止
       * platon 提供断言函数 platon_assert，断言失败会退出合约，此时会花费掉实际执行消耗的 gas。
       * VM returned with error err="execute function code: exec: transaction panic"
       */  
      ACTION void panic_contract(std::string name, uint64_t value){
            stringstorage.self() = name;
            ::platon_panic;
            stringstorage1.self() = name;
      }

      ACTION void set_string_storage(std::string &name){
            stringstorage.self() = name;
      }

      CONST std::string get_string_storage(){
          return stringstorage.self();
      }

      CONST std::string get_string_storage1(){
          return stringstorage1.self();
      }

   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
      platon::StorageType<"sstorage1"_n, std::string> stringstorage1;
};

PLATON_DISPATCH(ContractPanic, (init)(panic_contract)(set_string_storage)(get_string_storage)(get_string_storage1))
