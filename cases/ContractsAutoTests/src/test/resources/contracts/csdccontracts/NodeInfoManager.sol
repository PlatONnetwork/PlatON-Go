pragma solidity ^0.4.12;
/**
* @file NodeInfoManager.sol
* @author Jungle
* @time 2017-5-20 16:06:53
* @desc
*/

import "./library/LibNodeInfo.sol";
import "./sysbase/OwnerNamed.sol";
import "./interfaces/IUserManager.sol";
import "./interfaces/IDepartmentManager.sol";
import "./interfaces/INodeInfoManager.sol";

contract NodeInfoManager is OwnerNamed, INodeInfoManager {

    using LibNodeInfo for *;
    using LibString for *;
    using LibInt for *;
    using LibJson for *;
    
    LibNodeInfo.NodeInfo[] nodeInfos;
    LibNodeInfo.NodeInfo[] tmpNodeInfos;
    LibNodeInfo.NodeInfo tmpNodeInfo;

    uint revision;

    enum NodeInfoError {
        NO_ERROR,
        BAD_PARAMETER,
        NO_PERMISSION,
        ID_NOT_EXISTS,
        ADMIN_NOT_MEMBER,
        ID_CONFLICTED,
        PUBKEY_CONFLICTED
    }

    event Notify(uint _errno, string _info);

    function NodeInfoManager() {
        revision = 0;
        register("SystemModuleManager","0.0.1.0","NodeInfoManager", "0.0.1.0");
    }

    function saveNodeInfo(string _pubkey, string _nodeId) constant private returns (uint _ret) {
        log("NodeInfoManager.sol", "saveNodeInfo");
        log(_pubkey, _nodeId);
        _ret = writedb("nodeInfo|add", _pubkey, _nodeId);
        if (0 != _ret)
            log("NodeInfoManager.sol", "saveNodeInfo failed.");
        else
            log("NodeInfoManager.sol", "saveNodeInfo success.");
    }

    /**
    * @dev insert an node info
    * @param _json The json string described the object
    * @return errno , 0 for true
    */
    function insert(string _json) public returns(uint) {
        log("insert NodeInfo", "NodeInfoManager");
        uint prefix_error = 95270;
        // Decode json
        if (!tmpNodeInfo.fromJson(_json)) {
            log("json invalid", "NodeInfoManager"); 
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "json invalid");
            return errno;
        }
        if (getIndexById(tmpNodeInfo.nodeId) != uint(-1)) {
            log("nodeId conflicted", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.ID_CONFLICTED);
            Notify(errno, "nodeId conflicted");
            return errno;
        }
        if (bytes(tmpNodeInfo.nodeName).length == 0) {
            log("nodeName is empty", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "nodeName is empty");
            return errno;
        }
        if (bytes(tmpNodeInfo.deptId).length == 0) {
            log("deptId is empty", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "deptId is empty");
            return errno;
        }
        if (bytes(tmpNodeInfo.deptCN).length == 0) {
            log("deptCN is empty", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "deptCN is empty");
            return errno;
        }
        if (tmpNodeInfo.deptLevel != uint(1)) {
            log("deptLevel invalid ..", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "deptLevel invalid");
            return errno;
        }
        // check deptId and deptCN 
        if(__departmentExists(tmpNodeInfo.deptId) == 0) {
            log("deptId not exists in DepartmentManager", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "deptId not exists in DepartmentManager");
            return errno;
        }
        if(__departmentExistsByCN(tmpNodeInfo.deptCN) == 0) {
            log("deptCN not exists in DepartmentManager", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "deptCN not exists in DepartmentManager");
            return errno;
        }
        if (getIndexByNodeName(tmpNodeInfo.nodeName) != uint(-1)) {
            log("nodeName conflicted", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "nodeName conflicted");
            return errno;
        }
        if (tmpNodeInfo.nodeAdmin != 0) {
            if (!__getUserDepartmentId(tmpNodeInfo.nodeAdmin).equals(tmpNodeInfo.deptId)) {
                log("admin not member of department", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.ADMIN_NOT_MEMBER);
                Notify(errno, "admin is not a member of current node`department");
                return errno;
            }
        }

        // check pubkey is valid
        /* if (bytes(tmpNodeInfo.pubkey).length == 0 || bytes(tmpNodeInfo.pubkey).length > 128
                    || bytes(tmpNodeInfo.ip).length == 0
                    || tmpNodeInfo.port == 0 || tmpNodeInfo.port > 0xFFFF) {
                    log("node pubkey,ip,port invalid", "NodeInfoManager");
                    errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
                    Notify(errno, "node pubkey or ip or port format invalid");
                    return errno;
        }
        if(bytes(tmpNodeInfo.pubkey).length < 128 ) {
            uint lessCount = 128 - bytes(tmpNodeInfo.pubkey).length;
            for (uint j = 0; j < lessCount; ++j) {
                tmpNodeInfo.pubkey = "0".concat(tmpNodeInfo.pubkey);                
            }
        } */
        // check pubkey is exists
        /* uint pubkeyIndex = getIndexByPubkey(tmpNodeInfo.pubkey);
        if(pubkeyIndex != uint(-1)) {               
            log("node pubkey repeat..", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "node pubkey repeat..");
            return errno;
        } */
        // update struct 
        if(bytes(tmpNodeInfo.nodeNAT.ip).length == 0) {
            log("insert->node ip is empty");
            errno = prefix_error + uint256(NodeInfoError.BAD_PARAMETER);
            Notify(errno,"node ip is empty");
            return errno;
        }
        if (bytes(tmpNodeInfo.nodeNAT.pubkey).length == 0 
                    || tmpNodeInfo.nodeNAT.nodeAddress == 0
                    || tmpNodeInfo.nodeNAT.p2pPort == 0 
                    || tmpNodeInfo.nodeNAT.p2pPort > 0xFFFF
                    || tmpNodeInfo.nodeNAT.tPort == 0
                    || tmpNodeInfo.nodeNAT.tPort > 0xFFFF
                    || tmpNodeInfo.nodeNAT.rpcPort == 0
                    || tmpNodeInfo.nodeNAT.rpcPort > 0xFFFF) {
            log("nodeIP pubkey,nodeAddress,ip,P2PPort,RPCPort,TPort invalid", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "nodeIP ip,P2PPort,RPCPort,TPort invalid");
            return errno;
        }
        // check pubkey length
        if (bytes(tmpNodeInfo.nodeNAT.pubkey).length < 128) {
            uint lessCount = 128 - bytes(tmpNodeInfo.nodeNAT.pubkey).length;
            for (uint j=0; j < lessCount; ++j) {
                tmpNodeInfo.nodeNAT.pubkey = "0".concat(tmpNodeInfo.nodeNAT.pubkey);
            }
        }
        // check pubkey is conflict
        for(uint i = 0 ; i < nodeInfos.length; ++i) {
            if(nodeInfos[i].deleted) {
                continue;
            }
            if (tmpNodeInfo.nodeNAT.pubkey.equalsNoCase(nodeInfos[i].nodeNAT.pubkey)) {
                log("pubkey conflict other node", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.PUBKEY_CONFLICTED);
                Notify(errno, "pubkey conflict other node");
                return errno;
            }
        }
        /* for(uint i = 0 ; i < tmpNodeInfo.nodeIPList.length ; ++i) {

            if(bytes(tmpNodeInfo.nodeIPList[i].ip).length == 0){
                log("insert->node ip is empty");
                errno = prefix_error + uint256(NodeInfoError.BAD_PARAMETER);
                Notify(errno,"node ip is empty");
                return errno;
            }
            if(tmpNodeInfo.nodeIPList[i]._type == uint(1) || tmpNodeInfo.nodeIPList[i]._type == uint(0)) {
                if (bytes(tmpNodeInfo.nodeIPList[i].ip).length == 0
                        || tmpNodeInfo.nodeIPList[i].P2PPort == 0 
                        || tmpNodeInfo.nodeIPList[i].P2PPort > 0xFFFF
                        || tmpNodeInfo.nodeIPList[i].RPCPort == 0 
                        || tmpNodeInfo.nodeIPList[i].RPCPort > 0xFFFF
                        || tmpNodeInfo.nodeIPList[i].TPort == 0 
                        || tmpNodeInfo.nodeIPList[i].TPort > 0xFFFF) {
                    log("nodeIP ip,P2PPort,RPCPort,TPort invalid", "NodeInfoManager");
                    errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
                    Notify(errno, "nodeIP ip,P2PPort,RPCPort,TPort invalid");
                    return errno;
                }
            }
        } */

        if (tx.origin != owner) {
            log("sender not owner: ", tx.origin);
            /* if (__checkWritePermission(tx.origin, tmpNodeInfo.deptId) == 0
                && checkWritePermission(tx.origin,tmpNodeInfo.nodeId) == 0) {
                log("user not dept admin ..", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.NO_PERMISSION);
                Notify(errno, "user is not the admin of department or parent department");
                return errno;
            } */
        }

        tmpNodeInfo.createTime = now * 1000;
        tmpNodeInfo.updateTime = now * 1000;
        if(tx.origin != owner) {
            tmpNodeInfo.nodeNAT.activated = 0;
        }
        //tmpNodeInfo.createAddr = tx.origin;

        // Insert into department list
        nodeInfos.push(tmpNodeInfo);
        errno = uint(NodeInfoError.NO_ERROR);
        revision++;

        saveNodeInfo(tmpNodeInfo.nodeNAT.pubkey, tmpNodeInfo.nodeId);

        log("insert node succ", "NodeInfoManager");
        Notify(errno, "insert nodeInfo succ");
        return errno;
    }

    function update(string _json) public returns(uint){
        log("update action", "NodeInfoManager");
        uint prefix_error = 95270;
        if (!tmpNodeInfo.fromJson(_json)) {
            log("json invalid", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "json invalid");
            return errno;
        }

        uint index = getIndexById(tmpNodeInfo.nodeId);
        if (index == uint(-1)) {
            log("node not exists", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.ID_NOT_EXISTS);
            Notify(errno, "node id dose not exists");
            return errno;
        }

        // ====================== begin check =====================
        if (_json.keyExists("nodeName")) {
            if (bytes(tmpNodeInfo.nodeName).length == 0) {
                log("nodeName is empty", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
                Notify(errno, "nodeName is empty");
                return errno;
            }
            uint nodeNameIndex = getIndexByNodeName(tmpNodeInfo.nodeName);
            if (nodeNameIndex != uint(-1) && nodeNameIndex != index) {
                log("nodeName conflicted", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
                Notify(errno, "nodeName conflicted");
                return errno;
            }
        }

        // Check if admin is a valid user
        if (tmpNodeInfo.nodeAdmin != 0) {
            if (!__getUserDepartmentId(tmpNodeInfo.nodeAdmin).equals(nodeInfos[index].deptId)) {
                log("admin not member of department", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.ADMIN_NOT_MEMBER);
                Notify(errno, "admin is not a member of current node`department");
                return errno;
            }
        }

        if (tx.origin != owner && tx.origin != nodeInfos[index].nodeAdmin) {
            log("sender not owner: ", tx.origin);
            /* if (__checkWritePermission(tx.origin, nodeInfos[index].deptId) == 0
                && checkWritePermission(tx.origin,nodeInfos[index].nodeId) == 0) {
                log("user not dept admin ..", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.NO_PERMISSION);
                Notify(errno, "user is not the admin of department or parent department");
                return errno;
            } */
        }

        // update 
        string memory _tmpStr = "";
        if(_json.keyExists("nodeNAT")) {
            _tmpStr = _json.getObjectValueByKey("nodeNAT");
            if (_tmpStr.keyExists("pubkey")) {
                if (bytes(tmpNodeInfo.nodeNAT.pubkey).length == 0 || bytes(tmpNodeInfo.nodeNAT.pubkey).length > 128
                        || bytes(tmpNodeInfo.nodeNAT.ip).length == 0
                        || tmpNodeInfo.nodeNAT.p2pPort == 0 || tmpNodeInfo.nodeNAT.p2pPort > 0xFFFF
                        || tmpNodeInfo.nodeNAT.tPort == 0 || tmpNodeInfo.nodeNAT.tPort > 0xFFFF
                        || tmpNodeInfo.nodeNAT.rpcPort == 0 || tmpNodeInfo.nodeNAT.rpcPort > 0xFFFF) {
                        log("node pubkey,ip,port invalid", "NodeInfoManager");
                        errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
                        Notify(errno, "node pubkey or ip or port format invalid");
                        return errno;
                }
                if(bytes(tmpNodeInfo.nodeNAT.pubkey).length < 128 ) {
                    uint lessCount = 128 - bytes(tmpNodeInfo.nodeNAT.pubkey).length;
                    for (uint j = 0; j < lessCount; ++j) {
                        tmpNodeInfo.nodeNAT.pubkey = "0".concat(tmpNodeInfo.nodeNAT.pubkey);
                    }
                }

                // check pubkey is exists
                uint pubkeyIndex = getIndexByPubkey(tmpNodeInfo.nodeNAT.pubkey);
                if(pubkeyIndex != uint(-1) && pubkeyIndex != index) {
                    log("node pubkey repeat..", "NodeInfoManager");
                    errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
                    Notify(errno, "node pubkey repeat..");
                    return errno;
                }
            }
        }


        //TODO: 节点一旦创建，归属部门不可更改
        // 当前节点的部门必须为一级机构 departmentLevel == 1

        // Check nodeIPList list
        /* if (_json.keyExists("nodeIPList")) {
            for(uint i = 0 ; i < tmpNodeInfo.nodeIPList.length ; ++i) {
                if(bytes(tmpNodeInfo.nodeIPList[i].ip).length == 0){
                    log("update->node ip is empty");
                    errno = prefix_error + uint256(NodeInfoError.BAD_PARAMETER);
                    Notify(errno,"node ip is empty");
                    return errno;
                }
                if(tmpNodeInfo.nodeIPList[i]._type == uint(1) || tmpNodeInfo.nodeIPList[i]._type == uint(0)) {
                    if (bytes(tmpNodeInfo.nodeIPList[i].ip).length == 0
                            || tmpNodeInfo.nodeIPList[i].P2PPort == 0 
                            || tmpNodeInfo.nodeIPList[i].P2PPort > 0xFFFF
                            || tmpNodeInfo.nodeIPList[i].RPCPort == 0 
                            || tmpNodeInfo.nodeIPList[i].RPCPort > 0xFFFF
                            || tmpNodeInfo.nodeIPList[i].TPort == 0 
                            || tmpNodeInfo.nodeIPList[i].TPort > 0xFFFF) {
                        log("update->nodeIP ip,P2PPort,RPCPort,TPort invalid", "NodeInfoManager");
                        errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
                        Notify(errno, "nodeIP ip,P2PPort,RPCPort,TPort invalid");
                        return errno;
                    }
                }
            }
	    } */
        if (_json.keyExists("nodeName")) {
            nodeInfos[index].nodeName = tmpNodeInfo.nodeName;
        }
        if (_json.keyExists("nodeShortName")) {
            nodeInfos[index].nodeShortName = tmpNodeInfo.nodeShortName;
        }
        if (_json.keyExists("nodeAdmin") && tmpNodeInfo.nodeAdmin != 0) {
            nodeInfos[index].nodeAdmin = tmpNodeInfo.nodeAdmin;
        }
        if (_json.keyExists("nodeDescription")) {
            nodeInfos[index].nodeDescription = tmpNodeInfo.nodeDescription;
        }
        /* if (_json.keyExists("pubkey")) {
            nodeInfos[index].pubkey = tmpNodeInfo.pubkey;
        } */
        /* if (_json.keyExists("ip")) {
            nodeInfos[index].ip = tmpNodeInfo.ip;
        } */
        /* if (_json.keyExists("port")) {
            nodeInfos[index].port = tmpNodeInfo.port;
        } */
        /* if (_json.keyExists("nodeIPList")) {
            nodeInfos[index].nodeIPList = tmpNodeInfo.nodeIPList;
        } */
        if(_json.keyExists("nodeNAT")) {
            nodeInfos[index].nodeNAT = tmpNodeInfo.nodeNAT;
        }
        if(_json.keyExists("nodeLAN")) {
            nodeInfos[index].nodeLAN = tmpNodeInfo.nodeLAN;
        }
        nodeInfos[index].updateTime = uint(now) * 1000;
        errno = uint(NodeInfoError.NO_ERROR);
        revision++;
        log("update nodeInfo succ", "NodeInfoManager");
        Notify(errno, "update node info succ");
        return errno;
    }

    function updateState(string _nodeId,uint _state) public returns(uint){
        log("updateState action", "NodeInfoManager");
        uint prefix_error = 95270;

        uint index = getIndexById(_nodeId);
        if (index == uint(-1)) {
            log("node not exists", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.ID_NOT_EXISTS);
            Notify(errno, "node id dose not exists");
            return errno;
        }

        // ====================== begin check =====================
        if(_state != 0 && _state != 1 && _state != 2){
            log("_state error", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "error");
            return errno;
        }

        nodeInfos[index].state = _state;
        nodeInfos[index].updateTime = uint(now) * 1000;
        errno = uint(NodeInfoError.NO_ERROR);
        revision++;
        log("updateState nodeInfo succ", "NodeInfoManager");
        Notify(errno, "updateState node info succ");
        return errno;
    }

    /**
    * @dev Get enode list
    * @return return enodes in json string
    */
    function getEnodeList() constant public returns(string _json) {
        string memory bookkeeper;
        string memory follower;

        for (uint i = 0; i < nodeInfos.length ; ++i) {
            if (!nodeInfos[i].deleted && nodeInfos[i].deptLevel == 1) {
                if (nodeInfos[i]._type == 1) {
                    if (bytes(bookkeeper).length > 0) {
                         bookkeeper = bookkeeper.concat(",");
                    }

                    bookkeeper = bookkeeper.concat("{");
                    bookkeeper = bookkeeper.concat(nodeInfos[i].nodeId.toKeyValue("nodeId"), ",");
                    bookkeeper = bookkeeper.concat(nodeInfos[i].nodeNAT.pubkey.toKeyValue("pubkey"), ",");
                    bookkeeper = bookkeeper.concat(nodeInfos[i].nodeNAT.ip.toKeyValue("ip"), ",");
                    bookkeeper = bookkeeper.concat(uint(nodeInfos[i].nodeNAT.p2pPort).toKeyValue("port"), ",");
                    bookkeeper = bookkeeper.concat(uint(nodeInfos[i].nodeNAT.activated).toKeyValue("activated"), ",");
                    bookkeeper = bookkeeper.concat(uint(nodeInfos[i].disabled).toKeyValue("disabled"));
                    bookkeeper = bookkeeper.concat("}");
                } else {
                    if (bytes(follower).length > 0) {
                        follower = follower.concat(",");
                    }

                    follower = follower.concat("{");
                    follower = follower.concat(nodeInfos[i].nodeId.toKeyValue("nodeId"), ",");
                    follower = follower.concat(nodeInfos[i].nodeNAT.pubkey.toKeyValue("pubkey"), ",");
                    follower = follower.concat(nodeInfos[i].nodeNAT.ip.toKeyValue("ip"), ",");
                    follower = follower.concat(uint(nodeInfos[i].nodeNAT.p2pPort).toKeyValue("port"));
                    follower = follower.concat(uint(nodeInfos[i].nodeNAT.activated).toKeyValue("activated"), ",");
                    follower = follower.concat(uint(nodeInfos[i].disabled).toKeyValue("disabled"));
                    follower = follower.concat("}");
                }
            }
        }
        
        _json = "{\"ret\":0,\"data\":{\"bookkeeper\":[";
        _json = _json.concat(bookkeeper);
        _json = _json.concat("],\"follower\":[");
        _json = _json.concat(follower);
        _json = _json.concat("]}}");
    }

    /**
    * @dev Activate enode for consensus
    * @param _pubkey enode pubkey in hex string mode
    */
    function ActivateEnode(string _pubkey) public {
        for (uint i=0; i< nodeInfos.length; ++i) {
            if (!nodeInfos[i].deleted && nodeInfos[i].deptLevel == 1) {
                if (nodeInfos[i]._type == 1 && nodeInfos[i].nodeNAT.pubkey.equalsNoCase(_pubkey)) {
                    if (tx.origin == owner || tx.origin == nodeInfos[i].nodeNAT.nodeAddress) {
                        nodeInfos[i].nodeNAT.activated = 1;
                        return;
                    }
                }
            }
        }
    }

    /**
    * @dev If the specified name and ip is in the database
    * @param _commonName Object name
    * @param _ip ip of int
    * @return _json true or false in json string
    */
    function isInWhiteList(string _commonName, string _ip) constant public returns (string _json) {
        // 取到 type = 2 的 uintIP 值，基础值
        // 取到 type = 3 的 uintIP 值,掩码
        // maxInt : 2147483647
        int _uintIP = _ip.toInt();
        log("_uintIP:",_uintIP);
        log("into->isInWhiteList(_common,_uintIP):_commonName",_commonName);
        log("into->isInWhiteList(_common,_uintIP):_uintIP",_uintIP);
        bool isIn = false;
        for (uint i=0; i< nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted) {
                continue;
            }
            if (nodeInfos[i].deptCN.equals(_commonName)) {

                // get uintIP for type in(2,3)
                int baseUintIP = 0;
                int maskUintIP = 0;
                int startUintIP = 0;
                int endUintIP = 0;
                // mask
                if(nodeInfos[i].nodeLAN.goByInt != 0) {
                    baseUintIP = nodeInfos[i].nodeLAN.goByInt;
                }
                if(nodeInfos[i].nodeLAN.mastInt != 0) {
                    maskUintIP = nodeInfos[i].nodeLAN.mastInt;
                }
                if(nodeInfos[i].nodeLAN.startInt != 0) {
                    startUintIP = nodeInfos[i].nodeLAN.startInt;
                }
                if(nodeInfos[i].nodeLAN.endInt != 0) {
                    endUintIP = nodeInfos[i].nodeLAN.endInt;
                }
                // mask check
                if(baseUintIP != 0 && maskUintIP != 0 && (baseUintIP & maskUintIP) == (_uintIP & maskUintIP)) {
                    isIn = true;
                    break;
                }
                // ip range check
                if( startUintIP != 0 && endUintIP != 0 ) {
                    // < 0 deal with 
                    if(startUintIP < 0 || endUintIP < 0) {
                        if(endUintIP <= _uintIP && _uintIP <= startUintIP) {
                            isIn = true;
                            break;
                        }
                    }
                    if(startUintIP > 0) {
                        if(startUintIP <= _uintIP && _uintIP <= endUintIP) {
                            isIn = true;
                            break;
                        }
                    }
                }
            }
        }
        _json = _json.concat("{");
        _json = _json.concat(uint(0).toKeyValue("ret"), ",");
        if (isIn) {
            _json = _json.concat("\"data\":true");
        } else {
            _json = _json.concat("\"data\":false");
        }
        _json = _json.concat("}");
    }

    /**
    * @dev Get Node admin
    * @param _nodeId deparment id
    * @return Node admin address
    */
    function getNodeAdmin(string _nodeId) constant public returns(uint _admin) {
        uint n = 0;
        for (uint i=0; i< nodeInfos.length; ++i) {
            if (nodeInfos[i].nodeId.equals(_nodeId)) {
                if (!nodeInfos[i].deleted) {
                    return uint(nodeInfos[i].nodeAdmin);
                } else {
                    return 0;
                }
            }
        }
        return 0;
    }

    /**
    * @dev Set node admin
    * @param _nodeId The dest node id to check
    * @param _adminAddr The new amdin address
    * @return errno , 0 for true
    */
    function setAdmin(string _nodeId, address _adminAddr) public returns(uint) {
        log("setAdmin", "NodeInfoManager");
        uint prefix_error = 95270;

        uint index = getIndexById(_nodeId);
        if (index == uint(-1)) {
            log("node id not exists", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.ID_NOT_EXISTS);
            return errno;
        }
        /* if (__checkWritePermission(tx.origin, nodeInfos[index].deptId) == 0
            && checkWritePermission(tx.origin,nodeInfos[index].nodeId) == 0) {
            log("No permisson", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.NO_PERMISSION);
            return errno;
        } */
        // 检测用户归属的部门是否为当前节点归属的部门
        if (!__getUserDepartmentId(_adminAddr).equals(nodeInfos[index].deptId)) {
            log("admin not belong to department.", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.ADMIN_NOT_MEMBER);
            return errno;
        }

        // Assign the new address to the dest node
        nodeInfos[index].nodeAdmin = _adminAddr;

        log("setAdmin OK", "NodeInfoManager");
        errno = uint(NodeInfoError.NO_ERROR);
        revision++;
        return errno;
    }

    /**
    * @dev Erase node admin if admin address equals specified address
    * @param _userAddr The amdin address
    * @return errno , 0 for success
    */
    function eraseAdminByAdd(address _userAddr) public returns(uint) {
        log("eraseAdminByAddress", "NodeInfoManager");
        uint prefix_error = 95270;
        uint index = uint(-1);
        for (uint i = 0; i< nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
                continue;
            if (nodeInfos[i].nodeAdmin == _userAddr) {
                index = i;
                break;
            }
        }
        if (index == uint(-1)) {
            log("node id not exists", "NodeInfoManager");
            errno = prefix_error + uint256(NodeInfoError.ID_NOT_EXISTS);
            Notify(errno, "node id not exists");
            return errno;
        }
        if (tx.origin != owner && tx.origin != nodeInfos[index].nodeAdmin) {
            log("tx.origin is not owner: ", tx.origin);
	        /* if (tx.origin != rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0")) {
	            if (__checkWritePermission(tx.origin, nodeInfos[index].deptId) == 0
                    && checkWritePermission(tx.origin,nodeInfos[index].nodeId) == 0) {
	                log("no permission");
	                errno = prefix_error + uint256(NodeInfoError.NO_PERMISSION);
	                Notify(errno, "no permission");
	                return errno;
	            }
        	} */
        }
		nodeInfos[index].nodeAdmin = address(0);

        log("eraseAdminByAddress OK", "NodeInfoManager");
        errno = uint256(NodeInfoError.NO_ERROR);
        revision++;
        Notify(errno,"eraseAdmin success...");
        return errno;
    }

    /**
    * @dev Check if a node exists
    * @param _nodeId The node id to check
    * @return If exists return 1 else return 0
    */
    function nodeInfoExists(string _nodeId) constant public returns (uint _exists) {
        for (uint i=0; i< nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
                continue;
            if (nodeInfos[i].nodeId.equals(_nodeId)) {
                return uint(1);
            }
        }
        return 0;
    }

    /**
    * @dev Delete a department (must be empty)
    * @param _nodeId The department id to check
    * @return If Delete succ return 1 else return 0
    */
    function deleteById(string _nodeId) public {
        uint index = getIndexById(_nodeId);
        uint prefix_error = 95270;
        if (index == uint(-1)) {
            log("nodeId id not exists", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.ID_NOT_EXISTS);
            Notify(errno, "nodeId id not exists");
            return;
        }

        if (tx.origin != owner && tx.origin != nodeInfos[index].nodeAdmin) {
            log("msg.sender is not owner: ", msg.sender);
            log("tx.origin is not nodeAdmin:",tx.origin);
            // check user is belong to department or parent department admin 
            if (__checkWritePermission(tx.origin, nodeInfos[index].deptId) == 0
                && checkWritePermission(tx.origin,nodeInfos[index].nodeId) == 0) {
                log("no permission", "NodeInfoManager");
                errno = prefix_error + uint(NodeInfoError.NO_PERMISSION);
                Notify(errno, "no permission");
                return;
            }
        }

        nodeInfos[index].deleted = true;
        log("delete succ", "NodeInfoManager");
        errno = uint(NodeInfoError.NO_ERROR);
        Notify(errno, "delete succ");
        revision++;
    }

    function getRevision() constant public returns (uint _ret) {
        return revision;
    }

    /**
    * @dev check if the IP is used by any node
    * @param _ip The role id
    * @return _used If contains return 1, else return 0
    */
    function IPUsed(string _ip) constant public returns (uint _used) {
        for (uint i=0; i < nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted) {
                continue;
            }
            if(nodeInfos[i].nodeNAT.ip.equals(_ip)) {
                return 1;
            }
        }
        return 0;
    }

    /**
    * @dev List the all objects
    * @return No return
    */
    function listAll() constant public returns (string _json) {
        _json = listBy(0, "");
    }

    /**
    * @dev Find object by id
    * @param _id Object id
    * @return _json Objects in json string
    */
    function findById(string _id) constant public returns(string _json) {
        _json = listBy(1, _id);
    }

    /**
    * @dev Find object by name
    * @param _name Object name
    * @return _json Objects in json string
    */
    function findByName(string _name) constant public returns(string _json) {
        _json = listBy(2, _name);
    }

    /**
    * @dev Find object by department id
    * @param _departmentId Object prarent id
    * @return _json Objects in json string
    */
    function findByDepartmentId(string _departmentId) constant public returns(string _json) {
        _json = listBy(3, _departmentId);
    }

    /**
    * @dev Find object by nodeAdmin id
    * @param _nodeAdmin Object prarent id
    * @return _json Objects in json string
    */
    function findByNodeAdmin(address _nodeAdmin) constant public returns(string _json) {
        string memory tmp = uint(_nodeAdmin).toAddrString();
        _json = listBy(4, tmp);
    }

    /**
    * @dev Find object by pubkey
    * @param _pubkey Object prarent id
    * @return _json Objects in json string
    */
    function findByPubkey(string _pubkey) constant public returns(string _json) {
        _json = listBy(5, _pubkey);
    }

    /**
    * @dev list elem by condition
    * @param _cond for the condition
    *        0 for all
    *        1 for id
    *        2 for name
    *        3 for deptId
    *        4 for nodeAdmin
    *        5 for pubkey
    * @param _value The condition value
    * @return json string
    */
    function listBy(uint _cond, string _value) constant private returns(string _json) {
    	uint tatal = 0;
        for (uint i=0; i < nodeInfos.length; ++i) {
            if (!nodeInfos[i].deleted) {
                tatal++;
            }
        }

        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.append(",\"data\":{");
        len = LibStack.appendKeyValue("total", tatal);
        len = LibStack.append(",\"items\":[");
 
        uint n = 0;
        for (i=0; i < nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
                continue;

            bool suitable = false;
            if (_cond == 0) {
                suitable = true;
            } else if (_cond == 1) {
                if (nodeInfos[i].nodeId.equals(_value))
                    suitable = true;
            } else if (_cond == 2) {
                if (nodeInfos[i].nodeName.equals(_value))
                    suitable = true;
            } else if (_cond == 3) {
                if (nodeInfos[i].deptId.equals(_value))
                    suitable = true;
            } else if (_cond == 4) {
                if (uint(nodeInfos[i].nodeAdmin).toAddrString().equals(_value))
                    suitable = true;
            } else if (_cond == 5) {
                if (nodeInfos[i].nodeNAT.pubkey.equals(_value))
                    suitable = true;
            }

            if (suitable) {
                if (n > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(nodeInfos[i].toJson());
                n++;
            }
        }
        len = LibStack.append("]}}");
        _json = LibStack.popex(len);
    }

    function listByStateAndTypeAndName(uint _state, uint _type,string _nodeName,uint _pageNum,uint _pageSize ) constant public returns(string _json) {

        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.append(",\"data\":{");
        len = LibStack.append("\"items\":[");

        uint total = 0;
        uint n = 0;
        uint m =0;
        for (uint i = 0; i < nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
            continue;

            if ((bytes(_nodeName).length == 0 || nodeInfos[i].nodeName.indexOf(_nodeName) != - 1)) {

                if(_state == 2 || nodeInfos[i].state == _state){

                   if(_type==2 || nodeInfos[i]._type == _type){
                       if (n >= _pageNum*_pageSize && n < (_pageNum+1)*_pageSize) {
                           if (m > 0) {
                               len = LibStack.append(",");
                           }
                           len = LibStack.append(nodeInfos[i].toJson());
                           m++;
                       }

                       if (n >= (_pageNum+1)*_pageSize) {
                           break;
                       }
                       n++;
                   }
                }
            }
        }

        for (i = 0; i < nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
            continue;

            if ((bytes(_nodeName).length == 0 || nodeInfos[i].nodeName.indexOf(_nodeName) != - 1)) {

                if(_state == 2 || nodeInfos[i].state == _state){

                    if(_type==2 || nodeInfos[i]._type == _type){
                        total++;
                    }
                }
            }
        }
        len = LibStack.append("]");
        len = LibStack.appendKeyValue("total", total);
        len = LibStack.append("}}");
        _json = LibStack.popex(len);

    }

    function checkWritePermission(address _addr, string _nodeInfoId) constant public returns (uint _ret) {
        // Is current caller is super admin
        if (_addr == owner) {
            return 1;
        }

        uint nodeIndex = getIndexById(_nodeInfoId);
        if (nodeIndex == uint(-1)) {
            return 0;
        }

        if (_addr == nodeInfos[nodeIndex].nodeAdmin) {
            return 1;
        }
    }

    // TODO: updated
    function getIndexByPubkey(string _pubkey) constant private returns (uint _index) {
        for (uint i=0; i < nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
                continue;
            if (nodeInfos[i].nodeNAT.pubkey.equals(_pubkey))
                return i;
        }
        return uint(-1);
    }

    function getIndexByNodeName(string _nodeName) constant private returns (uint _index) {
        for (uint i=0; i< nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
                continue;
            if (nodeInfos[i].nodeName.equals(_nodeName))
                return i;
        }
        return uint(-1);
    }

    function getIndexById(string _nodeId) constant private returns (uint _index) {
        for (uint i=0; i< nodeInfos.length; ++i) {
            if (nodeInfos[i].deleted)
                continue;
            if (nodeInfos[i].nodeId.equals(_nodeId))
                return i;
        }
        return uint(-1);
    }

    /**
    * check the user is have permission to operator nodeInfo
    * @param _addr user address
    * @param _departmentId department id
    * @return json string
    */
    function __checkWritePermission(address _addr, string _departmentId) constant private returns (uint _ret) {
        if(_addr == owner) {
            return 1;
        }
        address deptAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager","0.0.1.0");
        if(deptAddr == 0) {
            log("__checkWritePermission,deptAddr is null ","NodeInfoManager");
            return 0;
        }
        IDepartmentManager dm = IDepartmentManager(deptAddr);
        return dm.checkWritePermission(_addr,_departmentId);
    }

    function __getUserDepartmentId(address _userAddr) constant internal returns (string _ret) {
        address userManagerAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","UserManager", "0.0.1.0");
        if (userManagerAddr == 0) {
            return uint(0).recoveryToString();
        }
        IUserManager userManager = IUserManager(userManagerAddr);
        return userManager.getUserDepartmentId(_userAddr).recoveryToString();
    }

    function __departmentExists(string _departmentId) constant private returns (uint _ret) {
        address deptAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager","0.0.1.0");
        if(deptAddr == 0) {
            log("__departmentExists, deptAddr is null ","NodeInfoManager");
            return 0;
        }
        IDepartmentManager dm = IDepartmentManager(deptAddr);
        return dm.departmentExists(_departmentId);
    }

    function __departmentExistsByCN(string _commonName) constant private returns (uint _ret) {
        address deptAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","DepartmentManager","0.0.1.0");
        if(deptAddr == 0) {
            log("__departmentExistsByCN, deptAddr is null ","NodeInfoManager");
            return 0;
        }
        IDepartmentManager dm = IDepartmentManager(deptAddr);
        return dm.departmentExistsByCN(_commonName);
    }

    function updateDisabled(string _nodeId,uint _disabled) public returns (uint) {
        log("updateDisabled action", "NodeInfoManager");
        uint prefix_error = 95270;

        uint index = getIndexById(_nodeId);
        if (index == uint(-1)) {
            log("node not exists", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.ID_NOT_EXISTS);
            Notify(errno, "node id dose not exists");
            return errno;
        }

        // ====================== begin check =====================
        if (_disabled != 0 && _disabled != 1) {
            log("_disabled error", "NodeInfoManager");
            errno = prefix_error + uint(NodeInfoError.BAD_PARAMETER);
            Notify(errno, "error");
            return errno;
        }

        nodeInfos[index].disabled = _disabled;
        nodeInfos[index].updateTime = uint(now) * 1000;
        errno = uint(NodeInfoError.NO_ERROR);
        revision++;
        log("updateDisabled nodeInfo succ", "NodeInfoManager");
        Notify(errno, "updateDisabled node info succ");
        return errno;
    }


}
