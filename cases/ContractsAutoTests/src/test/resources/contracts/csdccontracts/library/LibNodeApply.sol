pragma solidity ^0.4.12;
/**
 *
 * @file LibNodeApply.sol
 * @author Jungle
 * @time 2017-4-13 22:16:20
 * @desc Lib for node apply
 *
 */

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibJson.sol";
import "../utillib/LibStack.sol";

library LibNodeApply {

    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibNodeApply for *;

    struct ApplyNodeIP {
        string      ip;                 // ip地址
        int         uintIP;             // IP转为int的值
        uint        P2PPort;            // p2p端口
        uint        RPCPort;            // RPC端口
        uint        TPort;              // tomcat端口
        uint        _type;              // 内外网标识  0内网、1外网
    }

    struct ApplyNodeInfo {
        string          nodeId;             // 节点ID
        string          nodeName;           // 节点名称
        string          nodeShortName;      // 节点简称
        address         nodeAdmin;          // 节点管理员地址
        string          nodeDescription;    // 节点描述
        uint            state;              // 节点状态：0 无效 1 有效
        uint            _type;              // 1 记账节点 0 非记账节点
        string          deptId;             // 归属机构ID
        uint            deptLevel;          // 机构级别
        string          deptCN;             // 归属机构证书CN 
        string          pubkey;             // 节点公钥
        address         nodeAddress;        // 节点地址
        string          ip;                 // 节点P2P通信IP
        uint            port;               // 节点P2P通信端口
    }

    // 审核信息
    struct AuditData {
        string      applyId;             // node apply id
        string      parentId;            // department parent id
        uint        departmentLevel;     // departmentLevel
        uint        state;               // audit state 0 init ,1 wait ,2 success ,3 fail
        address     admin;               // the admin address of department 
        string      reason;              // the reason of autiting result
        string[]    roleIdList;          // the roleid
        string      cipherGroupKey;      // 用户群私钥
    }

    // 节点信息之 - 用户信息
    struct ApplyUser {
        address     userAddr;           // 用户钱包文件地址address
        string      name;               // 用户名称
        string      account;            // 用户账户名
        string      email;              // 用户邮箱
        uint        passwordStatus;     // 密码状态
        uint        accountStatus;      // 账户状态
        uint        deleteStatus;       // 
        string      mobile;             // 用户手机号
        string      departmentId;       // 归属部门ID（uuid）
        string      uuid;               // 唯一标识
        string      publicKey;          // 用户公钥  " " 
        string      cipherGroupKey;     // 密码组key " "
    } 

    // 节点申请之 - 部门信息
    struct ApplyDepartment {
        string id;                  // 机构Id [UUID由业务随机生成（保持唯一）]             
        string name;                // 机构名称
        uint departmentLevel;       // 机构层级（树形结构Level，0 开始）
        string parentId;            // 父机构Id
        string description;         // 机构描述
        string commonName;          // 证书CN
        string stateName;           // 省名称
        string countryName;         // 国家名称
        address admin;              // 部门管理员地址
        string orgaShortName;       // 组织简称
    }

    enum ApplyState {
        APPLY_INIT,
        APPLY_WAIT,
        APPLY_SUCCESS,
        APPLY_FAIL
    }

    // 申请信息记录
    struct NodeApply {
        string          id;                     // 申请ID（uuid）
        uint            applyTime;              // 申请时间戳
        uint            createTime;             // 创建时间
        uint            updateTime;             // 修改时间
        ApplyDepartment applyDepartment;        // 申请的部门信息
        ApplyUser       applyUser;              // 申请的用户信息
        bool            deleted;                // 标识当前对象是否有效
        uint            state;                  // 申请状态：0 初始化 1 等待审核 2 申请 3 申请失败
        string          reason;                 // 审核备注（通过或不通过原因）
        address         creator;                // 申请者地址
        address         auditUAddr;             // 审核人地址
        ApplyNodeInfo   applyNodeInfo;          // 申请节点信息
        ApplyNodeIP[]   applyNodeIPList;        // 申请节点的IP地址信息
    }

    function toJson(NodeApply storage _self) internal constant returns(string _json) {
        uint len = 0;
        len =  LibStack.push("{");
        len =  LibStack.appendKeyValue("id",_self.id);
        len =  LibStack.appendKeyValue("applyTime",_self.applyTime);
        len =  LibStack.appendKeyValue("createTime",_self.createTime);
        len =  LibStack.appendKeyValue("updateTime",_self.updateTime);
        len =  LibStack.appendKeyValue("applyDepartment",_self.applyDepartment.toJson());
        len =  LibStack.appendKeyValue("applyUser",_self.applyUser.toJson());
        len =  LibStack.appendKeyValue("state",_self.state);
        len =  LibStack.appendKeyValue("reason",_self.reason);
        len =  LibStack.appendKeyValue("creator",_self.creator);
        len =  LibStack.appendKeyValue("auditUAddr",_self.auditUAddr);
        len =  LibStack.appendKeyValue("applyNodeInfo",_self.applyNodeInfo.toJson());
        len =  LibStack.appendKeyValue("applyNodeIPList",_self.applyNodeIPList.toJsonArray());
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJson(NodeApply storage _self, string _json) internal returns(bool) {
        // internal : Internal visible only
        // constant : Diasllows modification of state
        if(bytes(_json).length == 0){
            return false;
        }
       LibJson.push(_json);
       _self.clear();
       _self.id = _json.jsonRead("id");
       _self.applyTime = _json.jsonRead("applyTime").toUint();
       _self.createTime = _json.jsonRead("createTime").toUint();
       _self.applyDepartment.fromJson(_json.jsonRead("applyDepartment"));
       _self.state = _json.jsonRead("state").toUint();
       _self.reason = _json.jsonRead("reason");
       _self.creator = _json.jsonRead("creator").toAddress();
       _self.auditUAddr = _json.jsonRead("auditUAddr").toAddress();
       _self.applyUser.fromJson(_json.jsonRead("applyUser"));
       _self.applyNodeInfo.fromJson(_json.jsonRead("applyNodeInfo"));
       if(bytes(_json.jsonRead("applyNodeIPList")).length > 0){
            _self.applyNodeIPList.fromJsonArray(_json.jsonRead("applyNodeIPList"));
        }
       LibJson.pop();
       return true;
    }

    function toJsonArray(NodeApply[] storage _self) internal constant returns(string _json){
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

    function fromJsonArray(NodeApply[] storage _self, string _json) internal returns(bool succ) {
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


    function clear(NodeApply storage _self) internal{
        _self.id = "";
        _self.applyTime = 0;
        _self.createTime = 0;
        _self.updateTime = 0;
        delete _self.applyDepartment;
        delete _self.applyUser; 
        _self.state = 0;
        _self.reason = "";
        _self.creator = address(0);
        _self.auditUAddr = address(0);
        delete _self.applyNodeInfo;
        _self.deleted = false;
        delete _self.applyNodeIPList;
    }

     //////////////////////////子结构——————ApplyDepartment//////////////////////////

     function toJson(ApplyDepartment storage _self) internal constant returns(string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("id",_self.id);
        len = LibStack.appendKeyValue("name",_self.name);
        len = LibStack.appendKeyValue("departmentLevel",_self.departmentLevel);
        len = LibStack.appendKeyValue("parentId",_self.parentId);
        len = LibStack.appendKeyValue("description",_self.description);
        len = LibStack.appendKeyValue("commonName",_self.commonName);
        len = LibStack.appendKeyValue("stateName",_self.stateName);
        len = LibStack.appendKeyValue("countryName",_self.countryName);
        len = LibStack.appendKeyValue("admin",_self.admin);
        len = LibStack.appendKeyValue("orgaShortName",_self.orgaShortName);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
     }


     function fromJson(ApplyDepartment storage _self, string _json) internal  returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        _self.clear();
        _self.id = _json.jsonRead("id");
        _self.name = _json.jsonRead("name");
        _self.departmentLevel = _json.jsonRead("departmentLevel").toUint();
        _self.parentId = _json.jsonRead("parentId");
        _self.description = _json.jsonRead("description");
        _self.commonName = _json.jsonRead("commonName");
        _self.stateName = _json.jsonRead("stateName");
        _self.countryName = _json.jsonRead("countryName");
        _self.admin = _json.jsonRead("admin").toAddress();
        _self.orgaShortName = _json.jsonRead("orgaShortName");
        LibJson.pop();
        return true;
     }

    function clear(ApplyDepartment storage _self) internal{
        _self.id = "";
        _self.name = "";
        _self.departmentLevel = 0;
        _self.parentId = "";
        _self.description = "";
        _self.commonName = "";
        _self.stateName = "";
        _self.countryName = "";
        _self.orgaShortName = "";
        _self.admin = address(0);
    }   

     //////////////////////////子结构——————ApplyUser//////////////////////////

     function toJson(ApplyUser storage _self) internal constant returns(string _json){
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("userAddr",_self.userAddr);
        len = LibStack.appendKeyValue("name",_self.name);
        len = LibStack.appendKeyValue("account",_self.account);
        len = LibStack.appendKeyValue("email",_self.email);
        len = LibStack.appendKeyValue("accountStatus",_self.accountStatus);
        len = LibStack.appendKeyValue("deleteStatus",_self.deleteStatus);
        len = LibStack.appendKeyValue("mobile",_self.mobile);
        len = LibStack.appendKeyValue("departmentId",_self.departmentId);
        len = LibStack.appendKeyValue("publicKey",_self.publicKey);
        len = LibStack.appendKeyValue("cipherGroupKey",_self.cipherGroupKey);
        //len = LibStack.appendKeyValue("uuid",_self.uuid);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
     }

     function fromJson(ApplyUser storage _self, string _json) internal  returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        _self.clear();
        _self.userAddr = _json.jsonRead("userAddr").toAddress();
        _self.name = _json.jsonRead("name");
        _self.account = _json.jsonRead("account");
        _self.email = _json.jsonRead("email");
        _self.accountStatus = _json.jsonRead("accountStatus").toUint();
        _self.deleteStatus = _json.jsonRead("deleteStatus").toUint();
        _self.mobile = _json.jsonRead("mobile");
        _self.departmentId = _json.jsonRead("departmentId");
        _self.publicKey = _json.jsonRead("publicKey");
        _self.cipherGroupKey = _json.jsonRead("cipherGroupKey");
        LibJson.pop();
        return true;
     }


    function clear(ApplyUser storage _self) internal {
        _self.userAddr = address(0);
        _self.name = "";
        _self.account = "";
        _self.email = "";
        _self.accountStatus = 0;
        _self.deleteStatus = 0;
        _self.mobile = "";
        _self.departmentId = "";
        _self.publicKey = "";
        _self.cipherGroupKey = "";
        _self.uuid = "";
    }

    //////////////////////////子结构——————ApplyNodeInfo//////////////////////////

    function toJson(ApplyNodeInfo storage _self) internal constant returns(string _json){
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("nodeId",_self.nodeId);
        len = LibStack.appendKeyValue("nodeName",_self.nodeName);
        len = LibStack.appendKeyValue("nodeShortName",_self.nodeShortName);
        len = LibStack.appendKeyValue("nodeAdmin",_self.nodeAdmin);
        len = LibStack.appendKeyValue("nodeDescription",_self.nodeDescription);
        len = LibStack.appendKeyValue("type",_self._type);
        len = LibStack.appendKeyValue("deptId",_self.deptId);
        len = LibStack.appendKeyValue("deptLevel",_self.deptLevel);
        len = LibStack.appendKeyValue("deptCN",_self.deptCN);
        len = LibStack.appendKeyValue("pubkey",_self.pubkey);
        len = LibStack.appendKeyValue("nodeAddress",_self.nodeAddress);
        len = LibStack.appendKeyValue("ip",_self.ip);
        len = LibStack.appendKeyValue("port",_self.port);
        len = LibStack.append("}"); 
        _json = LibStack.popex(len);
    }

    function fromJson(ApplyNodeInfo storage _self, string _json) internal returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }        
        LibJson.push(_json);
        _self.clear();
        _self.nodeId = _json.jsonRead("nodeId");
        _self.nodeName = _json.jsonRead("nodeName");
        _self.nodeShortName = _json.jsonRead("nodeShortName");
        _self.nodeAdmin = _json.jsonRead("nodeAdmin").toAddress();
        _self.nodeDescription = _json.jsonRead("nodeDescription");
        _self._type = _json.jsonRead("type").toUint();
        _self.deptId = _json.jsonRead("deptId");
        _self.deptLevel = _json.jsonRead("deptLevel").toUint();
        _self.deptCN = _json.jsonRead("deptCN");
        _self.pubkey = _json.jsonRead("pubkey");
        _self.nodeAddress = _json.jsonRead("nodeAddress").toAddress();
        _self.ip = _json.jsonRead("ip");
        _self.port = _json.jsonRead("port").toUint();
        LibJson.pop();
        return true;
    }

    function clear(ApplyNodeInfo storage _self)internal{
        _self.nodeId = "";
        _self.nodeName ="";
        _self.nodeShortName = "";
        _self.nodeAdmin = address(0);
        _self.nodeDescription = "";
        _self._type = 0;
        _self.deptId = "";
        _self.deptLevel = 0;
        _self.deptCN = "";
        _self.pubkey = "";
        _self.nodeAddress = address(0);
        _self.ip = "";
        _self.port = 0;
    }

    //////////////////////////子结构——————applyNodeIPList//////////////////////////

    function toJson(ApplyNodeIP storage _self)internal constant returns(string _json){
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ip",_self.ip);
        len = LibStack.appendKeyValue("uintIP",_self.uintIP);
        len = LibStack.appendKeyValue("P2PPort",_self.P2PPort);
        len = LibStack.appendKeyValue("RPCPort",_self.RPCPort);
        len = LibStack.appendKeyValue("TPort",_self.TPort);
        len = LibStack.appendKeyValue("type",_self._type);    
        len = LibStack.append("}");
        _json = LibStack.popex(len);  
    }


     function fromJson(ApplyNodeIP storage _self, string _json) internal  returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }        
        LibJson.push(_json);
        _self.clear();
        _self.ip = _json.jsonRead("ip");
        _self.uintIP = _json.jsonRead("uintIP").toInt();
        _self.P2PPort = _json.jsonRead("P2PPort").toUint();
        _self.RPCPort = _json.jsonRead("RPCPort").toUint();
        _self.TPort = _json.jsonRead("TPort").toUint();
        _self._type = _json.jsonRead("type").toUint();
        LibJson.pop();
        return true;
     }

    function toJsonArray(ApplyNodeIP[] storage _self) internal constant returns(string _json){
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

    function fromJsonArray(ApplyNodeIP[] storage _self, string _json) internal returns(bool succ) {
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


    function clear(ApplyNodeIP storage _self)internal{
        _self.ip = "";
        _self.uintIP =0;
        _self.P2PPort = 0;
        _self.RPCPort = 0;
        _self.TPort = 0;
        _self._type = 0;
    }

      ////////////////////////结构——AuditData////////////////////////////
    function toJsonForAuditData(AuditData storage _self) internal constant returns(string _json){
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("applyId",_self.applyId);
        len = LibStack.appendKeyValue("parentId",_self.parentId);
        len = LibStack.appendKeyValue("departmentLevel",_self.departmentLevel);
        len = LibStack.appendKeyValue("state",_self.state);
        len = LibStack.appendKeyValue("reason",_self.reason);
        len = LibStack.appendKeyValue("roleIdList", _self.roleIdList.toJsonArray()); 
        len = LibStack.appendKeyValue("cipherGroupKey",_self.cipherGroupKey); 
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJsonForAuditData(AuditData storage _self, string _json)internal returns(bool){
        if(bytes(_json).length == 0){
            return false;
        }        
        LibJson.push(_json);
        _self.clear();
        _self.applyId = _json.jsonRead("applyId");
        _self.parentId = _json.jsonRead("parentId");
        _self.departmentLevel = _json.jsonRead("departmentLevel").toUint();
        _self.state = _json.jsonRead("state").toUint();
        _self.admin = _json.jsonRead("admin").toAddress();
        _self.reason = _json.jsonRead("reason");
        if(bytes(_json.jsonRead("roleIdList")).length > 0){
            _self.roleIdList.fromJsonArray(_json.jsonRead("roleIdList"));
        }
        _self.cipherGroupKey = _json.jsonRead("cipherGroupKey");
        LibJson.pop();
        return true;
   }

    function toJsonArray(AuditData[] storage _self) internal constant returns(string _json){
        uint len = 0;
        len = LibStack.push("[");
        for (uint i=0; i<_self.length; ++i) {
            if (i > 0)
                len = LibStack.append(",");
            len = LibStack.append(_self[i].toJsonForAuditData());
        }
        len = LibStack.append("]");
        _json = LibStack.popex(len);
    }

    function fromJsonArray(AuditData[] storage _self, string _json) internal returns(bool succ) {
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
            _self[_self.length-1].fromJsonForAuditData(_json.jsonRead(key));
        }
        LibJson.pop();
        return true;
    }

    function clear(AuditData storage _self) internal {
        _self.applyId = "";
        _self.parentId = "";
        _self.departmentLevel = 0;
        _self.state = 0;
        _self.admin = address(0);
        _self.reason = "";
        delete _self.roleIdList;
        _self.cipherGroupKey = "";
    }
}

