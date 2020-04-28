pragma solidity ^0.5.0;
/**
 * 验证0.5.0版本构造函数可见性必须显示声明public或者internal
 *
 * @author hudenian
 * @dev 2019/12/20 11:09
 */


contract ConstructorPublicDeclaraction {

    uint count = 0;

    //构造函数必须显示声明为public(0.4.x可以不显式声明或者用同名构造函数)
    constructor(uint _count) public {
        count = _count;
    }

    function update(uint amount) public returns (address, uint){
        count += amount;
        return (msg.sender, count);
    }

    //5.查询balance
    function getCount() public view returns (uint){
        return count;
    }
}