#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

CONTRACT ContractEmitEvent1ComplexParam : public platon::Contract{
   public:
      PLATON_EVENT1(transfer,std::string,uint32_t,std::list<std::string>)

      ACTION void init(){}
 
      ACTION void one_emit_event1(std::string name,uint32_t value,const std::list<std::string>  &inList){
           stringstorage.self() = name;
           PLATON_EMIT_EVENT1(transfer,name,value,inList);
      }


      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
      platon::StorageType<"listvar"_n, std::list<std::string>> sList;
};

PLATON_DISPATCH(ContractEmitEvent1ComplexParam, (init)(one_emit_event1))
