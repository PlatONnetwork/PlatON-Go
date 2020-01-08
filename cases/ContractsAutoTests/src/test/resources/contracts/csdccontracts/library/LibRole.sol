pragma solidity ^0.4.12;
/**
* @file LibRole.sol
* @author liaoyan
* @time 2016-11-29
* @desc the defination of Role library
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibRole {

    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibRole for *;

    struct Role {
        string moduleName;
        string moduleVersion;
        string id;
        string name;
        uint status;
        string moduleId;
        string contractId;             // -
        string description;
        uint creTime;
        uint updTime;
        bool deleted;
        string[] actionIdList;
        address creator;
    }

    function fromJson(Role storage _self, string _json) internal returns(bool succ) {
        _self.clear();
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        _self.moduleName = _json.jsonRead("moduleName");
        _self.moduleVersion = _json.jsonRead("moduleVersion");
        _self.id = _json.jsonRead("id");
        _self.name = _json.jsonRead("name");
        _self.status = _json.jsonRead("status").toUint();
        _self.moduleId = _json.jsonRead("moduleId");
        _self.contractId = _json.jsonRead("contractId");
        _self.description = _json.jsonRead("description");
        _self.creTime = _json.jsonRead("creTime").toUint();
        _self.updTime = _json.jsonRead("updTime").toUint();
        if(bytes(_json.jsonRead("actionIdList")).length > 0){
            _self.actionIdList.fromJsonArray(_json.jsonRead("actionIdList"));
        }
        _self.creator = _json.jsonRead("creator").toAddress();
        LibJson.pop();
        return true;
    }

    function toJson(Role storage _self) internal constant returns (string _json) {
        uint len = 0;
        string memory strAddr = "0x";
        strAddr = strAddr.concat(_self.creator.addrToAsciiString());
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("moduleId",_self.moduleId);
        len = LibStack.appendKeyValue("moduleName",_self.moduleName);
        len = LibStack.appendKeyValue("moduleVersion",_self.moduleVersion);
        len = LibStack.appendKeyValue("id",_self.id);
        len = LibStack.appendKeyValue("name",_self.name);
        len = LibStack.appendKeyValue("status",uint(_self.status));
        len = LibStack.appendKeyValue("description",_self.description);
        len = LibStack.appendKeyValue("creTime",uint(_self.creTime));
        len = LibStack.appendKeyValue("updTime",uint(_self.updTime));
        len = LibStack.appendKeyValue("actionIdList",_self.actionIdList.toJsonArray());
        len = LibStack.appendKeyValue("creator", strAddr);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJsonArray(Role[] storage _self, string _json) internal returns(bool succ) {
        LibJson.push(_json);
        _self.length = 0;

        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }

        while (true) {
            string memory key = "[".concat(_self.length.toString(), "]");
            if (!_json.jsonKeyExists(key))
                break;

            _self.length++;
            _self[_self.length-1].fromJson(_json.jsonRead(key));
        }

        LibJson.pop();
        return true;
    }

    function toJsonArray(Role[] storage _self) internal constant returns(string _json) {
        uint len = 0;
        len = LibStack.push("[");
        for (uint i=0; i<_self.length; ++i) {
            if (i > 0)
                len = LibStack.append(",");
            len = LibStack.append(_self[i].toJson());
        }
        len = LibStack.append("]");
        _json = LibStack.popex(len);
    }

    function update(Role storage _self, string _json) internal returns(bool succ) {
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }

        if (_json.jsonKeyExists("moduleName"))
            _self.moduleName = _json.jsonRead("moduleName");
        if (_json.jsonKeyExists("moduleVersion"))
            _self.moduleVersion = _json.jsonRead("moduleVersion");
        if (_json.jsonKeyExists("id"))
            _self.id = _json.jsonRead("id");
        if (_json.jsonKeyExists("name"))
            _self.name = _json.jsonRead("name");
        if (_json.jsonKeyExists("status"))
            _self.status = _json.jsonRead("status").toUint();
        if (_json.jsonKeyExists("moduleId"))
            _self.moduleId = _json.jsonRead("moduleId");
        if (_json.jsonKeyExists("contractId"))
            _self.contractId = _json.jsonRead("contractId");
        if (_json.jsonKeyExists("description"))
            _self.description = _json.jsonRead("description");
        if (_json.jsonKeyExists("creTime"))
            _self.creTime = _json.jsonRead("creTime").toUint();
        if (_json.jsonKeyExists("updTime"))
            _self.updTime = _json.jsonRead("updTime").toUint();
        if (_json.jsonKeyExists("actionIdList"))
            _json.getArrayValueByKey("actionIdList",_self.actionIdList);
        if (_json.jsonKeyExists("creator"))
            _self.creator = _json.jsonRead("creator").toAddress();

        LibJson.pop();
        return true;
    }

    function clear(Role storage _self) internal {
        delete _self.moduleName;
        delete _self.moduleVersion;
        delete _self.id;
        delete _self.name;
        delete _self.status;
        delete _self.moduleId;
        delete _self.contractId;
        delete _self.description;
        delete _self.creTime;
        delete _self.updTime;
        delete _self.deleted;
        delete _self.actionIdList;
        delete _self.creator;
    }
}