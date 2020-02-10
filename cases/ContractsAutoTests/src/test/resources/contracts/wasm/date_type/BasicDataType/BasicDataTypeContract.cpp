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
 * 7、枚举类型
 *
 * 
 * */

/*class clothes{
       enum color {red,orange, yellow, green};
       public:
         clothes(){}
         // clothes(const color &c):color1(c){}
         color colors;
         PLATON_SERIALIZE(clothes, (colors))
};*/

extern char const int8_a[] = "int8_t";
extern char const int16_b[] = "int16_t";
extern char const uint8_c[] = "uint8_t";
extern char const uint8_d[] = "uint8_t";
extern char const uint8_e[] = "uint8_t";
extern char const int8_f[] = "int8_t";
extern char const u160_g[] = "u160";
extern char const bigInt_h[] = "bigint";

extern char const byte_a[] = "byte";
extern char const bool_a[] = "bool";
extern char const string_a[] = "string";
extern char const string_b[] = "string";
extern char const string_c[] = "string";
extern char const enum_a[] = "enum";
extern char const address_a[] = "address";
extern char const float_a[] = "float";
extern char const double_a[] = "double";



CONTRACT basicDataTypeContract : public platon::Contract{    

    private:
       platon:: StorageType<int8_a,int8_t> aInt8;
       platon:: StorageType<int8_f,int8_t> fInt8;
       platon:: StorageType<int16_b,int16_t> bInt16;
       platon:: StorageType<uint8_c,uint8_t> cUint8;
       platon:: StorageType<uint8_d,uint8_t> dUint8;
       platon:: StorageType<uint8_e,uint8_t> eUint8;
      // platon:: StorageType<u160_g,u160> gU160;
      // platon:: StorageType<bigInt_h,bigint> hBigInt;

     /*  platon:: StorageType<byte_a,byte> aByte;
       platon:: StorageType<bool_a,bool> aBool;
       platon:: StorageType<bool_a,bool> bBool;

       platon:: StorageType<string_a,std::string> aString;
       platon:: StorageType<string_b,std::string> bString;
       platon:: StorageType<string_c,std::string> cString;
       // platon:: StorageType<enum_a,enum> aEnum;
       platon:: StorageType<address_a,Address> contractAddress;
       platon:: StorageType<float_a,float> aFloat;
       platon:: StorageType<double_a,double> aDouble;*/

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
           aInt8.self() = 1;//正常
           bInt16.self() = -10;//正常
           cUint8.self() = 3;//正常
           // dUint8 = -4;//异常，值范围：0~255，估无符号编译异常
        }

       //2)、验证无符号整数位数
       ACTION void setUint(){
            //8位无符号整数取值范围0~255
            cUint8.self() = 1;//正常
            dUint8.self() = 255;//正常
            //eUint8.self() = 256;//异常，8位数无符号整数溢出，编译报错
       }

       //3)、验证有符号整数位数
     /* ACTION void setInt(){
            //8位有符号整数取值范围-128~127
            aInt8.self() = 1;//正常
            aInt8.self() = 127;//正常
            //aInt8.self() = 128;//异常，整数溢出，编译报错

            fInt8.self() =  -1;//正常
            fInt8.self() =  -128;//正常
            //fInt8.self() =  -129;//异常，整数溢出，编译报错
      }*/

      //4)、大位数整型赋值
    /*  ACTION void setBigInt(){
           gU160.self() = 99999999999999;
           hBigInt.self() = 99999999999999999;
      }*/

      /**
       * 2、布尔型(bool)
       *   取值常量true、false
       *
       **/
     /* ACTION void setBool(){
          aBool.self() = true;
          bBool.self() = false;
      }

       CONST bool getBool(){
           return aBool.self();
       }*/

      /**
       * 3、字节类型（byte）
       *   byte相当于uint8_t
       **/
    /*  ACTION void setByte(){
          aByte.self() = 100;//正常
      }*/

      /**
       * 4、字符串(string)
       *    字符串赋值、拼接、字符串.size()
       **/
      /* ACTION void setString(){
           aString.self() = "A";
           bString.self() = "B";
           cString.self() = "C" +  bString.self();
       }

       CONST std::string getString(){
           return cString.self();
       }

       CONST void getStringLength(){
           cString.self().size();
       }*/

     /**
      * 5、浮点类型(float、double)
      *
      **/
     /* ACTION void setFloat(){
          aFloat.self() = 1.0;
          aDouble.self() = 2.56;
       }*/

      /**
       * 6、地址类型(Address)
       * 
       **/
     /* ACTION void setContractAddress(){
           contractAddress.self() = platon_caller();//获取交易发起者地址
      }

      CONST std::string getContractAddress(){
          return contractAddress.self().toString();
      }*/



        /**
              * 7、枚举(enum)
              **/
          /*  ACTION void setEnum(const clothes &c){
                aEnum.self() = c;
            }

            ACTION clothes getEnum(){
                return  aEnum.self();
            }*/

};

/*PLATON_DISPATCH(basicDataTypeContract,(init)(setUint)(setBool)(getBool)
               (setByte)(setString)(getString)(getStringLength)(setFloat))*/
PLATON_DISPATCH(basicDataTypeContract,(init)(setUint))