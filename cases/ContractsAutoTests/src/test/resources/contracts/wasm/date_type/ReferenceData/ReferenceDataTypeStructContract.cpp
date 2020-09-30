#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型结构体(struct类型)：允许存储不同类型的数据项

 * 测试验证功能点：
 * 1、定义struct类型包含基本类型（赋值、取值）
 * */

struct person{
    public:
       std::string name;
       uint64_t age;
       person(){}
       person(const std::string &my_name,uint64_t &my_age):name(my_name),age(my_age){}
       PLATON_SERIALIZE(person,(name)(age))
};

CONTRACT ReferenceDataTypeStructContract : public platon::Contract{

    private:
       platon::StorageType<"person"_n,person> storage_struct_person;
       platon::StorageType<"name"_n,std::string> storage_string_name;
       platon::StorageType<"age"_n,uint64_t> storage_uint64_age;
     //  platon::StorageType<"storage_struct_group"_n,group> storage_struct_group;

    public:
        ACTION void init(){}
         /**
         * 1、定义struct类型包含基本类型（赋值、取值）
         *    赋值/取值
         **/
         //1)、赋值
        ACTION void setStructPersonA(const std::string &my_name,uint64_t &my_age){
             storage_struct_person.self() = person(my_name,my_age);
        }
         ACTION void setStructPersonB(){
              storage_string_name.self() = "张三";
              storage_uint64_age.self() = 20;
              storage_struct_person.self() = person(storage_string_name.self(),storage_uint64_age.self());
          }
        //2)、取值
        CONST std::string getPersonName(){
             return storage_struct_person.self().name;
        }

};

PLATON_DISPATCH(ReferenceDataTypeStructContract, (init)(setStructPersonA)(setStructPersonB)
               (getPersonName))

