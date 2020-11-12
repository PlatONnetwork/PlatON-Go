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

CONTRACT origin_type : public platon::Contract{
   public:
      PLATON_EVENT1(hello_event, std::string, std::string, uint32_t)

      origin_type(){
         size_t vect_info_len = platon_get_state_length((uint8_t *)key.data(), key_len);
         if (0 == vect_info_len){
            return;
         }
         std::vector<byte> vect_value;
         vect_value.resize(vect_info_len);
         platon_get_state((uint8_t *)key.data(), key_len, vect_value.data(), vect_info_len);
         fetch(RLP(vect_value), vect_info);
      }

      ACTION void init(){}

      ACTION std::vector<my_message> add_message(const my_message &one_message){
            PLATON_EMIT_EVENT1(hello_event, "add_message", "event1", 1);
            DEBUG("origin_type add_message");
            vect_info.push_back(one_message);
            RLPStream stream;
            stream << vect_info;
            platon::bytesRef result = stream.out();
//            std::vector<byte> result = stream.out();
            platon_set_state((uint8_t *)key.data(), key_len, result.data(), result.size());
            return vect_info;
      }
      CONST std::vector<my_message> get_message(const std::string &name){
            PLATON_EMIT_EVENT1(hello_event, "get_message", "event2", 2);
            DEBUG("origin_type get_message");
            return vect_info;
      }

      CONST uint64_t get_vector_size(){
            DEBUG("origin_type get_vector_size");
            return vect_info.size();
      }

   private:
      std::vector<my_message> vect_info;
      std::string key = "info";
      size_t key_len = 4;
};

PLATON_DISPATCH(origin_type, (init)(add_message)(get_message)(get_vector_size))
