#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;


CONTRACT cross_delegate_call : public platon::Contract {
     public:
      PLATON_EVENT1(myevent, std::string, std::string)

     ACTION void init(){}
     ACTION uint64_t delegate_call_add_message(const std::string &target_address,
     std::string &name, uint64_t gas) {
         PLATON_EMIT_EVENT1(myevent, "delegate_call_add_message", name);
         platon::bytes params = platon::cross_call_args("set_string", name);
         if (platon_delegate_call(Address(target_address), params, gas)) {
              DEBUG("Delegate call contract success", "address", target_address, "name", name);
         } else {
             DEBUG("Delegate call contract fail", "address", target_address, "name", name);
         }
         return 0;
     }

      CONST const std::string get_string(){
          PLATON_EMIT_EVENT1(myevent, "get_string", str.self());
          DEBUG("cross_delegate_call get_string", "name:", str.self());
          return str.self();
      }

   private:
      platon::StorageType<"str"_n, std::string> str;

};

PLATON_DISPATCH(cross_delegate_call, (init)(get_string))