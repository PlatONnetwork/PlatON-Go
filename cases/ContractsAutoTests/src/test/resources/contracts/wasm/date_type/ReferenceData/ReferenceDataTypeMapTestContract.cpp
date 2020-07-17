#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

class Person{
    public:
       std::string name;
       uint64_t age;
       Person(){}
       Person(const std::string &my_name,uint64_t &my_age):name(my_name),age(my_age){}
       PLATON_SERIALIZE(Person,(name)(age))
};
CONTRACT ReferenceDataTypeMapTestContract : public platon::Contract{

     private:
       platon::StorageType<"stringmap"_n, std::map<std::string,std::string>> storage_string_map;
       platon::StorageType<"mapuint"_n,std::map<uint8_t,std::string>> storage_map_uint;
       platon::StorageType<"amapuint"_n,std::map<uint8_t,std::string>> a_storage_map_uint;
       platon::StorageType<"mapbool"_n,std::map<bool,std::string>> storage_map_bool;
       platon::StorageType<"mapstring"_n,std::map<std::string,std::string>> storage_map_string;
       platon::StorageType<"amapstring"_n, std::map<std::string,std::string>> a_storage_map_string;
       platon::StorageType<"mapperson"_n,std::map<uint8_t,Person>> storage_map_person;
       platon::StorageType<"mapaddress"_n,std::map<Address,std::string>> storage_map_address;

    public:
    ACTION void init(){}

    /**
     * 1、定义map类型
     *    1)、map中的key与value可以是任意类型
     *    2)、map类型赋值&取值
     **/

     //1)、验证map中的key与value可以是任意类型
    ACTION void setMapKeyType(){
        storage_map_uint.self()[0] = "test1";//正常
        storage_map_bool.self()[true] = "test2";//正常
        storage_map_string.self()["1"] = "test3";//正常
        //key为Address类型
        Address address = platon_caller();//获取交易发起者地址
        storage_map_address.self()[address] = "test4";//正常
    }

    //2)、string类型map容器赋值&取值
    ACTION void addMapString(const std::string &one_key,const std::string &one_value){
             storage_string_map.self()[one_key]= one_value;
     }
    CONST uint64_t getMapStringSize(){
        return storage_string_map.self().size();
    }
    CONST std::string getMapValueByString(const std::string &key){
        return storage_string_map.self()[key];
    }

   //3)、person类型map容器赋值&取值
    ACTION void addMapByPerson(const uint8_t &key,const Person &person){
         storage_map_person.self()[key] = person;
    }
    CONST uint64_t getMapByPersonSize(){
          return storage_map_person.self().size();
    }
   CONST std::string getMapByPerson(uint8_t key){
         return storage_map_person.self()[key].name;
    }

};

PLATON_DISPATCH(ReferenceDataTypeMapTestContract, (init)(setMapKeyType)(addMapString)(getMapStringSize)(getMapValueByString)
               (addMapByPerson)(getMapByPersonSize)(getMapByPerson))
