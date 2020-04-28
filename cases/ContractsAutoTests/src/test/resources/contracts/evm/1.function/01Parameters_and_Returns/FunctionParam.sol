pragma solidity 0.5.13;

/**
 * 验证入参是函数的使用,来自官方case
 * @author liweic
 * @dev 2020/01/11 20:09
 */

library G {
    function g(function() internal returns (uint) _t) internal returns (uint) {
        return _t();
    }
}
contract FunctionParam {
    using G for *;
    function g(function() internal returns (uint) _t) internal returns (uint) {
        return _t();
    }

    function f() public returns (uint) {
        return t.g();
    }
    function t() public pure returns (uint)  {
        return 7;
    }
}