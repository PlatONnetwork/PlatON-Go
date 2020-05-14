#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author qudong
 * 测试合约数据类型转换
 * 主要分为两种类型：隐式类型转换、显示类型转换
 * 1、隐式类型转换：编译器默认进行的类型转换行为
 * 2、显示类型转换(强制类型转换)
 *    强制类型转换操作符：static_cast、const_cast、dynamic_cast
 * */

class People {
  public:
    virtual void setAge(uint8_t age){
        this->age = age;
    }
    uint8_t getAge(){
        return this->age;
    }
  private:
    uint8_t age;
    PLATON_SERIALIZE(People,(age))
};
class Student : public People {
     private:
        uint64_t sAge;//学生年龄
      public:
         void setSAge(uint64_t sAge) {
             this->sAge = sAge;
         };
         uint64_t getSAge() {
             return this->sAge;
         }
         PLATON_SERIALIZE_DERIVED(Student,People,(sAge))
};

CONTRACT TypeConversionContract : public platon::Contract{

    private:
      platon::StorageType<"vector1"_n, std::vector<std::string>> storage_vector_string;

    public:
       ACTION void init(){
       }

    //1、隐式类型转换
    //混合类型的算术运算表达式中(由低精度向高精度的转换)
    CONST auto get_add(uint8_t a,uint64_t b){
       //a会被自动转换成double
        auto c = a + b;
        return c;
    }

    //不同类型的赋值操作时(由低精度向高精度的转换)
    CONST uint64_t get_different_type_(bool a){
       //bol类型被转换成uint64
        uint64_t b = a;
        return b;
    }

    //函数参数传值时类型转换(由低精度向高精度的转换)
    CONST auto get_pram_type(){
        auto test = [](uint32_t a){
             return a;
        };
        //调用函数整数值被转换成double
        return test(1);
    }

    //函数返回值时(由低精度向高精度的转换)
    CONST uint64_t get_pram_return(const uint8_t a,uint8_t b){
       //运算结果会被隐式转换为double类型返回
        return (a + b);
    }

    //2、显示类型转换(强制类型转换)
    //显式类型转换(由高精度向低精度的转换)
    CONST uint8_t get_convert(){
        uint32_t a = 1002;
        uint8_t b = (uint8_t)a;
        return b;
    }

    // static_cast:非多态类型转换(静态转换)主要用于内置数据类型之间的相互转换
    //不能在没有派生关系的两个类类型之间转换
   CONST uint8_t get_convert_static_cast(){
          uint32_t aValue = 500;
          uint8_t bValue = static_cast<uint8_t>(aValue);
          return bValue;
    }

    //const_cast:删除变量的const属性，方便再次赋值,只能转换指针或者引用
    CONST uint8_t get_convert_const_cast(){
         const uint8_t i = 10;
         uint8_t *p = const_cast<uint8_t*>(&i);
         return i;
    }

    //dynamic_cast:执行派生类指针或引用与基类指针或引用之间的转换
    //其转换是运行时处理的，不能用于内置基本类型的强制转换
    //编译异常，编译器不支持
    /*CONST bool get_convert_dynamic_cast(){
         People* p = new People();
         Student* s = nullptr;
         s = dynamic_cast<Student*>(p);//向下转型
         if(nullptr == s){
            return false;
         }else{
            return true;
         }
         return false;
    }*/



};

PLATON_DISPATCH(TypeConversionContract,(init)
(get_add)(get_different_type_)(get_pram_type)
(get_pram_return)(get_convert)
(get_convert_static_cast)(get_convert_const_cast)
//(get_convert_dynamic_cast)
)
