pragma solidity ^0.4.12;

contract IRoleManager {
   
    /**
     * @dev add new actionId to role
     * @param _roleId role id
     * @param _actionId action id
     * @return no return
     */
    function addActionToRole(string _roleId, string _actionId) public returns(uint);

    /**
    * get role list by module id
    * @param _moduleId module id
    * @return _json role info list
    */
    function getRoleListByModuleId(string _moduleId) constant public returns(string _json) ;

    /**
    * get role list by contract id
    * @param _contract contractId
    * @return _json role info list
    */
    function getRoleListByContractId(string _contract) constant public returns(string _json) ;
    
    /**
    * get role list by moduleName and moduleVersion
    * @param _moduleName module name
    * @param _moduleVersion module version
    * @return _json role info list
    */
    function getRoleListByModuleName(string _moduleName, string _moduleVersion) constant public returns(string _json) ;
  
    /**
    * @dev insert an object
    * @param _json The json string described the object
    * @return No return
    */
    function insert(string _json) public returns(uint) ;
        
    /**
    * @dev update an object
    * @param _json The json string described the object
    * @return No return
    */
    function update(string _json) public ;
       
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
    * @dev check if a role id contains a action id
    * @param _roleId The role id
    * @param _actionId The action id
    * @return _ret If contains return 1, else return 0
    */
    function checkRoleAction(string _roleId, string _actionId) constant public returns (uint _ret) ;
      
    /**
    * @dev check if a role id exists in the roles
    * @param _roleId The role id
    * @return _ret If contains return 1, else return 0
    */
    function roleExists(string _roleId) constant public returns (uint _ret) ;
      
    /**
    * @dev check if a action is used by any role
    * @param _actionId The role id
    * @return _used If contains return 1, else return 0
    */
    function actionUsed(string _actionId) constant public returns (uint _used) ;
       
    /**
    * @dev Get a specified page of object by name
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByName(string _name, uint _pageNum, uint _pageSize) constant public returns (string _json) ;
        
    /**
    * @dev Get a specified page of object by name and moduleId
    * @param _moduleId the id of module
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByNameAndModuleId(string _moduleId,string _name, uint _pageNum, uint _pageSize) constant public returns (string _json) ;
      
    /**
    * @dev Get a specified page of object by name and moduleName
    * @param _moduleName the name of module
    * @param _name Object name
    * @param _pageNum The current page num, 0 for the first
    * @param _pageSize The current page size
    * @return _json Objects in json string
    */
    function pageByNameAndModuleName(string _moduleName,string _name, uint _pageNum, uint _pageSize) constant public returns (string _json) ;
       
    function checkRoleActionWithKey(string _roleId, address _resKey, string _opKey) constant public returns (uint _ret) ;
       
    /**
    * @dev remove role by id
    * @param _roleId role id
    * @return No return
    */
    function deleteById(string _roleId) public ;
       
    /**
    * @dev Get role id by action id and index. For UserManager contract
    * @param _actionId The action id
    * @param _index index
    * @return _roleId
    */
    function getRoleIdByActionIdAndIndex(string _actionId, uint _index) constant public returns (uint _roleId) ;
           
    function getRoleModuleId(string _roleId) constant public returns (uint _ret) ;

    function getRoleModuleName(string _roleId) constant public returns (uint _ret) ;

    function getRoleModuleVersion(string _roleId) constant public returns (uint _ret) ;
}
