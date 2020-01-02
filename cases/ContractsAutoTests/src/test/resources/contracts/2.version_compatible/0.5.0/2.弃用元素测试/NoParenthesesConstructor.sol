pragma solidity ^0.5.0;
/**
 * 不允许调用没有括号的基类构造函数
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

 1
import "./ErrorParamConstructorBase.sol";

contract Base {
    uint x;
    constructor() public {}
}

/**
 * 0.5.0版本不允许调用没有括号的基类构造函数（编译不通过）
 * 0.4.25版本可以调用基类构造函数
 */
contract NoParenthesesConstructor is ErrorParamConstructorBase {
    constructor(uint _b) errorParamConstructorBase(1) public {}
}

//contract noParenthesesConstructor is errorParamConstructorBase {
//    constructor(uint _b) errorParamConstructorBase() public {}
//}
