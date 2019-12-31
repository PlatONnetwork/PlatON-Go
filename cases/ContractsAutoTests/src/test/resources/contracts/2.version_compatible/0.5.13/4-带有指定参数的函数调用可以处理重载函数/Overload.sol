pragma solidity ^0.5.13;
/**
 * 带有指定参数的函数调用可以处理重载函数
 * Function calls with named arguments now work with overloaded functions.
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */

contract Overload {

    uint public re;

    function f(uint a,uint b) public returns (uint sum) {
        sum = a+b;
        return sum;
    }

    function f(uint a) public returns (uint sum) {
        sum = a;
        return sum;
    }

    function g() public{
        //0.4.0版本会出错
        uint tp;
        tp = f({a:2,b:3});
        tp = f({a:2});
        re= tp;
        // return re;
    }

    function getRe() view public returns(uint){
        return re;
    }

}