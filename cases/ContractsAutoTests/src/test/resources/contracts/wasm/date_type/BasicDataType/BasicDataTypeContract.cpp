#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author qudong
 * 测试合约基本数据类型
 * 1、整型
 * 2、布尔型
 * 3、字节类型（byte）
 * 4、字符串类型
 * 5、浮点类型(float、double)
 * 6、地址类型
 * */


CONTRACT BasicDataTypeContract : public platon::Contract{

    private:
       platon::StorageType<"bytekey"_n,byte> byte_v;
       platon::StorageType<"boolkey"_n,bool> bool_v;
       platon::StorageType<"strkey"_n,std::string> string_v;
       platon::StorageType<"addrkey"_n,Address> address_v;
//       platon::StorageType<"floatkey"_n,float> float_v;
//      platon::StorageType<"doublekey"_n,double> double_v;
       platon::StorageType<"long"_n,long> long_v;
       platon::StorageType<"long2"_n,long long> long_long_v;

    public:
       ACTION void init(){
       }
      /**
       * 1、布尔型(bool)
       *   取值常量true、false
       **/
      ACTION void set_bool(const bool &value){
          bool_v.self() = value;
      }
      CONST bool get_bool(){
          return bool_v.self();
      }

      /**
       * 2、字节类型（byte）
       *   byte相当于uint8_t
       **/
     ACTION void set_byte(const byte &value){
          byte_v.self() = value;
      }
      CONST byte get_byte(){
          return byte_v.self();
      }

      /**
       * 3、字符串(string)
       *    字符串赋值、字符串.size()
       **/
       ACTION void set_string(const std::string &value){
           string_v.self() = value;
       }
       CONST std::string get_string(){
           return string_v.self();
       }
       CONST uint8_t get_string_length(){
           return string_v.self().size();
       }

     /**
      * 4、浮点类型(float、double)
      * 验证结果：浮点型编译不通过
      **/
      //1)、float入参带&引用值
   /*  ACTION void set_float(const float &value){
          float_v.self() = value;
     }
     ACTION void set_float_one(const float value){
          float_v.self() = value;
     }
    CONST float get_float(){
          return float_v.self();
     }*/
     //2)、double入参带&引用值
   /*ACTION void set_double(const double &value){
          double_v.self() = value;
      }
     ACTION void set_double_one(const double value){
            double_v.self() = value;
     }
     CONST double get_double(){
         return double_v.self();
     }*/
     //3)、定义局部变量浮点型
     //验证结果：可以编译通过，脚本调用正常未见gas不足异常
    ACTION void set_float_type_local(){
           float a = 1.3f;
           double b = 2.56;
           double c = a + 3;
    }

      /**
       * 5、地址类型(Address)
       * 
       **/
     ACTION void set_address(){
           address_v.self() = platon_caller();//获取交易发起者地址
      }
      CONST std::string get_address(){
         return address_v.self().toString();
     }

  /**
   * 6、long类型
   */
   ACTION void set_long(const long &value){
      long_v.self()=value;
    }
    CONST long get_long(){
       return long_v.self();
    }

    ACTION void set_long_long(const long long &value){
       long_long_v.self()=value;
    }
    CONST long long get_long_long(){
       return long_long_v.self();
    }

    //7、枚举:由枚举常量构成，
    /**
     * 1)、枚举定义：只能以标识符形式表示，而不能以整型、字符型等文字常量。
     * 2)、枚举赋值
     * 3)、枚举类型限定作用域及非限定作用域
     *1、限定作用域的枚举类型：枚举成员的名字遵循常规的作用域准则，在枚举类型的作用域外是不可以访问的
     * 2、不限定作用域枚举类型：枚举成员的作用域和枚举类型本身作用域相同
     */
      //1)、枚举定义
      ACTION void set_enum_validity(){
         enum Weekday{SUN,MON,TUE,WED,THU,FRI,SAT};//合法编译正常
         //enum book{'a','b','c','d'};//非法编译异常
         //enum year{1998,1999,2010,2012};//非法编译异常
      }
      //2)、枚举赋值
      CONST uint8_t set_enum_assignment(){
          //枚举类型在声明之后具有默认值0->4值递增
          enum Animal{dog,cat,pig,chicken,duck};
          //自定义赋值,未赋值在此基础递增
          enum AnimalEnum{dogEnum = 7,catEnum = 2,pigEnum,chickenEnum,duckEnum};
          //对枚举元素不能对它赋予常量值
          //dogEnum = 1;//赋值非法编译
          return uint8_t(pigEnum);//返回值为3
      }
      //3)、枚举类型限定作用域及非限定作用域
      //A、限定作用域的枚举类型：枚举成员的名字遵循常规的作用域准则，在枚举类型的作用域外是不可以访问的
      //B、不限定作用域枚举类型：枚举成员的作用域和枚举类型本身作用域相同
      ACTION void set_enum_scope(){
         //定义非限定作用域的枚举类型
         enum color {red,yellow,green};
        // enum stoplight {red,yellow,green};//编译异常，作用范围内重复定义枚举成员
         //定义限定作用域的枚举类型
         enum class open_modes{input,output,append};
         enum class open_modes_enum{input,output,append};//编译正常，限定作用范围
      }
      //4)、枚举赋值
     CONST uint8_t set_enum_class_assignment(){
            enum size{x,xx,xxx,xxxl,xxxxl};
            enum class sizeEnum{xEnum,xxEnum,xxxEnum,xxxlEnum,xxxxlEnum};
            //赋值
            size a = size::xxx;
            sizeEnum b = sizeEnum::xxxEnum;
            return uint8_t(b);//返回值为2
        }



};

PLATON_DISPATCH(BasicDataTypeContract,(init)
(set_bool)(get_bool)
(set_byte)(get_byte)
(set_string)(get_string)(get_string_length)
(set_address)(get_address)
//(set_float)
//(set_float_one)
//(get_float)
//(set_double)
//(set_double_one)
//(get_double)
(set_float_type_local)
(set_long)(get_long)(set_long_long)(get_long_long)
(set_enum_validity)(set_enum_assignment)(set_enum_scope)
(set_enum_class_assignment)
)
