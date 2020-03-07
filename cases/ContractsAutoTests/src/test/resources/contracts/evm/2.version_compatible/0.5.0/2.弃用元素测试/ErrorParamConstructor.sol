pragma solidity ^0.5.0;
/**
 * 不允许调用带参数但具有错误参数计数的构造函数。
 * 如果只想在不提供参数的情况下指定继承关系，请不要提供括号
 *
 * @author hudenian
 * @dev 2019/12/20 09:57
 */


import "./ErrorParamConstructorBase.sol";


/**
 * 0.5.0版本可以不指定参数errorParamConstructorBase编译可以通过但是部署失败，
 * 指定参数errorParamConstructorBase(10)编译部署成功
 * 0.4.x版本可以不指定参数，并且带括号如：errorParamConstructorBase()
 */
contract ErrorParamConstructor is ErrorParamConstructorBase(10) {
    uint public b ;
    constructor(uint _b) public {
        b = _b;
    }

    function update(uint amount) public returns (address, uint){
        b += amount;
        return (msg.sender, b);
    }

    //查询a
    function getA() public view returns (uint){
        return a;
    }

    //查询b
    function getB() public view returns (uint){
        return b;
    }
}