pragma solidity ^0.4.12;

import "./Sequence.sol";
import "./OrderDao.sol";
import "./PaymentManager.sol";

contract NewSecPledgeApplyManager  is OwnerNamed  {
    using LibSecPledgeApply for *;
    using LibString for *;
    using LibInt for *;
    using LibBiz for *;
    using LibSecPledge for *;
    using LibJson for *;
    using LibPayment for *;

    OrderDao orderDao;
    Sequence sq;
    //inner setting member
    LibSecPledgeApply.SecPledgeApply internal tmp_SecPledgeApply;
    LibBiz.Biz internal tmp_Biz;
    LibSecPledge.SecPledge internal tmp_SecPledge;
    LibPayment.Payment tmp_payment;

    PaymentManager paymentManager;
    

    event Notify(uint _errorno, string _info);
 /** @brief errno for test case */
    enum SecPledgeApplyError {
        NO_ERROR,
        BAD_PARAMETER,
        OPERATE_NOT_ALLOWED,
        ID_EMPTY,   
        PAYERTYPE_EMPTY,
        INSERT_FAILED,
        USER_STATUS_ERROR,
        APPLIEDSECURITIES_ERROR,
        PLEDGOR_PLEDGEE_SAME,
        BUSINESSNO_ERROR //业务编号获取失败
    }

    enum OperateCode{
        NONE, 
        PASS, //通过
        REJECT, //拒绝
        WAIT //等待处理
    }

    enum OperError {
        NO_ERROR,
        BAD_PARAMETER,
        DAO_ERROR
    }

    uint errno_prefix = 10000;

    StateTrans[] state_trans_array; //保存状态变迁的表格
    function NewSecPledgeApplyManager() {
        register("CsdcModule", "0.0.1.0", "NewSecPledgeApplyManager", "0.0.1.0");
        
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
        orderDao = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
        paymentManager = PaymentManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PaymentManager", "0.0.1.0"));
    }

    modifier getOrderDao(){ 
      orderDao = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
      _;
    }

    modifier getPaymentManager(){ 
      paymentManager = PaymentManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PaymentManager", "0.0.1.0"));
      _;
    }

    function insertSecPledgeApply(string _json) getOrderDao {
        LibLog.log("insert: ", _json);
        if(!tmp_SecPledgeApply.fromJson(_json)) {
            LibLog.log("json invalid");
            Notify(errno_prefix + uint(OperError.BAD_PARAMETER), "json invalid");
            return;
        }
        tmp_SecPledgeApply.id = sq.getSeqNo("Biz.id");
        tmp_SecPledgeApply.bizId = tmp_SecPledgeApply.id;
        tmp_SecPledgeApply.secPledgeId = tmp_SecPledgeApply.id;
        tmp_SecPledgeApply.businessNo = sq.genBusinessNo("S", "01").recoveryToString();
        if (orderDao.insert_SecPledgeApply(tmp_SecPledgeApply.toJson()) != 0) {
          Notify(errno_prefix + uint(OperError.DAO_ERROR), "call dao error");
          return;
        }
        string memory success = "";
        uint _total = 1;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString() );
        success = success.concat(",\"id\":",  tmp_SecPledgeApply.id.toString());
       // success = success.concat(",\"businessNo\":",  businessNo);
        success = success.concat(",\"businessNo\":\"",  tmp_SecPledgeApply.businessNo,"\"");
        success = success.concat(",\"items\":[");
        success = success.concat("{\"id\":", tmp_SecPledgeApply.id.toString(),  "}]}}");
        Notify(0, success);
    }

    function updateSecPledgeApply(string _json) getOrderDao {
        LibLog.log("update_SecPledgeApply: ", _json);
        if (orderDao.update_SecPledgeApply(_json) != 0) {
          Notify(errno_prefix + uint(OperError.DAO_ERROR), "call dao error");
          return;
        }
        Notify(0, "success");
    }

    function addStateTrans(LibBiz.BizStatus current, OperateCode operate, LibBiz.BizStatus next) returns(bool _ret) {
        StateTrans memory stateTrans = 
          StateTrans(LibBiz.BizStatus.PLEDGE_PLEDGOR_UNPAID,
          OperateCode.PASS,
            LibBiz.BizStatus.PLEDGE_PLEDGOR_UNPAID

        );
        state_trans_array.push(stateTrans);

    }
    

    //17.1.质押申请保存(券商)
    //状态为BROKER_PLEDGE_STASHED
    function cacheSecPledgeApplyByBroker(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("cacheSecPledgeApplyByBroker: ", _json);
        _saveSecPledgeApplyByBroker(_json, LibBiz.BizStatus.BROKER_PLEDGE_STASHED);
    }

    //17.3.质押申请创建(券商)
    //状态为BROKER_PLEDGE_CREATED，产生businessNo
    function createSecPledgeApplyByBroker(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("createSecPledgeApplyByBroker: ", _json);
        _saveSecPledgeApplyByBroker(_json, LibBiz.BizStatus.BROKER_PLEDGE_CREATED);
    }

    //17.5.质押申请提交(券商)
    function saveSecPledgeApplyByBroker(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("saveSecPledgeApplyByBroker: ", _json);
        _saveSecPledgeApplyByBroker(_json, LibBiz.BizStatus.BROKER_PLEDGE_UNAUDITED);
    }

    //内部用于存储券商申请的方法
    function _saveSecPledgeApplyByBroker(string _json, LibBiz.BizStatus _status) internal returns(bool _ret) {
        LibLog.log("_saveSecPledgeApplyByBroker: ", _json);

        //封装Biz
        LibLog.log("封装Biz");
        tmp_Biz.reset();
        tmp_Biz.fromJson(_json);
        tmp_Biz.updateTime = now*1000;
        tmp_Biz.status = _status;

        //封装apply      
        LibLog.log("封装apply"); 
        tmp_SecPledgeApply.reset();       
        tmp_SecPledgeApply.fromJson(_json); 
       
        //分配一致的id，如id已存在，则选择更新
        LibLog.log("分配ID"); 
        uint id = 0;
        if (tmp_SecPledgeApply.id != 0) {
            id = tmp_SecPledgeApply.id;
        } else {
            id  = sq.getSeqNo("Biz.id");
        }
        tmp_SecPledgeApply.id = id;
        tmp_Biz.id = id;
        tmp_Biz.relatedId = id;
        tmp_SecPledgeApply.bizId = id;

        //证券数组对象加上编号
        LibLog.log("证券数组对象加上编号"); 
        for (uint i = 0; i < tmp_SecPledgeApply.appliedSecurities.length; i++ ) {
            tmp_SecPledgeApply.appliedSecurities[i].id = i+1;
        }

        //分配businessNo
        LibLog.log("分配businessNo"); 
        string memory businessNo = tmp_SecPledgeApply.businessNo;
        if (tmp_SecPledgeApply.businessNo.equals("")) {
            businessNo = sq.genBusinessNo("Q", "01").recoveryToString();
            tmp_SecPledgeApply.businessNo = businessNo;
            tmp_Biz.businessNo = businessNo;
        }

        if (isExistBiz(id)) { //更新数据
            //申请入库
            tmp_Biz.reset();
            tmp_Biz.fromJson(getBizJson(id));
            tmp_Biz.pledgors = tmp_SecPledgeApply.pledgors;
            tmp_Biz.pledgee = tmp_SecPledgeApply.pledgee;
            tmp_Biz.tradeOperator = tmp_SecPledgeApply.tradeOperator;
            tmp_Biz.updateTime = now*1000;
            tmp_Biz.status = _status;
            
            uint daoResult = orderDao.update_SecPledgeApply(_json);
            LibLog.log("_saveSecPledgeApplyByBroker: update_SecPledgeApply: ", daoResult);
            
            //biz入库
            LibLog.log("update_Biz: ", tmp_Biz.toJson());
            daoResult = orderDao.update_Biz(tmp_Biz.toJson());
            LibLog.log("_saveSecPledgeApplyByBroker: update_Biz: ", daoResult);
        } else { //插入新数据
            //申请入库
            uint insertResult = orderDao.insert_SecPledgeApply(tmp_SecPledgeApply.toJson());
            LibLog.log("_saveSecPledgeApplyByBroker: insert_SecPledgeApply: ", insertResult);
            
            //biz入库
            tmp_Biz.createTime = now*1000;
            tmp_Biz.startTime = now*1000;
            tmp_Biz.bizType = LibBiz.BizType.PLEDGE_BIZ;
            tmp_Biz.channelType = uint(LibBiz.ChannelType.BY_BROKER);
            LibLog.log("insert_Biz: ", tmp_Biz.toJson());
            insertResult = orderDao.insert_Biz(tmp_Biz.toJson());
            LibLog.log("_saveSecPledgeApplyByBroker: insert_Biz: ", insertResult);
        }

        //提交申请时将已提交状态的办理人（经办人）写到mapping中
        if(_status == LibBiz.BizStatus.BROKER_PLEDGE_UNAUDITED) { 
            orderDao.set_audit_of_status(tmp_Biz.id, uint(LibBiz.BizStatus.BROKER_PLEDGE_CREATED), tmp_SecPledgeApply.tradeOperator.id);
        }


        string memory success = "";
        uint _total = 1;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString() );
        success = success.concat(",\"businessNo\":\"",  businessNo,"\"");
        success = success.concat(",\"items\":[]");
        success = success.concat(",\"id\":", id.toString(),  "}}");
        Notify(0, success);
        
        return true;

    }

    //17.7.质押申请撤销(券商)
    function cancelSecPledgeApplyByBroker(uint id) getOrderDao returns(bool _ret) {
        if (checkAndDeleteCache(id)) {
            notify(0);
        } else {
            notify(errno_prefix + uint(OperError.DAO_ERROR));
        }
    }

    //检验某个biz是否存在，如存在则删除
    function checkAndDeleteCache(uint id) internal returns(bool _ret) {
        _ret = false;
        string memory _json = findBizById(id);
        LibLog.log(_json);

        LibJson.push(_json);
        if (_json.jsonKeyExists("data.items[0]")) {
            LibLog.log("key exists");
            uint deleteResult = orderDao.delete_Biz_byId(id);
            LibLog.log("checkAndDeleteCache: " , deleteResult);
            deleteResult = orderDao.delete_SecPledgeApply_byId(id);
            LibLog.log("checkAndDeleteCache: " , deleteResult);
            
            _ret = true;
        }
        LibJson.pop();
    }

    //检验某个记录是否存在
    function isExistBiz(uint id) internal returns(bool _ret) {
        _ret = false;
        string memory _json = findBizById(id);
        LibLog.log(_json);

        LibJson.push(_json);
        if (_json.jsonKeyExists("data.items[0]")) {
            LibLog.log("key exists");
            _ret = true;
        }
        LibJson.pop();
    }

    function getBizJson(uint id) getOrderDao constant returns (string _ret) {
        _ret = findBizById(id);
        LibJson.push(_ret);
        _ret = _ret.jsonRead("data.items[0]");
        LibJson.pop();
    }


    // 根据id查询申请单
    function findById(uint id) getOrderDao constant returns (string _ret) {
       uint len = orderDao.select_SecPledgeApply_byId(id);
       _ret = LibStack.popex(len);
    }

    // 根据id查询Biz单据
    function findBizById(uint id) getOrderDao constant returns (string _ret) {
       uint len = orderDao.select_Biz_byId(id);
       _ret = LibStack.popex(len);
    }

    //17.2.质押申请保存(中证登)
    //生成businessNo并状态置为CSDC_PLEDGE_CREATED
    function cacheSecPledgeApplyByCsdc(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("cacheSecPledgeApplyByCsdc: ", _json);
        return _saveSecPledgeApplyByCsdc(_json, LibBiz.BizStatus.CSDC_PLEDGE_CREATED);
   }

  
    //17.4质押申请创建(中证登，作废)
    function createSecPledgeApplyByCsdc(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("cacheSecPledgeApplyByCsdc: ", _json);
        return _saveSecPledgeApplyByCsdc(_json, LibBiz.BizStatus.CSDC_PLEDGE_CREATED);
    }

 
    //17.6.质押申请提交(中证登)
    //状态改为待审核
    function saveSecPledgeApplyByCsdc(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("cacheSecPledgeApplyByCsdc: ", _json);
        return _saveSecPledgeApplyByCsdc(_json, LibBiz.BizStatus.CSDC_PLEDGE_UNAUDITED);
    }

    //用于保存中证登录入的内部方法
    function _saveSecPledgeApplyByCsdc(string _json, LibBiz.BizStatus _status) internal returns(bool _ret) {
        LibLog.log("cacheSecPledgeApplyByCsdc: ", _json);

        //封装Biz
        tmp_Biz.reset();
        tmp_Biz.fromJson(_json);
        tmp_Biz.updateTime = now*1000;
        tmp_Biz.status = _status;

        //封装apply       
        tmp_SecPledgeApply.reset();       
        tmp_SecPledgeApply.fromJson(_json); 
       
        //分配一致的id，如id已存在，则选择更新
        uint id = 0;
        if (tmp_SecPledgeApply.id != 0) {
            id = tmp_SecPledgeApply.id;
        } else {
            id  = sq.getSeqNo("Biz.id");
        }
        tmp_SecPledgeApply.id = id;
        tmp_Biz.id = id;
        tmp_Biz.relatedId = id;
        tmp_SecPledgeApply.bizId = id;

        //证券数组对象加上编号
        for (uint i = 0; i < tmp_SecPledgeApply.appliedSecurities.length; i++ ) {
            tmp_SecPledgeApply.appliedSecurities[i].id = i+1;
        }

        //分配businessNo
        string memory businessNo = tmp_SecPledgeApply.businessNo;
        if (tmp_SecPledgeApply.businessNo.equals("")) {
            businessNo = sq.genBusinessNo("G", "01").recoveryToString();
            tmp_SecPledgeApply.businessNo = businessNo;
            tmp_Biz.businessNo = businessNo;
        }

        if (isExistBiz(id)) { //更新数据
            tmp_Biz.reset();
            tmp_Biz.fromJson(getBizJson(id));
            tmp_Biz.pledgors = tmp_SecPledgeApply.pledgors;
            tmp_Biz.pledgee = tmp_SecPledgeApply.pledgee;
            tmp_Biz.tradeOperator = tmp_SecPledgeApply.tradeOperator;
            tmp_Biz.updateTime = now*1000;
            tmp_Biz.status = _status;
            // tmp_Biz.id = tmp_SecPledgeApply.id;

            //申请入库
            uint daoResult = orderDao.update_SecPledgeApply(_json);
            LibLog.log("cacheSecPledgeApplyByCsdc: update_SecPledgeApply: ", daoResult);
            
            //biz入库
            LibLog.log("update_Biz: ", tmp_Biz.toJson());
            daoResult = orderDao.update_Biz(tmp_Biz.toJson());
            LibLog.log("cacheSecPledgeApplyByCsdc: update_Biz: ", daoResult);
        } else { //插入新数据
            //申请入库
            uint insertResult = orderDao.insert_SecPledgeApply(tmp_SecPledgeApply.toJson());
            LibLog.log("cacheSecPledgeApplyByCsdc: insert_SecPledgeApply: ", insertResult);
            
            //biz入库
            tmp_Biz.createTime = now*1000;
            tmp_Biz.startTime = now*1000;
            tmp_Biz.bizType = LibBiz.BizType.PLEDGE_BIZ;
            tmp_Biz.channelType = uint(LibBiz.ChannelType.BY_CSDC);

            LibLog.log("insert_Biz: ", tmp_Biz.toJson());
            insertResult = orderDao.insert_Biz(tmp_Biz.toJson());
            LibLog.log("cacheSecPledgeApplyByCsdc: insert_Biz: ", insertResult);
        }

        //提交申请时将已提交状态的办理人（经办人）写到mapping中
        if(_status == LibBiz.BizStatus.CSDC_PLEDGE_UNAUDITED) { 
            orderDao.set_audit_of_status(tmp_Biz.id, uint(LibBiz.BizStatus.CSDC_PLEDGE_CREATED), tmp_SecPledgeApply.tradeOperator.id);
        }

        string memory success = "";
        uint _total = 1;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString() );
        success = success.concat(",\"businessNo\":\"",  businessNo,"\"");
        success = success.concat(",\"items\":[]");
        success = success.concat(",\"id\":", id.toString(),  "}}");
        Notify(0, success);
        
        return true;

    }



    //17.8.质押申请撤销(中证登)
    function cancelSecPledgeApplyByCsdc(uint id) getOrderDao returns(bool _ret) {
        if (checkAndDeleteCache(id)) {
            notify(0);
        } else {
            notify(errno_prefix + uint(OperError.DAO_ERROR));
        }
    }

    //17.10.券商复审
    function reviewByBroker(uint id, uint operateCode, string auditInfoJson) getOrderDao returns(bool _ret) {
        if (operateCode == uint(OperateCode.PASS) ) {
            orderDao.update_Biz_Status(id, uint(LibBiz.BizStatus.BROKER_PLEDGE_CSDC_UNAUDITED));
            orderDao.add_Biz_Audit(id, auditInfoJson);
            Notify(0, 'success');
        } else if (operateCode == uint(OperateCode.REJECT)) {
            //质权人已拒绝
            orderDao.update_Biz_Status(id, uint(LibBiz.BizStatus.BROKER_PLEDGE_UNRETYPED));
            orderDao.add_Biz_Audit(id, auditInfoJson);
            orderDao.update_Biz_rejectStatus(id, uint(LibBiz.RejectStatus.BROKER_REJECT));
            Notify(0, 'success');
        }
        
    }

    //17.11.中证登审核券商订单
    function reviewByCsdcFromBroker(uint id, uint operateCode, uint status, string auditInfoJson) getOrderDao returns(bool _ret) {
        uint result = orderDao.update_Biz_Status(id, status);
        if (result != 0) {
            notify(result);
            return;
        }
        result = orderDao.add_Biz_Audit(id, auditInfoJson);
        if(operateCode == uint(OperateCode.REJECT)) {
            orderDao.update_Biz_rejectStatus(id, uint(LibBiz.RejectStatus.CSDC_REJECT));
        }
        notify(result);
    }

    //17.12.中证登审核订单
    function reviewByCsdc(uint id, uint operateCode, uint status, string auditInfoJson) getOrderDao returns(bool _ret) {
        uint result = orderDao.update_Biz_Status(id, status);
        if (result != 0) {
            notify(result);
            return;
        }
        result = orderDao.add_Biz_Audit(id, auditInfoJson);
        notify(result);
    }

    //17.13券商查看订单结果
    function checkByBroker(uint id, string auditInfoJson) getOrderDao returns(bool _ret) {
        uint result = orderDao.update_Biz_CheckStatus(id, uint(LibBiz.CheckStatus.BROKER_READ_CSDC_UNREAD));
        if (result != 0) {
            notify(result);
            return;
        }
        result = orderDao.add_Biz_Audit(id, auditInfoJson);
        notify(result);
    }

    //17.14.中证登查看订单结果
    function checkByCsdc(uint id, string auditInfoJson) getOrderDao returns(bool _ret) {
        uint result = orderDao.update_Biz_CheckStatus(id, uint(LibBiz.CheckStatus.BROKER_UNREAD_CSDC_READ));
        if (result != 0) {
            notify(result);
            return;
        }
        result = orderDao.add_Biz_Audit(id, auditInfoJson);
        notify(result);
    }



    //18.14.生成付款通知(中证登)
    function generatePaymentNotice(uint id, string _json) getOrderDao getPaymentManager returns(bool _ret) {
        LibLog.log("NewSecPledgeApplyManager generatePaymentNotice invoked.");
        uint daoResult = orderDao.update_SecPledgeApply_ById(id, _json);
        if (daoResult != 0) {
            notify(daoResult);
            return;
        }
        daoResult = orderDao.update_Biz_Status(id, uint(LibBiz.BizStatus.CSDC_PLEDGE_PAYNOTIFY_CREATED));

        tmp_SecPledgeApply.fromJson(_json);

        string memory _json1 = getPaymentById(id);
        if(_json1.equals("")) {
            //不存在付费信息
            LibLog.log("the payment does not exist");

            tmp_payment.reset();
            tmp_payment.id = tmp_SecPledgeApply.id;
            tmp_payment.relatedId = tmp_SecPledgeApply.id;
            tmp_payment.amount = tmp_SecPledgeApply.payAmount;
            daoResult = paymentManager.createPayment(tmp_payment.toJson());

        }else{
            //存在付费信息
            LibLog.log("the payment exists");

            tmp_payment.fromJson(_json1);
            tmp_payment.amount = tmp_SecPledgeApply.payAmount;
            tmp_payment.status = uint(LibPayment.PaymentStatus.UNPAID);
            daoResult = orderDao.update_Payment(tmp_payment.toJson());
        }
        LibLog.log("daoResult: ", daoResult.toString());

        notify(daoResult);
     }

    //18.15.设置付款信息(中证登)
    function setPaymentInfo(uint id, string _json) getOrderDao getPaymentManager returns(bool _ret) {
        LibLog.log("NewSecPledgeApplyManager setPaymentInfo invoked.");
        uint daoResult = orderDao.update_SecPledgeApply_ById(id, _json);
        if (daoResult != 0) {
            notify(daoResult);
            return;
        }
        daoResult = orderDao.update_Biz_Status(id, uint(LibBiz.BizStatus.CSDC_PLEDGE_PAYMENT_CONFIRMED));

        tmp_SecPledgeApply.fromJson(_json);

        string memory _json1 = getPaymentById(id);
        if(_json1.equals("")) {
            //不存在付费信息
            LibLog.log("the payment does not exist");

            tmp_payment.reset();
            tmp_payment.id = tmp_SecPledgeApply.id;
            tmp_payment.relatedId = tmp_SecPledgeApply.id;
            tmp_payment.amount = tmp_SecPledgeApply.payAmount;

            tmp_payment.time = now*1000;
            tmp_payment.status = uint(LibPayment.PaymentStatus.PAID);
            
            tmp_payment.flow = tmp_SecPledgeApply.payFlow;
            tmp_payment.paymentType = tmp_SecPledgeApply.payType;
            // tmp_payment.account = tmp_SecPledgeApply.payerName;
            if(tmp_SecPledgeApply.receivedAmount == 0) {
                tmp_payment.receivedAmount = tmp_payment.amount;
            } else {
                tmp_payment.receivedAmount = tmp_SecPledgeApply.receivedAmount;
            }
            daoResult = orderDao.insert_Payment(tmp_payment.toJson());

        }else{
            //存在付费信息
            LibLog.log("the payment exists");

            tmp_payment.fromJson(_json1);
            tmp_payment.time = now*1000;
            tmp_payment.status = uint(LibPayment.PaymentStatus.PAID);
            
            tmp_payment.flow = tmp_SecPledgeApply.payFlow;
            tmp_payment.paymentType = tmp_SecPledgeApply.payType;
            // tmp_payment.account = tmp_SecPledgeApply.payerName;
            if(tmp_SecPledgeApply.receivedAmount == 0) {
                tmp_payment.receivedAmount = tmp_payment.amount;
            } else {
                tmp_payment.receivedAmount = tmp_SecPledgeApply.receivedAmount;
            }

            daoResult = orderDao.update_Payment(tmp_payment.toJson());
        }
        LibLog.log("daoResult: ", daoResult.toString());

        notify(daoResult);
    }

    //增加后端附件
    function addBizAttachment(uint _bizId, string _json) getOrderDao {
        uint result = orderDao.add_Biz_Attachment(_bizId, _json);
        notify(result);
    }

    //增加前端附件
    function addApplyAttachment(uint _bizId, string _json) getOrderDao {
        uint result = orderDao.add_SecpledgeApply_Attachment(_bizId, _json);
        notify(result);
    }

    //17.9.AS400冻结解冻处理
    function updatePledgeApplyByAS400ResultNew(uint id, uint operateCode, uint status, string auditInfoJson, string secPledgeJson, string invoiceJson) getOrderDao returns(bool _ret) {
        LibLog.log("updatePledgeApplyByAS400ResultNew: auditInfoJson", auditInfoJson);
        LibLog.log("updatePledgeApplyByAS400ResultNew: secPledgeJson", secPledgeJson);
        LibLog.log("updatePledgeApplyByAS400ResultNew: invoiceJson", invoiceJson);
        
        uint result = orderDao.update_Biz_Status(id, status);
        if (result != 0) {
            notify(result);
            return;
        }
        if (!auditInfoJson.equals("")) {
            result = orderDao.add_Biz_Audit(id, auditInfoJson);
            LibLog.log("updatePledgeApplyByAS400ResultNew: add_Biz_Audit: ", result);
        }
        
        // tmp_SecPledge.reset();
        // tmp_SecPledge.fromJson(secPledgeJson);
        if (!secPledgeJson.equals("")) {
            result = orderDao.insert_SecPledge(secPledgeJson);
            LibLog.log("updatePledgeApplyByAS400ResultNew: insert_SecPledge: ", result);
        }
 
        notify(result);
    }

    function notify(uint result) internal {
        if (result == 0) {
            Notify(result, "success");
        } else {
            Notify(result, "error");
        }
    }
    struct StateTrans{
        LibBiz.BizStatus current;
        OperateCode operateCode;
        LibBiz.BizStatus next;
    }

    function listAll() getOrderDao constant returns(string) {
        uint len = orderDao.select_SecPledgeApply_all();
        return LibStack.popex(len);
    }

    function getPaymentById(uint _id) constant internal returns(string _ret) {
       uint len = orderDao.select_Payment_byId(_id);
       _ret = LibStack.popex(len);

       LibJson.push(_ret);
       _ret = _ret.jsonRead("data.items[0]");
       LibJson.pop();
    }
}
