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

class sub_my_message : public my_message {
   public:
      std::string from;
      std::string to;
      PLATON_SERIALIZE_DERIVED(sub_my_message, my_message,(from)(to))
};

//extern char const sub_my_message_vector[] = "info";

CONTRACT TwoInherit : public platon::Contract{
   public:
      ACTION void init(){
      }

      ACTION void add_sub_my_message(const sub_my_message &sub_one_message){
          info.self().push_back(sub_one_message);
      }

      CONST uint8_t get_sub_my_message_size(){
          return info.self().size();
      }

      CONST std::string get_sub_my_message_head(const uint8_t index){
          return info.self()[index].head;
      }

      CONST std::string get_sub_my_message_from(const uint8_t index){
          return info.self()[index].from;
      }


   private:
      platon::StorageType<"svector"_n, std::vector<sub_my_message>> info;
};

PLATON_DISPATCH(TwoInherit, (init)(add_sub_my_message)(get_sub_my_message_size)(get_sub_my_message_head)(get_sub_my_message_from))
