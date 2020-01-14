pragma solidity ^0.5.13;

import "./SafeMathLibrary.sol";
/**
 * 安全算术运算验证
 *（1）最大值
 *（2）最小值
 *（3）平均值
 *（4）无符号整型相加
 *（5）无符号整型相减
 *（6）无符号整型相乘
 *（7）无符号整型相除
 *（8）无符号整型取模运算
 *
 * @author Albedo
 * @dev 2020/01/08
 **/
contract SafeMathMock {
    //最大值验证
    function max(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.max(a, b);
    }
    //最小值验证
    function min(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.min(a, b);
    }
    //平均值验证
    function average(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.average(a, b);
    }
    //乘运算
    function mul(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.mul(a, b);
    }
    //除运算
    function div(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.div(a, b);
    }
    //减法运算
    function sub(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.sub(a, b);
    }
    //加法运算
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.add(a, b);
    }
    //取模运算
    function mod(uint256 a, uint256 b) public pure returns (uint256) {
        return SafeMathLibrary.mod(a, b);
    }
}