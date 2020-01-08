pragma solidity ^0.4.12;

contract IFileServerManager {
    
    ////Note: insert a file server 
    function insert(string _json) public returns(bool) ;
        
    ////Note: delete fileServer by id
    function deleteById(string _serverId) public returns(bool) ;
        
    ////Note: update server info 
    function update(string _json) public returns(bool) ;
        
    ////Note: enable or disable the file service of the node, 0 is disable, 1 or else is enable 
    function enable(string _serverId, uint256 _enable) public returns(bool) ;
        
    ////Note: find a file by fileId and container 
    function find(string _serverId) constant public returns(string _ret) ;

    ////Note: find a file by fileId and container 
    function isServerEnable(string _serverId) constant public returns(uint256 _ret) ;
       
    ////Note: get the count of all servers includes disable servers 
    function getCount() constant public returns(uint256 _total) ;
       
    ////Note: list all the servers info 
    function listAll() constant public returns(string _ret);
        
    ////Note: list server by group name 
    function listByGroup(string _group) constant public returns(string _ret);
    
    ////Note: get fileServer by host and port
    function findIdByHostPort(string _host, uint256 _port) constant public returns(string _ret) ;
        
}