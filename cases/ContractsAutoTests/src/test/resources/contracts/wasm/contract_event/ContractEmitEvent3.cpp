#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";
 /**
   * PLATON_EVENT 测试一个主题
   * 编译：./platon-cpp ContractEmitEvent3.cpp -std=c++17
   * 部署：deploy --wasm ContractEmitEvent3.wasm
   * 调用：invoke --addr 0xC13977ba78b5e42c910777fF13d34567b1a36D9B --func three_emit_event3 --params {"name":"5","nationality":"china","city":"hangzhou","value":4}
   * 调用：invoke --addr 0xC13977ba78b5e42c910777fF13d34567b1a36D9B --func three_emit_event3_args2 --params {"name":"5","nationality":"china","city":"hangzhou","value":4}
   * 查询：call --addr 0xBD891449A2403DF312572E9a40F161547819dD71 --func get_string 
   */
CONTRACT ContractEmitEvent3 : public platon::Contract{
   public:
      PLATON_EVENT3(transfer,std::string,std::string,std::string,uint32_t)
      PLATON_EVENT3(transfer2,std::string,std::string,std::string,uint32_t,uint32_t,std::string,std::string)

      ACTION void init(){}
 
      ACTION void three_emit_event3(std::string name,std::string nationality,std::string city,uint32_t value){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT3(transfer,name,nationality,city,value);
      }

      ACTION void three_emit_event3_args4(std::string name,std::string nationality,std::string city,uint32_t value1,uint32_t value2,std::string name1,std::string name2){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT3(transfer2,name,nationality,city,value1,value2,name1,name2);
      }

      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractEmitEvent3, (init)(three_emit_event3)(three_emit_event3_args4)(get_string))
