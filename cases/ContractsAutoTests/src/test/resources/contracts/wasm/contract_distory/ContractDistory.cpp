#include <platon/platon.hpp>
#include <string>
using namespace platon;

//extern char const string_storage[] = "stringstorage";

CONTRACT ContractDistory : public platon::Contract{
   public:
      ACTION void init(){
      }

      /**
       *  deploy --wasm ContractDistory.wasm
       *  //未传值前调用 返回0x80
       *  call --addr 0xA23c6D3f6D53B42fC62F0C6B50113AF43201978B --func get_string
       *  //正常传值
       *  invoke --addr 0xA23c6D3f6D53B42fC62F0C6B50113AF43201978B --func set_string --params {"name":"1"}
       *  //查询返回值 0x31
       *  call --addr 0xA23c6D3f6D53B42fC62F0C6B50113AF43201978B --func get_string
       *  //调用销毁合约
       *  invoke --addr 0xA23c6D3f6D53B42fC62F0C6B50113AF43201978B --func distory_contract
       *  //查询返回值 0x0
       *  call --addr 0xA23c6D3f6D53B42fC62F0C6B50113AF43201978B --func get_string
       *
       *  问题：合约销毁后
       */  
      ACTION int32_t distory_contract(){
            Address platon_address;
            platon_origin_caller(platon_address);
            return platon_destroy_contract(platon_address);
      }

      ACTION void set_string(std::string &name){
            stringstorage.self() = name;
      }

      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"string_storage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractDistory, (init)(distory_contract)(set_string)(get_string))
