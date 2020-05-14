#define TESTNET
#include <platon/platon.hpp>
#include <string>

using namespace platon;



CONTRACT receiver_noret : public platon::Contract{
   public:

      ACTION void init(){}

  
      ACTION void info (){
            
           for (uint8_t i = 0; i < 2; i++) {
             // uintmap.self().insert({i, i});
            uintmap.self()[i] = i;
           }
      }

      CONST uint8_t get_value(const uint8_t key){

              auto iter = uintmap.self().find(key);
               
              if(iter != uintmap.self().end())
                  return iter->second;
              else
                  return 0;

      }

   private:
      platon::StorageType<"uintmap"_n,std::map<uint8_t, uint8_t>> uintmap;
};

PLATON_DISPATCH(receiver_noret, (init)(info)(get_value))
