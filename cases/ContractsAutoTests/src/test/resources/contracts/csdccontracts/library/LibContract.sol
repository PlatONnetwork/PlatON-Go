pragma solidity ^0.4.12;
/**
*@file      LibContract.sol
*@author    anrui
*@time      2017-05-08
*@desc      the defination of LibContract
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../interfaces/IRegisterManager.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibContract{

	using LibInt for *;
	using LibString for *;
	using LibContract for *;
    	using LibJson for *;


	struct OpenAction{
		string id;
		string funcHash;
	}

	struct ContractInfo{
		string 					contractId;
		string 					contractName;
		string 					contractVersion;
	}

	struct Contract{
		string 				moduleName;
		string 				moduleVersion;
		string 				moduleId;
		string 		   		cctId;
		string 		   		cctName;
		string 		   		cctVersion;
		bool				deleted;
		uint 		   		enable;
		string 		   		description;
		uint 		   		createTime;
		uint 		   		updateTime;
		address		  		creator;
		string[] 	   		actionIdList;
		//string[]  	  		openMenuIdList;
		OpenAction[]   		openActionList;
		uint				blockNum;
	}

	function toJson(Contract storage _self) internal constant returns (string _json) {
		uint len = 0;
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("moduleId",_self.moduleId);
		len = LibStack.appendKeyValue("cctId",_self.cctId);
		len = LibStack.appendKeyValue("moduleName",_self.moduleName);
		len = LibStack.appendKeyValue("moduleVersion",_self.moduleVersion);
		len = LibStack.appendKeyValue("cctName",_self.cctName);
		len = LibStack.appendKeyValue("cctVersion",_self.cctVersion);
		len = LibStack.appendKeyValue("enable",_self.enable);
		len = LibStack.appendKeyValue("description",_self.description);
		len = LibStack.appendKeyValue("createTime",_self.createTime);
		len = LibStack.appendKeyValue("updateTime",_self.updateTime);
		len = LibStack.appendKeyValue("creator",_self.creator);
        if (_self.deleted) 
            len = LibStack.append(",\"deleted\":true");
        else
            len = LibStack.append(",\"deleted\":false");
		IRegisterManager rm = IRegisterManager(0x0000000000000000000000000000000000000011);
		address cctAddr = rm.getContractAddress(_self.moduleName, _self.moduleVersion, _self.cctName, _self.cctVersion);
		len = LibStack.appendKeyValue("cctAddr",cctAddr);
		len = LibStack.appendKeyValue("actionIdList", _self.actionIdList.toJsonArray());
		len = LibStack.appendKeyValue("blockNum",_self.blockNum);
		len = LibStack.appendKeyValue("openActionList",_self.openActionList.toJsonArray());
		len = LibStack.append("}");
        _json = LibStack.popex(len);
	}

	function fromJson(Contract storage _self, string _json) internal returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
		_self.clear();
		_self.moduleId = _json.jsonRead("moduleId");
		_self.cctId = _json.jsonRead("cctId");
		_self.moduleName = _json.jsonRead("moduleName");
		_self.moduleVersion = _json.jsonRead("moduleVersion");
		_self.cctName = _json.jsonRead("cctName");
		_self.cctVersion = _json.jsonRead("cctVersion");
		_self.enable = _json.jsonRead("enable").toUint();
		_self.description = _json.jsonRead("description");
		_self.createTime = _json.jsonRead("createTime").toUint();
		_self.updateTime = _json.jsonRead("creator").toUint();
		_self.blockNum = _json.jsonRead("blockNum").toUint();
		_self.creator = _json.jsonRead("creator").toAddress();
        if(bytes(_json.jsonRead("actionIdList")).length > 0){
            _self.actionIdList.fromJsonArray(_json.jsonRead("actionIdList"));
        }
        LibJson.pop();
		return true;
	}

	function toJsonArray(Contract[] storage _self)internal constant returns(string _json){
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

    function fromJsonArray(Contract[] storage _self, string _json) internal returns(bool succ) {
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

    function clear(Contract storage _self) internal{
        _self.moduleId = "";
        _self.cctId = "";
        _self.moduleName = "";
        _self.moduleVersion = "";
        _self.cctName = "";
        _self.cctVersion = "";
        _self.enable = 0;
        _self.deleted = false;
        _self.description = "";
        _self.createTime = 0;
        _self.updateTime = 0;
        _self.creator = address(0);
		delete _self.actionIdList.length;
        delete _self.openActionList.length;
    }

     //////////////////////////子结构——————openActionList//////////////////////////

    function toJson(OpenAction storage _self)internal constant returns(string _json){
     	uint len = 0;
     	len = LibStack.push("{");
     	len = LibStack.appendKeyValue("id",_self.id);
     	len = LibStack.appendKeyValue("funcHash",_self.funcHash);
     	len = LibStack.append("}");
     	_json = LibStack.popex(len);
     }

    function fromJson(OpenAction storage _self,string _json)internal returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
     	_self.id = _json.jsonRead("id");
     	_self.funcHash = _json.jsonRead("funcHash");
        LibJson.pop();
     	return true;
     }

    function toJsonArray(OpenAction[] storage _self)internal constant returns(string _json){
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

    function fromJsonArray(OpenAction[] storage _self, string _json) internal returns(bool succ) {
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


    function clear(OpenAction storage _self) internal{
        _self.id = "";
        _self.funcHash = "";
    }


     //////////////////////////结构——————ContractInfo//////////////////////////
    function toJson(ContractInfo storage _self)internal constant returns(string _json){
     	uint len = 0;
     	len = LibStack.push("{");
     	len = LibStack.appendKeyValue("contractId",_self.contractId);
     	len = LibStack.appendKeyValue("contractName",_self.contractName);
     	len = LibStack.appendKeyValue("contractVersion",_self.contractVersion);
    	len = LibStack.append("}");
    	_json = LibStack.popex(len);		
    }
     	
    function fromJson(ContractInfo storage _self,string _json)internal returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
    	_self.contractId = _json.jsonRead("contractId");
    	_self.contractName = _json.jsonRead("contractName");
    	_self.contractVersion = _json.jsonRead("contractVersion");
        LibJson.pop();
    	return true;
    }

    function toJsonArray(ContractInfo[] storage _self)internal constant returns(string _json){
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

    function fromJsonArray(ContractInfo[] storage _self, string _json) internal returns(bool succ) {
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
	
    function clear(ContractInfo storage _self) internal{
        _self.contractId = "";
        _self.contractName = "";
        _self.contractVersion = "";
    }
}