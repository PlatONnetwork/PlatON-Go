#undef NDEBUG
#define TESTNET
#include <platon/platon.hpp>

using namespace platon;

CONTRACT UpdateContract : public platon::Contract {
    public:
        CONST void init(){}

        CONST int32_t get_contract_length(const Address & contract_address){
            size_t code_length = platon_contract_code_length(contract_address.data());
            if(0 == code_length) return 0;
            bytes contract_code(code_length);
            int32_t result = platon_contract_code(contract_address.data(), contract_code.data(), code_length);
            return result;
        }

        ACTION void deploy_contract(const Address & contract_address){
            size_t code_length = platon_contract_code_length(contract_address.data());
            if(0 == code_length) return;
            bytes contract_code(code_length);
            int32_t result = platon_contract_code(contract_address.data(), contract_code.data(), code_length);
            if(result <= 0) return;
            auto info = platon_create_contract(contract_code, 0U, 0U);
            if(info.second) set_state("deploy", info.first);
        }

        CONST Address get_deploy_address(){
            Address result;
            get_state("deploy", result);
            return result;
        }

        ACTION void clone_contract(const Address & contract_address){
            auto info = platon_create_contract(contract_address,  0U, 0U);
            if(info.second) set_state("clone", info.first);
        }

        CONST Address get_clone_address(){
            Address result;
            get_state("clone", result);
            return result;
        }

        ACTION void set_simple_address(const Address& addr){
            set_state("simple", addr);
        }

        CONST Address get_simple_address(){
            Address result;
            get_state("simple", result);
            return result;
        }

        ACTION void migrate(const Address & contract_address){
            size_t code_length = platon_contract_code_length(contract_address.data());
            if(0 == code_length) return;
            bytes contract_code(code_length);
            int32_t get_result = platon_contract_code(contract_address.data(), contract_code.data(), code_length);
            if(get_result <= 0 ) return;

            // value and gas
            bytes value_bytes = value_to_bytes(0U);
            bytes gas_bytes = value_to_bytes(0U);

            // deploy agrs
            std::vector<byte> magic_number = {0x00, 0x61, 0x73, 0x6d};
            bytes init_rlp = cross_call_args("init");
            RLPSize rlps;
            rlps << contract_code << init_rlp;
            RLPStream stream;
            stream.appendPrefix(magic_number);
            stream.reserve(rlps.size() + 4);
            stream.appendList(2);
            stream << contract_code << init_rlp;
            bytesRef result = stream.out();

            // migrate
            Address return_address;
            platon_migrate(return_address.data(), result.data(), result.size(), value_bytes.data(), value_bytes.size(), gas_bytes.data(), gas_bytes.size());

            // // set new address
            Address simple_address = get_simple_address();
            platon_call(simple_address, 0U, 0U, "set_address", return_address);
        }

        ACTION void migrate_clone(const Address & contract_address){
            // value and gas
            bytes value_bytes = value_to_bytes(0U);
            bytes gas_bytes = value_to_bytes(0U);

            // init args
            bytes init_rlp = cross_call_args("init");

            // migrate
            Address return_address;
            platon_clone_migrate(contract_address.data(), return_address.data(), init_rlp.data(), init_rlp.size(), value_bytes.data(), value_bytes.size(), gas_bytes.data(), gas_bytes.size());

            // set new address
            Address simple_address = get_simple_address();
            platon_call(simple_address, 0U, 0U, "set_address", return_address);
        }
};

PLATON_DISPATCH(UpdateContract, (init)(get_contract_length)(deploy_contract)(get_deploy_address)(clone_contract)(get_clone_address)
    (set_simple_address)(get_simple_address)(migrate)(migrate_clone))