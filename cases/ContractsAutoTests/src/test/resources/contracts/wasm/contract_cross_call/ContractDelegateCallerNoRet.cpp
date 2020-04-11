#undef NDEBUG
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;



CONTRACT delegate_caller_noret : public platon::Contract {
    public:
        ACTION void init(){}

        ACTION void callFeed(std::string target_address, uint64_t gasValue) {

            uint64_t transfer_value = 0;

            platon::bytes params = platon::cross_call_args("info");

            if (platon_call(Address(target_address), params, transfer_value, gasValue)) {
                 status.self() = 0; // successed

                 DEBUG("delegate_caller_noret call receiver_noret info has successed!")
             } else {
                 status.self() = 1; //failed

                 DEBUG("delegate_caller_noret call receiver_noret info has failed!")
             }

        }
       CONST uint64_t get_status(){
          return  status.self();
       }

       private:
           platon::StorageType<"status"_n, uint64_t> status;

};

PLATON_DISPATCH(delegate_caller_noret, (init)(callFeed)(get_status))