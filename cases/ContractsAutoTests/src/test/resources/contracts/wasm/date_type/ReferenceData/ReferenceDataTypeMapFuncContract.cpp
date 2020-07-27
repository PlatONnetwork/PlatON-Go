#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型(map类型)
 * 测试验证功能点：
 * 1、map属性方法
 *     增加、删除、容器集合大小
 * */

CONTRACT ReferenceDataTypeMapFuncContract : public platon::Contract{

    private:
       platon::StorageType<"map1"_n,std::map<uint8_t,std::string>> storage_map_uint;
       platon::StorageType<"map2"_n,std::map<std::string,std::string>> storage_map_string;


    public:
        ACTION void init(){}
        /**
         * 1、map属性方法
         *    增加、删除、容器集合大小
         **/
         //1)、验证map容器新增
         ACTION void addMapByUint(const uint8_t n){
               for(uint8_t i = 0; i <= n; i++){
                   storage_map_uint.self()[i] = std::to_string(i);
               }
          }
         //2)、验证map容器大小
         CONST uint8_t getMapBySize(){
               return storage_map_uint.self().size();
         }
         //3)、验证map容器删除
         ACTION void deleteMapByIndex(const uint8_t &key){
                storage_map_uint.self().erase(key);
          }
         //4)、验证map容器插入方法insert()
         ACTION void insertMapUint(const uint8_t &key,const std::string &value){
                storage_map_uint.self().insert(pair<uint8_t,std::string>(key,value));
         }
         //5)、验证map容器清空集合clear()
          ACTION void clearMapUint(){
                storage_map_uint.self().clear();
          }
         //6)、验证map容器判断空方法empty()
          CONST bool getMapIsEmpty(){
               return storage_map_uint.self().empty();
           }

         //7)、验证map容器迭代器
         /* CONST std::map<uint8_t,std::string> getMapUintIterator(){
              std::map<uint8_t,std::string> ::iterator iter = storage_map_uint.self().begin();
              while(iter != storage_map_uint.self().end()) {
                      a_storage_map_uint.self()[iter->first] = iter->second;
                      iter++;
                  }
              return a_storage_map_uint.self();
          }*/



};

PLATON_DISPATCH(ReferenceDataTypeMapFuncContract, (init)(addMapByUint)(getMapBySize)(deleteMapByIndex)
               (insertMapUint)(clearMapUint)(getMapIsEmpty))
