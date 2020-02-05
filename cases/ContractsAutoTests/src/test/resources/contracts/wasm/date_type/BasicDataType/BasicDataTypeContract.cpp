#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author qudong
 * 测试合约基本数据类型
 * 1、整型
 * 2、布尔型
 * 3、地址类型
 * 4、枚举类型
 * 5、字符串类型
 * 
 * */

class clothes{
    public:
      enum color {red,orange, yellow, green};
      clothes(){}
      clothes(const color &c):color1(c){}
       
    private:
       color color1;
       PLATON_SERIALIZE(clothes, (color1)) 
};


extern char const int8_a[] = "int8_t";
extern char const int32_b[] = "int32_t";
extern char const uint8_c[] = "uint8_t";
extern char const uint8_d[] = "uint8_t";
extern char const uint8_e[] = "uint8_t";

extern char const bool_a[] = "bool";

extern char const string_a[] = "string";
extern char const string_b[] = "string";
extern char const string_c[] = "string";

extern char const enum_a[] = "enum";

//extern char const int128_var[] = "int128_t";
//extern char const int256_var[] = "int256_t";


CONTRACT basicDataTypeContract : public platon::Contract{    

    private:
      //platon:: StorageType<int8_a,int8_t> aInt8;
      //platon:: StorageType<int32_b,int> bInt32;
       platon:: StorageType<uint8_c,uint8_t> cUint8;
       platon:: StorageType<uint8_d,uint8_t> dUint8;
       platon:: StorageType<uint8_e,uint8_t> eUint8;
       
       platon:: StorageType<bool_a,bool> aBool;
       platon:: StorageType<bool_a,bool> bBool;

       platon:: StorageType<string_a,std::string> aString;
       platon:: StorageType<string_b,std::string> bString;
       platon:: StorageType<string_c,std::string> cString;

       platon:: StorageType<enum_a,clothes> aEnum;
    public:

       ACTION void init(){
       }

       /**
         * 1、整型
         * 1)、有符号整型int
         * 2)、无符号整型uint
         **/

       //1)、验证有符号/无符号整型
       ACTION void setSignedInt(){
           // aInt8.self() = 1;//异常，不支持int --------??待开发协助
          // bInt32.self() = -10;
           cUint8.self() = 3;//正常
          // dUint8 = -4;//异常，值范围：0~255，估无符号编译异常
       } 

       //2)、验证无符号整数位数
       ACTION void setUint(){
            cUint8.self() = 1;//正常
            dUint8.self() = 255;//正常，8位无符号整数取值范围0~255
            //eUint8.self() = 256;//异常，8位数无符号整数溢出，编译报错
       }

       //3)、验证有符号整数位数
      ACTION void setInt(){
          //整型类型暂时不支持，待支持后编写
          // todo.....
      }

      /**
       * 2、布尔型(bool)
       *   取值常量true、false
       *
       **/
      ACTION void setBool(){
          aBool.self() = true;
          bBool.self() = false;
      }

       CONST bool getBool(){
           return aBool.self();
       } 
      
      /**
       * 3、字符串(string)
       *    字符串赋值、拼接、字符串.size()
       **/
       ACTION void setString(){
           aString.self() = "A";
           bString.self() = "B";
           cString.self() = "C" +  bString.self();
       }

       CONST std::string getString(){
           return cString.self();
       }

       CONST void getStringLength(){
           cString.self().size();
       }
       
       /**
        * 4、枚举(enum)
        **/
      ACTION void setEnum(const clothes &c){
          aEnum.self() = c; 
      }

      ACTION clothes getEnum(){
          return  aEnum.self();
      }

      /**
       * 5、地址类型(Address)
       * 
       **/







      
};

PLATON_DISPATCH(basicDataTypeContract,(init)(setSignedInt)(setUint)(setBool)(setString))
