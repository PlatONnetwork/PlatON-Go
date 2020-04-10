#undef NDEBUG
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;



CONTRACT cross_caller_byret : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION void callFeed(std::string target_address, uint64_t gasValue) {

            uint64_t transfer_value = 0;
            auto result = platon::platon_call<uint8_t>(Address(target_address), transfer_value, gasValue, "info");
            if(result.second){
                status = 0; // successed

                DEBUG("cross_caller_byret call receiver_byret info has successed!")
            } else {
                status = 1; //failed

                DEBUG("cross_caller_byret call receiver_byret info has failed!")
            }

        }

        CONST uint64_t get_status(){
           return  status;
        }

    private:
       uint64_t status = 0;
};

PLATON_DISPATCH(cross_caller_byret, (init)(callFeed)(get_status))