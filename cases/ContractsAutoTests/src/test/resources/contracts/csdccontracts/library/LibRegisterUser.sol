pragma solidity ^0.4.12;

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibRegisterUser {

    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibRegisterUser for *;

    struct RegisterUser {
        string uuid;
        string userId;                      //用户ID
        address userAddr;                   //钱包地址
        string account;                     //用户账户
        string name;                        //用户名称
        string orgId;                       //所属组织ID
        string orgName;                     //组织名称
        uint certType;                      //账户类型  1 文件证书  2 U-key证书
        string mobile;                      //手机号
        string email;                       //用户邮箱
        string desc;                        //备注
        uint accountStatus;                 //账号状态 1待激活 2已激活 3已拒绝
        string publicKey;                   //用户公钥
        string cipherGroupKey;              //用户群私钥
    }

    function fromJson(RegisterUser storage _self, string _json) internal returns (bool succ) {
        _self.reset();
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }
        _self.uuid = _json.jsonRead("uuid");
        _self.userId = _json.jsonRead("userId");
        _self.userAddr = _json.jsonRead("userAddr").toAddress();
        _self.account = _json.jsonRead("account");
        _self.name = _json.jsonRead("name");
        _self.orgId = _json.jsonRead("orgId");
        _self.orgName = _json.jsonRead("orgName");
        _self.certType = _json.jsonRead("certType").toUint();
        _self.mobile = _json.jsonRead("mobile");
        _self.email = _json.jsonRead("email");
        _self.desc = _json.jsonRead("desc");
        _self.accountStatus = _json.jsonRead("accountStatus").toUint();
        if (_self.accountStatus == 0) _self.accountStatus == 1;
        _self.publicKey = _json.jsonRead("publicKey");
        _self.cipherGroupKey = _json.jsonRead("cipherGroupKey");
        LibJson.pop();
        return true;
    }

    function toJson(RegisterUser storage _self) internal constant returns (string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("uuid", _self.uuid);
        len = LibStack.appendKeyValue("userId", _self.userId);
        len = LibStack.appendKeyValue("userAddr", _self.userAddr);
        len = LibStack.appendKeyValue("account", _self.account);
        len = LibStack.appendKeyValue("name", _self.name);
        len = LibStack.appendKeyValue("orgId", _self.orgId);
        len = LibStack.appendKeyValue("orgName", _self.orgName);
        len = LibStack.appendKeyValue("certType", _self.certType);
        len = LibStack.appendKeyValue("mobile", _self.mobile);
        len = LibStack.appendKeyValue("email", _self.email);
        len = LibStack.appendKeyValue("desc", _self.desc);
        len = LibStack.appendKeyValue("accountStatus", _self.accountStatus);
        len = LibStack.appendKeyValue("publicKey", _self.publicKey);
        len = LibStack.appendKeyValue("cipherGroupKey", _self.cipherGroupKey);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJsonArray(RegisterUser[] storage _self, string _json) internal returns (bool succ) {
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
            _self[_self.length - 1].fromJson(_json.jsonRead(key));
        }
        LibJson.pop();
        return true;
    }

    function toJsonArray(RegisterUser[] storage _self) internal constant returns (string _json) {
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

    function update(RegisterUser storage _self, string _json) internal returns (bool succ) {
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }

        if (_json.jsonKeyExistsEx("uuid")!= uint(0))
        _self.uuid = _json.jsonRead("uuid");
        if (_json.jsonKeyExistsEx("userId")!= uint(0))
        _self.userId = _json.jsonRead("userId");
        if (_json.jsonKeyExistsEx("userAddr")!= uint(0))
        _self.userAddr = _json.jsonRead("userAddr").toAddress();
        if (_json.jsonKeyExistsEx("account")!= uint(0))
        _self.account = _json.jsonRead("account");
        if (_json.jsonKeyExistsEx("name")!= uint(0))
        _self.name = _json.jsonRead("name");
        if (_json.jsonKeyExistsEx("orgId")!= uint(0))
        _self.orgId = _json.jsonRead("orgId");
        if (_json.jsonKeyExistsEx("orgName")!= uint(0))
        _self.orgName = _json.jsonRead("orgName");
        if (_json.jsonKeyExistsEx("certType")!= uint(0))
        _self.certType = _json.jsonRead("certType").toUint();
        if (_json.jsonKeyExistsEx("mobile")!= uint(0))
        _self.mobile = _json.jsonRead("mobile");
        if (_json.jsonKeyExistsEx("email")!= uint(0))
        _self.email = _json.jsonRead("email");
        if (_json.jsonKeyExistsEx("desc")!= uint(0))
        _self.desc = _json.jsonRead("desc");
        if (_json.jsonKeyExistsEx("accountStatus")!= uint(0))
        _self.accountStatus = _json.jsonRead("accountStatus").toUint();
        if (_json.jsonKeyExistsEx("publicKey")!= uint(0))
        _self.publicKey = _json.jsonRead("publicKey");
        if (_json.jsonKeyExistsEx("cipherGroupKey")!= uint(0))
        _self.cipherGroupKey = _json.jsonRead("cipherGroupKey");
        LibJson.pop();
        return true;
    }

    function reset(RegisterUser storage _self) internal {
        delete _self.uuid;
        delete _self.userId;
        delete _self.userAddr;
        delete _self.account;
        delete _self.name;
        delete _self.orgId;
        delete _self.orgName;
        delete _self.certType;
        delete _self.mobile;
        delete _self.email;
        delete _self.desc;
        delete _self.accountStatus;
        delete _self.publicKey;
        delete _self.cipherGroupKey;
    }
}
