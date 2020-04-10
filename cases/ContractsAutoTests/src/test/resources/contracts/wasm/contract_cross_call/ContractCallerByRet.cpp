#undef NDEBUG
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;



CONTRACT cross_caller_byret : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION uint8_t callFeed(std::string target_address, uint64_t gasValue) {

            uint64_t transfer_value = 0;
            auto result = platon::platon_call<uint8_t>(Address(target_address), transfer_value, gasValue, "info");
            if(!result.secod){
               return 0; // successed
            }

            return 1; //failed
        }
       
};

PLATON_DISPATCH(cross_caller_byret, (init)(callFeed))