#define TESTNET
#include <platon/platon.hpp>
#include <string>
#include <map>
using namespace platon;

class message {
   public:
      std::string head;
      uint32_t age;
      uint64_t money;
      PLATON_SERIALIZE(message, (head)(age)(money))
};

class my_message : public message {
   public:
      std::string body;
      std::string end;
      PLATON_SERIALIZE_DERIVED(my_message, message, (body)(end))
};

CONTRACT OneInheritWithMultiDataType : public platon::Contract{
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

      CONST uint32_t get_my_message_age(const uint8_t index){
          return info.self()[index].age;
      }

      CONST uint64_t get_my_message_money(const uint8_t index){
          return info.self()[index].money;
      }

      CONST std::string get_my_message_body(const uint8_t index){
          return info.self()[index].body;
      }


   private:
      platon::StorageType<"mymvector"_n, std::vector<my_message>> info;
};

PLATON_DISPATCH(OneInheritWithMultiDataType, (init)(add_my_message)(get_my_message_size)(get_my_message_head)(get_my_message_age)(get_my_message_money)(get_my_message_body))
