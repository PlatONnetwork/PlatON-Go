#define TESTNET
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
            auto address_info = make_address(target_address);
            if(address_info.second){
                auto result = platon::platon_call<uint8_t>(address_info.first, transfer_value, gasValue, "info");
            if(result.second){
                status.self() = 0; // successed

                DEBUG("cross_caller_byret call receiver_byret info has successed!")
            } else {
                status.self() = 1; //failed

                DEBUG("cross_caller_byret call receiver_byret info has failed!")
            }
            }
        }

        CONST uint64_t get_status(){
           return  status.self();
        }

    private:
       platon::StorageType<"status"_n, uint64_t> status;
};

PLATON_DISPATCH(cross_caller_byret, (init)(callFeed)(get_status))