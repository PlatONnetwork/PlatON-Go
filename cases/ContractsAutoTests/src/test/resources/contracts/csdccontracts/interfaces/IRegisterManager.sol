pragma solidity ^0.4.12;

contract IRegisterManager {

    ////Note: 注册模块
    function register(string _moduleName, string _moduleVersion) returns(uint) ;

    ////Note: 注册合约
    function register(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) returns(uint) ;

    ////Note: 注销
    function unRegister() ;

    ////Note: 获取模块地址
    function getModuleAddress(string _moduleName, string _moduleVersion) constant returns (address _address) ;

    ////Note: 获取合约地址
    function getContractAddress(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) constant returns (address _address) ;

    ////Note: 检查模块地址是否已经注册
    function IfModuleRegist(address _moduleAddr) constant returns(bool) ;

    ////Note: 检查模块是否已经注册
    function IfModuleRegist(string _moduleName, string _moduleVersion) constant returns(bool) ;

    ////Note: 检查合约地址是否已经注册
    function IfContractRegist(address _contractAddr) constant returns(bool) ;

    ////Note: 检查合约是否已经注册
    function IfContractRegist(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) constant returns(bool) ;

    ////Note: 查询合约名字 (合约名字仅限于32个字符)，入参为合约地址
    function findResNameByAddress(address _addr) constant public returns(uint _contractName) ;

    ////Note: 查询合约版本号 (合约名字仅限于32个字符)，入参为合约地址
    function findContractVersionByAddress(address _addr) constant public returns(uint _contractVersion);

    ////Note: 查询模块名字 (模块名字仅限于32个字符)，入参为模块地址或模块内合约地址
    function findModuleNameByAddress(address _addr) constant public returns(uint _moduleName) ;

    ////Note: 查询模块版本号 (模块名字仅限于32个字符)，入参为模块地址或模块内合约地址
    function findModuleVersionByAddress(address _addr) constant public returns(uint _moduleVersion);

    ////Note: 获取所有已注册的合约列表
    function getRegisteredContract(uint _pageNum, uint _pageSize) constant public returns(string _json) ;

    ////Note: 合约数据迁移
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
        string _toModuleNameAndVersion, string _toNameAndVersion, string _signString) public returns (uint _errno);

    ////Note: 修改模块拥有者
    function changeModuleRegisterOwner(string _moduleName, string _moduleVersion, address _newOwner) returns(uint);
    ////Note: 修改合约拥有者
    function changeContractRegisterOwner(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, address _newOwner) returns(uint);
}
