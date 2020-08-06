#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";
 /**
   * PLATON_EVENT 测试一个主题
   * 编译：./platon-cpp ContractEmitEvent2.cpp -std=c++17
   * 部署：deploy --wasm ContractEmitEvent2.wasm
   * 调用：invoke --addr 0x42495B0a691061BBda4F3caAe8721D7CFD3d7d55 --func two_emit_event2 --params {"name":"5","nationality":"china","value":4}
   * 调用：invoke --addr 0x334Bb5c07103cD54fB564655616D2C5194E7725a --func two_emit_event2_args2 --params {"name":"5","nationality":"china","value":4}
   * 查询：call --addr 0x42495B0a691061BBda4F3caAe8721D7CFD3d7d55 --func get_string 
   */
CONTRACT ContractEmitEvent2 : public platon::Contract{
   public:
      PLATON_EVENT2(transfer,std::string,std::string,uint32_t)
      PLATON_EVENT2(transfer2,std::string,std::string,uint32_t,uint32_t,std::string,std::string)

      ACTION void init(){}
 
      ACTION void two_emit_event2(std::string name,std::string nationality,uint32_t value){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT2(transfer,name,nationality,value);
      }

      //两个topic 4个参数
      ACTION void two_emit_event2_args4(std::string name,std::string nationality,uint32_t value1,uint32_t value2,std::string name1,std::string name2){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT2(transfer2,name,nationality,value1,value2,name1,name2);
      }

      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractEmitEvent2, (init)(two_emit_event2)(two_emit_event2_args4)(get_string))
