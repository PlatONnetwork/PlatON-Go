pragma solidity ^0.4.12;
/**
* @file RoleManager.sol
* @author liaoyan
* @time 2016-11-29
* @desc the defination of RoleManager contract
*/

import "./library/LibRole.sol";
import "./sysbase/OwnerNamed.sol";
import "./interfaces/IRoleManager.sol";
import "./interfaces/IActionManager.sol";
import "./interfaces/IUserManager.sol";
import "./interfaces/IDepartmentManager.sol";

contract RoleManager is OwnerNamed, IRoleManager {
    using LibRole for *;
    using LibString for *;
    using LibInt for *;
    
    LibRole.Role[] roles;
    LibRole.Role[] tmpRoles;
    LibRole.Role tmpRole;
    string[] tmpRoleIdList;
    string[] tmpActionIdList;

    enum RoleError {
        NO_ERROR,
        BAD_PARAMETER,
        NAME_EMPTY,
        ID_NOT_EXISTS,
        ID_CONFLICTED,
        NO_PERMISSION,
        USER_INVALID,
        ACTION_ID_INVALID,
        ACTION_ID_EXCEEDED_DEPT,
        ACTION_ID_EXCEEDED_USER,
        ROLE_USED
    }

    event Notify(uint _errno, string _info);

    /**
    * @dev Contruction
    */
    function RoleManager() {
        register("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0");
    }

    /**
     * @dev add new actionId to role
     * @param _roleId role id
     * @param _actionId action id
     * @return no return
     */
    function addActionToRole(string _roleId, string _actionId) public returns(uint) {
        log("into addActionToRole.", "RoleManager");
        uint prefix = 95270;
        if (tx.origin != owner) {
            log("not root user, can not invoke addActionToRole function.", "RoleManager");
            errno = prefix + uint(RoleError.NO_PERMISSION);
            Notify(errno, "only root user can addActionToRole function.");
            return errno;
        }

        // Find the index of current rold
        uint index = uint(-1);
        for (uint i = 0; i < roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                index = i;
                break;
            }
        }

        if (index == uint(-1)) {
            log("addActionToRole->roid id does not exists", "RoleManager");
            errno = prefix + uint(RoleError.ID_NOT_EXISTS);
            Notify(errno, "roid id does not exists");
            return errno;
        }

        // Check if every new action id is in ActionManager
        if (__actionExists(_actionId) == 0) {
            log("addActionToRole -> action not in ActionManager", "RoleManager"); 
            errno = prefix + uint(RoleError.ACTION_ID_INVALID);
            Notify(errno, "action not exists in ActionManager");
            return errno;
        }
        if(_actionId.inArray(roles[index].actionIdList)){
            log("addActionToRole -> actionId already exists for role.", "RoleManager"); 
            errno = prefix + uint(RoleError.ACTION_ID_INVALID);
            Notify(errno, "actionId already exists for role.");
            return errno;
        }
        
        roles[index].actionIdList.push(_actionId);
        roles[index].updTime = uint(now) * 1000;

        errno = uint(RoleError.NO_ERROR);

        log("addActionToRole succ", "RoleManager");
        Notify(errno, "addActionToRole succ");
        return errno;
    }

    /**
    * get role list by module id
    * @param _moduleId the moduleId
    * @return _json role info list
    */
    function getRoleListByModuleId(string _moduleId) constant public returns(string _json) {
        _json = listByUK(3,_moduleId);
    }

    /**
    * get role list by contract id
    * @param _contract the contract id
    * @return _json role info list
    */
    function getRoleListByContractId(string _contract) constant public returns(string _json) {
        _json = listByUK(4,_contract);
    }
    
    function getRoleListByModuleName(string _moduleName, string _moduleVersion) constant public returns(string _json) {
        _json = listBy(3,_moduleName,_moduleVersion);
    }

    /**
    * @dev insert an object
    * @param _json The json string described the object
    * @return No return
    */
    function insert(string _json) public returns(uint) {
        log("insert a role", "RoleManager");

        if (!tmpRole.fromJson(_json)) {
            log("json invalid", "RoleManager");
            errno = 15300 + uint(RoleError.BAD_PARAMETER);
            Notify(errno, "json invalid");
            return errno;
        }

        /*
        * A new role must matches following conditions:
        * 1. role id DOSE NOT exists,
        */

        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(tmpRole.id)) {
                log("id conflicted", "RoleManager");
                errno = 15300 + uint(RoleError.ID_CONFLICTED);
                Notify(errno, "id conflicted");
                return errno;
            }
        }

        if (bytes(tmpRole.name).length == 0) {
            log("name empty", "RoleManager");
            errno = 15300 + uint(RoleError.NAME_EMPTY);
            Notify(errno, "name empty");
            return errno;
        }

        if (tx.origin == owner) {
            log("msg.sender is owner:", tx.origin); //debug
            // Check if every new action id is in ActionManager
            for (i=0; i<tmpRole.actionIdList.length; ++i) {
                if (__actionExists(tmpRole.actionIdList[i]) == 0) {
                    log("action not exists in ActionManager: ", tmpRole.actionIdList[i]); //debug
                    errno = 15300 + uint(RoleError.ACTION_ID_INVALID);
                    Notify(errno, "action not exists in ActionManager");
                    return errno;
                }
            }
        } else {
            log("msg.sender is user:", tx.origin); //debug

            // Get user's department id from UserManager
            /*string memory departmentId = __getUserDepartmentId(tx.origin);
            if (bytes(departmentId).length == 0) {
                log("user department not exists"); //debug
                errno = 15300 + uint(RoleError.USER_INVALID);
                Notify(errno, "user department does not exists");
                return errno;
            }

            delete tmpRoleIdList;
            delete tmpActionIdList;

            // Get department's role id list from DepartmentManager
            __getDepartmentRoleIdList(departmentId, tmpRoleIdList);

            // Get department's action id list
            getActionIdListByRoleIdList(tmpRoleIdList, tmpActionIdList);
            //log("department action id list: ", tmpActionIdList.toKeyValue("ids")); //debug

            // Check if every new action id is in the department's action id list
            for (i=0; i<tmpRole.actionIdList.length; ++i) {
                if (!tmpRole.actionIdList[i].inArray(tmpActionIdList)) {
                    log("action not in dept actionlist", "RoleManager"); //debug
                    errno = 15300 + uint(RoleError.ACTION_ID_EXCEEDED_DEPT);
                    Notify(errno, "new action not in department action id list");
                    return errno;
                }
            }*/

            delete tmpRoleIdList;
            delete tmpActionIdList;

            // Get use's role id list from UserManager
            __getUserRoleIdList(msg.sender, tmpRoleIdList);
            //log("user role id list: ", tmpRoleIdList.toKeyValue("ids")); //debug
            
            // Get user's action id list
            getActionIdListByRoleIdList(tmpRoleIdList, tmpActionIdList);
            //log("user action id list: ", tmpActionIdList.toKeyValue("ids")); //debug

            // Check if every new action id is in the user's action id list
            for (i=0; i<tmpRole.actionIdList.length; ++i) {
                if (!tmpRole.actionIdList[i].inArray(tmpActionIdList)) {
                    log("action not in user actionlist", "RoleManager"); //debug
                    log("valid action is:",tmpRole.actionIdList[i]);
                    errno = 15300 + uint(RoleError.ACTION_ID_EXCEEDED_USER);
                    Notify(errno, "new action not in user action id list");
                    return errno;
                }
            }
        }

        tmpRole.creator = tx.origin;
        tmpRole.creTime = uint(now)*1000;
        tmpRole.updTime = uint(now)*1000;

        roles.push(tmpRole);

        errno = uint(RoleError.NO_ERROR);

        log("insert succ", "RoleManager"); //debug
        Notify(errno, "insert succ");
        return errno;
    }

    function update(string _json) public {
        log("update role", "RoleManager");
       
        if (tx.origin != owner) {
            log("not root user, no update right", "RoleManager");
            errno = 15300 + uint(RoleError.NO_PERMISSION);
            Notify(errno, "only root user can update role");
            return;
        }

        if (!tmpRole.fromJson(_json)) {
            log("json invalid", "RoleManager");
            errno = 15300 + uint(RoleError.BAD_PARAMETER);
            Notify(errno, "json invalid");
            return;
        }

        log("after fromJson", "RoleManager"); //debug

        // Find the index of current rold
        uint index = uint(-1);
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(tmpRole.id)) {
                index = i;
            }
        }

        if (index == uint(-1)) {
            log("roid id does not exists", "RoleManager");
            errno = 15300 + uint(RoleError.ID_NOT_EXISTS);
            Notify(errno, "roid id does not exists");
            return;
        }

        // Check if every new action id is in ActionManager
        for (i=0; i<tmpRole.actionIdList.length; ++i) {
            if (__actionExists(tmpRole.actionIdList[i]) == 0) {
                log("action not in ActionManager", "RoleManager"); //debug
                errno = 15300 + uint(RoleError.ACTION_ID_INVALID);
                Notify(errno, "action not exists in ActionManager");
                return;
            }
        }

        if (_json.keyExists("name")) {
            roles[index].name = tmpRole.name;
        }
        if (_json.keyExists("status")) {
            roles[index].status = tmpRole.status;
        }
        if (_json.keyExists("description")) {
            roles[index].description = tmpRole.description;
        }
        if (_json.keyExists("actionIdList")) {
            roles[index].actionIdList = tmpRole.actionIdList;
        }
        roles[index].updTime = uint(now)*1000;

        errno = uint(RoleError.NO_ERROR);

        log("update role succ", "RoleManager"); //debug
        Notify(errno, "update succ");
    }

    /**
    * @dev List the all objects
    * @return No return
    */
    function listAll() constant public returns (string _json) {
        _json = listByUK(0, "");
    }

    /**
    * @dev Find object by id
    * @param _id Object id
    * @return _json Objects in json string
    */
    function findById(string _id) constant public returns(string _json) {
        _json = listByUK(1, _id);
    }

    /**
    * @dev Find object by name
    * @param _name Object name
    * @return _json Objects in json string
    */
    function findByName(string _name) constant public returns(string _json) {
        _json = listByUK(2, _name);
    }

    /**
    * @dev check if a role id contains a action id
    * @param _roleId The role id
    * @param _actionId The action id
    * @return _ret If contains return 1, else return 0
    */
    function checkRoleAction(string _roleId, string _actionId) constant public returns (uint _ret) {
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                for (uint j=0; j<roles[i].actionIdList.length; ++j) {
                    if (roles[i].actionIdList[j].equals(_actionId)) {
                        return 1;
                    }
                }
                return 0;
            }
        }

        return 0;
    }

    /**
    * @dev check if a role id exists in the roles
    * @param _roleId The role id
    * @return _ret If contains return 1, else return 0
    */
    function roleExists(string _roleId) constant public returns (uint _ret) {
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                return 1;
            }
        }

        return 0;
    }

    /**
    * @dev check if a role id exists in the roles
    * @param _roleId The role id
    * @return _ret If contains return 1, else return 0
    */
    function roleExistsEx(string _roleId) constant public returns (uint _ret) {
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                log("in roleExistsEx, creator:origin:", uint(roles[i].creator).toAddrString(), uint(tx.origin).toAddrString());
                if (roles[i].creator != tx.origin)
                  return uint(-2);
                else
                  return 1;
            }
        }

        return 0;
    }

    /**
    * @dev check if a action is used by any role
    * @param _actionId The action id
    * @return _used If contains return 1, else return 0
    */
    function actionUsed(string _actionId) constant public returns (uint _used) {
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].deleted) {
                continue;
            }

            for (uint j=0; j<roles[i].actionIdList.length; ++j) {
                if (roles[i].actionIdList[j].equals(_actionId)) {
                    return 1;
                }
            }
        }

        return 0;
    }

    /**
    * @dev list elem by condition
    * @param _cond for the condition
    *        0 for all
    *        1 for id
    *        2 for name
    *        3 for moduleId
    *        4 for contractId
    * @param _value The condition value
    * @return _json Objects in json string
    */
    function listByUK(uint _cond, string _value) constant private returns (string _json) {
        
        uint len = 0;
        uint n = 0;
        len = LibStack.push("");
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].deleted)
                continue;

            bool suitable = false;
            if (_cond == 0) {
                suitable = true;
            } else if (_cond == 1) {
                if (roles[i].id.equals(_value))
                    suitable = true;
            } else if (_cond == 2) {
                if (roles[i].name.equals(_value))
                    suitable = true;
            } else if (_cond == 3) {
                if(roles[i].moduleId.equals(_value)){
                    suitable = true;
                }
            } else if (_cond == 4) {
                if(roles[i].contractId.equals(_value)) {
                    suitable = true;
                }
            }

            if (suitable) {
                if (n > 0) {
                    len = LibStack.append(",");
                }

                len = LibStack.append(roles[i].toJson());
                n++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),roles.length);
        _json = LibStack.popex(_retLen);
    }

    /**
    * @dev list elem by condition
    * @param _cond for the condition
    *        0 for all
    *        1 for id
    *        2 for name
    *        3 for moduleName
    * @param _name The condition value
    * @param _version The condition version
    * @return _json Objects in json string
    */
    function listBy(uint _cond, string _name, string _version) constant private returns (string _json) {
        uint len = 0;
        uint n = 0;
        len = LibStack.push("");
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].deleted)
                continue;

            bool suitable = false;
            if (_cond == 0) {
                suitable = true;
            } else if (_cond == 1) {
                if (roles[i].id.equals(_name))
                    suitable = true;
            } else if (_cond == 2) {
                if (roles[i].name.equals(_name))
                    suitable = true;
            } else if (_cond == 3) {
                if(roles[i].moduleName.equals(_name) && roles[i].moduleVersion.equals(_version)){
                    suitable = true;
                }
            }

            if (suitable) {
                if (n > 0) {
                    len = LibStack.append(",");
                }

                len = LibStack.append(roles[i].toJson());
                n++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),roles.length);
        _json = LibStack.popex(_retLen);
    }

    /**
    * @dev Get a specified page of object by name
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByName(string _name, uint _pageNum, uint _pageSize) constant public returns (string _json) {
        uint len = 0;
        uint n = 0;
        uint m = 0;
        len = LibStack.push("");
        for (uint i=0; i<roles.length; ++i) {
            if ((bytes(_name).length == 0 || roles[i].name.equals(_name))
                && !roles[i].deleted) {
                if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                    if (m > 0) {
                        len = LibStack.append(",");
                    }
                    len = LibStack.append(roles[i].toJson());
                    m++;
                }
                if (n >= (_pageNum+1)*_pageSize) {
                    break;
                }
                n++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),roles.length);
        _json = LibStack.popex(_retLen);
    }

    /**
    * @dev Get a specified page of object by name and moduleId
    * @param _moduleId the id of module
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByNameAndModuleId(string _moduleId,string _name, uint _pageNum, uint _pageSize) constant public returns (string _json) {
        log("into pageByNameAndModuleId...","RoleManager");
        uint len = 0;
        uint n = 0;
        uint m = 0;
        uint total = 0;
        len = LibStack.push("");
        for (uint i=0; i<roles.length; ++i) {
            if(bytes(_moduleId).length != 0 && !roles[i].moduleId.equals(_moduleId)){
                continue;
            }
            if ((bytes(_name).length == 0 || roles[i].name.indexOf(_name) != -1)
                && !roles[i].deleted) {
                total++;
            }
        }
        for (i=0; i<roles.length; ++i) {
            if( bytes(_moduleId).length != 0 && !roles[i].moduleId.equals(_moduleId)){
                continue;
            }
            if ((bytes(_name).length == 0 || roles[i].name.indexOf(_name) != -1)
                && !roles[i].deleted) {
                if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                    if (m > 0) {
                        len = LibStack.append(",");
                    }
                    len = LibStack.append(roles[i].toJson());
                    m++;
                }
                if (n >= (_pageNum+1)*_pageSize) {
                    break;
                }
                n++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),total);
        _json = LibStack.popex(_retLen);
    }

    /**
    * @dev Get a specified page of object by name and moduleName
    * @param _moduleName the name of module
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByNameAndModuleName(string _moduleName,string _name, uint _pageNum, uint _pageSize) constant public returns (string _json) {
        log("into pageByNameAndModuleName...","RoleManager");
        uint len = 0;
        uint n = 0;
        uint m = 0;
        uint total = 0;
        len = LibStack.push("");
        for (uint i=0; i<roles.length; ++i) {
            if(bytes(_moduleName).length != 0 && !roles[i].moduleName.equals(_moduleName)){
                continue;
            }
            if ((bytes(_name).length == 0 || roles[i].name.indexOf(_name) != -1)
                && !roles[i].deleted) {
                total++;
            }
        }
        for (i=0; i<roles.length; ++i) {
            if( bytes(_moduleName).length != 0 && !roles[i].moduleName.equals(_moduleName)){
                continue;
            }
            if ((bytes(_name).length == 0 || roles[i].name.indexOf(_name) != -1)
                && !roles[i].deleted) {
                if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                    if (m > 0) {
                        len = LibStack.append(",");
                    }
                    len = LibStack.append(roles[i].toJson());
                    m++;
                }
                if (n >= (_pageNum+1)*_pageSize) {
                    break;
                }
                n++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),total);
        _json = LibStack.popex(_retLen);
    }

    function checkRoleActionWithKey(string _roleId, address _resKey, string _opKey) constant public returns (uint _ret) {
        for (uint i = 0; i < roles.length; ++i) {
            if (!(roles[i].deleted) && roles[i].id.equals(_roleId)) {
                //iterate the action list 
                for (uint j = 0; j < roles[i].actionIdList.length; ++j) {
                    IActionManager am = IActionManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","ActionManager", "0.0.1.0"));
                    if(am.checkActionWithKey(roles[i].actionIdList[j], _resKey, _opKey) == 1) {
                        return 1;
                    }
                }

                return 0;
            }
        }

        return 0;
    }

    function deleteById(string _roleId) public {
        log("deleteById", "RoleManager");
        if (__roleUsed(_roleId) == 1) {
            errno = 15300 + uint(RoleError.ROLE_USED);
            Notify(errno, "role id used");
            return;
        }

        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                roles[i].deleted = true;
                errno = uint(RoleError.NO_ERROR);
                Notify(errno, "delete succ");
                return;
            }
        }

        log("delete role succ", "RoleManager");
        errno = 15300 + uint(RoleError.ID_NOT_EXISTS);
        Notify(errno, "role id not exists");
        return;
    }

    /**
    * @dev Get role id by action id and index. For UserManager contract
    * @param _actionId action id
    * @param _index index
    * @return _roleId role id
    */
    function getRoleIdByActionIdAndIndex(string _actionId, uint _index) constant public returns (uint _roleId) {
        uint sn = 0;
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].deleted) {
                continue;
            }

            if (_actionId.inArray(roles[i].actionIdList)) {
                if (sn == _index) {
                    return roles[i].id.storageToUint();
                }
                sn++;
            }
        }

        return 0;
    }

    function __roleUsed(string _roleId) constant internal returns (uint _ret) {
        // Get UserManager address
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return 0;
        }

        // Get UserManager token
        IUserManager userManager = IUserManager(userManagerAddr);

        if (userManager.roleUsed(_roleId) == 1) {
            return 1;
        }

        // Get DepartmentManager address
        address departmentManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0");
        if (departmentManagerAddr == 0) {
            return 0;
        }

        // Get DepartmentManager token
        /* IDepartmentManager departmentManager = IDepartmentManager(departmentManagerAddr);

        if (departmentManager.roleUsed(_roleId) == 1) {
            return 1;
        } */

        return 0;
    }

    function __actionExists(string _actionId) constant internal returns (uint _ret) {
        // Get ActionManager address
        address actionManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","ActionManager", "0.0.1.0");
        if (actionManagerAddr == 0) {
            return;
        }

        // Get ActionManager token
        IActionManager actionManager = IActionManager(actionManagerAddr);

        return actionManager.actionExists(_actionId);
    }

    function __userExists(address _userAddr) constant internal returns (uint _ret) {
        // Get UserManager address
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return;
        }

        // Get UserManager token
        IUserManager userManager = IUserManager(userManagerAddr);

        return userManager.userExists(_userAddr);
    }

    function __getUserDepartmentId(address _userAddr) constant internal returns (string _ret) {
        // Get UserManager address
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return;
        }

        // Get UserManager token
        IUserManager userManager = IUserManager(userManagerAddr);

        // Get user department id
        return userManager.getUserDepartmentId(_userAddr).recoveryToString();
    }

    function __getUserRoleIdList(address _userAddr, string[] storage roleIdList) constant internal {
        // Get UserManager address
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return;
        }

        // Get UserManager token
        IUserManager userManager = IUserManager(userManagerAddr);

        // msg.sender to userAddr
        address __userAddr = __getUserAddrByAddr(_userAddr);    

        uint i = 0;
        while (true) {
            uint ret = userManager.getUserRoleId(__userAddr, i);
            if (ret == 0) {
                return;
            } else {
                roleIdList.push(ret.recoveryToString());
            }
            i++;
        }
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

    function getActionIdListByRoleIdList(string[] storage roleIdList, string[] storage actionIdList) constant internal {
        for (uint i=0; i<roleIdList.length; ++i) {
            for (uint j=0; j<roles.length; ++j) {
                if (roles[j].id.equals(roleIdList[i])) {
                    for (uint k=0; k<roles[j].actionIdList.length; ++k) {
                        actionIdList.push(roles[j].actionIdList[k]);
                    }
                }
            }
        }
    }
    
    function getRoleModuleId(string _roleId) constant public returns (uint _ret) {
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                return roles[i].moduleId.storageToUint();
            }
        }

        return 0;
    }
    
    function getRoleModuleName(string _roleId) constant public returns (uint _ret) {
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                return roles[i].moduleName.storageToUint();
            }
        }

        return 0;
    }

    function getRoleModuleVersion(string _roleId) constant public returns (uint _ret) {
        for (uint i=0; i<roles.length; ++i) {
            if (roles[i].id.equals(_roleId)) {
                return roles[i].moduleVersion.storageToUint();
            }
        }

        return 0;
    }

    function __getUserAddrByAddr(address _userAddr) constant internal returns(address _address) {
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return;
        }
        IUserManager userManager = IUserManager(userManagerAddr);
        return userManager.getUserAddrByAddr(_userAddr);
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
