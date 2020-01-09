pragma solidity ^0.4.12;
/**
* @file LibStack.sol
* @author liaoyan
* @time 2017-07-01
* @desc The defination of LibExtfuncs library.
*       All functions in this library are extension functions,
*       call eth functionalities through assembler commands.
*
* @funcs
*       stackPush(string _data) internal constant returns(bool _ret)
*       stackPop() internal constant returns(string _ret)
*       stackTop() internal constant returns(string _ret)
*       stackSize() internal constant returns(uint _ret)
*       stackClear() internal constant returns(uint _ret)
*       append(string _data) internal constant returns(bool _ret)
*       appendKeyValue(string _key, string _val) internal constant returns (bool _ret)
*       appendKeyValue(string _key, uint _val) internal constant returns (bool _ret)
*       appendKeyValue(string _key, int _val) internal constant returns (bool _ret)
*       appendKeyValue(string _key, address _val) internal constant returns (bool _ret)
*
* @usage
*       1) import "LibStack.sol";
*/


import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";

library LibStack {
    using LibInt for *;
    using LibString for *;

    function push(string _data) internal constant returns(uint _len) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][09StackPush]";
        arg = arg.concat(_data);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        _len = uint(b32);
    }

    function pop() internal constant returns(string _ret) {
        uint i = 0;
        uint stack_size = size();
        while (true) {
            string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][08StackPop]";
            arg = arg.concat(uint(i*32).toString(), "|$%&@*^#!|", uint(32).toString());

            uint argptr;
            uint arglen = bytes(arg).length;

            bytes32 b32;
            assembly {
                argptr := add(arg, 0x20)
                b32 := sha3(argptr, arglen)
            }

            string memory r = uint(b32).recoveryToString();
            _ret = _ret.concat(r);
            if (bytes(r).length < 32 ||  stack_size != size())
                break;

            ++i;
        }
    }

    function popex(uint _len) internal constant returns(string) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][10StackPopEx]";
        string memory result = new string(_len);

        uint strptr;
        assembly {
            strptr := add(result, 0x20)
        }
        cmd = cmd.concat(strptr.toString());

        bytes32 hash;
        uint strlen = bytes(cmd).length;
        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }

        string memory errRet = "";
        uint ret = uint(hash);
        if (ret != 0) {
            return errRet;
        }
        
        return result;
    }

    function top() internal constant returns(string _ret) {
        uint i = 0;
        while (true) {
            string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][08StackTop]";
            arg = arg.concat(uint(i*32).toString(), "|$%&@*^#!|", uint(32).toString());

            uint argptr;
            uint arglen = bytes(arg).length;

            bytes32 b32;
            assembly {
                argptr := add(arg, 0x20)
                b32 := sha3(argptr, arglen)
            }

            string memory r = uint(b32).recoveryToString();
            _ret = _ret.concat(r);
            if (bytes(r).length < 32)
                break;

            ++i;
        }
    }

    function topex(uint _len) internal constant returns(string _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][10StackTopEx]";
        string memory result = new string(_len);

        uint strptr;
        assembly {
            strptr := add(result, 0x20)
        }
        cmd = cmd.concat(strptr.toString());

        bytes32 hash;
        uint strlen = bytes(cmd).length;
        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }

        string memory errRet = "";
        uint ret = uint(hash);
        if (ret != 0) {
            return errRet;
        }
        
        return result;
    }

    function size() internal constant returns(uint _ret) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][09StackSize]";
        arg = arg.concat(""); //don't delete this line

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        return uint(b32);
    }
    
    function clear() internal constant returns(uint _ret) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][10StackClear]";
        arg = arg.concat(""); //don't delete this line

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        return uint(b32);
    }

    function append(string _data) internal constant returns(uint _len) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][11StackAppend]";
        arg = arg.concat(_data);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        _len = uint(b32);
    }

    function appendKeyValue(string _key, string _val) internal constant returns (uint _len) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][19StackAppendKeyValue]";
        arg = arg.concat(_key, "|$%&@*^#!|", _val);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        _len = uint(b32);
    }

    function appendKeyValue(string _key, uint _val) internal constant returns (uint _len) {
        _len = appendKeyValue(_key, _val.toString());
    }

    function appendKeyValue(string _key, int _val) internal constant returns (uint _len) {
        _len = appendKeyValue(_key, _val.toString());
    }

    function appendKeyValue(string _key, address _val) internal constant returns (uint _len) {
        _len = appendKeyValue(_key, uint(_val).toAddrString());
    }
}
