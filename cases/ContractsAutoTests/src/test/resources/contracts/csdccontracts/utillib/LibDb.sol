pragma solidity ^0.4.12;
/**
* @file LibDb.sol
* @author lixiaowen
* @time 2017-07-01
* @desc The definition of LibDb library.
*       All functions in this library are extension functions,
*       call eth functionalities through assembler commands.
*
* @funcs
*       writedb(string _name, string _key, string _value) public constant returns(uint _ret)
*
* @usage
*       1) import "LibDb.sol";
*/


import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";

library LibDb {
    using LibInt for *;
    using LibString for *;

    function toWidthString(string _self, uint _width) constant private returns (string _ret) {
        _ret = bytes(_self).length.toString(_width).concat(_self);
    }

    function writedb(string _name, string _key, string _value) public constant returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][07writedb]";
        cmd = cmd.concat(toWidthString(_name, 2));
        cmd = cmd.concat(toWidthString(_key, 4));
        cmd = cmd.concat(toWidthString(_value, 9));

        uint strptr;
        uint strlen = bytes(cmd).length;
        bytes32 hash;
        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = uint(hash);
    }
}