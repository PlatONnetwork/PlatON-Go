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
       public:
         //原生写法
         //enum color {red, orange, yellow, green};
         //改成这样吗？
         enum color { red, green, blue } colors;
         clothes(){}
         PLATON_SERIALIZE(clothes, (colors))
};*/
extern char const storage_byte[] = "byte_storage";
extern char const storage_bool[] = "bool_storage";
extern char const storage_string[] = "string_storage";
extern char const storage_float[] = "float_storage";
extern char const storage_double[] = "double_storage";
extern char const storage_enum[] = "enum_storage";
extern char const storage_address[] = "address_storage";

CONTRACT BasicDataTypeContract : public platon::Contract{

    private:
       platon:: StorageType<storage_byte,byte> storage_byte;
       platon:: StorageType<storage_bool,bool> a_storage_bool;
       platon:: StorageType<storage_bool,bool> b_storage_bool;
       platon:: StorageType<storage_string,std::string> a_storage_string;
       platon:: StorageType<storage_string,std::string> b_storage_string;
       platon:: StorageType<storage_string,std::string> c_storage_string;
       platon:: StorageType<storage_address,Address> storage_address;
       platon:: StorageType<storage_float,float> storage_float;
       platon:: StorageType<storage_double,double> storage_double;
       //platon:: StorageType<storage_enum,clothes> storage_enum_clothes;

    public:
       ACTION void init(){
       }
      /**
       * 1、布尔型(bool)
       *   取值常量true、false
       *
       **/
      ACTION void setBool(){
          a_storage_bool.self() = true;
          b_storage_bool.self() = false;
      }

       CONST bool getBool(){
           return a_storage_bool.self();
       }

      /**
       * 2、字节类型（byte）
       *   byte相当于uint8_t
       **/
      ACTION void setByte(){
          storage_byte.self() = 100;//正常
      }

      /**
       * 3、字符串(string)
       *    字符串赋值、拼接、字符串.size()
       **/
       ACTION void setString(std::string &str){
           a_storage_string.self() = str;
       }

       CONST std::string getString(){
           return a_storage_string.self();
       }

       CONST uint8_t getStringLength(){
           return a_storage_string.self().size();
       }
     /**
      * 4、浮点类型(float、double)
      *  浮点型暂时不支持，后续测试
      **/
     /* ACTION void setFloat(){
          storage_float.self() = 1.0;
          storage_double.self() = 2.56;
       }*/

     /* CONST float getFloat(){
          return storage_float.self();
      }*/

      /**
       * 5、地址类型(Address)
       * 
       **/
      ACTION void setContractCallAddress(){
           storage_address.self() = platon_caller();//获取交易发起者地址
      }

      CONST std::string getContractCallAddress(){
          return storage_address.self().toString();
      }

      /**
       * 6、枚举(enum)
       **/
      /* ACTION void setEnum(){
            storage_enum_clothes.self().colors = yellow;
        }
        ACTION colors getEnum(){
            return  storage_enum_clothes.self().colors;
        }*/
};

PLATON_DISPATCH(BasicDataTypeContract,(init)(setBool)(getBool)(setByte)(setString)(getString)(getStringLength)
               (setContractCallAddress)(getContractCallAddress))
