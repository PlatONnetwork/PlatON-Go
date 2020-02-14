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


CONTRACT cross_delegate_call : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION uint64_t delegate_call_add_message(const std::string &target_address,
        const my_message &one_message, uint64_t gas) {

            platon::bytes params = platon::cross_call_args("add_message", one_message);
            if (platon_delegate_call(Address(target_address), params, gas)) {
                 DEBUG("Delegate call contract success", "address", target_address);
            } else {
                DEBUG("Delegate call contract fail", "address", target_address);
            }
            return 0;
        }

         CONST uint64_t get_vector_size(){
             platon::StorageType<"info_arr"_n, std::vector<my_message>> arr; // Must use local definitions for manipulating the corresponding keys in the account space
             return arr.self().size();
         }     

};

PLATON_DISPATCH(cross_delegate_call, (init)(delegate_call_add_message)(get_vector_size))