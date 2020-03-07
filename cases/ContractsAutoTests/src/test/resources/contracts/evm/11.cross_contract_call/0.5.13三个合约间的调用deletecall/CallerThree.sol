pragma solidity ^0.5.13;
/**
 * 0.5.13跨合约被调用者合约三
 *
 * @author hudenian
 * @dev 2020/1/9 18:57
 */

contract CallerThree {
    uint256 public x;

    event EventName(address seder,uint256 x);

    function inc() public  {
        x++;
        emit EventName(msg.sender,x);
    }

    function getCalleeThreeX() public view returns(uint256){
        return x;
    }
}