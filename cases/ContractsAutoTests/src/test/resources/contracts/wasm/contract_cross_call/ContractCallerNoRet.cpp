#undef NDEBUG
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;



CONTRACT cross_caller_noret : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION uint8_t callFeed(std::string target_address, uint64_t gasValue) {

            uint64_t transfer_value = 0;

            platon::bytes params = platon::cross_call_args("info");

            if (platon_call(Address(target_address), params, transfer_value, gasValue)) {

                 return 1; // successed
             }
             return 0; // failed
        }
       
};

PLATON_DISPATCH(cross_caller_noret, (init)(callFeed))