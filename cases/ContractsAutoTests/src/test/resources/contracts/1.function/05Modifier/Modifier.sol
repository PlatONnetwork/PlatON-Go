pragma solidity 0.5.13;
/**
 *1.验证单一修饰器
 *2.验证特殊_的用法，符合函数修饰器定义的条件，才可以执行函数体内容
 *3.验证修饰器可以接收参数
 *4.验证合约继承情况下的修饰器的使用
 * @author liweic
 * @dev 2019/12/26 11:10
 */
contract Modifier {
   uint a = 10;

   // 定义修饰符 mf 带参数
   modifier mf (uint b) {
      if (b >= a) {
         _;
      }
   }
}

contract Inheritance is Modifier {

   // 使用修饰符 mf
   function inheritance(uint c) public mf(c) {
      a = 1;
   }

   function getA() public view returns (uint) {
      return a;
   }
}