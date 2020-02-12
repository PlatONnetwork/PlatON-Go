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



extern char const storage_int8[] = "int8storage";
extern char const storage_int16[] = "int16storage";
extern char const storage_uint8[] = "uint8storage";
extern char const storage_u160[] = "u160storage";
extern char const storage_bigInt[] = "bigintstorage";

CONTRACT BasicDataIntegerTypeContract : public platon::Contract{

    private:
       platon:: StorageType<storage_int8,int8_t> a_storage_int8;
       platon:: StorageType<storage_int16,int16_t> b_storage_int16;
       platon:: StorageType<storage_uint8,uint8_t> c_storage_uint8;
       platon:: StorageType<storage_uint8,uint8_t> d_storage_uint8;
       platon:: StorageType<storage_uint8,uint8_t> e_storage_uint8;
       platon:: StorageType<storage_u160,u160> g_storage_u160;
       platon:: StorageType<storage_bigInt,bigint> h_storage_bigint;
       platon:: StorageType<storage_int8,int8_t> f_storage_int8;
       platon:: StorageType<storage_uint8,uint8_t> i_storage_uint8;
       platon:: StorageType<storage_uint8,uint8_t> l_storage_uint8;
       platon:: StorageType<storage_uint8,uint8_t> j_storage_uint8;

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
       a_storage_int8.self() = 1;//正常
       b_storage_int16.self() = -10;//正常
       c_storage_uint8.self() = 3;//正常
       // d_storage_uint8 = -4;//异常，值范围：0~255，估无符号编译异常
    }

    CONST uint8_t getStorageUint8(){
        return c_storage_uint8.self();
    }

    CONST int8_t getStorageInt8(){
         return a_storage_int8.self();
     }

    CONST int16_t getStorageInt16(){
         return b_storage_int16.self();
    }

   //2)、验证无符号整数位数
   ACTION void setUint(){
        //8位无符号整数取值范围0~255
        i_storage_uint8.self() = 1;//正常
        l_storage_uint8.self() = 255;//正常
        //j_storage_uint8.self() = 256;//异常，8位数无符号整数溢出，编译报错
   }

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

/*PLATON_DISPATCH(BasicDataIntegerTypeContract,(init)(setSignedInt)(getStorageUint8)(setUint)(setInt)
                                          (setBigInt))*/

PLATON_DISPATCH(BasicDataIntegerTypeContract,(init)(setSignedInt)(getStorageUint8)(getStorageInt8)(getStorageInt16)
                (setUint))

/*
PLATON_DISPATCH(BasicDataIntegerTypeContract,(init)(setSignedInt))
*/
