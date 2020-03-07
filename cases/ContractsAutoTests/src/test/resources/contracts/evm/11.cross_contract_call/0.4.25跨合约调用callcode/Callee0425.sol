pragma solidity ^0.4.25;
/**
 * 0.4.25跨合约被调用者
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */


contract Callee0425 {
    uint256 public x;

    event EventName(address seder,uint256 x);

    function inc() public {
        x++;
        emit EventName(msg.sender,x);
    }

    function getCalleeX() public view returns(uint256){
        return x;
    }
}