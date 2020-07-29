#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型数组(array类型)：
 * 测试验证功能点：
 * 1、验证数组属性
 *    1)、size():返回数组的大小
 *    2)、empty():判断数组是否为空
 *    3)、at()和[]操作运算符实现一样的功能。
 *    4)、front():获取数组的第一个元素
 *    5)、fill():将数组每个元素都固定填一个值
 * */

CONTRACT ReferenceDataTypeArrayFuncContract:public platon::Contract{

    private:
       platon::StorageType<"arr"_n,std::array<std::string,10>> storage_array_string;

    public:
        ACTION void init(){}
         //1)、数据源初始化准备
         ACTION void setInitArrayDate(){
            storage_array_string.self()[0] = "a";
            storage_array_string.self()[1] = "b";
            storage_array_string.self()[2] = "c";
            storage_array_string.self()[3] = "d";
            storage_array_string.self()[4] = "e";
        }

        //2)、验证数组empty()
         CONST bool getArrayIsEmpty(){
            return storage_array_string.self().empty();
         }
       //3)、验证数组at()，获取数组值
         CONST std::string getArrayValueIndex(const uint32_t &index){
           return storage_array_string.self().at(index);
         }
        //4)、验证数组front()
         CONST std::string getArrayFirstValue(){
           return storage_array_string.self().front();
         }
        //5)、验证数组fill()
        ACTION void setArrayFill(std::string str){
            storage_array_string.self().fill(str);
        }
};

PLATON_DISPATCH(ReferenceDataTypeArrayFuncContract,(init)(setInitArrayDate)(getArrayIsEmpty)(getArrayValueIndex)
               (getArrayFirstValue)(setArrayFill))

