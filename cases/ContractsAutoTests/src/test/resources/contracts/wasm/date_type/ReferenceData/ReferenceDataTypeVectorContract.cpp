#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 测试合约引用类型(verctor类型)
 * vector是一个能够存放任意类型的动态数组，能够增加和压缩数据
 * 
 * */

class clothes {
    public:
        std::string color;
        clothes(){}
        clothes(const std::string &my_color):color(my_color){}
        PLATON_SERIALIZE(clothes, (color))
};

CONTRACT ReferenceDataTypeVectorContract : public platon::Contract{
    public:
    ACTION void init(){}

     //新增方式一
    ACTION void setClothesColorOne(const clothes &myClothes){
        vector_clothes.self().push_back(myClothes);
    }
    //新增方式二
      ACTION void setClothesColorTwo(const std::string &my_color){
         vector_clothes.self().push_back(clothes(my_color));
      }
     //取值
     CONST std::string getClothesColorIndex(){
         return vector_clothes.self()[0].color;
     }
    //vector大小
    CONST uint64_t getClothesColorLength(){
        return vector_clothes.self().size();
    }

    private:
    platon::StorageType<"vector1"_n, std::vector<clothes>> vector_clothes;
};

PLATON_DISPATCH(ReferenceDataTypeVectorContract, (init)(setClothesColorOne)(setClothesColorTwo)(getClothesColorIndex)
               (getClothesColorLength))
