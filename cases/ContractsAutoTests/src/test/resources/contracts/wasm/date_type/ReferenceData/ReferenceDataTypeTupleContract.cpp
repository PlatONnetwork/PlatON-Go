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
 * 定义tuple类型、初始化、取值
 *
 * */
extern char const tuple_a[] = "tuple_a";
extern char const tuple_b[] = "tuple_b";

CONTRACT ReferenceDataTypeTupleContract : public platon::Contract{

    private:
       platon::StorageType<tuple_a,std::tuple<std::string,uint8_t,std::vector<string>>> tuple_a;
       platon::StorageType<tuple_b,std::tuple<bool,std::string,uint8_t>> tuple_b;

    public:
        ACTION void init(){}
         /**
         * 1、定义类型
         *    赋值/取值
         **/

         //1)、初始化赋值
        ACTION void setInitTuple(){
            tuple_a.self() = {"test",2,{"1","2","3"}};
        }
        //2)、生成tuple对象，使用make_tuple函数
        ACTION void setTupleObject(){
            auto tupleObj = make_tuple(true,"1",1);//此对象类似于std::tuple<bool,std::string,uint8_t>
            tuple_b.self() = tupleObj;
        }
        //3)、tuple根据索引取值
        CONST std::string getTupleValueIndex(){
            return get<0>(tuple_a.self());
        }
};

PLATON_DISPATCH(ReferenceDataTypeTupleContract, (init)(setInitTuple)(setTupleObject)(getTupleValueIndex))
