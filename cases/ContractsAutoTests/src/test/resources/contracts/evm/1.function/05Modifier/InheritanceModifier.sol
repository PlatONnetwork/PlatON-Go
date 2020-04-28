pragma solidity 0.5.13;

import "./Modifier.sol";

contract InheritanceModifier is Modifier {

    // 使用修饰符 mf
    function inheritance(uint c) public mf(c) {
        a = 1;
    }

    function getA() public view returns (uint) {
        return a;
    }
}