#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";
 /**
   * PLATON_EVENT 测试一个主题
   * 编译：./platon-cpp ContractEmitEvent1.cpp -std=c++17
   * 部署：deploy --wasm ContractEmitEvent1.wasm
   * 调用：invoke --addr 0xd002CD0427bE17B0671B84A2221834116431aC29 --func one_emit_event1 --params {"name":"5","value":4}
   * 调用：invoke --addr 0x21cC984a2dbD9431F7b2ebfd564Ff6034b5887c2 --func one_emit_event1_args2 --params {"name":"5","value":4}
   * 查询：call --addr 0xd002CD0427bE17B0671B84A2221834116431aC29 --func get_string 
   */
CONTRACT ContractEmitEvent1ComplexParam : public platon::Contract{
   public:
      PLATON_EVENT1(transfer,std::string,uint32_t,std::list<std::string>)

      ACTION void init(){}
 
      ACTION void one_emit_event1(std::string name,uint32_t value,const std::list<std::string>  &inList){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT1(transfer,name,value,inList);
      }


      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
      platon::StorageType<"listvar"_n, std::list<std::string>> sList;
};

PLATON_DISPATCH(ContractEmitEvent1ComplexParam, (init)(one_emit_event1))
