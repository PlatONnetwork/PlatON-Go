pragma solidity ^0.5.13;

/**
 * 安全算术运算
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
library SafeMathLibrary {

    /**
     * 获取最大值
     */
    function max(uint256 a, uint256 b) internal pure returns (uint256) {
        return a >= b ? a : b;
    }

    /**
     * 获取最小值
     */
    function min(uint256 a, uint256 b) internal pure returns (uint256) {
        return a < b ? a : b;
    }

    /**
     * 获取平均值
     */
    function average(uint256 a, uint256 b) internal pure returns (uint256) {
        return (a / 2) + (b / 2) + ((a % 2 + b % 2) / 2);
    }
    /**
     * .两个无符号整型相加
     */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a, "SafeMath: addition overflow");
        return c;
    }

    /**
     * 两个无符号整型相减
     */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a, "SafeMath: subtraction overflow");
        uint256 c = a - b;

        return c;
    }


    /**
     * 两个无符号整型相乘
     */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b, "SafeMath: multiplication overflow");

        return c;
    }

    /**
     * 两个无符号整型相除，并检测除数为0异常
     * a 被除数
     * b 除数
     */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b > 0, "SafeMath: division by zero");
        uint256 c = a / b;
        return c;
    }

    /**
     * 两个无符号整型相除，并自定义除数为0异常信息
     * a 被除数
     * b 除数
     */
    function div(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        require(b > 0, errorMessage);
        uint256 c = a / b;
        return c;
    }

    /**
     * 两个无符号整型除余运算，并检测除数为0异常
     * a 被除数
     * b 除数
     */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b != 0, "SafeMath: modulo by zero");
        return a % b;
    }

    /**
      * 两个无符号整型除余运算，并自定义除数为0异常信息
      * a 被除数
      * b 除数
      */
    function mod(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        require(b != 0, errorMessage);
        return a % b;
    }
}