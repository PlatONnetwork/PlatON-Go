pragma solidity ^0.4.12;
/**
* @file      SystemModuleManager.sol
* @author    Jungle
* @time      2017-5-18 10:34:17
* @desc      the module of system
*/

import "./sysbase/OwnerNamed.sol";
import "./sysbase/BaseModule.sol";

import "./library/LibModule.sol";
import "./library/LibContract.sol";

contract SystemModuleManager is BaseModule {

    using LibModule for *;
    using LibContract for *;
    using LibString for *;
    using LibInt for *;

    LibModule.Module tmpModule;
    LibContract.Contract tmpContract;

    LibContract.OpenAction[] tmpOpenActionList;

    uint reversion;
    uint nowTime;

    string[] tmpArr;

    enum MODULE_ERROR {
        NO_ERROR
    }

    // define role
    enum ROLE_DEFINE {
        ROLE_SUPER,
        ROLE_ADMIN,
        ROLE_PLAIN
    }

    event Notify(uint _error, string _info);

    // module : predefined data
    function SystemModuleManager(){
        uint ret = 0;
        reversion = 0;
        register("SystemModuleManager","0.0.1.0");

        // insert module data
        nowTime = now * 1000;
        tmpModule.moduleName = "SystemModuleManager";
        tmpModule.moduleVersion = "0.0.1.0";
        tmpModule.moduleEnable = 1;
        tmpModule.moduleDescription = "内置合约权限配置器";
        tmpModule.moduleText = "DAPP-控台";
        tmpModule.moduleCreateTime = nowTime;
        tmpModule.moduleUpdateTime = nowTime;
        tmpModule.moduleCreator = msg.sender;
        tmpModule.icon = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACgAAAAoCAYAAACM/rhtAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAA4ZpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuNi1jMDY3IDc5LjE1Nzc0NywgMjAxNS8wMy8zMC0yMzo0MDo0MiAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wTU09Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9tbS8iIHhtbG5zOnN0UmVmPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvc1R5cGUvUmVzb3VyY2VSZWYjIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDoyNTBjNTFiOS1jMDc2LWM4NGEtOWEyOS1mMTNjMGY3NGExNjIiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6NDg5QkFEN0E3QTk0MTFFNzlBRjVBNzc2RDE2RThBQUQiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6NDg5QkFENzk3QTk0MTFFNzlBRjVBNzc2RDE2RThBQUQiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTUgKE1hY2ludG9zaCkiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDoyYWQ1MGExYy1lY2Y4LTRiYTYtYmZjNy02Y2VlMmVhMDQ2YWIiIHN0UmVmOmRvY3VtZW50SUQ9ImFkb2JlOmRvY2lkOnBob3Rvc2hvcDowZWU4MWJlYS1iZTI2LTExN2EtYmY4YS05ZWFhNmVmNDVmM2QiLz4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz7MPbc/AAAEF0lEQVR42sSZWUhVURSGz71p5FhmFlYmRTlkUtANw7c0IvIhJMiiwayHspF8DDLKnoIeRH2IsjkbycKKoMHHMk0qG9SHCLWssMk0QUv7F/wnDod7z3TvsQUf595z9t3r9+y911576WloaFAcmBcsAkuAD6SAaSCKz/vBe9AOmkA9aATDdh2F2Ww/HewA68AP8ABcBa2gE/SxXTRIAqkgG5wA48EFUAW6Qi0wHhwCa8FpkAdaDNp/Jc/BFd7LBEXgBbgISsEXK0NlZqvBK/CHQ1liIi6QtfC3c9iX9Flg9iOPwRyUt1sBcjmkjUpoTeZuDafJLvDbzhuMBLUgmR2FWpzCxeOjj1r6tCQwjHNEVuJK0Ku4Z7300U+fYVYEVsrQg41gSHHfhujLQ9+GAgs459aDQWX0bJA+cxkp/AqUUFIO1rg8rEbDvYYaEvwJPMyg+9Sph4iIiMngLhgGH8EGm12I70ugTB9mZoBnjHM9QQi8hssqzS3Z2jIGBgZabXQziVvkAtChvsFicCYYcZrYpp9CPpt99HC3KlY78HKCVocotim6N9jkoJ9qavKqWYls/C9DIHA7uAdGwCfZe20Or2qyDX4XbRIYc9ipkzk3WzIWiLgt33H9jMsy3Pfi83CQf+x9SefkDS4Ejx2Kk+G7hc/7tM9CIE6hJp+XOdsbm+Ik37vBHE+sDPeWB2ibCZpAsU2BoilVBCaCDzbEeTiJM3SrtQbPZunaTsSljqNUge9ZNgRKRp4oHceAnzZ+WMIcUW9x4DpERFLcGCYAyXwu38/hfpRFP6IpRo2DIxbfnpxBjhg0mQ+O8fNBWTC655KsHrUzzrKTSGBMMwvSECdnjGZGejO7bJIt52Eh3bFwzGiTN9gNppqIGyfDZ1GcYiGVP4k+E0zayCmxWwS2gXSTxlUOtiwjm6KZCoEsTX2DkkEsNnh7W3HZ7EJ6lY++iwyei6YmLw/VSwOIy+LByS0rh4+ZAZ6JpnoR+ARM0MU11U6BcBcFSog77ud+BsNWo5cZx3mwxU/DdMV9y/Fzbws1DXs1i6DQzyptHgWB7/wkrIXU9C/l72KqXaprvMlhFcGqvWU5RGv7WS7p0lcWJDC+BiuCOZcEabJnSwCfq9ZttIcmubGH+2fsfxAXS997tUUl/blYhvmhbOour169hdNnPes1hpWF3byeHSWRqjiFRSTT0scgT/eSjN50ebhj6SOaPgetVrd+sajToalChdp87Ft85NOn5fKbWtTZxtBTx5JEfAiEqSUW6fMAfQQsUlmpsMrCkfLtWJ74JeGc50BYBn/bzr4yuWpNE1Y7TqSIvpOH6m88Gj5iytapOTrIHpvEUko2N/44bl+Vio0iuifIf0PkMLimUHw0n/dRRDuDvoQuR/+G+CvAAGjg+nJ1aHTpAAAAAElFTkSuQmCC";
        tmpModule.publishTime = nowTime;
        ret = addModule(tmpModule.toJson());
        if (ret != 0) {
            log("addModule SystemModuleManager failed.");
            log("abort!!!");
            return;
        }

        ret = initContractData();
        if (ret != 0) {
            log("initContractData failed.");
            log("abort!!!");
            return;
        }

        ret = initMenuData();
        if (ret != 0) {
            log("initMenuData failed.");
            log("abort!!!");
            return;
        }

        ret = initActionData();
        if (ret != 0) {
            log("initActionData failed.");
            log("abort!!!");
            return;
        }

        ret = initRoleData();
        if (ret != 0) {
            log("initRoleData failed.");
            log("abort!!!");
            return;
        }
    }

    function initContractData() private returns(uint) {
        uint ret = 0;

        // set shared variables
        tmpContract.moduleName = "SystemModuleManager";
        tmpContract.moduleVersion = "0.0.1.0";
        tmpContract.cctVersion = "0.0.1.0";
        tmpContract.deleted = false;
        tmpContract.enable = 1;
        tmpContract.createTime = nowTime;
        tmpContract.updateTime = nowTime;
        tmpContract.creator = msg.sender;
        tmpContract.blockNum = block.number;

        // 1.UserManager insert
        tmpContract.cctName = "UserManager";//
        tmpContract.description = "用户管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract UserManager failed.");
            return ret;
        }

        // 2.DepartmentManager insert
        tmpContract.cctName = "DepartmentManager";//
        tmpContract.description = "部门管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract DepartmentManager failed.");
            return ret;
        }

        // 3.RoleManager insert
        tmpContract.cctName = "RoleManager";//
        tmpContract.description = "角色管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract RoleManager failed.");
            return ret;
        }

        // 4.ActionManager insert
        tmpContract.cctName = "ActionManager";//
        tmpContract.description = "权限管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract ActionManager failed.");
            return ret;
        }

        // 5.FileInfoManager insert
        tmpContract.cctName = "FileInfoManager";//
        tmpContract.description = "文件信息管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract FileInfoManager failed.");
            return ret;
        }

        // 6.FileServerManager insert
        tmpContract.cctName = "FileServerManager";//
        tmpContract.description = "文件服务管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract FileServerManager failed.");
            return ret;
        }

        // 7.NodeApplyManager insert
        tmpContract.cctName = "NodeApplyManager";//
        tmpContract.description = "节点申请管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract NodeApplyManager failed.");
            return ret;
        }

        // 8.RoleFilterManager insert
        tmpContract.cctName = "RoleFilterManager";//
        tmpContract.description = "角色过滤器管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract RoleFilterManager failed.");
            return ret;
        }

        // 9.SystemConfig insert
        tmpContract.cctName = "SystemConfig";//
        tmpContract.description = "系统配置管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract SystemConfig failed.");
            return ret;
        }

        // 10.NodeInfoManager insert
        tmpContract.cctName = "NodeInfoManager";//
        tmpContract.description = "节点信息管理";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract NodeInfoManager failed.");
            return ret;
        }

        // 11.RegisterApplyManager insert
        tmpContract.cctName = "RegisterApplyManager";
        tmpContract.description = "注册申请管理合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract RegisterApplyManager failed.");
            return ret;
        }

        log("init contract success...","System Manager");
        return 0;
    }

    function initActionData() private returns(uint) {
        log("init action data ", "SystemModuleManager");
        uint ret = 0;
        string memory jsonStr;

        // ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ UserManager ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1001\",\"enable\":1,\"name\":\"getAccountState\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getAccountState\",\"resKey\":\"UserManager\",\"opKey\":\"getAccountState(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1002\",\"enable\":1,\"name\":\"findByDepartmentIdTree\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"user/findByDepartmentIdTree\",\"description\":\"findByDepartmentIdTree\",\"resKey\":\"UserManager\",\"opKey\":\"findByDepartmentIdTree(string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1003\",\"enable\":1,\"name\":\"findByMobile\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByMobile\",\"resKey\":\"UserManager\",\"opKey\":\"findByMobile(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1004\",\"enable\":1,\"name\":\"userExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"userExists\",\"resKey\":\"UserManager\",\"opKey\":\"userExists(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1005\",\"enable\":1,\"name\":\"getErrno\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getErrno\",\"resKey\":\"UserManager\",\"opKey\":\"getErrno()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1006\",\"enable\":1,\"name\":\"checkUserRole\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkUserRole\",\"resKey\":\"UserManager\",\"opKey\":\"checkUserRole(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1007\",\"enable\":1,\"name\":\"updatePasswordStatus\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"updatePasswordStatus\",\"resKey\":\"UserManager\",\"opKey\":\"updatePasswordStatus(address,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1008\",\"enable\":1,\"name\":\"actionUsed\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"actionUsed\",\"resKey\":\"UserManager\",\"opKey\":\"actionUsed(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1009\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"UserManager\",\"opKey\":\"log()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1010\",\"enable\":1,\"name\":\"login\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"login\",\"resKey\":\"UserManager\",\"opKey\":\"login(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1011\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"UserManager\",\"opKey\":\"log(string,int256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1012\",\"enable\":1,\"name\":\"pageByAccountStatus\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"user/pageByAccountStatus\",\"description\":\"pageByAccountStatus\",\"resKey\":\"UserManager\",\"opKey\":\"pageByAccountStatus(uint256,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1013\",\"enable\":0,\"name\":\"update\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"user/update\",\"description\":\"update\",\"resKey\":\"UserManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1014\",\"enable\":1,\"name\":\"register\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"register\",\"resKey\":\"UserManager\",\"opKey\":\"register(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1015\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"UserManager\",\"opKey\":\"log(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1016\",\"enable\":1,\"name\":\"getUserState\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserState\",\"resKey\":\"UserManager\",\"opKey\":\"getUserState(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1017\",\"enable\":1,\"name\":\"kill\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"kill\",\"resKey\":\"UserManager\",\"opKey\":\"kill()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1018\",\"enable\":1,\"name\":\"checkWritePermission\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkWritePermission\",\"resKey\":\"UserManager\",\"opKey\":\"checkWritePermission(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1019\",\"enable\":1,\"name\":\"findByEmail\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByEmail\",\"resKey\":\"UserManager\",\"opKey\":\"findByEmail(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1020\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"UserManager\",\"opKey\":\"log(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1021\",\"enable\":1,\"name\":\"checkRoleAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkRoleAction\",\"resKey\":\"UserManager\",\"opKey\":\"checkRoleAction(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1022\",\"enable\":1,\"name\":\"findByLoginName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByLoginName\",\"resKey\":\"UserManager\",\"opKey\":\"findByLoginName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1023\",\"enable\":1,\"name\":\"clearLog\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"clearLog\",\"resKey\":\"UserManager\",\"opKey\":\"clearLog()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1024\",\"enable\":1,\"name\":\"getSender\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getSender\",\"resKey\":\"UserManager\",\"opKey\":\"getSender()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1025\",\"enable\":1,\"name\":\"eraseAdminByAddress\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"eraseAdminByAddress\",\"resKey\":\"UserManager\",\"opKey\":\"eraseAdminByAddress(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1026\",\"enable\":1,\"name\":\"getChildIdByIndex\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getChildIdByIndex\",\"resKey\":\"UserManager\",\"opKey\":\"getChildIdByIndex(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1027\",\"enable\":1,\"name\":\"checkDepartmentAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkDepartmentAction\",\"resKey\":\"UserManager\",\"opKey\":\"checkDepartmentAction(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1028\",\"enable\":1,\"name\":\"departmentExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"departmentExists\",\"resKey\":\"UserManager\",\"opKey\":\"departmentExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1029\",\"enable\":1,\"name\":\"listAll\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"UserManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1030\",\"enable\":1,\"name\":\"getUserRoleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserRoleId\",\"resKey\":\"UserManager\",\"opKey\":\"getUserRoleId(address,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1031\",\"enable\":1,\"name\":\"getUserCountByDepartmentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserCountByDepartmentId\",\"resKey\":\"UserManager\",\"opKey\":\"getUserCountByDepartmentId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1032\",\"enable\":1,\"name\":\"getOwner\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getOwner\",\"resKey\":\"UserManager\",\"opKey\":\"getOwner()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1033\",\"enable\":1,\"name\":\"checkActionExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkActionExists\",\"resKey\":\"UserManager\",\"opKey\":\"checkActionExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1034\",\"enable\":1,\"name\":\"deleteByAddress\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"user/deleteByAddress\",\"description\":\"deleteByAddress\",\"resKey\":\"UserManager\",\"opKey\":\"deleteByAddress(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1035\",\"enable\":1,\"name\":\"addUserRole\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"addUserRole\",\"resKey\":\"UserManager\",\"opKey\":\"addUserRole(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1036\",\"enable\":1,\"name\":\"getDepartmentRoleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getDepartmentRoleId\",\"resKey\":\"UserManager\",\"opKey\":\"getDepartmentRoleId(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1037\",\"enable\":1,\"name\":\"getLog\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getLog\",\"resKey\":\"UserManager\",\"opKey\":\"getLog()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1038\",\"enable\":1,\"name\":\"findByAddress\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"user/findByAddress\",\"description\":\"findByAddress\",\"resKey\":\"UserManager\",\"opKey\":\"findByAddress(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1039\",\"enable\":1,\"name\":\"findByAccount\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByAccount\",\"resKey\":\"UserManager\",\"opKey\":\"findByAccount(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1040\",\"enable\":1,\"name\":\"roleUsed\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleUsed\",\"resKey\":\"UserManager\",\"opKey\":\"roleUsed(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1041\",\"enable\":1,\"name\":\"actionExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"actionExists\",\"resKey\":\"UserManager\",\"opKey\":\"actionExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1042\",\"enable\":1,\"name\":\"getUserDepartmentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserDepartmentId\",\"resKey\":\"UserManager\",\"opKey\":\"getUserDepartmentId(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1043\",\"enable\":1,\"name\":\"checkUserAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkUserAction\",\"resKey\":\"UserManager\",\"opKey\":\"checkUserAction(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1044\",\"enable\":1,\"name\":\"insert\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"user/insert\",\"description\":\"insert\",\"resKey\":\"UserManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1045\",\"enable\":1,\"name\":\"getUserCount\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserCount\",\"resKey\":\"UserManager\",\"opKey\":\"getUserCount()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1046\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"UserManager\",\"opKey\":\"log(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1047\",\"enable\":1,\"name\":\"checkDepartmentRole\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkDepartmentRole\",\"resKey\":\"UserManager\",\"opKey\":\"checkDepartmentRole(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1048\",\"enable\":1,\"name\":\"roleExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleExists\",\"resKey\":\"UserManager\",\"opKey\":\"roleExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1049\",\"enable\":1,\"name\":\"findByDepartmentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByDepartmentId\",\"resKey\":\"UserManager\",\"opKey\":\"findByDepartmentId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1050\",\"enable\":1,\"name\":\"checkUserPrivilege\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkUserPrivilege\",\"resKey\":\"UserManager\",\"opKey\":\"checkUserPrivilege(address,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1051\",\"enable\":1,\"name\":\"getAdmin\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getAdmin(string)\",\"resKey\":\"UserManager\",\"opKey\":\"getAdmin(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1052\",\"enable\":1,\"name\":\"checkActionWithKey\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkActionWithKey\",\"resKey\":\"UserManager\",\"opKey\":\"checkActionWithKey(string,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1053\",\"enable\":1,\"name\":\"findByLoginName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByLoginName\",\"resKey\":\"UserManager\",\"opKey\":\"findByLoginName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1054\",\"enable\":1,\"name\":\"updateAccountStatus\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"updateAccountStatus\",\"resKey\":\"UserManager\",\"opKey\":\"updateAccountStatus(address,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1055\",\"enable\":1,\"name\":\"findByDepartmentIdTreeAndContion\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByDepartmentIdTreeAndContion\",\"resKey\":\"UserManager\",\"opKey\":\"findByDepartmentIdTreeAndContion(uint256,string,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1056\",\"enable\":1,\"name\":\"findByRoleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByRoleId\",\"resKey\":\"UserManager\",\"opKey\":\"findByRoleId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1057\",\"enable\":1,\"name\":\"updateUserStatus\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"updateUserStatus\",\"resKey\":\"UserManager\",\"opKey\":\"updateUserStatus(address,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1058\",\"enable\":1,\"name\":\"getUserCountByActionId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserCountByActionId\",\"resKey\":\"UserManager\",\"opKey\":\"getUserCountByActionId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1059\",\"enable\":1,\"name\":\"getUserCountMappingByRoleIds\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserCountMappingByRoleIds\",\"resKey\":\"UserManager\",\"opKey\":\"getUserCountMappingByRoleIds(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1060\",\"enable\":1,\"name\":\"findByUuid\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByUuid\",\"resKey\":\"UserManager\",\"opKey\":\"findByUuid(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1061\",\"enable\":1,\"name\":\"resetPasswd\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"resetPasswd\",\"resKey\":\"UserManager\",\"opKey\":\"resetPasswd(address,address,string,string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1062\",\"enable\":1,\"name\":\"getRevision\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRevision\",\"resKey\":\"UserManager\",\"opKey\":\"getRevision()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1063\",\"enable\":1,\"name\":\"getUserAddrByAddr\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserAddrByAddr\",\"resKey\":\"UserManager\",\"opKey\":\"getUserAddrByAddr(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action1064\",\"enable\":1,\"name\":\"getOwnerAddrByAddr\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getOwnerAddrByAddr\",\"resKey\":\"UserManager\",\"opKey\":\"getOwnerAddrByAddr(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);


        // +++++++++++++++++++++++++++++++ DepartmentManager ++++++++++++++++++++++++++++++
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2001\",\"enable\":1,\"name\":\"setAdmin\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"setAdmin\",\"resKey\":\"DepartmentManager\",\"opKey\":\"setAdmin(string,address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2002\",\"enable\":1,\"name\":\"userExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"userExists\",\"resKey\":\"DepartmentManager\",\"opKey\":\"userExists(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2003\",\"enable\":1,\"name\":\"departmentEmpty\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"departmentEmpty\",\"resKey\":\"DepartmentManager\",\"opKey\":\"departmentEmpty(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2004\",\"enable\":1,\"name\":\"getErrno\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getErrno\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getErrno()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2005\",\"enable\":1,\"name\":\"deleteById\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"department/deleteById\",\"description\":\"deleteById\",\"resKey\":\"DepartmentManager\",\"opKey\":\"deleteById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2006\",\"enable\":1,\"name\":\"actionUsed\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"actionUsed\",\"resKey\":\"DepartmentManager\",\"opKey\":\"actionUsed(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2007\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"DepartmentManager\",\"opKey\":\"log(string,address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2008\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"DepartmentManager\",\"opKey\":\"log(string,int256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2009\",\"enable\":1,\"name\":\"update\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"department/update\",\"description\":\"update\",\"resKey\":\"DepartmentManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2010\",\"enable\":1,\"name\":\"register\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"register\",\"resKey\":\"DepartmentManager\",\"opKey\":\"register(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2011\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"DepartmentManager\",\"opKey\":\"log(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2012\",\"enable\":1,\"name\":\"findByParentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByParentId\",\"resKey\":\"DepartmentManager\",\"opKey\":\"findByParentId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2013\",\"enable\":1,\"name\":\"kill\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"kill\",\"resKey\":\"DepartmentManager\",\"opKey\":\"kill()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2014\",\"enable\":1,\"name\":\"checkWritePermission\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkWritePermission\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkWritePermission(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2015\",\"enable\":1,\"name\":\"pageByName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"pageByName\",\"resKey\":\"DepartmentManager\",\"opKey\":\"pageByName(string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2016\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"DepartmentManager\",\"opKey\":\"log(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2017\",\"enable\":1,\"name\":\"findById\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"department/findById\",\"description\":\"findById\",\"resKey\":\"DepartmentManager\",\"opKey\":\"findById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2018\",\"enable\":1,\"name\":\"checkRoleAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkRoleAction\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkRoleAction(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2019\",\"enable\":1,\"name\":\"findByName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByName\",\"resKey\":\"DepartmentManager\",\"opKey\":\"findByName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2020\",\"enable\":1,\"name\":\"clearLog\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"clearLog\",\"resKey\":\"DepartmentManager\",\"opKey\":\"clearLog()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2021\",\"enable\":1,\"name\":\"getSender\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getSender\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getSender()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2022\",\"enable\":1,\"name\":\"eraseAdminByAddress\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"eraseAdminByAddress\",\"resKey\":\"DepartmentManager\",\"opKey\":\"eraseAdminByAddress(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2023\",\"enable\":1,\"name\":\"getChildIdByIndex\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getChildIdByIndex\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getChildIdByIndex(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2024\",\"enable\":1,\"name\":\"isInWhiteList\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"isInWhiteList\",\"resKey\":\"DepartmentManager\",\"opKey\":\"isInWhiteList(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2025\",\"enable\":1,\"name\":\"checkDepartmentAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkDepartmentAction\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkDepartmentAction(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2026\",\"enable\":1,\"name\":\"departmentExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"departmentExists\",\"resKey\":\"DepartmentManager\",\"opKey\":\"departmentExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2027\",\"enable\":1,\"name\":\"listAll\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"department/listAll\",\"description\":\"listAll\",\"resKey\":\"DepartmentManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2028\",\"enable\":1,\"name\":\"getUserRoleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserRoleId\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getUserRoleId(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2029\",\"enable\":1,\"name\":\"getUserCountByDepartmentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserCountByDepartmentId\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getUserCountByDepartmentId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2030\",\"enable\":1,\"name\":\"getOwner\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getOwner\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getOwner()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2031\",\"enable\":1,\"name\":\"checkActionExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkActionExists\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkActionExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2032\",\"enable\":1,\"name\":\"getDepartmentRoleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getDepartmentRoleId\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getDepartmentRoleId(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2033\",\"enable\":1,\"name\":\"getLog\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getLog\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getLog()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2034\",\"enable\":1,\"name\":\"roleUsed\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleUsed\",\"resKey\":\"DepartmentManager\",\"opKey\":\"roleUsed(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2035\",\"enable\":1,\"name\":\"actionExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"actionExists\",\"resKey\":\"DepartmentManager\",\"opKey\":\"actionExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2036\",\"enable\":1,\"name\":\"getUserDepartmentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserDepartmentId\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getUserDepartmentId(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2037\",\"enable\":1,\"name\":\"checkUserAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkUserAction\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkUserAction(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2038\",\"enable\":1,\"name\":\"insert\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"department/insert\",\"description\":\"insert\",\"resKey\":\"DepartmentManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2039\",\"enable\":1,\"name\":\"log\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"log\",\"resKey\":\"DepartmentManager\",\"opKey\":\"log(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2040\",\"enable\":1,\"name\":\"checkDepartmentRole\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkDepartmentRole\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkDepartmentRole(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2041\",\"enable\":1,\"name\":\"getMiningInfo\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getMiningInfo\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getMiningInfo()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2042\",\"enable\":1,\"name\":\"roleExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleExists\",\"resKey\":\"DepartmentManager\",\"opKey\":\"roleExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2043\",\"enable\":1,\"name\":\"checkUserPrivilege\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkUserPrivilege\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkUserPrivilege(address,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2044\",\"enable\":1,\"name\":\"getAdmin\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getAdmin(string)\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getAdmin(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2045\",\"enable\":1,\"name\":\"checkActionWithKey\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkActionWithKey\",\"resKey\":\"DepartmentManager\",\"opKey\":\"checkActionWithKey(string,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2046\",\"enable\":1,\"name\":\"setDepartmentStatus\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"setDepartmentStatus\",\"resKey\":\"DepartmentManager\",\"opKey\":\"setDepartmentStatus(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2047\",\"enable\":1,\"name\":\"departmentExistsByCN\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"departmentExistsByCN\",\"resKey\":\"DepartmentManager\",\"opKey\":\"departmentExistsByCN(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2048\",\"enable\":1,\"name\":\"pageByNameAndStatus\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"pageByNameAndStatus\",\"resKey\":\"DepartmentManager\",\"opKey\":\"pageByNameAndStatus(string,uint256,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action2049\",\"enable\":1,\"name\":\"getRevision\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRevision\",\"resKey\":\"DepartmentManager\",\"opKey\":\"getRevision()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);

        // ++++++++++++++++++++++++++++++++++++++ ActionManager +++++++++++++++++++++++++++++++++++++++++++
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3000\",\"enable\":1,\"name\":\"actionExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"actionExists\",\"resKey\":\"ActionManager\",\"opKey\":\"actionExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3001\",\"enable\":1,\"name\":\"findActionByKey\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findActionByKey\",\"resKey\":\"ActionManager\",\"opKey\":\"findActionByKey(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3002\",\"enable\":1,\"name\":\"findActionById\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findActionById\",\"resKey\":\"ActionManager\",\"opKey\":\"findActionById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3003\",\"enable\":1,\"name\":\"listAll\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"ActionManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3004\",\"enable\":1,\"name\":\"checkActionWithKey\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkActionWithKey\",\"resKey\":\"ActionManager\",\"opKey\":\"checkActionWithKey(string,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3005\",\"enable\":1,\"name\":\"insert\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"insert\",\"resKey\":\"ActionManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3006\",\"enable\":1,\"name\":\"deleteById\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"deleteById\",\"resKey\":\"ActionManager\",\"opKey\":\"deleteById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3007\",\"enable\":1,\"name\":\"getCount\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getCount\",\"resKey\":\"ActionManager\",\"opKey\":\"getCountgetCount()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3008\",\"enable\":1,\"name\":\"findById\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findById\",\"resKey\":\"ActionManager\",\"opKey\":\"findById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3009\",\"enable\":1,\"name\":\"listContractActions\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listContractActions\",\"resKey\":\"ActionManager\",\"opKey\":\"listContractActions(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3010\",\"enable\":1,\"name\":\"findByKey\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByKey\",\"resKey\":\"ActionManager\",\"opKey\":\"findByKey(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3011\",\"enable\":1,\"name\":\"getActionListByModuleName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getActionListByModuleName\",\"resKey\":\"ActionManager\",\"opKey\":\"getActionListByModuleName(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3012\",\"enable\":1,\"name\":\"getActionListByContractId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getActionListByContractId\",\"resKey\":\"ActionManager\",\"opKey\":\"getActionListByContractId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3013\",\"enable\":1,\"name\":\"getActionListByContractName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getActionListByContractName\",\"resKey\":\"ActionManager\",\"opKey\":\"getActionListByContractName(string,string,string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3014\",\"enable\":1,\"name\":\"getActionListByModuleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getActionListByModuleId\",\"resKey\":\"ActionManager\",\"opKey\":\"getActionListByModuleId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3015\",\"enable\":1,\"name\":\"listByForUK\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listByForUK\",\"resKey\":\"ActionManager\",\"opKey\":\"listByForUK(uint256,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3016\",\"enable\":1,\"name\":\"listBy\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listBy\",\"resKey\":\"ActionManager\",\"opKey\":\"listBy(uint256,string,string,string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3017\",\"enable\":1,\"name\":\"queryActionEnable\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"queryActionEnable\",\"resKey\":\"ActionManager\",\"opKey\":\"queryActionEnable(string,bool)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3018\",\"enable\":1,\"name\":\"update\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"update\",\"resKey\":\"ActionManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action3019\",\"enable\":1,\"name\":\"update\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findActionByType\",\"resKey\":\"ActionManager\",\"opKey\":\"findActionByType(uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);


        // ++++++++++++++++++++++++++++++++++++++++++ RoleManager +++++++++++++++++++++++++++++++++++++++
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4000\",\"enable\":1,\"name\":\"insert\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"role/insert\",\"description\":\"insert\",\"resKey\":\"RoleManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4001\",\"enable\":1,\"name\":\"update\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"role/update\",\"description\":\"update\",\"resKey\":\"RoleManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4002\",\"enable\":1,\"name\":\"listAll\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"RoleManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4003\",\"enable\":1,\"name\":\"findById\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"role/findById\",\"description\":\"findById\",\"resKey\":\"RoleManager\",\"opKey\":\"findById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4004\",\"enable\":1,\"name\":\"findByName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByName\",\"resKey\":\"RoleManager\",\"opKey\":\"findByName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4005\",\"enable\":1,\"name\":\"checkRoleAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkRoleAction\",\"resKey\":\"RoleManager\",\"opKey\":\"checkRoleAction(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4006\",\"enable\":1,\"name\":\"roleExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleExists\",\"resKey\":\"RoleManager\",\"opKey\":\"roleExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4007\",\"enable\":1,\"name\":\"pageByName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"role/pageByName\",\"description\":\"pageByName\",\"resKey\":\"RoleManager\",\"opKey\":\"pageByName(string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4008\",\"enable\":1,\"name\":\"deleteById\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"role/deleteById\",\"description\":\"deleteById\",\"resKey\":\"RoleManager\",\"opKey\":\"deleteById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4009\",\"enable\":1,\"name\":\"getName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getName\",\"resKey\":\"RoleManager\",\"opKey\":\"getName()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4010\",\"enable\":1,\"name\":\"getVersion\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getVersion\",\"resKey\":\"RoleManager\",\"opKey\":\"getVersion()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4011\",\"enable\":1,\"name\":\"userExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"userExists\",\"resKey\":\"RoleManager\",\"opKey\":\"userExists(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4012\",\"enable\":1,\"name\":\"actionUsed\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"actionUsed\",\"resKey\":\"RoleManager\",\"opKey\":\"actionUsed(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4013\",\"enable\":1,\"name\":\"register\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"register\",\"resKey\":\"RoleManager\",\"opKey\":\"register(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4014\",\"enable\":1,\"name\":\"checkWritePermission\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkWritePermission\",\"resKey\":\"RoleManager\",\"opKey\":\"checkWritePermission(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4015\",\"enable\":1,\"name\":\"checkRoleActionWithKey\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkRoleActionWithKey\",\"resKey\":\"RoleManager\",\"opKey\":\"checkRoleActionWithKey(string,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4016\",\"enable\":1,\"name\":\"eraseAdminByAddress\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"eraseAdminByAddress\",\"resKey\":\"RoleManager\",\"opKey\":\"eraseAdminByAddress(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4017\",\"enable\":1,\"name\":\"getChildIdByIndex\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getChildIdByIndex\",\"resKey\":\"RoleManager\",\"opKey\":\"getChildIdByIndex(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4018\",\"enable\":1,\"name\":\"checkDepartmentAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkDepartmentAction\",\"resKey\":\"RoleManager\",\"opKey\":\"checkDepartmentAction(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4019\",\"enable\":1,\"name\":\"departmentExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"departmentExists\",\"resKey\":\"RoleManager\",\"opKey\":\"departmentExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4020\",\"enable\":1,\"name\":\"getUserRoleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserRoleId\",\"resKey\":\"RoleManager\",\"opKey\":\"getUserRoleId(address,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4021\",\"enable\":1,\"name\":\"getUserCountByDepartmentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserCountByDepartmentId\",\"resKey\":\"RoleManager\",\"opKey\":\"getUserCountByDepartmentId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4022\",\"enable\":1,\"name\":\"getOwner\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getOwner\",\"resKey\":\"RoleManager\",\"opKey\":\"getOwner()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4023\",\"enable\":1,\"name\":\"checkActionExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkActionExists\",\"resKey\":\"RoleManager\",\"opKey\":\"checkActionExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4024\",\"enable\":1,\"name\":\"getDepartmentRoleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getDepartmentRoleId\",\"resKey\":\"RoleManager\",\"opKey\":\"getDepartmentRoleId(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4025\",\"enable\":1,\"name\":\"roleUsed\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleUsed\",\"resKey\":\"RoleManager\",\"opKey\":\"roleUsed(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4026\",\"enable\":1,\"name\":\"actionExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"actionExists\",\"resKey\":\"RoleManager\",\"opKey\":\"actionExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4027\",\"enable\":1,\"name\":\"getUserDepartmentId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getUserDepartmentId\",\"resKey\":\"RoleManager\",\"opKey\":\"getUserDepartmentId(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4028\",\"enable\":1,\"name\":\"checkUserAction\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkUserAction\",\"resKey\":\"RoleManager\",\"opKey\":\"checkUserAction(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4029\",\"enable\":1,\"name\":\"checkDepartmentRole\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkDepartmentRole\",\"resKey\":\"RoleManager\",\"opKey\":\"checkDepartmentRole(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4030\",\"enable\":1,\"name\":\"roleExists\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleExists\",\"resKey\":\"RoleManager\",\"opKey\":\"roleExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4031\",\"enable\":1,\"name\":\"checkUserPrivilege\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkUserPrivilege\",\"resKey\":\"RoleManager\",\"opKey\":\"checkUserPrivilege(address,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4032\",\"enable\":1,\"name\":\"getAdmin\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getAdmin\",\"resKey\":\"RoleManager\",\"opKey\":\"getAdmin(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4033\",\"enable\":1,\"name\":\"checkActionWithKey\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"checkActionWithKey\",\"resKey\":\"RoleManager\",\"opKey\":\"checkActionWithKey(string,address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4034\",\"enable\":1,\"name\":\"getRoleListByModuleName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRoleListByModuleName\",\"resKey\":\"RoleManager\",\"opKey\":\"getRoleListByModuleName(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4035\",\"enable\":1,\"name\":\"getRoleListByContractId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRoleListByContractId\",\"resKey\":\"RoleManager\",\"opKey\":\"getRoleListByContractId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4036\",\"enable\":1,\"name\":\"pageByNameAndModuleName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"pageByNameAndModuleName\",\"resKey\":\"RoleManager\",\"opKey\":\"pageByNameAndModuleName(string,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4037\",\"enable\":1,\"name\":\"addActionToRole\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"addActionToRole\",\"resKey\":\"RoleManager\",\"opKey\":\"addActionToRole(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4038\",\"enable\":1,\"name\":\"getRoleListByModuleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRoleListByModuleId\",\"resKey\":\"RoleManager\",\"opKey\":\"getRoleListByModuleId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4039\",\"enable\":1,\"name\":\"roleExistsEx\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"roleExistsEx\",\"resKey\":\"RoleManager\",\"opKey\":\"roleExistsEx(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4040\",\"enable\":1,\"name\":\"pageByNameAndModuleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"pageByNameAndModuleId\",\"resKey\":\"RoleManager\",\"opKey\":\"pageByNameAndModuleId(string,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4041\",\"enable\":1,\"name\":\"getRoleIdByActionIdAndIndex\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRoleIdByActionIdAndIndex\",\"resKey\":\"RoleManager\",\"opKey\":\"getRoleIdByActionIdAndIndex(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4042\",\"enable\":1,\"name\":\"getRoleModuleId\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRoleModuleId\",\"resKey\":\"RoleManager\",\"opKey\":\"getRoleModuleId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4043\",\"enable\":1,\"name\":\"getRoleModuleName\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRoleModuleName\",\"resKey\":\"RoleManager\",\"opKey\":\"getRoleModuleName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action4044\",\"enable\":1,\"name\":\"getRoleModuleVersion\",\"type\":2,\"level\":3,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getRoleModuleVersion\",\"resKey\":\"RoleManager\",\"opKey\":\"getRoleModuleVersion(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);


        // +++++++++++++++++++++++++++++++++++++++++ NodeApplyManager ++++++++++++++++++++++++++++++++++++++++++++
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5000\",\"enable\":1,\"name\":\"getErrno\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"getErrno\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"getErrno()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5001\",\"enable\":1,\"name\":\"register\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"register\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"register(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5002\",\"enable\":1,\"name\":\"kill\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"kill\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"kill()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5003\",\"enable\":1,\"name\":\"getOwner\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"getOwner\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"getOwner()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5004\",\"enable\":1,\"name\":\"getSender\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"getSender\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"getSender()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5005\",\"enable\":1,\"name\":\"insert\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"insert\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5006\",\"enable\":1,\"name\":\"update\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"update\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5007\",\"enable\":1,\"name\":\"nodeApplyExists\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"nodeApplyExists\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"nodeApplyExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5008\",\"enable\":1,\"name\":\"auditing\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"auditing\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"auditing(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5009\",\"enable\":1,\"name\":\"pageByNameAndStatus\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"pageByNameAndStatus\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"pageByNameAndStatus(uint256,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5010\",\"enable\":1,\"name\":\"listAll\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5011\",\"enable\":1,\"name\":\"findByState\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"findByState\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"findByState(uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5012\",\"enable\":1,\"name\":\"findByApplyId\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"findByApplyId\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"findByApplyId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action5013\",\"enable\":1,\"name\":\"deleteById\",\"level\":3,\"type\":2,\"parentId\":\"action100018\",\"url\":\"\",\"description\":\"deleteById\",\"resKey\":\"NodeApplyManager\",\"opKey\":\"deleteById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);

        // ++++++++++++++++++++++++++++++++++++ RoleFilterManager ++++++++++++++++++++++++++++++++++++++++++
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6000\",\"enable\":1,\"name\":\"getErrno\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getErrno\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"getErrno()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6001\",\"enable\":1,\"name\":\"register\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"register\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"register(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6002\",\"enable\":1,\"name\":\"kill\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"kill\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"kill()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6003\",\"enable\":1,\"name\":\"getOwner\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getOwner\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"getOwner()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6004\",\"enable\":1,\"name\":\"getSender\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getSender\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"getSender()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6005\",\"enable\":1,\"name\":\"authorizeProcessor\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"authorizeProcessor\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"authorizeProcessor(address,address,string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6006\",\"enable\":1,\"name\":\"addModule\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"addModule\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"addModule(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6007\",\"enable\":1,\"name\":\"updModule\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"updModule\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"updModule(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6008\",\"enable\":1,\"name\":\"delModule\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"delModule\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"delModule(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6009\",\"enable\":1,\"name\":\"setModuleEnable\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"setModuleEnable\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"setModuleEnable(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6010\",\"enable\":1,\"name\":\"setConntractEnable\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"setConntractEnable\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"setConntractEnable(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6011\",\"enable\":1,\"name\":\"addContract\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"addContract\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"addContract(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6012\",\"enable\":1,\"name\":\"addMenu\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"addMenu\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"addMenu(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6013\",\"enable\":1,\"name\":\"addAction\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"addAction\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"addAction(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6014\",\"enable\":1,\"name\":\"addRole\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"addRole\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"addRole(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6015\",\"enable\":1,\"name\":\"listAll\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6016\",\"enable\":1,\"name\":\"qryModules\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"qryModules\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"qryModules()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6017\",\"enable\":1,\"name\":\"qryModuleDetail\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"qryModuleDetail\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"qryModuleDetail(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6018\",\"enable\":1,\"name\":\"listContractByModuleName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"listContractByModuleName\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"listContractByModuleName(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6019\",\"enable\":1,\"name\":\"changeModuleOwner\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"changeModuleOwner\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"changeModuleOwner(string,string,address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6020\",\"enable\":1,\"name\":\"qryModuleDetail\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"qryModuleDetail\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"qryModuleDetail(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6021\",\"enable\":1,\"name\":\"findByName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findByName\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"findByName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6022\",\"enable\":1,\"name\":\"findByModuleText\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findByModuleText\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"findByModuleText(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6023\",\"enable\":1,\"name\":\"findContractByModName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findContractByModName\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"findContractByModName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6024\",\"enable\":1,\"name\":\"findContractByModText\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findContractByModText\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"findContractByModText(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6025\",\"enable\":1,\"name\":\"findContractByModName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findContractByModName\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"findContractByModName(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6026\",\"enable\":1,\"name\":\"listContractByModNameAndCttName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"listContractByModNameAndCttName\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"listContractByModNameAndCttName(string,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6027\",\"enable\":1,\"name\":\"listContractByModTextAndCttName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"listContractByModTextAndCttName\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"listContractByModTextAndCttName(string,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6028\",\"enable\":1,\"name\":\"getModuleCount\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getModuleCount\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"getModuleCount()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6029\",\"enable\":1,\"name\":\"listContractByModuleId\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"listContractByModuleId\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"listContractByModuleId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6030\",\"enable\":1,\"name\":\"addFilter\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"addFilter\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"addFilter(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6031\",\"enable\":1,\"name\":\"moduleIsExist\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"moduleIsExist\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"moduleIsExist(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action6032\",\"enable\":1,\"name\":\"addActionToRole\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"addActionToRole\",\"resKey\":\"RoleFilterManager\",\"opKey\":\"addActionToRole(string,string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);

        // +++++++++++++++++++++++++++++++++++ NodeInfoManager ++++++++++++++++++++++++++++++++++++++
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7000\",\"enable\":1,\"name\":\"getErrno\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getErrno\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"getErrno()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7001\",\"enable\":1,\"name\":\"register\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"register\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"register(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7002\",\"enable\":1,\"name\":\"kill\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"kill\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"kill()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7003\",\"enable\":1,\"name\":\"getOwner\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getOwner\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"getOwner()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7004\",\"enable\":1,\"name\":\"getSender\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getSender\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"getSender()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7005\",\"enable\":1,\"name\":\"insert\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"insert\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7006\",\"enable\":1,\"name\":\"update\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"update\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7007\",\"enable\":1,\"name\":\"getEnodeList\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getEnodeList\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"getEnodeList()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7008\",\"enable\":0,\"name\":\"ActivateEnode\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"ActivateEnode\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"ActivateEnode(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7009\",\"enable\":1,\"name\":\"isInWhiteList\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"isInWhiteList\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"isInWhiteList(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7010\",\"enable\":1,\"name\":\"getNodeAdmin\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getNodeAdmin\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"getNodeAdmin(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7011\",\"enable\":1,\"name\":\"setAdmin\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"setAdmin\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"setAdmin(string,address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7012\",\"enable\":1,\"name\":\"eraseAdminByAdd\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"eraseAdminByAdd\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"eraseAdminByAdd(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7013\",\"enable\":1,\"name\":\"nodeInfoExists\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"nodeInfoExists\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"nodeInfoExists(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7014\",\"enable\":1,\"name\":\"deleteById\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"deleteById\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"deleteById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7015\",\"enable\":1,\"name\":\"getRevision\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"getRevision\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"getRevision()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7016\",\"enable\":1,\"name\":\"IPUsed\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"IPUsed\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"IPUsed(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7017\",\"enable\":1,\"name\":\"listAll\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7018\",\"enable\":1,\"name\":\"findById\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findById\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"findById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7019\",\"enable\":1,\"name\":\"findByName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findByName\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"findByName(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7020\",\"enable\":1,\"name\":\"findByDepartmentId\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findByDepartmentId\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"findByDepartmentId(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7021\",\"enable\":1,\"name\":\"findByNodeAdmin\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findByNodeAdmin\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"findByNodeAdmin(address)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7022\",\"enable\":1,\"name\":\"checkWritePermission\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"checkWritePermission\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"checkWritePermission(address,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7023\",\"enable\":1,\"name\":\"findByPubkey\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"findByPubkey\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"findByPubkey(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7024\",\"enable\":1,\"name\":\"updateState\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"updateState\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"updateState(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7025\",\"enable\":1,\"name\":\"listByStateAndTypeAndName\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"listByStateAndTypeAndName\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"listByStateAndTypeAndName(uint256,uint256,string,uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action7026\",\"enable\":1,\"name\":\"updateDisabled\",\"level\":3,\"type\":2,\"parentId\":\"action100021\",\"url\":\"\",\"description\":\"updateDisabled\",\"resKey\":\"NodeInfoManager\",\"opKey\":\"updateDisabled(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        

        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action8000\",\"enable\":1,\"name\":\"writeConfig\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"writeConfig\",\"resKey\":\"SystemConfig\",\"opKey\":\"writeConfig(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action8001\",\"enable\":1,\"name\":\"readConfig\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"readConfig\",\"resKey\":\"SystemConfig\",\"opKey\":\"readConfig(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);


        // ============================================ FileInfoManager =====================================
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_getErrno\",\"enable\":1,\"name\":\"getErrno\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getErrno\",\"resKey\":\"FileInfoManager\",\"opKey\":\"getErrno()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_stringToUint\",\"enable\":1,\"name\":\"stringToUint\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"stringToUint\",\"resKey\":\"FileInfoManager\",\"opKey\":\"stringToUint(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_listByGroup\",\"enable\":1,\"name\":\"listByGroup\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listByGroup\",\"resKey\":\"FileInfoManager\",\"opKey\":\"listByGroup(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_deleteById\",\"enable\":1,\"name\":\"deleteById\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"deleteById\",\"resKey\":\"FileInfoManager\",\"opKey\":\"deleteById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_pageByGroup\",\"enable\":1,\"name\":\"pageByGroup\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"pageByGroup\",\"resKey\":\"FileInfoManager\",\"opKey\":\"pageByGroup(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_update\",\"enable\":1,\"name\":\"update\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"update\",\"resKey\":\"FileInfoManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_register\",\"enable\":1,\"name\":\"register\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"register\",\"resKey\":\"FileInfoManager\",\"opKey\":\"register(string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_generateFileID\",\"enable\":1,\"name\":\"generateFileID\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"generateFileID\",\"resKey\":\"FileInfoManager\",\"opKey\":\"generateFileID(string,string,string,string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_listAll\",\"enable\":1,\"name\":\"listAll\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"FileInfoManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_getCurrentPageCount\",\"enable\":1,\"name\":\"getCurrentPageCount\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getCurrentPageCount\",\"resKey\":\"FileInfoManager\",\"opKey\":\"getCurrentPageCount()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_find\",\"enable\":1,\"name\":\"find\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"find\",\"resKey\":\"FileInfoManager\",\"opKey\":\"find(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_getGroupPageCount\",\"enable\":1,\"name\":\"getGroupPageCount\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getGroupPageCount\",\"resKey\":\"FileInfoManager\",\"opKey\":\"getGroupPageCount(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_getCount\",\"enable\":1,\"name\":\"getCount\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getCount\",\"resKey\":\"FileInfoManager\",\"opKey\":\"getCount()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_insert\",\"enable\":1,\"name\":\"insert\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"insert\",\"resKey\":\"FileInfoManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_pageFiles\",\"enable\":1,\"name\":\"pageFiles\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"pageFiles\",\"resKey\":\"FileInfoManager\",\"opKey\":\"pageFiles(uint256,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_getCurrentPageSize\",\"enable\":1,\"name\":\"getCurrentPageSize\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getCurrentPageSize\",\"resKey\":\"FileInfoManager\",\"opKey\":\"getCurrentPageSize()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileInfoManager_getGroupFileCount\",\"enable\":1,\"name\":\"getGroupFileCount\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getGroupFileCount\",\"resKey\":\"FileInfoManager\",\"opKey\":\"getGroupFileCount(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_getErrno\",\"enable\":1,\"name\":\"getErrno\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getErrno\",\"resKey\":\"FileServerManager\",\"opKey\":\"getErrno()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_listByGroup\",\"enable\":1,\"name\":\"listByGroup\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listByGroup\",\"resKey\":\"FileServerManager\",\"opKey\":\"listByGroup(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_deleteById\",\"enable\":1,\"name\":\"deleteById\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"deleteById\",\"resKey\":\"FileServerManager\",\"opKey\":\"deleteById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_update\",\"enable\":1,\"name\":\"update\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"update\",\"resKey\":\"FileServerManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_listAll\",\"enable\":1,\"name\":\"listAll\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listAll\",\"resKey\":\"FileServerManager\",\"opKey\":\"listAll()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_find\",\"enable\":1,\"name\":\"find\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"find\",\"resKey\":\"FileServerManager\",\"opKey\":\"find(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_isServerEnable\",\"enable\":1,\"name\":\"isServerEnable\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"isServerEnable\",\"resKey\":\"FileServerManager\",\"opKey\":\"isServerEnable(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_getCount\",\"enable\":1,\"name\":\"getCount\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getCount\",\"resKey\":\"FileServerManager\",\"opKey\":\"getCount()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_insert\",\"enable\":1,\"name\":\"insert\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"insert\",\"resKey\":\"FileServerManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_enable\",\"enable\":1,\"name\":\"enable\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"enable\",\"resKey\":\"FileServerManager\",\"opKey\":\"enable(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"FileServerManager_findIdByHostPort\",\"enable\":1,\"name\":\"findIdByHostPort\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findIdByHostPort\",\"resKey\":\"FileServerManager\",\"opKey\":\"findIdByHostPort(string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        
        // ============================================ RegisterApplyManager =====================================
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9000\",\"enable\":0,\"name\":\"insert\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"insert\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"insert(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9001\",\"enable\":0,\"name\":\"update\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"update\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"update(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9002\",\"enable\":1,\"name\":\"audit\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"audit\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"audit(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9003\",\"enable\":0,\"name\":\"findById\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findById\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"findById(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9004\",\"enable\":0,\"name\":\"findByUuid\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"findByUuid\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"findByUuid(string)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9005\",\"enable\":0,\"name\":\"listByCondition\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"listByCondition\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"listByCondition(string,string,uint256,uint256,uint256,string,uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9006\",\"enable\":0,\"name\":\"getAutoAuditSwitch\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"getAutoAuditSwitch\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"getAutoAuditSwitch()\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        jsonStr = "{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action9007\",\"enable\":0,\"name\":\"updateAutoAuditSwitch\",\"level\":3,\"type\":2,\"parentId\":\"action100011\",\"url\":\"\",\"description\":\"updateAutoAuditSwitch\",\"resKey\":\"RegisterApplyManager\",\"opKey\":\"updateAutoAuditSwitch(uint256)\",\"version\":\"0.0.1.0\"}";
        addAction(jsonStr);
        log("init FileServerManager complete...");
        return 0;
    }

    // ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
    // +++++++++++++++++++++++ init menu ++++++++++++++++++++++++++++++++++++++++
    // ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
    // +++++++++++++++++++++++ init menu ++++++++++++++++++++++++++++++++++++++++
    function initMenuData() private returns(uint) {
        log("init menu data ","SystemModuleManager");
        string memory jsonStr;
        jsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action100011\",\"name\":\"链数据管理\",\"enable\":1,\"level\":1,\"parentId\":\"0\",\"url\":\"system\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"1\",\"type\":1,\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);
        jsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action100018\",\"name\":\"审核管理\",\"enable\":1,\"level\":2,\"parentId\":\"action100011\",\"url\":\"audit/list.do\",\"description\":\"\",\"resKey\":\"系统管理\",\"opKey\":\"5\",\"type\":1,\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);
        jsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action100019\",\"name\":\"权限管理\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"rolemgr\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\",\"type\":1,\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);
        jsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action100020\",\"name\":\"角色\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action100019\",\"url\":\"web/html/organUser/role.html\",\"description\":\"\",\"resKey\":\"角色\",\"opKey\":\"2\",\"type\":1,\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);
        jsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action100021\",\"name\":\"组织用户与权限\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action100019\",\"url\":\"web/html/organUser/organUser.html\",\"description\":\"\",\"resKey\":\"组织用户与权限\",\"opKey\":\"3\",\"type\":1,\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);
        log("init menu data complete..","SystemModuleManager");
        return 0;
    }

    // ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
    // +++++++++++++++++++++++ init role ++++++++++++++++++++++++++++++++++++++++
    function initRoleData() private returns(uint) {
        log("init role data ","SystemModuleManager");
        string memory roleJsonStr;
        roleJsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"role100000\",\"name\":\"节点管理员\",\"status\":1,\"description\":\"节点管理员\",\"actionIdList\":[\"action7000\",\"action7001\",\"action7002\",\"action7003\",\"action7004\",\"action7005\",\"action7006\",\"action7007\",\"action7008\",\"action7009\",\"action7010\",\"action7011\",\"action7012\",\"action7013\",\"action7014\",\"action7015\",\"action7016\",\"action7017\",\"action7018\",\"action7019\",\"action7020\",\"action7021\",\"action7022\",\"action7023\",\"action7024\",\"action7025\",\"action7026\",\"action1001\",\"action1002\",\"action1003\",\"action1004\",\"action1005\",\"action1006\",\"action1007\",\"action1008\",\"action1009\",\"action1010\",\"action1011\",\"action1012\",\"action1013\",\"action1014\",\"action1015\",\"action1016\",\"action1017\",\"action1018\",\"action1019\",\"action1020\",\"action1021\",\"action1022\",\"action1023\",\"action1024\",\"action1025\",\"action1026\",\"action1027\",\"action1028\",\"action1029\",\"action1030\",\"action1031\",\"action1032\",\"action1033\",\"action1034\",\"action1035\",\"action1036\",\"action1037\",\"action1038\",\"action1039\",\"action1040\",\"action1041\",\"action1042\",\"action1043\",\"action1044\",\"action1045\",\"action1046\",\"action1047\",\"action1048\",\"action1049\",\"action1050\",\"action1051\",\"action1052\",\"action1053\",\"action1054\",\"action1055\",\"action1056\",\"action1057\",\"action1058\",\"action1059\",\"action1060\",\"action1061\",\"action1062\",\"action1063\",\"action1064\",\"action2001\",\"action2002\",\"action2003\",\"action2004\",\"action2005\",\"action2006\",\"action2007\",\"action2008\",\"action2009\",\"action2010\",\"action2011\",\"action2012\",\"action2013\",\"action2014\",\"action2015\",\"action2016\",\"action2017\",\"action2018\",\"action2019\",\"action2020\",\"action2021\",\"action2022\",\"action2023\",\"action2024\",\"action2025\",\"action2026\",\"action2027\",\"action2028\",\"action2029\",\"action2030\",\"action2031\",\"action2032\",\"action2033\",\"action2034\",\"action2035\",\"action2036\",\"action2037\",\"action2038\",\"action2039\",\"action2040\",\"action2041\",\"action2042\",\"action2043\",\"action2044\",\"action2045\",\"action2046\",\"action2047\",\"action2048\",\"action2049\",\"action3000\",\"action3001\",\"action3002\",\"action3003\",\"action3004\",\"action3005\",\"action3006\",\"action3007\",\"action3008\",\"action3009\",\"action3010\",\"action3011\",\"action3012\",\"action3013\",\"action3014\",\"action3015\",\"action3016\",\"action3017\",\"action3018\",\"action3019\",\"action4000\",\"action4001\",\"action4002\",\"action4003\",\"action4004\",\"action4005\",\"action4006\",\"action4007\",\"action4008\",\"action4009\",\"action4010\",\"action4011\",\"action4012\",\"action4013\",\"action4014\",\"action4015\",\"action4016\",\"action4017\",\"action4018\",\"action4019\",\"action4020\",\"action4021\",\"action4022\",\"action4023\",\"action4024\",\"action4025\",\"action4026\",\"action4027\",\"action4028\",\"action4029\",\"action4030\",\"action4031\",\"action4032\",\"action4033\",\"action4034\",\"action4035\",\"action4036\",\"action4037\",\"action4038\",\"action4039\",\"action4040\",\"action4041\",\"action4042\",\"action4043\",\"action4044\",\"action5000\",\"action5001\",\"action5002\",\"action5003\",\"action5004\",\"action5005\",\"action5006\",\"action5007\",\"action5008\",\"action5009\",\"action5010\",\"action5011\",\"action5012\",\"action5013\"]}";
        addRole(roleJsonStr);
        roleJsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"role100001\",\"name\":\"链管理员\",\"status\":1,\"description\":\"链管理员\",\"actionIdList\":[\"action100011\",\"action100018\",\"action5000\",\"action5001\",\"action5002\",\"action5003\",\"action5004\",\"action5005\",\"action5006\",\"action5007\",\"action5008\",\"action5009\",\"action5010\",\"action5011\",\"action5012\",\"action5013\",\"FileInfoManager_getErrno\",\"FileInfoManager_stringToUint\",\"FileInfoManager_listByGroup\",\"FileInfoManager_deleteById\",\"FileInfoManager_pageByGroup\",\"FileInfoManager_update\",\"FileInfoManager_register\",\"FileInfoManager_generateFileID\",\"FileInfoManager_listAll\",\"FileInfoManager_getCurrentPageCount\",\"FileInfoManager_find\",\"FileInfoManager_getGroupPageCount\",\"FileInfoManager_getCount\",\"FileInfoManager_insert\",\"FileInfoManager_pageFiles\",\"FileInfoManager_getCurrentPageSize\",\"FileInfoManager_getGroupFileCount\",\"FileServerManager_getErrno\",\"FileServerManager_listByGroup\",\"FileServerManager_deleteById\",\"FileServerManager_update\",\"FileServerManager_listAll\",\"FileServerManager_isServerEnable\",\"FileServerManager_getCount\",\"FileServerManager_find\",\"FileServerManager_insert\",\"FileServerManager_enable\",\"FileServerManager_findIdByHostPort\",\"action1001\",\"action1002\",\"action1003\",\"action1004\",\"action1005\",\"action1006\",\"action1007\",\"action1008\",\"action1009\",\"action1010\",\"action1011\",\"action1012\",\"action1013\",\"action1014\",\"action1015\",\"action1016\",\"action1017\",\"action1018\",\"action1019\",\"action1020\",\"action1021\",\"action1022\",\"action1023\",\"action1024\",\"action1025\",\"action1026\",\"action1027\",\"action1028\",\"action1029\",\"action1030\",\"action1031\",\"action1032\",\"action1033\",\"action1034\",\"action1035\",\"action1036\",\"action1037\",\"action1038\",\"action1039\",\"action1040\",\"action1041\",\"action1042\",\"action1043\",\"action1044\",\"action1045\",\"action1046\",\"action1047\",\"action1048\",\"action1049\",\"action1050\",\"action1051\",\"action1052\",\"action1053\",\"action1054\",\"action1055\",\"action1056\",\"action1057\",\"action1058\",\"action1059\",\"action1060\",\"action1061\",\"action1062\",\"action1063\",\"action1064\",\"action2001\",\"action2002\",\"action2003\",\"action2004\",\"action2005\",\"action2006\",\"action2007\",\"action2008\",\"action2009\",\"action2010\",\"action2011\",\"action2012\",\"action2013\",\"action2014\",\"action2015\",\"action2016\",\"action2017\",\"action2018\",\"action2019\",\"action2020\",\"action2021\",\"action2022\",\"action2023\",\"action2024\",\"action2025\",\"action2026\",\"action2027\",\"action2028\",\"action2029\",\"action2030\",\"action2031\",\"action2032\",\"action2033\",\"action2034\",\"action2035\",\"action2036\",\"action2037\",\"action2038\",\"action2039\",\"action2040\",\"action2041\",\"action2042\",\"action2043\",\"action2044\",\"action2045\",\"action2046\",\"action2047\",\"action2048\",\"action2049\",\"action3000\",\"action3001\",\"action3002\",\"action3003\",\"action3004\",\"action3005\",\"action3006\",\"action3007\",\"action3008\",\"action3009\",\"action3010\",\"action3011\",\"action3012\",\"action3013\",\"action3014\",\"action3015\",\"action3016\",\"action3017\",\"action3018\",\"action3019\",\"action4000\",\"action4001\",\"action4002\",\"action4003\",\"action4004\",\"action4005\",\"action4006\",\"action4007\",\"action4008\",\"action4009\",\"action4010\",\"action4011\",\"action4012\",\"action4013\",\"action4014\",\"action4015\",\"action4016\",\"action4017\",\"action4018\",\"action4019\",\"action4020\",\"action4021\",\"action4022\",\"action4023\",\"action4024\",\"action4025\",\"action4026\",\"action4027\",\"action4028\",\"action4029\",\"action4030\",\"action4031\",\"action4032\",\"action4033\",\"action4034\",\"action4035\",\"action4036\",\"action4037\",\"action4038\",\"action4039\",\"action4040\",\"action4041\",\"action4042\",\"action4043\",\"action4044\"]}";
        addRole(roleJsonStr);
        roleJsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"role100002\",\"name\":\"系统管理员\",\"status\":1,\"description\":\"系统管理员\",\"actionIdList\":[\"action8000\",\"action8001\",\"action1001\",\"action1002\",\"action1003\",\"action1004\",\"action1005\",\"action1006\",\"action1007\",\"action1008\",\"action1009\",\"action1010\",\"action1011\",\"action1012\",\"action1013\",\"action1014\",\"action1015\",\"action1016\",\"action1017\",\"action1018\",\"action1019\",\"action1020\",\"action1021\",\"action1022\",\"action1023\",\"action1024\",\"action1025\",\"action1026\",\"action1027\",\"action1028\",\"action1029\",\"action1030\",\"action1031\",\"action1032\",\"action1033\",\"action1034\",\"action1035\",\"action1036\",\"action1037\",\"action1038\",\"action1039\",\"action1040\",\"action1041\",\"action1042\",\"action1043\",\"action1044\",\"action1045\",\"action1046\",\"action1047\",\"action1048\",\"action1049\",\"action1050\",\"action1051\",\"action1052\",\"action1053\",\"action1054\",\"action1055\",\"action1056\",\"action1057\",\"action1058\",\"action1059\",\"action1060\",\"action1061\",\"action1062\",\"action1063\",\"action1064\",\"action2001\",\"action2002\",\"action2003\",\"action2004\",\"action2005\",\"action2006\",\"action2007\",\"action2008\",\"action2009\",\"action2010\",\"action2011\",\"action2012\",\"action2013\",\"action2014\",\"action2015\",\"action2016\",\"action2017\",\"action2018\",\"action2019\",\"action2020\",\"action2021\",\"action2022\",\"action2023\",\"action2024\",\"action2025\",\"action2026\",\"action2027\",\"action2028\",\"action2029\",\"action2030\",\"action2031\",\"action2032\",\"action2033\",\"action2034\",\"action2035\",\"action2036\",\"action2037\",\"action2038\",\"action2039\",\"action2040\",\"action2041\",\"action2042\",\"action2043\",\"action2044\",\"action2045\",\"action2046\",\"action2047\",\"action2048\",\"action2049\",\"action3000\",\"action3001\",\"action3002\",\"action3003\",\"action3004\",\"action3005\",\"action3006\",\"action3007\",\"action3008\",\"action3009\",\"action3010\",\"action3011\",\"action3012\",\"action3013\",\"action3014\",\"action3015\",\"action3016\",\"action3017\",\"action3018\",\"action3019\",\"action4000\",\"action4001\",\"action4002\",\"action4003\",\"action4004\",\"action4005\",\"action4006\",\"action4007\",\"action4008\",\"action4009\",\"action4010\",\"action4011\",\"action4012\",\"action4013\",\"action4014\",\"action4015\",\"action4016\",\"action4017\",\"action4018\",\"action4019\",\"action4020\",\"action4021\",\"action4022\",\"action4023\",\"action4024\",\"action4025\",\"action4026\",\"action4027\",\"action4028\",\"action4029\",\"action4030\",\"action4031\",\"action4032\",\"action4033\",\"action4034\",\"action4035\",\"action4036\",\"action4037\",\"action4038\",\"action4039\",\"action4040\",\"action4041\",\"action4042\",\"action4043\",\"action4044\"]}";
        addRole(roleJsonStr);    
        roleJsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"role100003\",\"name\":\"权限管理员\",\"status\":1,\"description\":\"机构、用户权限管理员\",\"actionIdList\":[\"action2001\",\"action2002\",\"action2003\",\"action2004\",\"action2005\",\"action2006\",\"action2007\",\"action2008\",\"action2009\",\"action2010\",\"action2011\",\"action2012\",\"action2013\",\"action2014\",\"action2015\",\"action2016\",\"action2017\",\"action2018\",\"action2019\",\"action2020\",\"action2021\",\"action2022\",\"action2023\",\"action2024\",\"action2025\",\"action2026\",\"action2027\",\"action2028\",\"action2029\",\"action2030\",\"action2031\",\"action2032\",\"action2033\",\"action2034\",\"action2035\",\"action2036\",\"action2037\",\"action2038\",\"action2039\",\"action2040\",\"action2041\",\"action2042\",\"action2043\",\"action2044\",\"action2045\",\"action2046\",\"action2047\",\"action2048\",\"action2049\",\"action1001\",\"action1002\",\"action1003\",\"action1004\",\"action1005\",\"action1006\",\"action1007\",\"action1008\",\"action1009\",\"action1010\",\"action1011\",\"action1012\",\"action1013\",\"action1014\",\"action1015\",\"action1016\",\"action1017\",\"action1018\",\"action1019\",\"action1020\",\"action1021\",\"action1022\",\"action1023\",\"action1024\",\"action1025\",\"action1026\",\"action1027\",\"action1028\",\"action1029\",\"action1030\",\"action1031\",\"action1032\",\"action1033\",\"action1034\",\"action1035\",\"action1036\",\"action1037\",\"action1038\",\"action1039\",\"action1040\",\"action1041\",\"action1042\",\"action1043\",\"action1044\",\"action1045\",\"action1046\",\"action1047\",\"action1048\",\"action1049\",\"action1050\",\"action1051\",\"action1052\",\"action1053\",\"action1054\",\"action1055\",\"action1056\",\"action1057\",\"action1058\",\"action1059\",\"action1060\",\"action1061\",\"action1062\",\"action1063\",\"action1064\",\"action100019\",\"action100020\",\"action100021\",\"action6000\",\"action6001\",\"action6002\",\"action6003\",\"action6004\",\"action6005\",\"action6006\",\"action6007\",\"action6008\",\"action6009\",\"action6010\",\"action6011\",\"action6012\",\"action6013\",\"action6014\",\"action6015\",\"action6016\",\"action6017\",\"action6018\",\"action6019\",\"action6020\",\"action6021\",\"action6022\",\"action6023\",\"action6024\",\"action6025\",\"action6026\",\"action6027\",\"action6028\",\"action6029\",\"action6030\",\"action6031\",\"action6032\",\"action3000\",\"action3001\",\"action3002\",\"action3003\",\"action3004\",\"action3005\",\"action3006\",\"action3007\",\"action3008\",\"action3009\",\"action3010\",\"action3011\",\"action3012\",\"action3013\",\"action3014\",\"action3015\",\"action3016\",\"action3017\",\"action3018\",\"action3019\",\"action4000\",\"action4001\",\"action4002\",\"action4003\",\"action4004\",\"action4005\",\"action4006\",\"action4007\",\"action4008\",\"action4009\",\"action4010\",\"action4011\",\"action4012\",\"action4013\",\"action4014\",\"action4015\",\"action4016\",\"action4017\",\"action4018\",\"action4019\",\"action4020\",\"action4021\",\"action4022\",\"action4023\",\"action4024\",\"action4025\",\"action4026\",\"action4027\",\"action4028\",\"action4029\",\"action4030\",\"action4031\",\"action4032\",\"action4033\",\"action4034\",\"action4035\",\"action4036\",\"action4037\",\"action4038\",\"action4039\",\"action4040\",\"action4041\",\"action4042\",\"action4043\",\"action4044\",\"action9000\",\"action9001\",\"action9002\",\"action9003\",\"action9004\",\"action9005\",\"action9006\",\"action9007\"]}";
        addRole(roleJsonStr);
        roleJsonStr="{\"moduleName\":\"SystemModuleManager\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"role100004\",\"name\":\"普通用户\",\"status\":1,\"description\":\"普通用户\",\"actionIdList\":[\"action1001\",\"action1002\",\"action1003\",\"action1004\",\"action1005\",\"action1006\",\"action1008\",\"action1009\",\"action1010\",\"action1011\",\"action1012\",\"action1014\",\"action1015\",\"action1016\",\"action1017\",\"action1018\",\"action1019\",\"action1020\",\"action1021\",\"action1022\",\"action1024\",\"action1026\",\"action1027\",\"action1028\",\"action1029\",\"action1030\",\"action1031\",\"action1032\",\"action1033\",\"action1036\",\"action1037\",\"action1038\",\"action1039\",\"action1040\",\"action1041\",\"action1042\",\"action1043\",\"action1045\",\"action1046\",\"action1047\",\"action1048\",\"action1049\",\"action1050\",\"action1051\",\"action1052\",\"action1053\",\"action1055\",\"action1056\",\"action1058\",\"action1059\",\"action1060\",\"action1062\",\"action1063\",\"action1064\",\"action2002\",\"action2003\",\"action2004\",\"action2006\",\"action2007\",\"action2008\",\"action2011\",\"action2012\",\"action2014\",\"action2015\",\"action2016\",\"action2017\",\"action2018\",\"action2019\",\"action2021\",\"action2023\",\"action2024\",\"action2025\",\"action2026\",\"action2027\",\"action2028\",\"action2029\",\"action2030\",\"action2031\",\"action2032\",\"action2033\",\"action2034\",\"action2035\",\"action2036\",\"action2037\",\"action2039\",\"action2040\",\"action2041\",\"action2042\",\"action2043\",\"action2044\",\"action2045\",\"action2047\",\"action2048\",\"action2049\",\"action3000\",\"action3001\",\"action3002\",\"action3003\",\"action3004\",\"action3007\",\"action3008\",\"action3009\",\"action3010\",\"action3011\",\"action3012\",\"action3013\",\"action3014\",\"action3015\",\"action3016\",\"action3017\",\"action4002\",\"action4003\",\"action4004\",\"action4005\",\"action4006\",\"action4007\",\"action4009\",\"action4010\",\"action4011\",\"action4012\",\"action4014\",\"action4015\",\"action4017\",\"action4018\",\"action4019\",\"action4020\",\"action4021\",\"action4022\",\"action4023\",\"action4024\",\"action4025\",\"action4026\",\"action4027\",\"action4028\",\"action4029\",\"action4030\",\"action4031\",\"action4032\",\"action4033\",\"action4034\",\"action4035\",\"action4036\",\"action4038\",\"action4039\",\"action4040\",\"action4041\",\"action4042\",\"action4043\",\"action4044\",\"action5000\",\"action5003\",\"action5004\",\"action5007\",\"action5009\",\"action5010\",\"action5011\",\"action5012\",\"action6000\",\"action6003\",\"action6004\",\"action6015\",\"action6016\",\"action6017\",\"action6018\",\"action6020\",\"action6021\",\"action6022\",\"action6023\",\"action6024\",\"action6025\",\"action6026\",\"action6027\",\"action6028\",\"action6029\",\"action6031\",\"action7000\",\"action7003\",\"action7004\",\"action7007\",\"action7009\",\"action7010\",\"action7013\",\"action7015\",\"action7016\",\"action7017\",\"action7018\",\"action7019\",\"action7020\",\"action7021\",\"action7022\",\"action7023\",\"action7025\",\"action7026\",\"action8001\",\"FileInfoManager_getErrno\",\"FileInfoManager_stringToUint\",\"FileInfoManager_listByGroup\",\"FileInfoManager_pageByGroup\",\"FileInfoManager_listAll\",\"FileInfoManager_getCurrentPageCount\",\"FileInfoManager_find\",\"FileInfoManager_getGroupPageCount\",\"FileInfoManager_getCount\",\"FileInfoManager_pageFiles\",\"FileInfoManager_getCurrentPageSize\",\"FileInfoManager_getGroupFileCount\",\"FileServerManager_getErrno\",\"FileServerManager_listByGroup\",\"FileServerManager_listAll\",\"FileServerManager_find\",\"FileServerManager_isServerEnable\",\"FileServerManager_getCount\",\"FileServerManager_findIdByHostPort\",\"action9003\",\"action9004\",\"action9005\",\"action9006\",\"action3019\"]}";
        addRole(roleJsonStr);
        log("init role data completed...");
        return 0;
    }
}
