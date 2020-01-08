pragma solidity ^0.4.12;
/**
* @file      RegisterManager.sol
* @author    yanze
* @time      2016-11-29
* @desc      the defination of OwnerNamed contract
*/

import "./utillib/LibString.sol";
import "./utillib/LibInt.sol";
import "./utillib/LibDecode.sol";
import "./utillib/LibLog.sol";
import "./interfaces/IRegisterManager.sol";

contract RegisterManager is IRegisterManager {

    using LibString for *;
    using LibInt for *;

    mapping(address => Register) registerMap;
    mapping(string => address) keyMap;
    address[] contractAddrList;
    address[] tmpContractAddrList;
    address   owner;

    /** 注册信息  */
    struct Register {
        string moduleName;
        string moduleVersion;
        string contractName;
        string contractVersion;
        address addr;
        address origin;
        bool isModule;//true:module, false:contract
        uint createTime;
    }

    function RegisterManager(){
        owner = msg.sender;
    }

    enum RegisterManagerError{
        NO_ERROR,
        NO_SAME_ORIGIN,
        SIGN_FAILED,
        TRANSFER_CONTRACT_FAILED
    }

    /** 注册模块 */
    function register(string _moduleName, string _moduleVersion) returns(uint) {
        LibLog.log("register->msg.sender(new module address):",msg.sender);
        LibLog.log("register->owner(old owner):",owner);
        LibLog.log("register->tx.origin:",tx.origin);
        address sendAddr = msg.sender;

        string memory key = _moduleName.concat(_moduleVersion);

        if(keyMap[key] != address(0) && registerMap[keyMap[key]].origin != tx.origin) {
            LibLog.log("old module address, keyMap[key]:",uint(keyMap[key]).toAddrString());
            if (keyMap[key] != address(0))
                LibLog.log("tx.origin != registerMap[keyMap[key]].origin(old module origin)",uint(tx.origin).toAddrString(),uint(registerMap[keyMap[key]].origin).toAddrString());
            LibLog.log("register fail, tx.origin not publisher.",key);
            return uint(-1);
        }

        Register register = registerMap[sendAddr];
        register.moduleName = _moduleName;
        register.moduleVersion = _moduleVersion;
        register.addr = sendAddr;
        register.origin= tx.origin;
        register.isModule = true;
        register.createTime=now*1000;
        keyMap[key] = sendAddr;
        contractAddrList.push(sendAddr);
        LibLog.log("IRegisterManager -> register success.",key);
        return 0;
    }

    /** 注册合约 */
    function register(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) returns(uint) {
        LibLog.log("register->msg.sender(new contract address):",msg.sender);
        LibLog.log("register->owner(old owner):",owner);
        LibLog.log("register->tx.origin:",tx.origin);
        address sendAddr = msg.sender;

        if (_moduleName.equals(_contractName)) {
            LibLog.log("register fail, moduleName equals contractName.");
            return uint(-1);
        }

        string memory key = _moduleName.concat(_moduleVersion, _contractName, _contractVersion);

        if(keyMap[key] != address(0) && registerMap[keyMap[key]].origin != tx.origin) {
            LibLog.log("old contract address, keyMap[key]:",uint(keyMap[key]).toAddrString());
            if (keyMap[key] != address(0))
                LibLog.log("tx.origin != registerMap[keyMap[key]].origin(old contract origin)",uint(tx.origin).toAddrString(),uint(registerMap[keyMap[key]].origin).toAddrString());
            LibLog.log("register fail, tx.origin not publisher.",key);
            return uint(-1);
        }

        Register register = registerMap[sendAddr];
        register.moduleName = _moduleName;
        register.moduleVersion = _moduleVersion;
        register.contractName = _contractName;
        register.contractVersion = _contractVersion;
        register.addr = sendAddr;
        register.origin= tx.origin;
        register.isModule = false;
        register.createTime=now*1000;
        keyMap[key] = sendAddr;
        contractAddrList.push(sendAddr);
        LibLog.log("IRegisterManager -> register success.",key);
        return 0;
    }

    function changeModuleRegisterOwner(string _moduleName, string _moduleVersion, address _newOwner) returns(uint) {
        LibLog.log("in changeModuleRegisterOwner", _moduleName, _moduleVersion);
        LibLog.log("in changeModuleRegisterOwner newOwner:", _newOwner);
        if(bytes(_moduleName).length == 0){
          LibLog.log("moduleName is null..");
          return uint(-1);
        }
        if(bytes(_moduleVersion).length == 0){
          LibLog.log("moduleVersion is null..");
          return uint(-1);
        }
        // no address check
        // check module
        string memory key = _moduleName.concat(_moduleVersion);
        if (keyMap[key] == address(0)) {
            LibLog.log("module: is unregistered", _moduleName, _moduleVersion);
            return uint(-1);
        }
        Register register = registerMap[keyMap[key]];
        if (!register.isModule) {
            LibLog.log("module: is not module", _moduleName, _moduleVersion);
            return uint(-1);
        }

        LibLog.log("change module owner from->to:", uint(register.origin).toAddrString(), uint(_newOwner).toAddrString());
        register.origin = _newOwner;
        return 0;
    }

    function changeContractRegisterOwner(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, address _newOwner) returns(uint) {
        LibLog.log("in changeContractRegisterOwner", _moduleName, _moduleVersion);
        LibLog.log("in changeContractRegisterOwner", _contractName, _contractVersion);
        LibLog.log("in changeContractRegisterOwner newOwner:", _newOwner);
        if(bytes(_moduleName).length == 0){
          LibLog.log("moduleName is null..");
          return uint(-1);
        }
        if(bytes(_moduleVersion).length == 0){
          LibLog.log("moduleVersion is null..");
          return uint(-1);
        }
        if(bytes(_contractName).length == 0){
          LibLog.log("contractName is null..");
          return uint(-1);
        }
        if(bytes(_contractVersion).length == 0){
          LibLog.log("contractVersion is null..");
          return uint(-1);
        }
        // no address check
        // check module
        string memory key = _moduleName.concat(_moduleVersion);
        if (keyMap[key] == address(0)) {
            LibLog.log("module: is unregistered", _moduleName, _moduleVersion);
            return uint(-1);
        }
        Register register = registerMap[keyMap[key]];
        if (!register.isModule) {
            LibLog.log("module: is not module", _moduleName, _moduleVersion);
            return uint(-1);
        }
        // check contract
        key = key.concat(_contractName, _contractVersion);
        if (keyMap[key] == address(0)) {
            LibLog.log("contract: is unregistered", _contractName, _contractVersion);
            return uint(-1);
        }
        register = registerMap[keyMap[key]];
        if (register.isModule) {
            LibLog.log("contract: is not contract", _contractName, _contractVersion);
            return uint(-1);
        }

        LibLog.log("change contract owner from->to:", uint(register.origin).toAddrString(), uint(_newOwner).toAddrString());
        register.origin = _newOwner;
        return 0;
    }

    /** 注销 */
    function unRegister() {
        address sendAddr = msg.sender;
        Register register = registerMap[sendAddr];
        string memory key;
        if (register.isModule)
            key = register.moduleName.concat(register.moduleVersion);
        else
            key = register.moduleName.concat(register.moduleVersion, register.contractName, register.contractVersion);
        delete registerMap[sendAddr];
        delete keyMap[key];

        for (uint i=0; i<contractAddrList.length; ++i) {
            if (contractAddrList[i] != sendAddr) {
                tmpContractAddrList.push(contractAddrList[i]);
            }
        }
        delete contractAddrList;
        for (i=0; i<tmpContractAddrList.length; ++i) {
            contractAddrList.push(tmpContractAddrList[i]);
        }
        delete tmpContractAddrList;
    }

    /** 获取模块地址 */
    function getModuleAddress(string _moduleName, string _moduleVersion) constant returns (address _address) {
        string memory key = _moduleName.concat(_moduleVersion);
        _address = keyMap[key];
    }

    /** 获取合约地址 */
    function getContractAddress(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) constant returns (address _address) {
        string memory key = _moduleName.concat(_moduleVersion, _contractName, _contractVersion);
        _address = keyMap[key];
    }

    /** 检查模块地址是否已经注册  */
    function IfModuleRegist(address _moduleAddr) constant returns(bool) {
        Register register = registerMap[_moduleAddr];
        if (register.addr == 0x0) {
            return false;
        }
        if (!register.isModule) {
            return false;
        }
        string memory key = register.moduleName.concat(register.moduleVersion);
        if (keyMap[key] != _moduleAddr) {
            return false;
        }
        return true;
    }

    /** 检查模块是否已经注册  */
    function IfModuleRegist(string _moduleName, string _moduleVersion) constant returns(bool) {
        string memory key = _moduleName.concat(_moduleVersion);
        if (keyMap[key] == address(0)) {
            LibLog.log("module: is unregistered", _moduleName, _moduleVersion);
            return false;
        }
        Register register = registerMap[keyMap[key]];
        if (!register.isModule) {
            LibLog.log("module: is not module", _moduleName, _moduleVersion);
            return false;
        }
        return true;
    }

    /** 检查合约地址是否已经注册  */
    function IfContractRegist(address _contractAddr) constant returns(bool) {
        Register register = registerMap[_contractAddr];
        if (register.addr == 0x0) {
            return false;
        }
        if (register.isModule) {
            return false;
        }
        string memory key = register.moduleName.concat(register.moduleVersion, register.contractName, register.contractVersion);
        if (keyMap[key] != _contractAddr) {
            return false;
        }
        return true;
    }

    /** 检查合约是否已经注册  */
    function IfContractRegist(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) constant returns(bool) {
        string memory key = _moduleName.concat(_moduleVersion);
        if (keyMap[key] == address(0)) {
            LibLog.log("module: is unregistered", _moduleName, _moduleVersion);
            return false;
        }
        Register register = registerMap[keyMap[key]];
        if (!register.isModule) {
            LibLog.log("module: is not module", _moduleName, _moduleVersion);
            return false;
        }

        key = key.concat(_contractName, _contractVersion);
        if (keyMap[key] == address(0)) {
            LibLog.log("contract: is unregistered", _contractName, _contractVersion);
            return false;
        }
        register = registerMap[keyMap[key]];
        if (register.isModule) {
            LibLog.log("contract: is not contract", _contractName, _contractVersion);
            return false;
        }
        return true;
    }

    /** 查询合约名字 (合约名字仅限于32个字符)，入参为合约地址 */
    function findResNameByAddress(address _addr) constant public returns(uint _contractName) {
        _contractName = 0;
        Register register = registerMap[_addr];
        if (register.addr == 0x0) {
            return;
        }
        if (register.isModule) {
            LibLog.log("address belongs to module:, not contract", register.moduleName);
            return;
        }
        string memory key = register.moduleName.concat(register.moduleVersion, register.contractName, register.contractVersion);
        if (keyMap[key] == _addr) {
            _contractName = register.contractName.storageToUint();
            return;
        }
    }

    /** 查询合约版本号 (合约名字仅限于32个字符)，入参为合约地址 */
    function findContractVersionByAddress(address _addr) constant public returns(uint _contractVersion) {
        _contractVersion = 0;
        Register register = registerMap[_addr];
        if (register.addr == 0x0) {
            return;
        }
        if (register.isModule) {
            LibLog.log("address belongs to module:, not contract", register.moduleName);
            return;
        }
        string memory key = register.moduleName.concat(register.moduleVersion, register.contractName, register.contractVersion);
        if (keyMap[key] == _addr) {
            _contractVersion = register.contractVersion.storageToUint();
            return;
        }
    }

    /** 查询模块名字 (模块名字仅限于32个字符)，入参为模块地址或模块内合约地址 */
    function findModuleNameByAddress(address _addr) constant public returns(uint _moduleName) {
        _moduleName = 0;
        Register register = registerMap[_addr];
        if (register.addr == 0x0) {
            return;
        }
        string memory key;
        if (register.isModule)
            key = register.moduleName.concat(register.moduleVersion);
        else
            key = register.moduleName.concat(register.moduleVersion, register.contractName, register.contractVersion);
        if (keyMap[key] == _addr) {
            _moduleName = register.moduleName.storageToUint();
            return;
        }
    }

    /** 查询模块版本号 (模块版本号仅限于32个字符)，入参为模块地址或模块内合约地址 */
    function findModuleVersionByAddress(address _addr) constant public returns(uint _moduleVersion) {
        _moduleVersion = 0;
        Register register = registerMap[_addr];
        if (register.addr == 0x0) {
            return;
        }
        string memory key;
        if (register.isModule)
            key = register.moduleName.concat(register.moduleVersion);
        else
            key = register.moduleName.concat(register.moduleVersion, register.contractName, register.contractVersion);
        if (keyMap[key] == _addr) {
            _moduleVersion = register.moduleVersion.storageToUint();
            return;
        }
    }

    function getRegisteredContract(uint _pageNum, uint _pageSize) constant public returns(string _json) {
        _json = _json.concat("{");
        _json = _json.concat(uint(0).toKeyValue("ret"), ",");
        _json = _json.concat("\"data\":{");
        _json = _json.concat(uint(contractAddrList.length).toKeyValue("total"), ",");
        _json = _json.concat("\"items\":[");

        uint start = _pageNum * _pageSize;
        uint end = (_pageNum+1) * _pageSize;
        for (uint i=start; i<end && i<contractAddrList.length; ++i) {
            if (i > start) {
                _json = _json.concat(",");
            }
            _json = _json.concat("{");
            _json = _json.concat(registerMap[contractAddrList[i]].moduleName.toKeyValue("moduleName"), ",");
            _json = _json.concat(registerMap[contractAddrList[i]].moduleVersion.toKeyValue("moduleVersion"), ",");
            _json = _json.concat(uint(registerMap[contractAddrList[i]].createTime).toKeyValue("createTime"), ",");
            if (!registerMap[contractAddrList[i]].isModule) {
                _json = _json.concat(registerMap[contractAddrList[i]].contractName.toKeyValue("contractName"), ",");
                _json = _json.concat(registerMap[contractAddrList[i]].contractVersion.toKeyValue("contractVersion"), ",");
            }
            _json = _json.concat(uint(registerMap[contractAddrList[i]].addr).toAddrString().toKeyValue("address"), ",");
            if (registerMap[contractAddrList[i]].isModule)
                _json = _json.concat("\"isModule\":true");
            else
                _json = _json.concat("\"isModule\":false");
            _json = _json.concat("}");
        }

        _json = _json.concat("]}}");
    }

    /**
    * transfer 2 contracts
    * @param _fromModuleNameAndVersion = fromModuleName.concat(fromModuleVersion)
    * @param _fromNameAndVersion = fromContractName.concat(fromContractVersion)
    * @param _toModuleNameAndVersion = toModuleName.concat(toModuleVersion)
    * @param _toNameAndVersion = toContractName.concat(toContractVersion)
    * @param _signString sign
    * @return return errno
    */
     function transferContract(string _fromModuleNameAndVersion, string _fromNameAndVersion,
        string _toModuleNameAndVersion, string _toNameAndVersion, string _signString) public returns (uint _errno) {
        string memory fromKey = _fromModuleNameAndVersion.concat(_fromNameAndVersion);
        string memory toKey = _toModuleNameAndVersion.concat(_toNameAndVersion);
        string memory fromStr = registerMap[keyMap[fromKey]].addr.addrToAsciiString();
        string memory toStr = registerMap[keyMap[toKey]].addr.addrToAsciiString();
        _errno = uint256(RegisterManagerError.NO_ERROR);
		    //uint errno_prefix = 95270;
        if (!isSameOrigin(fromKey, toKey)) {
			      _errno = 95270 + uint256(RegisterManagerError.NO_SAME_ORIGIN);
            return _errno;
        }
        //string memory sha3Data = fromStr.concat(toStr);
        address addr = LibDecode.decode(_signString, sha3(fromStr.concat(toStr)));

        if (registerMap[keyMap[fromKey]].origin != addr){
			      _errno = 95270 + uint256(RegisterManagerError.SIGN_FAILED);
            return _errno;
        }

        _errno = uint(executeTransferContract(fromStr, toStr));
        if (_errno != 0){
            _errno = 95270 + uint256(RegisterManagerError.TRANSFER_CONTRACT_FAILED);
            return _errno;
        }
        return _errno;
    }
    function executeTransferContract(string _fromStr, string _toStr) internal returns(bytes32){
        string memory cmd = "[69d98d6a04c41b4605aacb7bd2f74bee][16transfercontract]";
        cmd = cmd.concat("|", _fromStr);
        cmd = cmd.concat("|", _toStr);
        uint strptr;
        uint strlen = bytes(cmd).length;
        bytes32 hash;
        assembly {
            strptr := add(cmd, 0x20)
            hash := sha3(strptr, strlen)
        }
        return hash;
    }
    function isSameOrigin(string _fromKey, string _toKey) constant internal returns (bool){
        address fromOrigin = registerMap[keyMap[_fromKey]].origin;
        address toOrigin = registerMap[keyMap[_toKey]].origin;
        return fromOrigin == toOrigin;
    }
}
