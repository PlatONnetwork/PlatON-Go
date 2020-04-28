pragma solidity ^0.5.13;
/**
 * 0.5.13跨合约调用的调用者
 * delegatecall将会改变delegatecaller中的x值，而不会改变delegatecallCallee中的x值
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */


contract DelegatecallCaller {
    uint256 public x;

    /**
     * 入参为callee合约的地址 0xdc544654fefd1a458eb24064a6c958b14e579154
     *
     */
    function inc_delegatecall(address _contractAddress) public {
        _contractAddress.delegatecall(abi.encodePacked(keccak256("inc()")));
    }

    function getCallerX() public view returns(uint256){
        return x;
    }
}