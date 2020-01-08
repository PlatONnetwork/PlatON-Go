pragma solidity ^0.4.12;
/**
*@file      PerUserManager.sol
*@author    xuhui
*@time      2016-3-6
*@desc      the defination of PerUserManager
*/

import "./utillib/LibStack.sol";
import "./sysbase/OwnerNamed.sol";
import "./csdc_library/LibPerUser.sol";
import "./csdc_base/CsdcBaseInterface.sol";

import "./UserManager.sol";
import "./OrgUserManager.sol";


/* @brief PerUserManager definition */
contract PerUserManager is OwnerNamed, CsdcBaseInterface{
    using LibPerUser for *;
    using LibString for *;
    using LibInt for *;
    
    mapping(address=>LibPerUser.PerUser)                perUserMap;
    address[]                                           addrList;
    address[]                                           tempList;

    //inner setting member
    LibPerUser.PerUser internal                         per_User;

    /** @brief errno for test case */
    enum PerUserError {
        NO_ERROR,
        BAD_PARAMETER,
        ID_EMPTY,   
        LOGINNAME_EMPTY,
        IDTYPE_EMPTY,
        IDNO_EMPTY,
        EMAIL_EMPTY,
        MOBILE_EMPTY,
        NAME_EMPTY,
        IDVALIDITY_EMPTY,
        ADDR_EMPTY,
        POSTCODE_EMPTY,
        USER_NOT_EXISTS,
        ADDRESS_ALREADY_EXISTS,
        LOGINNAME_ALREADY_EXISTS,
        IDENTITY_ALREADY_EXISTS,
        INSERT_FAILED,
        USER_STATUS_ERROR,
        USER_HAS_TODO,
        MOBILE_ALREADY_EXISTS
    }
    uint errno_prefix = 16200;

    modifier getOm(){ 
        om = OrgUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrgUserManager","0.0.1.0")); 
        _;
    }

    function PerUserManager() {
        register("CsdcModule", "0.0.1.0", "PerUserManager", "0.0.1.0");
    }

    CsdcBaseInterface bi;
    OrgUserManager om;

    /**
     * @dev 新增个人用户
     * @param _json _json表示的用户资料
        id, loginName, idType, idNo, email, mobile
     * @return _ret 返回值 true false
     */
    function insertPerUser(string _json) returns(bool _ret) {
        //判断json格式
        if(!per_User.fromJson(_json)){
            Notify(errno_prefix+uint(PerUserError.BAD_PARAMETER),"json invalid");
            return;
        }

        if(per_User.id == address(0)) {
            Notify(errno_prefix+uint(PerUserError.ID_EMPTY), "id cannot be empty.");
            return;
        }

        if(perUserMap[per_User.id].id != address(0)) {
            Notify(errno_prefix+uint(PerUserError.ADDRESS_ALREADY_EXISTS), "address already exists.");
            return;
        }

        if(per_User.loginName.equals("")) {
            Notify(errno_prefix+uint(PerUserError.LOGINNAME_EMPTY), "login name cannot be empty.");
            return;
        }

        if(uint(per_User.idType) == 0) {
            Notify(errno_prefix+uint(PerUserError.IDTYPE_EMPTY), "idType cannot be empty.");
            return;
        }

        if(per_User.idNo.equals("")) {
            Notify(errno_prefix+uint(PerUserError.IDNO_EMPTY), "idNo cannot be empty.");
            return;
        }

        if(per_User.email.equals("")) {
            Notify(errno_prefix+uint(PerUserError.EMAIL_EMPTY), "email cannot be empty.");
            return;
        }

        if(per_User.mobile.equals("")) {
            Notify(errno_prefix+uint(PerUserError.MOBILE_EMPTY), "mobile cannot be empty.");
            return;
        }

        //loginName 不能重复
        string memory _loginName = per_User.loginName;
        for(uint i = 0; i < addrList.length; i++){
            if(perUserMap[addrList[i]].loginName.equals(_loginName)) {
                Notify(errno_prefix+uint(PerUserError.LOGINNAME_ALREADY_EXISTS), "login name already exists.");
                return;
            }
        }

        //身份证件不能重复
        address _userId = findIdByIdNo(per_User.idNo, LibCommonEnum.IdType(per_User.idType));
        if(_userId != address(0) && per_User.id != _userId) {
            Notify(errno_prefix+uint(PerUserError.IDENTITY_ALREADY_EXISTS), "identity already exists");
            return;
        }

        //手机号不能重复
        _userId = findIdByMobile(per_User.mobile);
        if(_userId != address(0) && per_User.id != _userId) {
            //若有未认证用户已使用该手机号，则将该未认证用户置为无效
            if(perUserMap[_userId].status == LibPerUser.PerUserStatus.INITIAL) {
                perUserMap[_userId].status = LibPerUser.PerUserStatus.INVALID;
            } else {
                Notify(errno_prefix+uint(PerUserError.MOBILE_ALREADY_EXISTS), "mobile already exists");
                return; 
            }
        }

        address _id = per_User.id;

        UserManager userMgr = UserManager(rm.getContractAddress("SystemModuleManager", "0.0.1.0", "UserManager","0.0.1.0"));
        //拼接传入底层的json
        string memory _userJson = "{";
        _userJson = _userJson.concat(uint(_id).toAddrString().toKeyValue("userAddr"),",");
        _userJson = _userJson.concat(_loginName.toKeyValue("account"),",");
        _userJson = _userJson.concat(_loginName.toKeyValue("name"),",");
        // _userJson = _userJson.concat(per_User.mobile.toKeyValue("mobile"),",");
        // _userJson = _userJson.concat(per_User.email.getStringValueByKey("email").toKeyValue("email"), ",");
        _userJson = _userJson.concat("\"accountStatus\": 1", ",");
        _userJson = _userJson.concat("\"departmentId\": \"peruserdpt\"", ",");
        _userJson = _userJson.concat("\"roleIdList\": [\"csdc_role_301\",\"csdc_role_302\"]");
        // _userJson = _userJson.concat("\"roleIdList\": \"['peruserrole']\"");
        _userJson = _userJson.concat("}");
        LibLog.log("insertUser:", _userJson);
        uint _errno = userMgr.insert(_userJson);
        if(_errno == 0){
            perUserMap[_id].id = _id;
            perUserMap[_id].loginName = per_User.loginName;
            perUserMap[_id].idType = per_User.idType;
            perUserMap[_id].idNo = per_User.idNo;
            perUserMap[_id].email = per_User.email;
            perUserMap[_id].mobile = per_User.mobile;
            perUserMap[_id].status = LibPerUser.PerUserStatus.INITIAL;
            addrList.push(_id);
            per_User.reset();
        }else{
            Notify(errno_prefix+uint(PerUserError.INSERT_FAILED), "insert a user from usermanager failed.");
            return;
        }
        Notify(uint(PerUserError.NO_ERROR), "insert a user success");
        return true;
    }
    
    /**
     * @dev 认证个人用户
     * @param _json _json表示的认证参数
            id, idValidity, phone, addr, postCode
     * @return _ret 返回值 true false
     */
    function authPerUser(string _json) returns (bool _ret){
        //判断json格式
        if(!per_User.fromJson(_json)){
            Notify(errno_prefix+uint(PerUserError.BAD_PARAMETER),"json invalid");
            return;
        }

        if(perUserMap[per_User.id].status != LibPerUser.PerUserStatus.INITIAL) {
            Notify(errno_prefix+uint(PerUserError.USER_STATUS_ERROR), "The user status is not initial.");
            return;
        }   

        if(per_User.name.equals("")) {
            Notify(errno_prefix+uint(PerUserError.NAME_EMPTY), "name cannot be empty.");
            return;
        }

        if(perUserMap[per_User.id].id == address(0)) {
            Notify(errno_prefix+uint(PerUserError.USER_NOT_EXISTS), 'The user dose not exist.');
            return;
        }

        if(per_User.idValidity == 0) {
            Notify(errno_prefix+uint(PerUserError.IDVALIDITY_EMPTY), "idValidity cannot be empty.");
            return;
        }

        // if(per_User.phone.equals("")) {
        //     Notify(errno_prefix+uint(PerUserError.IDNO_EMPTY), "idNo cannot be empty.");
        //     return;
        // }

        if(per_User.addr.equals("")) {
            Notify(errno_prefix+uint(PerUserError.ADDR_EMPTY), "addr cannot be empty.");
            return;
        }

        if(per_User.postCode.equals("")) {
            Notify(errno_prefix+uint(PerUserError.POSTCODE_EMPTY), "postCode cannot be empty.");
            return;
        }

        perUserMap[per_User.id].name = per_User.name;
        perUserMap[per_User.id].idValidity = per_User.idValidity;
        perUserMap[per_User.id].phone = per_User.phone;
        perUserMap[per_User.id].addr = per_User.addr;
        perUserMap[per_User.id].postCode = per_User.postCode;

        perUserMap[per_User.id].status = LibPerUser.PerUserStatus.VALID;
        Notify(uint(PerUserError.NO_ERROR), "auth the user success");
        return true;
    }

    /**
     * @dev 修改用户信息
     * @param _json 修改的参数
            name , idValidity, mobile, phone, email, addr, postCode
     * @return _ret 返回值 true false
     */
    function updatePerUser(string _json) returns(bool _ret){

        //判断json格式
        if(!per_User.fromJson(_json)){
            Notify(errno_prefix+uint(PerUserError.BAD_PARAMETER),"json invalid");
            return;
        }

        if(perUserMap[per_User.id].id == address(0)) {
            Notify(errno_prefix+uint(PerUserError.USER_NOT_EXISTS), 'The user dose not exist.');
            return;
        }

        //手机号不能重复
        address _userId = findIdByMobile(per_User.mobile);
        if(_userId != address(0) && per_User.id != _userId) {
            Notify(errno_prefix+uint(PerUserError.MOBILE_ALREADY_EXISTS), "mobile already exists");
            return;
        }

        if (!per_User.name.equals("")) {
            perUserMap[per_User.id].name = per_User.name;
        }

        if (per_User.idValidity != 0) {
            perUserMap[per_User.id].idValidity = per_User.idValidity;
        }

        if (!per_User.mobile.equals("")) {
            perUserMap[per_User.id].mobile = per_User.mobile;
        }

        if (!per_User.phone.equals("")) {
            perUserMap[per_User.id].phone = per_User.phone;
        }

        if (!per_User.email.equals("")) {
            perUserMap[per_User.id].email = per_User.email;
        }

        if (!per_User.addr.equals("")) {
            perUserMap[per_User.id].addr = per_User.addr;
        }

        if (!per_User.postCode.equals("")) {
            perUserMap[per_User.id].postCode = per_User.postCode;
        }

        _ret = true;
        Notify(uint(PerUserError.NO_ERROR), "update the user success");
        per_User.reset();
    }

    /**
     * @dev 根据地址获取个人信息
     * @return _ret 返回值 true false
     */
    function findByPerId(address _id) constant returns(string _ret){
        _ret = "{\"ret\":0, \"message\": \"success\", \"data\":{\"total\":0,\"items\":[]}}";

        if(perUserMap[_id].id == address(0)) {
            return;
        }
        _ret = "{\"ret\": 0, \"message\": \"success\", \"data\": {\"total\": 1, \"items\":[";
        _ret = _ret.concat(perUserMap[_id].toJson(),"]}}");
    }

    /**
     * @dev 根据身份证件号码查询用户
     * @param _json idType-身份证件类型，idNo-身份证件号码
     * @return _ret 
     */
    function findByIdNo(string _json) constant returns(string _ret){
        _ret = "{\"ret\": 0, \"message\": \"success\", \"data\": {\"total\": 0, \"items\":[]}}";

        uint _idType = uint(_json.getIntValueByKey("idType"));
        string memory _idNo = _json.getStringValueByKey("idNo") ;
        if(_idType == 0 || _idNo.equals("")) {
            return;
        }

        address _userId = findIdByIdNo(_idNo, LibCommonEnum.IdType(_idType));
        if(_userId == address(0)) {
            return;
        }

        _ret = "{\"ret\": 0, \"message\": \"success\", \"data\": {\"total\": 1, \"items\":[";
        _ret = _ret.concat(perUserMap[_userId].toJson(),"]}}");
    }

    /**
     * @dev 根据id查询用户姓名
     * @param _id 用户address
     * @return _ret 
     */
    function findNameById(address _id) getOm constant returns (uint){
        string memory name;
        if(userExists(_id) == 1) {
            return LibStack.push(perUserMap[_id].name);
        }
        om = OrgUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrgUserManager", "0.0.1.0"));
        return om.findNameById(_id);
    }

    function findById(address _id) getOm constant returns (uint){
        string memory json = "";
        if(userExists(_id) == 1) {
            json = perUserMap[_id].toJson();
        }
        return LibStack.push(json);
    }

    /**
     * @dev 根据登录名查询用户
     */
    function findByPerLoginName(string _loginName) constant returns(string _ret){
        _ret = "{\"ret\": 0, \"message\": \"success\", \"data\": {\"total\": 0, \"items\":[]}}";
        for(uint i = 0; i < addrList.length; i++){
            if(perUserMap[addrList[i]].loginName.equals(_loginName)){
                _ret = "{\"ret\": 0, \"message\": \"success\", \"data\": {\"total\": 1, \"items\":[";
                _ret = _ret.concat(perUserMap[addrList[i]].toJson(), "]}}");
                return;
            }
        }
    }

    /* 根据身份证号/手机号查询用户 */
    function findByIdNoOrMobile(string _no) constant returns (string _ret) {
        _ret = "{\"ret\": 0, \"message\": \"success\", \"data\": {\"total\": 0, \"items\":[]}}";
        address _userId = findIdByIdNo(_no, LibCommonEnum.IdType.IDENTITY_CARD);    //作为身份证号查找
        if(_userId == address(0)) {
            _userId = findIdByMobile(_no);  //作为手机号查找
        }
        if(_userId == address(0)) {
            return;
        }
        _ret = "{\"ret\": 0, \"message\": \"success\", \"data\": {\"total\": 1, \"items\":[";
        _ret = _ret.concat(perUserMap[_userId].toJson(),"]}}");
    }

    /**
     * @dev 根据用户状态分页获取用户
     * @param _json 查询条件 status-账户状态, pageSize-页面大小, pageNo-页面号
     * @return _ret 
     */
    function pageByStatus(string _json) constant returns(string _ret){
        _ret = "{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":0,\"items\":[]}}";
        if (addrList.length <= 0) {
            return;
        }

        uint _pageSize = uint(_json.getIntValueByKey("pageSize"));
        uint _pageNo = uint(_json.getIntValueByKey("pageNo"));
        uint _status = uint(_json.getIntValueByKey("status"));

        uint _startIndex = _pageSize * _pageNo;

        if (_startIndex >= addrList.length) {
          return;
        }

        _ret = "";

        string memory _data;
        uint _count = 0; //满足条件的消息数量
        uint _total = 0; //满足条件的指定页数的消息数量
        for (uint i = 0; i < addrList.length; i++) {
            per_User = perUserMap[addrList[i]];
            if (_status != 0 && uint(per_User.status) != _status) {
              continue;
            }

            if (_count++ < _startIndex) {
                continue;
            }

            if (_total < _pageSize) {
                if (_total > 0) {
                  _data = _data.concat(",");
                }
                _total ++;
                _data = _data.concat(per_User.toJson());
            }
        }
        _ret = _ret.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _count.toString(), ",\"items\":[");
        _ret = _ret.concat(_data, "]}}");
    }

    /**
     * @dev 判断用户是否是个人用户
     * @param _id 用户地址
     * @return _ret 
     */
    function userExists(address _id) constant returns(uint) {
        if(perUserMap[_id].status == LibPerUser.PerUserStatus.VALID) {
            return 1;
        }
    }

    /**
     * @dev 注销用户
     * @param _id 用户地址
     * @return _ret 
     */
    function delPerUser(address _id) {
        if(perUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(PerUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        if(perUserMap[_id].status != LibPerUser.PerUserStatus.VALID) {
            Notify(errno_prefix+uint(PerUserError.USER_STATUS_ERROR), "The user status is not valid.");
            return;
        }
        bi = CsdcBaseInterface(rm.getContractAddress("CsdcModule", "0.0.1.0", "BizManager", "0.0.1.0"));
        if (bi.hasTodo(_id)) {
            Notify(errno_prefix+uint(PerUserError.USER_HAS_TODO), "The user has to-do business.");
            return;
        }
        perUserMap[_id].status = LibPerUser.PerUserStatus.INVALID;
        Notify(uint(PerUserError.NO_ERROR), "delete the user success");
        // UserManager userMgr = UserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "UserManager","0.0.1.0"));
        // _ret = userMgr.updateAccountStatus(_id,0);
    }

    /**
     * @dev 锁定用户
     * @param _id 用户地址
     * @return _ret 
     */
    function lockPerUser(address _id) {
        if(perUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(PerUserError.USER_STATUS_ERROR), "The user is not existed");
            return;
        }

        if(perUserMap[_id].status != LibPerUser.PerUserStatus.VALID) {
            Notify(errno_prefix+uint(PerUserError.USER_STATUS_ERROR), "The user status is not valid.");
            return;
        } 
        perUserMap[_id].status = LibPerUser.PerUserStatus.LOCKED;
        Notify(uint(PerUserError.NO_ERROR), "lock the user success");
        // UserManager userMgr = UserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "UserManager","0.0.1.0"));
        // _ret = userMgr.updateAccountStatus(_id,2);
    }

    /**
     * @dev 解锁用户
     * @param _id 用户地址
     * @return _ret 
     */
    function unlockPerUser(address _id) {
        if(perUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(PerUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        if(perUserMap[_id].status != LibPerUser.PerUserStatus.LOCKED) {
            Notify(errno_prefix+uint(PerUserError.USER_STATUS_ERROR), "The user status is not locked.");
            return;
        } 
        perUserMap[_id].status = LibPerUser.PerUserStatus.VALID;
        Notify(uint(PerUserError.NO_ERROR), "unlock the user success");
        // UserManager userMgr = UserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "UserManager","0.0.1.0"));
        // _ret = userMgr.updateAccountStatus(_id,1);
    }


    /* for test only */
    function deleteByAddress(address _userAddr) {
        // first update id list
        delete tempList;
        for (uint i = 0; i < addrList.length; i++)  {
            if(_userAddr == perUserMap[addrList[i]].id) {
                continue;
            }
            else {
                tempList.push(addrList[i]);
            }
        }

        //copy elements
        delete addrList;
        for (uint j = 0; j < tempList.length; ++j) {
            addrList.push(tempList[j]);
        }

        delete perUserMap[_userAddr];
    }

    /* 以下是internal方法 */
    /* 根据身份证号查询address */
    function findIdByIdNo(string _idNo, LibCommonEnum.IdType _idType) internal constant returns (address) {
        for(uint i = 0; i<addrList.length; i++){
            if(perUserMap[addrList[i]].idNo.equals(_idNo) && perUserMap[addrList[i]].idType == _idType
            && perUserMap[addrList[i]].status != LibPerUser.PerUserStatus.NONE && perUserMap[addrList[i]].status != LibPerUser.PerUserStatus.INVALID){
                return addrList[i];
            }
        }
    }

    function findIdByMobile(string _mobile) internal constant returns (address) {
        for(uint i = 0; i<addrList.length; i++){
            if(perUserMap[addrList[i]].mobile.equals(_mobile)
            && perUserMap[addrList[i]].status != LibPerUser.PerUserStatus.NONE && perUserMap[addrList[i]].status != LibPerUser.PerUserStatus.INVALID){
                return addrList[i];
            }
        }
    }

    event Notify(uint _errno, string _info);

    /* for CsdcBaseInterface */
    function hasTodo(address _userId) constant returns (bool) {}
    function updatePledgeSecurityOfOneDisSecPledgeApply(uint disSecPedgeApplyId, uint pledgeSecuityId, uint profitAmount) returns(bool _ret) {
        
    }
}