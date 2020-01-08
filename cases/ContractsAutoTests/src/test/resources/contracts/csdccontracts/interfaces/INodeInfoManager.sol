pragma solidity ^0.4.12;

contract INodeInfoManager {

    /**
    * @dev insert node info
    * @param _json The json string described the object
    * @return errno , 0 for true
    */
    function insert(string _json) public returns(uint) ;
       
    /**
    * @dev update node info
    * @param _json The json string described the object
    * @return errno , 0 for true
    */
    function update(string _json) public returns(uint);
      
    /**
    * @dev Get enode list
    * @return return enodes in json string
    */
    function getEnodeList() constant public returns(string _json) ;
       
    /**
    * @dev Activate enode for consensus
    * @param _pubkey enode pubkey in hex string mode
    */
    function ActivateEnode(string _pubkey) public ;
        
    /**
    * @dev If the specified name and ip is in the database
    * @param _commonName Object name
    * @param _ip ip of int
    * @return _json true or false in json string
    */
    function isInWhiteList(string _commonName, string _ip) constant public returns (string _json) ;
     
    /**
    * @dev Get Node admin
    * @param _nodeId deparment id
    * @return Node admin address
    */
    function getNodeAdmin(string _nodeId) constant public returns(uint _admin) ;
     
    /**
    * @dev Set node admin
    * @param _nodeId The dest node id to check
    * @param _adminAddr The new amdin address
    * @return errno , 0 for true
    */
    function setAdmin(string _nodeId, address _adminAddr) public returns(uint) ;
       
    /**
    * @dev Erase node admin if admin address equals specified address
    * @param _userAddr The amdin address
    * @return errno , 0 for success
    */
    function eraseAdminByAdd(address _userAddr) public returns(uint) ;
       
    /**
    * @dev Check if a node exists
    * @param _nodeId The node id to check
    * @return If exists return 1 else return 0
    */
    function nodeInfoExists(string _nodeId) constant public returns (uint _exists) ;
       
    /**
    * @dev Delete a department (must be empty)
    * @param _nodeId The department id to check
    * @return If Delete succ return 1 else return 0
    */
    function deleteById(string _nodeId) public ;
        
    function getRevision() constant public returns (uint _ret) ;
     
    /**
    * @dev check if the IP is used by any node
    * @param _ip The role id
    * @return _used If contains return 1, else return 0
    */
    function IPUsed(string _ip) constant public returns (uint _used) ;
       
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
    * @dev Find object by department id
    * @param _departmentId Object prarent id
    * @return _json Objects in json string
    */
    function findByDepartmentId(string _departmentId) constant public returns(string _json) ;
    
    /**
    * @dev Find object by nodeAdmin id
    * @param _nodeAdmin Object prarent id
    * @return _json Objects in json string
    */
    function findByNodeAdmin(address _nodeAdmin) constant public returns(string _json) ;
 
    /**
    * @dev Find object by pubkey
    * @param _pubkey Object prarent id
    * @return _json Objects in json string
    */
    function findByPubkey(string _pubkey) constant public returns(string _json) ;
   
    function checkWritePermission(address _addr, string _nodeInfoId) constant public returns (uint _ret) ;
    
}
