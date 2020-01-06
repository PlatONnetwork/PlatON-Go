pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 *
 *测试引用数据类型
 *1、数组（Array）
 *2、结构体（Struct）
 */

contract ReferenceDataTypeArrayContract {

   /**
    *2、数组（Array）：一种数据结构，它是存储同类元素的有序集合
    * 1)、数组的声明及初始化及取值(定长数组、可变数组)
    * 2)、多维数组
    * 3)、数组成员方法
    *     length属性、push()方法
    */

   /**
    *验证：1)、数组的声明及初始化及取值(定长数组、可变数组)
    *
    */

    //声明定长数组
    uint[5]  numArray = [1,2,3,4,5];

    //声明可变数组
    string[] numArray1 = ["1","2","3","4","5","6"];

    //new关键字创建动态数组
    uint[] a = new uint[](5);

    function setArray(uint index,uint value) public {
         numArray[index] = value;
    }

    function getArray(uint index) public view returns (uint value) {
          return numArray[index];
    }

  /**
    *验证：2)、多维数组声明及初始化及取值
    *
    */
    //声明二维数组并赋值、取值
    uint[][]  multiArray = [[0,0],[0,1],[0,2],[1,0],[1,1],[1,2]];

    function setMultiArray() public {
       multiArray[1][0] = 100;
    }

    function getMultiArray() public view returns (uint value,uint lenght) {
      return (multiArray[1][0],multiArray[0].length);
    }

  /**
    *验证：3)、数组的属性及方法
    *  length属性、push()方法
    */
    function setArrayPush(string memory x) public {
         numArray1.push(x);
    }

    function getArrayLength() public view returns (uint) {
          return numArray1.length;
    }
}