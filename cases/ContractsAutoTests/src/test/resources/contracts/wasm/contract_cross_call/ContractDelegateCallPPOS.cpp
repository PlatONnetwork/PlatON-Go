#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

CONTRACT delegate_call_ppos : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION uint64_t delegate_call_ppos_send (std::string target_addr, std::string &in, uint64_t gas) {
            platon::bytes  input = fromHex(in);

            auto address_info = make_address(target_addr);
            if(address_info.second){
                if (platon_delegate_call(address_info.first, input, gas)) {
                DEBUG("delegate call contract delegate_call_ppos_send success", "address", target_addr);
                return 0;
            }
            }

            DEBUG("delegate call contract delegate_call_ppos_send fail", "address", target_addr);
            return 1;
        }

        CONST const std::string  delegate_call_ppos_query (std::string target_addr, std::string &in, uint64_t gas) {
            platon::bytes  input = fromHex(in);

            auto address_info = make_address(target_addr);
            if(address_info.second){
                if (platon_delegate_call(address_info.first, input, gas)) {
                DEBUG("delegate call contract delegate_call_ppos_query success", "address", target_addr);
                platon::bytes ret;
                size_t len = platon_get_call_output_length();
                ret.resize(len);
                platon_get_call_output(ret.data());
                std::string str = toHex(ret);
                DEBUG("delegate call contract delegate_call_ppos_query success", "ret", str);
                return str;
            }
            }

            DEBUG("delegate call contract delegate_call_ppos_query fail", "address", target_addr);
            return "";
        }

};

PLATON_DISPATCH(delegate_call_ppos, (init)(delegate_call_ppos_send)(delegate_call_ppos_query))