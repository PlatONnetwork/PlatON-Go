#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author qudong
 * 测试合约bind函数
 * bind函数可以看做是一个通用的函数适配器，它可以接受一个可调用的对象，并生成一个新的可调用对象。
 * bind函数的形式：auto newCallable = bind(callable,arg_list);
 * newCallable本身是一个可调用对象，arg_list是一个逗号分隔的参数列表，对应给定的callable的参数
 * 即，当我们调用newCallable时，newCallable会调用callable,并传给它arg_list中的参数。
 * */

class Plus{
   public:
        uint32_t plus(uint32_t a,uint32_t b){
            return a + b;
        }
};

class PlusTwo{
	public:
		static int plus(int a,int b)
		{
		    return a + b;
		}
};

CONTRACT BindFunctionContract : public platon::Contract{
    private:
      platon::StorageType<"vector1"_n, std::vector<std::string>> storage_vector_string;
    public:
       ACTION void init(){
       }

       static uint8_t plus(uint8_t a,uint8_t b){
          return a + b;
       }

       //1、bind绑定普通函数
       CONST uint8_t get_bind_function(){
             //1)、通过bind函数，绑定plus()参数由funcPlus1来指定
             std::function<uint8_t(uint8_t,uint8_t)> funcPlus1 = std::bind(plus,std::placeholders::_1,std::placeholders::_2);
             //2)、通过bind函数，绑定plus()参数进行赋值
             auto funcPlus2 = std::bind(plus,1, 2);
             return funcPlus1(1,5);
       }

       //2、 bind绑定类的成员函数(指针形式调用成员函数)
      CONST uint32_t get_bind_class_function(){
            Plus p;
            //1)、
            std::function<uint32_t(uint32_t,uint32_t)> funcPlus1 = std::bind(&Plus::plus,&p,std::placeholders::_1,std::placeholders::_2);
            return funcPlus1(1,5);
      }

       //3、 bind绑定类的成员函数(对象形式调用成员函数)
     CONST uint32_t get_bind_class_function_one(){
          Plus p;
          std::function<uint32_t(uint32_t,uint32_t)> funcPlus = std::bind(&Plus::plus,p,std::placeholders::_1,std::placeholders::_2);
          return funcPlus(1,5);
     }

    //4、 bind绑定类静态成员函数
     CONST uint32_t get_bind_static_function(){
          Plus p;
          std::function<uint32_t(uint32_t,uint32_t)> funcPlus = std::bind(&PlusTwo::plus,std::placeholders::_1,std::placeholders::_2);
          return funcPlus(1,5);
     }



};

PLATON_DISPATCH(BindFunctionContract,(init)
(get_bind_function)
(get_bind_class_function)
(get_bind_class_function_one)
(get_bind_static_function)

)
