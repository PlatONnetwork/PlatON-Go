pragma solidity ^0.4.12;

contract IUserManager {
 
	////Note: find user by uuid
    function findByUuid(string _uuid) constant public returns(string _strjson) ;

	////Note: reset password,
    function resetPasswd(address _userAddr, address _ownerAddr, string _publilcKey, string _cipherGroupKey, string _uuid) public returns(bool _ret) ;
        
	////Note: query user`s status , _state : 0 invalid ,1 valid , 2 logined
    function getUserState(address _userAddr) constant public returns (uint _state) ;
       
	////Note: query accounts state of user , _state : 0 invalid ,1 valid , 2 logined
    function getAccountState(string _account) constant public returns (uint _state) ;
      
	////Note: get users by accountStatus
    function pageByAccountStatus(uint _accountStatus, uint _pageNo, uint _pageSize) public constant returns(string _strjson) ;
     
    ////Note: find user by userAddr
    function findByAddress(address _userAddr) constant public returns(string _ret) ;
      
	////Note: find user by loginName
    function findByLoginName(string _name) constant public returns(string _strjson) ;
       
    ////Note: find user by account
    function findByAccount(string _account) constant public returns(string _strjson);
	
	////Note: find user by mobile
    function findByMobile(string _mobile) constant public returns(string _strjson) ;
      
    ////Note: find user by email
    function findByEmail(string _email) constant public returns(string _strjson);
       
	////Note: find user by departmentId
    function findByDepartmentId(string _departmentId) constant public returns(string _strjson) ;
      
	////Note: find all user by departmentId
    function findByDepartmentIdTree(string _departmentId, uint _pageNum, uint _pageSize) constant public returns(string _strjson) ;
       
	////Note: find all user by condition(_status,_name,_departmentId,_pageNum,_pageSize)
    function findByDepartmentIdTreeAndContion(uint _status,string _name,string _departmentId, uint _pageNum, uint _pageSize) constant public returns(string _strjson) ;
       
	////Note: find user by roleId
    function findByRoleId(string _roleId) constant public returns(string _strjson) ;
       
	////Note: find user`s departmentId by userAddr
    function getUserDepartmentId(address _userAddr) constant returns(uint _departId) ;
       
	////Note: check a userAddr is belong to role
    function checkUserRole(address _userAddr, string _roleId) constant public returns(uint _ret) ;
        
	////Note:
    function checkUserAction(address _userAddr, string _actionId) constant public returns (uint _ret) ;
        
	////Note: check the user is have privilege
    function checkUserPrivilege(address _userAddr, address _contractAddr, string _funcSha3) constant public returns (uint _ret) ;
        
	////Note: check user is exists
    function userExists(address _userAddr) constant public returns(uint _ret) ;

	////Note: insert new user
    function insert(string _userJson) public returns(uint);
       
	////Note: update user info
    function update(string _userJson) public returns(uint) ;

	////Note: update user`s status , _status: 0 off , 1 on
    function updateUserStatus(address _userAddr, uint _status) public returns(bool _ret) ;
    
	////Note: update account`s status of user, _status: 0 invalid ,1 valid
    function updateAccountStatus(address _userAddr, uint _status) public returns(bool _ret) ;
       
	////Note:
    function updatePasswordStatus(address _userAddr, uint _status) public returns(bool _ret) ;
       
	////Note:
    function addUserRole(address _userAddr, string _roleId) returns(uint) ;

	////Note: remove user by userAddr
    function deleteByAddress(address _userAddr) public ;
       
	////Note:
    function login(string _account) public returns(string _json) ;
    
	////Note:
    function listAll() constant public returns(string _userListJson) ;

	////Note:
    function getUserCountByDepartmentId(string _departmentId) constant public returns(uint _count) ;
	
	////Note:
    function getUserCountByActionId(string _actionId) constant public returns(uint _count) ;

	////Note: check roleId is exists ; _used : 1 exists , 0 not exists
    function roleUsed(string _roleId) constant public returns (uint _used) ;
	
	////Note:
    function getUserCountMappingByRoleIds(string _roleIds) constant public returns(string _json) ;
       
    ////Note:
    function getUserAddrByAddr(address _userAddr) constant public returns (address _address) ;

    ////Note:
    function getOwnerAddrByAddr(address _userAddr) constant public returns (address _address) ;

    ////Note:
    function getUserRoleId(address _userAddr, uint _index) constant returns (uint _ret);

    ////Note: check user repetitive
    function isRepetitive(string _mobile, string _email, address _userAddr, string _uuid, string _publicKey, string _account) public constant returns (uint);

    ////Note: check email repetitive
    function checkEmailUniqueness(string _email,string _mobile) public constant returns(uint); 
}