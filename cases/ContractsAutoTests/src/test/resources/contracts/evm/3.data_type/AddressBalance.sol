pragma solidity ^0.4.12;

contract AddressBalance {
    function balanceOfPlatON(address user) public constant returns (uint256) {
        return user.balance;
    }
}
