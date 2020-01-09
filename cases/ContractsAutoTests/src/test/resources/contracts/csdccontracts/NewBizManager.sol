pragma solidity ^0.4.12;

import "./Sequence.sol";
import "./OrderDao.sol";

contract NewBizManager is OwnerNamed {
    using LibInt for *;
    using LibString for *;
    using LibBiz for *;
    using LibAudit for *;
    
    mapping(uint=>LibBiz.Biz) bizs;
    uint[] bizIds;
    
    Sequence sq;
    OrderDao od;
    
    uint errno_prefix = 17000;

    enum BizError {
      NONE,
      DAO_ERROR
    }

    modifier getOrderDao(){ 
      od = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
      _;
    }

    function NewBizManager() {
        register("CsdcModule", "0.0.1.0", "NewBizManager", '0.0.1.0');
    }

    function updateBiz(string _json) getOrderDao returns (uint) {
        LibLog.log("NewBizManager updateBiz invoked.");
        LibLog.log("json: ", _json);
        if (od.update_Biz(_json) != 0) {
          Notify(errno_prefix + uint(BizError.DAO_ERROR), "call dao error");
        }
        Notify(0, "success");
    }

    function changeStatus(uint _bizId, uint _status) getOrderDao returns (uint) {
        LibLog.log("NewBizManager changeStatus invoked.");
        LibLog.log("bizId: ", _bizId.toString());
        LibLog.log("status: ", _status.toString());
        if (od.update_Biz_Status(_bizId, _status) != 0) {
          Notify(errno_prefix + uint(BizError.DAO_ERROR), "call dao error");
        }
        Notify(0, "success");
    }

    function update_Biz_CheckStatus(uint _bizId, uint _status) getOrderDao returns (uint) {
        LibLog.log("NewBizManager update_Biz_CheckStatus invoked.");
        LibLog.log("bizId: ", _bizId.toString());
        LibLog.log("status: ", _status.toString());
        if (od.update_Biz_CheckStatus(_bizId, _status) != 0) {
          Notify(errno_prefix + uint(BizError.DAO_ERROR), "call dao error");
        }
        Notify(0, "success");
    }

    function addAudit(uint _bizId, string _json) getOrderDao returns(uint _ret) {
        if (od.add_Biz_Audit(_bizId, _json) != 0) {
          Notify(errno_prefix + uint(BizError.DAO_ERROR), "call dao error");
        }
        Notify(0, "success");
    }

    function addBizAttachment(uint _bizId, string _json) getOrderDao {
        if (od.add_Biz_Attachment(_bizId, _json) != 0) {
          Notify(errno_prefix + uint(BizError.DAO_ERROR), "call dao error");
        }
        Notify(0, "success");
    }

    function getStatusById(uint _bizId) getOrderDao constant returns(uint) {
        return od.select_Biz_Status_ById(_bizId);
    }

    function pageBizForCsdc(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager pageBizForCsdc invoked.");
        LibLog.log("json: ", _json);
        uint len = od.pageBizForCsdc(_json);
        return LibStack.popex(len);
    }

    function pageBizForBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager pageBizForBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.pageBizForBroker(_json);
        return LibStack.popex(len);
    }

    function findToDoForBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findToDoForBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findToDoForBroker(_json);
        return LibStack.popex(len);
    }

    function findHandledForBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findHandledForBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findHandledForBroker(_json);
        return LibStack.popex(len);
    }

    function findAllBizForCsdcByBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findAllBizForCsdcByBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findAllBizForCsdcByBroker(_json);
        return LibStack.popex(len);
    }

    function findUnAuditedBizForCsdcByBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findUnAuditedBizForCsdcByBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findUnAuditedBizForCsdcByBroker(_json);
        return LibStack.popex(len);
    }

    function findUnReviewedBizForCsdcByBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findUnReviewedBizForCsdcByBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findUnReviewedBizForCsdcByBroker(_json);
        return LibStack.popex(len);
    }

    function findUnCheckedBizForCsdcByBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findUnCheckedBizForCsdcByBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findUnCheckedBizForCsdcByBroker(_json);
        return LibStack.popex(len);
    }

    function findUnReviewedBizForCsdcLeaderByBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findUnReviewedBizForCsdcLeaderByBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findUnReviewedBizForCsdcLeaderByBroker(_json);
        return LibStack.popex(len);
    }

    function findClosedBizForCsdcByBroker(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager findClosedBizForCsdcByBroker invoked.");
        LibLog.log("json: ", _json);
        uint len = od.findClosedBizForCsdcByBroker(_json);
        return LibStack.popex(len);
    }

    function findById(uint _id) getOrderDao constant returns(string) {
        uint len = od.select_Biz_byId(_id);
        return LibStack.popex(len);
    }

    function extractPledgeByPledgeRegisterNo(string _no) getOrderDao constant returns(string) {
        uint len = od.extractPledgeByPledgeRegisterNo(_no);
        return LibStack.popex(len);
    }

    function listAll() getOrderDao constant returns(string) {
        uint len = od.select_Biz_all();
        return LibStack.popex(len);
    }

    function getToDoIdListByBroker(address _userId, uint _brokerId, uint _role) getOrderDao constant returns(string) {
        uint len = od.getToDoIdListByBroker(_userId, _brokerId, _role);
        return LibStack.popex(len);
    }

    function getToDoIdListByCsdcFromBroker(address _userId) getOrderDao constant returns(string) {
        uint len = od.getToDoIdListByCsdcFromBroker(_userId);
        return LibStack.popex(len);
    }

    function pageBizForProxy(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager pageBizForProxy invoked.");
        LibLog.log("json: ", _json);
        uint len = od.pageBizForProxy(_json);
        return LibStack.popex(len);
    }

    function pageBiz(string _json) getOrderDao constant returns(string) {
        LibLog.log("NewBizManager pageBiz invoked.");
        LibLog.log("json: ", _json);
        uint len = od.pageBiz(_json);
        return LibStack.popex(len);
    }
    
    event Notify(uint _errorno, string _info);

  //   function findNameById(address _id) constant returns (uint) {}
}