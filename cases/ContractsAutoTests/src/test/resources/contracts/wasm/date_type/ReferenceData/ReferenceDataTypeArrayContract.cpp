#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型数组(array类型)：数组可以容纳相同类型的数据表
 * 测试验证功能点：
 * 1、定义array类型
 *    1)、基本类型数组
 *    2)、自定义类型数组----异常，待与开发沟通
 * 2、字节数组bytes
 *
 * */

/*class Person{
    public:
       std::string name;
       uint64_t age;
       Person(){}
       Person(const std::string &my_name,uint64_t &my_age):name(my_name),age(my_age){}
       PLATON_SERIALIZE(Person,(name)(age))
};*/
CONTRACT ReferenceDataTypeArrayContract : public platon::Contract{

    private:
       platon::StorageType<"a"_n,std::array<std::string,10>> storage_array_string;
       platon::StorageType<"b"_n,std::array<uint8_t,10>> storage_array_uint8;
       platon::StorageType<"c"_n,bytes> storage_array_bytes;
      //platon::StorageType<"storage_array_peron"_n,std::array<Person,5>> storage_array_peron;
      //platon::StorageType<"storage_array_bool"_n,std::array<bool,5>> storage_array_bool;


    public:
        ACTION void init(){}
         /**
         * 1、定义array类型
         *    赋值/取值
         **/
         //1)、验证数组赋值
         ACTION void setInitArray(){
            //赋值方式一
            storage_array_uint8.self() = {1,2,3,4,5};
            //赋值方式二
            storage_array_string.self()[0] = "a";
            storage_array_string.self()[1] = "b";
            storage_array_string.self()[2] = "c";
        }

        //2)、验证数组取值
         CONST std::string getArrayStringIndex(const uint32_t &index){
            return storage_array_string.self()[index];
         }

        //3)、验证数组大小
         CONST uint8_t getArrayUintSize(){
             return storage_array_uint8.self().size();
         }

        //4)、验证数组定义person类型
      /* ACTION void setArrayPerson(const Person &person){
              storage_array_peron.self()[0] = person;
        }
        CONST std::string getArrayPersonNameIndex(){
              return storage_array_peron.self()[0].name;
        }*/

        //5)、定长array:验证定长数组赋值超出存储空间、赋值错误类型值
      /*   ACTION void setArrayOver(){
               storage_array_bool.self() = {true,false};//正常
               //storage_array_bool.self() = {1,2,3};//异常，赋值错误类型值
               //storage_array_bool.self() = {true,false,true,false,true,false,true,false};//异常，赋值超过存储空间
         }*/

      /**
        * 2、字节数组bytes
        **/
         ACTION void setBytesArray(){
             storage_array_bytes.self() = {1,2,3,4,5};
         }

         CONST uint8_t getBytesArrayIndex(const uint32_t &index){
              return storage_array_bytes.self()[index];
         }
};

PLATON_DISPATCH(ReferenceDataTypeArrayContract,(init)(setInitArray)(getArrayStringIndex)(getArrayUintSize)
               (setBytesArray)(getBytesArrayIndex))

