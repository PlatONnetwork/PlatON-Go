#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型链表：是由一系列连接在一起的结点构成，其中每个结点是一个数据结构。
 * 其中，结点包含一个或多个数据成员，除数据之外，每个结点还包含一个后续指针指向链表的下一个结点，
 * 以此类推，构成链表。
 * 【说明】：链表的结点通常是动态分配、使用、删除的，如果需要将新的信息添加到链表中，
 *  则只需分配另一个结点并将其插入到系列中。需要删除特定的信息，则只需删除包含信息的结点。
 *
 * 测试验证功能点：
 * 1、定义链表类型
 * 2、链表取值
 * */

 struct listNode{
     public:
        std::string name;
        listNode *nextPointer;
        listNode(){}
        listNode(const std::string &my_name,listNode *my_nextPointer):name(my_name),nextPointer(my_nextPointer){}
        PLATON_SERIALIZE(listNode, (name))
 };


CONTRACT linkedListContractTest : public platon::Contract{

    private:
       platon::StorageType<"storage_listnode_a"_n,listNode> storage_listnode_a;
       platon::StorageType<"storage_listnode_b"_n,listNode> storage_listnode_b;
       platon::StorageType<"storage_listnode_c"_n,listNode> storage_listnode_c;
       platon::StorageType<"storage_listnode_d"_n,listNode> storage_listnode_d;
       //platon::StorageType<"storage_listnode_vector"_n, std::vector<listNode>> listnode_vector;

    public:
        ACTION void init(){}

         /**
         * 1、定义链表类型
         **/
         //1)、定义单向链表
        ACTION void setListNode(){
            //第一个结点
            storage_listnode_a.self().name = "a";
            //第二个结点
            storage_listnode_b.self().name = "b";
            storage_listnode_b.self().nextPointer = &storage_listnode_a.self();
            //第三个结点
            storage_listnode_c.self().name = "c";
            storage_listnode_c.self().nextPointer = &storage_listnode_b.self();
        }
        //2)、增加节点
     /*   ACTION void addListNodeVector(const std::string &my_name,listNode *my_nextPointer)){
            listnode_vector.self().push_back(listNode(my_name,my_nextPointer));
        }*/

        ACTION void addListNode(){
            for(int i = 0; i < 5; i++){
                storage_listnode_d.self().name = "Lucy";
                storage_listnode_d.self().nextPointer = &storage_listnode_d.self();
            }
        }

        CONST uint64_t getListNode(){
           //nextPointer：是指针类型
            return (uint64_t)storage_listnode_d.self().nextPointer;
        }
};

PLATON_DISPATCH(linkedListContractTest,(init)(setListNode)(addListNode)(getListNode))
