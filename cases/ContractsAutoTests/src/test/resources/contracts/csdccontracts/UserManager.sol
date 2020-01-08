pragma solidity ^0.4.12;
/**
*@file      UserManager.sol
*@author    kelvin
*@time      2016-11-29
*@desc      the defination of UserManager contract
*/

import "./library/LibUser.sol";
import "./sysbase/OwnerNamed.sol";
import "./interfaces/IRoleManager.sol";
import "./interfaces/IUserManager.sol";
import "./interfaces/IDepartmentManager.sol";

contract UserManager is OwnerNamed, IUserManager {
 
    using LibUser for *;
    using LibString for *;
    using LibInt for *;

    event Notify(uint _errno, string _info);
    
    mapping(address=>LibUser.User)      userMap;
    address[]                           addrList;
    address[]                           tempList;

    //inner setting member
    LibUser.User internal               m_User;

    //temp roleIdList
    string[] tmpArray;
    
    LibUser.ModuleRoles[] tmpModuleRoles;
    LibUser.ModuleRoles tmpModuleRole;

    uint revision;

    enum UserError {
        NO_ERROR,
        BAD_PARAMETER,
        NAME_EMPTY,
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
        EMAIL_EXISTS,
        MOBILE_EXISTS

    }
    

    function UserManager() {
        revision = 0;
        register("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");

    }    

    function updateUserInfoDb(address _userAddr, uint _accountStatus, uint _deleteStatus, uint _state) constant private returns (uint _ret) {
            log("UserManager.sol", "updateUserInfoDb");
            address __userAddr = getOwnerAddrByAddr(_userAddr);
            log("_userAddr->__userAddr", uint(_userAddr).toAddrString(), uint(__userAddr).toAddrString());
            string memory value = uint(__userAddr).toAddrString();//0x002ed552b386464419bff1d4f45bbaec2ea24177,1,0,1
            value = value.concat("|", _accountStatus.toString());
            value = value.concat("|", _deleteStatus.toString());
            value = value.concat("|", _state.toString());
            log("user info:", value);
            _ret = writedb("userInfo|update", uint(__userAddr).toAddrString(), value);
            if (0 != _ret)
                log("UserManager.sol", "updateUserInfoDb failed.");
            else
                log("UserManager.sol", "updateUserInfoDb success.");
    }

    /**
    * 根据uuid获取用户信息
    * since client v1.0.0
    * @param _uuid 文件证书为钱包uid，u-key为喷码
    * @return 用户信息
    */
    function findByUuid(string _uuid) constant public returns(string _strjson) {
        string memory items;
        for (uint i = 0; i < addrList.length; i++)
        {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && _uuid.equals(userMap[addrList[i]].uuid))
            {
                items = userMap[addrList[i]].toJson();
                break;
            }
        }
        uint len = itemStackPush(items, items.equals("") ? uint(0) : uint(1));
        _strjson = LibStack.popex(len);
    }

    /**
    * 重置用户密码（重置数据）
    * since client v1.0.0
    * @param _userAddr 用户地址
    * @param _ownerAddr 用户状态 
    * @return true success | false fail
    */
    function resetPasswd(address _userAddr,address _ownerAddr,string _publilcKey,string _cipherGroupKey,string _uuid) public returns(bool _ret) {
        _ret = false;
        if(_userAddr == 0 || _ownerAddr == 0){
            log("resetPasswd _oldUserAddr is null..", "UserManager");
            errno = 15200 + uint(UserError.BAD_PARAMETER);
            Notify(errno, "operator no permission");
            return;
        }
        if (msg.sender != owner) {
            log("resetPasswd msg.sender is not owner", "UserManager");
            if (__checkWritePermission(msg.sender, userMap[_userAddr].departmentId) == 0) {
                log("resetPasswd user operator no permission", "UserManager");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return;
            }
        }
        userMap[_userAddr].ownerAddr = _ownerAddr;
        userMap[_userAddr].publicKey = _publilcKey;
        userMap[_userAddr].cipherGroupKey = _cipherGroupKey;
        userMap[_userAddr].uuid = _uuid;
        userMap[_userAddr].updateTime = now * 1000;

        // insert new record to leveldb
        splitRolesByModuleName(_userAddr);
        for (uint i = 0 ; i < tmpModuleRoles.length; ++i) {
            log("resetPasswd->save Data to DB,userAddr:",uint(userMap[_userAddr].userAddr).toAddrString());
            log("resetPasswd->save Data to DB,roleId:",tmpModuleRoles[i].roleIds);
            updateAuthorizeUserRole(tmpModuleRoles[i].moduleName, tmpModuleRoles[i].moduleVersion, userMap[_userAddr].userAddr, tmpModuleRoles[i].roleIds);
        }
        updateUserInfoDb(_userAddr, userMap[_userAddr].accountStatus, userMap[_userAddr].deleteStatus, uint(userMap[_userAddr].state));

        _ret = true;
        errno = uint(UserError.NO_ERROR);
        revision++;
        log("resetPasswd success ...", "UserManager");
        Notify(errno, "update user status success");
        return;
    }

    /**
    *@desc  按照用户地址获取用户状态
    *@param _userAddr 用户地址
    *@_ret  0: 非法用户，1：合法用户，2： 已登录用户
    */
    function getUserState(address _userAddr) constant public returns (uint _state) {
        _state = 0;
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID) {
            return;
        }

        _state = uint(userMap[_userAddr].state);
    }

    /**
    *@desc  按照账号名字获取状态
    *@param _userAddr 用户地址
    *@_ret  0: 非法用户，1：合法用户，2： 已登录用户
    */
    function getAccountState(string _account) constant public returns (uint _state) {
        _state = 0;
        for (uint i = 0; i < addrList.length; i++)
        {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && _account.equals(userMap[addrList[i]].account))
            {
                _state = uint(userMap[addrList[i]].state);
                return;
            }
        }
    }


    /**
    *@desc  按照登录名字和账号状态获取分页用户列表
    *@param _accountStatus 账号状态_pageNo 页面号, _pageSize 页面大小
    *@_ret  响应json字符串
    */
    function pageByAccountStatus(uint _accountStatus, uint _pageNo, uint _pageSize) public constant returns (string _strjson) {
        uint startIndex = uint(_pageNo * _pageSize);
        uint endIndex = uint(startIndex + _pageSize - 1);
        bool flag = true;
        uint len = 0;
        uint count = 0;
        if (startIndex >= addrList.length)
        {
            flag = false;
            return;
        }

        if (endIndex >= addrList.length) {
            endIndex = addrList.length - 1;
        }
        len = LibStack.push("");
        for (uint index = startIndex; flag && index <= endIndex; index++) {
            if (userMap[addrList[index]].state != LibUser.UserState.USER_INVALID &&
            userMap[addrList[index]].accountStatus == _accountStatus) {
                if (count > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(userMap[addrList[index]].toJson());
                count++;
            }
        }
        len = itemStackPush(LibStack.popex(len), getUserCount());
        _strjson = LibStack.popex(len);
    }

    // added by liaoyan. 2016-12-27
    function findByAddress(address _userAddr) constant public returns(string _ret) {
        string memory items;
        if (userMap[_userAddr].state != LibUser.UserState.USER_INVALID) {
            items = userMap[_userAddr].toJson();
        }
        uint len = itemStackPush(items, getUserCount());
        _ret = LibStack.popex(len);
    }
    
    /**
    *@desc  按照名字查询用户
    *@param _name 用户名字
    *@_ret  响应json字符串
    */
    function findByLoginName(string _name) constant public returns(string _strjson) {
        if (addrList.length <= 0) {
            _strjson = "{\"ret\":0,\"data\":{\"total\":0,\"items\":[]}}";
            return;
        }

        uint len = 0;
        uint counter = 0;
        len = LibStack.push("");
        for (uint i = 0; i < addrList.length; i++)
        {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && _name.equals(userMap[addrList[i]].name))
            {
                if (counter > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(userMap[addrList[i]].toJson());
                counter++;
            }
        }

        len = itemStackPush(LibStack.popex(len), getUserCount());
        _strjson = LibStack.popex(len);
    }
    
    /**
    *@desc  按照账号查询用户
    *@param _account 用户账号
    *@_ret  响应json字符串
    */
    function findByAccount(string _account) constant public returns(string _strjson){
        string memory items;
        for (uint i = 0; i < addrList.length; i++)
        {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && _account.equals(userMap[addrList[i]].account))
            {
                items = userMap[addrList[i]].toJson();
                break;
            }
        }
        uint len = itemStackPush(items, items.equals("") ? uint(0) : uint(1));
        _strjson = LibStack.popex(len);
    }
    
    /**
    *@desc  按照电话查询用户
    *@param _mobile 用户电话
    *@_ret  响应json字符串
    */
    function findByMobile(string _mobile) constant public returns(string _strjson) {
        uint len = 0;
        uint counter = 0;
        len = LibStack.push("");
        for (uint i = 0; i < addrList.length; i++)
        {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && _mobile.equals(userMap[addrList[i]].mobile))
            {
                if (counter > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(userMap[addrList[i]].toJson());
                counter++;
            }
        }
        len = itemStackPush(LibStack.popex(len), getUserCount());
        _strjson = LibStack.popex(len);
    }
    
    /**
    *@desc  按照email查询用户
    *@param _email 用户email
    *@_ret  响应json字符串
    */
    function findByEmail(string _email) constant public returns(string _strjson){
        uint len = 0;
        uint counter = 0;
        len = LibStack.push("");
        for (uint i = 0; i < addrList.length; i++)
        {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && _email.equals(userMap[addrList[i]].email)) 
            {
                if (counter > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(userMap[addrList[i]].toJson());
                counter++;
            }
        }
        len = itemStackPush(LibStack.popex(len), getUserCount());
        _strjson = LibStack.popex(len);
    }

    /**
    *@desc  按照部门id查询用户列表   
    *@param departmentId 
    *@ret   响应json字符串
    */
    function findByDepartmentId(string _departmentId) constant public returns(string _strjson) {
        uint len = 0;
        uint counter = 0;
        len = LibStack.push("");
        for (uint i = 0; i < addrList.length; i++) {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && _departmentId.equals(userMap[addrList[i]].departmentId)) {
                if (counter > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(userMap[addrList[i]].toJson());
                counter++;
            }
        }

        len = itemStackPush(LibStack.popex(len), getUserCount());
        _strjson = LibStack.popex(len);
    }

    /**
    *@desc  按照部门id查询当前部门树下的所有用户列表   
    *@param departmentId 
    *@ret   响应json字符串
    */
    function findByDepartmentIdTree(string _departmentId, uint _pageNum, uint _pageSize) constant public returns(string _strjson) {
        delete tmpArray;
        __getDepartmentIdTree(_departmentId, tmpArray);

        uint len = 0;
        uint n = 0;
        uint m = 0;
        len = LibStack.push("");
        for (uint i = 0; i < addrList.length; i++) {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && userMap[addrList[i]].departmentId.inArray(tmpArray)) {
                if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                    if (m > 0) {
                        len = LibStack.append(",");
                    }
                    len = LibStack.append(userMap[addrList[i]].toJson());
                    m++;
                }
                if (n >= (_pageNum+1)*_pageSize) {
                    break;
                }
                n++;
            }
        }

        len = itemStackPush(LibStack.popex(len), getUserCount());
        _strjson = LibStack.popex(len);
    }

    /**
    * 根据部门ID及检索条件查看用户列表
    * @param _status 用户状态 0 all ,1 disabled,2 enabled
    * @param _name 搜索关键字
    * @param _departmentId 部门ID
    * @return   响应json字符串
    */
    function findByDepartmentIdTreeAndContion(uint _status,string _name,string _departmentId, uint _pageNum, uint _pageSize) constant public returns(string _strjson) {
        // 0 "",admin
        delete tmpArray;
        __getDepartmentIdTree(_departmentId, tmpArray);

        uint len = 0;

        uint n = 0;
        uint m = 0;
        uint total = 0;
        for (uint i = 0; i < addrList.length; i++) {
            uint tmpStatus1 = 0;
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && userMap[addrList[i]].departmentId.inArray(tmpArray)) {
                if((bytes(_name).length != 0 && userMap[addrList[i]].name.indexOf(_name) == -1)){
                    continue;
                }
                if(_status == 1) tmpStatus1 = 0;
                if(_status == 2) tmpStatus1 = 1;
                if(_status != 0 && userMap[addrList[i]].status != tmpStatus1) {
                    continue;
                }
                total++;
            }
        }
        len = LibStack.push("");
        for (i = 0; i < addrList.length; i++) {
            uint tmpStatus = 0;
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID && userMap[addrList[i]].departmentId.inArray(tmpArray)) {
                // name and status
                if((bytes(_name).length != 0 && userMap[addrList[i]].name.indexOf(_name) == -1)){
                    continue;
                }
                if(_status == 1) tmpStatus = 0;
                if(_status == 2) tmpStatus = 1;
                if(_status != 0 && userMap[addrList[i]].status != tmpStatus) {
                    continue;
                }
                if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                    if (m > 0) {
                        len = LibStack.append(",");
                    }
                    len = LibStack.append(userMap[addrList[i]].toJson());
                    m++;
                }
                if (n >= (_pageNum+1)*_pageSize) {
                    break;
                }
                n++;
            }
        }

        len = itemStackPush(LibStack.popex(len), total);
        _strjson = LibStack.popex(len);
    }

    /**
    * get user info of list by role id   
    * @param _roleId the id of role
    * @return _strjson
    */
    function findByRoleId(string _roleId) constant public returns(string _strjson) {
        log("into findByRoleId...");
        if (addrList.length <= 0) {
            _strjson = "{\"ret\":0,\"data\":{\"total\":0,\"items\":[]}}";
            return;
        }

        uint len = 0;
        len = LibStack.push("");
        uint counter = 0;
        for (uint i = 0; i < addrList.length; i++)
        {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID 
                && _roleId.inArray(userMap[addrList[i]].roleIdList))
            {
                if (counter > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(userMap[addrList[i]].toJson());
                counter++;
            }
        }

        len = itemStackPush(LibStack.popex(len), counter);
        _strjson = LibStack.popex(len);
    }

    /**
    *@desc  获取用户部门
    *@param _userAddr 用户地址
    *@_ret  失败返回0， 成功返回部门id
    */
    function getUserDepartmentId(address _userAddr) constant returns(uint _departId) {
        _departId = 0;
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID)
        {
            return;
        }

        _departId = userMap[_userAddr].departmentId.storageToUint();
    }
    
    /**
    *@desc  验证用户角色
    *@param _userAddr 用户地址，_roleId 角色ID
    *@_ret  true: 成功, false: 失败
    */
    function checkUserRole(address _userAddr, string _roleId) constant public returns(uint _ret) {
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID && userMap[_userAddr].roleIdList.length <= 0) {
            return 0;
        }

        /*if (__getDepartmentAdmin(userMap[_userAddr].departmentId) == _userAddr) {
            delete tmpArray;
            __getDepartmentRoleIdList(userMap[_userAddr].departmentId, tmpArray);

            for (uint i = 0; i < tmpArray.length; ++i) {
                if (tmpArray[i].equals(_roleId)) {
                    return 1;
                }
            }
        } else {*/
            for (uint i = 0; i < userMap[_userAddr].roleIdList.length; ++i) {
                if (userMap[_userAddr].roleIdList[i].equals(_roleId)) {
                    return 1;
                }
            }
        //}

        return 0;
    }

    /**
    *@desc  验证用户权限ID
    *@param _userAddr 用户地址，_actionId 权限ID
    *@_ret  1: 成功, 0: 失败
    */
    function checkUserAction(address _userAddr, string _actionId) constant public returns (uint _ret) {
        // add for passwd reset
        address userAddr = getUserAddrByAddr(_userAddr);
        LibUser.User user = userMap[userAddr];
        if (user.state == LibUser.UserState.USER_INVALID) {
            return 0;
        }
        IRoleManager role = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
        for (uint i = 0; i < user.roleIdList.length; ++i) {
            if (role.checkRoleAction(user.roleIdList[i], _actionId) == 1) {
                return 1;
            }
        }
        return 0;
    }

    /**
    *@desc  验证用户权限
    *@param _userAddr 用户地址，_resKey 合约地址，_opKey 方法名的sha3签名
    *@_ret  1: 成功, 0: 失败
    */
    function checkUserPrivilege(address _userAddr, address _contractAddr, string _funcSha3) constant public returns (uint _ret) {
        if (_userAddr == owner) {
            log1("user is owner", "UserManager");
            return 1;
        }
        // add for passwd reset
        address userAddr = getUserAddrByAddr(_userAddr);
        //check if exist the user
        if (userMap[userAddr].state == LibUser.UserState.USER_INVALID) {
            log1("user is not exist", "UserManager");
            return 0;
        }

        // check the account status
        if (userMap[userAddr].accountStatus == uint(LibUser.AccountState.INVALID) || 
            userMap[userAddr].accountStatus == uint(LibUser.AccountState.LOCKED)) {
            log1("user accountStatus invalid", "UserManager");
            return 0;
        }

        // check department in DepartmentManager
        IDepartmentManager departMgr = IDepartmentManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0"));
        if (0 == departMgr.departmentExists(userMap[userAddr].departmentId)) {
            log1("department id not exists", "UserManager");
            return 0;
        }

        IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
        //iterate the role List

        /*if (__getDepartmentAdmin(userMap[_userAddr].departmentId) == _userAddr) {
            delete tmpArray;
            __getDepartmentRoleIdList(userMap[_userAddr].departmentId, tmpArray);

            for (uint i = 0; i < tmpArray.length; ++i) {
                //scan the role list in RoleManager
                if (roleMgr.checkRoleActionWithKey(tmpArray[i], _contractAddr, _funcSha3) == 1) {
                    return 1;
                }
            }
        } else {*/
            for (uint i = 0; i < userMap[userAddr].roleIdList.length; ++i) {
                //scan the role list in RoleManager
                if (roleMgr.checkRoleActionWithKey(userMap[userAddr].roleIdList[i], _contractAddr, _funcSha3) == 1) {
                    return 1;
                }
            }
        //}

        log1("no privileges, check roles", "UserManager");
        return 0;
    }

    // added by liaoyan. 2016-12-20
    function userExists(address _userAddr) constant public returns(uint _ret) {
        if (_userAddr == owner) {
            return 1;
        }
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID) {
            return 0;
        } else {
            return 1;
        }
    }

    function splitRolesByModuleName(address _userAddr) private returns (uint) {
        delete tmpModuleRoles;
        IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
        for (uint i=0; i < userMap[_userAddr].roleIdList.length; ++i) {
            uint _moduleName = roleMgr.getRoleModuleName(userMap[_userAddr].roleIdList[i]);
            string memory moduleName = _moduleName.recoveryToString();
            uint _moduleVersion = roleMgr.getRoleModuleVersion(userMap[_userAddr].roleIdList[i]);
            string memory moduleVersion = _moduleVersion.recoveryToString();
            if (bytes(moduleName).length == 0 || bytes(moduleVersion).length == 0) {
                continue;
            }
            string memory _roleId = "|";
            _roleId = _roleId.concat(userMap[_userAddr].roleIdList[i], "|");
            uint j;
            for (j=0; j<tmpModuleRoles.length; ++j) {
                if (moduleName.equals(tmpModuleRoles[j].moduleName) && moduleVersion.equals(tmpModuleRoles[j].moduleVersion)) {
                    break;
                }
            }
            if (j >= tmpModuleRoles.length) {
                tmpModuleRole.moduleName = moduleName;
                tmpModuleRole.moduleVersion = moduleVersion;
                tmpModuleRole.roleIds = _roleId;
                tmpModuleRoles.push(tmpModuleRole);
            }
            else {
                if (tmpModuleRoles[j].roleIds.indexOf(_roleId) == -1) {
                    tmpModuleRoles[j].roleIds = tmpModuleRoles[j].roleIds.concat(userMap[_userAddr].roleIdList[i], "|");
                }
                else {
                    log("ignore exist roleid:",userMap[_userAddr].roleIdList[i]);
                }
            }
        }
       return 0;
    }

    function splitRolesByModuleName() private returns (uint) {
        delete tmpModuleRoles;
        IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
        for (uint i=0; i<m_User.roleIdList.length; ++i) {
            uint _moduleName = roleMgr.getRoleModuleName(m_User.roleIdList[i]);
            string memory moduleName = _moduleName.recoveryToString();
            uint _moduleVersion = roleMgr.getRoleModuleVersion(m_User.roleIdList[i]);
            string memory moduleVersion = _moduleVersion.recoveryToString();
            if (bytes(moduleName).length == 0 || bytes(moduleVersion).length == 0) {
                continue;
            }
            string memory _roleId = "|";
            _roleId = _roleId.concat(m_User.roleIdList[i], "|");
            uint j;
            for (j=0; j<tmpModuleRoles.length; ++j) {
                if (moduleName.equals(tmpModuleRoles[j].moduleName) && moduleVersion.equals(tmpModuleRoles[j].moduleVersion)) {
                    break;
                }
            }
            if (j >= tmpModuleRoles.length) {
                tmpModuleRole.moduleName = moduleName;
                tmpModuleRole.moduleVersion = moduleVersion;
                tmpModuleRole.roleIds = _roleId;
                tmpModuleRoles.push(tmpModuleRole);
            }
            else {
                if (tmpModuleRoles[j].roleIds.indexOf(_roleId) == -1) {
                    tmpModuleRoles[j].roleIds = tmpModuleRoles[j].roleIds.concat(m_User.roleIdList[i], "|");
                }
                else {
                    log("ignore exist roleid:",m_User.roleIdList[i]);
                }
            }
        }
        /*log("tmpModuleRoles.length=",tmpModuleRoles.length.toString());
       for (i=0; i<tmpModuleRoles.length; ++i) {
        log(i.toString(), tmpModuleRoles[i].moduleName, tmpModuleRoles[i].roleIds);
       }*/
       return 0;
    }

    function updateAuthorizeUserRole(string _moduleName, string _moduleVersion, address _userAddr, string _roleIds) constant private returns (uint _ret) {//roleIds:|0001|0002|0003|
            log("UserManager.sol", "updateAuthorizeUserRole");
            log(_moduleName,_moduleVersion);
            log(uint(_userAddr).toAddrString(),_roleIds);
            address __userAddr = getOwnerAddrByAddr(_userAddr);
            log("_userAddr->__userAddr", uint(_userAddr).toAddrString(), uint(__userAddr).toAddrString());
            string memory key = uint(__userAddr).toAddrString().concat(_moduleName,_moduleVersion);
            _ret = writedb("authorizeUserRole|update", key, _roleIds);
            if (0 != _ret)
                log("UserManager.sol", "updateAuthorizeUserRole failed.");
            else
                log("UserManager.sol", "updateAuthorizeUserRole success.");
        }

    /**
    *@desc 添加用户
    *@param _userJson 用户json对象
    *@ret   true: 成功, false: 失败
    */
    function insert(string _userJson) public returns(uint) {
        //_ret = false;
        log("insert", "UserManager");
        //log(_userJson);
        if (m_User.jsonParse(_userJson) == false) {
            log("insert bad json", "UserManager");
            errno = 15200 + uint(UserError.BAD_PARAMETER);
            Notify(errno, "insert bad json");
            return errno;
        }
        if (m_User.name.equals("")) {
            log("user name is invalid", "UserManager");
            errno = 15200 + uint(UserError.NAME_EMPTY);
            Notify(errno, "user name is invalid");
           return errno;
        }

        m_User.createTime = now*1000;
        m_User.updateTime = now*1000;
        m_User.loginTime = 0;
        m_User.status = uint(1);
        //m_User.creator = msg.sender;
        // check if user already exists
        m_User.state = LibUser.UserState.USER_VALID;
        if (userMap[m_User.userAddr].state != LibUser.UserState.USER_INVALID) {
            log("address aready exists", "UserManager");
            errno = 15200 + uint(UserError.ADDRESS_ALREADY_EXISTS);
            Notify(errno, "address aready exists");
            return errno;
        }
        for (uint i = 0; i < addrList.length; i++) {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID) {
                if (m_User.account.equals(userMap[addrList[i]].account)) {
                    log("account aready exists", "UserManager");
                    errno = 15200 + uint(UserError.ACCOUNT_ALREDY_EXISTS);
                    Notify(errno, "account aready exists");
                    return errno;
                }
            }
        }
        if (_userJson.keyExists("email") && bytes(m_User.email).length > 0) {
            for ( i = 0; i < addrList.length; i++) {
                if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID) {
                    if (m_User.email.equals(userMap[addrList[i]].email)) {
                        log("email aready exists", "UserManager");
                        errno = 15200 + uint(UserError.EMAIL_EXISTS);
                        Notify(errno, "email aready exists");
                        return errno;
                    }
                }
            }
        }
        if (_userJson.keyExists("mobile") && bytes(m_User.mobile).length > 0) {
            for ( i = 0; i < addrList.length; i++) {
                if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID) {
                    if (m_User.mobile.equals(userMap[addrList[i]].mobile)) {
                        log("mobile aready exists", "UserManager");
                        errno = 15200 + uint(UserError.MOBILE_EXISTS);
                        Notify(errno, "mobile aready exists");
                        return errno;
                    }
                }
            }
        }
        if (tx.origin != owner) {
            log("sender is no owner:", tx.origin);
            /* if (__checkWritePermission(tx.origin, m_User.departmentId) == 0) {
                log("operator no permission");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return errno;
            } */
        }

        // check the role list
        IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
        if(m_User.roleIdList.length == 0){
            m_User.roleIdList.push("role100004");
        }
        for (i = 0; i < m_User.roleIdList.length; ++i) {
            if (roleMgr.roleExists(m_User.roleIdList[i]) == 0) {
                log("insert with invalid role", "UserManager");
                errno = 15200 + uint(UserError.ROLE_ID_INVALID);
                Notify(errno, "insert a user with invalid role");
                return errno;
            }
        }
        // check department in DepartmentManager
        IDepartmentManager departMgr = IDepartmentManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0"));
        if (0 == departMgr.departmentExists(m_User.departmentId)) {
            log("department not exists", "UserManager");
            errno = 15200 + uint(UserError.DEPT_NOT_EXISTS);
            Notify(errno, "department not exists");
            return errno;
        }

        // insert a user 
        m_User.ownerAddr = m_User.userAddr;
        userMap[m_User.userAddr] = m_User;
        addrList.push(m_User.userAddr);

        splitRolesByModuleName();
        for (i=0; i<tmpModuleRoles.length; ++i) {
            updateAuthorizeUserRole(tmpModuleRoles[i].moduleName, tmpModuleRoles[i].moduleVersion, m_User.userAddr, tmpModuleRoles[i].roleIds);
        }
        updateUserInfoDb(m_User.userAddr, m_User.accountStatus, m_User.deleteStatus, uint(m_User.state));

        errno = uint(UserError.NO_ERROR);
        revision++;

        log("insert a user succcess", "UserManager");
        Notify(errno, "insert a user succcess");
        m_User.reset();
        return errno;
    }

    /**
    *@desc 更新用户信息 （json中空值或0，对应为不设置，非空或非0为设置该字段）
    *@param _userJson 用户json对象
    *@ret   true: 成功, false: 失败
    */
    function update(string _userJson) public returns(uint) {
        log("update user", "UserManager");
        //log(_userJson);
        
        if (m_User.jsonParse(_userJson) == false) {
            log("insert bad json", "UserManager");
            errno = 15200 + uint(UserError.BAD_PARAMETER);
            Notify(errno, "insert bad json");
            return errno;
        }

        // check if user exists
        m_User.state = LibUser.UserState.USER_VALID;
        if (userMap[m_User.userAddr].state == LibUser.UserState.USER_INVALID) {
            log("user not exists", "UserManager");
            errno = 15200 + uint(UserError.USER_NOT_EXISTS);
            Notify(errno, "user not exists");
            return errno;
        }

        if (tx.origin != owner) {
            log("msg.sender is not owner", "UserManager");
            /* if (__checkWritePermission(tx.origin, userMap[m_User.userAddr].departmentId) == 0) {
                log("operator no permission", "UserManager");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return errno;
            } */
        }

        for (uint i = 0; i < addrList.length; i++) {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID) {
                if(userMap[addrList[i]].userAddr == m_User.userAddr){
                    continue;
                }
                if (m_User.email.equals(userMap[addrList[i]].email)) {
                    log("email aready exists", "UserManager");
                    errno = 15200 + uint(UserError.EMAIL_EXISTS);
                    Notify(errno, "email aready exists");
                    return errno;
                }
            }
        }

        for ( i = 0; i < addrList.length; i++) {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID) {
                if(userMap[addrList[i]].userAddr == m_User.userAddr){
                    continue;
                }
                if (m_User.mobile.equals(userMap[addrList[i]].mobile)) {
                    log("mobile aready exists", "UserManager");
                    errno = 15200 + uint(UserError.MOBILE_EXISTS);
                    Notify(errno, "mobile aready exists");
                    return errno;
                }
            }
        }

        // check department
        if (!m_User.departmentId.equals("") && !m_User.departmentId.equals(userMap[m_User.userAddr].departmentId)) {
            log("not permit update department", "UserManager");
            errno = 15200 + uint(UserError.DEPT_CANNOT_UPDATE);
            Notify(errno, "can not update department");
            return errno;
        }

        // check account
        if (!m_User.account.equals("") && !m_User.account.equals(userMap[m_User.userAddr].account)) {
            log("account can not update", "UserManager");
            errno = 15200 + uint(UserError.ACCOUNT_CANNOT_UPDATE);
            Notify(errno, "account can not update");
            return errno;
        }

        // check the role list
        if (m_User.roleIdList.length > 0) {
            IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
            //BaseInterface departMgr = BaseInterface(rm.getContractAddress("DepartmentManager", "0.0.1.0"));
            for ( i = 0; i < m_User.roleIdList.length; ++i) {
                // check role in RoleManager
                if (roleMgr.roleExists(m_User.roleIdList[i]) == 0) {
                    log("input invalid role", "UserManager");
                    errno = 15200 + uint(UserError.ROLE_ID_INVALID);
                    Notify(errno, "input invalid role");
                    return errno;
                }
                // check department role
                /*if (departMgr.checkDepartmentRole(userMap[m_User.userAddr].departmentId, m_User.roleIdList[i]) == 0) {
                    log("bad role for department", "UserManager");
                    errno = 15200 + uint(UserError.ROLE_ID_EXCEED_DEPT);
                    return;
                }*/
            }
        }

        if (tx.origin != owner && tx.origin != userMap[m_User.userAddr].ownerAddr) {
            log("exec updata sender is must be owner or belong to self.", tx.origin);
            log("operator no permission");
            errno = 15200 + uint(UserError.NO_PERMISSION);
            Notify(errno, "operator no permission");
            return errno;
        }

        m_User.updateTime = now*1000;
        if (_userJson.keyExists("name")) {
            userMap[m_User.userAddr].name = m_User.name;
        }

        /* if (_userJson.keyExists("age")) {
            userMap[m_User.userAddr].age = m_User.age;
        } */

        /* if (_userJson.keyExists("sex")) {
            userMap[m_User.userAddr].sex = m_User.sex;
        } */

        /* if (_userJson.keyExists("birthday")) {
            userMap[m_User.userAddr].birthday = m_User.birthday;
        } */

        if (_userJson.keyExists("account")) {
            userMap[m_User.userAddr].account = m_User.account;
        }

        if (_userJson.keyExists("email")) {
            userMap[m_User.userAddr].email = m_User.email;
        }

        if (_userJson.keyExists("mobile")) {
            userMap[m_User.userAddr].mobile = m_User.mobile;
        }

        if (_userJson.keyExists("account")) {
            userMap[m_User.userAddr].account = m_User.account;
        }

        if (_userJson.keyExists("accountStatus")) {
            userMap[m_User.userAddr].accountStatus = m_User.accountStatus;
        }

        if (_userJson.keyExists("passwordStatus")) {
            userMap[m_User.userAddr].passwordStatus = m_User.passwordStatus;
        }

        if (_userJson.keyExists("deleteStatus")) {
            userMap[m_User.userAddr].deleteStatus = m_User.deleteStatus;
        }

        if (_userJson.keyExists("remark")) {
            userMap[m_User.userAddr].remark = m_User.remark;
        }

        if (_userJson.keyExists("icon")) {
            userMap[m_User.userAddr].icon = m_User.icon;
        }

        /* if (_userJson.keyExists("tokenSeed")) {
            userMap[m_User.userAddr].tokenSeed = m_User.tokenSeed;
        } */

        if (_userJson.keyExists("uuid")) {
            userMap[m_User.userAddr].uuid = m_User.uuid;
        }

        if (_userJson.keyExists("publicKey")) {
            userMap[m_User.userAddr].publicKey = m_User.publicKey;
        }

        if (_userJson.keyExists("cipherGroupKey")) {
            userMap[m_User.userAddr].cipherGroupKey = m_User.cipherGroupKey;
        }

        if (_userJson.keyExists("status")) {
            userMap[m_User.userAddr].status = m_User.status;
        }


        if (_userJson.keyExists("roleIdList")) {
            delete userMap[m_User.userAddr].roleIdList;
            for (uint index = 0 ; index < m_User.roleIdList.length; ++index) {
                userMap[m_User.userAddr].roleIdList.push(m_User.roleIdList[index]);

                splitRolesByModuleName();
                for (i=0; i<tmpModuleRoles.length; ++i) {
                    updateAuthorizeUserRole(tmpModuleRoles[i].moduleName, tmpModuleRoles[i].moduleVersion, m_User.userAddr, tmpModuleRoles[i].roleIds);
                }
            }
        }
        
        updateUserInfoDb(m_User.userAddr, userMap[m_User.userAddr].accountStatus, userMap[m_User.userAddr].deleteStatus, uint(userMap[m_User.userAddr].state));
        errno = uint(UserError.NO_ERROR);
        revision++;
        log("update a user succcess", "UserManager");
        Notify(errno, "update a user succcess");
        m_User.reset();
        return errno;
    }

    /**
    * 更新用户状态 0 禁用 1 启用
    * @param _userAddr 用户地址
    * @param _status 用户状态 
    * @return true success | false fail
    */
    function updateUserStatus(address _userAddr, uint _status) public returns(bool _ret) {
        _ret = false;
        if(_status != 0 && _status != 1) {
            log("user status invalid", "UserManager");
            errno = 15200 + uint(UserError.BAD_PARAMETER);
            Notify(errno, "user status invalid");
            return;
        }
        // TODO: 此处是否要加入对角色的校验
        // 仅有部门管理员才可更改此操作
        if (msg.sender != owner) {
            log("msg.sender is not owner", "UserManager");
            /* if (__checkWritePermission(msg.sender, userMap[_userAddr].departmentId) == 0) {
                log("operator no permission", "UserManager");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return;
            } */
        }

        userMap[_userAddr].status = _status;
        userMap[_userAddr].updateTime = now * 1000;
        
        _ret = true;
        errno = uint(UserError.NO_ERROR);
        revision++;
        log("update user status OK", "UserManager");
        Notify(errno, "update user status success");
        return;
    }

    /**
    *desc   更新账号状态
    *param  _status  账号状态，0 非法，1 正常
    *@ret   true: 成功, false: 失败
    */
    function updateAccountStatus(address _userAddr, uint _status) public returns(bool _ret) {
        _ret = false;
        
        if (_status > uint(LibUser.AccountState.LOCKED)) {
            log("account status invalid", "UserManager");
            errno = 15200 + uint(UserError.BAD_PARAMETER);
            Notify(errno, "account status invalid");
            return;
        }

        // check if user exists
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID) {
            log("user not exists", "UserManager");
            errno = 15200 + uint(UserError.USER_NOT_EXISTS);
            Notify(errno, "user not exists");
            return;
        }

        if (msg.sender != owner) {
            log("msg.sender is not owner", "UserManager");
            /* if (__checkWritePermission(msg.sender, userMap[_userAddr].departmentId) == 0) {
                log("operator no permission", "UserManager");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return;
            } */
        }

        userMap[_userAddr].accountStatus = uint(_status);
        userMap[_userAddr].updateTime = now*1000;
        updateUserInfoDb(_userAddr, userMap[_userAddr].accountStatus, userMap[_userAddr].deleteStatus, uint(userMap[_userAddr].state));
        
        _ret = true;
        errno = uint(UserError.NO_ERROR);
        revision++;
        log("update user accountStatus OK", "UserManager");
        Notify(errno, "update accountStatus success");
        return;
    }

    /**
    *desc   更新账号状态
    *param  _status  
    *@ret   true: 成功, false: 失败
    */
    function updatePasswordStatus(address _userAddr, uint _status) public returns(bool _ret) {
        _ret = false;
        
        // check if user exists
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID) {
            log("user not exists", "UserManager");
            errno = 15200 + uint(UserError.USER_NOT_EXISTS);
            return;
        }

        if (msg.sender != owner) {
            log("msg.sender is not owner", "UserManager");
            /* if (__checkWritePermission(msg.sender, userMap[_userAddr].departmentId) == 0) {
                log("operator no permission", "UserManager");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return;
            } */
        }

        userMap[_userAddr].passwordStatus = uint(_status);
        userMap[_userAddr].updateTime = now*1000;
        
        _ret = true;
        errno = uint(UserError.NO_ERROR);
        revision++;

        log("update user account status OK", "UserManager");
        Notify(errno, "update user account status OK");
        return;
    }

    /**
    *@desc 添加用户角色
    *@param _userAddr 用户地址， _roleId 用户角色id
    *@ret   true: 成功, false: 失败
    */
    function addUserRole(address _userAddr, string _roleId) returns(uint) {

        log("add user role", "UserManager");
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID) {
            log("user not exists: ", _userAddr);
            errno = 15200 + uint(UserError.USER_NOT_EXISTS);
            Notify(errno,"user not exists...");
            return errno;
        }

        if (tx.origin != owner) {
            log("msg.sender is not owner", "UserManager");
            /* if (__checkWritePermission(tx.origin, userMap[_userAddr].departmentId) == 0) {
                log("operator no permission", "UserManager");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return errno;
            } */
        }

        // if exists do not add again
        for (uint i = 0; i < userMap[_userAddr].roleIdList.length; ++i) {
            if (userMap[_userAddr].roleIdList[i].equals(_roleId)) {
                errno = 15200 + uint(UserError.ROLE_ID_ALREADY_EXISTS);
                log("role already exists", "UserManager");
                Notify(errno,"role already exists");
                return errno;
            }
        }
        
        // check the role
        IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
        if (roleMgr.roleExists(_roleId) == 0) {
            log("add a bad role", "UserManager");
            log("bad roleId : ",_roleId);
            errno = 15200 + uint(UserError.ROLE_ID_INVALID);
            Notify(errno,"add a bad role");
            return errno;
        }

        // check department role
        /*IDepartmentManager departMgr = IDepartmentManager(rm.getContractAddress("DepartmentManager", "0.0.1.0"));
        if (departMgr.checkDepartmentRole(userMap[_userAddr].departmentId, _roleId) == 0) {
            log("role exceed department", "UserManager");
            errno = 15200 + uint(UserError.ROLE_ID_EXCEED_DEPT);
            return errno;
        }*/

        log("add role success", "UserManager");
        userMap[_userAddr].roleIdList.push(_roleId);
        errno = uint(UserError.NO_ERROR);

        // save data to db
        splitRolesByModuleName(_userAddr);
        for (i = 0 ; i < tmpModuleRoles.length; ++i) {
            log("addUserRole->save Data to DB,userAddr:",uint(userMap[_userAddr].userAddr).toAddrString());
            log("addUserRole->save Data to DB,roleId:",tmpModuleRoles[i].roleIds);
            updateAuthorizeUserRole(tmpModuleRoles[i].moduleName, tmpModuleRoles[i].moduleVersion, userMap[_userAddr].userAddr, tmpModuleRoles[i].roleIds);
        }

        revision++;
        Notify(errno,"add a userRole success");
        return errno;
    }

    /**
    *@desc 按照用户地址删除用户
    *@param _userAddr 用户
    *@ret   true: 成功, false: 失败
    */
    function deleteByAddress(address _userAddr) public {
        log("deleteByAddress", "UserManager");

        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID) {
            log("user not exists: ", _userAddr);
            errno = 15200 + uint(UserError.USER_NOT_EXISTS);
            Notify(errno, "user not exists");
            return;
        } 

        if (msg.sender != owner) {
            log("msg.sender is not owner", "UserManager");
            /* if (__checkWritePermission(msg.sender, userMap[_userAddr].departmentId) == 0) {
                log("operator no permission");
                errno = 15200 + uint(UserError.NO_PERMISSION);
                Notify(errno, "operator no permission");
                return;
            } */
        }

        // erase department admin address, if this address is admin
        IDepartmentManager departMgr = IDepartmentManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0"));
        departMgr.eraseAdminByAddress(_userAddr);

        // first update id list
        delete tempList;
        for (uint i = 0; i < addrList.length; i++)  {
            if(_userAddr == userMap[addrList[i]].userAddr) {
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

        userMap[_userAddr].state = LibUser.UserState.USER_INVALID;
        updateUserInfoDb(_userAddr, userMap[_userAddr].accountStatus, userMap[_userAddr].deleteStatus, uint(userMap[_userAddr].state));
        log("delete user success", "UserManager");
        errno = uint(UserError.NO_ERROR);
        revision++;
        Notify(errno, "delete user success");
    }

    /**
    *@desc 用户登录
    *@param _account 用户账号
    *@ret   返回json字符串
    */
    function login(string _account) public returns(string _json) {
        _json = "{\"ret\":-1,\"data\":false}";

        for (uint i = 0; i < addrList.length; i++) {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID 
                && _account.equals(userMap[addrList[i]].account) 
                && userMap[addrList[i]].accountStatus == uint(LibUser.AccountState.VALID) 
                && userMap[addrList[i]].deleteStatus == 0) {
                userMap[addrList[i]].state = LibUser.UserState.USER_LOGIN;
                userMap[addrList[i]].loginTime = now;

                uint len = itemStackPush(userMap[addrList[i]].toJson(), getUserCount());
                _json = LibStack.popex(len);
                errno = uint(UserError.NO_ERROR);
                Notify(errno, "login user success");
                log("user login success, account :", _account);
                return;
            }
        }
        
        errno = 15200 + uint(UserError.USER_LOGIN_FAILED);
        Notify(errno, "login user failed");
        return;
    }

    /**
    *@desc 展示所有用户
    *@ret  返回json
    */
    function listAll() constant public returns(string _userListJson) {
        uint len = 0;
        uint counter = 0;
        len = LibStack.push("");
        for (uint index = 0; index < addrList.length; index++) {
            if (userMap[addrList[index]].state != LibUser.UserState.USER_INVALID) {
                if (counter > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(userMap[addrList[index]].toJson());
                counter++;
            }
        }
        len = itemStackPush(LibStack.popex(len), counter);
        _userListJson = LibStack.popex(len);
    }

    /**
    *@desc 用户登录
    *@param _account 用户账号
    *@ret   返回json字符串
    */
    function getUserCount() constant public returns(uint _count) {
        for(uint i = 0 ; i < addrList.length ; ++i) {
            if (userMap[addrList[i]].state != LibUser.UserState.USER_INVALID) {
                _count++;
            }
        }
        //_count = addrList.length;
    }

    /**
    *@desc 获取用户id的int值
    *@param _userAddr 用户地址，index 索引
    *@ret   id的int序列值
    */
    function getUserRoleId(address _userAddr, uint _index) constant returns (uint _ret) {
        if (userMap[_userAddr].state == LibUser.UserState.USER_INVALID) {
            return 0;
        }

        if (_index >= 0 && _index < userMap[_userAddr].roleIdList.length) {
            return userMap[_userAddr].roleIdList[_index].storageToUint();
        }

        return 0;
    }

    /**
    * @dev Get data revision
    * @return revision id
    */
    function getRevision() constant public returns(uint _revision) {
        return revision;
    }

    /**
    *@desc 获取指定部门下用户数量（不含子部门）
    *@param _departmentId 指定的部门ID
    *@ret _count 用户数量
    */
    function getUserCountByDepartmentId(string _departmentId) constant public returns(uint _count) {
        _count = 0;
        for (uint i = 0; i < addrList.length; i++) {
            if(userMap[addrList[i]].state != LibUser.UserState.USER_INVALID 
                && userMap[addrList[i]].departmentId.equals(_departmentId)) {
                _count++;
            }
        }
    }

    /**
    *@desc 获取拥有指定权限的用户数量
    *@param _actionId 指定的权限ID
    *@ret _count 用户数量
    */
    function getUserCountByActionId(string _actionId) constant public returns(uint _count) {
        _count = 0;
        
        //获取actionId对应的所有的roldId
        delete tmpArray;
        __getRoleIdListByActionId(_actionId, tmpArray);

        //检查每个User是否拥有至少1个上述RoleId
        for (uint i=0; i<addrList.length; ++i) {
            if (userMap[addrList[i]].state == LibUser.UserState.USER_INVALID) {
                continue;
            }

            for (uint j=0; j<tmpArray.length; ++j) {
                if (tmpArray[j].inArray(userMap[addrList[i]].roleIdList)) {
                    _count++;
                    break;
                }
            }
        }
    }

    /**
    *@desc 检查role是否被使用
    *@param _roleId 角色Id
    *@ret   1 存在，0 不存在 
    */
    function roleUsed(string _roleId) constant public returns (uint _used) {
        _used = 0;
        for (uint i = 0; i < addrList.length; ++i) {
            if (userMap[addrList[i]].state == LibUser.UserState.USER_INVALID) {
                continue;
            }

            for (uint j = 0; j < userMap[addrList[i]].roleIdList.length; ++j) {
                if (userMap[addrList[i]].roleIdList[j].equals(_roleId)) {
                    _used = 1;
                    return;
                }
            }
        }
    }

    /**
    *@desc 根据角色ID获取用户数量
    *@param _roleId 角色Id
    *@ret   用户数量
    */
    function getUserCountMappingByRoleIds(string _roleIds) constant public returns(string _json) {
        // roleIds = roleId0001,roleId0002
        log('getCountMapping->',_roleIds);
        delete tmpArray;
        _roleIds.split(",",tmpArray);
        uint len = 0;
        len = LibStack.push("");
        // {"ret":0,"data":{"items":[{"rooleId":1},{}]}}
        for(uint j = 0 ; j < tmpArray.length; ++j){
            uint roleIdCount = 0;
            if(j > 0) {
                len = LibStack.append(",");
            }
            len = LibStack.append("{");
            for (uint i = 0; i < addrList.length; i++) {
                if(userMap[addrList[i]].state == LibUser.UserState.USER_INVALID) {
                    continue;
                }
                if(tmpArray[j].inArray(userMap[addrList[i]].roleIdList)){ 
                    roleIdCount++;
                }
            }
            len = LibStack.appendKeyValue(tmpArray[j], roleIdCount);
            len = LibStack.append("}");
        }
        len = itemStackPush(LibStack.popex(len), uint(tmpArray.length));
        _json = LibStack.popex(len);
    }

    function __getDepartmentAdmin(string departmentId) constant internal returns(address admin) {
        // Get DepartmentManager address
        address departmentManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0");
        if (departmentManagerAddr == 0) {
            return 0;
        }
 
        // Get DepartmentManager token
        IDepartmentManager departmentManager = IDepartmentManager(departmentManagerAddr);

        return address(departmentManager.getAdmin(departmentId));
    }

    /*function __getDepartmentRoleIdList(string _departmentId, string[] storage roleIdList) constant internal {
        // Get DepartmentManager address
        address departmentManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0");
        if (departmentManagerAddr == 0) {
            return;
        }

        // Get DepartmentManager token
        IDepartmentManager departmentManager = IDepartmentManager(departmentManagerAddr);
            
        uint i = 0;
        while (true) {
            uint ret = departmentManager.getDepartmentRoleId(_departmentId, i);
            if (ret == 0) {
                return;
            } else {
                roleIdList.push(ret.recoveryToString());
            }
            i++;
        }
    }*/

    function __checkWritePermission(address _addr, string _departmentId) constant internal returns (uint _ret) {
        // Get DepartmentManager address
        address departmentManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0");
        if (departmentManagerAddr == 0) {
            return 0;
        }

        // Get DepartmentManager token
        IDepartmentManager departmentManager = IDepartmentManager(departmentManagerAddr);
        return departmentManager.checkWritePermission(_addr, _departmentId);
    }

    function __getDepartmentIdTree(string _departmentId, string[] storage _departmentIdTree) constant internal {
        // Get DepartmentManager address
        address departmentManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0");
        if (departmentManagerAddr == 0) {
            return;
        }

        // Get DepartmentManager token
        IDepartmentManager departmentManager = IDepartmentManager(departmentManagerAddr);

        _departmentIdTree.push(_departmentId);

        for (uint i=_departmentIdTree.length-1; i<_departmentIdTree.length; ++i) {
            uint index = 0;
            while (true) {
                uint id = departmentManager.getChildIdByIndex(_departmentIdTree[i], index);
                if (id == 0) {
                    break;
                }

                _departmentIdTree.push(id.recoveryToString());
                index++;
            }
        }
    }

    function __getRoleIdListByActionId(string _actionId, string[] storage roleIdList) constant internal {
        // Get DepartmentManager address
        address roleManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0");
        if (roleManagerAddr == 0) {
            return;
        }

        // Get DepartmentManager token
        IRoleManager roleManager = IRoleManager(roleManagerAddr);
            
        uint i = 0;
        while (true) {
            uint ret = roleManager.getRoleIdByActionIdAndIndex(_actionId, i);
            if (ret == 0) {
                return;
            } else {
                roleIdList.push(ret.recoveryToString());
            }
            i++;
        }
    }

    // 根据update获取当前用户的userAddr,仅针对使用update进行调用场景，eth底层直接调用的函数
    // since client v1.0.0
    function getUserAddrByAddr(address _userAddr) constant public returns (address _address) {
        // 如果用户地址等于发送者，则不获取新的
        _address = address(0);
        for (uint i = 0; i < addrList.length; ++i) {
            if (userMap[addrList[i]].state == LibUser.UserState.USER_INVALID) {
                continue;
            }
            if(userMap[addrList[i]].ownerAddr == _userAddr){
                _address = userMap[addrList[i]].userAddr;
                return;
            }
        }
        return _userAddr;
    }

    // since client v1.0.0
    function getOwnerAddrByAddr(address _userAddr) constant public returns (address _address) {
        return userMap[_userAddr].ownerAddr;
    }

    function checkEmailUniqueness(string _email,string _mobile) public constant returns(uint){
        for (uint index = 0 ; index < addrList.length; index++){
            if (userMap[addrList[index]].state == LibUser.UserState.USER_INVALID) {
                continue;
            }
            if (userMap[addrList[index]].mobile.equals(_mobile)){
                return 2;
            }
            if(userMap[addrList[index]].email.equals(_email)){
                return 1;
            }else {
                return 0;
            }
        }
        return 0;
    }



    //检查重复
    function isRepetitive(string _mobile, string _email, address _userAddr, string _uuid, string _publicKey, string _account) public constant returns (uint) {
        for (uint index = 0; index < addrList.length; index++) {
            if (userMap[addrList[index]].state == LibUser.UserState.USER_INVALID) {
                continue;
            }
            if (userMap[addrList[index]].account.equals(_account)) return 5;
            if (userMap[addrList[index]].mobile.equals(_mobile)) return 6;
            if (userMap[addrList[index]].email.equals(_email)) return 7;
            if (userMap[addrList[index]].userAddr == _userAddr ||
            userMap[addrList[index]].uuid.equals(_uuid) ||
            userMap[addrList[index]].publicKey.equals(_publicKey)) {
                return 4;
            }
        }
        return 0;
    }

    //items入栈
    function itemStackPush(string _items, uint _total) constant private returns (uint len){
        len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.append(",");
        len = LibStack.append("\"data\":{");
        len = LibStack.appendKeyValue("total", _total);
        len = LibStack.append(",");
        len = LibStack.append("\"items\":[");
        len = LibStack.append(_items);
        len = LibStack.append("]}}");
        return len;
    }
}