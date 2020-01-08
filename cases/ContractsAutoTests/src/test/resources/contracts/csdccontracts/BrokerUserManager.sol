pragma solidity ^0.4.12;
/**
*@file      BrokerUserManager.sol
*@author    xuhui
*@time      2016-4-14
*@desc      the defination of BrokerUserManager
*/

import "./csdc_library/LibBrokerUser.sol";
import "./sysbase/OwnerNamed.sol";
import "./UserManager.sol";

contract BrokerUserManager is OwnerNamed {
    using LibBrokerUser for *;
    using LibString for *;
    using LibInt for *;

    event Notify(uint _errorno, string _info);
    
    mapping(address=>LibBrokerUser.BrokerUser)          brokerUserMap;
    address[]                                       	addrList;
    address[]                                       	tempList;

    LibBrokerUser.BrokerUser internal                    broker_User;

    /** @brief errno for test case */
    enum BrokerUserError {
        NO_ERROR,
        BAD_PARAMETER,
        ID_EMPTY,
        LOGINNAME_EMPTY,
        PASSWORD_EMPTY,
        NAME_EMPTY,
        AGE_EMPTY,
        SEX_EMPTY,
        BIRTHDAY_EMPTY,
        EMAIL_EMPTY,
        MOBILE_EMPTY,
        LOGINNAME_ALREADY_EXISTS,
        INSERT_FAILED,
        DEPT_NOT_EXISTS,
        ROLE_ID_INVALID,
        ROLE_ID_EXCEED_DEPT,
        USER_NOT_EXISTS,
        ROLE_ID_ALREADY_EXISTS,
        ADDRESS_ALREADY_EXISTS,
        ACCOUNT_ALREDY_EXISTS,
        ACCOUNT_CANNOT_UPDATE,
        USER_LOGIN_FAILED,
        DEPT_CANNOT_UPDATE,
        NO_PERMISSION,
        USER_STATUS_ERROR
    }

    uint errno_prefix = 16500;

    function BrokerUserManager() {
        register("CsdcModule", "0.0.1.0", "BrokerUserManager", "0.0.1.0");
    }

    /**
     * @dev 新增券商用户
     * @param _strJson _strJson表示的用户资料
     * @return _ret 返回值 true false
     */
    
    function insertBrokerUser(string _strJson) public returns (bool _ret){
    	//判断json格式
        if(!broker_User.fromJson(_strJson)){
            Notify(errno_prefix+uint(BrokerUserError.BAD_PARAMETER),"json invalid");
            return;
        }
        if(broker_User.id == address(0)) {
            Notify(errno_prefix+uint(BrokerUserError.ID_EMPTY), "id cannot be empty.");
            return;
        }
        if(brokerUserMap[broker_User.id].id != address(0)) {
            Notify(errno_prefix+uint(BrokerUserError.ADDRESS_ALREADY_EXISTS), "address already exists.");
            return;
        }
        if(broker_User.loginName.equals("")) {
            Notify(errno_prefix+uint(BrokerUserError.LOGINNAME_EMPTY), "login name cannot be empty.");
            return;
        }
        if(broker_User.password.equals("")) {
            Notify(errno_prefix+uint(BrokerUserError.PASSWORD_EMPTY), "password cannot be empty.");
            return;
        }
        if(broker_User.name.equals("")) {
            Notify(errno_prefix+uint(BrokerUserError.NAME_EMPTY), "name cannot be empty.");
            return;
        }
        //loginName 不能重复
        string memory _loginName = broker_User.loginName;
        for(uint i = 0; i < addrList.length; i++){
            if(brokerUserMap[addrList[i]].loginName.equals(_loginName)) {
                Notify(errno_prefix+uint(BrokerUserError.LOGINNAME_ALREADY_EXISTS), "login name already exists.");
                return;
            }
        }

        address _id = broker_User.id;
        UserManager userMgr = UserManager(rm.getContractAddress("SystemModuleManager", "0.0.1.0", "UserManager","0.0.1.0"));
        //拼接传入底层的json
        string memory _userJson = "{";
        _userJson = _userJson.concat(uint(_id).toAddrString().toKeyValue("userAddr"),",");
        _userJson = _userJson.concat(_loginName.toKeyValue("account"),",");
        _userJson = _userJson.concat(_loginName.toKeyValue("name"),",");
        _userJson = _userJson.concat("\"accountStatus\": 1", ",");
        _userJson = _userJson.concat("\"departmentId\": \"brokeruserdpt\"", ",");
        if(broker_User.role == 1){
            _userJson = _userJson.concat("\"roleIdList\": [\"csdc_role_201\", \"csdc_role_202\"]");
        }
        if(broker_User.role == 2){
            _userJson = _userJson.concat("\"roleIdList\": [\"csdc_role_203\", \"csdc_role_204\"]");
        }
        _userJson = _userJson.concat("}");
        if(userMgr.insert(_userJson)==0){
            brokerUserMap[_id] = broker_User;
            addrList.push(_id);
            broker_User.reset();
        }else{
            Notify(errno_prefix+uint(BrokerUserError.INSERT_FAILED), "insert a user from usermanager failed.");
            return;
        }
        Notify(uint(BrokerUserError.NO_ERROR), "insert a user success");
        return true;
    }
    /**
     * @dev 修改券商用户信息
     * @param _strJson 修改信息字段
     * @return _ret 返回信息
     */
    function updateUser(string _strJson) public returns(bool _ret){
        //address _id = msg.sender;
        if(!broker_User.fromJson(_strJson)) {
            Notify(uint(BrokerUserError.BAD_PARAMETER), "update the user success");
        }
        if(brokerUserMap[broker_User.id].id == address(0)) {
            Notify(errno_prefix+uint(BrokerUserError.USER_NOT_EXISTS), 'The user dose not exist.');
            return;
        }
        brokerUserMap[broker_User.id].update(_strJson);
        Notify(uint(BrokerUserError.NO_ERROR), "insert a user success");
        return true;        
    }
    /**
     * @dev 根据用户id地址获取用户信息
     * @param _id 用户地址
     * @return _ret 返回机构信息
     */
    function findById(address _id) constant public returns(string _ret){
        _ret = "{\"ret\":0,\"data\":{\"total\":0,\"items\":[]}}";

        if(brokerUserMap[_id].id == address(0)) {
            return;
        }
        _ret = "{\"ret\": 0, \"data\": {\"total\": 1, \"items\":[";
        _ret = _ret.concat(brokerUserMap[_id].toJson(),"]}}");
    }
    /**
     * @dev 根据登录名获取中证登用户信息
     * @param _loginName 登录名
     * @return _ret 返回机构信息
     */
    function findByBrokerLoginName(string _loginName) constant public returns(string _ret){
        bool _flag = false;
        _ret = "{\"ret\": 0, \"data\": {\"total\": 1, \"items\":[";
        for(uint i = 0; i < addrList.length; i++){
            if(brokerUserMap[addrList[i]].loginName.equals(_loginName)){
                _ret = _ret.concat(brokerUserMap[addrList[i]].toJson(),"]}}");
                _flag = true;
                return;
            }
        }
        if(!_flag){
            _ret = "{\"ret\": 0, \"data\": {\"total\": 0, \"items\":[]}}";
        }
    }

    /**
    * @dev 根据状态分页显示用户
    * @ _status 账户状态
    * @return _ret
    */
    function pageByStatus(string _strJson) constant public returns (string _ret){
        _ret = "{\"ret\":0,\"data\":{\"total\":0,\"items\":[]}}";
        if (addrList.length <= 0) {
            return;
        }

        uint _pageSize = uint(_strJson.getIntValueByKey("pageSize"));
        uint _pageNo = uint(_strJson.getIntValueByKey("pageNo"));
        uint _status = uint(_strJson.getIntValueByKey("status"));
        uint _role = uint(_strJson.getIntValueByKey("role"));
        uint _brokerId = uint(_strJson.getIntValueByKey("brokerId"));

        uint _startIndex = _pageSize * _pageNo;

        if (_startIndex >= addrList.length) {
          return;
        }

        _ret = "";

        string memory _data;
        uint _count = 0; //满足条件的消息数量
        uint _total = 0; //满足条件的指定页数的消息数量
        for (uint i = 0; i < addrList.length; i++) {
            broker_User = brokerUserMap[addrList[i]];
            if (_status != 0 && uint(broker_User.status) != _status) {
              continue;
            }

            if (_role != 0 && uint(broker_User.role) != _role) {
              continue;
            }

            if (_brokerId != 0 && uint(broker_User.brokerId) != _brokerId) {
              continue;
            }

            if (_count < _startIndex) {
              _count ++;
              continue;
            }

            if (_total > 0) {
              _data = _data.concat(",");
            }
            _count ++;
            _total ++;
            _data = _data.concat(broker_User.toJson());

            if (_total == _pageSize) {
              break;
            }
        }
        _ret = _ret.concat("{\"ret\":0,\"data\":{\"total\":", _count.toString(), ",\"items\":[");
        _ret = _ret.concat(_data, "]}}");
    }

    /**
     * @dev 注销用户
     * @param _id 用户地址
     * @return _ret
     */
    function delBrokerUser(address _id) public returns (bool _ret){
        if(brokerUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(BrokerUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        if(brokerUserMap[_id].status != LibBrokerUser.AccountStatus.VALID) {
            Notify(errno_prefix+uint(BrokerUserError.USER_STATUS_ERROR), "The user status is not valid.");
            return;
        }
        brokerUserMap[_id].status = LibBrokerUser.AccountStatus.INVALID;
        Notify(uint(BrokerUserError.NO_ERROR), "delete the user success");
        return true;
    }

    /**
     * @dev 锁定账户
     * @param _id 机构id
     * @return _ret 返回值 true false
     */
    function lockBrokerUser(address _id) public returns (bool _ret){
        if(brokerUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(BrokerUserError.USER_STATUS_ERROR), "The user is not existed");
            return;
        }

        if(brokerUserMap[_id].status != LibBrokerUser.AccountStatus.VALID) {
            Notify(errno_prefix+uint(BrokerUserError.USER_STATUS_ERROR), "The user status is not valid.");
            return;
        } 
        brokerUserMap[_id].status = LibBrokerUser.AccountStatus.LOCKED;
        Notify(uint(BrokerUserError.NO_ERROR), "lock the user success");
        return true;
    }

    /**
     * @dev 解锁中证登用户
     * @param _id 机构id
     * @return _ret 返回值 true false
     */
    function unlockBrokerUser(address _id) public returns (bool _ret){
        if(brokerUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(BrokerUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        if(brokerUserMap[_id].status != LibBrokerUser.AccountStatus.LOCKED) {
            Notify(errno_prefix+uint(BrokerUserError.USER_STATUS_ERROR), "The user status is not locked.");
            return;
        } 
        brokerUserMap[_id].status = LibBrokerUser.AccountStatus.VALID;
        Notify(uint(BrokerUserError.NO_ERROR), "unlock the user success");
        return true;
    }

    /**
    * @dev 根据用户地址获取用户密码
    * @param _id 用户id
    * @return _ret 返回值
    */
    function findPasswordByUserId(address _id) constant public returns(string _ret){

        if(brokerUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(BrokerUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        _ret = brokerUserMap[_id].password;
        Notify(uint(BrokerUserError.NO_ERROR), "find the user password success");
        
    }
    /**
     * @dev 根据登录名获取机构用户密码
     * @param _loginName 登录名
     * @return _ret 返回值
     */
    function findPasswordByLoginName(string _loginName) constant public returns(string _ret){
        
        for(uint i = 0; i < addrList.length; i++){
            if(brokerUserMap[addrList[i]].loginName.equals(_loginName)){
                _ret = brokerUserMap[addrList[i]].password;
                Notify(errno_prefix+uint(BrokerUserError.NO_ERROR), "find the user password success");
                return;
            }
        }
        Notify(errno_prefix+uint(BrokerUserError.USER_NOT_EXISTS), "The user is not existed");
        
    }
}