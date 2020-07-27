#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;



/**
 * @author qudong
 * 合约引用类型(std::unordered_set类型)：基于哈希表无序不重复set容器,
 * */

CONTRACT ReferenceDataTypeUnorderSetContract : public platon::Contract{

     private:
       platon::StorageType<"set"_n, std::unordered_set<std::string>> storage_string_set;

    public:
       ACTION void init(){}

     // 1)、unorederd_set类型类型初始化赋值，编译异常，暂不支持此类型
    ACTION void init_unorder_set(){
          storage_string_set.self() = {"one","two","three"};
    }
    //2)、unorederd_set创建，编译异常，暂不支持此类型
    ACTION void create_unorder_set(){
         //1、创建空元素容器
         std::unordered_set<std::string> values1;
         //2、创建一个容器并赋值
         std::unordered_set<std::string> values2{"aaa","bbb","ccc"};
    }


};

PLATON_DISPATCH(ReferenceDataTypeUnorderSetContract,(init)
(init_unorder_set)(create_unorder_set)
)
