pragma solidity 0.5.13;
/**
 * 验证getter(访问器)函数,编译器会自动为所有 public 状态变量创建 getter 函数,编译器会为我们自动生成的data()函数
 * 在合约内，可以直接操作及访问data状态变量，但在合约外我们只能用data()的方式来访问
 * 在合约内，不能直接访问data()，因为访问器函数的可见性是external
 * @author liweic
 * @dev 2019/12/27 14:10
 */
contract Getter{
  uint public data = 10;

  function f() public returns (uint, uint){
    //分别以internal,external的方式访问
    return (data, this.data());
  }
}