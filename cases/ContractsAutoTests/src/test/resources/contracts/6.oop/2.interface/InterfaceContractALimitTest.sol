pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约接口功能点
 *接口(interface)简述：接口类似于抽象合约，但是其不能实现任何函数。
 *-----------------  测试点   ------------------------------
 *1、接口限制点测试(5.0以后版本)
 *1)、 不能声明构造函数
 *2)、 不能声明状态变量
 *3)、 接口的函数只能声明外部类型(external)
 */

interface InterfaceContractParent  {

   /**
    * 验证：1、接口不能声明状态变量
    *-----------------------------
    * 验证结果：接口中定义状态变量，会编译失败
    * Variables cannot be declared in interfaces.solc
    */
     //  uint a = -1;

   /**
    * 验证：2、接口不能声明构造函数
    *-----------------------------
    * 验证结果：接口中定义构造函数，会编译失败
    */
     /*constructor() public{
          a += 2;
     }*/

     
   /**
    * 验证：3、接口的函数只能声明外部类型(external)
    *------------------------------
    * 验证结果：在5.0以后版本，接口的函数只能声明外部类型(external)，否则会编译失败；
    *          在5.0以前的版本，使用public
    */

   //5.0以后版本
   function sumExternal(int a,int b)  external view  returns (int);

   //5.0以前版本,否则Error: Functions in interfaces must be declared external
   //function sumPublic(int a,int b) public returns (int);
}

contract InterfaceContractParentTest is InterfaceContractParent {

   function sumExternal(int a,int b)  external view returns (int) {
      return a + b;
   }
  
}







