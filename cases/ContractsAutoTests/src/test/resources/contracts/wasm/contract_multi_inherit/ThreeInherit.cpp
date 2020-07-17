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

class greate_sub_my_message : public sub_my_message {
   public:
      std::string level;
      std::string desc;
      PLATON_SERIALIZE_DERIVED(greate_sub_my_message, sub_my_message,(level)(desc))
};


CONTRACT ThreeInherit : public platon::Contract{
   public:
      ACTION void init(){
      }

      ACTION void add_greate_sub_my_message(const greate_sub_my_message &my_greate_sub_one_message){
          info.self().push_back(my_greate_sub_one_message);
      }

      CONST uint8_t get_greate_sub_my_message_size(){
          return info.self().size();
      }

      CONST std::string get_greate_sub_my_message_head(const uint8_t index){
          return info.self()[index].head;
      }

      CONST std::string get_greate_sub_my_message_desc(const uint8_t index){
          return info.self()[index].desc;
      }


   private:
      platon::StorageType<"gsmvector"_n, std::vector<greate_sub_my_message>> info;
};

PLATON_DISPATCH(ThreeInherit, (init)(add_greate_sub_my_message)(get_greate_sub_my_message_size)(get_greate_sub_my_message_head)(get_greate_sub_my_message_desc))
