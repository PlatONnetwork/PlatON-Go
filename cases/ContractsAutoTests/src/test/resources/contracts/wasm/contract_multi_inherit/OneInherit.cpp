#define TESTNET
#include <platon/platon.hpp>
#include <string>
#include <map>
using namespace platon;

class message {
   public:
      std::string head;
      PLATON_SERIALIZE(message, (head))
};

class my_message : public message {
   public:
      std::string body;
      std::string end;
      PLATON_SERIALIZE_DERIVED(my_message, message, (body)(end))
};

//extern char const my_message_vector[] = "info";
 /**
   * 单继承测试
   * 编译：./platon-cpp bbb.cpp
   * 部署：deploy --wasm StorageType_uintts8t.cpp.wasm
   * 调用：invoke --addr 0xa4E7351A774f8e48a3302b0972a8844454932Ffa --func set_int --params {"a":12}
   * 查询：call --addr 0xa4E7351A774f8e48a3302b0972a8844454932Ffa --func get_int
   */
CONTRACT OneInherit : public platon::Contract{
   public:
      ACTION void init(){
      }

      ACTION void add_my_message(const my_message &one_message){
          info.self().push_back(one_message);
      }

      CONST uint8_t get_my_message_size(){
          return info.self().size();
      }

      CONST std::string get_my_message_head(const uint8_t index){
          return info.self()[index].head;
      }

      CONST std::string get_my_message_body(const uint8_t index){
          return info.self()[index].body;
      }


   private:
      platon::StorageType<"mymvector"_n, std::vector<my_message>> info;
};

PLATON_DISPATCH(OneInherit, (init)(add_my_message)(get_my_message_size)(get_my_message_head)(get_my_message_body))
