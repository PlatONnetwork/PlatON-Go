//#define TESTNET
//#include <platon/platon.hpp>
//#include <string>
//using namespace platon;
//
////extern char const string_storage[] = "stringstorage";
// /**
//   * PLATON_EVENT 测试一个主题
//   * 编译：./platon-cpp ContractEmitEvent4.cpp -std=c++17
//   * 部署：deploy --wasm ContractEmitEvent4.wasm
//   * 调用：invoke --addr 0xC822Ed460348C8F84dd8bad979A01e51C0673370 --func four_emit_event4 --params {"name":"5","nationality":"china","city":"hangzhou","village":"yuhang","value":4}
//   * 调用：invoke --addr 0xBcad5dE91De1845Fe95812C206ED01d21fF7393F --func four_emit_event4_args2 --params {"name":"5","nationality":"china","city":"hangzhou","village":"yuhang","value":4}
//   * 查询：call --addr 0xC822Ed460348C8F84dd8bad979A01e51C0673370 --func get_string
//   */
//CONTRACT ContractEmitEvent4 : public platon::Contract{
//   public:
//      PLATON_EVENT4(transfer,std::string,std::string,std::string,std::string,uint32_t)
//      PLATON_EVENT4(transfer2,std::string,std::string,std::string,std::string,uint32_t,uint32_t,std::string,std::string)
//
//      ACTION void init(){}
//
//      ACTION void four_emit_event4(std::string name,std::string nationality,std::string city,std::string village,uint32_t value){
//           stringstorage.self() = name;
//           PLATON_EMIT_EVENT4(transfer,name,nationality,city,village,value);
//      }
//
//      ACTION void four_emit_event4_args4(std::string name,std::string nationality,std::string city,std::string village,uint32_t value1,uint32_t value2,std::string name1,std::string name2){
//           stringstorage.self() = name;
//           PLATON_EMIT_EVENT4(transfer2,name,nationality,city,village,value1,value2,name1,name2);
//      }
//
//      CONST std::string get_string(){
//          return stringstorage.self();
//      }
//   private:
//      platon::StorageType<"sstorage"_n, std::string> stringstorage;
//};
//
//PLATON_DISPATCH(ContractEmitEvent4, (init)(four_emit_event4)(four_emit_event4_args4)(get_string))
