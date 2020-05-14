#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;



CONTRACT cross_call_storage_str : public platon::Contract {
    public:

        ACTION void init(){}

        ACTION uint64_t call_set_string(const std::string &target_address, const std::string &name,
          uint64_t value, uint64_t gas) {

            DEBUG("Call contract start", "address", target_address, "name", name);
            platon::bytes params = platon::cross_call_args("set_string", name);

            if (platon_call(Address(target_address), params, value, gas)) {
                 DEBUG("Call contract success", "address", target_address);
             } else {
                 DEBUG("Call contract fail", "address", target_address);
             }
            return 0;
        }

       CONST const std::string get_string(){

          DEBUG("cross_call_storage_str get_string", "name:", str.self());
          return str.self();
      }


    private:
       platon::StorageType<"str"_n, std::string> str;
};

PLATON_DISPATCH(cross_call_storage_str, (init)(call_set_string)(get_string))