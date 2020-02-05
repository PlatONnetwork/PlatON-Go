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
        clothes(){}
        clothes(const std::string &my_color):color(my_color){}
    private:
        std::string color;
        PLATON_SERIALIZE(clothes, (color))
};

extern char const vector_clothes[] = "vector_clothes";

CONTRACT clothesVector : public platon::Contract{
    public:
    ACTION void init(){}

     //新增
    ACTION void setClothesColor(const clothes &myClothes){
        vector_clothes.self().push_back(myClothes);
    }
    //vector大小
    CONST uint64_t getClothesColor(){
        return vector_clothes.self().size();
    }

    
    private:
    platon::StorageType<vector_clothes, std::vector<clothes>> vector_clothes;
};

PLATON_DISPATCH(clothesVector, (init)(setClothesColor)(getClothesColor))
