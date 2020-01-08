pragma solidity ^0.4.12;

contract IRoleFilterManager {


	function addActionToRole(string _moduleId, string _roleId, string _actionId) public returns(uint);

	/**
	* execute authorization
	* @param _from : user addr
	* @param _to : contract addr
	* @param _funcHash : openkey sha3 value
	* @param _extraData : additional data -> [;"moduleId":""}],指明要进行验证的模块
	* @return 0 for false ,1 for true
	*/
	function authorizeProcessor(address _from, address _to, string _funcHash, string _extraData) public constant returns(uint);
		
	////Note: add open action
	function addOpenAction(string _moduleId, address _contractAddr, string _funcHash) constant private returns (uint _ret) ;//just one funcHash:f6bfc763
	
	////Note: add open action
	function addOpenAction(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, string _actionId) constant private returns (uint _ret);
		
	////Note: delete the open action from record
	function deleteOpenAction(string _moduleId, address _contractAddr, string _funcHash) constant private returns (uint _ret) ;
	
	function deleteOpenAction(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, string _actionId) constant private returns (uint _ret);
	
	////Note: add role 
	function addAuthorizeRole(string _moduleId, string _roleId, string _actionInfo) constant private returns (uint _ret) ;//just one actionInfo:0000000000000000000000000000000000000012,78a9eeed
	
	////Note: add module
	function addModule(string _json) public returns(uint) ;
	
	////Note: update module
	function updModule(string _json) public returns(uint) ;

	////Note: remove module by id
	function delModule(string _moduleId) public returns(uint) ;
		
	/**
	* enable the module switch
	* if set enalbe to 0,then do not authentication 
	* if set enable to 1,then do authentication
	* @param _moduleId module id
	* @param _enable 1 for on,0 for off
	* @return 0 for false,1 for true
	*/
	//function setModuleEnable(string _moduleId,uint _enable) public returns(uint) ;
		
	/**
	* enable the contract switch
	* if set enalbe to 0,then do not authentication 
	* if set enable to 1,then do authentication
	* @param _contractId module id
	* @param _enable 1 for on,0 for off
	* @return 0 for false,1 for true
	*/
	//function setConntractEnable(string _contractId,uint _enable) public returns(uint) ;
		
	/**
	* Add Contract
	* 
	* @param _json the json struct for contract
	* @return 0 for false,other for true (eg:contractId)
	*/
	function addContract(string _json) public returns(uint) ;
		
	/**
	* Add Menu
	* 
	* @param _json the json struct for action
	* @return 0 for false,other for true(eg:menuId)
	*/
	function addMenu(string _json) public returns(uint) ;
		
	/**
	* Add Action
	* 
	* @param _json the json struct for action
	* @return 0 for false,other for true(eg:menuId)
	*/
	function addAction(string _json) public returns(uint) ;
		
	/**
	* Add Role
	* 
	* @param _json the json struct for Role
	* @return 0 for false,other for true(eg:roleId)
	*/
	function addRole(string _json) public returns(uint) ;
	
	////Note: get all modules
	function listAll() constant public returns (string _json) ;
	
	////Note: get all module
	function qryModules() constant public returns (string _json) ;

	////Note: get module by id
	function qryModuleDetail(string _moduleId) constant public returns (string _json) ;

	////Note: get module by id
	function qryModuleDetail(string _moduleName, string _moduleVersion) constant public returns (string _json) ;

	////Note: get module by moduleText
	function findByModuleText(string _moduleText) constant public returns (string _json) ;
	
	////Note: get module by moduleName
	function findByName(string _name) constant public returns (string _json) ;
	
	////Note: get constrasts by moduleName
	//function findContractByModName(string _moduleName, string _moduleVersion) constant public returns (string _json) ;

	////Note: get constrasts by moduleText
	function findContractByModText(string _moduleText) constant public returns (string _json) ;

	////Note: get constrasts by moduleText and contract name
	function listContractByModTextAndCttName(string _moduleText,string _cttName,uint _pageNum,uint _pageSize) constant public returns (string _json);
	
	////Note: get module count
	function getModuleCount() constant public returns (uint _count) ;
	
	////Note: get all contracts by module id
	function listContractByModuleId(string _moduleId) constant public returns (string _json) ;

	////Note: get all contracts by moduleName
	function listContractByModuleName(string _moduleName, string _moduleVersion) constant public returns (string _json) ;

	/**
	* change module owner, and all contracts owner
	* @param _moduleName moduleName
	* @param _moduleVersion moduleVersion
	* @param _newOwner new owner
	* @return 0 for success, else -1 for failed
	*/
	function changeModuleOwner(string _moduleName, string _moduleVersion, address _newOwner) public returns(uint) ;

	function moduleIsExist(string _moduleId) public constant returns(uint);
		
}
