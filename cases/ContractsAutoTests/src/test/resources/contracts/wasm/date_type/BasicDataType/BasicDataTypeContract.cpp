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
 * */


CONTRACT BasicDataTypeContract : public platon::Contract{

    private:
       platon::StorageType<"bytekey"_n,byte> byte_v;
       platon::StorageType<"boolkey"_n,bool> bool_v;
       platon::StorageType<"strkey"_n,std::string> string_v;
       platon::StorageType<"addrkey"_n,Address> address_v;
      // platon::StorageType<"floatkey"_n,float> float_v;
      // platon::StorageType<"doublekey"_n,double> double_v;

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
      **/
      ACTION void set_float(const float &value){
          float_v.self() = value;
       }
      CONST float get_float(){
          return float_v.self();
      }
      ACTION void set_double(const double &value){
          double_v.self() = value;
      }
      CONST double get_double(){
           return double_v.self();
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


};

PLATON_DISPATCH(BasicDataTypeContract,(init)
(set_bool)(get_bool)(set_byte)(get_byte)
(set_string)(get_string)(get_string_length)
(set_address)(get_address)
(set_float)(get_float)(set_double)(get_double)
)
