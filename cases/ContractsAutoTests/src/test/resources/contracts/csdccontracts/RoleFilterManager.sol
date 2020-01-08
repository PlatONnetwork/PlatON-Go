/**
 * @file RoleFilterManager.sol
 * @author Jungle
 * @time 2017年5月13日17:31:01
 * @desc 角色合约过滤器
 *
 */
pragma solidity ^0.4 .2;

import "./sysbase/OwnerNamed.sol";
import "./library/LibModule.sol";
import "./library/LibFilter.sol";
import "./library/LibAction.sol";
import "./library/LibRole.sol";
import "./library/LibContract.sol";

import "./interfaces/FilterBase.sol";
import "./interfaces/IRoleManager.sol";
import "./interfaces/IActionManager.sol";
import "./interfaces/IUserManager.sol";
import "./interfaces/IRoleFilterManager.sol";

// the filter for role mgr
contract RoleFilterManager is OwnerNamed, FilterBase, IRoleFilterManager {

	using LibString for * ;
	using LibInt for * ;
	using LibModule for * ;
	using LibJson for * ;
	using LibFilter for * ;
	using LibAction for * ;
	using LibRole for * ;
	using LibContract for * ;

	uint revision;
	string[] tmpArray; // use for storage

	// the contractId mapping to Contract info
	//key: moduleName.concat(moduleVersion,contractName,contractVersion);
	mapping(string => LibContract.Contract) contractMapping;

	LibModule.Module[] modules;
	LibModule.Module tmpModule;
	LibFilter.Filter tmpFilter;

	LibAction.Action tmpAction;
	LibRole.Role tmpRole;
	LibContract.Contract tmpContract;
	LibContract.Contract tmpContractFind;
	LibContract.OpenAction tmpOpenAction;

	enum ROLE_FILTER_ERRORCODE {
		NO_ERROR,
		JSON_INVALID,
		INVALID_PARAM,
		SYS_EXCEPTION,
		MODULE_NOT_EXIST
	}

	event Notify(uint _error, string _info);

	function setAdministratorAddress(address _addr) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "setAdministratorAddress", uint(_addr).toAddrString());
		_ret = writedb("administrator|update", "address", uint(_addr).toAddrString());
		if(0 != _ret)
			log("RoleFilterManager.sol", "setAdministratorAddress failed.");
		else
			log("RoleFilterManager.sol", "setAdministratorAddress success.");
	}

	//Notice : AuthorizeFilters init must be in front of RoleFilterManager
	function RoleFilterManager() {
		revision = 0;
		register("SystemModuleManager", "0.0.1.0", "RoleFilterManager", "0.0.1.0");
		setAdministratorAddress(owner);

		FilterBase af = FilterBase(rm.getContractAddress("SystemModuleManager", "0.0.1.0", "AuthorizeFilters", "0.0.1.0"));
		tmpFilter = LibFilter.Filter({
			id: "0",
			name: "RoleFilterManager",
			version: "0.0.1.0",
			_type: 1,
			state: 1,
			enable: 1,
			desc: "角色过滤器",
			addr: msg.sender
		});
		uint ret = af.addFilter(tmpFilter.toJson());
		log("addFilter success, filterId:", ret.recoveryToString());
	}

	/**
	 * execute authorization
	 * @param _from : user addr
	 * @param _to : contract addr
	 * @param _funcHash : openkey sha3 value
	 * @param _extraData : additional data -> [{"moduleId":""}],指明要进行验证的模块
	 * @return 0 for false ,1 for true
	 */
	function authorizeProcessor(address _from, address _to, string _funcHash, string _extraData) public constant returns(uint) {
		log("into authorizeProcessor...", "RoleFilterManager");
		log("authProcess->from:", _from);
		log("authProcess->to:", _to);
		log("authProcess->funcHash:", _funcHash);
		log("authProcess->extraData:", _extraData);
		log("authProcess->msg.sender:", msg.sender);
		log("authProcess->owner:", owner);

		if(_to == address(0)) {
			log("deploying contract, pass");
			return 1;
		}
		if(_from == address(0)) {
			log("from is invalid", "RoleFilterManager");
			return 0;
		}
		if(_from == owner) {
			log("user is owner,valid succ..", "RoleFilterManager");
			return 1;
		}
		if(__userExists(_from) == uint(0)) {
			log("user is not exists...", "RoleFilterManager");
			return 0;
		}

		string memory moduleName = rm.findModuleNameByAddress(_to).recoveryToString();
		string memory moduleVersion = rm.findModuleVersionByAddress(_to).recoveryToString();
		string memory contractName = rm.findResNameByAddress(_to).recoveryToString();
		string memory contractVersion = rm.findContractVersionByAddress(_to).recoveryToString();

		string memory contractKey = moduleName.concat(moduleVersion, contractName, contractVersion);

		uint index = getModuleIndexByNameAndVersion(moduleName, moduleVersion, false);
		if(index == uint(-1)) {
			log("module not exists...", "RoleFilterManager");
			return 0;
		}

		// check whether the module switch is on or not
		if(modules[index].moduleEnable == 0) {
			log("module switch is off.", "RoleFilterManager");
			return 1;
		}

		// check whether the contract switch is on or not
		if(contractMapping[contractKey].deleted) {
			log("contract is invalid ...", "RoleFilterManager");
			return 0;
		}
		if(contractMapping[contractKey].enable == 0) {
			log("contract switch is off..", "RoleFilterManager");
			return 1;
		}

		return __authorizeProcessor(_from, _to, _funcHash, contractKey, index);
	}

	function __authorizeProcessor(address _from, address _to, string _funcHash, string _contractKey, uint _index) private constant returns(uint) { //for kill: Compiler error: Stack too deep, try removing local variables.
		// check the _opSha3Key is exists in openAction of Contract
		// check contract privilege = openAction
		string memory strOpKey = _funcHash.substr(0, 8).toLower();
		for(uint i = 0; i < contractMapping[_contractKey].openActionList.length; ++i) {
			if(contractMapping[_contractKey].openActionList[i].funcHash.substr(0, 8).toLower().equals(strOpKey)) {
				log("funcHash belong openAction..", "RoleFilterManager");
				return 1;
			}
		}

		// check eht action is have to contract actionId
		IActionManager am = IActionManager(rm.getContractAddress("SystemModuleManager", "0.0.1.0", "ActionManager", "0.0.1.0"));
		bool isFlag = false;
		for(i = 0; i < contractMapping[_contractKey].actionIdList.length; ++i) {
			if(am.checkActionWithKey(contractMapping[_contractKey].actionIdList[i], _to, _funcHash) == 1) {
				isFlag = true;
			}
		}
		if(!isFlag) {
			log("action not belong the contract...", "RoleFilterManager");
			return 0;
		}

		// check module is have privilege , then check user is have privilege
		delete tmpArray;
		__getUserRoleIdList(_from, tmpArray);
		log("authpro->tmpArray.length:", tmpArray.length);
		IRoleManager roleMgr = IRoleManager(rm.getContractAddress("SystemModuleManager", "0.0.1.0", "RoleManager", "0.0.1.0"));
		for(i = 0; i < modules[_index].roleIds.length; ++i) {
			if(roleMgr.checkRoleActionWithKey(modules[_index].roleIds[i], _to, _funcHash) != 1) {
				continue;
			}
			// check the action for the roleIds of user
			for(uint j = 0; j < tmpArray.length; ++j) {
				if(roleMgr.checkRoleActionWithKey(tmpArray[j], _to, _funcHash) == 1) {
					log("auth,check in role action ,succ..");
					return 1;
				}
			}
		}
		log("authProcess end,auth fail..", "RoleFilterManager");
		return 0;
	}

	function setModuleEnableDb(string _id, uint _enable) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "setModuleEnableDb");
		log(_id, _enable.toString());
		_ret = writedb("moduleEnable|update", _id, _enable.toString());
		if(0 != _ret)
			log("RoleFilterManager.sol", "setModuleEnableDb failed.");
		else
			log("RoleFilterManager.sol", "setModuleEnableDb success.");
	}

	function setModuleEnableDb(string _moduleName, string _moduleVersion, uint _moduleEnable) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "setModuleEnableDb");
		log(_moduleName, _moduleVersion, _moduleEnable.toString());
		_ret = writedb("moduleEnable|update", _moduleName.concat(_moduleVersion), _moduleEnable.toString());
		if(0 != _ret)
			log("RoleFilterManager.sol", "setModuleEnableDb failed.");
		else
			log("RoleFilterManager.sol", "setModuleEnableDb success.");
	}

	function setContractEnableDb(address _contractAddr, string _moduleId, uint _enable) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "setContractEnableDb");
		log(uint(_contractAddr).toAddrString(), _moduleId, _enable.toString());
		_ret = writedb("contractEnable|update", uint(_contractAddr).toAddrString(), _moduleId.concat("|", _enable.toString()));
		if(0 != _ret)
			log("RoleFilterManager.sol", "setContractEnableDb failed.");
		else
			log("RoleFilterManager.sol", "setContractEnableDb success.");
	}

	function setContractEnableDb(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, uint _contractEnable) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "setContractEnableDb");
		log(_moduleName, _moduleVersion);
		log(_contractName, _contractVersion, _contractEnable.toString());
		string memory key = _moduleName.concat(_moduleVersion, _contractName, _contractVersion);
		_ret = writedb("contractEnable|update", key, _contractEnable.toString());
		if(0 != _ret)
			log("RoleFilterManager.sol", "setContractEnableDb failed.");
		else
			log("RoleFilterManager.sol", "setContractEnableDb success.");
	}

	function addOpenAction(string _moduleId, address _contractAddr, string _funcHash) constant private returns(uint _ret) { //just one funcHash:f6bfc763
		log("RoleFilterManager.sol", "addOpenAction");
		log(_moduleId, uint(_contractAddr).toAddrString(), _funcHash);
		_ret = writedb("openAction|add", _moduleId.concat("|", uint(_contractAddr).toAddrString()), _funcHash);
		if(0 != _ret)
			log("RoleFilterManager.sol", "addOpenAction failed.");
		else
			log("RoleFilterManager.sol", "addOpenAction success.");
	}

	function addOpenAction(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, string _actionId) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "addOpenAction");
		log(_moduleName, _moduleVersion);
		log(_contractName, _contractVersion, _actionId);
		string memory key = _moduleName.concat(_moduleVersion, _contractName, _contractVersion);
		_ret = writedb("openAction|add", key, _actionId);
		if(0 != _ret)
			log("RoleFilterManager.sol", "addOpenAction failed.");
		else
			log("RoleFilterManager.sol", "addOpenAction success.");
	}

	function deleteOpenAction(string _moduleId, address _contractAddr, string _funcHash) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "deleteOpenAction");
		log(_moduleId, uint(_contractAddr).toAddrString(), _funcHash);
		_ret = writedb("openAction|delete", _moduleId.concat("|", uint(_contractAddr).toAddrString()), _funcHash);
		if(0 != _ret)
			log("RoleFilterManager.sol", "deleteOpenAction failed.");
		else
			log("RoleFilterManager.sol", "deleteOpenAction success.");
	}

	function deleteOpenAction(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, string _actionId) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "addOpenAction");
		log(_moduleName, _moduleVersion);
		log(_contractName, _contractVersion, _actionId);
		string memory key = _moduleName.concat(_moduleVersion, _contractName, _contractVersion);
		_ret = writedb("openAction|delete", key, _actionId);
		if(0 != _ret)
			log("RoleFilterManager.sol", "addOpenAction failed.");
		else
			log("RoleFilterManager.sol", "addOpenAction success.");
	}

	function addAuthorizeRole(string _moduleId, string _roleId, string _actionInfo) constant private returns(uint _ret) { //just one actionInfo:0000000000000000000000000000000000000012,78a9eeed
		log("RoleFilterManager.sol", "addAuthorizeRole");
		log(_moduleId, _roleId, _actionInfo);
		_ret = writedb("authorizeRole|add", _moduleId.concat("|", _roleId), _actionInfo);
		if(0 != _ret)
			log("RoleFilterManager.sol", "addAuthorizeRole failed.");
		else
			log("RoleFilterManager.sol", "addAuthorizeRole success.");
	}

	function deleteAuthorizeRole(string _moduleId, string _roleId, string _actionInfo) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "deleteAuthorizeRole");
		log(_moduleId, _roleId, _actionInfo);
		_ret = writedb("authorizeRole|delete", _moduleId.concat("|", _roleId), _actionInfo);
		if(0 != _ret)
			log("RoleFilterManager.sol", "deleteAuthorizeRole failed.");
		else
			log("RoleFilterManager.sol", "deleteAuthorizeRole success.");
	}

	function updateFuncInfoDb(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion, string _func, string _actionId) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "updateFuncInfoDb");
		log(_moduleName, _moduleVersion, _contractName);
		log(_contractVersion, _func, _actionId);
		string memory key = _moduleName.concat(_moduleVersion, _contractName, _contractVersion);
		key = key.concat(_func);
		_ret = writedb("funcInfo|update", key, _actionId);
		if(0 != _ret)
			log("RoleFilterManager.sol", "updateFuncInfoDb failed.");
		else
			log("RoleFilterManager.sol", "updateFuncInfoDb success.");
	}

	function updateAuthorizeRole(string _moduleId, string _roleId, string _actionInfos) constant private returns(uint _ret) { //actionInfos:0000000000000000000000000000000000000012,78a9eeed|0000000000000000000000000000000000000013,78a9eeee
		log("RoleFilterManager.sol", "updateAuthorizeRole");
		log(_moduleId, _roleId, _actionInfos);
		_ret = writedb("authorizeRole|update", _moduleId.concat("|", _roleId), _actionInfos);
		if(0 != _ret)
			log("RoleFilterManager.sol", "updateAuthorizeRole failed.");
		else
			log("RoleFilterManager.sol", "updateAuthorizeRole success.");
	}

	function updateAuthorizeRole(string _moduleName, string _moduleVersion, string _roleId, string _actionInfos) constant private returns(uint _ret) {
		log("RoleFilterManager.sol", "updateAuthorizeRole");
		log(_moduleName, _moduleVersion);
		log(_roleId, _actionInfos);
		string memory key = _moduleName.concat(_moduleVersion, _roleId);
		_ret = writedb("authorizeRole|update", key, _actionInfos);
		if(0 != _ret)
			log("RoleFilterManager.sol", "updateAuthorizeRole failed.");
		else
			log("RoleFilterManager.sol", "updateAuthorizeRole success.");
	}

	/**
	 * Add a new module
	 * @param _json module info
	 * @return 0 for success, else -1 for failed
	 */
	function addModule(string _json) public returns(uint) {
		log("addModule->into...", "RoleFilterManager");
		log("addModule->msg.sender:", msg.sender);
		log("addModule->owner:", owner);
		log("addModule->address:", this);
		log("addModule->tx.origin:", tx.origin);
		uint errno_prefix = 95370;
		uint retUint = uint(-1);
		if(!tmpModule.fromJson(_json)) {
			log("addModule->json invalid", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.JSON_INVALID);
			Notify(errno, "json invalid");
			return retUint;
		}
		if(bytes(tmpModule.moduleName).length == 0) {
			log("addModule->moduleName is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleName is null..");
			return retUint;
		}
		if(bytes(tmpModule.moduleText).length == 0) {
			log("addModule->moduleText is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleText is null..");
			return retUint;
		}
		if(bytes(tmpModule.moduleVersion).length == 0) {
			log("addModule->moduleVersion is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleVersion is null..");
			return retUint;
		}
		uint index = getModuleIndexByNameAndVersion(tmpModule.moduleName, tmpModule.moduleVersion, true);
		if(index == uint(-2)) {
			log("module already exists, but isn't yours, change module owner via admin please.");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module already exists, but isn't yours..");
			return retUint;
		} else if(index == uint(-1)) {
			log("exec insert mode.");
			tmpModule.moduleId = tmpModule.moduleName.concat("_", tmpModule.moduleVersion);
			tmpModule.moduleCreator = tx.origin;
			tmpModule.moduleCreateTime = now * 1000;
			tmpModule.moduleUpdateTime = now * 1000;
			tmpModule.publishTime = now * 1000;
			modules.push(tmpModule);
			setModuleEnableDb(tmpModule.moduleName, tmpModule.moduleVersion, tmpModule.moduleEnable);
			log("addModule success, moduleName:moduleVersion:", tmpModule.moduleName, tmpModule.moduleVersion);
		} else {
			log("exec update mode.");
			if(_json.keyExists("moduleEnable") && modules[index].moduleEnable != tmpModule.moduleEnable) {
				log("update module enable old->new:",
					modules[index].moduleEnable.toString(),
					tmpModule.moduleEnable.toString());
				modules[index].moduleEnable = tmpModule.moduleEnable;
				setModuleEnableDb(tmpModule.moduleName, tmpModule.moduleVersion, tmpModule.moduleEnable);
			}
			if(_json.keyExists("moduleDescription") && !modules[index].moduleDescription.equals(tmpModule.moduleDescription)) {
				log("update module description old->new:",
					modules[index].moduleDescription,
					tmpModule.moduleDescription);
				modules[index].moduleDescription = tmpModule.moduleDescription;
			}
			if(_json.keyExists("moduleText") && !modules[index].moduleText.equals(tmpModule.moduleText)) {
				log("update moduleText old->new:",
					modules[index].moduleText,
					tmpModule.moduleText);
				modules[index].moduleText = tmpModule.moduleText;
			}
			modules[index].moduleUpdateTime = now * 1000;
			log("addModule(update) success, moduleName:moduleVersion:", tmpModule.moduleName, tmpModule.moduleVersion);
		}

		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		Notify(errno, "add module success...");
		return 0;
	}

	/**
	 * @dev update module info
	 * @param _json the json struct for module
	 * @return 0 for false, other for true
	 */
	function updModule(string _json) public returns(uint) {
		log("updModule->into...", "RoleFilterManager");
		log("updModule->msg.sender:", msg.sender);
		log("updModule->owner:", owner);
		log("updModule->address:", this);

		uint errno_prefix = 95370;
		uint retUint = 0;

		if (tx.origin != owner) {
            log("not root user, can not exec updModule func.", "RoleFilterManager");
            errno = errno_prefix + uint(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
            Notify(errno, "only root user can updModule func.");
            return errno;
        }
		if(!tmpModule.fromJson(_json)) {
			log("updModule->json invalid", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.JSON_INVALID);
			Notify(errno, "json invalid");
			return retUint;
		}
		if(bytes(tmpModule.moduleName).length == 0) {
			log("updModule->moduleName invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleName is null...");
			return retUint;
		}

		if(bytes(tmpModule.moduleVersion).length == 0) {
			log("updModule->moduleVersion invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleVersion is null...");
			return retUint;
		}

		tmpModule.moduleId = tmpModule.moduleName.concat("_",tmpModule.moduleVersion);

		// check the moduleId is exists...
		uint index = getIndexById(tmpModule.moduleId);
		if(index == uint(-1)) {
			log("updModule id not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleId is null...");
			return retUint;
		}
		if(_json.keyExists("moduleEnable") && tmpModule.moduleEnable != 0 && tmpModule.moduleEnable != 1) {
			log("updModule->moduleEnable invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleEnable invalid...");
			return retUint;
		} 
		if(_json.keyExists("roleIds")) {
			modules[index].roleIds = tmpModule.roleIds;
		}
		if(_json.keyExists("moduleEnable")) {
			modules[index].moduleEnable = tmpModule.moduleEnable;
			setModuleEnableDb(tmpModule.moduleName, tmpModule.moduleVersion, modules[index].moduleEnable);
		}
		if(_json.keyExists("moduleDescription")) {
			modules[index].moduleDescription = tmpModule.moduleDescription;
		}
		if(_json.keyExists("moduleName")) {
			modules[index].moduleName = tmpModule.moduleName;
		}
		if(_json.keyExists("moduleText")) {
			modules[index].moduleText = tmpModule.moduleText;
		}
		if(_json.keyExists("moduleType")) {
			modules[index].moduleType = tmpModule.moduleType;
		}
		if(_json.keyExists("icon")) {
			modules[index].icon = tmpModule.icon;
		}
		if(_json.keyExists("moduleUrl")) {
			modules[index].moduleUrl = tmpModule.moduleUrl;
		}
		if(_json.keyExists("moduleUpdateTime")) {
			modules[index].moduleUpdateTime = tmpModule.moduleUpdateTime;
		} else {
			modules[index].moduleUpdateTime = now * 1000;
		}
		log("updModule success..", "RoleFilterManager");
		revision++;
		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		Notify(errno, "updModule success..");
		return tmpModule.moduleId.storageToUint();
	}

	/**
	 * @dev delete module info
	 * @param _moduleId module id
	 * @return 0 for false,other for true(eg:moduleId)
	 */
	function delModule(string _moduleId) public returns(uint) {
		log("delModule->into...", "RoleFilterManager");
		log("delModule->msg.sender:", msg.sender);
		log("delModule->owner:", owner);
		log("delModule->address:", this);
		log("delModule->moduleId:", _moduleId);

		uint errno_prefix = 95370;
		uint retUint = 0;
		if(bytes(_moduleId).length == 0) {
			log("delModule->_moduleId is null", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "param invalid");
			return retUint;
		}

		// check the moduleId is exists...
		uint index = getIndexById(_moduleId);
		if(index == uint(-1)) {
			log("delModule id not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleId is not exists ...");
			return retUint;
		}

		// do delete
		modules[index].deleted = true;
		setModuleEnableDb(modules[index].moduleId, 0);
		revision++;
		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		log("delModule success..", "RoleFilterManager");

		Notify(errno, "del module success...");
		return _moduleId.storageToUint();
	}

	/**
	 * enable the module switch
	 * if set enalbe to 0,then do not authentication
	 * if set enable to 1,then do authentication
	 * @param _moduleId module id
	 * @param _enable 1 for on,0 for off
	 * @return 0 for false,1 for true
	 */
	/* function setModuleEnable(string _moduleId, uint _enable) public returns(uint) {
		log("setModuleEnable->into...", "RoleFilterManager");
		log("setModuleEnable->moduleId:", _moduleId);

		uint errno_prefix = 95370;
		uint retUint = 0;
		if(bytes(_moduleId).length == 0) {
			log("setModuleEnable->_moduleId is null");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "param invalid");
			return retUint;
		}

		// check the moduleId is exists...
		uint index = getIndexById(_moduleId);
		if(index == uint(-1)) {
			log("setModuleEnable id not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleId is not exists ...");
			return retUint;
		}

		// check enable is invalid
		if(_enable != 0 && _enable != 1) {
			log("setModuleEnable value invalid...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "enable value invalid ...");
			return retUint;
		}

		// do delete
		modules[index].moduleEnable = _enable;
		setModuleEnableDb(modules[index].moduleId, modules[index].moduleEnable);
		revision++;
		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		log("set enable success..", "RoleFilterManager");

		Notify(errno, "set enable success..");
		return 1;
	} */

	/**
	 * enable the contract switch
	 * if set enalbe to 0,then do not authentication
	 * if set enable to 1,then do authentication
	 * @param _contractId module id
	 * @param _enable 1 for on,0 for off
	 * @return 0 for false,1 for true
	 */
	/* function setConntractEnable(string _contractId, uint _enable) public returns(uint) {
		log("setConntractEnable->into...", "RoleFilterManager");
		log("setConntractEnable->_contractId:", _contractId);

		uint errno_prefix = 95370;
		uint retUint = 0;
		if(bytes(_contractId).length == 0) {
			log("setConntractEnable->_mod_contractIduleId is null");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "param invalid");
			return retUint;
		}

		// check enable is invalid
		if(_enable != 0 && _enable != 1) {
			log("setConntractEnable value invalid...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "enable value invalid ...");
			return retUint;
		}

		// set enable
		contractMapping[_contractId].enable = _enable;
		setContractEnableDb(contractMapping[_contractId].cctAddr, contractMapping[_contractId].moduleId, _enable);
		revision++;
		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		log("setContract enable success..", "RoleFilterManager");

		Notify(errno, "setContract enable success..");
		return 1;
	} */

	/**
	 * @dev Add Contract
	 * @param _json the json struct for contract
	 * @return 0 for false,other for true (eg:contractId)
	 */
	function addContract(string _json) public returns(uint) {
		log("addContract->into...", "RoleFilterManager");
		log("addContract->msg.sender:", msg.sender);
		log("addContract->owner:", owner);
		log("addContract->address:", this);
		log("addContract->tx.origin:", tx.origin);

		uint errno_prefix = 95370;
		uint retUint = uint(-1);
		if(!tmpContract.fromJson(_json)) {
			log("addContract->json invalid", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.JSON_INVALID);
			Notify(errno, "json invalid");
			return retUint;
		}
		if(bytes(tmpContract.moduleName).length == 0) {
			log("addContract->moduleName is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleName is null..");
			return retUint;
		}
		if(bytes(tmpContract.moduleVersion).length == 0) {
			log("addContract->moduleVersion is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleVersion is null..");
			return retUint;
		}
		if(bytes(tmpContract.cctName).length == 0) {
			log("addContract-> cctName is null...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "cctName is null...");
			return retUint;
		}
		if(bytes(tmpContract.cctVersion).length == 0) {
			log("addContract-> cctVersion is null...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "cctVersion is null...");
			return retUint;
		}

		// check openAction is repeat
		for(uint i = 0; i < tmpContract.openActionList.length; ++i) {
			for(uint j = 0; j < tmpContract.openActionList.length; ++i) {
				if(i == j) {
					continue;
				}
				if(tmpContract.openActionList[i].id.equals(tmpContract.openActionList[j].id)) {
					log("openAction id repeat...", "RoleFilterManager");
					errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
					Notify(errno, "openAction id repeat is null...");
					return retUint;
				}
			}
		}

		string memory contractKey = tmpContract.moduleName.concat(tmpContract.moduleVersion, tmpContract.cctName, tmpContract.cctVersion);

		uint index = getModuleIndexByNameAndVersion(tmpContract.moduleName, tmpContract.moduleVersion, true);
		if(index == uint(-2)) {
			log("module already exists, but isn't yours, change module owner via admin please.");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module already exists, but isn't yours..");
			return retUint;
		} else if(index == uint(-1)) {
			log("addContract module not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module is not exists...");
			return retUint;
		} else { //module exist, and can be insert or update contract
			uint enable = getContractEnableByNameAndVersion(tmpContract.moduleName, tmpContract.moduleVersion, tmpContract.cctName, tmpContract.cctVersion);
			if(enable == uint(-2)) {
				log("contract already exists, but isn't yours, change contract owner via admin please.");
				errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
				Notify(errno, "contract already exists, but isn't yours..");
				return retUint;
			} else if(enable == uint(-1)) {
				log("exec insert mode.");
				tmpContract.moduleId = tmpContract.moduleName.concat("_",tmpContract.moduleVersion);
				tmpContract.cctId = tmpContract.moduleName.concat("_",tmpContract.moduleVersion).concat("_",tmpContract.cctName).concat("_",tmpContract.cctVersion);
				tmpContract.creator = tx.origin;
				tmpContract.createTime = now * 1000;
				tmpContract.updateTime = now * 1000;

				// save contract info
				contractMapping[contractKey] = tmpContract;
				setContractEnableDb(tmpContract.moduleName, tmpContract.moduleVersion, tmpContract.cctName, tmpContract.cctVersion, tmpContract.enable);

				// save contract to modules
				LibContract.ContractInfo memory contractInfo;
				contractInfo.contractId = tmpContract.cctId;
				contractInfo.contractName = tmpContract.cctName;
				contractInfo.contractVersion = tmpContract.cctVersion;
				modules[index].contractInfos.push(contractInfo);
				log("addContract success, contractName:contractVersion:", tmpContract.cctName, tmpContract.cctVersion);
			} else {
				log("exec update mode.");
				if(_json.keyExists("enable") && enable != tmpContract.enable) {
					log("update contract enable old->new:",
						enable.toString(),
						tmpContract.enable.toString());
					contractMapping[contractKey].enable = tmpContract.enable;
					setContractEnableDb(tmpContract.moduleName, tmpContract.moduleVersion, tmpContract.cctName, tmpContract.cctVersion, tmpContract.enable);
				}
				contractMapping[contractKey].updateTime = now * 1000;
				log("addContract(update) success, contractName:contractVersion:", tmpContract.cctName, tmpContract.cctVersion);
			}

			errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
			Notify(errno, "add contract success...");
			return 0;
		}
	}

	/**
	 * Add Menu
	 * @param _json the json struct for action
	 * @return 0 for false,other for true(eg:menuId)
	 */
	function addMenu(string _json) public returns(uint) {
		log("addMenu->into...", "RoleFilterManager");
		log("addMenu->msg.sender:", msg.sender);
		log("addMenu->owner:", owner);
		log("addMenu->address:", this);
		log("addMenu->tx.origin:", tx.origin);

		uint errno_prefix = 95370;
		uint retUint = uint(-1);
		// get action obj and set actionId
		tmpAction.jsonParse(_json);
		if(bytes(tmpAction.moduleName).length == 0) {
			log("addMenu->moduleName is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleName is null..");
			return retUint;
		}
		if(bytes(tmpAction.moduleVersion).length == 0) {
			log("addMenu->moduleVersion is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleVersion is null..");
			return retUint;
		}
		/* if(tmpAction.level != 1 && tmpAction.level != 2) {
			log("addMenu->level is invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "level is invalid..");
			return retUint;
		} */
		if(tmpAction.Type != 1) {
			log("addMenu->Type is invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "Type is invalid..");
			return retUint;
		}

		uint index = getModuleIndexByNameAndVersion(tmpAction.moduleName, tmpAction.moduleVersion, true);
		if(index == uint(-2)) {
			log("module already exists, but isn't yours, change module owner via admin please.");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module already exists, but isn't yours..");
			return retUint;
		} else if(index == uint(-1)) {
			log("addMenu module not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module is not exists...");
			return retUint;
		} else { //module exist, and can be insert or update menu
			// do insert oper - invoke ActionManager.insert();
			address amAddress = rm.getContractAddress("SystemModuleManager", "0.0.1.0", "ActionManager", "0.0.1.0");
			if(amAddress == 0) {
				log("addMenu->ActionManager addr not exists", "RoleFilterManager");
				errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
				Notify(errno, "ActionManager addr not exists");
				return retUint;
			}
			IActionManager am = IActionManager(amAddress);

			uint isExists = am.actionExists(tmpAction.id);
			if(isExists == 0) {
				log("exec insert mode.");
				tmpAction.moduleId = tmpAction.moduleName.concat("_",tmpAction.moduleVersion);
				tmpAction.state = LibAction.ActionState.VALID;
				tmpAction.createTime = now * 1000;
				tmpAction.updateTime = now * 1000;
				tmpAction.creator = tx.origin;
				if(!am.insert(tmpAction.toJson())) {
					log("addMenu->am insert fail..", "RoleFilterManager");
					errno = uint256(ROLE_FILTER_ERRORCODE.SYS_EXCEPTION);
					Notify(errno, "am insert fail..");
					return retUint;
				}
				log("add menu success...");
			} else {
				log("exec update mode.");
				tmpAction.state = LibAction.ActionState.VALID;
				tmpAction.updateTime = now * 1000;
				if(!am.update(tmpAction.toJson())) {
					log("addMenu->am update fail..", "RoleFilterManager");
					errno = uint256(ROLE_FILTER_ERRORCODE.SYS_EXCEPTION);
					Notify(errno, "am update fail..");
					return retUint;
				}
				log("add menu(update) success...");
			}
		}

		revision++;

		log("add menu success...");
		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		Notify(errno, "add menu success...");
		return 0;
	}

	/**
	 * Add Action
	 * @param _json the json struct for action
	 * @return 0 for false,other for true(eg:menuId)
	 */
	function addAction(string _json) public returns(uint) {
		log("addAction->into...", "RoleFilterManager");
		log("addAction->msg.sender:", msg.sender);
		log("addAction->owner:", owner);
		log("addAction->address:", this);
		log("addAction->tx.origin:", tx.origin);
		uint errno_prefix = 95370;
		uint retUint = uint(-1);
		// get action obj and set actionId
		tmpAction.jsonParse(_json);
		log("+++++++ addAction ,the action ID +++++++++:" , tmpAction.id);
		if(bytes(tmpAction.moduleName).length == 0) {
			log("addAction->moduleName is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleName is null..");
			return retUint;
		}
		if(bytes(tmpAction.moduleVersion).length == 0) {
			log("addAction->moduleVersion is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleVersion is null..");
			return retUint;
		}
		if(tmpAction.level != 3) {
			log("addAction->level is invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "level is invalid..");
			return retUint;
		}
		if(tmpAction.Type != 2) {
			log("addAction->Type is invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "Type is invalid..");
			return retUint;
		}
		if(tmpAction.enable != 0 && tmpAction.enable != 1) {
			log("addAction->enable is invalid..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "enable is invalid..");
			return retUint;
		}
		if(bytes(tmpAction.resKey).length == 0) {
			log("addAction->contractName is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "contractName is null..");
			return retUint;
		}
		if(bytes(tmpAction.version).length == 0) {
			log("addAction->contractVersion is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "contractVersion is null..");
			return retUint;
		}
		if(bytes(tmpAction.opKey).length == 0) {
			log("addAction->opKey is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "opKey is null..");
			return retUint;
		}

		uint index = getModuleIndexByNameAndVersion(tmpAction.moduleName, tmpAction.moduleVersion, true);
		if(index == uint(-2)) {
			log("module already exists, but isn't yours, change module owner via admin please.");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module already exists, but isn't yours..");
			return retUint;
		} else if(index == uint(-1)) {
			log("addAction module not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module is not exists...");
			return retUint;
		} else { //module exist, and can be insert or update action
			// do insert oper - invoke ActionManager.insert();
			address amAddress = rm.getContractAddress("SystemModuleManager", "0.0.1.0", "ActionManager", "0.0.1.0");
			if(amAddress == 0) {
				log("addAction->ActionManager addr not exists", "RoleFilterManager");
				errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
				Notify(errno, "ActionManager addr not exists");
				return retUint;
			}
			IActionManager am = IActionManager(amAddress);

			string memory contractKey = tmpAction.moduleName.concat(tmpAction.moduleVersion, tmpAction.resKey, tmpAction.version);

			string memory funNameSha3 = uint(sha3(tmpAction.opKey)).toHexString64().toLower().substr(2, 8);

			uint oldEnable = am.queryActionEnable(tmpAction.id, false);
			if(oldEnable == uint(-1)) {
				log("exec insert mode.");
				//tmpAction.state = LibAction.ActionState.VALID;
				//tmpAction.createTime = now * 1000;
				//tmpAction.updateTime = now * 1000;
				//tmpAction.creator = tx.origin;
				//if(!am.insert(tmpAction.toJson())) {
				tmpAction.moduleId = tmpAction.moduleName.concat("_",tmpAction.moduleVersion);
				tmpAction.contractId = tmpAction.moduleName.concat("_",tmpAction.moduleVersion).concat("_",tmpAction.resKey).concat("_",tmpAction.version);
				if(!am.insert(_json)) {
					log("addAction->am insert fail..", "RoleFilterManager");
					errno = uint256(ROLE_FILTER_ERRORCODE.SYS_EXCEPTION);
					Notify(errno, "am insert fail..");
					return retUint;
				}

				// set action to actionIdList of contract
				contractMapping[contractKey].actionIdList.push(tmpAction.id);

				// set action to openAction
				if(tmpAction.enable == 0) {
					log("set openAction..", "RoleFilterManager");
					tmpOpenAction.id = tmpAction.id;
					tmpOpenAction.funcHash = funNameSha3;
					contractMapping[contractKey].openActionList.push(tmpOpenAction);
					addOpenAction(tmpAction.moduleName, tmpAction.moduleVersion, tmpAction.resKey, tmpAction.version, tmpAction.id);
					log("set openAction succ..");
				}
				log("add action success...");
			} else {
				log("exec update mode.");
				//tmpAction.state = LibAction.ActionState.VALID;
				//tmpAction.updateTime = now * 1000;
				//if(!am.update(tmpAction.toJson())) {
				if(!am.update(_json)) {
					log("addAction->am update fail..", "RoleFilterManager");
					errno = uint256(ROLE_FILTER_ERRORCODE.SYS_EXCEPTION);
					Notify(errno, "am update fail..");
					return retUint;
				}
				//just update enable!!!
				if(oldEnable != tmpAction.enable) {
					log("update action enable old->new:",
						oldEnable.toString(),
						tmpAction.enable.toString());
					if(0 == tmpAction.enable) { //1->0, close to open
						log("addAction->set openAction..", "RoleFilterManager");
						tmpOpenAction.id = tmpAction.id;
						tmpOpenAction.funcHash = funNameSha3;
						contractMapping[contractKey].openActionList.push(tmpOpenAction);
						addOpenAction(tmpAction.moduleName, tmpAction.moduleVersion, tmpAction.resKey, tmpAction.version, tmpAction.id);
						log("set openAction succ..");
					} else { //0->1, open to close
						log("addAction->unset openAction..", "RoleFilterManager");
						uint len = contractMapping[contractKey].openActionList.length;
						uint idx = 0;
						for(uint i = 0; i < len; ++i) {
							if(!contractMapping[contractKey].openActionList[i].id.equals(tmpAction.id)) {
								contractMapping[contractKey].openActionList[idx++].id = contractMapping[contractKey].openActionList[i].id;
								contractMapping[contractKey].openActionList[idx].funcHash = contractMapping[contractKey].openActionList[i].funcHash;
							}
						}
						for(i = 0; i < (len - idx); ++i) {
							delete contractMapping[contractKey].openActionList[len - 1 - i];
						}
						contractMapping[contractKey].openActionList.length = contractMapping[contractKey].openActionList.length - (len - idx);
						deleteOpenAction(tmpAction.moduleName, tmpAction.moduleVersion, tmpAction.resKey, tmpAction.version, tmpAction.id);
						log("unset openAction succ..");
					}
				}
				log("add action(update) success...");
			}
			updateFuncInfoDb(tmpAction.moduleName, tmpAction.moduleVersion, tmpAction.resKey, tmpAction.version, funNameSha3, tmpAction.id);
		}

		revision++;
		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		Notify(errno, "add action success...");
		return 0;
	}

	/**
	 * Add Role
	 * @param _json the json struct for Role
	 * @return 0 for false,other for true(eg:roleId)
	 */
	function addRole(string _json) public returns(uint) {
		log("addRole->into...", "RoleFilterManager");
		log("addRole->msg.sender:", msg.sender);
		log("addRole->owner:", owner);
		log("addRole->address:", this);
		log("addRole->tx.origin:", tx.origin);

		uint errno_prefix = 95370;
		uint retUint = uint(-1);
		if(!tmpRole.fromJson(_json)) {
			log("addRole->json invalid", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.JSON_INVALID);
			Notify(errno, "json invalid");
			return retUint;
		}
		if(bytes(tmpRole.moduleName).length == 0) {
			log("addRole->moduleName is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleName is null..");
			return retUint;
		}
		if(bytes(tmpRole.moduleVersion).length == 0) {
			log("addRole->moduleVersion is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleVersion is null..");
			return retUint;
		}
		tmpRole.moduleId = tmpRole.moduleName.concat("_", tmpRole.moduleVersion);
		uint index = getModuleIndexByNameAndVersion(tmpRole.moduleName, tmpRole.moduleVersion, true);
		if(index == uint(-2)) {
			log("module already exists, but isn't yours, change module owner via admin please.");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module already exists, but isn't yours..");
			return retUint;
		} else if(index == uint(-1)) {
			log("addRole module not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module is not exists...");
			return retUint;
		} else { //module exist, and can be insert or update role
			address roleAddress = rm.getContractAddress("SystemModuleManager", "0.0.1.0", "RoleManager", "0.0.1.0");
			if(roleAddress == 0) {
				log("addRole,roleAddress is null ");
				errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
				Notify(errno, "roleAddress is null");
				return retUint;
			}
			IRoleManager roleRmanager = IRoleManager(roleAddress);

			bool isExists = tmpRole.id.inArray(modules[index].roleIds);
			uint isExistsEx = roleRmanager.roleExists(tmpRole.id); //roleRmanager.roleExistsEx(tmpRole.id);
			if(!isExists) {
				if(isExistsEx != 0) {
					log("addRole, role existence conflicted");
					errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
					Notify(errno, "role existence conflicted");
					return retUint;
				}
				log("exec insert mode.");
				// set role state = 1 valid
				tmpRole.status = 1;
				uint ret = roleRmanager.insert(tmpRole.toJson());
				if(ret != 0) {
					log("addRole,role insert fail..");
					errno = uint256(ROLE_FILTER_ERRORCODE.SYS_EXCEPTION);
					Notify(errno, "role insert fail..");
					return retUint;
				}
				// set roleID to modules
				modules[index].roleIds.push(tmpRole.id);
				log("add role success...");
			} else {
				if(isExistsEx == 0) {
					log("addRole, role existence conflicted222");
					errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
					Notify(errno, "role existence conflicted");
					return retUint;
				} else {
					log("exec update mode.");
					roleRmanager.update(tmpRole.toJson());
					log("add role(update) success...");
				}
			}
			// writedb role define
			string memory actionInfos;
			log("prepare updateAuthorizeRole", tmpRole.id, tmpRole.actionIdList.length.toString());
			for(uint i = 0; i < tmpRole.actionIdList.length; ++i) {
				if(i > 0)
					actionInfos = actionInfos.concat("|");
				actionInfos = actionInfos.concat(tmpRole.actionIdList[i]);
			}
			//log("actionInfos:", actionInfos);
			updateAuthorizeRole(tmpRole.moduleName, tmpRole.moduleVersion, tmpRole.id, actionInfos);

			revision++;
			errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
			Notify(errno, "add role success...");
			return 0;
		}
	}

	/**
	 * change module owner, and all contracts owner
	 * @param _moduleName moduleName
	 * @param _moduleVersion moduleVersion
	 * @param _newOwner new owner
	 * @return 0 for success, else -1 for failed
	 */
	function changeModuleOwner(string _moduleName, string _moduleVersion, address _newOwner) public returns(uint) {
		log("changeModuleOwner->into...", "RoleFilterManager");
		log("changeModuleOwner->moduleName and version:", _moduleName, _moduleVersion);
		log("changeModuleOwner->owner:", owner);
		log("changeModuleOwner->newOwner:", _newOwner);
		log("changeModuleOwner->tx.origin:", tx.origin);

		uint errno_prefix = 96370;
		uint retUint = uint(-1);

		if(bytes(_moduleName).length == 0) {
			log("changeModuleOwner->moduleName is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleName is null..");
			return retUint;
		}
		if(bytes(_moduleVersion).length == 0) {
			log("changeModuleOwner->moduleVersion is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "moduleVersion is null..");
			return retUint;
		}

		uint index = getModuleIndexByNameAndVersion(_moduleName, _moduleVersion, true);
		if(index == uint(-1)) {
			log("changeModuleOwner module not exists...", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "module is not exists...");
			return retUint;
		} else {
			if(tx.origin != owner) {
				log("tx.origin isn't administrator, cann't change module owner.", "RoleFilterManager");
				errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
				Notify(errno, "tx.origin isn't administrator, cann't change module owner.");
				return retUint;
			}
			if(modules[index].moduleCreator == _newOwner) {
				log("changeModuleOwner owner no change.", "RoleFilterManager");
				errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
				Notify(errno, "owner no change.");
				return 0;
			}

			log("change module owner from->to:", uint(modules[index].moduleCreator).toAddrString(), uint(_newOwner).toAddrString());
			modules[index].moduleCreator = _newOwner;
			changeModuleRegisterOwner(_moduleName, _moduleVersion, _newOwner);

			string memory contractKey;
			for(uint i = 0; i < modules[index].contractInfos.length; ++i) {
				contractKey = _moduleName.concat(_moduleVersion, modules[index].contractInfos[i].contractName, modules[index].contractInfos[i].contractVersion);
				log("contractKey:", contractKey);
				log("change contract owner from->to:", uint(contractMapping[contractKey].creator).toAddrString(), uint(_newOwner).toAddrString());
				contractMapping[contractKey].creator = _newOwner;
				changeContractRegisterOwner(_moduleName, _moduleVersion, modules[index].contractInfos[i].contractName, modules[index].contractInfos[i].contractVersion, _newOwner);
			}

			revision++;
			errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
			Notify(errno, "changeModuleOwner success...");
			return 0;
		}
	}

	// =========================== the api for DAPP ========================

	function listAll() constant public returns(string _json) {
		_json = listBy(0, "", "");
	}

	function qryModules() constant public returns(string _json) {
		_json = listBy(0, "", "");
	}

	function qryModuleDetail(string _moduleId) constant public returns(string _json) {
		_json = listBy(1, _moduleId);
	}

	function qryModuleDetail(string _moduleName, string _moduleVersion) constant public returns(string _json) {
		_json = listBy(1, _moduleName, _moduleVersion);
	}

	function findByName(string _name) constant public returns (string _json) {
		return findByModuleText(_name);
	}

	function findByModuleText(string _moduleText) constant public returns(string _json) {
		log("findByModuleText->_moduleText:", _moduleText);
		uint len = 0;
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret",uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.appendKeyValue("total", getModuleCount());
		len = LibStack.append(",\"items\":[");
		uint n = 0;
		for(uint i = 0; i < modules.length; ++i) {
			if((bytes(_moduleText).length == 0 || modules[i].moduleText.indexOf(_moduleText) != -1) && !modules[i].deleted) {
				if(n > 0) {
					len = LibStack.append(",");
				}
				len = LibStack.append(modules[i].toJson());
				n++;
			}
		}
		len = LibStack.append("]}}");
		_json = LibStack.popex(len);
	}



	function findContractByModName(string _name) constant public returns (string _json) {
		return findContractByModText(_name);
	}

	function findContractByModText(string _moduleText) constant public returns(string _json) {
		log("findContractByModText->_moduleText:", _moduleText);
		uint len = 0;
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.append("\"items\":[");
		uint total = 0;
		uint n = 0;
		string memory contractKey;
		for(uint i = 0; i < modules.length; ++i) {
			if((bytes(_moduleText).length == 0 || modules[i].moduleText.indexOf(_moduleText) != -1) && !modules[i].deleted) {

				for(uint j = 0; j < modules[i].contractInfos.length; ++j) {
					contractKey = modules[i].moduleName.concat(modules[i].moduleVersion, modules[i].contractInfos[j].contractName, modules[i].contractInfos[j].contractVersion);
					if(contractMapping[contractKey].deleted) {
						continue;
					}
					if(n > 0) {
						len = LibStack.append(",");
					}
					len = LibStack.append(contractMapping[contractKey].toJson());
					total++;
					n++;
				}
			}
		}
		len = LibStack.append("]");
		len = LibStack.appendKeyValue("total",total);
		len = LibStack.append("}}");
		_json = LibStack.popex(len);
	}

	function findContractByModName(string _moduleName, string _moduleVersion) constant public returns(string _json) {
		log("findContractByModName->_moduleName:", _moduleName);
		uint len = 0;
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.appendKeyValue("total",getModuleCount() );
		len = LibStack.append(",\"items\":[");
		uint n = 0;
		string memory contractKey;
		for(uint i = 0; i < modules.length; ++i) {
			if((bytes(_moduleName).length == 0 || modules[i].moduleName.indexOf(_moduleName) != -1) && !modules[i].deleted) {
				for(uint j = 0; j < modules[i].contractInfos.length; ++j) {
					contractKey = _moduleName.concat(_moduleVersion, modules[i].contractInfos[j].contractName, modules[i].contractInfos[j].contractVersion);
					if(contractMapping[contractKey].deleted) {
						continue;
					}
					if(n > 0) {
						len = LibStack.append(",");
					}
					len = LibStack.append(contractMapping[contractKey].toJson());
					n++;
				}
			}
		}
		len = LibStack.append("]}}");
		_json = LibStack.popex(len);
	}

	function listContractByModNameAndCttName(string _moduleName,string _cttName,uint _pageNum,uint _pageSize) constant public returns (string _json) {
		return listContractByModTextAndCttName(_moduleName,_cttName,_pageNum, _pageSize);
	}

	/**
	 * list elem by condition
	 * @param _moduleText module text
	 * @param _cttName contract name
	 * @return _json string
	 */
	function listContractByModTextAndCttName(string _moduleText, string _cttName, uint _pageNum, uint _pageSize) constant public returns(string _json) {
		log("listContractByModTextAndCttName->_moduleText:", _moduleText);
		log("listContractByModTextAndCttName->_cttName:", _cttName);
		uint len = 0;
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.append("\"items\":[");
		uint total = 0;
		uint n = 0;
		uint m = 0;
		string memory contractKey;
		for(uint i = 0; i < modules.length; ++i) {
			if((bytes(_moduleText).length == 0 || modules[i].moduleText.indexOf(_moduleText) != -1) && !modules[i].deleted) {

				for(uint j = 0; j < modules[i].contractInfos.length; ++j) {
					contractKey = modules[i].moduleName.concat(modules[i].moduleVersion, modules[i].contractInfos[j].contractName, modules[i].contractInfos[j].contractVersion);
					if((bytes(_cttName).length == 0 || contractMapping[contractKey].cctName.indexOf(_cttName) != -1) && !contractMapping[contractKey].deleted) {

						if(n >= _pageNum * _pageSize && n < (_pageNum + 1) * _pageSize) {
							if(m > 0) {
								len = LibStack.append(",");
							}
							len = LibStack.append(contractMapping[contractKey].toJson());

							m++;
						}

						if(n >= (_pageNum + 1) * _pageSize) {
							break;
						}
						n++;
					}
				}
			}
		}

		for(i = 0; i < modules.length; ++i) {
			if((bytes(_moduleText).length == 0 || modules[i].moduleText.indexOf(_moduleText) != -1) && !modules[i].deleted) {

				for(j = 0; j < modules[i].contractInfos.length; ++j) {
					contractKey = modules[i].moduleName.concat(modules[i].moduleVersion, modules[i].contractInfos[j].contractName, modules[i].contractInfos[j].contractVersion);
					if((bytes(_cttName).length == 0 || contractMapping[contractKey].cctName.indexOf(_cttName) != -1) && !contractMapping[contractKey].deleted) {
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

	/**
	 * @dev 获取有效的DAPP数量
	 * @return _count 返回总数
	 */
	function getModuleCount() constant public returns(uint _count) {
		for(uint i = 0; i < modules.length; ++i) {
			if(!modules[i].deleted) {
				_count++;
			}
		}
	}

	/**
	 * list elem by condition
	 * @param _cond for the condition
	 *        0 for all
	 *        1 for id
	 *        2 for name
	 * @param _value The condition value
	 * @return _json string
	 */
	function listBy(uint _cond, string _value) constant private returns(string _json) {
		uint tatal = 0;
		uint len = 0;
		for(uint i = 0; i < modules.length; ++i) {
			if(!modules[i].deleted) {
				tatal++;
			}
		}
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.appendKeyValue("total", tatal);
		len = LibStack.append(",\"items\":[");

		uint n = 0;
		for(i = 0; i < modules.length; ++i) {
			if(modules[i].deleted) {
				continue;
			}
			bool suitable = false;
			if(_cond == 0) {
				suitable = true;
			} else if(_cond == 1) {
				if(modules[i].moduleId.equals(_value)) {
					suitable = true;
				}
			}
			if(suitable) {
				if(n > 0) {
					len = LibStack.append(",");
				}
				len = LibStack.append(modules[i].toJson());
				n++;
			}
		}
		len = LibStack.append("]}}");
		_json = LibStack.popex(len);
	}

	function listBy(uint _cond, string _value, string _ver) constant private returns(string _json) {
		uint tatal = 0;
		uint len = 0;
		for(uint i = 0; i < modules.length; ++i) {
			if(!modules[i].deleted) {
				tatal++;
			}
		}
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.appendKeyValue("total", tatal);
		len = LibStack.append(",\"items\":[");

		uint n = 0;
		for(i = 0; i < modules.length; ++i) {
			if(modules[i].deleted) {
				continue;
			}
			bool suitable = false;
			if(_cond == 0) {
				suitable = true;
			} else if(_cond == 1) {
				if(modules[i].moduleName.equals(_value) && modules[i].moduleVersion.equals(_ver)) {
					suitable = true;
				}
			}
			if(suitable) {
				if(n > 0) {
					len = LibStack.append(",");
				}
				len = LibStack.append(modules[i].toJson());
				n++;
			}
		}
		len = LibStack.append("]}}");
		_json = LibStack.popex(len);
	}

	function listContractByModuleId(string _moduleId) constant public returns(string _json) {
		uint len = 0;
		uint index = getIndexById(_moduleId);
		if(index == uint(-1)) {
			return __getNullResult();
		}
		string memory contractKey;
		uint tatal = 0;
		for(uint i = 0; i < modules[index].contractInfos.length; ++i) {
			contractKey = modules[index].moduleName.concat(modules[index].moduleVersion, modules[index].contractInfos[i].contractName, modules[index].contractInfos[i].contractVersion);
			if(!contractMapping[contractKey].deleted) {
				tatal++;
			}
		}
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.appendKeyValue("total", tatal);
		len = LibStack.append(",\"items\":[");

		uint n = 0;
		for(i = 0; i < modules[index].contractInfos.length; ++i) {
			contractKey = modules[index].moduleName.concat(modules[index].moduleVersion, modules[index].contractInfos[i].contractName, modules[index].contractInfos[i].contractVersion);
			if(contractMapping[contractKey].deleted) {
				continue;
			}
			if(n > 0) {
				len = LibStack.append(",");
			}
			len = LibStack.append(contractMapping[contractKey].toJson());
			n++;
		}
		len = LibStack.append("]}}");
		_json = LibStack.popex(len);
	}

	function listContractByModuleName(string _moduleName, string _moduleVersion) constant public returns(string _json) {
		uint index = getModuleIndexByNameAndVersion(_moduleName, _moduleVersion, false);
		uint len = 0;
		if(index == uint(-1)) {
			return __getNullResult();
		}
		string memory contractKey;
		uint tatal = 0;
		for(uint i = 0; i < modules[index].contractInfos.length; ++i) {
			contractKey = _moduleName.concat(_moduleVersion, modules[index].contractInfos[i].contractName, modules[index].contractInfos[i].contractVersion);
			if(!contractMapping[contractKey].deleted) {
				tatal++;
			}
		}
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.appendKeyValue("total", tatal);
		len = LibStack.append(",\"items\":[");

		uint n = 0;
		for(i = 0; i < modules[index].contractInfos.length; ++i) {
			contractKey = _moduleName.concat(_moduleVersion, modules[index].contractInfos[i].contractName, modules[index].contractInfos[i].contractVersion);
			if(contractMapping[contractKey].deleted) {
				continue;
			}
			if(n > 0) {
				len = LibStack.append(",");
			}
			len = LibStack.append(contractMapping[contractKey].toJson());
			n++;
		}
		len = LibStack.append("]}}");
		_json = LibStack.popex(len);
	}

	// --------------------------- invoke for current contract -----------------------------

	function getIndexById(string _id) constant private returns(uint) {
		for(uint i = 0; i < modules.length; ++i) {
			if(modules[i].deleted) {
				continue;
			}
			if(modules[i].moduleId.equals(_id)) {
				return i;
			}
		}
		return uint(-1);
	}

	function getModuleIndexByNameAndVersion(string _name, string _version, bool _checkOwner) constant private returns(uint) {
		for(uint i = 0; i < modules.length; ++i) {
			if(modules[i].deleted) {
				continue;
			}
			if(modules[i].moduleName.equals(_name) && modules[i].moduleVersion.equals(_version)) {
				if(_checkOwner && modules[i].moduleCreator != tx.origin) {
					log("module index:moduleCreator: module isn't yours.", i.toString(), uint(modules[i].moduleCreator).toAddrString());
					return uint(-2);
				} else {
					return i;
				}
			}
		}
		return uint(-1);
	}

	function getContractEnableByNameAndVersion(string _moduleName, string _moduleVersion, string _contractName, string _contractVersion) constant private returns(uint) {
		string memory contractKey = _moduleName.concat(_moduleVersion, _contractName, _contractVersion);
		tmpContractFind = contractMapping[contractKey];
		if(bytes(tmpContractFind.cctName).length == 0) {
			return uint(-1);
		} else if(tmpContractFind.creator != tx.origin) {
			return uint(-2);
		} else {
			return tmpContractFind.enable;
		}
	}

	function __userExists(address _userAddr) constant internal returns(uint) {
		address userManagerAddr = rm.getContractAddress("SystemModuleManager", "0.0.1.0", "UserManager", "0.0.1.0");
		if(userManagerAddr == 0) {
			return 0;
		}
		IUserManager userManager = IUserManager(userManagerAddr);
		uint ret = userManager.userExists(_userAddr);
		log("__userExists(string):", ret);
		return ret;
	}

	function __getUserRoleIdList(address _userAddr, string[] storage roleIdList) constant internal {
		address userManagerAddr = rm.getContractAddress("SystemModuleManager", "0.0.1.0", "UserManager", "0.0.1.0");
		if(userManagerAddr == 0) {
			return;
		}

		// get UserManager instance
		IUserManager userManager = IUserManager(userManagerAddr);
		uint i = 0;
		while(true) {
			uint ret = userManager.getUserRoleId(_userAddr, i);
			if(ret == 0) {
				return;
			} else {
				roleIdList.push(ret.recoveryToString());
			}
			i++;
		}
	}

	// get null data
	function __getNullResult() constant internal returns(string _json) {
		uint len = 0;
		len = LibStack.push("{");
		len = LibStack.appendKeyValue("ret", uint(0));
		len = LibStack.append(",\"data\":{");
		len = LibStack.appendKeyValue("total", uint(0));
		len = LibStack.append(",\"items\":[]}}");
		_json = LibStack.popex(len);
	}

	function addFilter(string _filterJson) public returns(uint _filterId) {}

	function moduleIsExist(string _moduleId) public constant returns(uint) {
		uint errno_prefix = 95370;
		uint index = getModuleIndexById(_moduleId);
		if(index == uint(-1)) {
			log("module not exists", "RoleFilterManager");
			errno = errno_prefix + uint(ROLE_FILTER_ERRORCODE.MODULE_NOT_EXIST);
			Notify(errno, "module does not exists");
			return 1;
		}
		log("moduleIsExist(string):", index);
		return 0;
	}

	function getModuleIndexById(string _moduleId) constant private returns(uint _index) {
		for(uint i = 0; i < modules.length; ++i) {
			if(modules[i].deleted)
				continue;
			if(modules[i].moduleId.equals(_moduleId))
				return i;
		}
		return uint(-1);
	}

	function addActionToRole(string _moduleId, string _roleId, string _actionId) public returns(uint) {
		log("addActionToRole->into...", "RoleFilterManager");
		log("addActionToRole->msg.sender:", msg.sender);
		log("addActionToRole->owner:", owner);
		log("addActionToRole->address:", this);
		uint errno_prefix = 95370;
		uint retUint = 0;

		if(bytes(_moduleId).length == 0) {
			log("addActionToRole->moduleId is null..", "RoleFilterManager");
			errno = uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "addActionToRole->moduleId is null.");
			return errno;
		}

		// get index from the modules by moduleId
		uint mindex = getIndexById(_moduleId);

		// do insert : invoke RoleManager.insert();
		address roleAddress = rm.getContractAddress("SystemModuleManager","0.0.1.0","RoleManager", "0.0.1.0");
		if(roleAddress == 0) {
			log("addRoleToAction,roleAddress is null ");
			errno = errno_prefix + uint256(ROLE_FILTER_ERRORCODE.INVALID_PARAM);
			Notify(errno, "addRoleToAction roleAddress is null");
			return errno;
		}
		IRoleManager roleRmanager = IRoleManager(roleAddress);
		// set role state = 1 valid
		tmpRole.status = 1;
		uint ret = roleRmanager.addActionToRole(_roleId, _actionId);
		if(ret != 0) {
			log("addRoleToAction,role insert fail..");
			errno = errno_prefix + uint256(ROLE_FILTER_ERRORCODE.SYS_EXCEPTION);
			Notify(errno, "role insert fail..");
			return errno;
		}

		revision++;

		log("addRoleToAction  success...");
		errno = uint256(ROLE_FILTER_ERRORCODE.NO_ERROR);
		Notify(errno, "addRoleToAction  success.");
		return errno;
	}
}
