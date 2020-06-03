#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>
#undef NDEBUG

using namespace platon;

CONTRACT call_ppos : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION uint64_t cross_call_ppos_send (std::string target_addr, std::string &in, uint64_t value, uint64_t gas) {
            platon::bytes input = fromHex(in);

            auto address_info = make_address(target_addr);

            printf("address_info.first:%s\t\n", address_info.first.toString().c_str());
            printf("address_info.second:%s\t\n", address_info.second ? "true" : "false");

            if(address_info.second){
                if (platon_call(address_info.first, input, value, gas)) {
                    platon_call(address_info.first, input, value, gas);
                    platon_call(address_info.first, input, value, gas);
                    DEBUG("cross call contract cross_call_ppos_send success", "address", target_addr);
                    return 0;
                }
            }
            DEBUG("cross call contract cross_call_ppos_send fail", "address", target_addr);
            return 1;
        }

        CONST const std::string  cross_call_ppos_query (std::string target_addr, std::string &in, uint64_t value, uint64_t gas) {
            platon::bytes  input = fromHex(in);

            auto address_info = make_address(target_addr);
            if(address_info.second){
                if (platon_call(address_info.first, input, value, gas)) {
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