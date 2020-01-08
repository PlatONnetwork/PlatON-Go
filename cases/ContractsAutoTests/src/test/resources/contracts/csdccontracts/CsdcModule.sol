pragma solidity ^0.4.12;
/**
*@file    LibCsdc.sol
*@author  huanggaofeng
*@time    2017-07-20
*@desc    the definition of Csdc Privileges module
*/

import "./sysbase/OwnerNamed.sol";
import "./sysbase/BaseModule.sol";

import "./library/LibModule.sol";
import "./library/LibContract.sol";

contract CsdcModule is BaseModule {
    using LibModule for *;
    using LibContract for *;
    using LibString for *;
    using LibInt for *;
    using LibLog for *;

    LibModule.Module tmpModule;
    LibContract.Contract tmpContract;
    string regModuleId;
    uint nowTime;

    //模块构造函数
    function CsdcModule(){
        uint ret = 0;
        reversion = 0;
        register("CsdcModule","0.0.1.0");

        // insert module data
        nowTime = now * 1000;
        tmpModule.moduleName = "CsdcModule";
        tmpModule.moduleVersion = "0.0.1.0";
        tmpModule.moduleEnable = 0;
        tmpModule.moduleDescription = "中国结算权限控制模块";
        tmpModule.moduleText = "DAPP-中国结算";
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

		//合约注册
    function initContractData() private returns(uint) {


        uint ret = 0;

        // set shared variables
        tmpContract.moduleName = "CsdcModule";
        tmpContract.moduleVersion = "0.0.1.0";
        tmpContract.cctVersion = "0.0.1.0";
        tmpContract.deleted = false;
        tmpContract.enable = 0;
        tmpContract.createTime = nowTime;
        tmpContract.updateTime = nowTime;
        tmpContract.creator = msg.sender;
        tmpContract.blockNum = block.number;

        // 1.UserManager insert
        tmpContract.cctName = "Sequence";//
        tmpContract.description = "计数合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract Sequence failed.");
            return ret;
        }

        // 2.BrokerUserManager insert
        tmpContract.cctName = "BrokerUserManager";//
        tmpContract.description = "券商用户合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract BrokerUserManager failed.");
            return ret;
        }

        // 3.BrokerDao insert
        tmpContract.cctName = "BrokerDao";//
        tmpContract.description = "券商dao";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract BrokerDao failed.");
            return ret;
        }

        // 4.BrokerManager insert
        tmpContract.cctName = "BrokerManager";//
        tmpContract.description = "券商合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract BrokerManager failed.");
            return ret;
        }

        // 5.OrderDao insert
        tmpContract.cctName = "OrderDao";//
        tmpContract.description = "OrderDao";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract OrderDao failed.");
            return ret;
        }

        // 6.NewSecPledgeApplyManager insert
        tmpContract.cctName = "NewSecPledgeApplyManager";//
        tmpContract.description = "柜面/券商质押申请合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract NewSecPledgeApplyManager failed.");
            return ret;
        }

        // 7.NewBizManager insert
        tmpContract.cctName = "NewBizManager";//
        tmpContract.description = "柜面/券商业务合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract NewBizManager failed.");
            return ret;
        }

        // 8.NewSecPledgeManager insert
        tmpContract.cctName = "NewSecPledgeManager";//
        tmpContract.description = "柜面/券商质物状态合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract NewSecPledgeManager failed.");
            return ret;
        }

        // 9.NewDisSecPledgeApplyManager insert
        tmpContract.cctName = "NewDisSecPledgeApplyManager";//
        tmpContract.description = "柜面/券商解除质押申请合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract NewDisSecPledgeApplyManager failed.");
            return ret;
        }

        // 10.PerUserManager insert
        tmpContract.cctName = "PerUserManager";//
        tmpContract.description = "个人用户合约";//
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract PerUserManager failed.");
            return ret;
        }

        // 11.OrgUserManager insert
        tmpContract.cctName = "OrgUserManager";
        tmpContract.description = "机构用户合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract OrgUserManager failed.");
            return ret;
        }

        // 12.InvoiceManager insert
        tmpContract.cctName = "InvoiceManager";
        tmpContract.description = "发票合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract InvoiceManager failed.");
            return ret;
        }

        // 13.EvidenceManager insert
        tmpContract.cctName = "EvidenceManager";
        tmpContract.description = "证明文件合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract EvidenceManager failed.");
            return ret;
        }

        // 14.BizManager insert
        tmpContract.cctName = "BizManager";
        tmpContract.description = "在线业务合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract BizManager failed.");
            return ret;
        }

        // 15.PaymentManager insert
        tmpContract.cctName = "PaymentManager";
        tmpContract.description = "付款信息合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract PaymentManager failed.");
            return ret;
        }

        // 16.SecPledgeManager insert
        tmpContract.cctName = "SecPledgeManager";
        tmpContract.description = "在线质物状态合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract SecPledgeManager failed.");
            return ret;
        }


        // 17.SecPledgeApplyManager insert
        tmpContract.cctName = "SecPledgeApplyManager";
        tmpContract.description = "在线质押申请合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract SecPledgeApplyManager failed.");
            return ret;
        }


        // 18.DisSecPledgeApplyManager insert
        tmpContract.cctName = "DisSecPledgeApplyManager";
        tmpContract.description = "在线解除质押申请合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract DisSecPledgeApplyManager failed.");
            return ret;
        }


        // 19.SupplyDemandManager insert
        tmpContract.cctName = "SupplyDemandManager";
        tmpContract.description = "供需合约";
        tmpContract.enable = 0;
        ret = addContract(tmpContract.toJson());
        if (ret != 0) {
            log("addContract SupplyDemandManager failed.");
            return ret;
        }

        log("init contract success...","System Manager");
        return 0;
    }

    function initActionData() public returns(uint) {
        LibLog.log("init CsdcModule action data...");

        // --------- CsdcService -----------
        //action1001 action1002
        
        //add action control here
        //string memory jsonStr = "";
        //jsonStr = '{"id":"CsdcService_NY_select_Csdc_all","moduleId":"CsdcMudule_v0.0.1.0","contractId":"CsdcService_v0.0.1.0","name":"NY_select_Csdc_all","level":3,"type":2,"parentId":"action100011","url":"","description":"NY_select_Csdc_all","resKey":"CsdcService","opKey":"NY_select_Csdc_all()","version":"0.0.1.0"}';
        // addAction(jsonStr);

        LibLog.log("init CsdcModule action data complete.");
        
        return 0;
    }

    // +++++++++++++++++++++++ init menu ++++++++++++++++++++++++++++++++++++++++
    function initMenuData() public returns(uint) {
        LibLog.log("init CsdcModule menu data...");
        LibLog.log(jsonStr);

        // no menu for Csdc message
        string memory jsonStr="";
                
        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100020\",\"name\":\"业务审核\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"biz\",\"description\":\"\",\"resKey\":\"业务审核\",\"opKey\":\"0\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100020_1\",\"name\":\"证券质押审核\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100020\",\"url\":\"/csdc/static/html/brokerPledgeThirdAuditList.html\",\"description\":\"\",\"resKey\":\"业务审核\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100020_3\",\"name\":\"解除证券质押审核\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100020\",\"url\":\"/csdc/static/html/removePledgeAuditList.html\",\"description\":\"\",\"resKey\":\"业务审核\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100020_4\",\"name\":\"发票初审\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100020\",\"url\":\"/csdc/static/html/invoiceFirst.html\",\"description\":\"\",\"resKey\":\"业务审核\",\"opKey\":\"4\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100020_5\",\"name\":\"发票复核\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100020\",\"url\":\"/csdc/static/html/invoiceAgainList.html\",\"description\":\"\",\"resKey\":\"业务审核\",\"opKey\":\"5\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100020_10\",\"name\":\"黑名单申诉审核\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100020\",\"url\":\"/csdc/static/html/blackListAudit.html\",\"description\":\"\",\"resKey\":\"业务审核\",\"opKey\":\"10\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPost100000\",\"name\":\"邮寄管理\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"\",\"description\":\"\",\"resKey\":\"邮寄管理\",\"opKey\":\"0\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPost100000_1\",\"name\":\"所有业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPost100000\",\"url\":\"/csdc/static/html/documentList.html\",\"description\":\"\",\"resKey\":\"邮寄管理\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPost100000_2\",\"name\":\"证券质押登记证明\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPost100000\",\"url\":\"/csdc/static/html/pledgeDocument.html\",\"description\":\"\",\"resKey\":\"邮寄管理\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPost100000_3\",\"name\":\"证券质押登记证明<br>（部分解除质押）\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPost100000\",\"url\":\"/csdc/static/html/removePledgeDocument.html\",\"description\":\"\",\"resKey\":\"邮寄管理\",\"opKey\":\"2\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPost100000_4\",\"name\":\"普通发票\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPost100000\",\"url\":\"/csdc/static/html/plainInvoice.html\",\"description\":\"\",\"resKey\":\"邮寄管理\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPost100000_5\",\"name\":\"增值税专用发票\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPost100000\",\"url\":\"/csdc/static/html/fraudulentInvoice.html\",\"description\":\"\",\"resKey\":\"邮寄管理\",\"opKey\":\"4\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPost100000_6\",\"name\":\"解除质押登记通知\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPost100000\",\"url\":\"/csdc/static/html/cancellationRegistration.html\",\"description\":\"\",\"resKey\":\"邮寄管理\",\"opKey\":\"5\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);



        // // 待添加url
        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionAccounting100000\",\"name\":\"账务处理\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"0\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionAccounting100000_1\",\"name\":\"付款信息查询\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionAccounting100000\",\"url\":\"/csdc-bar/web/html/financing/PayInfoQuery/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionAccounting100000_2\",\"name\":\"在线业务对账\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionAccounting100000\",\"url\":\"/csdc-bar/web/html/financing/accountCheck/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionAccounting100000_3\",\"name\":\"退款管理\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionAccounting100000\",\"url\":\"\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionAccounting100000_3_1\",\"name\":\"退款初审\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"actionAccounting100000_3\",\"url\":\"/csdc-bar/web/html/financing/refundManager/refundFirstAudit/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionAccounting100000_3_2\",\"name\":\"退款复核\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"actionAccounting100000_3\",\"url\":\"/csdc-bar/web/html/financing/refundManager/refundFinalAudit/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);



        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200040\",\"name\":\"代理机构业务审核\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200050\",\"name\":\"所有业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200040\",\"url\":\"/csdc-bar/web/html/proxy/allBusi/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200051\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200050\",\"url\":\"/csdc-bar/web/html/proxy/allBusi/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200052\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200050\",\"url\":\"/csdc-bar/web/html/proxy/allBusi/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200060\",\"name\":\"待审核业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200040\",\"url\":\"/csdc-bar/web/html/proxy/pendBusi/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}"; 
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200061\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200060\",\"url\":\"/csdc-bar/web/html/proxy/pendBusi/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200062\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200060\",\"url\":\"/csdc-bar/web/html/proxy/pendBusi/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200070\",\"name\":\"待复核业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200040\",\"url\":\"/csdc-bar/web/html/proxy/reviewBusi/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200071\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200070\",\"url\":\"/csdc-bar/web/html/proxy/reviewBusi/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200072\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200070\",\"url\":\"/csdc-bar/web/html/proxy/reviewBusi/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200080\",\"name\":\"待领导审批业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200040\",\"url\":\"/csdc-bar/web/html/proxy/leaderBusi/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}"; 
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200081\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200080\",\"url\":\"/csdc-bar/web/html/proxy/leaderBusi/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200082\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200080\",\"url\":\"/csdc-bar/web/html/proxy/leaderBusi/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200090\",\"name\":\"待查看业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200040\",\"url\":\"/csdc-bar/web/html/proxy/viewBusi/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200091\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200090\",\"url\":\"/csdc-bar/web/html/proxy/viewBusi/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200092\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200090\",\"url\":\"/csdc-bar/web/html/proxy/viewBusi/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200100\",\"name\":\"已办结业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200040\",\"url\":\"/csdc-bar/web/html/proxy/doneBusi/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}"; 
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200101\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200100\",\"url\":\"/csdc-bar/web/html/proxy/doneBusi/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200102\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200100\",\"url\":\"/csdc-bar/web/html/proxy/doneBusi/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);




        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200000\",\"name\":\"柜台业务办理\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200001\",\"name\":\"所有业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200000\",\"url\":\"/csdc-bar/web/html/counters/allState/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200002\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200001\",\"url\":\"/csdc-bar/web/html/counters/allState/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200003\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200001\",\"url\":\"/csdc-bar/web/html/counters/allState/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200010\",\"name\":\"未提交复核业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200000\",\"url\":\"/csdc-bar/web/html/counters/noReviewed/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200011\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200010\",\"url\":\"/csdc-bar/web/html/counters/noReviewed/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200012\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200010\",\"url\":\"/csdc-bar/web/html/counters/noReviewed/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200020\",\"name\":\"待复核业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200000\",\"url\":\"/csdc-bar/web/html/counters/alReviewed/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200021\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200020\",\"url\":\"/csdc-bar/web/html/counters/alReviewed/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200022\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200020\",\"url\":\"/csdc-bar/web/html/counters/alReviewed/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200030\",\"name\":\"已办结业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200000\",\"url\":\"/csdc-bar/web/html/counters/alDone/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200031\",\"name\":\"证券质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200030\",\"url\":\"/csdc-bar/web/html/counters/alDone/pledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200032\",\"name\":\"证券解除质押\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"3\",\"parentId\":\"action200030\",\"url\":\"/csdc-bar/web/html/counters/alDone/rePledge.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"action200110\",\"name\":\"退款申请\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"action200000\",\"url\":\"/csdc-bar/web/html/counters/refundApply/index.html\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPub100000\",\"name\":\"公共功能\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"/public\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"0\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPub100000_1\",\"name\":\"所有业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPub100000\",\"url\":\"/public/all-business\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPub100000_2\",\"name\":\"我的待办业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPub100000\",\"url\":\"/public/my-todo-business\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPub100000_3\",\"name\":\"我的经手业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPub100000\",\"url\":\"/public/my-handled-business\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionPub100000_4\",\"name\":\"消息管理\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionPub100000\",\"url\":\"/public/message-management\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100001\",\"name\":\"投资人业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"/investor-business\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"0\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100001_1\",\"name\":\"代理机构信息维护\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100001\",\"url\":\"/investor-business/agency-information-maintenance\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100001_2\",\"name\":\"质押业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100001\",\"url\":\"/investor-business/pledge-business\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100001_3\",\"name\":\"解除质押业务\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100001\",\"url\":\"/investor-business/lifting-of-pledged-business\",\"description\":\"\",\"resKey\":\"broker-web\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100002\",\"name\":\"邮寄管理\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"1\",\"parentId\":\"0\",\"url\":\"/post-manager\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"0\", \"type\":\"1\",\"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100002_1\",\"name\":\"邮寄管理\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100002\",\"url\":\"/post-manage/post-manage\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100002_2\",\"name\":\"证券质押登记证明\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100002\",\"url\":\"/post-manage/pledge-evidence\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"1\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100002_3\",\"name\":\"解除质押登记通知\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100002\",\"url\":\"/post-manage/displedge-notice\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"3\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"actionBiz100002_4\",\"name\":\"证券质押登记证明<br>（部分解除质押）\",\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"enable\":1,\"level\":\"2\",\"parentId\":\"actionBiz100002\",\"url\":\"/post-manage/partial-pledge-evidence\",\"description\":\"\",\"resKey\":\"\",\"opKey\":\"2\",\"type\":\"1\", \"version\":\"0.0.1.0\"}";
        addMenu(jsonStr);


        LibLog.log("init CsdcModule menu data complete.");
        return 0;
    } 

    // +++++++++++++++++++++++ init role ++++++++++++++++++++++++++++++++++++++++
    function initRoleData() public returns(uint) {
        LibLog.log("init CsdcModule role data ");

        string memory jsonStr = "";

        /*中国结算用户权限*/    
        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_101\",\"name\":\"在线质押审核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"在线质押审核人员\",\"actionIdList\":[\"actionBiz100020\",\"actionBiz100020_1\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_102\",\"name\":\"在线解押审核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"在线解押审核人员\",\"actionIdList\":[\"actionBiz100020\",\"actionBiz100020_3\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_103\",\"name\":\"供需黑名单审核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"供需黑名单审核人员\",\"actionIdList\":[\"actionBiz100020\",\"actionBiz100020_10\"]}"; 
        addRole(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_104\",\"name\":\"发票初审\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"发票初审人员\",\"actionIdList\":[\"actionBiz100020\",\"actionBiz100020_4\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_105\",\"name\":\"发票复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"发票复核人员\",\"actionIdList\":[\"actionBiz100020\",\"actionBiz100020_5\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_106\",\"name\":\"邮寄管理\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"邮寄管理人员\",\"actionIdList\":[\"actionPost100000\",\"actionPost100000_1\",\"actionPost100000_2\",\"actionPost100000_3\",\"actionPost100000_4\",\"actionPost100000_5\",\"actionPost100000_6\"]}";
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_107\",\"name\":\"付款信息查询\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"付款信息查询人员\",\"actionIdList\":[\"actionAccounting100000\",\"actionAccounting100000_1\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_108\",\"name\":\"在线业务对账\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"在线业务对账人员\",\"actionIdList\":[\"actionAccounting100000\",\"actionAccounting100000_2\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_109\",\"name\":\"退款初审\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"退款初审人员\",\"actionIdList\":[\"actionAccounting100000\",\"actionAccounting100000_3\",\"actionAccounting100000_3_1\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_110\",\"name\":\"退款复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"退款复核人员\",\"actionIdList\":[\"actionAccounting100000\",\"actionAccounting100000_3\",\"actionAccounting100000_3_2\"]}"; 
        addRole(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_111\",\"name\":\"代理质押经办\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"代理质押经办人员\",\"actionIdList\":[\"action200040\",\"action200050\",\"action200051\",\"action200060\",\"action200061\",\"action200090\",\"action200091\",\"action200100\",\"action200101\",\"actionPost100000\",\"actionPost100000_1\",\"actionPost100000_2\",\"actionPost100000_3\",\"actionPost100000_4\",\"actionPost100000_5\",\"actionPost100000_6\"]}"; //actionPost100000_5
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_112\",\"name\":\"代理质押复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"代理质押复核人员\",\"actionIdList\":[\"action200040\",\"action200050\",\"action200051\",\"action200070\",\"action200071\",\"action200100\",\"action200101\",\"actionPost100000\",\"actionPost100000_1\",\"actionPost100000_2\",\"actionPost100000_3\",\"actionPost100000_4\",\"actionPost100000_5\",\"actionPost100000_6\"]}"; //actionPost100000_5
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_113\",\"name\":\"代理解押经办\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"代理解押经办人员\",\"actionIdList\":[\"action200040\",\"action200050\",\"action200052\",\"action200060\",\"action200062\",\"action200090\",\"action200092\",\"action200100\",\"action200102\",\"actionPost100000\",\"actionPost100000_1\",\"actionPost100000_2\",\"actionPost100000_3\",\"actionPost100000_4\",\"actionPost100000_5\",\"actionPost100000_6\"]}"; //actionPost100000_5
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_114\",\"name\":\"代理解押复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"代理解押复核人员\",\"actionIdList\":[\"action200040\",\"action200050\",\"action200052\",\"action200070\",\"action200072\",\"action200100\",\"action200102\",\"actionPost100000\",\"actionPost100000_1\",\"actionPost100000_2\",\"actionPost100000_3\",\"actionPost100000_4\",\"actionPost100000_5\",\"actionPost100000_6\"]}"; //actionPost100000_5
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_115\",\"name\":\"投业部领导\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"投业部领导\",\"actionIdList\":[\"action200040\",\"action200080\", \"action200081\", \"action200082\"]}"; 
        addRole(jsonStr);


        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_116\",\"name\":\"柜台质押经办\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台质押经办人员\",\"actionIdList\":[\"action200000\",\"action200001\",\"action200002\",\"action200010\",\"action200011\",\"action200020\",\"action200021\",\"action200030\",\"action200031\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_117\",\"name\":\"柜台质押复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台质押复核人员\",\"actionIdList\":[\"action200000\",\"action200001\",\"action200002\",\"action200020\",\"action200021\",\"action200030\",\"action200031\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_118\",\"name\":\"柜台解押经办\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台解押经办人员\",\"actionIdList\":[\"action200000\",\"action200001\",\"action200003\",\"action200010\",\"action200012\",\"action200020\",\"action200022\",\"action200030\",\"action200032\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_119\",\"name\":\"柜台解押复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台解押复核人员\",\"actionIdList\":[\"action200000\",\"action200001\",\"action200003\",\"action200020\",\"action200022\",\"action200030\",\"action200032\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_120\",\"name\":\"柜台退款申请\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台退款申请人员\",\"actionIdList\":[\"action200000\", \"action200110\"]}"; 
        addRole(jsonStr);


        /*券商用户权限*/
        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_201\",\"name\":\"券商质押经办\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台质押经办人员\",\"actionIdList\":[\"actionPub100000\",\"actionPub100000_1\",\"actionPub100000_2\",\"actionPub100000_3\",\"actionPub100000_4\",\"actionBiz100001\",\"actionBiz100001_2\",\"actionBiz100002\",\"actionBiz100002_1\",\"actionBiz100002_2\",\"actionBiz100002_3\",\"actionBiz100002_4\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_202\",\"name\":\"券商解押经办\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台质押复核人员\",\"actionIdList\":[\"actionPub100000\",\"actionPub100000_1\",\"actionPub100000_2\",\"actionPub100000_3\",\"actionPub100000_4\",\"actionBiz100001\",\"actionBiz100001_3\",\"actionBiz100002\",\"actionBiz100002_1\",\"actionBiz100002_2\",\"actionBiz100002_3\",\"actionBiz100002_4\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_203\",\"name\":\"券商质押复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台解押经办人员\",\"actionIdList\":[\"actionPub100000\",\"actionPub100000_1\",\"actionPub100000_2\",\"actionPub100000_3\",\"actionPub100000_4\",\"actionBiz100001\",\"actionBiz100001_1\",\"actionBiz100001_2\",\"actionBiz100002\",\"actionBiz100002_1\",\"actionBiz100002_2\",\"actionBiz100002_3\",\"actionBiz100002_4\"]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_204\",\"name\":\"券商解押复核\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"柜台解押复核人员\",\"actionIdList\":[\"actionPub100000\",\"actionPub100000_1\",\"actionPub100000_2\",\"actionPub100000_3\",\"actionPub100000_4\",\"actionBiz100001\",\"actionBiz100001_1\",\"actionBiz100001_3\",\"actionBiz100002\",\"actionBiz100002_1\",\"actionBiz100002_2\",\"actionBiz100002_3\",\"actionBiz100002_4\"]}"; 
        addRole(jsonStr);


        /*投资者用户权限*/
        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_301\",\"name\":\"在线质押\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"投资者用户（在线质押）\",\"actionIdList\":[]}"; 
        addRole(jsonStr);

        jsonStr="{\"moduleName\":\"CsdcModule\",\"moduleVersion\":\"0.0.1.0\",\"id\":\"csdc_role_302\",\"name\":\"在线解押\",\"status\":1,\"moduleId\":\"CsdcModule_v0.0.1.0\",\"contractId\":\"\",\"description\":\"投资者用户（在线解押）\",\"actionIdList\":[]}"; 
        addRole(jsonStr);


        LibLog.log("init CsdcModule role data completed.");
        return 0;
    }
}