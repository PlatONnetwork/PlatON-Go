pragma solidity ^0.4.25;

/**
 * 1. 0.4.25版本验证使用constrictor关键字定义构造函数，但是不强制声明可见性(默认为public可见性）
 * 2. 0.4.25版本同一继承层次结构中允许多次指定基类构造函数参数验证:
 *  (1) 允许合约直接声明构造函数 ———— is Base(7)
 * （2）子合约构造函数继承父合约构造函数———— constructor(uint _y) Base(_y * _y)
 * 两种引用构造函数方式共存时，合约优先选择（2）方式
 * @author Albedo
 * @dev 2019/12/23
 **/
contract BaseDefault {
    uint x;
    constructor(uint _x) public { x = _x; }
}
//(1) 允许合约直接声明构造函数Base(7)
contract ConstructorDefaultVisibility is BaseDefault(7) {
    uint outI;
    //（2）子合约构造函数继承父合约构造函数———— constructor(uint _y) Base(_y * _y)
    constructor(uint _y) BaseDefault(_y * _y) {outI=_y;}

    function getOutI() public view returns (uint) {
        return outI;
    }

}
