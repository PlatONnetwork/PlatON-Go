pragma solidity ^0.4.12;
/**
* @file LibFilter.sol
* @author ximi
* @time 2017-05-08
* @desc the defination of Filter library
*/

import "../library/Errnos.sol";
import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibJson.sol";
import "../utillib/LibStack.sol";

library LibFilter {
    
    using Errnos for *;
    using LibInt for *;
    using LibString for *;
    using LibFilter for *;
    using LibJson for *;
    
    struct Filter {
        string id;
        string name;
        string version;
        uint _type;
        uint state;
        uint enable;
        string desc;
        address addr;
    }
    
    function toJson(Filter storage _self) internal constant returns(string _json) {
        uint len = 0 ;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("id", _self.id);
        len = LibStack.appendKeyValue("name", _self.name);
        len = LibStack.appendKeyValue("version", _self.version);
        len = LibStack.appendKeyValue("type", _self._type);
        len = LibStack.appendKeyValue("state", _self.state);
        len = LibStack.appendKeyValue("enable", _self.enable);   
        len = LibStack.appendKeyValue("desc", _self.desc);     
        len = LibStack.appendKeyValue("addr", _self.addr);    
        len = LibStack.append("}");    
        _json = LibStack.popex(len);
    }
    
    function fromJson(Filter storage _self, string _json) internal constant returns(bool _succ) {
        if(bytes(_json).length == 0){
            return false;
        }
        _self.clear();
        LibJson.push(_json);
        if(!_json.isJson()){
            LibJson.pop();
            return false;
        }
        _self.id = _json.jsonRead("id");
        _self.name = _json.jsonRead("name");
        _self.version = _json.jsonRead("version");
        _self._type = _json.jsonRead("type").toUint();
        _self.state = _json.jsonRead("state").toUint();
        _self.enable = _json.jsonRead("enable").toUint();
        _self.desc = _json.jsonRead("desc");
        _self.addr = _json.jsonRead("addr").toAddress();
        LibJson.pop();
        if (bytes(_self.id).length == 0) {
            return false;
        }
        return true;
    }
    
    function clear(Filter storage _self) internal {
        _self.id = "";
        _self.name = "";
        _self.version = "";
        _self._type = uint(Errnos.FilterType.FILTER_TYPE_START);
        _self.state = uint(Errnos.State.STATE_INVALID);
        _self.enable = 0;
        _self.desc = "";
        _self.addr = 0;
    }
}