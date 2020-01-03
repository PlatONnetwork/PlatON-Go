pragma solidity 0.5.13;


/**
 *
 * 
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试类型转换
 *
 *类型简述：
 *1、基本类型之间的转换
 *
 *1)、隐式转换：
 *    A、在进行运算符运算时，两个不同类型的变量之间，那么编译器将隐式地将其中一个类型转换为另一个类型。
 *    B、不同类型之间的赋值
 *    只要值类型之间的转换在语义上行得通，而且转换的过程中没有信息丢失，那么隐式转换基本都是可以实现的，
 *    隐式转换可以将一个类型转换成与它相当大小，或者更大的类型，反正不行。
 *
 *2)、显式转换：
 *    编译器不会将语法上不可转换的类型进行隐式转换，此时我们要通过显式转换的方式。
 *    A、一个类型显式转换成更小的类型，相应的高位将被舍弃 
 *    B、将一个类型显式转换为更大的类型，则将填充左侧（即在更高阶的位置）
 *

 *-----------------  测试点   ------------------------------
 *1、基本类型隐式转换
 *   1)、进行运算符操作
 *   2)、进行赋值操作
 *   3)、不同类型转换 
 *2、基本类型显式转换
 *   1)、整数操作：大类型<-->小类型
 *   2)、字节操作：大类型<-->小类型
 *
 */

contract TypeConversionContractTest {


   /**
    *1、基本类型隐式转换
    *1)、进行运算符操作
    *2)、进行赋值操作
    */

    int8 a = 2;
    int16 b = 100;

    //运算符操作隐式转换
    function sum() public view returns(int16) {

        return  a + b;
    }

    //赋值操作隐式转换
    function conversion() public view returns(uint16) {

       uint8 a = 10;
       uint16 b = a;
       return b;
    }



   /**
    *2、基本类型显示转换
    *
    *
    */

    //无符合与有符号转换
    function displayConversion() public view returns (int8) {

        uint8 a = 1;
        int8 b = int8(a);
        return b;
    }
    
    //转换成更小的类型，会丢失高位
    function displayConversion1() public view returns (uint16,bytes2) {

        uint32 a = 0x12345678; //二进制：//‭0001 0010 0011 0100   0101 0110 0111 1000，十进制：305,419,896
        uint16 b = uint16(a);  //转换高位截取丢失，即 0101 0110 0111 1000，十进制：22136
        return (b,bytes2(b));
    }

    //转换成更大的类型，将向左侧添加填充位
    function displayConversion2() public view returns (uint32,bytes4) {

        uint16 a = 0x1234; //二进制：0001 0010 0011 0100‬，十进制：4660
        uint32 b = uint32(a);//转换
        return (b,bytes4(b));
    }

    //转换到更小的字节类型，会丢失后面数据
    function displayConversion3() public view returns (bytes1) {

        bytes2 a = 0x1234;
        bytes1 b = bytes1(a);
        return b;
    }

    //转换为更大的字节类型时，向右添加填充位
    function displayConversion4() public view returns (bytes4) {

        bytes2 a = 0x1234;
        bytes4 b = bytes4(a);
        return b;
    }

    //只有当字节类型和int类型大小相同时，才可以进行转换
    function displayConversion5() public view  returns (uint32,uint32,uint8,uint8) {

        bytes2 a = 0x1234;
        uint32 b = uint16(a); // b = 0x00001234
        uint32 c = uint32(bytes4(a)); // c = 0x12340000
        uint8 d = uint8(uint16(a)); // d = 0x34
        uint8 e = uint8(bytes1(a)); // e = 0x12
        return (b,c,d,e);
    }














}
