pragma solidity 0.5.13;
/**
 * 验证public (或 external) 函数有一个特殊的成员selector, 它对应一个ABI 函数选择器.
 * @author liweic
 * @dev 2020/01/11 20:30
 */
contract Selector {
    function h() payable external {
    }
    function f() view external returns (bytes4) {
        function () external g = this.h;
        return g.selector;
    }
}