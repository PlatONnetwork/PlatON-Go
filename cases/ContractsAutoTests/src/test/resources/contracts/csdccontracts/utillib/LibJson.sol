pragma solidity ^0.4.12;
/**
* @file LibJson.sol
* @author liaoyan
* @time 2017-07-01
* @desc The definition of LibJson library.
*       All functions in this library are extension functions,
*       call eth functionalities through assembler commands.
*
* @funcs
*       isJson(string _json) internal constant returns(bool _ret)
*       jsonRead(string _json, string _keyPath) internal constant returns(string _ret)
*       jsonKeyExists(string _json, string _keyPath) internal constant returns(bool _ret)
*       jsonUpdate(string _json, string _key, string _value) internal constant returns (string _ret)
*
* @usage
*       1) import "LibJson.sol";
*       2) using LibJson for *;
*/


import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibLog.sol";

library LibJson {
    using LibInt for *;
    using LibString for *;
    using LibJson for *;

    function isJson(string _json) internal constant returns(bool _ret) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][09JsonParse]";

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        if (uint(b32) != 0)
            return true;
        else
            return false;
    }

    /* function jsonRead(string _json, string _keyPath) internal constant returns(string _ret) {
        uint i = 0;
        while (true) {
            string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][08JsonRead]";
            arg = arg.concat(_keyPath, "|$%&@*^#!|");
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
    } */

    function jsonRead(string _json, string _keyPath) internal constant returns(string _ret) {
        _ret = "";
        uint len = jsonKeyExistsEx(_json, _keyPath);
        if (len > 0) {
            _ret = new string(len);
            uint strptr;
            assembly {
                strptr := add(_ret, 0x20)
            }

            string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][10JsonReadEx]";
            arg = arg.concat(_keyPath);
            arg = arg.concat("|$%&@*^#!|", strptr.toString());

            uint argptr;
            uint arglen = bytes(arg).length;
            bytes32 b32;
            assembly {
                argptr := add(arg, 0x20)
                b32 := sha3(argptr, arglen)
            }

            string memory errRet = "";
            uint ret = uint(b32);
            if (ret != 0) {
                return errRet;
            }
        }
    }

    function jsonKeyExists(string _json, string _keyPath) internal constant returns(bool _ret) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][13JsonKeyExists]";
        arg = arg.concat(_keyPath);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;
        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        if (uint(b32) != 0)
            return true;
        else
            return false;
    }

    function jsonKeyExistsEx(string _json, string _keyPath) internal constant returns(uint _len) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][15JsonKeyExistsEx]";
        arg = arg.concat(_keyPath);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;
        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        _len = uint(b32);
    }

    function jsonArrayLength(string _json, string _keyPath) internal constant returns(uint _ret) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][15jsonArrayLength]";
        arg = arg.concat(_keyPath);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;
        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        return uint(b32);
    }

    function jsonUpdate(string _json, string _keyPath, string _value) internal constant returns (string _ret) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][10JsonUpdate]";
        arg = arg.concat(_keyPath, "|$%&@*^#!|", _value);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;
        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        if (uint(b32) != 0)
            //return lastJson();
            return lastJsonEx();
        else
            return "";
    }

    function lastJson() internal constant returns(string _ret) {
        uint i = 0;
        while (true) {
            string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][08LastJson]";
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

    function lastJsonEx() internal constant returns(string _ret) {
        _ret = "";
        uint len = lastJsonLength();
        if (len > 0) {
            _ret = new string(len);
            uint strptr;
            assembly {
                strptr := add(_ret, 0x20)
            }

            string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][10LastJsonEx]";
            arg = arg.concat(strptr.toString());

            uint argptr;
            uint arglen = bytes(arg).length;
            bytes32 b32;
            assembly {
                argptr := add(arg, 0x20)
                b32 := sha3(argptr, arglen)
            }

            string memory errRet = "";
            uint ret = uint(b32);
            if (ret != 0) {
                return errRet;
            }
        }
    }

    function lastJsonLength() internal constant returns(uint _len) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][14LastJsonLength]";
        arg = arg.concat(""); //don't delete this line

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;
        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        _len = uint(b32);
    }

    function jsonCat(string _self, string _str) internal constant returns (string _ret) {
        _ret = _self.trim().concat(_str.trim());
    }

    function jsonCat(string _self, string _key, string _val) internal constant returns (string _ret) {
        _ret = _self.trim();
        if (bytes(_ret).length > 0 && (bytes(_ret)[bytes(_ret).length-1] != '{' && bytes(_ret)[bytes(_ret).length-1] != '[')) {
            _ret = _ret.concat(",");
        }

        bool isJson = false;
        if (bytes(_val).length > 0 && bytes(_val)[0] == '{' && bytes(_val)[bytes(_val).length-1] == '}') {
            isJson = true;
        }
        if (bytes(_val).length > 0 && bytes(_val)[0] == '[' && bytes(_val)[bytes(_val).length-1] == ']') {
            isJson = true;
        }

        if (isJson) {
            _ret = _ret.concat("\"", _key, "\":");
            _ret = _ret.concat(_val);
        } else {
            _ret = _ret.concat("\"", _key, "\":");
            _ret = _ret.concat("\"", _val, "\"");
        }

        "}"; //compiler bug
    }

    function jsonCat(string _self, string _key, uint _val) internal constant returns (string _ret) {
        _ret = _self.trim();
        if (bytes(_ret).length > 0 && (bytes(_ret)[bytes(_ret).length-1] != '{' && bytes(_ret)[bytes(_ret).length-1] != '[')) {
            _ret = _ret.concat(",");
        }

        _ret = _ret.concat("\"", _key, "\":");
        _ret = _ret.concat(_val.toString());

        "}"; //compiler bug
    }

    function jsonCat(string _self, string _key, int _val) internal constant returns (string _ret) {
        _ret = _self.trim();
        if (bytes(_ret).length > 0 && (bytes(_ret)[bytes(_ret).length-1] != '{' && bytes(_ret)[bytes(_ret).length-1] != '[')) {
            _ret = _ret.concat(",");
        }

        _ret = _ret.concat("\"", _key, "\":");
        _ret = _ret.concat(_val.toString());

        "}"; //compiler bug
    }

    function jsonCat(string _self, string _key, address _val) internal constant returns (string _ret) {
        _ret = _self.trim();
        if (bytes(_ret).length > 0 && (bytes(_ret)[bytes(_ret).length-1] != '{' && bytes(_ret)[bytes(_ret).length-1] != '[')) {
            _ret = _ret.concat(",");
        }

        _ret = _ret.concat("\"", _key, "\":");
        _ret = _ret.concat("\"", uint(_val).toAddrString(), "\"");

        "}"; //compiler bug
    }

    function toJsonArray(uint[] storage _self) internal constant returns(string _json) {
        _json = _json.concat("[");
        for (uint i=0; i<_self.length; ++i) {
            if (i == 0)
                _json = _json.concat(_self[i].toString());
            else
                _json = _json.concat(",", _self[i].toString());
        }
        _json = _json.concat("]");
    }

    function toJsonArray(string[] storage _self) internal constant returns(string _json) {
        _json = _json.concat("[");
        for (uint i=0; i<_self.length; ++i) {
            if (i == 0)
                _json = _json.concat("\"", _self[i], "\"");
            else
                _json = _json.concat(",\"", _self[i], "\"");
        }
        _json = _json.concat("]");
    }

    function fromJsonArray(uint[] storage _self, string _json) internal returns(bool succ) {
        if(push(_json) == 0) {
            return false;
        }
        _self.length = 0;

        if (!isJson(_json)) {
            pop();
            return false;
        }

        while (true) {
            string memory key = "[".concat(_self.length.toString(), "]");
            if (!jsonKeyExists(_json, key))
                break;

            _self.length++;
	    //_self[_self.length-1] = jsonRead(_json, key).toUint();
            _self[_self.length-1] = jsonRead(_json, key).toUint();
        }

        pop();
        return true;
    }

    function fromJsonArray(string[] storage _self, string _json) internal returns(bool succ) {
        if(push(_json) == 0) {
            return false;
        }
        _self.length = 0;

        if (!isJson(_json)) {
            pop();
            return false;
        }

        while (true) {
            string memory key = "[".concat(_self.length.toString(), "]");
            if (!jsonKeyExists(_json, key))
                break;

            _self.length++;
	    //_self[_self.length-1] = jsonRead(_json, key);
            _self[_self.length-1] = jsonRead(_json, key);
        }

        pop();
        return true;
    }
    
    //old libJSON
    // ???JSON
    function getObjectValueByKey(string _self, string _key) internal returns (string _ret) {
        int pos = -1;
        uint searchStart = 0;
        while (true) {
            pos = _self.indexOf("\"".concat(_key, "\""), searchStart);
            if (pos == -1) {
                pos = _self.indexOf("'".concat(_key, "'"), searchStart);
                if (pos == -1) {
                    return;
                }
            }

            pos += int(bytes(_key).length+2);
            // pos ??????{
            bool colon = false;
            while (uint(pos) < bytes(_self).length) {
                if (bytes(_self)[uint(pos)] == ' ' || bytes(_self)[uint(pos)] == '\t' 
                    || bytes(_self)[uint(pos)] == '\r' || bytes(_self)[uint(pos)] == '\n') {
                    pos++;
                } else if (bytes(_self)[uint(pos)] == ':') {
                    pos++;
                    colon = true;
                    break;
                } else {
                    break;
                }
            }

            if(uint(pos) == bytes(_self).length) {
                return;
            }

            if (colon) {
                break;
            } else {
                searchStart = uint(pos);
            }
        }

        int start = _self.indexOf("{", uint(pos));
        if (start == -1) {
            return;
        }
        //start += 1;
        
        int end = _self.indexOf("}", uint(pos));
        if (end == -1) {
            return;
        }
        end +=1 ;
        _ret = _self.substr(uint(start), uint(end-start));
    }

    function getIntArrayValueByKey(string _self, string _key, uint[] storage _array) internal {
         for (uint i=0; i<10; ++i) {
            //delete _array[i];
            _array[i] = i;
        }
        //_array.length = 0;

        /*
        int pos = -1;
        uint searchStart = 0;
        while (true) {
            pos = _self.indexOf("\"".concat(_key, "\""), searchStart);
            if (pos == -1) {
                pos = _self.indexOf("'".concat(_key, "'"), searchStart);
                if (pos == -1) {
                    return;
                }
            }

            pos += int(bytes(_key).length+2);

            bool colon = false;
            while (uint(pos) < bytes(_self).length) {
                if (bytes(_self)[uint(pos)] == ' ' || bytes(_self)[uint(pos)] == '\t' 
                    || bytes(_self)[uint(pos)] == '\r' || bytes(_self)[uint(pos)] == '\n') {
                    pos++;
                } else if (bytes(_self)[uint(pos)] == ':') {
                    pos++;
                    colon = true;
                    break;
                } else {
                    break;
                }
            }

            if(uint(pos) == bytes(_self).length) {
                return;
            }

            if (colon) {
                break;
            } else {
                searchStart = uint(pos);
            }
        }

        int start = _self.indexOf("[", uint(pos));
        if (start == -1) {
            return;
        }
        start += 1;
        
        int end = _self.indexOf("]", uint(pos));
        if (end == -1) {
            return;
        }

        string memory vals = _self.substr(uint(start), uint(end-start)).trim(" \t\r\n");

        if (bytes(vals).length == 0) {
            return;
        } */
        
        

        // string[] memory _strArray ;
        // vals.split(",", _strArray);

        // for (uint i=0; i<_strArray.length; ++i) {
        //     _array[i] = _strArray[i].trim(" \t\r\n");
        //     _array[i] = _strArray[i].trim("'\"");
        // }
    }

    function push(string _json) internal constant returns(uint _len) {
        if (bytes(_json).length == 0) {
            LibLog.log("ERROR: LibJson.push an empty json!!!");
            return 0;
        }
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][08JsonPush]";
        arg = arg.concat(_json);

        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        _len = uint(b32);
    }

    function pop() internal constant {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][07JsonPop]";
        if (bytes(arg).length == 0) {
            LibLog.log("ERROR: LibJson.pop an args invalid!!!");
        }
        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }
    }

    function size() internal constant returns(uint _ret) {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][13JsonStackSize]";
        if (bytes(arg).length == 0) {
            LibLog.log("ERROR: LibJson.size an args invalid!!!");
        }
        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }

        return uint(b32);
    }

    function clear() internal constant {
        string memory arg = "[69d98d6a04c41b4605aacb7bd2f74bee][14JsonStackClear]";
        if (bytes(arg).length == 0) {
            LibLog.log("ERROR: LibJson.clear an args invalid!!!");
        }
        uint argptr;
        uint arglen = bytes(arg).length;

        bytes32 b32;

        assembly {
            argptr := add(arg, 0x20)
            b32 := sha3(argptr, arglen)
        }
    }
} 
