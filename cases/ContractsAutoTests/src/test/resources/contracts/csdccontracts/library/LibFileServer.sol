pragma solidity ^0.4.12;
/**
*@file      LibFileServer.sol
*@author    kelvin
*@time      2016-11-29
*@desc      the defination of LibFileServer
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibJson.sol";
import "../utillib/LibStack.sol";

library LibFileServer {

    using LibInt for *;
    using LibString for *;
    using LibFileServer for *;
    using LibJson for *;


    /** @brief User state: invalid, valid */
    enum ServerState { INVALID, VALID}

    /** @brief File Structure */
    struct FileServerInfo {
        string      id;
        string      host;
        uint256     port;
        uint256     updateTime;
        string      organization;
        string      position;
        string      group;
        string      info;
        uint256     enable;
        ServerState   state;
    }


    function toJson(FileServerInfo storage _self) internal returns(string _strjson){

        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("id", _self.id);
        len = LibStack.appendKeyValue("host", _self.host);
        len = LibStack.appendKeyValue("port", _self.port);
        len = LibStack.appendKeyValue("updateTime", _self.updateTime);
        len = LibStack.appendKeyValue("position", _self.position);
        len = LibStack.appendKeyValue("group", _self.group);
        len = LibStack.appendKeyValue("organization", _self.organization);
        len = LibStack.appendKeyValue("enable", _self.enable);
        len = LibStack.appendKeyValue("info", _self.info);
        len = LibStack.append("}");
        _strjson = LibStack.popex(len);
    
    }

    function jsonParse(FileServerInfo storage _self, string _json) internal returns(bool) {
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
        _self.host = _json.jsonRead("host");
        _self.port = _json.jsonRead("port").toUint();
        _self.position = _json.jsonRead("position");
        _self.group = _json.jsonRead("group");
        _self.organization = _json.jsonRead("organization");
        _self.info = _json.jsonRead("info");
        _self.enable = _json.jsonRead("enable").toUint();
        LibJson.pop();
        if (bytes(_self.id).length == 0) {
            return false;
        }
        return true;
    }

    function reset(FileServerInfo storage _self) internal {
        _self.id = "";
        _self.host = "";
        _self.port = 0;
        _self.updateTime = 0;
        _self.position = "";
        _self.organization = "";
        _self.info = "";
        _self.group = "";
        _self.enable = 0;

        _self.state = ServerState.INVALID;
    }
}
