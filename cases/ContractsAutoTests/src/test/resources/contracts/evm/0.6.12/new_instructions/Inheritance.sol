// SPDX-License-Identifier: MIT
pragma solidity ^0.6.12;

/* Graph of inheritance
    A
   / \
  B   C
 / \ /
F  D,E

*/

contract InheritanceA {
    function foo() public pure virtual returns (string memory) {
        return "InheritanceA";
    }
}

// Contracts inherit other contracts by using the keyword 'is'.
contract InheritanceB is InheritanceA {
    // Override A.foo()
    function foo() public pure virtual override returns (string memory) {
        return "InheritanceB";
    }
}

contract InheritanceC is InheritanceA {
    // Override A.foo()
    function foo() public pure virtual override returns (string memory) {
        return "InheritanceC";
    }
}

// Contracts can inherit from multiple parent contracts.
// When a function is called that is defined multiple times in
// different contracts, parent contracts are searched from
// right to left, and in depth-first manner.

contract InheritanceD is InheritanceB, InheritanceC {
    // D.foo() returns "C"
    // since C is the right most parent contract with function foo()
    function foo() public pure override(InheritanceB, InheritanceC) returns (string memory) {
        return super.foo();
    }
}

contract InheritanceE is InheritanceC, InheritanceB {
    // E.foo() returns "B"
    // since B is the right most parent contract with function foo()
    function foo() public pure override(InheritanceC, InheritanceB) returns (string memory) {
        return super.foo();
    }
}

// Inheritance must be ordered from “most base-like” to “most derived”.
// Swapping the order of A and B will throw a compilation error.
contract InheritanceF is InheritanceA, InheritanceB {
    function foo() public pure override(InheritanceA, InheritanceB) returns (string memory) {
        return super.foo();
    }
}