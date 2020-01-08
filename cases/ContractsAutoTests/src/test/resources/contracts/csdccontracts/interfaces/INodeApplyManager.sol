pragma solidity ^0.4.12;

contract INodeApplyManager {
	
	////Note: add nodeApply record
	function insert(string _json) public returns(uint);
		
	////Note: update nodeApply record
	function update(string _json) public returns(uint);

	////Note: to judge the nodeApplyId is exists
	function nodeApplyExists(string _nodeApplyId) constant public returns(uint _exists);
	
	////Note:
	function auditing(string _json) public returns(uint) ;
	
	////Note:	page query by name and status
	function pageByNameAndStatus(uint _status, string _deptName, uint _pageNum, uint _pageSize) constant public returns(string _json) ;
      
    ////Note: list all apply record
	function listAll() constant public returns (string _json);

	////Note: get apply by state
	function findByState(uint _state) constant public returns(string _strjson) ;
    
    ////Note: get apply record by id
	function findByApplyId(string _nodeApplyId) constant public returns (string _json);

	////Note: remove apply record by id
	function deleteById(string _nodeApplyId) public returns(bool);
	
}