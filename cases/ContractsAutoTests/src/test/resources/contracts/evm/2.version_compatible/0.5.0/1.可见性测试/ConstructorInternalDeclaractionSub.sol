pragma solidity ^0.5.0;
/**
 * 验证0.5.0版本构造函数可见性必须显示声明public或者internal
 * https://github.com/ethereum/solidity/issues/2638(internal类型主要用于子合约使用)
 *
 * @author hudenian
 * @dev 2019/12/20 11:09
 */


import "./ConstructorInternalDeclaraction.sol";

contract  ConstructorInternalDeclaractionSub is  ConstructorInternalDeclaraction{

    //构造函数必须显式声明为internal(0.4.x可以不显式声明或者用同名构造函数)
    constructor(uint _count) ConstructorInternalDeclaraction(_count) public{}


    function update(uint amount) public returns (address, uint){
        count += amount;
        return (msg.sender, count);
    }

    //5.查询count
    function getCount() public view returns (uint){
        return count;
    }

}
