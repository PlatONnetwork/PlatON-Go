#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT hello : public platon::Contract{
   public:
      PLATON_EVENT1(myevent, std::string, std::string)

     ACTION void init(){}
      
      ACTION void set_string(const std::string &name){
          PLATON_EMIT_EVENT1(myevent, "set_string", name);
          DEBUG("hello set_string", "name:", name);
          str.self() = name;
      }
      CONST const std::string get_string(){
          PLATON_EMIT_EVENT1(myevent, "get_string", str.self());
          DEBUG("hello get_string", "name:", str.self());
          return str.self();
      }

      


   private:
      platon::StorageType<"str"_n, std::string> str;
};

PLATON_DISPATCH(hello, (init)(set_string)(get_string))
