pragma solidity ^0.4.12;
/**
*@file      LibUser.sol
*@author    kelvin
*@time      2016-11-29
*@desc      the defination of LibUser
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";



library LibUser {

    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibUser for *;


    
    enum UserState { USER_INVALID, USER_VALID, USER_LOGIN }

    enum AccountState { INVALID, VALID, LOCKED }

    struct User {
         address    userAddr;
         string     name;
         string     account;                
         string     email;              
         string     mobile;             
         string     departmentId;       
         uint       accountStatus;       
         uint       passwordStatus;          // -  
         uint       deleteStatus;            // -         
         string     uuid;                   
         string     publicKey;              
         string     cipherGroupKey;
         uint       createTime;
         uint       updateTime;
         uint       loginTime;                // -
         string[]   roleIdList;
         UserState  state;                  // 标记用户数据是否有效
         uint       status;                 // 是否禁用用户：0 禁用 1 激活
         address    ownerAddr;              // 用户真实地址
         uint       certType;               // 证书类型 1 文件证书 2 U-KEY证书
         string     remark;
         string     icon;                   // icon base64编码

    }
    
    struct ModuleRoles {                   
    	string moduleId;
	    string moduleName;
      	string moduleVersion;
    	string roleIds;
    }

    function toJson(User storage _self) internal returns(string _strjson){
        string memory strAddr = "0x";
        strAddr = strAddr.concat(_self.userAddr.addrToAsciiString());
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("userAddr", strAddr);
        len = LibStack.appendKeyValue("name", _self.name);
        len = LibStack.appendKeyValue("account", _self.account);
        len = LibStack.appendKeyValue("email", _self.email);
        len = LibStack.appendKeyValue("mobile", _self.mobile);
        len = LibStack.appendKeyValue("departmentId", _self.departmentId);
        len = LibStack.appendKeyValue("accountStatus", _self.accountStatus);
        len = LibStack.appendKeyValue("passwordStatus", _self.passwordStatus);
        len = LibStack.appendKeyValue("deleteStatus", _self.deleteStatus);
        len = LibStack.appendKeyValue("uuid", _self.uuid);
        len = LibStack.appendKeyValue("publicKey", _self.publicKey);
        len = LibStack.appendKeyValue("cipherGroupKey", _self.cipherGroupKey);
        len = LibStack.appendKeyValue("createTime", _self.createTime);
        len = LibStack.appendKeyValue("updateTime", _self.updateTime);
        len = LibStack.appendKeyValue("loginTime", _self.loginTime);
        len = LibStack.appendKeyValue("roleIdList", _self.roleIdList.toJsonArray());
        len = LibStack.appendKeyValue("state", uint(_self.state));
        len = LibStack.appendKeyValue("status", _self.status);
        len = LibStack.appendKeyValue("ownerAddr", _self.ownerAddr);
        len = LibStack.appendKeyValue("certType", _self.certType);
        len = LibStack.appendKeyValue("remark", _self.remark);
        len = LibStack.appendKeyValue("icon", _self.icon);
        len = LibStack.append("}");
        _strjson = LibStack.popex(len);
    }

    function jsonParse(User storage _self, string _userJson) internal returns(bool) {
        if(bytes(_userJson).length == 0){
            return false;
        }
        _self.reset();
        LibJson.push(_userJson);
        _self.userAddr = _userJson.jsonRead("userAddr").toAddress();
        _self.name = _userJson.jsonRead("name");
        _self.mobile = _userJson.jsonRead("mobile");
        _self.account = _userJson.jsonRead("account");
        _self.email = _userJson.jsonRead("email");
        _self.departmentId = _userJson.jsonRead("departmentId");
        _self.passwordStatus = uint(_userJson.jsonRead("passwordStatus").toUint());
        _self.accountStatus  = uint(_userJson.jsonRead("accountStatus").toUint());
        _self.deleteStatus = uint(_userJson.jsonRead("deleteStatus").toUint());
        _self.status = uint(_userJson.jsonRead("status").toUint());
        _self.uuid = _userJson.jsonRead("uuid");
        _self.publicKey = _userJson.jsonRead("publicKey");
        _self.cipherGroupKey = _userJson.jsonRead("cipherGroupKey");
        _self.icon = _userJson.jsonRead("icon");
        _self.remark = _userJson.jsonRead("remark");
        _self.certType = uint(_userJson.jsonRead("certType").toUint());
        if(bytes(_userJson.jsonRead("roleIdList")).length > 0){
            _self.roleIdList.fromJsonArray(_userJson.jsonRead("roleIdList"));
        }
        _self.state = UserState.USER_INVALID;
        LibJson.pop();
        return true;
    }

    function fromJsonArray(User[] storage _self, string _json) internal returns(bool succ) {
        _self.length = 0;
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }

        while (true) {
            string memory key = "[".concat(_self.length.toString(), "]");
            if (_json.jsonKeyExistsEx(key) == uint(0))
            break;

            _self.length++;
            _self[_self.length-1].jsonParse(_json.jsonRead(key));
        }
        LibJson.pop();
        return true;
    }

    function toJsonArray(User[] storage _self) internal constant returns(string _json) {
        uint len = 0;
        len = LibStack.push("[");
        for (uint i = 0; i < _self.length; ++i) {
            if (i > 0)
            len = LibStack.append(",");
            len = LibStack.append(_self[i].toJson());
        }
        len = LibStack.append("]");
        _json = LibStack.popex(len);
    }

    function reset(User storage _self) internal {
        _self.userAddr = address(0);
        _self.name = "";
        _self.account = "";
        _self.email = "";
        _self.mobile = "";
        _self.departmentId = "";
        _self.accountStatus = 0;
        _self.passwordStatus = 0;
        _self.deleteStatus = 0;
        _self.uuid = "";
        _self.publicKey = "";
        _self.cipherGroupKey = "";
        _self.createTime = 0;
        _self.updateTime = 0;
        _self.loginTime = 0;
        delete _self.roleIdList;
        _self.state = UserState.USER_INVALID;
        _self.status = 1;
        _self.ownerAddr = address(0);
        _self.remark = "";
        _self.icon = "";
    }
}