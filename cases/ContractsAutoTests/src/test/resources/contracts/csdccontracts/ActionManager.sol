pragma solidity ^0.4.12;
/**
*@file      CommonWidget.sol
*@author    kelvin
*@time      2016-11-29
*@desc      the defination of ActionManager
*/

import "./library/LibAction.sol";
import "./sysbase/OwnerNamed.sol";
import "./interfaces/IActionManager.sol";
import "./interfaces/IRoleManager.sol";

contract ActionManager is OwnerNamed,IActionManager {
    
    using LibAction for *;
    using LibString for *;
    using LibInt for *;

    event Notify(uint _errno, string _info);

    mapping(string=>LibAction.Action)           actionMap;//<id, Action>
    string[]                                    actionIdList;
    string[]                                    tempActionIdList;

    mapping(string=>string)                     keyMap; //<resKey_opKey, id>
    mapping(string=>string)                     opKeyMap;//<opKeySha3, opKey>

    LibAction.Action        internal            m_Action;

    enum ActionError {
        NO_ERROR,
        BAD_PARAMETER,
        NAME_EMPTY,
        ACTION_NOT_EXISTS,
        ACTION_ID_ALREADY_EXISTS,
        ACTION_KEY_ALREADY_EXISTS,
        CONTRACT_NOT_REGISTER,
        ACTION_USED
    }

    function ActionManager() {
        register("SystemModuleManager","0.0.1.0","ActionManager", "0.0.1.0");
    }

    function findActionByType(uint256 _type) constant public returns(string _json) {

        uint len = 0;
    
        uint counter = 0;
        len = LibStack.push("");
        for (uint index = 0; index < actionIdList.length; index++) {
            if (actionMap[actionIdList[index]].state == LibAction.ActionState.INVALID) {
                continue;
            }

            if (actionMap[actionIdList[index]].Type == _type) {
                if (index == actionIdList.length-1) {
		    len = LibStack.append(actionMap[actionIdList[index]].toJson());
                } else {
		    len = LibStack.append(actionMap[actionIdList[index]].toJson());
		    len = LibStack.append(",");
                }
                counter++;
             }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),counter);
        _json = LibStack.popex(_retLen);
    }

    function getActionListByModuleName(string _moduleName, string _moduleVersion) constant public returns(string _json) {
        _json = listBy(3,_moduleName,_moduleVersion,"","");
    }

    function getActionListByContractName(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) constant public returns(string _json) {
        _json = listBy(4,_moduleName,_moduleVersion,_contractName,_contractVersion);
    }

    function getActionListByModuleId(string _moduleId) constant public returns(string _json) {
        _json = listByForUK(3,_moduleId);
    }

    function getActionListByContractId(string _contractId) constant public returns(string _json) {
        _json = listByForUK(4,_contractId);
    }

    /**
    * @dev list elem by condition
    * @param _cond for the condition
    *        0 for all
    *        1 for id
    *        2 for name
    *        3 for moduleId
    *        4 for contractId
    *        5 for type
    * @param _value The condition value
    * @return _json Objects in json string
    */
    function listByForUK(uint _cond,string _value) constant public returns(string _json) {
        uint tatal = 0;     
        for(uint i = 0 ; i < actionIdList.length; ++i){
            if(actionMap[actionIdList[i]].state != LibAction.ActionState.INVALID){
                tatal++;
            }
        }

        uint len = 0;
        len = LibStack.push("");
        uint n = 0;
        for(i  = 0 ; i < actionIdList.length; ++i){
            if(actionMap[actionIdList[i]].state == LibAction.ActionState.INVALID){
                continue;
            }
            bool suitable = false;
            if(_cond == 0){
                suitable = true;
            }else if(_cond == 1){
                if(actionMap[actionIdList[i]].id.equals(_value)){
                    suitable = true;
                }
            }else if(_cond == 2){
                if(actionMap[actionIdList[i]].name.equals(_value)){
                    suitable = true;
                }
            }else if(_cond == 3){
                if(actionMap[actionIdList[i]].moduleId.equals(_value)){
                    suitable = true;
                }
            }else if(_cond == 4){
                if(actionMap[actionIdList[i]].contractId.equals(_value)){
                    suitable = true;
                }
            }
            if(suitable){
                if(n > 0){
                    len = LibStack.append(",");
                }
                len = LibStack.append(actionMap[actionIdList[i]].toJson());
                n++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),tatal);
        _json = LibStack.popex(_retLen);
    }

    /**
    * @dev list elem by condition
    * @param _cond for the condition
    *        0 for all
    *        1 for id(action id)
    *        2 for name(action name)
    *        3 for moduleName,moduleVersion
    *        4 for moduleName,moduleVersion,contractName,contractVersion
    * @return _json Objects in json string
    */
    function listBy(uint _cond,string _name1,string _version1,string _name2,string _version2) constant public returns(string _json) {
        uint tatal = 0;
        for(uint i = 0 ; i < actionIdList.length; ++i){
            if(actionMap[actionIdList[i]].state != LibAction.ActionState.INVALID){
                tatal++;
            }
        }
        uint len = 0;
        len = LibStack.push("");
        uint n = 0;
        for(i  = 0 ; i < actionIdList.length; ++i){
            if(actionMap[actionIdList[i]].state == LibAction.ActionState.INVALID){
                continue;
            }
            bool suitable = false;
            if(_cond == 0){
                suitable = true;
            }else if(_cond == 1){
                if(actionMap[actionIdList[i]].id.equals(_name1)){
                    suitable = true;
                }
            }else if(_cond == 2){
                if(actionMap[actionIdList[i]].name.equals(_name1)){
                    suitable = true;
                }
            }else if(_cond == 3){
                if(actionMap[actionIdList[i]].moduleName.equals(_name1) && actionMap[actionIdList[i]].moduleVersion.equals(_version1)){
                    suitable = true;
                }
            }else if(_cond == 4){
                if(actionMap[actionIdList[i]].moduleName.equals(_name1) && actionMap[actionIdList[i]].moduleVersion.equals(_version1)
                  && actionMap[actionIdList[i]].resKey.equals(_name2) && actionMap[actionIdList[i]].version.equals(_version2)){
                    suitable = true;
                }
            }
            if(suitable){
                if(n > 0){
                    len = LibStack.append(",");
                }
                len = LibStack.append(actionMap[actionIdList[i]].toJson());
                n++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),tatal);
        _json = LibStack.popex(_retLen);
    }

    /**
    * check if action id exists
    * @param _actionId 
    * @return 0 , the action exists, else not exists
    */
    function actionExists(string _actionId) constant public returns(uint _ret){
        if (actionMap[_actionId].state == LibAction.ActionState.INVALID) {
            return 0;
        }
        return 1;
    }

    function queryActionEnable(string _actionId, bool _checkOwner) constant public returns(uint _ret){
        if (actionMap[_actionId].state == LibAction.ActionState.INVALID) {
            return uint(-1);
        }
        if (_checkOwner && actionMap[_actionId].creator != tx.origin)
            return uint(-2);
        else
            return actionMap[_actionId].enable;
    }

    /**
    * find the action by key
    * @param _resKey ,contract name
    * @param _opKey ,function brief name, eg. foo(uints)
    * @return return the result json, items contain the object
    */
    function findByKey(string _resKey, string _opKey) constant public returns(string _actionJson) {
        string memory strKey = _resKey;
        strKey = strKey.concat("_", _opKey);
        
        string memory actionId = keyMap[strKey];

        uint len = 0;
        len = LibStack.push("");
        if (actionMap[actionId].state != LibAction.ActionState.INVALID) {
            len = LibStack.append(actionMap[actionId].toJson());
        }
        uint _retLen = itemStackPush(LibStack.popex(len),getCount());
        _actionJson = LibStack.popex(_retLen);
    }

    /**
    * find the action by id
    * @param _actionId ,the action id
    * @return return the result json, items contain the object
    */
    function findById(string _actionId) constant public returns(string _actionJson) {
        uint len = 0;
        len = LibStack.push("");
        if (actionMap[_actionId].state != LibAction.ActionState.INVALID) {
            len = LibStack.append(actionMap[_actionId].toJson());
        }
        uint _retLen = itemStackPush(LibStack.popex(len),getCount());
        _actionJson = LibStack.popex(_retLen);
    }

    /**
    * list all actions
    * @return the result json, items contain the object
    */
    function listAll() constant public returns(string _actionListJson) {
        uint len = 0;
        len = LibStack.push("");
        for (uint index = 0; index < actionIdList.length; index++) {
            if (actionMap[actionIdList[index]].state != LibAction.ActionState.INVALID) {
                if (index == actionIdList.length-1) {
                    len = LibStack.append(actionMap[actionIdList[index]].toJson());
		} else {
                    len = LibStack.append(actionMap[actionIdList[index]].toJson());
		    len = LibStack.append(",");
                }
             }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),getCount());
        _actionListJson = LibStack.popex(_retLen);
    }

    /**
    * list actions of a contract
    * @return the result json, items contain the object
    */
    function listContractActions(string _contractName) constant public returns(string _actionListJson) {
       
        uint len = 0;
        uint counter = 0;
        len = LibStack.push("");
        for (uint index = 0; index < actionIdList.length; index++) {
            if (actionMap[actionIdList[index]].state != LibAction.ActionState.INVALID && 
                actionMap[actionIdList[index]].resKey.equals(_contractName)) {
                if (counter > 0) {
                    len = LibStack.append(",");
                }
                len = LibStack.append(actionMap[actionIdList[index]].toJson());
                counter++;
            }
        }
        uint _retLen = itemStackPush(LibStack.popex(len),getCount());
        _actionListJson = LibStack.popex(_retLen);
    }


    /**
    * check if the action id match the resKey and opKey
    * @param _actionId the id of action
    * @param _contractAddr contract address
    * @param _opSha3Key the sha3 value of brief function name
    * @return the result json, items contain the object
    */
    function checkActionWithKey(string _actionId, address _contractAddr, string _opSha3Key) constant public returns(uint _ret) {
        uint uintResKey = rm.findResNameByAddress(_contractAddr);
        string memory strResKey = uintResKey.recoveryToString();
        if (bytes(strResKey).length <= 0) {
            return 0;
        }
        string memory strKey = strResKey;
        string memory strOpKey = opKeyMap[_opSha3Key.substr(0, 8).toLower()];
        strKey = strKey.concat("_", strOpKey);

        string memory id = keyMap[strKey];
        if (actionMap[id].state == LibAction.ActionState.INVALID) {
            return 0;
        }
        if (_actionId.equals(id)) {
            return 1;
        }
        return 0;
    }

    /**
    * insert into a action object
    * @param _actionJson the action object json
    * @return true is success else failed
    */
    function insert(string _actionJson) public returns(bool _ret){
        log("insert a action", "ActionManager");

        _ret = false;
        if (!m_Action.jsonParse(_actionJson)) {
            log("action json is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "bad input json as action");
            return;
        }
        if (m_Action.moduleName.equals("")) {
            log("moduleName is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "moduleName is invalid");
            return;
        }
        if (m_Action.moduleVersion.equals("")) {
            log("moduleVersion is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "moduleVersion is invalid");
            return;
        }
        if (m_Action.id.equals("")) {
            log("action id is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "action id is invalid");
            return;
        }
        if (m_Action.name.equals("")) {
            log("action name is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.NAME_EMPTY);
            Notify(errno, "action name is invalid");
            return;
        }

        // check if the contract register
        address contractAddr = rm.getContractAddress(m_Action.moduleName, m_Action.moduleVersion, m_Action.resKey, m_Action.version);
        if (contractAddr == 0x0 && m_Action.Type != 1) {//not register
            log("contract not registered", "ActionManager");
            log("current actionId is",m_Action.id);
            log("current action type is",m_Action.Type.toString());
            errno = 15500 + uint(ActionError.CONTRACT_NOT_REGISTER);
            Notify(errno, "contract not registered");
            return;
        }
        
        // check if action already exists
        if (actionMap[m_Action.id].state != LibAction.ActionState.INVALID) {
            log("duplicate action id insert", "ActionManager");
            errno = 15500 + uint(ActionError.ACTION_ID_ALREADY_EXISTS);
            Notify(errno, "duplicate action id, cannot insert");
            return;
        }

        string memory key = m_Action.resKey;
        key = key.concat("_", m_Action.opKey);
        keyMap[key] = m_Action.id;

        // insert action 
        m_Action.createTime = now*1000;
        m_Action.updateTime = now*1000;
        m_Action.state = LibAction.ActionState.VALID;
        m_Action.creator = tx.origin;

        actionMap[m_Action.id] = m_Action;

        bool _isFind = false;
        for (uint index = 0; index < actionIdList.length; ++index) {
            if (actionIdList[index].equals(m_Action.id)) {
                _isFind = true;
                break;
            }
        }
        if (!_isFind) {
            actionIdList.push(m_Action.id);
        }
        
        uint sha3Value = uint(sha3(m_Action.opKey));//bytes32
        string memory funNameSha3 = sha3Value.toHexString64().toLower().substr(2, 8); //toHexString64();//TODO: 注意此处需要sha3 sha3Value.toHexString64();
        opKeyMap[funNameSha3] = m_Action.opKey;
        
        log("insert action success", "ActionManager");
        m_Action.resetAction();
        _ret = true;
        errno = uint(ActionError.NO_ERROR);
        Notify(errno, "insert action success");
        return;
    }

    /**
    * update a action object
    * @param _actionJson the action object json
    * @return true is success else failed
    */
    function update(string _actionJson) public returns(bool _ret){
        log("update a action", "ActionManager");

        _ret = false;
        if (!m_Action.jsonParse(_actionJson)) {
            log("action json is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "bad input json as action");
            return;
        }
        if (m_Action.moduleName.equals("")) {
            log("moduleName is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "moduleName is invalid");
            return;
        }
        if (m_Action.moduleVersion.equals("")) {
            log("moduleVersion is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "moduleVersion is invalid");
            return;
        }
        if (m_Action.id.equals("")) {
            log("action id is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.BAD_PARAMETER);
            Notify(errno, "action id is invalid");
            return;
        }
        if (m_Action.name.equals("")) {
            log("action name is invalid", "ActionManager");
            errno = 15500 + uint(ActionError.NAME_EMPTY);
            Notify(errno, "action name is invalid");
            return;
        }

        // check if the contract register
        address contractAddr = rm.getContractAddress(m_Action.moduleName, m_Action.moduleVersion, m_Action.resKey, m_Action.version);
        if (contractAddr == 0x0 && m_Action.Type != 1) {//not register
            log("contract not registered", "ActionManager");
            errno = 15500 + uint(ActionError.CONTRACT_NOT_REGISTER);
            Notify(errno, "contract not registered");
            return;
        }

        // check if action already exists
        if (actionMap[m_Action.id].state == LibAction.ActionState.INVALID) {
            log("action not exist", "ActionManager");
            errno = 15500 + uint(ActionError.ACTION_ID_ALREADY_EXISTS);
            Notify(errno, "action not exist");
            return;
        }

        // forbidden to change moduleName or moduleVersion
        if (!actionMap[m_Action.id].moduleName.equals(m_Action.moduleName) || !actionMap[m_Action.id].moduleVersion.equals(m_Action.moduleVersion)) {
            log("forbidden to change moduleName or moduleVersion", "ActionManager");
            errno = 15500 + uint(ActionError.ACTION_ID_ALREADY_EXISTS);
            Notify(errno, "forbidden to change moduleName or moduleVersion");
            return;
        }

        // update action
        if (_actionJson.keyExists("name") && !actionMap[m_Action.id].name.equals(m_Action.name)) {
        		log("update action name from->to:", actionMap[m_Action.id].name, m_Action.name);
            actionMap[m_Action.id].name = m_Action.name;
        }
        if (_actionJson.keyExists("level") && actionMap[m_Action.id].level != m_Action.level) {
        		log("update action level from->to:", actionMap[m_Action.id].level.toString(), m_Action.level.toString());
            actionMap[m_Action.id].level = m_Action.level;
        }
        if (_actionJson.keyExists("type") && actionMap[m_Action.id].Type != m_Action.Type) {
        		log("update action type from->to:", actionMap[m_Action.id].Type.toString(), m_Action.Type.toString());
            actionMap[m_Action.id].Type = m_Action.Type;
        }
        if (_actionJson.keyExists("enable") && actionMap[m_Action.id].enable != m_Action.enable) {
        		log("update action enable from->to:", actionMap[m_Action.id].enable.toString(), m_Action.enable.toString());
            actionMap[m_Action.id].enable = m_Action.enable;
        }
        if (_actionJson.keyExists("parentId") && !actionMap[m_Action.id].parentId.equals(m_Action.parentId)) {
        		log("update action parentId from->to:", actionMap[m_Action.id].parentId, m_Action.parentId);
            actionMap[m_Action.id].parentId = m_Action.parentId;
        }
        if (_actionJson.keyExists("url") && !actionMap[m_Action.id].url.equals(m_Action.url)) {
        		log("update action url from->to:", actionMap[m_Action.id].url, m_Action.url);
            actionMap[m_Action.id].url = m_Action.url;
        }
        if (_actionJson.keyExists("resKey") && !actionMap[m_Action.id].resKey.equals(m_Action.resKey)) {
        		log("update action resKey from->to:", actionMap[m_Action.id].resKey, m_Action.resKey);
            actionMap[m_Action.id].resKey = m_Action.resKey;
        }
        if (_actionJson.keyExists("opKey") && !actionMap[m_Action.id].opKey.equals(m_Action.opKey)) {
        		log("update action opKey from->to:", actionMap[m_Action.id].opKey, m_Action.opKey);
            actionMap[m_Action.id].opKey = m_Action.opKey;
        }
        if (_actionJson.keyExists("version") && !actionMap[m_Action.id].version.equals(m_Action.version)) {
        		log("update action version from->to:", actionMap[m_Action.id].version, m_Action.version);
            actionMap[m_Action.id].version = m_Action.version;
        }
        actionMap[m_Action.id].updateTime = now*1000;

        log("update action success", "ActionManager");
        m_Action.resetAction();
        _ret = true;
        errno = uint(ActionError.NO_ERROR);
        Notify(errno, "update action success");
        return;
    }

    function deleteById(string _actionId) public {
        log("delete by id", "ActionManager");
        log(_actionId, "ActionManager");

        IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0"));
        if (roleMgr.actionUsed(_actionId) == 1) {
            log("action used, not delete", "ActionManager");
            errno = 15500 + uint(ActionError.ACTION_USED);
            Notify(errno, "action used, not delete");
            return;
        }
        if (actionMap[_actionId].state == LibAction.ActionState.INVALID) {
            log("action not exists", "ActionManager");
            errno = 15500 + uint(ActionError.ACTION_NOT_EXISTS);
            Notify(errno, "action not exists");
            return;
        }

        //delete the element
        delete tempActionIdList;
        for (uint i = 0; i < actionIdList.length; ++i) {
            if (actionIdList[i].equals(_actionId)) {
                continue;
            }

            tempActionIdList.push(actionIdList[i]);
        }

        delete actionIdList;
        for (i = 0; i < tempActionIdList.length; ++i) {
            actionIdList.push(tempActionIdList[i]);
        }

        log("delete OK", "ActionManager");
        actionMap[_actionId].state = LibAction.ActionState.INVALID;

        errno = uint(ActionError.NO_ERROR);
        Notify(errno, "delete action success");
    }

    /**
    * get action count 
    * @return the length of actionIdList
    */
    function getCount() constant returns(uint _count) {
        _count = actionIdList.length;
    }

    //items入栈
    function itemStackPush(string _items, uint _total) constant private returns (uint len){
        len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.append(",");
        len = LibStack.append("\"data\":{");
        len = LibStack.appendKeyValue("total", _total);
        len = LibStack.append(",");
        len = LibStack.append("\"items\":[");
        len = LibStack.append(_items);
        len = LibStack.append("]}}");
        return len;
    }
}
