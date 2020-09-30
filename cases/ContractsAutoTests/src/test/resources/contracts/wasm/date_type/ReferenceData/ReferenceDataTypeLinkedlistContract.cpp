#define TESTNET
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
       //每一个结点都建立一个单链表来保存与其相邻的边权值和结点的信息
       platon::StorageType<"array1"_n,std::array<std::vector<Edge>,10>> array_vector;
       platon::StorageType<"intvalue"_n,int8_t> count_int;
       platon::StorageType<"vector1"_n, std::vector<std::string>> vector_string;
    public:
        ACTION void init(){}

        /**
         * 1、定义链表类型
         **/
         //1)、链表新增元素
        ACTION void insertNodeElement(const std::string &nodeData) {
            if(array_vector.self().size() == 0){
                count_int.self() = 0;
            }else{
                count_int.self() = count_int.self() + 1;
            }
            Edge tmp; //准备一个Edge结构体
            tmp.nextNode = count_int.self(); //下一结点编号
            tmp.nodeName = nodeData; //该结点数据
            array_vector.self()[0].push_back(tmp); //将改元素添加到结点的单链表中(Vector)
        }

        //查询指定元素
        CONST uint8_t getNodeElementIndex(const uint8_t &arrayIndex,const uint8_t &vectorIndex){
             return array_vector.self()[arrayIndex][vectorIndex].nextNode;
        }

        //清除数据
        ACTION void clearNodeElement(){
              for (uint8_t i = 0;i < array_vector.self().size();i ++) {
                   //遍历所有结点
                   array_vector.self()[i].clear(); //清空其单链表
              }
        }
        //查询某个结点的所有邻接信息
         CONST std::vector<std::string> findNodeElement(const uint8_t &index){
              for (int i = 0;i < array_vector.self()[index].size(); i ++) {
              	//对所有与结点index相邻的边进行遍历，
              	uint8_t nextNode = array_vector.self()[index][i].nextNode; //结点编号
              	std::string nodeName = array_vector.self()[index][i].nodeName; //结点数据
              	DEBUG("ReferenceDataTypeLinkedlistContract", "查询某个结点的所有邻接信息nextNode：", nextNode);
              	DEBUG("ReferenceDataTypeLinkedlistContract", "查询某个结点的所有邻接信息nodeName：", nodeName);
                vector_string.self().push_back(nodeName);
               }
               return vector_string.self();
         }

};

PLATON_DISPATCH(ReferenceDataTypeLinkedlistContract,(init)
(insertNodeElement)(getNodeElementIndex)(clearNodeElement)
(findNodeElement))
