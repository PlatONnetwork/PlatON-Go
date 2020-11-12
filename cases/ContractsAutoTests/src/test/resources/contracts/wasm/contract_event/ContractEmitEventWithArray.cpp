#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;
using myintArray = std::array<int8_t, 8>;
using myuintArray = std::array<uint8_t, 8>;


CONTRACT ContractEmitEventWithArray : public platon::Contract{
   public:
      PLATON_EVENT0(transfer0,std::vector<int8_t>,std::vector<uint8_t>,FixedHash<8>,myintArray,myuintArray,std::list<uint8_t>,std::list<int8_t>,std::string)
      PLATON_EVENT1(transfer1,std::list<int8_t>,std::vector<int8_t>,std::vector<uint8_t>,FixedHash<8>,myintArray,myuintArray,std::list<uint8_t>,std::string)
      PLATON_EVENT2(transfer2,std::vector<int8_t>,std::vector<uint8_t>,FixedHash<8>,myintArray,myuintArray,std::list<int8_t>,std::list<uint8_t>,std::string)
      PLATON_EVENT3(transfer3,FixedHash<8>,myintArray,myuintArray, std::vector<int8_t>,std::vector<uint8_t>,std::list<int8_t>,std::list<std::string>,std::string)

      ACTION void init(){}

      ACTION void zerotopic_eigthargs_event(std::vector<int8_t> _argsOne,std::vector<uint8_t> _argsTwo,FixedHash<8> _argsThree,myintArray _argsFour,myuintArray _argsFive,std::list<uint8_t> _argsSix,std::list<int8_t> _argsSeven,std::string _argsEight){
          stringstorage.self() = _argsEight;
          PLATON_EMIT_EVENT0(transfer0,_argsOne,_argsTwo,_argsThree,_argsFour,_argsFive,_argsSix,_argsSeven,_argsEight);
      }

      ACTION void onetopic_sevenargs_event(std::list<int8_t> _topicOne,std::vector<int8_t> _argsOne,std::vector<uint8_t> _argsTwo,FixedHash<8> _argsThree,myintArray _argsFour,myuintArray _argsFive,std::list<uint8_t> _argsSix,std::string _argsSeven){
          stringstorage.self() = _argsSeven;
          PLATON_EMIT_EVENT1(transfer1, _topicOne, _argsOne,_argsTwo,_argsThree,_argsFour,_argsFive,_argsSix,_argsSeven);
      }

      ACTION void twotopic_sixargs_event(std::vector<int8_t> _topicOne,std::vector<uint8_t> _topicTwo,FixedHash<8> _argsOne,myintArray _argsTwo,myuintArray _argsThree,std::list<int8_t> _argsFour,std::list<uint8_t> _argsFive,std::string _argsSix){
          stringstorage.self() = _argsSix;
          PLATON_EMIT_EVENT2(transfer2, _topicOne,_topicTwo, _argsOne,_argsTwo,_argsThree,_argsFour,_argsFive,_argsSix);
      }

      ACTION void threetopic_fiveargs_event(FixedHash<8> _topicOne,myintArray _topicTwo,myuintArray _topicThree, std::vector<int8_t> _argsOne,std::vector<uint8_t> _argsTwo,std::list<int8_t> _argsThree,std::list<std::string> _argsFour,std::string _argsFive){
          stringstorage.self() = _argsFive;
          PLATON_EMIT_EVENT3(transfer3, _topicOne,_topicTwo,_topicThree, _argsOne,_argsTwo,_argsThree,_argsFour,_argsFive);
      }

      ACTION void set_string(std::string _stringargs){
         stringstorage.self() = _stringargs;
      }

      CONST std::string get_string(){
          return stringstorage.self();
      }
   private:
      platon::StorageType<"sstorage"_n, std::string> stringstorage;
};

PLATON_DISPATCH(ContractEmitEventWithArray, (init)(zerotopic_eigthargs_event)(onetopic_sevenargs_event)(twotopic_sixargs_event)(threetopic_fiveargs_event)(set_string)(get_string))