#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author qudong
 * 测试合约类模板std :: function
 * 类模板std :: function:是一个通用的多态函数包装器；
 *  std :: function的实例可以存储，复制和调用任何可调用的目标。
 *
 * */
class Plus{
    public:
        static uint8_t plus(uint8_t a, uint8_t b){
            return a + b;
        }
};
class PlusAdd{
    public:
        uint8_t add(uint8_t a, uint8_t b){
            return a + b;
        }
};
CONTRACT FunctionTemplateContract : public platon::Contract{
    private:
      platon::StorageType<"vector1"_n, std::vector<std::string>> storage_vector_string;
    public:
       ACTION void init(){
       }

       //1、调用lambda表达式
       CONST uint8_t get_lambda_function(){
             auto add = [](uint8_t a,uint8_t b) -> uint8_t {
                          return a + b;
                      };
             std::function<uint8_t(uint8_t,uint8_t)>function = add;
             return function(1,5);
       }

       //2、调用普通函数
      static uint8_t addFunc(uint8_t a, uint8_t b){
           return a + b;
        }
       CONST uint8_t get_normal_function(){
              //定义函数
              std::function<uint8_t(uint8_t,uint8_t)>function = addFunc;
              return function(1,5);
       }

       //3、调用类静态成员函数
      CONST uint8_t get_class_static_function(){
             //定义函数
             std::function<uint8_t(uint8_t,uint8_t)>function = &Plus::plus;
             return function(1,5);
      }




};

PLATON_DISPATCH(FunctionTemplateContract,(init)
(get_lambda_function)(get_normal_function)
(get_class_static_function)
)
