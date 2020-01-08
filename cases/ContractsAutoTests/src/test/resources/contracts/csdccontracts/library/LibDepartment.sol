pragma solidity ^0.4.12;
/**
* @file LibDepartment.sol
* @author liaoyan
* @time 2016-12-01
* @desc the defination of Department library
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibDepartment {

    using LibInt for *;
    using LibString for *;
    using LibDepartment for *;
    using LibJson for *;

    struct Department {
        string id;
        string name;
        uint departmentLevel;
        string parentId;
        string description;
        uint creTime;
        uint updTime;
        string commonName;
        string stateName;
        string countryName;
        address admin;                  
        address creator;                
        bool deleted;
        string orgaShortName;           // -
        uint status;                    // 0 禁用 1 激活
        string groupPubkey;             // 群公钥
        string icon;                
    }

    function toJson(Department storage _self) internal constant returns(string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("id", _self.id);
        len = LibStack.appendKeyValue("name", _self.name);
        len = LibStack.appendKeyValue("departmentLevel", _self.departmentLevel);
        len = LibStack.appendKeyValue("parentId", _self.parentId);
        len = LibStack.appendKeyValue("description", _self.description);
        len = LibStack.appendKeyValue("creTime", _self.creTime);
        len = LibStack.appendKeyValue("updTime", _self.updTime);
        len = LibStack.appendKeyValue("commonName", _self.commonName);
        len = LibStack.appendKeyValue("stateName", _self.stateName);
        len = LibStack.appendKeyValue("countryName", _self.countryName);
        len = LibStack.appendKeyValue("orgaShortName", _self.orgaShortName);
        len = LibStack.appendKeyValue("admin", _self.admin);
        len = LibStack.appendKeyValue("status", _self.status);
        len = LibStack.appendKeyValue("icon", _self.icon);
        len = LibStack.appendKeyValue("groupPubkey", _self.groupPubkey);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJson(Department storage _self, string _json) internal constant returns(bool succ) {
        if(bytes(_json).length == 0){
            return false;
        }
        _self.clear();
        LibJson.push(_json);
        _self.id = _json.jsonRead("id");
        _self.name = _json.jsonRead("name");
        _self.departmentLevel = _json.jsonRead("departmentLevel").toUint();
        _self.parentId = _json.jsonRead("parentId");
        _self.description = _json.jsonRead("description");
        _self.creTime = _json.jsonRead("creTime").toUint();
        _self.updTime = _json.jsonRead("updTime").toUint();
        _self.commonName = _json.jsonRead("commonName");
        _self.stateName = _json.jsonRead("stateName");
        _self.countryName = _json.jsonRead("countryName");
        _self.orgaShortName = _json.jsonRead("orgaShortName");
        _self.admin = _json.jsonRead("admin").toAddress();
        _self.status = _json.jsonRead("status").toUint();
        _self.groupPubkey = _json.jsonRead("groupPubkey");
        LibJson.pop();
        if (bytes(_self.id).length == 0) {
            return false;
        }
        return true;
    }

    function clear(Department storage _self) internal {
        _self.admin = address(0);
        _self.id = "";
        _self.name = "";
        _self.departmentLevel = 0;
        _self.parentId = "";
        _self.description = "";
        _self.creTime = 0;
        _self.updTime = 0;
        _self.commonName = "";
        _self.stateName = "";
        _self.countryName = "";
        _self.orgaShortName = "";
        _self.deleted = false;
        _self.status=1;
        _self.groupPubkey = "";
        _self.icon = "";
    }
}
