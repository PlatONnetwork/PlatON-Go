#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT ContractEmitEvent1WithMap : public platon::Contract{
   public:
      PLATON_EVENT1(transfer,std::string,uint32_t,std::map<std::string,std::string>)

      ACTION void init(){}
 
      ACTION void one_emit_event1(std::string name,uint32_t value,const std::map<std::string,std::string>  &inMap){
           stringstorage.self() = name;
           for (auto iter=inMap.begin(); iter!=inMap.end(); iter++) {
              DEBUG("ContractEmitEvent1WithMap", "inMap", iter->first, iter->second);
           }
//           inMap.insert(std::pair<std::string, std::string>("key3", "value3"));
           PLATON_EMIT_EVENT1(transfer,name,value,inMap);
      }

      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
      platon::StorageType<"strmap"_n, std::map<std::string,std::string>> strmap;

};

PLATON_DISPATCH(ContractEmitEvent1WithMap, (init)(one_emit_event1))
