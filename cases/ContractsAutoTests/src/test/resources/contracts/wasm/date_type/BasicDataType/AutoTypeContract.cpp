#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author qudong
 * 测试合约anto关键字
 *  作用：auto可以在声明变量的时候根据初始值的类型自动为此类型变量匹配类型
 * */

CONTRACT AntoTypeContract : public platon::Contract{

    private:
      platon::StorageType<"vector1"_n, std::vector<std::string>> storage_vector_string;

    public:
       ACTION void init(){
       }

    //自动匹配int类型
    CONST auto get_anto_int(){
        auto a = 5;//值为整数
        return a;
    }
    //自动匹配
    CONST auto get_anto_int32(){
        auto a = -10;//值为负数
        return a;
    }
    //自动匹配int*类型--指针类型编译器不支持
  /*  CONST auto get_anto_int2(){
        auto b = new auto(1);//y是int*
        return b;
    }*/
    //自动匹配int*类型--指针类型编译器不支持
   /* CONST auto get_anto_int3(){
        auto a = 5;//x是int
        auto c = &a;//c是int*类型
        return c;
    }*/
    //自动匹配double类型(暂不支持浮点)
 /*   CONST auto get_anto_double(){
        auto d = 1.0;//d是double类型
        return d;
    }*/
    //自动匹配多个值类型
    CONST auto get_anto_multiple(){
        auto e = 10,f = 20,g = 30;//auto定义多个值时必须类型一致
        return g;
    }
     //自动匹配uint8类型
    CONST auto get_anto_uint8_t(){
        uint8_t a1 = 10;
        auto v1 = a1;//v1是uint8_t类型
        return v1;
    }
     //自动匹配数组类型--指针类型不支持
 /*   CONST auto get_anto_array(){
        int64_t arr[3] = {1,2,3};
        auto v2 = arr;//v2是int64_t
        return v2;
    }*/

    //自动匹配表达式(暂不支持浮点)
   /*CONST auto get_anto_express(){

     auto a = 5;
     auto b = 10.32;
     auto c = a + b;//c是double
     return c;
   }*/

      ACTION void set_anto_care_one(){
         //1、如果初始化表达式是引用，则去除引用语义
         int a = 10;
         int &b = a;
         auto c = b;//c的类型为int而非int&（去除引用）

        //2、如果初始化表达式为const或volatile（或者两者兼有），则除去const/volatile语义
        const int a1 = 10;
        auto  b1= a1; //b1的类型为int而非const int（去除const）
        const auto c1 = a1;//此时c1的类型为const int
      }


      /**
       * 函数参数不能被声明为auto
       * auto仅仅是一个占位符，它并不是一个真正的类型，
       * 不能使用一些以类型为操作数的操作符
       */
       /* ACTION void set_anto_func(auto a,auto b)
            auto c = a * b;
        }*/
       //auto应用迭代器中
      CONST uint8_t get_anto_iterator(){
          uint8_t count = 0;
         storage_vector_string.self().push_back("v1");
         storage_vector_string.self().push_back("v2");
         storage_vector_string.self().push_back("v3");
         for(auto it = storage_vector_string.self().begin(); it != storage_vector_string.self().end(); ++it)
         {
             count += 1;
         }
          return count;
      }

};

PLATON_DISPATCH(AntoTypeContract,(init)
(get_anto_int)(get_anto_int32)
//(get_anto_int2)(get_anto_int3)
//(get_anto_double)
(get_anto_multiple)
(get_anto_uint8_t)//(get_anto_array)
//(get_anto_express)
(set_anto_care_one)(get_anto_iterator)
)
