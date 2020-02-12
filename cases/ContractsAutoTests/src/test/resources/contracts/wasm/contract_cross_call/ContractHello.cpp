#include <platon/platon.hpp>
#include <string>
//#include <Contract_hello.hpp>
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

CONTRACT hello : public platon::Contract{
   public:
      PLATON_EVENT1(hello, std::string, std::string, uint32_t)

     ACTION void init(){}
      
      ACTION std::vector<my_message> add_message(const my_message &one_message){
          PLATON_EMIT_EVENT1(hello, "add_message", "event1", 1);
          DEBUG("hello add_message");
          arr.self().push_back(one_message);
          return arr.self();
      }
      CONST std::vector<my_message> get_message(const std::string &name){
          PLATON_EMIT_EVENT1(hello, "get_message", "event2", 2);
          DEBUG("hello get_message");
          return arr.self();
      }

      CONST uint64_t get_vector_size(){
          DEBUG("hello get_vector_size");
          return arr.self().size();
      }


   private:
      platon::StorageType<"info_arr"_n, std::vector<my_message>> arr;
};

PLATON_DISPATCH(hello, (init)(add_message)(get_message)(get_vector_size))
