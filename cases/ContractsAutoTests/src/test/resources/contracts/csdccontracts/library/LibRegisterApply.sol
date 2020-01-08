pragma solidity ^0.4.12;

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibJson.sol";
import "../utillib/LibLog.sol";
import "../utillib/LibStack.sol";
import "../library/LibRegisterUser.sol";
import "../library/LibAuditRecord.sol";

library LibRegisterApply {

    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibLog for *;
    using LibRegisterUser for *;
    using LibAuditRecord for *;
    using LibRegisterApply for *;

    struct RegisterApply {
        string applyId;                             //申请ID
        uint256 applyTime;                          //申请时间
        LibRegisterUser.RegisterUser registerUser;  //申请的用户信息
        LibAuditRecord.AuditRecord[] auditList;     //审核记录
        uint auditStatus;                           //审核状态【默认1】1 待审核  2 已同意 3 已拒绝
        uint256 createTime;                         //创建时间
        uint256 updateTime;                         //最后更新时间
    }

    function fromJson(RegisterApply storage _self, string _json) internal returns (bool succ) {
        _self.reset();
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }
        _self.applyId = _json.jsonRead("applyId");
        _self.applyTime = uint256(_json.jsonRead("applyTime").toUint());
        _self.registerUser.fromJson(_json.jsonRead("registerUser"));
        if(bytes(_json.jsonRead("auditList")).length > 0){
        _self.auditList.fromJsonArray(_json.jsonRead("auditList"));
        }
        _self.auditStatus = _json.jsonRead("auditStatus").toUint();
        _self.createTime = uint256(_json.jsonRead("createTime").toUint());
        _self.updateTime = uint256(_json.jsonRead("updateTime").toUint());
        LibJson.pop();
        return true;
    }

    function toJson(RegisterApply storage _self) internal constant returns (string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("applyId", _self.applyId);
        len = LibStack.appendKeyValue("applyTime", uint256(_self.applyTime));
        len = LibStack.appendKeyValue("registerUser", _self.registerUser.toJson());
        len = LibStack.appendKeyValue("auditList", _self.auditList.toJsonArray());
        len = LibStack.appendKeyValue("auditStatus", _self.auditStatus);
        len = LibStack.appendKeyValue("createTime", uint256(_self.createTime));
        len = LibStack.appendKeyValue("updateTime", uint256(_self.updateTime));
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJsonArray(RegisterApply[] storage _self, string _json) internal returns (bool succ) {
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

    function toJsonArray(RegisterApply[] storage _self) internal constant returns (string _json) {
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

    function update(RegisterApply storage _self, string _json) internal returns (bool succ) {
        LibJson.push(_json);
        if (!_json.isJson()){
            LibJson.pop();
            return false;
        }

        if (_json.jsonKeyExistsEx("applyId")!= uint(0))
        _self.applyId = _json.jsonRead("applyId");
        if (_json.jsonKeyExistsEx("applyTime")!= uint(0))
        _self.applyTime = uint256(_json.jsonRead("applyTime").toUint());
        if (_json.jsonKeyExistsEx("registerUser")!= uint(0))
        _self.registerUser.fromJson(_json.jsonRead("registerUser"));
        if (_json.jsonKeyExistsEx("auditList")!= uint(0))
        _self.auditList.fromJsonArray(_json.jsonRead("auditList"));
        if (_json.jsonKeyExistsEx("auditStatus")!= uint(0))
        _self.auditStatus = _json.jsonRead("auditStatus").toUint();
        _self.updateTime = uint256(now * 1000);
        LibJson.pop();
        return true;
    }

    function reset(RegisterApply storage _self) internal {
        delete _self.applyId;
        delete _self.applyTime;
        delete _self.registerUser;
        _self.auditList.length = 0;
        delete _self.auditStatus;
        delete _self.createTime;
        delete _self.updateTime;
    }

}
