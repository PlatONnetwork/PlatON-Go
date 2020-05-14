#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;


CONTRACT InitWithParams : public platon::Contract{
   public:
      ACTION void init(const std::map<std::string,std::string>  &inMap){
         strmap.self() = inMap;
      }

      ACTION void add_map(const std::map<std::string,std::string>  &inMap){
          strmap.self() = inMap;
      }
      CONST std::map<std::string,std::string> get_map(){
          return strmap.self();
      }

      ACTION void add_map_map(const std::map<std::string,std::map<std::string,std::string>>  &inMapmap){
          mapmap.self() = inMapmap;
      }
      CONST std::map<std::string,std::map<std::string,std::string>> get_map_map(){
          return mapmap.self();
      }

      ACTION void add_map_list(const std::map<std::string,std::list<std::string>>  &inMaplist){
          maplist.self() = inMaplist;
      }
      CONST std::map<std::string,std::list<std::string>> get_map_list(){
          return maplist.self();
      }

      //map add element
      ACTION void add_map_element(std::string &key,std::string &value){
          strmap.self().insert(pair<std::string, std::string>(key,value));
      }

      //map delete element
      ACTION void delete_map_element(std::string &key){
          strmap.self().erase(key);
      }

      //map find value by key
      CONST std::string find_element_bykey(std::string &key){
          if(strmap.self().count(key)>0){
                return strmap.self()[key];
          }
          return "";
      }

      //map size
      CONST uint8_t get_map_size(){
            return strmap.self().size();
      }

      //map add map
      ACTION void addMap(const std::map<std::string,std::string>& inMap) {
        for (auto iter=inMap.begin(); iter!=inMap.end(); iter++) {
        	DEBUG("InitWithParams", "inMap", iter->first, iter->second);
        	if(strmap.self().count(iter->first)>0){
        	    continue;
        	}else{
        	    strmap.self()[iter->first] = iter->second;
        	}
        }
      }


   private:
      platon::StorageType<"strmap"_n, std::map<std::string,std::string>> strmap;
      platon::StorageType<"mapmap"_n, std::map<std::string,std::map<std::string,std::string>>> mapmap;
      platon::StorageType<"maplist"_n, std::map<std::string,std::list<std::string>>> maplist;
};

PLATON_DISPATCH(InitWithParams, (init)(add_map)(get_map)(add_map_map)(get_map_map)(add_map_list)(get_map_list)(add_map_element)(delete_map_element)(find_element_bykey)(get_map_size)(addMap))
