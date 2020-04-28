pragma solidity ^0.5.13;
/**
 * 0.5.13跨合约调用合约一
 *
 * @author hudenian
 * @dev 2020/1/9 18:57
 */
import "./CallerTwo.sol";

contract CallerOne {
    uint256 public x;

    /**
     * 入参为callee合约的地址
     *
     */
    // function inc_delegatecall(address _contractAddress) public {
    function inc_delegatecall() public {
        CallerTwo callerTwo= new CallerTwo();
        address(callerTwo).delegatecall(abi.encodePacked(keccak256("inc()")));
    }

    function getCallerX() public view returns(uint256){
        return x;
    }
}