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
extern char const listNode_a[] = "listNode_a";
extern char const listNode_b[] = "listNode_b";
extern char const listNode_c[] = "listNode_c";


CONTRACT linkedListContractTest : public platon::Contract{

    private:
       platon::StorageType<listNode_a,listNode> listNode_a;
       platon::StorageType<listNode_b,listNode> listNode_b;
       platon::StorageType<listNode_c,listNode> listNode_c;

    public:
        ACTION void init(){}
         /**
         * 1、定义链表类型
         *
         **/

         //1)、定义单向链表
        ACTION void setListNode(){
            //第一个结点
            listNode_a.self().name = "a";
            //第二个结点
            listNode_b.self().name = "b";
            listNode_b.self().nextPointer = &listNode_a.self();
            //第三个结点
            listNode_c.self().name = "c";
            listNode_c.self().nextPointer = &listNode_b.self();
        }


};

PLATON_DISPATCH(linkedListContractTest, (init)(setListNode))
