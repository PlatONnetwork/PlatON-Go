#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";
 /**
   * PLATON_EVENT 测试零个主题
   * 编译：./platon-cpp ContractEmitEvent.cpp -std=c++17
   * 部署：deploy --wasm ContractEmitEvent.wasm
   * 调用：invoke --addr 0x1C31AE86DBDE69364a2cFBc90673df645e44e239 --func zero_emit_event_args2 --params {"name":"4"}
   * 查询：call --addr 0x1C31AE86DBDE69364a2cFBc90673df645e44e239 --func get_string 
   */
CONTRACT ContractEmitEvent : public platon::Contract{
   public:
      PLATON_EVENT0(transfer,std::string)
      PLATON_EVENT0(transfer2,std::string,std::string)
      PLATON_EVENT0(transfer3,std::string,std::string,uint32_t)

      ACTION void init(){}
 
      ACTION void zero_emit_event(std::string name){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT0(transfer,name);
      }

      ACTION void zero_emit_event_args2(std::string name){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT0(transfer2,name,name);
      }

      ACTION void zero_emit_event_args3(std::string name,uint32_t value){
         stringstorage.self() = name;
         PLATON_EMIT_EVENT0(transfer3,name,name,value);
      }

      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractEmitEvent, (init)(zero_emit_event)(zero_emit_event_args2)(zero_emit_event_args3)(get_string))
