pragma solidity ^0.5.0;
/**
 * 不允许调用带参数但具有错误参数计数的构造函数。
 * 如果只想在不提供参数的情况下指定继承关系，请不要提供括号
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */


contract  ErrorParamConstructorBase{
    uint public a ;

    //基类构造函数
    constructor (uint _a) public{
      a = _a;
    }
}