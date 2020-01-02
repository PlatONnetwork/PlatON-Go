pragma solidity ^0.4.25;
/**
 * 1. 0.4.25版本验证使用constructor关键字定义构造函数，使用internal声明可见性
 * 2. 0.4.25版本验证子合约直接声明父合约构造函数，但是构造函数参数与父合约不一致
 *  如：父合约：constructor(uint _x) 子合约：is Base()
 * @author Albedo
 * @dev 2019/12/23
 **/

contract BaseInternal {
    uint x;
    constructor(uint _x) public {x = _x;}
}
// 验证点：(1) 允许合约直接声明构造函数但参数与父合约不一致————子合约编译可以通过，但是在evm中无法部署
// 0.4.x Base(7) 如果改成 Base()，编译可以通过，但是无法部署；0.5.x Base(7) 如果改成 Base()，直接编译报错；
// --- 增强了编译阶段的语法校验能力
// contract ConstructorInternalVisibility is Base()
// evm返回 This contract may be abstract, not implement an abstract parent's methods completely
// or not invoke an inherited contract's constructor correctly.
contract ConstructorInternalVisibility is BaseInternal(7) {
    uint outI;
    //constructor声明构造函数，允许internal可见性，但部署报错
    constructor(uint _y) {outI = _y;}

    function getOutI() public view returns (uint) {
        return outI;
    }

}
