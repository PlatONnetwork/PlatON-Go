pragma solidity ^0.4.25;
/**
 * 0.4.25跨合约调用的调用者
 * 说明：DELEGATECALL会一直使用原始调用者的地址，而CALLCODE不会。
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */


contract Caller0425 {
    uint256 public x;

    /**
     * 改变的是被调用者的状态变量
     */
    function inc_call(address _contractAddress) public {
        _contractAddress.call(bytes4(keccak256("inc()")));
    }

    /**
     * 改变的是调用者的状态变量，msg.sender是当前合约
     */
    function inc_callcode(address _contractAddress) public {
        _contractAddress.callcode(bytes4(keccak256("inc()")));
    }

    /**
     * 改变的是调用者的状态变量，msg.sender是最初调用者
     */
    function inc_delegatecall(address _contractAddress) public {
        _contractAddress.delegatecall(bytes4(keccak256("inc()")));
    }

    function getCallerX() public view returns(uint256){
        return x;
    }
}