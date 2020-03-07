pragma solidity ^0.5.13;
/**
 * 0.5.13跨合约被调用者合约二
 *
 * @author hudenian
 * @dev 2020/1/9 18:57
 */

import "./CallerThree.sol";

contract CallerTwo {
    uint256 public x;

    event EventName(address seder,uint256 x);

    function inc() public  {
        CallerThree callerThree = new CallerThree();
        address(callerThree).delegatecall(abi.encodePacked(keccak256("inc()")));
        emit EventName(msg.sender,x);
    }

    function getCalleeX() public view returns(uint256){
        return x;
    }
}
