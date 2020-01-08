pragma solidity ^0.4.12;
/**
*@file      OrgUserManager.sol
*@author    xuhui
*@time      2016-3-6
*@desc      the defination of OrgUserManager
*/


import "./utillib/LibStack.sol";
import "./sysbase/OwnerNamed.sol";
import "./csdc_library/LibOrgUser.sol";

import "./UserManager.sol";


/* @brief OrgUserManager definition */
contract OrgUserManager is OwnerNamed/* , CsdcBaseInterface */ {
    using LibOrgUser for *;
    using LibString for *;
    using LibInt for *;

    event Notify(uint _errorno, string _info);
    
    mapping(address=>LibOrgUser.OrgUser)            OrgUserMap;
    address[]                                       addrList;
    address[]                                       keyInfoAddrList;
    address[]                                       tempList;

    LibOrgUser.OrgUser internal                     organ_User;

    /** @brief errno for test case */
    enum OrgUserError {
        NO_ERROR,
        BAD_PARAMETER,
        ID_EMPTY,   
        LOGINNAME_EMPTY,
        PASSWORD_EMPTY,
        ORGANFULLNAME_EMPTY,
        ORGANTYPE_EMPTY,
        BUSINESSLICENSENO_EMPTY,
        LEGALREPRESENTNAME_EMPTY,
        LEGALREPRESENTIDCERTTYPE_EMPTY,
        LEGALREPRESENTIDCERTNO_EMPTY,
        APPOINTERNAME_EMPTY,
        APPOINTERIDCERTTYPE_EMPTY,
        APPOINTERIDCERTNO_EMPTY,
        APPOINTERMOBILE_EMPTY,
        COMPANYNAME_EMPTY,
        ORGANIZATIONNO_EMPTY,
        CODEVALIDITY_EMPTY,
        ESTABLISHTIME_EMPTY,
        REGISTEREDCAPITAL_EMPTY,
        REGISTEREDADDR_EMPTY,
        LEGALREPRESENTTEL_EMPTY,
        LEGALREPRESENTMOBILE_EMPTY,
        LEGALREPRESENTMAIL_EMPTY,
        APPOINTERDEPARTMENT_EMPTY,
        APPOINTERPOST_EMPTY,
        APPOINTERTEL_EMPTY,
        APPOINTERFAX_EMPTY,
        APPOINTERMAIL_EMPTY,
        APPOINTERADDR_EMPTY,
        APPOINTERPOSTCODE_EMPTY,
        BUSINESSLICENSECOPYFILE_EMPTY,
        CERTIFICATEOFLEGALREPRESENTCOPYFILE_EMPTY,
        LEGALREPRESENTAUTHCOPYFILE_EMPTY,
        LEGALREPRESENTIDCOPYFILE_EMPTY,
        APPOINTERIDCOPYFILE_EMPTY,
        ORGANFULLNAME_ALREADY_EXISTS,
        BUSINESSLICENSENO_ALREADY_EXISTS,
        ORGANIZATIONNO_ALREADY_EXISTS,
        REGISTEREDADDR_ALREADY_EXISTS,
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
    uint errno_prefix = 16400;

    UserManager userMgr;

    function OrgUserManager() {
        register("CsdcModule", "0.0.1.0", "OrgUserManager", "0.0.1.0");
    }

    /**
     * @dev 新增机构用户
     * @param _strJson _strJson表示的用户资料
     * @return _ret 返回值 true false
     */
    function insertOrgUser(string _strJson) public returns(bool _ret) {
        //1.判断json格式
        if(!organ_User.fromJson(_strJson)){
            Notify(uint(OrgUserError.BAD_PARAMETER),"json invalid");
            _ret = false;
            return;
        }

        if(organ_User.id == address(0)) {
            Notify(errno_prefix+uint(OrgUserError.ID_EMPTY), "id cannot be empty.");
            return;
        }

        if(OrgUserMap[organ_User.id].id != address(0)) {
            Notify(errno_prefix+uint(OrgUserError.ADDRESS_ALREADY_EXISTS), "address already exists.");
            return;
        }

        for(uint i = 0; i < addrList.length; i++){
            if(OrgUserMap[addrList[i]].loginName.equals(organ_User.loginName)){
                Notify(uint(OrgUserError.ACCOUNT_ALREDY_EXISTS), "loginName already exists");
                _ret = false;
                return;
            }
        }

        //多个用户可以有一个businessNo
        // for(uint j = 0; j < addrList.length; j++){
        //     if(OrgUserMap[addrList[j]].businessLicenseNo.equals(organ_User.businessLicenseNo)){
        //         Notify(uint(OrgUserError.BUSINESSLICENSENO_ALREADY_EXISTS), "businesslicenseNo already exists");
        //         _ret = false;
        //         return;
        //     }
        // }

        address _id = organ_User.id;
        //2.拼接传入底层的json
        string memory _userJson = "{";
        _userJson = _userJson.concat(uint(_id).toAddrString().toKeyValue("userAddr"),",");
        _userJson = _userJson.concat(organ_User.loginName.toKeyValue("account"),",");
        _userJson = _userJson.concat(organ_User.loginName.toKeyValue("name"),",");
        _userJson = _userJson.concat("\"accountStatus\": 1", ",");
        _userJson = _userJson.concat("\"departmentId\": \"orguserdpt\"", ",");
        _userJson = _userJson.concat("\"roleIdList\": [\"csdc_role_301\", \"csdc_role_302\"]");
        // _userJson = _userJson.concat("\"roleIdList\": \"['orguserrole']\"");
        _userJson = _userJson.concat("}");
        LibLog.log("insertUser:", _userJson);
        //3.底层判断是否可以插入成功
        userMgr = UserManager(rm.getContractAddress("SystemModuleManager", "0.0.1.0", "UserManager", "0.0.1.0"));
        if(userMgr.insert(_userJson)==0){
            organ_User.status = LibOrgUser.AccountStatus.VALID;
            addrList.push(_id);
            OrgUserMap[_id] = organ_User;
            organ_User.reset();
            Notify(0,"insert a user success");
            return;
        }else{
            Notify(0,"insert a user error");
            return;
        }
    }

    /**
     * @dev 更新机构用户
     * @param _strJson _strJson表示的用户信息
     * @return _ret 返回值 true false
     */
    function updateOrgInfo(string _strJson) public returns(bool _ret){
        _ret = false;
        //判断json格式
        if(!organ_User.fromJson(_strJson)){
            Notify(0,"insert bad json");
            _ret = false;
            return;
        }

        if(OrgUserMap[organ_User.id].id == address(0)) {
            Notify(errno_prefix+uint(OrgUserError.USER_NOT_EXISTS), 'The user dose not exist.');
            return;
        }

        if(!organ_User.name.equals("")) {
            OrgUserMap[organ_User.id].name = organ_User.name;
        }

        // if(!organ_User.loginName.equals("")) {
        //     OrgUserMap[organ_User.id].loginName = organ_User.loginName;
        // }

        if(!organ_User.password.equals("")) {
            OrgUserMap[organ_User.id].password = organ_User.password;
        }

        if(!organ_User.organFullName.equals("")) {
            OrgUserMap[organ_User.id].organFullName = organ_User.organFullName;
        }

        if(uint(organ_User.participantsType) != 0) {
            OrgUserMap[organ_User.id].participantsType = organ_User.participantsType;
        }

        if(uint(organ_User.publisherType) != 0) {
            OrgUserMap[organ_User.id].publisherType = organ_User.publisherType;
        }

        if(uint(organ_User.organIdType) != 0) {
            OrgUserMap[organ_User.id].organIdType = organ_User.organIdType;
        }

        if(!organ_User.companyName.equals("")) {
            OrgUserMap[organ_User.id].companyName = organ_User.companyName;
        }

        if(!organ_User.businessLicenseNo.equals("")) {
            OrgUserMap[organ_User.id].businessLicenseNo = organ_User.businessLicenseNo;
        }

        if(organ_User.licenseValidity!=0) {
            OrgUserMap[organ_User.id].licenseValidity = organ_User.licenseValidity;
        }

        if(organ_User.licenseInvalidity!=0) {
            OrgUserMap[organ_User.id].licenseInvalidity = organ_User.licenseInvalidity;
        }

        if(!organ_User.organizationNo.equals("")) {
            OrgUserMap[organ_User.id].organizationNo = organ_User.organizationNo;
        }

        if(organ_User.codeValidity!=0) {
            OrgUserMap[organ_User.id].codeValidity = organ_User.codeValidity;
        }

        if(organ_User.codeInvalidity!=0) {
            OrgUserMap[organ_User.id].codeInvalidity = organ_User.codeInvalidity;
        }

        if(organ_User.establishTime!=0) {
            OrgUserMap[organ_User.id].establishTime = organ_User.establishTime;
        }

        if(organ_User.registeredCapital!=0) {
            OrgUserMap[organ_User.id].registeredCapital = organ_User.registeredCapital;
        }

        if(!organ_User.registeredAddr.equals("")) {
            OrgUserMap[organ_User.id].registeredAddr = organ_User.registeredAddr;
        }

        if(!organ_User.communicationAddr.equals("")) {
            OrgUserMap[organ_User.id].communicationAddr = organ_User.communicationAddr;
        }

        if(!organ_User.legalRepresentName.equals("")) {
            OrgUserMap[organ_User.id].legalRepresentName = organ_User.legalRepresentName;
        }

        if(uint(organ_User.legalRepresentIDCertType) != 0) {
            OrgUserMap[organ_User.id].legalRepresentIDCertType = organ_User.legalRepresentIDCertType;
        }

        if(!organ_User.legalRepresentIdCertNo.equals("")) {
            OrgUserMap[organ_User.id].legalRepresentIdCertNo = organ_User.legalRepresentIdCertNo;
        }

        if(!organ_User.legalRepresentTel.equals("")) {
            OrgUserMap[organ_User.id].legalRepresentTel = organ_User.legalRepresentTel;
        }

        if(!organ_User.legalRepresentMobile.equals("")) {
            OrgUserMap[organ_User.id].legalRepresentMobile = organ_User.legalRepresentMobile;
        }

        if(!organ_User.legalRepresentFax.equals("")) {
            OrgUserMap[organ_User.id].legalRepresentFax = organ_User.legalRepresentFax;
        }

        if(!organ_User.legalRepresentMail.equals("")) {
            OrgUserMap[organ_User.id].legalRepresentMail = organ_User.legalRepresentMail;
        }

        if(uint(organ_User.status) != 0) {
            OrgUserMap[organ_User.id].status = organ_User.status;
        }

        organ_User.reset();
        Notify(0,"update success");
        _ret = true;
    }

    /**
     * @dev 根据用户id地址获取机构用户信息
     * @param _id 用户地址
     * @return _ret 返回机构信息
     */
    function findByOrgId(address _id) constant public returns(string _ret){
        _ret = "{\"ret\":0,\"data\":{\"total\":0,\"items\":[]}}";

        if(OrgUserMap[_id].id == address(0)) {
            return;
        }
        _ret = "{\"ret\": 0, \"data\": {\"total\": 1, \"items\":[";
        _ret = _ret.concat(OrgUserMap[_id].toJson(),"]}}");
        Notify(0,_ret);
    }

    /**
     * @dev 根据登录名获取机构用户信息
     * @param _loginName 用户名
     * @return _ret 返回机构信息
     */
    function findByOrgLoginName(string _loginName) constant public returns(string _ret){
        bool _flag = false;
        _ret = "{\"ret\": 0, \"data\": {\"total\": 1, \"items\":[";
        for(uint i = 0; i < addrList.length; i++){
            if(OrgUserMap[addrList[i]].loginName.equals(_loginName)){
                _ret = _ret.concat(OrgUserMap[addrList[i]].toJson(),"]}}");
                _flag = true;
                return;
            }
        }
        if(!_flag){
            _ret = "{\"ret\": 0, \"data\": {\"total\": 0, \"items\":[]}}";
        }
    }

    /**
     * @dev 删除机构用户
     * @param _id 机构id
     * @return _ret 返回值 true false
     */
    function delOrgUser(address _id) public returns (bool _ret){
        if(OrgUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(OrgUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        if(OrgUserMap[_id].status != LibOrgUser.AccountStatus.VALID) {
            Notify(errno_prefix+uint(OrgUserError.USER_STATUS_ERROR), "The user status is not valid.");
            return;
        }
        OrgUserMap[_id].status = LibOrgUser.AccountStatus.INVALID;
        Notify(uint(OrgUserError.NO_ERROR), "delete the user success");
        return true;
    }

    /**
     * @dev 锁定机构用户
     * @param _id 机构id
     * @return _ret 返回值 true false
     */
    function lockOrgUser(address _id) public returns (bool _ret){
        if(OrgUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(OrgUserError.USER_STATUS_ERROR), "The user is not existed");
            return;
        }

        if(OrgUserMap[_id].status != LibOrgUser.AccountStatus.VALID) {
            Notify(errno_prefix+uint(OrgUserError.USER_STATUS_ERROR), "The user status is not valid.");
            return;
        }
        OrgUserMap[_id].status = LibOrgUser.AccountStatus.LOCKED;
        Notify(uint(OrgUserError.NO_ERROR), "lock the user success");
        return true;
    }

    /**
     * @dev 解锁机构用户
     * @param _id 机构id
     * @return _ret 返回值 true false
     */
    function unlockOrgUser(address _id) public returns (bool _ret){
        if(OrgUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(OrgUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        if(OrgUserMap[_id].status != LibOrgUser.AccountStatus.LOCKED) {
            Notify(errno_prefix+uint(OrgUserError.USER_STATUS_ERROR), "The user status is not locked.");
            return;
        } 
        OrgUserMap[_id].status = LibOrgUser.AccountStatus.VALID;
        Notify(uint(OrgUserError.NO_ERROR), "unlock the user success");
        return true;
    }

    /**
    * @dev 查询机构用户
    * @param _cond json表示的查询条件--id, organFullName, businessLicense, status ,pageNo , pageSize
    * @return _ret
    */
    LibOrgUser.Cond _cond;
    function pageByOrgUser(string _json) constant public returns(string _ret){
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": [] }}";
        if(addrList.length<=0){
            return;
        }
        _cond.organFullName = _json.getStringValueByKey("organFullName");
        _cond.businessLicenseNo = _json.getStringValueByKey("businessLicenseNo");
        _cond.organizationNo = _json.getStringValueByKey("organizationNo");
        _cond.status = LibOrgUser.AccountStatus(_json.getIntValueByKey("status"));
        _cond.pageSize = uint(_json.getIntValueByKey("pageSize"));
        _cond.pageNo = uint(_json.getIntValueByKey("pageNo"));

        uint _startIndex = _cond.pageSize * _cond.pageNo;
        if (_startIndex >= addrList.length) {
          return;
        }
        _ret = "";
        
        string memory _data;
        uint _count = 0; //满足条件的消息数量
        uint _total = 0; //满足条件的指定页数的消息数量

        for(uint i = 0; i< addrList.length; i++){
            organ_User.reset();
            organ_User = OrgUserMap[addrList[i]];
            if (!_cond.organFullName.equals("") && !_cond.organFullName.equals(organ_User.organFullName)) {
                continue;
            }  
            if (!_cond.businessLicenseNo.equals("") && !_cond.businessLicenseNo.equals(organ_User.businessLicenseNo)) {
                continue;
            }  
            if (!_cond.organizationNo.equals("") && !_cond.organizationNo.equals(organ_User.organizationNo)) {
                continue;
            }
            if(_cond.status != LibOrgUser.AccountStatus.NONE && _cond.status != organ_User.status) {
                continue;
            }
            if (_count++ < _startIndex) {
                continue;
            }

            if(_cond.pageSize == 0 && _cond.pageNo == 0){
                if (_total > 0) {
                      _data = _data.concat(",");
                    }
                    _total ++;
                    _data = _data.concat(organ_User.toJson());
            }else{
                if (_total < _cond.pageSize) {
                    if (_total > 0) {
                      _data = _data.concat(",");
                    }
                    _total ++;
                    _data = _data.concat(organ_User.toJson());
                }
            }
        }
        _ret = _ret.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _count.toString(), ",\"items\":[");
        _ret = _ret.concat(_data, "]}}");
    }
    /**
    * @dev 根据用户地址获取用户密码
    * @param _id 用户id
    * @return _ret 返回值
    */
    function findPasswordByUserId(address _id) constant public returns(string _ret){

        if(OrgUserMap[_id].id == address(0)) {
            Notify(errno_prefix+uint(OrgUserError.USER_NOT_EXISTS), "The user is not existed");
            return;
        }

        _ret = OrgUserMap[_id].password;
        Notify(uint(OrgUserError.NO_ERROR), "find the user password success");
        
    }
    /**
     * @dev 根据登录名获取机构用户密码
     * @param _loginName 登录名
     * @return _ret 返回值
     */
    function findPasswordByLoginName(string _loginName) constant public returns(string _ret){
        
        for(uint i = 0; i < addrList.length; i++){
            if(OrgUserMap[addrList[i]].loginName.equals(_loginName)){
                _ret = OrgUserMap[addrList[i]].password;
                Notify(errno_prefix+uint(OrgUserError.NO_ERROR), "find the user password success");
                return;
            }
        }
        Notify(errno_prefix+uint(OrgUserError.USER_NOT_EXISTS), "The user is not existed");
        
    }

    /**
     * @dev 根据id查询机构全称
     * @param _id 机构用户address
     * @return _ret 
     */
    function findNameById(address _id) constant returns (uint){
        string memory name = "";
        if(OrgUserMap[_id].id != address(0)) {
            name = OrgUserMap[_id].organFullName;
        }
        return LibStack.push(name);
    }

    function findById(address _id) constant returns (uint){
        string memory json = "";
        if(OrgUserMap[_id].id != address(0)) {
            json = OrgUserMap[_id].toJson();
        }
        return LibStack.push(json);
    }

    function userExists(address _userAddr) constant public returns(uint _ret) {
        if (OrgUserMap[_userAddr].status == LibOrgUser.AccountStatus.VALID) {
            return 1;
        } else {
            return 0;
        }
    }

    function deleteByAddress(address _userAddr) {
        delete tempList;
        for (uint i = 0; i < addrList.length; i++)  {
            if(_userAddr == OrgUserMap[addrList[i]].id) {
                continue;
            }
            else {
                tempList.push(addrList[i]);
            }
        }
        delete addrList;
        for (uint j = 0; j < tempList.length; ++j) {
            addrList.push(tempList[j]);
        }
        delete OrgUserMap[_userAddr];
    }

    function listAll() constant public returns(string _ret){
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": [] }}";
        if(addrList.length<=0){
            return;
        }
        _ret = "{\"ret\":0, \"message\": \"success\", \"data\":{";
        _ret = _ret.concat(uint(addrList.length).toKeyValue("total"), ",\"items\":[");
        for(uint i = 0; i < addrList.length; i++){    
            if(i==addrList.length-1){
                _ret = _ret.concat(OrgUserMap[addrList[i]].toJson());
            }else{
                _ret = _ret.concat(OrgUserMap[addrList[i]].toJson(),",");
            }
        }
        _ret = _ret.concat("]}}");
        Notify(0,_ret);
    }
}