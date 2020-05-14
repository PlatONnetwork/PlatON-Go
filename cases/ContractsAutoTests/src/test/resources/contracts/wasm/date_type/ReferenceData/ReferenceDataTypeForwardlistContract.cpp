#define TESTNET
#include <platon/platon.hpp>
#include <string>
#include <forward_list>
using namespace std;
using namespace platon;


/**
 * @author qudong
 * 合约引用类型(std::forward_list类型)：单向链表容器,不支持快速随机访问，
 * 不提供 size() 方法的容器。当不需要双向迭代时，具有比 std::list 更高的空间利用率。
 * 测试验证功能点：
 * 1、定义单链表容器类型
 *    1)、forward_list类型初始化赋值
 *
 * */

CONTRACT ReferenceDataTypeForwardlistContract : public platon::Contract{
     private:
       platon::StorageType<"forwardlist"_n, std::forward_list<uint8_t>> storage_int_forward_list;
    public:
    ACTION void init(){}

     //1)、forwardlist类型初始化赋值
    ACTION void init_forward_list(){
          storage_int_forward_list.self() = {1,2,3,4,5};
    }
    //2)、forwardlist创建，编译异常，暂不支持此类型
    ACTION void create_forward_list(){
         //1、创建空元素容器
         std::forward_list<uint8_t> values1;
         //2、创建一个包含10个元素的 容器
         std::forward_list<uint8_t> values2(10);
         //3、创建一个包含10个元素，并初始化值
         std::forward_list<uint8_t> values3(10, 5);
    }
};
PLATON_DISPATCH(ReferenceDataTypeForwardlistContract,(init)(init_forward_list)(create_forward_list))
