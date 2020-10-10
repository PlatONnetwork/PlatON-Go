pragma solidity ^0.7.1;

/**
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约基本数据类型
 *1、整型
 *2、浮点型
 *3、布尔型
 *4、字节数组
 */

contract BasicDataTypeContract {

    int a = -1;//正常
    int b = 1;//正常
    uint c = 2;//正常
    //uint d = -3;//异常，编译异常，无符号整数不能赋值有符号整数

    //2)、验证无符号整数位数
    uint8 x1 = 1;
    uint8 x2 = 255;//8位无符号整数取值范围0~255
    //uint8 x3 = 256;//8位数无符号整数溢出，编译报错
    uint16 y1 = 1;
    uint16 y2 = 65535;//16位无符号整数取值范围0~65535
    //uint16 y3 = 65535;//16位数溢出，编译报错

    //无符号8位整数溢出(输入最大值255+1，则溢出显示结果为0)
    function addUintOverflow(uint8 a) public pure returns(uint8) {
        return a + 1;
    }

    //3)、验证有符号整数位数
    int8 z1 = 1;
    int8 z2 = 127 ;//8位有符号整数取值范围-128~127
    //int8 z3 = 128;//8位数有符号整数溢出，编译报错
    int8 w1 = -1;
    int8 w2 = -128;
    //int8 w3 = -129;//8位数有符号整数溢出，编译报错

     //有符号8位整数溢出(输入最大值127+1，则溢出显示结果-128)
    function addIntOverflow(int8 a) public pure returns(int8) {
        return a + 1;
    }



  //1、fixed 有符号固定位浮点数，关键字为 fixedMxN
  //fixed f1 = 1.0;//编译异常

    bool h = true;
    bool j = false;
      
      //1)、验证定长字节数组
      bytes1 b1 = "a";
      bytes1 b2 = bytes1(uint8(1));
      bytes2 b3 = "ab";
      bytes3 b4 = "abc";
      //bytes3 b5 = "abcd";//位数超出，编译异常

      function getBytes1Length() public view returns (uint) {
           return  b3.length;
      }
      
      //2)、验证变长字节数组
      bytes k1 = "a";
      bytes k2 = "ab";
      bytes k3 = "abc";

      function getBytesLength() public view returns (uint) {
           return  k3.length;
      }
}