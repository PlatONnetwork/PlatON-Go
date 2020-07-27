#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

class message {
   public:
      message(){}
      message(const std::string &p_head):head(p_head){}
   private:
      std::string head;
      PLATON_SERIALIZE(message, (head))
};

class my_message : public message {
   public:
      my_message(){}
      my_message(const std::string &p_head, const std::string &p_body, const std::string &p_end):message(p_head), body(p_body), end(p_end){}
   private:
      std::string body;
      std::string end;
      PLATON_SERIALIZE_DERIVED(my_message, message, (body)(end))
};

CONTRACT InitWithObjectParams : public platon::Contract{
   public:
      ACTION void init(const my_message &one_message){
         info.self().push_back(one_message);
      }

      ACTION std::vector<my_message> add_message(const my_message &one_message){
          info.self().push_back(one_message);
          return info.self();
      }
      CONST std::vector<my_message> get_message(const std::string &name){
          return info.self();
      }

   private:
      platon::StorageType<"cinfo"_n, std::vector<my_message>> info;
};

PLATON_DISPATCH(InitWithObjectParams, (init)(add_message)(get_message))
