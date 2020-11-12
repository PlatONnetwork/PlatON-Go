#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型tuple：可以将一些不同类型的数据和成单一对象，可以用于函数返回多个值。
 * tuple所有成员是public，一个tuple可以有任意数量的成员；
 * 一个确定的tuple类型的成员数目是固定的，不能进行添加和删除等操作。
 *
 * 测试验证功能点：
 * 1、定义tuple类型、初始化、取值
 * 2、定义包含引用类型
 * */

struct person{
    public:
       std::string name;
       uint64_t age;
       person(){}
       person(const std::string &my_name,uint64_t &my_age):name(my_name),age(my_age){}
       PLATON_SERIALIZE(person,(name)(age))
};

CONTRACT ReferenceDataTypeTupleContract : public platon::Contract{

    private:
       platon::StorageType<"tuple1"_n,std::tuple<std::string,uint8_t,std::vector<string>>> storage_tuple_one;
       platon::StorageType<"tuple2"_n,std::tuple<std::string,uint8_t>> storage_tuple_two;
       platon::StorageType<"tuple3"_n,std::tuple<std::string,person,std::array<std::string,10>>> storage_tuple_three;
       platon::StorageType<"struct4"_n,person> storage_struct_person;
    public:
        ACTION void init(){}
         /**
         * 1、定义类型
         *    赋值/取值
         **/
         //1)、元组初始化赋值方式一
        ACTION void setInitTupleModeOne(){
            storage_tuple_one.self() = {"Lucy",2,{"1","2","3"}};
        }
        //2)、tuple根据索引取值
        CONST std::string getTupleValueIndex1(){
            return get<0>(storage_tuple_one.self());
        }
        CONST uint8_t getTupleValueIndex2(){
            return get<1>(storage_tuple_one.self());
        }

         //3)、元组初始化赋值方式二(使用make_tuple函数)
        ACTION void setInitTupleModeTwo(const std::string &a,const uint8_t &b){
            storage_tuple_two.self() = make_tuple(a,b);
        }
        CONST std::string getTupleValueIndex3(){
            return get<0>(storage_tuple_two.self());
        }

        //4)、定义包含引用类型
        ACTION void setInitTupleModeThree(const std::string &name,const uint64_t &age){
           std::string str = "Lili";
           std::array<std::string,10> array;
           array[0] = "a";
           array[1] = "b";
           array[2] = "c";
           storage_struct_person.self().name = name;
           storage_struct_person.self().age = age;
           storage_tuple_three.self() = make_tuple(str,storage_struct_person.self(),array);
        }
         CONST person getTupleValueIndex4(){
            return get<1>(storage_tuple_three.self());
         }




};

PLATON_DISPATCH(ReferenceDataTypeTupleContract,(init)(setInitTupleModeOne)(getTupleValueIndex1)
(getTupleValueIndex2)(setInitTupleModeTwo)(getTupleValueIndex3)(setInitTupleModeThree)(getTupleValueIndex4))
