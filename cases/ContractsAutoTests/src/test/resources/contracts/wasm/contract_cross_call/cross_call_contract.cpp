#include <platon/platon.hpp>
#include <vector>
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

CONTRACT user : public platon::Contract {
    public:
        ACTION void init(){}
        ACTION uint64_t call_add_message(const std::string &target_address, const my_message &one_message, 
            uint64_t transfer_value, uint64_t gas_value) {
            platon::bytes paras = platon::cross_call_args("add_message", one_message);
            int32_t return_vale = platon::platon_call(target_address, paras, transfer_value, gas_value);
            return 0;
        }

        ACTION std::vector<my_message> delegate_call_add_message(const std::string &target_address, const my_message &one_message,
            uint64_t gas_value) {
            platon::bytes paras = platon::cross_call_args("add_message", one_message);
            platon::platon_delegate_call(target_address, paras, gas_value);
            std::vector<my_message> return_value;
            get_call_output(return_value);
            return return_value;
        }

};

PLATON_DISPATCH(user, (init)(call_add_message)(delegate_call_add_message))