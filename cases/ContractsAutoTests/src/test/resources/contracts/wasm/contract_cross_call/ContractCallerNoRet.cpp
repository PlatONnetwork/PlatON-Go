#undef NDEBUG
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;



CONTRACT cross_caller_noret : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION void callFeed(std::string target_address, uint64_t gasValue) {

            uint64_t transfer_value = 0;

            platon::bytes params = platon::cross_call_args("info");

            if (platon_call(Address(target_address), params, transfer_value, gasValue)) {
                 status = 0; // successed

                 DEBUG("cross_caller_noret call receiver_noret info has successed!")
             } else {
                 status = 1; //failed

                 DEBUG("cross_caller_noret call receiver_noret info has failed!")
             }

        }
       CONST uint64_t get_status(){
          return  status;
       }

       private:
           uint64_t status = 0;

};

PLATON_DISPATCH(cross_caller_noret, (init)(callFeed)(get_status))