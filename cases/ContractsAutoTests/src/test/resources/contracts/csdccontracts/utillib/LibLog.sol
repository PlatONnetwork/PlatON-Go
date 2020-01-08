pragma solidity ^0.4.12;
/**
* @file LibLog.sol
* @author lixiaowen
* @time 2017-07-01
* @desc The defination of LibLog library.
*       All functions in this library are extension functions,
*       call eth functionalities through assembler commands.
*
* @funcs
*       log(string _str) internal constant returns(uint _ret)
*       log(string _str, string _str2) internal constant returns(uint _ret)
*       log(string _str, string _str2, string _str3) internal constant returns(uint _ret)
*       log(string _str, uint _ui) internal constant returns(uint _ret)
*       log(string _str, int _i) internal constant returns(uint _ret)
*       log(string _str, address _addr) internal constant returns(uint _ret)
*
* @usage
*       1) import "LibLog.sol";
*/


import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";

library LibLog {
    using LibInt for *;
    using LibString for *;

    function log(string _str) internal constant returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat(_str);

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
    }

    function log(string _str, string _str2) internal constant returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat(_str, " ", _str2);
 
        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
    }
    
    function log(string _str, string _str2, string _str3) internal constant returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat(_str, " ", _str2);
        cmd = cmd.concat(" ", _str3);

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
    }
    
    function log(string _str, uint _ui) internal constant returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat(_str, " ", _ui.toString());

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
    }
    
    function log(string _str, int _i) internal constant returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat(_str, " ", _i.toString());

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }

        _ret = 0;
    }

    function log(string _str, address _addr) internal constant returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat(_str, " ", uint(_addr).toAddrString());

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
    }
}
