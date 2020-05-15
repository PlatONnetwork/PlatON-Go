#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

CONTRACT call_ppos : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION uint64_t cross_call_ppos_send (std::string target_addr, std::string &in, uint64_t value, uint64_t gas) {
            platon::bytes  input = fromHex(in);

            auto address_info = make_address(target_address);
            if(address_info.second){
                if (platon_call(address_info.first, input, value, gas)) {
                DEBUG("cross call contract cross_call_ppos_send success", "address", target_addr);
                return 0;
            }
            }
            DEBUG("cross call contract cross_call_ppos_send fail", "address", target_addr);
            return 1;
        }

        CONST const std::string  cross_call_ppos_query (std::string target_addr, std::string &in, uint64_t value, uint64_t gas) {
            platon::bytes  input = fromHex(in);

            auto address_info = make_address(target_address);
            if(address_info.second){
            if (platon_call(Address(target_addr), input, value, gas)) {
                DEBUG("cross call contract cross_call_ppos_query success", "address", target_addr);
                platon::bytes ret;
                size_t len = platon_get_call_output_length();
                ret.resize(len);
                platon_get_call_output(ret.data());
                std::string str = toHex(ret);
                DEBUG("cross call contract cross_call_ppos_query success", "ret", str);
                return str;
            }
            }


            DEBUG("cross call contract cross_call_ppos_query fail", "address", target_addr);
            return "";
        }

};

PLATON_DISPATCH(call_ppos, (init)(cross_call_ppos_send)(cross_call_ppos_query))