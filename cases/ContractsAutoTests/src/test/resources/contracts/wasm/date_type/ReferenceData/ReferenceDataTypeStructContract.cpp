#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型结构体(struct类型)：允许存储不同类型的数据项

 * 测试验证功能点：
 * 1、定义struct类型
 *    赋值、取值
 *
 * */

struct person{
    public:
       std::string name;
       uint64_t age;
       person(){}
       person(const std::string &my_name,uint64_t &my_age):name(my_name),age(my_age){}
       PLATON_SERIALIZE(person,(name)(age))
};
extern char const struct_person[] = "struct_person";
extern char const string_name[] = "string_name";
extern char const uint64_age[] = "uint64_age";

CONTRACT structContractTest : public platon::Contract{

    private:
       platon::StorageType<struct_person,person> struct_person;
       platon::StorageType<string_name,std::string> string_name;
       platon::StorageType<uint64_age,uint64_t> uint64_age;

    public:
        ACTION void init(){}

         /**
         * 1、定义struct类型
         *    赋值/取值
         **/

         //1)、赋值
        ACTION void setStructPersonA(const std::string &my_name,uint64_t &my_age){
             struct_person.self() = person(my_name,my_age);
        }

         ACTION void setStructPersonB(){
              string_name.self() = "张三";
              uint64_age.self() = 20;
              struct_person.self() = person(string_name.self(),uint64_age.self());
          }

        //2)、取值
        CONST std::string getPersonName(){
             return struct_person.self().name;
        }
};

PLATON_DISPATCH(structContractTest, (init)(setStructPersonA)(setStructPersonB)(getPersonName))
