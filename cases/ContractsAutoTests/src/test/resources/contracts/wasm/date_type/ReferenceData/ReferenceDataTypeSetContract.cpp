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
CONTRACT ReferenceDataTypeSetContract : public platon::Contract{

     private:
       platon::StorageType<"intmap"_n, std::set<uint8_t>> storage_int_set;

    public:
    ACTION void init(){}

    /**
     *    set类型:无顺序不可重复容器
     *    1)、set类型初始化赋值
     *    2)、set插入元素
     *    3)、set类型存储值唯一性
     *    4)、访问元素
     *    5)、查找元素
     *    6)、删除元素
     *    7)、判断元素
     *    8)、清空元素
     *    9)、set容器元素数量
     **/

     // 1)、set类型初始化赋值
    ACTION void init_set(){
          storage_int_set.self() = {1,2,3,4,5,6};
    }
    //2)、set插入元素
    ACTION void insert_set(const uint8_t &value){
          storage_int_set.self().insert(value);
    }
    //4)、访问元素
    ACTION void iterator_set(){
        for (auto iter = storage_int_set.self().begin(); iter != storage_int_set.self().end(); iter++){
            DEBUG("ReferenceDataTypeSetContract", "setIterator", *iter);
        }
    }
    //5)、查找元素
    CONST uint8_t find_set(){
          uint8_t v = 0;
          auto iter = storage_int_set.self().find(3);
          if(iter != storage_int_set.self().end()){
              v = *iter;
          }
          return v;
     }
    //6)、删除元素
     ACTION void erase_set(const uint8_t &value){
            storage_int_set.self().erase(value);
     }
    //7)、判断元素
     CONST bool get_set_empty(){
           return storage_int_set.self().empty();
     }
   //8)、清空元素
   ACTION void clear_set(){
         storage_int_set.self().clear();
   }
   //9)、set容器元素数量
    CONST uint64_t get_set_size(){
        return storage_int_set.self().size();
    }

};

PLATON_DISPATCH(ReferenceDataTypeSetContract,(init)(init_set)(insert_set)(get_set_size)
(iterator_set)(find_set)(erase_set)(get_set_empty)(clear_set)(get_set_size))
