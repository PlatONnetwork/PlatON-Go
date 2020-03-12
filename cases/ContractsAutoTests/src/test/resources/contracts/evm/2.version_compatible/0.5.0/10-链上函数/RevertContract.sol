pragma solidity ^0.5.13;
/**
 * 10-链上函数
 * 函数 revert()
 *
 * @author hudenian
 * @dev 2020/1/8 09:57
 *
 */
contract RevertContract {

    address payable owner;
    uint result;

    constructor() public{
        owner = msg.sender;
    }
    /**
     * revert会把所有剩下的gas退回给caller
     * 验证输入参数合法性
     */
    function toSenderAmount(uint first,uint second) public returns(uint) {
        if(first < second){
            revert("first less than second");
        }
        result = first-second;
        return result;
    }


    function getResult() view public returns(uint){
        return result;
    }
}