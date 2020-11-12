pragma solidity ^0.5.5;

contract SimpleStorage {


    function hello() public {
        uint256 random = uint256(keccak256(abi.encodePacked(blockhash(block.number-200))));
        assert(random != 0);
    }

    function hash() public view returns(bytes32) {
        return blockhash(block.number-200);
    }

}