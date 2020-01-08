pragma solidity ^0.4.12;

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibAuditUser {
    
    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibAuditUser for *;

    struct AuditUser {
        string applyId;                 //申请记录ID
        string auditorId;               //审核人ID
        string auditorName;             //审核人名称
        string auditorMobile;           //审核人电话
        uint256 auditorTime;            //审核时间
        int auditResult;                //审核结果  1 同意  2 不同意
        string auditDesc;               //审核建议
        string cipherGroupKey;          //用户群私钥
    }

    function fromJson(AuditUser storage _self, string _json) internal returns(bool succ) {
        if(bytes(_json).length == 0){
            return false;
        }
        _self.reset();
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }

        _self.applyId = _json.jsonRead("applyId");
        _self.auditorId = _json.jsonRead("auditorId");
        _self.auditorName = _json.jsonRead("auditorName");
        _self.auditorMobile = _json.jsonRead("auditorMobile");
        _self.auditorTime = uint256(_json.jsonRead("auditorTime").toUint());
        _self.auditResult = _json.jsonRead("auditResult").toInt();
        _self.auditDesc = _json.jsonRead("auditDesc");
        _self.cipherGroupKey = _json.jsonRead("cipherGroupKey");
        LibJson.pop();
        return true;
    }

    function toJson(AuditUser storage _self) internal constant returns (string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("applyId", _self.applyId);
        len = LibStack.appendKeyValue("auditorId", _self.auditorId);
        len = LibStack.appendKeyValue("auditorName", _self.auditorName);
        len = LibStack.appendKeyValue("auditorMobile", _self.auditorMobile);
        len = LibStack.appendKeyValue("auditorTime", uint256(_self.auditorTime));
        len = LibStack.appendKeyValue("auditResult", _self.auditResult);
        len = LibStack.appendKeyValue("auditDesc", _self.auditDesc);
        len = LibStack.appendKeyValue("cipherGroupKey", _self.cipherGroupKey);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJsonArray(AuditUser[] storage _self, string _json) internal returns(bool succ) {
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
            _self[_self.length-1].fromJson(_json.jsonRead(key));
        }
        LibJson.pop();
        return true;
    }

    function toJsonArray(AuditUser[] storage _self) internal constant returns(string _json) {
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

    function update(AuditUser storage _self, string _json) internal returns(bool succ) {
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }

        if (_json.jsonKeyExistsEx("applyId")!= uint(0))
        _self.applyId = _json.jsonRead("applyId");
        if (_json.jsonKeyExistsEx("auditorId")!= uint(0))
        _self.auditorId = _json.jsonRead("auditorId");
        if (_json.jsonKeyExistsEx("auditorName")!= uint(0))
        _self.auditorName = _json.jsonRead("auditorName");
        if (_json.jsonKeyExistsEx("auditorMobile")!= uint(0))
        _self.auditorMobile = _json.jsonRead("auditorMobile");
        if (_json.jsonKeyExistsEx("auditorTime")!= uint(0))
        _self.auditorTime = uint256(_json.jsonRead("auditorTime").toUint());
        if (_json.jsonKeyExistsEx("auditResult")!= uint(0))
        _self.auditResult = _json.jsonRead("auditResult").toInt();
        if (_json.jsonKeyExistsEx("auditDesc")!= uint(0))
        _self.auditDesc = _json.jsonRead("auditDesc");
        if (_json.jsonKeyExistsEx("cipherGroupKey")!= uint(0))
        _self.cipherGroupKey = _json.jsonRead("cipherGroupKey");
        LibJson.pop();
        return true;
    }

    function reset(AuditUser storage _self) internal {
        delete _self.applyId;
        delete _self.auditorId;
        delete _self.auditorName;
        delete _self.auditorMobile;
        delete _self.auditorTime;
        delete _self.auditResult;
        delete _self.auditDesc;
        delete _self.cipherGroupKey;
    }
}
