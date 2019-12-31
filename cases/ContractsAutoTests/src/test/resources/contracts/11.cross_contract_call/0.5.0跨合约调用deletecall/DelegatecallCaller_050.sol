pragma solidity ^0.5.0;
/**
 * 0.5.0跨合约调用的调用者
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract DelegatecallCaller_050 {
    uint256 public x;


    /**
     * 入参为callee合约的地址 0xdc544654fefd1a458eb24064a6c958b14e579154
     * testaddress.delegatecall(bytes4(keccak256("test()")));
     */
    function inc_delegatecall(address _contractAddress) public {
        _contractAddress.delegatecall(abi.encodePacked(keccak256("inc()")));
    }

    function getCallerX() public view returns(uint256){
        return x;
    }
}