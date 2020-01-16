pragma solidity ^0.5.0;
/**
 * 0.5.0不允许在同一继承层次结构中多次指定基类构造函数参数
 * 0.4.x可以同时使用2种方式，但如果2种方式都存在，优先选择修饰符方式
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */


import "./ErrorParamConstructorBase.sol";

/**
* 方式一：contract mulicPointBaseConstructor is errorParamConstructorBase(10)
* (方式二参考mulicPointBaseConstructorWay2.sol)
*/
contract MulicPointBaseConstructorWay1 is ErrorParamConstructorBase(10) {
    constructor() public {}

    function update(uint amount) public returns (address, uint){
        a += amount;
        return (msg.sender, a);
    }
}