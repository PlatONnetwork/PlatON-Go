pragma solidity ^0.4.12;
/**
* @file LibNodeInfo.sol
* @author Jungle
* @time 2017-5-20 09:42:28
* @desc the defination of Department library
*/

import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibStack.sol";
import "../utillib/LibJson.sol";

library LibNodeInfo {
    
    using LibInt for *;
    using LibString for *;
    using LibNodeInfo for *;
    using LibJson for *;

    // 入口信息
    struct NodeNAT {
        string      ip;                     // 入口IP地址（外网）
        string      pubkey;
        address     nodeAddress;
        uint        activated;              // 0 禁用 1 启用                 
        uint        p2pPort;                // p2p端口
        uint        tPort;
        uint        rpcPort;
    }

    // 出口信息
    struct NodeLAN {
        string      maskIP;                 // 掩码IP
        int         mastInt;                // 掩码IP对应整数
        string      goByIP;                 // 掩码基础IP
        int         goByInt;                // 基础IP对应整数
        string      startIP;                // IP段的起始IP
        int         startInt;               
        string      endIP;                  // IP段对应结束IP
        int         endInt;
    }

    struct NodeInfo {
        string          nodeId;             // 节点ID
        string          nodeName;           // 节点名称
        string          nodeShortName;      // 节点简称             *
        address         nodeAdmin;          // 节点管理员地址      
        string          nodeDescription;    // 节点描述
        bool            deleted;            // 删除标记 true 已删除 false 未删除
        uint            state;              // 节点状态：0 无效 1 有效
        uint            _type;              // 1 记账节点 0 非记账节点
        string          deptId;             // 归属机构ID
        uint            deptLevel;          // 机构级别             - 
        string          deptCN;             // 归属机构证书CN 
        uint            createTime; 
        uint            updateTime;
        NodeLAN         nodeLAN;            // 局域网信息
        NodeNAT         nodeNAT;            // 公网信息
        uint            disabled;           // 隔离状态 0 禁用 1 启动
    }

    function toJson(NodeInfo storage _self) internal constant returns(string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("nodeId",_self.nodeId);
        len = LibStack.appendKeyValue("nodeName",_self.nodeName);
        len = LibStack.appendKeyValue("nodeShortName",_self.nodeShortName);
        len = LibStack.appendKeyValue("nodeAdmin",_self.nodeAdmin);
        len = LibStack.appendKeyValue("nodeDescription",_self.nodeDescription);
        len = LibStack.appendKeyValue("state",_self.state);
        len = LibStack.appendKeyValue("type",_self._type);
        len = LibStack.appendKeyValue("deptId",_self.deptId);
        len = LibStack.appendKeyValue("deptLevel",_self.deptLevel);
        len = LibStack.appendKeyValue("deptCN",_self.deptCN);
        len = LibStack.appendKeyValue("createTime",_self.createTime);
        len = LibStack.appendKeyValue("updateTime",_self.updateTime);
        len = LibStack.appendKeyValue("nodeLAN",_self.nodeLAN.toJson());
        len = LibStack.appendKeyValue("nodeNAT",_self.nodeNAT.toJson());
        len = LibStack.appendKeyValue("disabled",_self.disabled);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJson(NodeInfo storage _self, string _json) internal constant returns(bool succ) {
        _self.clear();
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        _self.nodeId = _json.jsonRead("nodeId");
        _self.nodeName = _json.jsonRead("nodeName");
        _self.nodeShortName = _json.jsonRead("nodeShortName");
        _self.nodeAdmin = _json.jsonRead("nodeAdmin").toAddress();
        _self.nodeDescription = _json.jsonRead("nodeDescription");
        _self.state = _json.jsonRead("state").toUint();
        _self._type = _json.jsonRead("type").toUint();
        _self.deptId = _json.jsonRead("deptId");
        _self.deptLevel = _json.jsonRead("deptLevel").toUint();
        _self.deptCN = _json.jsonRead("deptCN");
        _self.nodeLAN.fromJson(_json.jsonRead("nodeLAN"));
        _self.nodeNAT.fromJson(_json.jsonRead("nodeNAT"));
        _self.disabled = _json.jsonRead("disabled").toUint();
        LibJson.pop();
        if (bytes(_self.nodeId).length == 0) {
            return false;
        }
        return true;
    }

    function clear(NodeInfo storage _self) internal {
        _self.nodeId = "";
        _self.nodeName = "";       
        _self.nodeShortName = "";  
        _self.nodeAdmin = 0;      
        _self.nodeDescription = "";
        _self.deleted = false;        
        _self.state = 1;          
        _self._type = 0;          
        _self.deptId = "";
        _self.deptLevel = 0;    
        _self.deptCN = "";     
        delete _self.nodeLAN;
        delete _self.nodeNAT;
        _self.disabled = 1;      
    }


    //////////////////////////子结构操作函数//////////////////////////
    function toJson(NodeNAT storage _self) internal constant returns (string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ip",_self.ip);
        len = LibStack.appendKeyValue("pubkey",_self.pubkey);
        len = LibStack.appendKeyValue("nodeAddress",_self.nodeAddress);
        len = LibStack.appendKeyValue("activated",_self.activated);
        len = LibStack.appendKeyValue("p2pPort",_self.p2pPort);
        len = LibStack.appendKeyValue("tPort",_self.tPort);
        len = LibStack.appendKeyValue("rpcPort",_self.rpcPort);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJson(NodeNAT storage _self,string _json) internal returns (bool) {
        _self.reset();
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        _self.ip = _json.jsonRead("ip");
        _self.pubkey = _json.jsonRead("pubkey");
        _self.nodeAddress = _json.jsonRead("nodeAddress").toAddress();
        _self.activated = _json.jsonRead("activated").toUint();
        _self.p2pPort = _json.jsonRead("p2pPort").toUint();
        _self.tPort = _json.jsonRead("tPort").toUint();
        _self.rpcPort = _json.jsonRead("rpcPort").toUint();
        LibJson.pop();
        return true;
    }

    function toJson(NodeLAN storage _self) internal constant returns (string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("maskIP",_self.maskIP);
        len = LibStack.appendKeyValue("mastInt",_self.mastInt);
        len = LibStack.appendKeyValue("goByIP",_self.goByIP);
        len = LibStack.appendKeyValue("goByInt",_self.goByInt);
        len = LibStack.appendKeyValue("startIP",_self.startIP);
        len = LibStack.appendKeyValue("startInt",_self.startInt);
        len = LibStack.appendKeyValue("endIP",_self.endIP);
        len = LibStack.appendKeyValue("endInt",_self.endInt);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    function fromJson(NodeLAN storage _self,string _json) internal returns (bool) {
        _self.reset();
        if(bytes(_json).length == 0){
            return false;
        }
        LibJson.push(_json);
        _self.maskIP = _json.jsonRead("maskIP");
        _self.mastInt = _json.jsonRead("mastInt").toInt();
        _self.goByIP = _json.jsonRead("goByIP");
        _self.goByInt = _json.jsonRead("goByInt").toInt();
        _self.startIP = _json.jsonRead("startIP");
        _self.startInt = _json.jsonRead("startInt").toInt();
        _self.endIP = _json.jsonRead("endIP");
        _self.endInt = _json.jsonRead("endInt").toInt();
        LibJson.pop();
        return true;
    }

    function reset(NodeLAN storage _self) internal{
         _self.maskIP = "";
         _self.mastInt = 0;
         _self.goByIP = "";
         _self.goByInt = 0;
         _self.startIP = "";
         _self.startInt = 0;
         _self.endIP = "";
         _self.endInt = 0;
    }

    function reset(NodeNAT storage _self) internal{
         _self.ip = "";
         _self.pubkey = "";
         _self.nodeAddress = address(0);
         _self.activated = 0;
         _self.p2pPort = 0;
         _self.tPort = 0;
         _self.rpcPort = 0;
    }
}
