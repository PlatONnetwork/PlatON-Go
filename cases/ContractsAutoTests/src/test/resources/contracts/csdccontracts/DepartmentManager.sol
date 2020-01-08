pragma solidity ^0.4.12;
/**
* @file DepartmentManager.sol
* @author liaoyan
* @time 2016-11-29
* @desc the defination of DepartmentManager contract
*/

import "./library/LibDepartment.sol";
import "./sysbase/OwnerNamed.sol";
import "./interfaces/IDepartmentManager.sol";
import "./interfaces/IRoleManager.sol";
import "./interfaces/IUserManager.sol";

contract DepartmentManager is OwnerNamed, IDepartmentManager {
    using LibDepartment for *;
    using LibString for *;
    using LibInt for *;
    
    LibDepartment.Department[] departments;
    LibDepartment.Department[] tmpDepartments;
    LibDepartment.Department tmpDepartment;
    uint revision;

    enum DepartmentError {
        NO_ERROR,
        BAD_PARAMETER,
        NAME_EMPTY,
        ID_NOT_EXISTS,
        ID_CONFLICTED,
        COMMON_NAME_EMPTY,
        COMMON_NAME_CONFLICTED,
        PARENT_ID_NOT_EXISTS,
        PARENT_ID_IS_SELF,
        ADMIN_NOT_NULL,
        ADMIN_NOT_MEMBER,
        NO_PERMISSION,
        PARENT_ID_CANNOT_UPDATE,
        ENODE_LIST_EMPTY,
        ENODE_ARG_FORMAT_INVALID,
        ENODE_ID_DUPLICATED,
        ENODE_ID_CONFLICTED,
        NOT_EMPTY
    }

    event Notify(uint _errno, string _info);

    /**
    * @dev Contruction
    */
    function DepartmentManager() {
        revision = 0;
        register("SystemModuleManager","0.0.1.0","DepartmentManager", "0.0.1.0");
    }

    /**
    * @dev insert an object
    * @param _json The json string described the object
    * @return No return
    */
    function insert(string _json) public returns(uint){
        log("insert department", "DepartmentManager");
        
        // Decode json
        if (!tmpDepartment.fromJson(_json)) {
            log("json invalid", "DepartmentManager"); //debug
            errno = 15400 + uint(DepartmentError.BAD_PARAMETER);
            Notify(errno, "json invalid");
            return errno;
        }

        /*
        * A new department must matches following conditions:
        * 1. department id DOSE NOT exists,
        * 2. department parent id exists,
        * 3. department admin is a valid user,
        * 4. department admin NOT assigned to any department yet.
        */

        if (getIndexById(tmpDepartment.id) != uint(-1)) {
            log("id conflicted", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.ID_CONFLICTED);
            Notify(errno, "id conflicted");
            return errno;
        }

        if (bytes(tmpDepartment.commonName).length == 0) {
            log("commonName is empty", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.COMMON_NAME_EMPTY);
            Notify(errno, "commonName is empty");
            return errno;
        }

        if (getIndexByCommonName(tmpDepartment.commonName) != uint(-1)) {
            log("commonName conflicted", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.COMMON_NAME_CONFLICTED);
            Notify(errno, "commonName conflicted");
            return errno;
        }

        if (bytes(tmpDepartment.name).length == 0) {
            log("name is empty", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.NAME_EMPTY);
            Notify(errno, "name is empty");
            return errno;
        }

        uint parentIndex = getIndexById(tmpDepartment.parentId);
        if (parentIndex == uint(-1)) {
            if (!(msg.sender == owner && tmpDepartment.id.equals("admin") && tmpDepartment.parentId.equals(""))) {
                log("parent id not exists", "DepartmentManager");
                errno = 15400 + uint(DepartmentError.PARENT_ID_NOT_EXISTS);
                Notify(errno, "parent id not exists");
                return errno;
            }
        }

        if (tmpDepartment.admin != 0) {
            log("no member, cannot be admin", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.ADMIN_NOT_NULL);
            Notify(errno, "new department has no member, can not specify admin");
            return errno;
        }

        if (tx.origin != owner) {
            log("msg.sender is not owner ", "DepartmentManager"); 
            /* if (checkWritePermission(tx.origin, tmpDepartment.parentId) == 0) {
                log("user is not the new department's forefather");
                errno = 15400 + uint(DepartmentError.NO_PERMISSION);
                Notify(errno, "user is not the new department's forefather");
                return errno;
            } */
        }

        /*if (!tmpDepartment.parentId.equals("admin")) {
            // Check if every role is in parent department role id list
            for (i=0; i<tmpDepartment.roleIdList.length; ++i) {
                if (!tmpDepartment.roleIdList[i].inArray(departments[parentIndex].roleIdList)) {
                    log("role not in prarent department", "DepartmentManager"); //debug
                    errno = 15400 + uint(DepartmentError.ROLE_ID_EXCEEDED_DEPT);
                    Notify(errno, "role id does not exists in prarent department");
                    return errno;
                }
            }
        }*/

        // Set level, creation time, update time
        if (tx.origin == owner && tmpDepartment.id.equals("admin") && tmpDepartment.parentId.equals("")) {
            tmpDepartment.departmentLevel = 0;
        } else {
            tmpDepartment.departmentLevel = departments[parentIndex].departmentLevel+1;
        }
        tmpDepartment.creTime = now * 1000;
        tmpDepartment.updTime = now * 1000;
        tmpDepartment.creator = tx.origin;
        tmpDepartment.status = 1;

        // Insert into department list
        departments.push(tmpDepartment);
        errno = uint(DepartmentError.NO_ERROR);
        revision++;

        log("insert department succ", "DepartmentManager"); //debug
        Notify(errno, "insert succ");
        return errno;
    }

    function update(string _json) public returns(uint){
        log("update action", "DepartmentManager");

        // Decode json
        if (!tmpDepartment.fromJson(_json)) {
            log("json invalid", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.BAD_PARAMETER);
            Notify(errno, "json invalid");
            return errno;
        }

        uint index = getIndexById(tmpDepartment.id);
        if (index == uint(-1)) {
            log("department not exists", "DepartmentManager"); //debug
            errno = 15400 + uint(DepartmentError.ID_NOT_EXISTS);
            Notify(errno, "department id dose not exists");
            return errno;
        }

        if (_json.keyExists("commonName")) {
            if (bytes(tmpDepartment.commonName).length == 0) {
                log("commonName is empty", "DepartmentManager");
                errno = 15400 + uint(DepartmentError.COMMON_NAME_EMPTY);
                Notify(errno, "commonName is empty");
                return errno;
            }
        
            uint commonNameIndex = getIndexByCommonName(tmpDepartment.commonName);
            if (commonNameIndex != uint(-1) && commonNameIndex != index) {
                log("commonName conflicted", "DepartmentManager");
                errno = 15400 + uint(DepartmentError.COMMON_NAME_CONFLICTED);
                Notify(errno, "commonName conflicted");
                return errno;
            }
        }

        if (_json.keyExists("parentId") && !tmpDepartment.parentId.equals(departments[index].parentId)) {
            log("can not update parentId");
            errno = 15400 + uint(DepartmentError.PARENT_ID_CANNOT_UPDATE);
            Notify(errno, "can not update parentId");
            return errno;
        }

        // Check if admin is a valid user
        if (tmpDepartment.admin != 0 && tx.origin != owner) {
            if (!__getUserDepartmentId(tmpDepartment.admin).equals(departments[index].id)) {
                log("admin not member of department", "DepartmentManager");
                errno = 15400 + uint(DepartmentError.ADMIN_NOT_MEMBER);
                Notify(errno, "admin is not a member of current department");
                return errno;
            }
        }

        if (tx.origin != owner) {
            log("sender not owner: ", tx.origin); 
            /* if (checkWritePermission(tx.origin, departments[index].id) == 0) {
                log("user not ancestors", "DepartmentManager");
                errno = 15400 + uint(DepartmentError.NO_PERMISSION);
                Notify(errno, "user is not the new department");
                return errno;
            } */
        }

        if (_json.keyExists("name")) {
            departments[index].name = tmpDepartment.name;
        }
        if (_json.keyExists("description")) {
            departments[index].description = tmpDepartment.description;
        }
        if (_json.keyExists("commonName")) {
            departments[index].commonName = tmpDepartment.commonName;
        }
        if (_json.keyExists("stateName")) {
            departments[index].stateName = tmpDepartment.stateName;
        }
        if (_json.keyExists("countryName")) {
            departments[index].countryName = tmpDepartment.countryName;
        }
        if (_json.keyExists("admin")) {
            departments[index].admin = tmpDepartment.admin;
        }
        if (_json.keyExists("orgaShortName")) {
            departments[index].orgaShortName = tmpDepartment.orgaShortName;
        }
        departments[index].updTime = now*1000;

        errno = uint(DepartmentError.NO_ERROR);
        revision++;

        log("update department succ", "DepartmentManager"); 
        Notify(errno, "update succ");
        return errno;
    }

    /**
    * @dev Set department status
    * @param _departmentId The dest department id to check
    * @param _status status of dept , 0 disabled 1 enabled
    * @return No return
    */
    function setDepartmentStatus(string _departmentId, uint _status) public returns(uint){
        log("setDepartmentStatus", "DepartmentManager");
        log("_status:",_status);
        if(_status != 0 && _status != 1 ){
            log("_status not valid ,0 or 1 ", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.BAD_PARAMETER);
            Notify(errno, "_status not valid");
            return errno;
        }

        /* if (checkWritePermission(tx.origin, _departmentId) == 0) {
            log("No permisson", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.NO_PERMISSION);
            Notify(errno, "No permisson");
            return errno;
        } */

        uint index = getIndexById(_departmentId);
        if (index == uint(-1)) {
            log("department id not exists", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.ID_NOT_EXISTS);
            Notify(errno, "department id not exists");
            return errno;
        }

        departments[index].status = _status;

        log("setDepartmentStatus OK", "DepartmentManager");
        errno = uint(DepartmentError.NO_ERROR);
        revision++;
        Notify(errno, "setDepartmentStatus OK");
        return errno;
    }

    /**
    * @dev Set department admin
    * @param _departmentId The dest department id to check
    * @param _adminAddr The new amdin address
    * @return No return
    */
    function setAdmin(string _departmentId, address _adminAddr) public returns(uint){
        log("setAdmin", "DepartmentManager");

        if (checkWritePermission(tx.origin, _departmentId) == 0) {
            log("No permisson", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.NO_PERMISSION);
            return errno;
        }

        uint index = getIndexById(_departmentId);
        if (index == uint(-1)) {
            log("department id not exists", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.ID_NOT_EXISTS);
            return errno;
        }

        if (!__getUserDepartmentId(_adminAddr).equals(_departmentId)) {
            log("admin not member of department", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.ADMIN_NOT_MEMBER);
            return errno;
        }

        // Assign the new address to the dest department
        departments[index].admin = _adminAddr;

        log("setAdmin OK", "DepartmentManager");
        errno = uint(DepartmentError.NO_ERROR);
        revision++;
        return errno;
    }

    /**
    * @dev Erase department admin if admin address equals specified address
    * @param _userAddr The amdin address
    * @return No return
    */
    function eraseAdminByAddress(address _userAddr) public {
        log("eraseAdminByAddress", "DepartmentManager");

        uint index = uint(-1);
        for (uint i=0; i<departments.length; ++i) {
            if (departments[i].deleted)
                continue;
            if (departments[i].admin == _userAddr) {
                index = i;
                break;
            }
        }

        if (index == uint(-1)) {
            log("department id not exists", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.ID_NOT_EXISTS);
            return;
        }

        if (msg.sender != owner) {
            log("msg.sender is not owner: ", msg.sender);
	        if (msg.sender != rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0")) {
	            if (checkWritePermission(msg.sender, departments[index].id) == 0) {
	                log("no permission");
	                errno = 15400 + uint(DepartmentError.NO_PERMISSION);
	                Notify(errno, "no permission");
	                return;
	            }
        	}
        }

		departments[index].admin = address(0);

        log("eraseAdminByAddress OK", "DepartmentManager");
        errno = uint(DepartmentError.NO_ERROR);
        revision++;
    }

    /**
    * @dev Check if a department exists by CN
    * @param _commonName The department commonName to check
    * @return If exists return 1 else return 0
    */
    function departmentExistsByCN(string _commonName) constant public returns (uint _ret) {
        for (uint i = 0; i < departments.length; ++i) {
            if (departments[i].deleted)
                continue;
            if (departments[i].commonName.equals(_commonName)) {
                return 1;
            }
        }
        return 0;
    }

    /**
    * @dev Check if a department exists
    * @param _departmentId The department id to check
    * @return If exists return 1 else return 0
    */
    function departmentExists(string _departmentId) constant public returns (uint _exists) {
        for (uint i=0; i<departments.length; ++i) {
            if (departments[i].deleted)
                continue;
            if (departments[i].id.equals(_departmentId)) {
                return 1;
            }
        }

        return 0;
    }

    /**
    * @dev Check if a department is empty
    * @param _departmentId The department id to check
    * @return If empty return 1 else return 0
    */
    function departmentEmpty(string _departmentId) constant public returns (bool _empty) {
        for (uint i=0; i<departments.length; ++i) {
            if (!departments[i].deleted 
                && departments[i].parentId.equals(_departmentId)) {
                return false;
            }
        }

        if(__getUserCountByDepartmentId(_departmentId) > 0) {
            return false;
        }

        return true;
    }

    /**
    * @dev Delete a department (must be empty)
    * @param _departmentId The department id to check
    * @return If Delete succ return 1 else return 0
    */
    function deleteById(string _departmentId) public {
        uint index = getIndexById(_departmentId);
        if (index == uint(-1)) {
            log("department id not exists", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.ID_NOT_EXISTS);
            return;
        }

        if (msg.sender != owner) {
            log("msg.sender is not owner: ", msg.sender); 
            if (checkWritePermission(msg.sender, departments[index].parentId) == 0) {
                log("no permission", "DepartmentManager");
                errno = 15400 + uint(DepartmentError.NO_PERMISSION);
                Notify(errno, "no permission");
                return;
            }
        }

        if (!departmentEmpty(_departmentId)) {
        	log("department not empty", "DepartmentManager");
            errno = 15400 + uint(DepartmentError.NOT_EMPTY);
            Notify(errno, "department not empty");
            return;
        }

        departments[index].deleted = true;
        log("delete succ", "DepartmentManager");
        errno = uint(DepartmentError.NO_ERROR);
        Notify(errno, "delete succ");
        revision++;
    }

    /**
    * @dev List the all objects
    * @return No return
    */
    function listAll() constant public returns (string _json) {
        _json = listBy(0, "");
    }

    /**
    * @dev Find object by id
    * @param _id Object id
    * @return _json Objects in json string
    */
    function findById(string _id) constant public returns(string _json) {
        _json = listBy(1, _id);
    }

    /**
    * @dev Find object by name
    * @param _name Object name
    * @return _json Objects in json string
    */
    function findByName(string _name) constant public returns(string _json) {
        _json = listBy(2, _name);
    }

    /**
    * @dev Find object by prarent id
    * @param _parentId Object prarent id
    * @return _json Objects in json string
    */
    function findByParentId(string _parentId) constant public returns(string _json) {
        _json = listBy(3, _parentId);
    }

    /**
    * @dev list elem by condition
    * @param _cond for the condition
    *        0 for all
    *        1 for id
    *        2 for name
    *        3 for parent id
    * @param _value The condition value
    * @return json string
    */
    function listBy(uint _cond, string _value) constant private returns(string _json) {
    	uint tatal = 0;
        for (uint i=0; i<departments.length; ++i) {
            if (!departments[i].deleted) {
                tatal++;
            }
        }
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.append(",\"data\":{");
        len = LibStack.appendKeyValue("total", tatal);
        len = LibStack.append(",\"items\":[");
        uint n = 0;
        for (i=0; i<departments.length; ++i) {
            if (departments[i].deleted)
                continue;

            bool suitable = false;
            if (_cond == 0) {
                suitable = true;
            } else if (_cond == 1) {
                if (departments[i].id.equals(_value))
                    suitable = true;
            } else if (_cond == 2) {
                if (departments[i].name.equals(_value))
                    suitable = true;
            } else if (_cond == 3) {
                if (departments[i].parentId.equals(_value))
                    suitable = true;
            }

            if (suitable) {
                if (n > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(departments[i].toJson());
                n++;
            }
        }
        len = LibStack.append("]}}");
        _json = LibStack.popex(len);
    }

    /**
    * @dev Get department admin
    * @param _departmentId deparment id
    * @return Department admin address
    */
    function getAdmin(string _departmentId) constant public returns(uint _admin) {
        uint n = 0;
        for (uint i=0; i<departments.length; ++i) {
            if (departments[i].id.equals(_departmentId)) {
                if (!departments[i].deleted) {
                    return uint(departments[i].admin);
                } else {
                    return 0;
                }
            }
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
    * @dev Get a specified page of object by name
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByName(string _name, uint _pageNum, uint _pageSize) constant public returns(string _json) {
        
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.append(",\"data\":{");
        len = LibStack.appendKeyValue("total", uint(departments.length));
        len = LibStack.append(",\"items\":[");

        uint n = 0;
        uint m = 0;
        for (uint i=0; i<departments.length; ++i) {
            if ((bytes(_name).length == 0 || departments[i].name.equals(_name))
                && !departments[i].deleted) {
                if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                    if (m > 0) {
                        len = LibStack.append(",");
                    }
                    len = LibStack.append(departments[i].toJson());
                    m++;
                }
                if (n >= (_pageNum+1)*_pageSize) {
                    break;
                }
                n++;
            }
        }
        len = LibStack.append("]}}");
        _json = LibStack.popex(len);
    }

    /**
    * @dev Get a specified page of object by name and status
    * @param _parentId 上级部门id
    * @param _status 0 全部 1 禁用 2 启用
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByNameAndStatus(string _parentId,uint _status, string _name, uint _pageNum, uint _pageSize) constant public returns(string _json) {
        log("into pageByNameAndStatus...");
        // 0 全部 1 禁用 2 启用

        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.append(",\"data\":{");
        len = LibStack.append("\"items\":[");

        uint n = 0;
        uint m = 0;
        uint tmp = _status;
        uint total = 0;
        for (uint i=0; i<departments.length; ++i) {
            if ((bytes(_name).length == 0 || departments[i].name.indexOf(_name) != -1)
                && !departments[i].deleted
                && departments[i].parentId.equals(_parentId)) {
                if(_status != 0){
                    if(_status == 1) tmp = 0;
                    if(_status == 2) tmp = 1;
                    if(tmp != departments[i].status){
                        continue;
                    }
                }
                if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                    if (m > 0) {
                        len = LibStack.append(",");
                    }
                    len = LibStack.append(departments[i].toJson());
                    m++;
                }
                if (n >= (_pageNum+1)*_pageSize) {
                    break;
                }
                n++;
            }
        }
        for (i=0; i < departments.length; ++i) {
            if ((bytes(_name).length == 0 || departments[i].name.indexOf(_name) != -1) 
                && !departments[i].deleted
                && departments[i].parentId.equals(_parentId)) {
                if(_status != 0){
                    if(_status == 1) tmp = 0;
                    if(_status == 2) tmp = 1;
                    if(tmp != departments[i].status){
                        continue;
                    }
                }
                total++;
            }
        }
        len = LibStack.append("]");
        len = LibStack.appendKeyValue("total", total);
        len = LibStack.append("}}");
        _json = LibStack.popex(len);
    }

    /**
    * @dev Get child department id by index. For UserManager contract
    * @param _departmentId 部门id
    * @param _index 索引
    * @return _childDepartmentId
    */
    function getChildIdByIndex(string _departmentId, uint _index) constant public returns (uint _childDepartmentId) {
        uint sn = 0;
        for (uint i=0; i<departments.length; ++i) {
            if (departments[i].deleted) {
                continue;
            }

            if (departments[i].parentId.equals(_departmentId)) {
                if (sn == _index) {
                    return departments[i].id.storageToUint();
                }
                sn++;
            }
        }

        return 0;
    }

    function checkWritePermission(address _addr, string _departmentId) constant public returns (uint _ret) {
        address __addr = __getUserAddrByAddr(_addr);
        if (__addr == owner) { // Is current caller is super admin
            return 1;
        }

        uint parentIndex = getIndexById(_departmentId);
        if (parentIndex == uint(-1)) {
            return 0;
        }

        while (true) {
            // If current caller is parent department admin
            if (__addr == departments[parentIndex].admin) {
                return 1;
            }

            // Get parent's index
            parentIndex = getIndexById(departments[parentIndex].parentId);

            // If reach the admin department
            if (parentIndex == uint(-1)) {
                return 0;
            }
        }
    }

    function getIndexById(string _id) constant private returns (uint) {
        for (uint i=0; i<departments.length; ++i) {
            if (departments[i].deleted)
                continue;
            if (departments[i].id.equals(_id))
                return i;
        }

        return uint(-1);
    }

    function getIndexByCommonName(string _commonName) constant private returns (uint) {
        for (uint i=0; i<departments.length; ++i) {
            if (departments[i].deleted)
                continue;
            if (departments[i].commonName.equals(_commonName))
                return i;
        }

        return uint(-1);
    }

    function __roleExists(string _roleId) constant internal returns (uint _ret) {
        // Get RoleManager address
        address roleManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0");
        if (roleManagerAddr == 0) {
            return;
        }

        // Get RoleManager token
        IRoleManager roleManager = IRoleManager(roleManagerAddr);

        return roleManager.roleExists(_roleId);
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

        return userManager.getUserDepartmentId(_userAddr).recoveryToString();
    }

    function __getUserCountByDepartmentId(string _departmentId) constant internal returns(uint _count) {
        // Get UserManager address
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return;
        }

        // Get UserManager token
        IUserManager userManager = IUserManager(userManagerAddr);

        return userManager.getUserCountByDepartmentId(_departmentId);
    }

    function __getUserAddrByAddr(address _userAddr) constant internal returns(address _address) {
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return;
        }
        // Get UserManager token
        IUserManager userManager = IUserManager(userManagerAddr);
        return userManager.getUserAddrByAddr(_userAddr);
    }

    
}
