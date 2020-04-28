pragma solidity ^0.5.0;
/**
 * 08-合约和地址
 * 验证msg.sender属于address payable类型
 *
 * @author hudenian
 * @dev 2019/12/24 09:57
 */


contract MsgSenderBelongToPayable {
    function getContractAddr() view public returns (address){
        address payable  ad= msg.sender;//msg.sender属于address payable
        return address(this);
    }

    function getMsgSenderAddr() view public returns (address){
        return msg.sender;
    }
}

