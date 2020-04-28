pragma solidity ^0.5.13;
/**
 * 10-链上函数
 * 函数 assert()
 *
 * @author hudenian
 * @dev 2020/1/8 09:57
 *
 */
contract AssertContract {

    address payable owner;
    uint result;

    constructor() public{
        owner = msg.sender;
    }
    /**
     * revert会把所有剩下的gas退回给caller
     * 验证输入参数合法性
     */
    function toSenderAmount(uint frist,uint second) public returns(uint) {
        assert(frist >  second);

        result = frist+second;
        return result;
    }


    function getResult() view public returns(uint){
        return result;
    }
}