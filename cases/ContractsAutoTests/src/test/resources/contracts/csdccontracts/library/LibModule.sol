pragma solidity ^0.4.12;

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../library/LibContract.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibModule{

	using LibInt for *;
	using LibString for *;
	using LibModule for *;
	using LibContract for *;
	using LibJson for *;

	struct Module{
		string 					moduleId;
		string 					moduleName;
		string 					moduleVersion;
		bool 					deleted;
		uint 					moduleEnable;			// 0 false 1 true
		string 					moduleDescription;		//detail description
		string				    moduleText;				//usually chinese module name
		uint					moduleCreateTime;
		uint					moduleUpdateTime;
		address					moduleCreator;
		LibContract.ContractInfo[] 		contractInfos;
		string[] 				roleIds;
		string 					moduleUrl;
		string 					icon;
		uint 					publishTime;
		uint 					moduleType;				// 模块类型：1 系统模块
	}

	function toJson(Module storage _self) constant internal returns(string _json){
		uint len = 0;
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("moduleId",_self.moduleId);
		len = LibStack.appendKeyValue("moduleName",_self.moduleName);
		len = LibStack.appendKeyValue("moduleText",_self.moduleText);
		len = LibStack.appendKeyValue("moduleVersion",_self.moduleVersion);

		if (_self.deleted) 
			len = LibStack.append(",\"deleted\":true");
		else
			len = LibStack.append(",\"deleted\":false");

		len = LibStack.appendKeyValue("moduleEnable",_self.moduleEnable);
		len = LibStack.appendKeyValue("moduleDescription",_self.moduleDescription);
		len = LibStack.appendKeyValue("moduleCreateTime",_self.moduleCreateTime);
		len = LibStack.appendKeyValue("moduleUpdateTime",_self.moduleUpdateTime);
		len = LibStack.appendKeyValue("moduleCreator",_self.moduleCreator);
		len = LibStack.appendKeyValue("contractInfos",_self.contractInfos.toJsonArray());
		len = LibStack.appendKeyValue("roleIds",_self.roleIds.toJsonArray());
		len = LibStack.appendKeyValue("moduleUrl",_self.moduleUrl);
		len = LibStack.appendKeyValue("icon",_self.icon);
		len = LibStack.appendKeyValue("publishTime",_self.publishTime);
		len = LibStack.appendKeyValue("moduleType",_self.moduleType);
		len = LibStack.append("}");
		_json = LibStack.popex(len);
		
	}

	function fromJson(Module storage _self, string _json) constant internal returns(bool succ){
		if(bytes(_json).length == 0){
			return false;
		}
		LibJson.push(_json);
		_self.clear();
		_self.moduleId = _json.jsonRead("moduleId");
		_self.moduleName = _json.jsonRead("moduleName");
		_self.moduleText = _json.jsonRead("moduleText");
		_self.moduleVersion = _json.jsonRead("moduleVersion");
		_self.moduleEnable = _json.jsonRead("moduleEnable").toUint();
		_self.moduleDescription = _json.jsonRead("moduleDescription");
		_self.moduleCreateTime = _json.jsonRead("moduleCreateTime").toUint();
		_self.moduleUpdateTime = _json.jsonRead("moduleUpdateTime").toUint();
		_self.moduleCreator = _json.jsonRead("moduleCreator").toAddress();
		_self.moduleUrl = _json.jsonRead("moduleUrl");
		_self.icon = _json.jsonRead("icon");
		_self.publishTime = _json.jsonRead("publishTime").toUint();
		_self.moduleType = _json.jsonRead("moduleType").toUint();
		if(bytes(_json.jsonRead("contractInfos")).length > 0){
			_self.contractInfos.fromJsonArray(_json.jsonRead("contractInfos"));
		}
		if(bytes(_json.jsonRead("roleIds")).length > 0){
			_self.roleIds.fromJsonArray(_json.jsonRead("roleIds"));
		}
		LibJson.pop();
		return true;
	}

	function toJsonArray(Module[] storage _self) constant internal returns(string _json){
		uint len = 0;
		len = LibStack.push("[");
		for(uint i=0; i<_self.length;++i){
			if(i > 0){
				len = LibStack.append(",");
			}
			len = LibStack.append(_self[i].toJson());
		}
		len = LibStack.append("]");
		_json = LibStack.popex(len);
	}

	function fromJsonArray(Module[] storage _self, string _json) internal returns(bool){
		LibJson.push(_json);
		_self.length = 0;

		if(!_json.isJson()){
			LibJson.pop();
			return false;
		}

		while(true){
			string memory key = "[".concat(_self.length.toString(), "]");
			if(!_json.jsonKeyExists(key)){
				break;
			}
			_self.length++;
			_self[_self.length-1].fromJson(_json.jsonRead(key));
		}
		LibJson.pop();
		return true;
	}

	function clear(Module storage _self) internal{
		_self.moduleId = "";
		_self.moduleName = "";
		_self.moduleText = "";
		_self.moduleVersion = "";
		_self.deleted = false;
		_self.moduleEnable = 0;
		_self.moduleDescription = "";
		_self.moduleCreateTime = 0;
		_self.moduleUpdateTime = 0;
		_self.moduleCreator = address(0);
		delete _self.contractInfos;
		_self.moduleUrl = "";
		_self.icon = "";
		_self.publishTime = 0;
		_self.moduleType = 1;
		delete _self.roleIds;
	}
}