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
       platon:: StorageType<"a_storage_int8"_n,int8_t> a_storage_int8;
       platon:: StorageType<"b_storage_int16"_n,int16_t> b_storage_int16;
       platon:: StorageType<"c_storage_uint8"_n,uint8_t> c_storage_uint8;
       platon:: StorageType<"d_storage_uint8"_n,uint8_t> d_storage_uint8;
       platon:: StorageType<"e_storage_uint8"_n,uint8_t> e_storage_uint8;
       platon:: StorageType<"g_storage_u160"_n,u160> g_storage_u160;
       platon:: StorageType<"h_storage_bigint"_n,bigint> h_storage_bigint;
       platon:: StorageType<"f_storage_int8"_n,int8_t> f_storage_int8;
       platon:: StorageType<"i_storage_uint8"_n,uint8_t> i_storage_uint8;
       platon:: StorageType<"l_storage_uint8"_n,uint8_t> l_storage_uint8;
       platon:: StorageType<"j_storage_uint8"_n,uint8_t> j_storage_uint8;
       platon:: StorageType<"k_storage_int32"_n,int32_t> k_storage_int32;

    public:
       ACTION void init(){
       }
       /**
         * 1、整型
         * 1)、有符号整型int
         * 2)、无符号整型uint
         **/

    //赋值int8_t、int16_t、uint8_t
    ACTION void setStorageInt(int8_t a,int16_t b,int32_t c,uint8_t d){
          a_storage_int8.self() = a;
          b_storage_int16.self() = b;
          k_storage_int32.self() = c;
          c_storage_uint8.self() = d;

     }
    CONST int8_t getStorageInt8(){
         return a_storage_int8.self();
     }
    CONST uint8_t getStorageUint8(){
        return c_storage_uint8.self();
    }
    CONST int16_t getStorageInt16(){
         return b_storage_int16.self();
    }
    CONST int32_t getStorageInt32(){
          return k_storage_int32.self();
    }


    //赋值 u160、bigint
   /* ACTION void setStorageBigInt(u160 a,bigint b){
              g_storage_u160.self() = a;
              h_storage_bigint.self() = b;
    }
    CONST u160 getStorageU160(){
        return g_storage_u160.self();
    }
    CONST bigint getStorageBigInt(){
         return h_storage_bigint.self();
    }*/

   //2)、验证无符号整数位数
 /*  ACTION void setUint(){
        //8位无符号整数取值范围0~255
        i_storage_uint8.self() = 1;//正常
        l_storage_uint8.self() = 255;//正常
        //j_storage_uint8.self() = 256;//异常，8位数无符号整数溢出，编译报错
   }*/

     /*ACTION void setSignedInt(){
          a_storage_int8.self() = 1;//正常
          b_storage_int16.self() = -10;//正常
          c_storage_uint8.self() = 3;//正常
          // d_storage_uint8 = -4;//异常，值范围：0~255，估无符号编译异常
       }*/


       //3)、验证有符号整数位数 (负数SDK暂时不支持，后面测试)
  /* ACTION void setInt(){
        //8位有符号整数取值范围-128~127
        a_storage_int8.self() = 1;//正常
        a_storage_int8.self() = 127;//正常
        //a_storage_int8.self() = 128;//异常，整数溢出，编译报错

        f_storage_int8.self() =  -1;//正常
        f_storage_int8.self() =  -128;//正常
        //f_storage_int8.self() =  -129;//异常，整数溢出，编译报错
   }*/

      //4)、验证大位数整型赋值（SDK暂时不支持：u160,u256）
   /* ACTION void setBigInt(){
        g_storage_u160.self() = 99999999999999;
        h_storage_bigint.self() = 99999999999999999;
    }*/

   /* CONST u160 getStorageU160(){
         return g_storage_u160.self();
    }*/
};

PLATON_DISPATCH(BasicDataIntegerTypeContract,(init)(setStorageInt)(getStorageInt8)
               (getStorageUint8)(getStorageInt16))


