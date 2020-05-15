#define TESTNET
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


CONTRACT delegate_call_origin_type : public platon::Contract {
    public:
        ACTION void init(){}

        delegate_call_origin_type(){
            size_t vect_info_len = platon_get_state_length((uint8_t *)key.data(), key_len);
            if (0 == vect_info_len){
              return;
            }
            std::vector<byte> vect_value;
            vect_value.resize(vect_info_len);
            platon_get_state((uint8_t *)key.data(), key_len, vect_value.data(), vect_info_len);
            fetch(RLP(vect_value), vect_info);
        }

        ACTION uint64_t delegate_call_add_message(const std::string &target_address,
        const my_message &one_message, uint64_t gas) {
            platon::bytes params = platon::cross_call_args("add_message", one_message);

            auto address_info = make_address(target_address);
            if(address_info.second){
                if (platon_delegate_call(address_info.first, params, gas)) {
                    DEBUG("Delegate call contract success", "address", target_address);
            } else {
                DEBUG("Delegate call contract fail", "address", target_address);
            }
            }
            return 0;
        }

        CONST uint64_t get_vector_size(){
            return vect_info.size();
        } 

    private:
        std::vector<my_message> vect_info;
        std::string key = "info";
        size_t key_len = 4;
};

PLATON_DISPATCH(delegate_call_origin_type, (init)(delegate_call_add_message)(get_vector_size))