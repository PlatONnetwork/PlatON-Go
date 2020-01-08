pragma solidity ^0.4.12;
/**
*@file      LibAction.sol
*@author    kelvin
*@time      2016-11-29
*@desc      the defination of LibAction
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibAction {

    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibAction for *;

    enum ActionState { INVALID, VALID, ACTIVATE }

    struct Action {
        string      moduleName;
        string      moduleVersion;
        string      id;
        string      name;
        string      contractId;
        string      moduleId;
        uint        level;
        uint        Type;
        uint        enable;         // 0 false ,1 true
        string      parentId;
        string      url;
        string      description;
        string      resKey;
        string      opKey;
        string      version;
        uint        createTime;
        uint        updateTime;
        ActionState state;
        address     creator; 
    }

    function jsonParse(Action storage _self, string _json) internal returns(bool succ) {
        if(bytes(_json).length == 0){
            return false;
        }
        _self.resetAction();
        LibJson.push(_json);
        _self.moduleName = _json.jsonRead("moduleName");
        _self.moduleVersion = _json.jsonRead("moduleVersion");
        _self.id = _json.jsonRead("id");
        _self.name = _json.jsonRead("name");
        _self.contractId = _json.jsonRead("contractId");
        _self.moduleId = _json.jsonRead("moduleId");
        _self.level = _json.jsonRead("level").toUint();
        _self.Type = _json.jsonRead("type").toUint();
        _self.enable = _json.jsonRead("enable").toUint();
        _self.parentId = _json.jsonRead("parentId");
        _self.url = _json.jsonRead("url");
        _self.description = _json.jsonRead("description");
        _self.resKey = _json.jsonRead("resKey");
        _self.opKey = _json.jsonRead("opKey");
        _self.version = _json.jsonRead("version");
        _self.createTime = _json.jsonRead("createTime").toUint();
        _self.updateTime = _json.jsonRead("updateTime").toUint();
        _self.state = ActionState(_json.jsonRead("state").toUint());
        _self.creator = _json.jsonRead("creator").toAddress();
        LibJson.pop();
        return true;
    }

    function toJson(Action storage _self) internal constant returns (string _json) {
        uint len = 0;
        string memory creatorAddr = "0x";
        creatorAddr = creatorAddr.concat(_self.creator.addrToAsciiString());
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("moduleName", _self.moduleName);
        len = LibStack.appendKeyValue("moduleVersion", _self.moduleVersion);
        len = LibStack.appendKeyValue("id", _self.id);
        len = LibStack.appendKeyValue("name", _self.name);
        len = LibStack.appendKeyValue("contractId",_self.contractId);
        len = LibStack.appendKeyValue("moduleId",_self.moduleId);
        len = LibStack.appendKeyValue("level", _self.level);
        len = LibStack.appendKeyValue("type", _self.Type);
        len = LibStack.appendKeyValue("enable", _self.enable);
        len = LibStack.appendKeyValue("parentId", _self.parentId);
        len = LibStack.appendKeyValue("url", _self.url);
        len = LibStack.appendKeyValue("description", _self.description);
        len = LibStack.appendKeyValue("resKey", _self.resKey);
        len = LibStack.appendKeyValue("opKey", _self.opKey);
        len = LibStack.appendKeyValue("version", _self.version);
        len = LibStack.appendKeyValue("createTime", _self.createTime);
        len = LibStack.appendKeyValue("updateTime", _self.updateTime);
        len = LibStack.appendKeyValue("state", uint(_self.state));
        len = LibStack.appendKeyValue("creator", creatorAddr);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJsonArray(Action[] storage _self, string _json) internal returns(bool succ) {
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
            _self[_self.length-1].jsonParse(_json.jsonRead(key));
        }

        LibJson.pop();
        return true;
    }

    function toJsonArray(Action[] storage _self) internal constant returns(string _json) {
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

    function update(Action storage _self, string _json) internal returns(bool succ) {
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
        if (_json.jsonKeyExists("contractId"))
            _self.contractId = _json.jsonRead("contractId");
        if (_json.jsonKeyExists("moduleId"))
            _self.moduleId = _json.jsonRead("moduleId");
        if (_json.jsonKeyExists("level"))
            _self.level = _json.jsonRead("level").toUint();
        if (_json.jsonKeyExists("type"))
            _self.Type = _json.jsonRead("type").toUint();
        if (_json.jsonKeyExists("enable"))
            _self.enable = _json.jsonRead("enable").toUint();
        if (_json.jsonKeyExists("parentId"))
            _self.parentId = _json.jsonRead("parentId");
        if (_json.jsonKeyExists("url"))
            _self.url = _json.jsonRead("url");
        if (_json.jsonKeyExists("description"))
            _self.description = _json.jsonRead("description");
        if (_json.jsonKeyExists("resKey"))
            _self.resKey = _json.jsonRead("resKey");
        if (_json.jsonKeyExists("opKey"))
            _self.opKey = _json.jsonRead("opKey");
        if (_json.jsonKeyExists("version"))
            _self.version = _json.jsonRead("version");
        if (_json.jsonKeyExists("createTime"))
            _self.createTime = _json.jsonRead("createTime").toUint();
        if (_json.jsonKeyExists("updateTime"))
            _self.updateTime = _json.jsonRead("updateTime").toUint();
        if (_json.jsonKeyExists("state"))
            _self.state = ActionState(_json.jsonRead("state").toUint());
        if (_json.jsonKeyExists("creator"))
            _self.creator = _json.jsonRead("creator").toAddress();

        LibJson.pop();
        return true;
    }

    function resetAction(Action storage _self) internal {
        delete _self.moduleName;
        delete _self.moduleVersion;
        delete _self.id;
        delete _self.name;
        delete _self.contractId;
        delete _self.moduleId;
        delete _self.level;
        delete _self.Type;
        delete _self.enable;
        delete _self.parentId;
        delete _self.url;
        delete _self.description;
        delete _self.resKey;
        delete _self.opKey;
        delete _self.version;
        delete _self.createTime;
        delete _self.updateTime;
        delete _self.state;
        delete _self.creator;
    }

}
