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
* 方式二：通过派生的构造函数中用类似于修饰符方式调用父类构造函数（只能2选1）
* (方式一参考mulicPointBaseConstructorWay1.sol)
*/
contract MulicPointBaseConstructorWay2 is ErrorParamConstructorBase {
    constructor(uint _y) errorParamConstructorBase(_y * _y) public {}

    function update(uint amount) public returns (address, uint){
        a += amount;
        return (msg.sender, a);
    }

    //查询a
    function getA() public view returns (uint){
        return a;
    }

}