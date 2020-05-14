#define TESTNET
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

CONTRACT storage_vector : public platon::Contract{
   public:


     ACTION void init(){}
      
      ACTION std::vector<my_message> add_message(const my_message &one_message){

          DEBUG("storage_vector add_message");
          arr.self().push_back(one_message);
          return arr.self();
      }
      CONST std::vector<my_message> get_message(const std::string &name){

          DEBUG("storage_vector get_message");
          return arr.self();
      }

      CONST uint64_t get_vector_size(){
          DEBUG("storage_vector get_vector_size");
          return arr.self().size();
      }


   private:
      platon::StorageType<"arr"_n, std::vector<my_message>> arr;
};

PLATON_DISPATCH(storage_vector, (init)(add_message)(get_message)(get_vector_size))
