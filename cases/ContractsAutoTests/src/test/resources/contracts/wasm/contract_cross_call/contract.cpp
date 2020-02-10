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

extern char const contract_info[] = "info";

CONTRACT hello : public platon::Contract{
   public:
      PLATON_EVENT1(hello, std::string, std::string, uint32_t)
      ACTION void init(const my_message &one_message){
         info.self().push_back(one_message);
      }
      
      ACTION std::vector<my_message> add_message(const my_message &one_message){
          PLATON_EMIT_EVENT1(hello, "add_message", "event1", 1);
          info.self().push_back(one_message);
          return info.self();
      }
      CONST std::vector<my_message> get_message(const std::string &name){
          PLATON_EMIT_EVENT1(hello, "get_message", "event2", 2);
          return info.self();
      }

      ACTION void hello_abort(){platon_assert(0, "hello abort");}

      ACTION void hello_panic(){platon_panic();}

   private:
      platon::StorageType<contract_info, std::vector<my_message>> info;
};

PLATON_DISPATCH(hello, (init)(add_message)(get_message)(hello_abort)(hello_panic))
