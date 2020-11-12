// SPDX-License-Identifier: MIT
pragma solidity ^0.7.1;

contract ShadowingInheritedStateVariablesSuper {
    string public name = "Contract ShadowingInheritedStateVariablesSuper";

    function getName() public view returns (string memory) {
        return name;
    }
}

contract ShadowingInheritedStateVariablesSuperChild is ShadowingInheritedStateVariablesSuper {
    // This is the correct way to override inherited state variables.
    constructor() public {
        name = "Contract ShadowingInheritedStateVariablesSuperChild";
    }
}