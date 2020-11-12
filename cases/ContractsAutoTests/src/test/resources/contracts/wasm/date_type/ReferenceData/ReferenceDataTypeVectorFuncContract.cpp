#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 测试合约引用类型(verctor类型)属性/函数
 * 1)、增加函数
 * 2)、遍历函数
 * 3)、删除函数
 * 4)、判断函数
 * 5)、其他函数
 * vector是一个能够存放任意类型的动态数组，能够增加和压缩数据
 * 
 * */

CONTRACT ReferenceDataTypeVectorFuncContract : public platon::Contract{

    private:
      platon::StorageType<"vector1"_n, std::vector<std::string>> storage_vector_string;
    public:
    ACTION void init(){}


    //1)、增加函数：将值添加到向量begin()起始位置
    //A、指定位置插入数据，在a的第1个元素（从第0个算起）的位置插入数值5，如a为1,2,3,4，插入元素后为5,1,2,3,4
     ACTION void insertVectorValue(const std::string &my_value){
         storage_vector_string.self().insert(storage_vector_string.self().begin(),my_value);
     }
    //B、指定位置插入多个相同数据，在a的第1个元素（从第0个算起）的位置插入3个数，其值都为5
     ACTION void insertVectorMangValue(const uint64_t &num,const std::string &my_value){
          storage_vector_string.self().insert(storage_vector_string.self().begin(),num,my_value);
     }
     //取值
    /* CONST std::string getClothesColorIndex(){
                return vector_clothes.self()[0].color;
            }*/
     //vector大小
     CONST uint64_t getVectorLength(){
       return storage_vector_string.self().size();
     }


     //2)、遍历函数
     //A、at()函数，返回index位置元素的引用
    CONST std::string findVectorAt(const uint64_t &index){
           return storage_vector_string.self().at(index);
     }
    //B、front():返回首元素的引用
    CONST std::string findVectorFront(){
           return storage_vector_string.self().front();
    }
    //C、 back():返回尾元素的引用
    CONST std::string findVectorBack(){
           return storage_vector_string.self().back();
    }

     //3)、删除函数
     //A、pop_back()删除最后一个元素
     ACTION void deleteVectorPopBack(){
           storage_vector_string.self().pop_back();
     }
    //B、erase()删除指定元素,将起始位置的元素删除
     ACTION void deleteVectorErase(){
            storage_vector_string.self().erase(storage_vector_string.self().begin());
     }
    //C、clear()清空元素
    ACTION void deleteVectorClear(){
            storage_vector_string.self().clear();
     }

     //4)、 empty()判断函数
     CONST bool findVectorEmpty(){
           return storage_vector_string.self().empty();
     }

};

PLATON_DISPATCH(ReferenceDataTypeVectorFuncContract, (init)(insertVectorValue)(insertVectorMangValue)(getVectorLength)
               (findVectorAt)(findVectorFront)(findVectorBack)(deleteVectorPopBack)(deleteVectorErase)(deleteVectorClear)
               (findVectorEmpty))
