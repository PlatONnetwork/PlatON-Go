#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * @author qudong
 * 测试合约基本数据类型
 * 1、整型
 * 1)、验证有符号/无符号整型
 * 2)、验证无符号整数位数
 * 3)、验证有符号整数位数
 * 4)、验证大位数整型赋值
 * */

CONTRACT BasicDataIntegerTypeContract : public platon::Contract{


    private:
       platon::StorageType<"int1key"_n,int8_t> int8_v;
       platon::StorageType<"int2key"_n,int16_t> int16_v;
       platon::StorageType<"int3key"_n,int32_t> int32_v;
       platon::StorageType<"int4key"_n,int64_t> int64_v;
       platon::StorageType<"uint1key"_n,uint8_t> uint8_v;
       platon::StorageType<"uint2key"_n,uint16_t> uint16_v;
       platon::StorageType<"uint3key"_n,uint32_t> uint32_v;
       platon::StorageType<"uint4key"_n,uint64_t> uint64_v;

       platon::StorageType<"key1range"_n,int8_t> int8_range;
       platon::StorageType<"key2range"_n,int16_t> int16_range;
       platon::StorageType<"key3range"_n,int32_t> int32_range;
       platon::StorageType<"key4range"_n,int64_t> int64_range;
       platon::StorageType<"key5range"_n,uint8_t> uint8_range;
       platon::StorageType<"u1range"_n,uint16_t> uint16_range;
       platon::StorageType<"u2range"_n,uint32_t> uint32_range;
       platon::StorageType<"u3range"_n,uint64_t> uint64_range;
       platon::StorageType<"u1key"_n,u128> u128_v;
       //platon::StorageType<"u1key"_n,u160> u160_v;
       //platon::StorageType<"u2key"_n,u256> u256_v;
       //platon::StorageType<"bigintkey"_n,bigint> bigint_v;

    public:
       ACTION void init(){
       }
       /**
         * 1、整型
         * 1)、有符号整型int
         * 2)、无符号整型uint
         **/

    //赋值int8_t、int16_t、int32_t、uint8_t
    ACTION void set_uint8(const uint8_t &value){
   	     uint8_v.self()=value;
   	  }

   	  CONST uint8_t get_uint8(){
   	  	return uint8_v.self();
   	  }


   	  ACTION void set_uint16(const uint16_t &value){
   	     uint16_v.self()=value;
   	  }

   	  CONST uint16_t get_uint16(){
   	  	return uint16_v.self();
   	  }

   	  ACTION void set_uint32(const uint32_t &value){
   	     uint32_v.self()=value;
   	  }

   	  CONST uint32_t get_uint32(){
   	  	return uint32_v.self();
   	  }

   	  ACTION void set_uint64(const uint64_t &value){
   	     uint64_v.self()=value;
   	  }

   	  CONST uint64_t get_uint64(){
   	  	return uint64_v.self();
   	  }

   	  ACTION void set_int8(const int8_t &value){
   	     int8_v.self()=value;
   	  }

   	  CONST int8_t get_int8(){
   	  	return int8_v.self();
   	  }


   	  ACTION void set_int16(const int16_t &value){
   	     int16_v.self()=value;
   	  }

   	  CONST int16_t get_int16(){
   	  	return int16_v.self();
   	  }

   	  ACTION void set_int32(const int32_t &value){
   	     int32_v.self()=value;
   	  }

   	  CONST int32_t get_int32(){
   	  	return int32_v.self();
   	  }

   	  ACTION void set_int64(const int64_t &value){
   	     int64_v.self()=value;
   	  }

   	  CONST int64_t get_int64(){
   	  	return int64_v.self();
   	  }

     /**
      *  3)、验证大位数整型赋值
      */
      ACTION void set_u128(const uint64_t &value){
         u128_v.self()=value;
      }

      CONST std::string get_u128(){
         return std::to_string(u128_v.self());
      }


     /* CONST std::string get_u128(uint64_t input)
     {
        u128 u = u128(input);
        return std::to_string(u);
     }*/

     /* ACTION void set_u160(const uint64_t &value){
          u160_v.self() = u160(value);
      }
      CONST std::string get_u160(){
          return to_string(u160_v.self());
      }

      ACTION void set_u256(const uint64_t &value){
          u256_v.self() = u256(value);
      }
      CONST std::string get_u256(){
          return to_string(u256_v.self());
      }

      ACTION void set_bigint(const uint64_t &value){
          bigint_v.self() = bigint(value);
      }
      CONST std::string get_bigint(){
          return to_string(bigint_v.self());
      }*/

    /**
     * 4)、验证有符号整数位数
     */
     /* //int8()取值范围：-128~127
      ACTION void set_int8_range(){
         //int8_range.self() = -129;//异常，整数溢出，编译报错
         int8_range.self() = -128;  //正常
         int8_range.self() = -100;  //正常
         int8_range.self() = 1;     //正常
         int8_range.self() = 127;   //正常
         //int8_range.self() = 128; //异常，整数溢出，编译报错
      }

      //int16()取值范围：-32768~32767
      ACTION void set_int16_range(){
         //int16_range.self() = -32770;//异常，整数溢出，编译报错
         int16_range.self() = -32768;  //正常
         int16_range.self() = 1;       //正常
         int16_range.self() = 127;     //正常
         int16_range.self() = 32767;   //异常，整数溢出，编译报错
         //int8_range.self() = 32768;  //异常，整数溢出，编译报错
      }*/

      /**
       * 5)、验证无符号整数位数
       **/

      //uint8()取值范围：0~255
    /*  ACTION void set_uint8_range(){
          //uint8_range.self() = -10; //此处编译未报错？？待开发查看
          uint8_range.self() = 0;     //正常
          uint8_range.self() = 255;   //正常
         // uint8_range.self() = 256;  //异常警告，整数溢出
      }
      //uint16()取值范围：0~65535
     ACTION void set_uint16_range(){
         // uint16_range.self() = -1;  //此处编译未报错？？待开发查看
         uint16_range.self() = 0;      //正常
         uint16_range.self() = 65535;  //正常
         //uint16_range.self() = 65536;//异常，整数溢出，编译报错
     }
     //uint32()取值范围：0~4294967295
     ACTION void set_uint32_range(){
         //uint32_range.self() = -1;       //此处编译未报错？？待开发查看
         uint32_range.self() = 0;           //正常
         uint32_range.self() = 4294967295;  //正常
         //uint32_range.self() = 4294967296;//异常，整数溢出，编译报错
     }*/



};

PLATON_DISPATCH(BasicDataIntegerTypeContract,(init)(set_uint8)(get_uint8)(set_uint16)(get_uint16)(set_uint32)
(get_uint32)(set_uint64)(get_uint64)(set_int8)(get_int8)(set_int16)(get_int16)(set_int32)(get_int32)
(set_int64)(get_int64)(set_u128)(get_u128)
//(set_u256)(get_u256)(set_u160)(get_u160)(set_bigint)(get_bigint)
)
//(set_u160)(get_u160)(set_u256)(get_u256)(set_bigint)(get_bigint))
//(set_int8_range)(set_int16_range)(set_uint8_range)(set_uint16_range)(set_uint32_range)

