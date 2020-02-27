#include <platon/platon.hpp>
#include <string>
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

   private:
      platon::StorageType<"strmap"_n, std::map<std::string,std::string>> strmap;
      platon::StorageType<"mapmap"_n, std::map<std::string,std::map<std::string,std::string>>> mapmap;
      platon::StorageType<"maplist"_n, std::map<std::string,std::list<std::string>>> maplist;
};

PLATON_DISPATCH(InitWithParams, (init)(add_map)(get_map)(add_map_map)(get_map_map)(add_map_list)(get_map_list))
