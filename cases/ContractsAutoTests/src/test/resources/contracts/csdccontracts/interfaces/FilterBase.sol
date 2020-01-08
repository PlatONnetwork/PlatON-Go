pragma solidity ^0.4.12;

contract FilterBase {
	
    function authorizeProcessor(address _from, address _to, string _funcHash, string _extraData) public constant returns(uint _ret);

    function addFilter(string _filterJson) public returns(uint _filterId);
}

