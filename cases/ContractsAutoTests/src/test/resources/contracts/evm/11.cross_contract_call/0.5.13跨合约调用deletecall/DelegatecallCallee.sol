pragma solidity ^0.5.13;
/**
 * 0.5.13跨合约被调用者
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract DelegatecallCallee {
    uint256 public x;

    event EventName(address seder,uint256 x);

    function inc() public  {
        x++;
        emit EventName(msg.sender,x);
    }

    function getCalleeX() public view returns(uint256){
        return x;
    }
}