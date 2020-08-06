#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

CONTRACT ContractMigrate : public platon::Contract{
   public:
      PLATON_EVENT1(transfer,std::string,std::string)

      ACTION void init(){
        addressstorage.self()= platon_caller().toString();
        DEBUG("origin caller is:",platon_caller().toString());
      }

      /**
       * 合约升级
       * 内置 platon_migrate 函数 return_address 参数被写入新合约地址
       * init_arg 参数为  magic number +  RLP(code, RLP("init", init_paras...))
       * transfer_value 为转到新合约地址的金额，gas_value 为预估消耗的 gas
       */
      ACTION std::string migrate_contract(const bytes &init_arg, uint64_t transfer_value, uint64_t gas_value){
            if(addressstorage.self() != platon_caller().toString()){
                DEBUG("please use contract create to migrate !!!");
            }

            //输出init_arg参数
            DEBUG("init_arg is :", toHex(init_arg));
            Address return_address;
            platon_migrate_contract(return_address, init_arg, transfer_value, gas_value);
            PLATON_EMIT_EVENT1(transfer,return_address.toString(),return_address.toString());
            DEBUG("new contract address:", return_address.toString());
            return return_address.toString();
      }

      ACTION void set_string(const std::string  &one_name){
          DEBUG("set_string:", one_name);
          stringstorage.self()= one_name;
      }

      CONST std::string get_string(){
          DEBUG("get_string:", stringstorage.self());
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
      platon::StorageType<"sstorage"_n, std::string> addressstorage;
};

PLATON_DISPATCH(ContractMigrate, (init)(migrate_contract)(set_string)(get_string))
