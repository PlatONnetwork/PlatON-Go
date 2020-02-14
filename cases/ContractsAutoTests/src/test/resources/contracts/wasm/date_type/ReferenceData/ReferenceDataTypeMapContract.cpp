#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型(map类型)：map是一个key-value一对一存储的容器
 * 测试验证功能点：
 * 1、定义map类型
 *    key关键字只能在map出现一次，value值可以多次
 *    map中的key与value可以是任意类型
 * 2、map属性方法
 *     增加、删除、容器集合大小
 * */

CONTRACT mapContractTest : public platon::Contract{


    private:
      // platon::StorageType<map_int,std::map<int8_t,std::string>> map_int;
       platon::StorageType<"map_uint8"_n,std::map<uint8_t,std::string>> map_uint;
       platon::StorageType<"map_bool"_n,std::map<bool,std::string>> map_bool;
       platon::StorageType<"map_string"_n,std::map<std::string,std::string>> map_string;
       platon::StorageType<"map_address"_n,std::map<Address,std::string>> map_address;

    public:
        ACTION void init(){}

         /**
         * 1、定义map类型
         *    1)、map中的key与value可以是任意类型
         *    2)、key关键字只能在map出现一次，value值可以多次
         **/

         //1)、验证map中的key与value可以是任意类型
        ACTION void setMap(){
            map_uint.self()[0] = "test1";//正常
            map_bool.self()[true] = "test2";//正常
            map_string.self()["1"] = "test3";//正常
            //key为Address类型
            Address address = platon_caller();//获取交易发起者地址
            map_address.self()[address] = "test4";//正常
        }
        //2)、key关键字只能在map出现一次，value值可以多次
        ACTION void setSameKeyMap(){
              map_uint.self()[1] = "test1";
              map_uint.self()[2] = "test2";
              map_uint.self()[2] = "test3";
        }

        /**
         *2、map属性方法
         *   增加、删除、容器集合大小
         **/
         //验证新增
         ACTION void addMap(){
               map_string.self()["a"] = "a";
               map_string.self()["b"] = "b";
          }
          //验证map容器大小
          CONST uint8_t getMapSize(){
               return map_string.self().size();
          }
          //验证删除
         /*ACTION void deleteMap(){

          }*/

};

PLATON_DISPATCH(mapContractTest, (init)(setMap)(setSameKeyMap)(addMap)(getMapSize))
