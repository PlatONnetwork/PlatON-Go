pragma solidity ^0.4.12;

contract IFileInfoManager {
    
	////Note: add fileinfo 
    function insert(string _fileJson) public returns(bool) ;
    	
    ////Note: remove fileInfo by id
    function deleteById(string _fileId) public returns(bool) ;
       
    ////Note: update file info
    function update(string _fileJson) public returns(bool) ;
    	
    ////Note: find a file by fileId and container
    function find(string _fileId) constant public returns(string _ret) ;
     
    ////Note: get the size of FileInfo
    function getCount() constant public returns(uint _ret);
      
    ////Note: get all of fileInfo
    function listAll() constant public returns (string _ret);
       
    ////Note: list all file info by group 
    function listByGroup(string _group) constant public returns (string _ret);
        
    ////Note: list files info by page in pageSize
    function pageFiles(uint _pageNo, uint256 _pageSize) constant public returns(string _ret) ;
      
    ////Note: get the fixed group in a page, pagesize is default 16 */
    function pageByGroup(string _group, uint256 _pageNo) constant public returns(string _ret) ;
        
    ////Note:  get the total pages in a group 
    function getGroupPageCount(string _group) constant public returns(uint256 _ret) ;
        
    ////Note:  get the total files in one group  
    function getGroupFileCount(string _group) constant public returns(uint256 _ret) ;
       
    ////Note: get current page size 
    function getCurrentPageSize() constant public returns(uint256 _ret) ;
        
    ////Note: get default page count 
    function getCurrentPageCount() constant public returns(uint256 _ret) ;
       
    ////Note:  create a unqiue file ID by group, server, filename and time 
    function generateFileID(string _salt, string _groupID, string _serverId, string _filename) constant public returns(string _ret);
      

}
