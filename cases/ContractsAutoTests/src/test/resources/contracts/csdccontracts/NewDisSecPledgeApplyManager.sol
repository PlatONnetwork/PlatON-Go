pragma solidity ^0.4.12;

import "./Sequence.sol";
import "./OrderDao.sol";

contract NewDisSecPledgeApplyManager  is OwnerNamed  {
    using LibDisSecPledgeApply for *;
    using LibString for *;
    using LibInt for *;
    using LibBiz for *;
    using LibSecPledge for *;
    using LibJson for *;

    OrderDao orderDao;
    Sequence sq;
    //inner setting member
    LibDisSecPledgeApply.DisSecPledgeApply internal tmp_DisSecPledgeApply;
    
    LibBiz.Biz internal tmp_Biz;

    event Notify(uint _errorno, string _info);
    
    /** @brief errno for test case */
    enum DisSecPledgeApplyError {
        NO_ERROR,
        BAD_PARAMETER,
        DAO_ERROR,
        OPERATE_NOT_ALLOWED,
        ID_EMPTY,   
        PAYERTYPE_EMPTY,
        INSERT_FAILED,
        USER_STATUS_ERROR,
        APPLIEDSECURITIES_ERROR,
        PLEDGOR_PLEDGEE_SAME,
        BUSINESSNO_ERROR, //业务编号获取失败
        STATUS_NOT_ALLOWED  //状态不对，不允许解除
    }

    enum OperateCode{
        NONE, 
        PASS, //通过
        REJECT, //拒绝
        WAIT //等待处理
    }
    uint errno_prefix = 10000;

    enum OperError {
        NO_ERROR,
        BAD_PARAMETER,
        DAO_ERROR
    }

    function NewDisSecPledgeApplyManager() {
        register("CsdcModule", "0.0.1.0", "NewDisSecPledgeApplyManager", "0.0.1.0");
        
        orderDao = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
    }

    modifier getOrderDao(){ 
      orderDao = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
      _;
    }

    //19.1.解除质押申请保存(中证登)
    function cacheDisSecPledgeApplyByCsdc(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("cacheDisSecPledgeApplyByCsdc: ", _json);
        _saveDisSecPledgeApplyByCsdc(_json, LibBiz.BizStatus.CSDC_DISPLEDGE_STASHED);

    }
  
    //19.2.解除质押申请提交(中证登)
    function saveDisSecPledgeApplyByCsdc(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("cacheDisSecPledgeApplyByCsdc: ", _json);
        _saveDisSecPledgeApplyByCsdc(_json, LibBiz.BizStatus.CSDC_DISPLEDGE_CREATED);

    }

    //19.3 解除质押申请撤销(中证登)
    function cancelDisSecPledgeApplyByCsdc(uint id) getOrderDao returns(bool _ret) {
        if (checkAndDeleteCache(id)) {
            notify(0);
        } else {
            notify(errno_prefix + uint(OperError.DAO_ERROR));
        }

    }

    //19.4.中证登复审解除质押订单
    function reviewByCsdc(uint id, uint operateCode, uint status, string auditInfoJson) getOrderDao returns(bool _ret) {
        uint result = orderDao.update_Biz_Status(id, status);
        if (result != 0) {
            notify(result);
            return;
        }
        result = orderDao.add_Biz_Audit(id, auditInfoJson);
        notify(result);
    }

    //19.5.AS400解冻处理
    function updatePledgeApplyByAS400ResultNew(uint id, uint operateCode, uint status, string auditInfoJson, string secPledgeJson) getOrderDao returns(bool _ret) {
        LibLog.log("updatePledgeApplyByAS400ResultNew", id); //, status, auditInfoJson, secPledgeJson);
        LibLog.log("updatePledgeApplyByAS400ResultNew", operateCode); //, status, auditInfoJson, secPledgeJson);
        LibLog.log("updatePledgeApplyByAS400ResultNew", status); //, status, auditInfoJson, secPledgeJson);
        LibLog.log("updatePledgeApplyByAS400ResultNew", auditInfoJson); //, status, auditInfoJson, secPledgeJson);
        LibLog.log("updatePledgeApplyByAS400ResultNew", secPledgeJson); //, status, auditInfoJson, secPledgeJson);
        
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
            result = orderDao.update_SecPledge(secPledgeJson);
        }
        notify(result);

    }

    function findById(uint _id) getOrderDao constant returns(string _ret) {
        uint len = orderDao.select_DisSecPledgeApply_byId(_id);
        return LibStack.popex(len);
    }

    function updateDisSecPledgeApply(string _json) getOrderDao returns (uint) {
        if (orderDao.update_DisSecPledgeApply(_json) != 0) {
          Notify(errno_prefix + uint(DisSecPledgeApplyError.DAO_ERROR), "call dao error");
          return;
        }
        LibJson.push(_json);
        string memory _ret = "";
        _ret = _ret.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\": 1");
        _ret = _ret.concat(",\"businessNo\":\"",  _json.jsonRead("businessNo"), "\"");
        _ret = _ret.concat(",\"items\":[]");
        _ret = _ret.concat(",\"id\":", _json.jsonRead("id"), "}}");
        LibJson.pop();
        Notify(0, _ret);
    }



    /* ******************************************************* */
    //17.1.质押申请保存(券商)
    //状态为BROKER_DISPLEDGE_STASHED
    function cacheDisSecPledgeApplyByBroker(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("cacheSecPledgeApplyByBroker: ", _json);
        _saveDisSecPledgeApplyByBroker(_json, LibBiz.BizStatus.BROKER_DISPLEDGE_STASHED);
    }

    //17.3.质押申请创建(券商)
    //状态为BROKER_DISPLEDGE_CREATED，产生businessNo
    function createDisSecPledgeApplyByBroker(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("createSecPledgeApplyByBroker: ", _json);
        _saveDisSecPledgeApplyByBroker(_json, LibBiz.BizStatus.BROKER_DISPLEDGE_CREATED);
    }

    //17.5.质押申请提交(券商)
    function saveDisSecPledgeApplyByBroker(string _json) getOrderDao returns(bool _ret) {
        LibLog.log("saveSecPledgeApplyByBroker: ", _json);
        _saveDisSecPledgeApplyByBroker(_json, LibBiz.BizStatus.BROKER_DISPLEDGE_UNAUDITED);
    }

    //内部用于存储券商申请的方法
    function _saveDisSecPledgeApplyByBroker(string _json, LibBiz.BizStatus _status) internal returns(bool _ret) {
        LibLog.log("_saveDisSecPledgeApplyByBroker: ", _json);

        //封装Biz
        LibLog.log("封装Biz");
        tmp_Biz.reset();
        tmp_Biz.fromJson(_json);
        tmp_Biz.updateTime = now*1000;
        tmp_Biz.status = _status;

        //封装apply      
        LibLog.log("封装apply"); 
        tmp_DisSecPledgeApply.reset();       
        tmp_DisSecPledgeApply.fromJson(_json); 

        if(_status == LibBiz.BizStatus.BROKER_DISPLEDGE_UNAUDITED && !__checkSecPlegeStatus(tmp_DisSecPledgeApply.secPledgeId)) {
            Notify(uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED), "质押状态有误，不允许解除");
            return;
        } 
       
        //分配一致的id，如id已存在，则选择更新
        LibLog.log("分配ID"); 
        uint id = 0;
        if (tmp_DisSecPledgeApply.id != 0) {
            id = tmp_DisSecPledgeApply.id;
        } else {
            id  = sq.getSeqNo("Biz.id");
        }
        tmp_DisSecPledgeApply.id = id;
        tmp_Biz.id = id;
        tmp_Biz.relatedId = id;
        tmp_DisSecPledgeApply.bizId = id;

        //证券数组对象加上编号
        LibLog.log("证券数组对象加上编号"); 
        for (uint i = 0; i < tmp_DisSecPledgeApply.appliedSecurities.length; i++ ) {
            tmp_DisSecPledgeApply.appliedSecurities[i].id = i+1;
        }

        //分配businessNo
        LibLog.log("分配businessNo"); 
        string memory businessNo = tmp_DisSecPledgeApply.businessNo;
        if (tmp_DisSecPledgeApply.businessNo.equals("")) {
            businessNo = sq.genBusinessNo("Q", "02").recoveryToString();
            tmp_DisSecPledgeApply.businessNo = businessNo;
            tmp_Biz.businessNo = businessNo;
        }

        if (__isExistBiz(id)) { //更新数据
            //申请入库
            tmp_Biz.reset();
            tmp_Biz.fromJson(getBizJson(id));
            tmp_Biz.pledgors = tmp_DisSecPledgeApply.pledgors;
            tmp_Biz.pledgee = tmp_DisSecPledgeApply.pledgee;
            tmp_Biz.tradeOperator = tmp_DisSecPledgeApply.tradeOperator;
            tmp_Biz.updateTime = now*1000;
            tmp_Biz.status = _status;
            
            uint daoResult = orderDao.update_DisSecPledgeApply(_json);
            LibLog.log("_saveDisSecPledgeApplyByBroker: update_DisSecPledgeApply: ", daoResult);
            
            //biz入库
            LibLog.log("update_Biz: ", tmp_Biz.toJson());
            daoResult = orderDao.update_Biz(tmp_Biz.toJson());
            LibLog.log("_saveDisSecPledgeApplyByBroker: update_Biz: ", daoResult);
        } else { //插入新数据
            //申请入库
            uint insertResult = orderDao.insert_DisSecPledgeApply(tmp_DisSecPledgeApply.toJson());
            LibLog.log("_saveDisSecPledgeApplyByBroker insert_DisSecPledgeApply: ", insertResult);
            
            //biz入库
            tmp_Biz.createTime = now*1000;
            tmp_Biz.startTime = now*1000;
            tmp_Biz.bizType = LibBiz.BizType.DISPLEDGE_BIZ;
            tmp_Biz.channelType = uint(LibBiz.ChannelType.BY_BROKER);
            LibLog.log("insert_Biz: ", tmp_Biz.toJson());
            insertResult = orderDao.insert_Biz(tmp_Biz.toJson());
            LibLog.log("_saveDisSecPledgeApplyByBroker: insert_Biz: ", insertResult);
        }

        if(_status == LibBiz.BizStatus.BROKER_DISPLEDGE_UNAUDITED) {
            orderDao.set_audit_of_status(tmp_Biz.id, uint(LibBiz.BizStatus.BROKER_DISPLEDGE_CREATED), tmp_DisSecPledgeApply.tradeOperator.id);
            orderDao.update_SecPledge_status(tmp_DisSecPledgeApply.secPledgeId, uint(LibSecPledge.PledgeStatus.DISPLEGING));
            orderDao.add_SecPledge_disId(tmp_DisSecPledgeApply.secPledgeId, tmp_DisSecPledgeApply.id);
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
    function cancelDisSecPledgeApplyByBroker(uint id) getOrderDao returns(bool _ret) {
        if (checkAndDeleteCache(id)) {
            notify(0);
        } else {
            notify(errno_prefix + uint(OperError.DAO_ERROR));
        }
    }

    //17.10.券商复审
    function reviewByBroker(uint id, uint operateCode, string auditInfoJson) getOrderDao returns(bool _ret) {
        if (operateCode == uint(OperateCode.PASS) ) {
            orderDao.update_Biz_Status(id, uint(LibBiz.BizStatus.BROKER_DISPLEDGE_CSDC_UNAUDITED));                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   
            orderDao.add_Biz_Audit(id, auditInfoJson);
            Notify(0, 'success');
        } else if (operateCode == uint(OperateCode.REJECT)) {
            //质权人已拒绝
            orderDao.update_Biz_Status(id, uint(LibBiz.BizStatus.BROKER_DISPLEDGE_UNRETYPED));
            orderDao.add_Biz_Audit(id, auditInfoJson);
            orderDao.update_Biz_rejectStatus(id, uint(LibBiz.RejectStatus.BROKER_REJECT));

            uint ret = orderDao.undo_SecPledgeStatus_byDis(id);
            LibLog.log("undo_SecPledgeStatus_byDis: " , ret);
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
    function checkByCsdcFromBroker(uint id, string auditInfoJson) getOrderDao returns(bool _ret) {
        uint result = orderDao.update_Biz_CheckStatus(id, uint(LibBiz.CheckStatus.BROKER_UNREAD_CSDC_READ));
        if (result != 0) {
            notify(result);
            return;
        }
        result = orderDao.add_Biz_Audit(id, auditInfoJson);
        notify(result);
    }



    /* 以下是internal接口 */

    // 根据id查询Biz单据
    function __findBizById(uint id) constant internal returns (string _ret) {
       uint len = orderDao.select_Biz_byId(id);
       _ret = LibStack.popex(len);
    }

    function getBizJson(uint id) constant returns (string _ret) {
        _ret = __findBizById(id);
        LibJson.push(_ret);
        _ret = _ret.jsonRead("data.items[0]");
        LibJson.pop();
    }

    //检验某个记录是否存在
    function __isExistBiz(uint id) internal returns(bool _ret) {
        _ret = false;
        string memory _json = __findBizById(id);
        LibLog.log(_json);

        LibJson.push(_json);
        if (_json.jsonKeyExists("data.items[0]")) {
            LibLog.log("key exists");
            _ret = true;
        }
        LibJson.pop();
    }

    //内部用于存储柜面申请解质押的方法
    function _saveDisSecPledgeApplyByCsdc(string _json, LibBiz.BizStatus _status) internal returns(bool _ret) {
        LibLog.log("_saveDisSecPledgeApplyByCsdc: ", _json);

        //封装Biz
        tmp_Biz.reset();
        tmp_Biz.fromJson(_json);
        tmp_Biz.updateTime = now*1000;
        tmp_Biz.status = _status;

        //封装apply       
        tmp_DisSecPledgeApply.reset();       
        tmp_DisSecPledgeApply.fromJson(_json);

        if(_status == LibBiz.BizStatus.CSDC_DISPLEDGE_CREATED && !__checkSecPlegeStatus(tmp_DisSecPledgeApply.secPledgeId)) {
            Notify(uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED), "质押状态有误，不允许解除");
            return;
        }
       
        //分配一致的id，如id已存在，则选择更新
        uint id = 0;
        if (tmp_DisSecPledgeApply.id != 0) {
            id = tmp_DisSecPledgeApply.id;
        } else {
            id  = sq.getSeqNo("Biz.id");
        }
        tmp_DisSecPledgeApply.id = id;
        tmp_Biz.id = id;
        tmp_Biz.relatedId = id;
        tmp_DisSecPledgeApply.bizId = id;

        //证券数组对象加上编号
        for (uint i = 0; i < tmp_DisSecPledgeApply.appliedSecurities.length; i++ ) {
            tmp_DisSecPledgeApply.appliedSecurities[i].id = i+1;
        }

        //分配businessNo
        string memory businessNo = tmp_DisSecPledgeApply.businessNo;
        if (tmp_DisSecPledgeApply.businessNo.equals("")) {
            businessNo = sq.genBusinessNo("G", "02").recoveryToString();
            tmp_DisSecPledgeApply.businessNo = businessNo;
            tmp_Biz.businessNo = businessNo;
        }

        if (__isExistBiz(id)) { //更新数据
            //申请入库
            uint daoResult = orderDao.update_DisSecPledgeApply(_json);
            LibLog.log("_saveDisSecPledgeApplyByCsdc: update_DisSecPledgeApply: ", daoResult);

            LibLog.log("update_Biz: ", tmp_Biz.toJson());
            daoResult = orderDao.update_Biz(tmp_Biz.toJson());
            LibLog.log("_saveDisSecPledgeApplyByBroker: update_Biz: ", daoResult);
            

        } else { //插入新数据
            //申请入库
            uint insertResult = orderDao.insert_DisSecPledgeApply(tmp_DisSecPledgeApply.toJson());
            LibLog.log("_saveDisSecPledgeApplyByCsdc: insert_DisSecPledgeApply: ", insertResult);
            
            //biz入库
            tmp_Biz.createTime = now*1000;
            tmp_Biz.startTime = now*1000;
            tmp_Biz.bizType = LibBiz.BizType.DISPLEDGE_BIZ;
            tmp_Biz.channelType = uint(LibBiz.ChannelType.BY_CSDC);
            
            LibLog.log("insert_Biz: ", tmp_Biz.toJson());
            insertResult = orderDao.insert_Biz(tmp_Biz.toJson());
            LibLog.log("_saveDisSecPledgeApplyByCsdc: insert_Biz: ", insertResult);
        }

        if(_status == LibBiz.BizStatus.CSDC_DISPLEDGE_CREATED) {
            orderDao.update_SecPledge_status(tmp_DisSecPledgeApply.secPledgeId, uint(LibSecPledge.PledgeStatus.DISPLEGING));
            orderDao.add_SecPledge_disId(tmp_DisSecPledgeApply.secPledgeId, tmp_DisSecPledgeApply.id);
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

       //检验某个biz是否存在，如存在则删除
    function checkAndDeleteCache(uint id) internal returns(bool _ret) {
        _ret = false;
        string memory _json = __findBizById(id);
        LibLog.log(_json);

        LibJson.push(_json);
        if (_json.jsonKeyExists("data.items[0]")) {
            LibLog.log("key exists");
            uint deleteResult = orderDao.undo_SecPledgeStatus_byDis(id);
            LibLog.log("undo_SecPledgeStatus_byDis: " , deleteResult);
            
            deleteResult = orderDao.delete_Biz_byId(id);
            LibLog.log("delete_Biz_byId: " , deleteResult);
            
            deleteResult = orderDao.delete_DisSecPledgeApply_byId(id);
            LibLog.log("delete_DisSecPledgeApply_byId: " , deleteResult);
            
            _ret = true;
        }
        LibJson.pop();
    }

    function __checkSecPlegeStatus(uint _secPledgeId) internal returns(bool _ret) {
        uint _status = orderDao.select_SecPledge_status_ById(_secPledgeId);
        if(_status == uint(LibSecPledge.PledgeStatus.PLEDGING) || _status == uint(LibSecPledge.PledgeStatus.PARTIAL_PLEDGED)) {
            return true;
        }
    }

    function notify(uint result) internal {
        if (result == 0) {
            Notify(result, "success");
        } else {
            Notify(result, "error");
        }
    }

    function listAll() getOrderDao constant returns(string) {
        uint len = orderDao.select_DisSecPledgeApply_all();
        return LibStack.popex(len);
    }
}
