pragma solidity >=0.4.26;

contract Eventer {

    event TestInt8(int8 indexed out1, int8 indexed out2);
    event AnonEvent(address, address);

    function getEvent() public {
        // set to 2,3 for functioning filter
        emit TestInt8(-2, -3);
    }
}