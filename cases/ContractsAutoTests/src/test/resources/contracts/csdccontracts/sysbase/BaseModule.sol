pragma solidity ^0.4.12;
/**
* @file      BaseModule.sol
* @author    Jungle
* @time      2017-5-18 10:37:00
* @desc      define a base interface
*/

import "../sysbase/OwnerNamed.sol";
import "../interfaces/IRoleFilterManager.sol";

contract BaseModule is OwnerNamed {

    using LibString for *;
    using LibInt for *;

    IRoleFilterManager rfm;
    uint reversion ;

    function BaseModule() {
    	address rfmAddr = rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleFilterManager","0.0.1.0");
    	if(rfmAddr == 0) {
    		log("baseModule address is null.","BaseModule");
    	}
    	reversion = 0;
    	rfm = IRoleFilterManager(rfmAddr);
    }

    /**
    * Add a new module
    * @param _json module info
    * @return 0 for false or moduleId
    */
    function addModule(string _json) internal returns(uint _ret) {
        return rfm.addModule(_json);
    }

    /**
    * Add a new contract
    * @param _json contract info
    * @return 0 for false or moduleId
    */
    function addContract(string _json) internal returns(uint _ret) {
        return rfm.addContract(_json);
    }

    /**
    * Add a new menu
    * @param _json menu info
    * @return 0 for false or moduleId
    */
    function addMenu(string _json) internal returns(uint _ret) {
        return rfm.addMenu(_json);
    }

    /**
    * Add a new action
    * @param _json action info
    * @return 0 for false or moduleId
    */
    function addAction(string _json) internal returns(uint _ret) {
        return rfm.addAction(_json);
    }

    /**
    * Add a new role
    * @param _json role info
    * @return 0 for false or moduleId
    */
    function addRole(string _json) internal returns(uint _ret) {
        return rfm.addRole(_json);
    }

    function updModule(string _json) public returns(uint _ret) {
        return rfm.updModule(_json);
    }

    function moduleIsExist(string _json) public returns(uint _ret) {
        return rfm.moduleIsExist(_json);
    }

    function addActionToRole(string _moduleId, string _roleId, string _actionId) internal returns(uint _ret) {
        return rfm.addActionToRole(_moduleId,_roleId,_actionId);
    }

    /**
    * change module owner, and all contracts owner
    * @param _moduleName moduleName
    * @param _moduleVersion moduleVersion
    * @param _newOwner new owner
    * @return 0 for success, else -1 for failed
    */
    function changeModuleOwner(string _moduleName, string _moduleVersion, address _newOwner) internal returns(uint _ret) {
    	return rfm.changeModuleOwner(_moduleName,_moduleVersion,_newOwner);
    }
}
