pragma solidity ^0.4.12;

contract IDepartmentManager {
    
    /**
    * @dev insert an object
    * @param _json The json string described the object
    * @return No return
    */
    function insert(string _json) public returns(uint);
        
    /**
    * @dev update department
    * @param _json The json string described the object
    * @return No return
    */
    function update(string _json) public returns(uint);
       
    /**
    * @dev Set department status
    * @param _departmentId The dest department id to check
    * @param _status status of dept , 0 disabled 1 enabled
    * @return No return
    */
    function setDepartmentStatus(string _departmentId, uint _status) public returns(uint);
     
    /**
    * @dev Set department admin
    * @param _departmentId The dest department id to check
    * @param _adminAddr The new amdin address
    * @return No return
    */
    function setAdmin(string _departmentId, address _adminAddr) public returns(uint);
     
    /**
    * @dev Erase department admin if admin address equals specified address
    * @param _userAddr The amdin address
    * @return No return
    */
    function eraseAdminByAddress(address _userAddr) public ;
       
    /**
    * @dev Check if a department exists by CN
    * @param _commonName The department commonName to check
    * @return If exists return 1 else return 0
    */
    function departmentExistsByCN(string _commonName) constant public returns (uint _ret) ;
       
    /**
    * @dev Check if a department exists
    * @param _departmentId The department id to check
    * @return If exists return 1 else return 0
    */
    function departmentExists(string _departmentId) constant public returns (uint _exists) ;
       
    /**
    * @dev Check if a department is empty
    * @param _departmentId The department id to check
    * @return If empty return 1 else return 0
    */
    function departmentEmpty(string _departmentId) constant public returns (bool _empty) ;
       
    /**
    * @dev Delete a department (must be empty)
    * @param _departmentId The department id to check
    * @return If Delete succ return 1 else return 0
    */
    function deleteById(string _departmentId) public ;
       
    /**
    * @dev List the all objects
    * @return No return
    */
    function listAll() constant public returns (string _json) ;
     
    /**
    * @dev Find object by id
    * @param _id Object id
    * @return _json Objects in json string
    */
    function findById(string _id) constant public returns(string _json) ;
   
    /**
    * @dev Find object by name
    * @param _name Object name
    * @return _json Objects in json string
    */
    function findByName(string _name) constant public returns(string _json) ;
   
    /**
    * @dev Find object by prarent id
    * @param _parentId Object prarent id
    * @return _json Objects in json string
    */
    function findByParentId(string _parentId) constant public returns(string _json) ;
    
    /**
    * @dev Get department admin
    * @param _departmentId deparment id
    * @return Department admin address
    */
    function getAdmin(string _departmentId) constant public returns(uint _admin) ;
       
    /**
    * @dev Get data revision
    * @return revision id
    */
    function getRevision() constant public returns(uint _revision) ;
   
    /**
    * @dev Get a specified page of object by name
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByName(string _name, uint _pageNum, uint _pageSize) constant public returns(string _json) ;
     
    /**
    * @dev Get a specified page of object by name and status
    * @param _parentId Object prarent id
    * @param _status 0 全部 1 禁用 2 启用
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByNameAndStatus(string _parentId,uint _status, string _name, uint _pageNum, uint _pageSize) constant public returns(string _json) ;
       
    /**
    * @dev Get child department id by index. For UserManager contract
    * @param _departmentId departmentId
    * @param _index index
    * @return _childDepartmentId
    */
    function getChildIdByIndex(string _departmentId, uint _index) constant public returns (uint _childDepartmentId) ;

    function checkWritePermission(address _addr, string _departmentId) constant public returns (uint _ret) ;
}
