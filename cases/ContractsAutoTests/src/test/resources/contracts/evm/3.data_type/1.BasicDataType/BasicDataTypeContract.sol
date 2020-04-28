pragma solidity 0.5.2;

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

  /**
   * 1、整型 (1个数字1个字节代表8位)
        1)、int有符号整型，默认int256，包含位数：int8、int16、int24.....int256
        2)、uint：无符号整数，默认uint256，包含位数：uint8、uint16,uint24.....uint256
   */

    //1)、验证有符号/无符号整数
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

  /**
   * 2、浮点型
        可以用来声明变量，但是不可以用来赋值
      1)、fixed 有符号固定位浮点数，关键字为 fixedMxN
      2)、ufixed:无符号的固定位浮点数，关键字为 ufixedMxN
      （M 表示这个类型要占用的位数，以 8 步进，可为 8 到 256 位。N 表示小数点的个数，可为 0 到 80 之间。）
   *     
   */ 

  //1、fixed 有符号固定位浮点数，关键字为 fixedMxN
  //fixed f1 = 1.0;//编译异常

    /**
    *3、布尔型：取值为常量值true和false
    *  
    */
    bool h = true;
    bool j = false;

    /**
     *4、字节数组
     *1)、定长字节数组，关键字：bytes1、bytes2、bytes3...，byte32，byte默认代表bytes1。
     *2)、变长字节数组，关键字：bytes，动态分配大小字节数组
     */
      
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