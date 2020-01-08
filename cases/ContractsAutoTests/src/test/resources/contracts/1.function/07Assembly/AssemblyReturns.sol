pragma solidity 0.5.13;

contract AssemblyReturns {
    uint constant a = 2;
    bytes2 constant b = 0xabcd;
    bytes3 constant c = "abc";
    bool constant d = true;
    address payable constant e = 0x1212121212121212121212121212121212121212;
    function f() public pure returns (uint w, bytes2 x, bytes3 y, bool z, address t) {
        assembly {
            w := a
            x := b
            y := c
            z := d
            t := e
        }
    }
}