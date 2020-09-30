pragma solidity ^0.5.17;

contract Instructions {

    function getChainId() public view returns (uint) {
        uint256 chainId;
        assembly { chainId := chainid() }
        return chainId;
    }

    function getSelfBalance() public view returns (uint) {
        uint256 ret;
        assembly { ret := selfbalance() }
        return ret;
    }

    function test(uint x, uint y) public returns (uint) {
        return test_mul(2,3);
     }

    function test_mul(uint x, uint y) public returns (uint) {
        return multiply(x,y);
    }

    function multiply(uint x, uint y) public returns (uint) {
        return x * y;
    }
}