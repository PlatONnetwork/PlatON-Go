pragma solidity ^0.4.12;

/**
*@file      LibFile.sol
*@author    kelvin
*@time      2016-11-29
*@desc      the defination of LibFile
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibJson.sol";
import "../utillib/LibStack.sol";

library LibFile {
    using LibInt for *;
    using LibString for *;
    using LibFile for *;
    using LibJson for *;

    /** @brief User state: invalid, valid */
    enum FileState { FILE_INVALID, FILE_VALID}

    /** @brief File Structure */
    struct FileInfo {
        string      id;
        string      container;
        string      filename;
        uint256     updateTime;
        uint256     size;
        string      owner;
        string      file_hash;          // md5 hash
        string      src_node;           //server node id
        string      node_group;
        uint256     Type;               //status
        string      priviliges;
        string      info;
        //uint256     expire;//new add
        FileState   state;
    }


    function toJson(FileInfo storage _self) internal returns(string _strjson){

        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("id",_self.id);
        len = LibStack.appendKeyValue("container",_self.container);
        len = LibStack.appendKeyValue("filename",_self.filename);
        len = LibStack.appendKeyValue("updateTime",_self.updateTime);
        len = LibStack.appendKeyValue("size",_self.size);
        len = LibStack.appendKeyValue("file_hash",_self.file_hash);
        len = LibStack.appendKeyValue("type",_self.Type);
        len = LibStack.appendKeyValue("priviliges",_self.priviliges);
        len = LibStack.appendKeyValue("src_node",_self.src_node);
        len = LibStack.appendKeyValue("node_group",_self.node_group);
        len = LibStack.appendKeyValue("info",_self.info);
        len = LibStack.appendKeyValue("state", uint256(_self.state));
        len = LibStack.push("}");
        _strjson = LibStack.popex(len);
    }

    function jsonParse(FileInfo storage _self, string _json) internal returns(bool) {
        if(bytes(_json).length == 0){
            return false;
        }
        _self.reset();
        LibJson.push(_json);
        if(!_json.isJson()){
            LibJson.pop();
            return false;
        }
        _self.id = _json.jsonRead("id");
        _self.id = _json.jsonRead("id");
        _self.filename = _json.jsonRead("filename");
        _self.container = _json.jsonRead("container");
        _self.size = _json.jsonRead("size").toUint();
        _self.file_hash = _json.jsonRead("file_hash");
        _self.Type = _json.jsonRead("type").toUint();
        _self.priviliges = _json.jsonRead("priviliges");
        _self.info = _json.jsonRead("info");
        _self.src_node = _json.jsonRead("src_node");
        _self.node_group = _json.jsonRead("node_group");

        LibJson.pop();
        if (bytes(_self.id).length == 0) {
            return false;
        }

        return true;
    }

    function reset(FileInfo storage _self) internal {
        _self.id = "";
        _self.filename = "";
        _self.container = "";
        _self.size = 0;
        _self.file_hash = "";
        _self.node_group = "";
        _self.Type = 0;
        _self.updateTime = 0;
        _self.priviliges = "";
        _self.info = "";
        _self.src_node = "";

        _self.state = FileState.FILE_INVALID;
    }
}
