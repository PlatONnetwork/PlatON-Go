pragma solidity 0.5.13;
/**
 1.验证合约关键字this,表示当前合约，可以显示的转换为Address
 2.验证Fallback函数,调用了未命名的函数等方式
 * @author liweic
 * @dev 2019/12/26 15:30
 */

contract FallBack {
    uint a = 1;

    //定义一个回退函数
    function () external {
        a = 100; 
    }
    
    //调用一个不存在的函数将触发回退函数
    function CallFunctionNotExist() public {
        address(this).delegatecall("functionNotExist()");
    }
    
    function getA() view public returns (uint) {
        return a;
    }
}
