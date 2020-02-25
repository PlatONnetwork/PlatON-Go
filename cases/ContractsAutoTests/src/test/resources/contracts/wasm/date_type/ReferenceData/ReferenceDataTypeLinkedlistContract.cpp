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

 struct Edge{
        public:
           std::string nodeName;//结点数据
           uint8_t nextNode;//下一个结点编号
           Edge(){}
           PLATON_SERIALIZE(Edge,(nodeName)(nextNode))
};


CONTRACT ReferenceDataTypeLinkedlistContract : public platon::Contract{

    private:

     /*  platon::StorageType<"node_head"_n,listNode> storage_node_head;
       platon::StorageType<"node_tmp"_n,listNode> storage_node_tmp;
       platon::StorageType<"node_vector"_n, std::vector<listNode>> node_vector;
       platon::StorageType<"storage_array_uint8"_n,std::array<uint8_t,10>> storage_array_uint8;*/
       platon::StorageType<"array1"_n,std::array<std::vector<Edge>,10>> array_vector;
      // platon::StorageType<"int8value"_n,int8_t> count;
    public:
        ACTION void init(){}

         /**
         * 1、定义链表类型
         **/
         //1)、定义链表
        //初始化链表
        ACTION void insertNodeElement() {
           //std::vector<Edge> edgeArray[10];
           for (uint8_t i = 0;i < 10;i ++) {
           	    //遍历所有结点
           	    array_vector.self()[i].clear(); //清空其单链表
            }
            Edge tmp; //准备一个Edge结构体
            tmp.nextNode = 1; //下一结点编号为3
            tmp.nodeName = "one"; //该边权值为38
            array_vector.self()[0].push_back(tmp); //将该边加入结点1的单链表中
        }

        CONST uint8_t getNodeElement(const uint8_t &arrayIndex,const uint8_t &vectorIndex){

             return array_vector.self()[arrayIndex][vectorIndex].nextNode;
        }






};

PLATON_DISPATCH(ReferenceDataTypeLinkedlistContract,(init)
(insertNodeElement)(getNodeElement))
