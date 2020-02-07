#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
 * @author qudong
 * 合约引用类型数组(array类型)：数组可以容纳相同类型的数据表
 * 测试验证功能点：
 * 1、定义array类型
 *
 * 2、定长/变长array
 *
 * */
extern char const array_uint8[] = "array_uint8_t";
extern char const array_string[] = "array_string";
extern char const array_bool[] = "array_bool";

CONTRACT arrayContractTest : public platon::Contract{

    private:
       platon::StorageType<array_string,std::array<std::string, 10>> array_string;
       platon::StorageType<array_uint8,std::array<uint8_t,10>> array_uint8;
       platon::StorageType<array_bool,std::array<bool,5>> array_bool;
    public:
        ACTION void init(){}

         /**
         * 1、定义array类型
         *    赋值/取值
         **/

         //1)、验证数组赋值
        ACTION void setArray(){
           array_string.self()[0] = "a";
           array_string.self()[1] = "b";
           array_string.self()[2] = "c";
        }

         ACTION void setInitArray(){
           array_uint8.self() = {1,2,3,4,5};
        }

        //2)、验证数组取值
         CONST uint8_t getArrayIndex(){
            return array_uint8.self()[0];
         }

        //3)、验证数组大小
         CONST uint8_t getArraySize(){
             return array_string.self().size();
         }

        /**
         *2、定长/变长array
         *
         **/
         //1)、验证定长数组赋值超出存储空间、赋值错误类型值
         ACTION void setArrayOver(){
               array_bool.self() = {true,false};//正常
               //array_bool.self() = {1,2,3};//异常，赋值错误类型值
               //array_bool.self() = {true,false,true,false,true,false,true,false};//异常，赋值超过存储空间
         }
         //2)、变长数组---？？如何编写语法




};

PLATON_DISPATCH(arrayContractTest, (init)(setArray)(setInitArray)(getArrayIndex)(getArraySize))
