pragma solidity ^0.5.0;
/**
 * 0.5.0跨合约被调用者
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract DelegatecallCallee_050 {
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