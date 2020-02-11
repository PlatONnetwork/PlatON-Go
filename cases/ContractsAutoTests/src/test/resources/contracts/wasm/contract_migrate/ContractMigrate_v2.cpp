#include <platon/platon.hpp>
#include <string>
using namespace platon;

extern char const string_storage[] = "stringstorage";
extern char const string_contract_ower[] = "contract_ower";
extern char const uint8_storage[] = "uint8storage";

CONTRACT ContractMigrate : public platon::Contract{
   public:
      PLATON_EVENT1(transfer,std::string,std::string)

      ACTION void init(){
          platon::set_owner();
      }

      /**
       * 合约升级
       * 内置 platon_migrate 函数 return_address 参数被写入新合约地址
       * init_arg 参数为  magic number +  RLP(code, RLP("init", init_paras...))
       * transfer_value 为转到新合约地址的金额，gas_value 为预估消耗的 gas
       */  
      ACTION std::string migrate_contract(const bytes &init_arg, uint64_t transfer_value, uint64_t gas_value){
            if (platon::is_owner()) {
                return "invalid address";
            }

            Address return_address;
            platon_migrate_contract(return_address, init_arg, transfer_value, gas_value);
            PLATON_EMIT_EVENT1(transfer,return_address.toString(),return_address.toString());

            return return_address.toString();
      }

      CONST uint8_t get_balance(){
          return uint8balance.self();
      }
   private:
      platon::StorageType<uint8_storage, uint8_t> uint8balance;  

};

PLATON_DISPATCH(ContractMigrate, (init)(migrate_contract)(get_balance))
