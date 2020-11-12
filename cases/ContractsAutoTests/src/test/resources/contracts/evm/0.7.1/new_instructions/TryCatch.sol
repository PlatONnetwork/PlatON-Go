// SPDX-License-Identifier: MIT
pragma solidity ^0.7.1;

contract CalledContract {
    function someFunction() external pure{
        // Code that reverts
        revert();
    }
}


contract TryCatcher {

    event CatchEvent();
    event SuccessEvent();

    CalledContract public externalContract;

    constructor() public {
        externalContract = new CalledContract();
    }

    function execute() external {

    try externalContract.someFunction() {
        // Do something if the call succeeds
    emit SuccessEvent();
    } catch {
        // Do something in any other case
    emit CatchEvent();
    }
    }
}