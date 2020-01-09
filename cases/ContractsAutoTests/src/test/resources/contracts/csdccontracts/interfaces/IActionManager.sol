pragma solidity ^0.4.12;

/** @title ActionManager interface*/
contract IActionManager {

    function findActionByType(uint256 _type) constant public returns(string _json) ;

    /**@dev get action list by module name and version
     * @param _moduleName and _moduleVersion
     * @return _json return result.
     */
    function getActionListByModuleName(string _moduleName, string _moduleVersion) constant public returns(string _json) ;

    function getActionListByContractName(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) constant public returns(string _json) ;

    function getActionListByModuleId(string _moduleId) constant public returns(string _json) ;

    function getActionListByContractId(string _contractId) constant public returns(string _json) ;

    function actionExists(string _actionId) constant public returns(uint _ret);

    function queryActionEnable(string _actionId, bool _checkOwner) constant public returns(uint _ret);

    /**
    * find the action by key
    * @param _resKey ,contract name
    * @param _opKey ,function brief name, eg. foo(uints)
    * @return return the result json, items contain the object
    */
    function findByKey(string _resKey, string _opKey) constant public returns(string _actionJson) ;

    /**
    * find the action by id
    * @param _actionId ,the action id
    * @return return the result json, items contain the object
    */
    function findById(string _actionId) constant public returns(string _actionJson) ;

    /**
    * list all actions
    * @return the result json, items contain the object
    */
    function listAll() constant public returns(string _actionListJson) ;

    /**
    * list actions of a contract
    * @return the result json, items contain the object
    */
    function listContractActions(string _contractName) constant public returns(string _actionListJson) ;

    /**
    * check if the action id match the resKey and opKey
    * @param _actionId the id of action
    * @param _contractAddr contract address
    * @param _opSha3Key the sha3 value of brief function name
    * @return the result json, items contain the object
    */
    function checkActionWithKey(string _actionId, address _contractAddr, string _opSha3Key) constant public returns(uint _ret) ;

    /**
    * insert into a action object
    * @param _actionJson the action object json
    * @return true is success else failed
    */
    function insert(string _actionJson) public returns(bool _ret);

    /**
    * update a action object
    * @param _actionJson the action object json
    * @return true is success else failed
    */
    function update(string _actionJson) public returns(bool _ret);

    /**
    * remove action by actionId
    * @param _actionId the id of action
    * @return none
    */
    function deleteById(string _actionId) public ;

}
