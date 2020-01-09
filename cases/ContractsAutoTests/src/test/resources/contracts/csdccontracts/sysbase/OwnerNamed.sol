pragma solidity ^0.4.12;

import "../interfaces/IRegisterManager.sol";
import "../utillib/LibString.sol";
import "../utillib/LibInt.sol";
import "../utillib/LibLog.sol";
import "../utillib/LibJson.sol";
import "../utillib/LibStack.sol";

contract OwnerNamed {

    using LibString for *;
    using LibInt for *;

    address owner;
    IRegisterManager rm;

    function OwnerNamed() {
        owner = msg.sender;
        rm = IRegisterManager(0x0000000000000000000000000000000000000011);
    }

    function register(string _moduleName, string _moduleVersion) public {
        uint ret = rm.register(_moduleName,_moduleVersion);
        if (ret != 0) {
            log("register module: failed. abort!!!", _moduleName, _moduleVersion);
            return;
        }
    }

    function register(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) public {
        /*if (!_moduleName.equals("defaultModule")) {
            if (!IfModuleRegist(_moduleName, _moduleVersion)) {
                log("module: is unregistered.", _moduleName, _moduleVersion);
                log("contract: register failed.", _contractName, _contractVersion);
                return;
            }
        }*/
        uint ret = rm.register(_moduleName,_moduleVersion,_contractName,_contractVersion);
        if (ret != 0) {
            log("register contract: failed. abort!!!", _moduleName, _contractName);
            return;
        }

        updateContractAddr(_moduleName,_moduleVersion,_contractName,_contractVersion,this);
    }

    function changeModuleRegisterOwner(string _moduleName, string _moduleVersion, address _newOwner) public {
        uint ret = rm.changeModuleRegisterOwner(_moduleName,_moduleVersion,_newOwner);
        if (ret != 0) {
            log("register changeModuleRegisterOwner: failed. abort!!!", _moduleName, _moduleVersion);
            return;
        }
    }

    function changeContractRegisterOwner(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, address _newOwner) public {
        uint ret = rm.changeContractRegisterOwner(_moduleName,_moduleVersion,_contractName,_contractVersion,_newOwner);
        if (ret != 0) {
            log("register changeContractRegisterOwner: failed. abort!!!", _moduleName, _contractName);
            return;
        }
    }

    function kill() public {
        if (msg.sender != owner) {
            return;
        }

        rm.unRegister();
    }

    function getOwner() constant public returns (string _ret) {
        _ret = owner.addrToAsciiString();
    }

    function getSender() constant public returns (string _ret) {
        _ret = msg.sender.addrToAsciiString();
    }

    uint errno = 0;

    function getErrno() constant returns (uint) {
        return errno;
    }

    //old log functions
    function log(string _str) constant public returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat("|", _str);

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
   }
   function log(string _str, string _str2) constant public returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat("|", _str);
        cmd = cmd.concat("|", _str2);

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
   }
   function log(string _str, string _str2, string _str3) constant public returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat("|", _str);
        cmd = cmd.concat("|", _str2);
        cmd = cmd.concat("|", _str3);

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
   }
   function log(string _str, uint _ui) constant public returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat("|", _str);
        cmd = cmd.concat("|", _ui.toString());

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
   }
   function log(string _str, int _i) constant public returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat("|", _str);
        cmd = cmd.concat("|", _i.toString());

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }

        _ret = 0;
   }
   function log(string _str, address _addr) constant public returns(uint _ret) {
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][05vmlog]";
        cmd = cmd.concat("|", _str);
        cmd = cmd.concat("|", uint(_addr).toAddrString());

        uint strptr;
        uint strlen = bytes(cmd).length;

        bytes32 hash;

        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        _ret = 0;
   }

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

    function updateContractAddr(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, address __contractAddr) constant private returns (uint _ret) {
        string memory value = _moduleName.concat("|", _moduleVersion);
        value = value.concat("|", _contractName);
        value = value.concat("|", _contractVersion);
        _ret = writedb("contractAddr|update", uint(__contractAddr).toAddrString(), value);
        if (0 != _ret)
            log("OwnerNamed.sol", "update contractAddr failed.");
        else
            log("OwnerNamed.sol", "update contractAddr success.");
    }

}
