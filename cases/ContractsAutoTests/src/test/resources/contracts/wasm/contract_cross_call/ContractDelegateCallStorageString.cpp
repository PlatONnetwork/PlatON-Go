#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;


CONTRACT delegate_call_storage_str : public platon::Contract {
     public:

     ACTION void init(){}

     ACTION uint64_t delegate_call_set_string(const std::string &target_address, std::string &name, uint64_t gas) {
         
         DEBUG("Delegate call contract start", "address", target_address, "name", name);
         platon::bytes params = platon::cross_call_args("set_string", name);

         auto address_info = make_address(target_address);
         if(address_info.second){
             if (platon_delegate_call(address_info.first, params, gas)) {
                  DEBUG("Delegate call contract success", "address", target_address, "name", name);
             } else {
                  DEBUG("Delegate call contract fail", "address", target_address, "name", name);
             }
         }

         return 0;
     }

      CONST const std::string get_string(){
          
          platon::StorageType<"str"_n, std::string> str; // Must use local definitions for manipulating the corresponding keys in th
          return str.self();
      }



};

PLATON_DISPATCH(delegate_call_storage_str, (init)(delegate_call_set_string)(get_string))