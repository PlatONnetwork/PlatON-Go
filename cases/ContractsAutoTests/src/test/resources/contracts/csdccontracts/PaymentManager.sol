pragma solidity ^0.4.12;

import "./OrderDao.sol";
import "./Sequence.sol";

contract PaymentManager is OwnerNamed {

    using LibInt for *;
    using LibString for *;
    using LibPayment for *;
    using LibAudit for *;
    using LibJson for *;
    using LibRejectHistory for *;

    enum PaymentError {
        NO_ERROR,
        JSON_INVALID,
        BAD_PARAMETER,
        STATUS_ERROR,
        OPERATE_CODE_ERROR,
        DAO_ERROR
    }

    OrderDao od;
    Sequence sq;
    LibPayment.Payment t_payment;
    LibAudit.Audit t_audit;
    LibRejectHistory.RejectHistory t_rejectHistory;

    uint errno_prefix = 16600;

    modifier getSq(){ 
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0")); 
        _;
    }

    modifier getOrderDao(){ 
        od = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0")); 
        _;
    }

    function PaymentManager(){
        register("CsdcModule", "0.0.1.0", "PaymentManager", '0.0.1.0');
    }

    function createPayment(string _json) getSq getOrderDao returns (uint _ret) {
        LibLog.log("PaymentManager createPayment invoked.");
        LibLog.log("json", _json);

        if(!t_payment.fromJson(_json)) {
            LibLog.log("json invalid");
            return;
        }
        // t_payment.id = sq.getSeqNo("Payment.id");
        t_payment.status = uint(LibPayment.PaymentStatus.UNPAID);
        t_payment.refundApplyStatus = uint(LibPayment.RefundApplyStatus.UNAPPLIED);
        _ret = od.insert_Payment(t_payment.toJson());
        LibLog.log("call dao result:", _ret.toString());
    }

    function updatePaymentStatus(string _json) getOrderDao {
        LibLog.log("PaymentManager updatePaymentStatus invoked.");
        LibLog.log("json", _json);

        if(!t_payment.fromJson(_json)) {
            LibLog.log("json invalid");
            Notify(errno_prefix + uint(PaymentError.JSON_INVALID), "json invalid");
            return;
        }

        uint _ret = od.update_Payment(_json);

        if(_ret != 0) {
            LibLog.log("call dao error", _ret.toString());
            Notify(errno_prefix + uint(PaymentError.DAO_ERROR), "call dao error");
        } else {
            Notify(0, "success");
        }
    }

    function refundPayment(uint _id) returns (string _ret)  {
        // t_payment.reset();
        // t_payment = keyMap[_id];
        // if(t_payment.id == 0){
        //     Notify(10000, "payment not exist");
        //     return ;
        // }
        // t_payment.status  = LibPayment.PaymentStatus(3);
        // keyMap[_id] = t_payment;
        // Notify(0, "success");
    }

    //申请退款
    function applyRefund(uint _id, string _reason) getOrderDao {
        string memory _json = getPaymentById(_id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(PaymentError.BAD_PARAMETER), "the payment id does not exist");
            return;
        }
        t_payment.fromJson(_json);

        if(t_payment.status != uint(LibPayment.PaymentStatus.PAID) || 
           (t_payment.refundApplyStatus != uint(LibPayment.RefundApplyStatus.UNAPPLIED) && 
            t_payment.refundApplyStatus != uint(LibPayment.RefundApplyStatus.AUDIT_FAILED) && 
            t_payment.refundApplyStatus != uint(LibPayment.RefundApplyStatus.REVIEW_FAILED))
        ) {
            LibLog.log("payment status error");
            Notify(errno_prefix+uint(PaymentError.STATUS_ERROR), "payment status error");
            return;
        }

        t_payment.refundApplyReason = _reason;
        t_payment.refundApplyStatus = uint(LibPayment.RefundApplyStatus.UNAUDITED);
        t_payment.refundApplyTime = now*1000;
        t_payment.updateTime = now*1000;

        uint _ret = od.update_Payment(t_payment.toJson());

        if(_ret != 0) {
            LibLog.log("call dao error", _ret.toString());
            Notify(errno_prefix + uint(PaymentError.DAO_ERROR), "call dao error");
        } else {
            Notify(0, "success");
        }
    }

    //退款初审
    function auditRefund(uint _id, string _auditInfo) getOrderDao {
        LibLog.log("PaymentManager auditRefund invoked.");
        LibLog.log("id", _id);
        LibLog.log("auditInfo", _auditInfo);

        if(!t_audit.fromJson(_auditInfo)) {
            LibLog.log("auditorInfo invalid");
            Notify(errno_prefix+uint(PaymentError.JSON_INVALID), "auditorInfo invalid");
            return;
        }

        string memory _json = getPaymentById(_id);
        if(_json.equals("")) {
            LibLog.log("the payment id does not exist");
            Notify(errno_prefix+uint(PaymentError.BAD_PARAMETER), "the payment id does not exist");
            return;
        }
        t_payment.fromJson(_json);

        if(t_payment.refundApplyStatus != uint(LibPayment.RefundApplyStatus.UNAUDITED)) {
            LibLog.log("refundApplyStatus error");
            Notify(errno_prefix+uint(PaymentError.STATUS_ERROR), "refundApplyStatus error");
            return;
        }

        uint _status;
        if(t_audit.operateCode == LibAudit.OperateCode.PASS) {
            _status = uint(LibPayment.RefundApplyStatus.UNREVIEWED);
        } else if(t_audit.operateCode == LibAudit.OperateCode.FAIL) {
            _status = uint(LibPayment.RefundApplyStatus.AUDIT_FAILED);
        } else {
            LibLog.log("operateCode error");
            Notify(errno_prefix+uint(PaymentError.OPERATE_CODE_ERROR), "operateCode error");
            return;
        }
        t_audit.oldStatus = uint(t_payment.refundApplyStatus);
        t_audit.status = _status;
        t_audit.auditTime = now*1000;

        //写入拒绝历史记录
        if(t_audit.operateCode == LibAudit.OperateCode.FAIL) {
            t_rejectHistory.fromJson(t_audit.toJson());
            t_rejectHistory.applyTime = t_payment.refundApplyTime;
            t_rejectHistory.applyReason = t_payment.refundApplyReason;
            t_payment.rejectHistory.push(t_rejectHistory);
        }

        t_payment.audits.push(t_audit);
        t_payment.refundApplyStatus = _status;
        t_payment.updateTime = now*1000;

        uint _ret = od.update_Payment(t_payment.toJson());
        if(_ret != 0) {
            LibLog.log("call dao error", _ret.toString());
            Notify(errno_prefix + uint(PaymentError.DAO_ERROR), "call dao error");
        } else {
            Notify(0, "success");
        }
    }

    //退款复核
    function reviewRefund(uint _id, string _auditInfo) getOrderDao {
        LibLog.log("PaymentManager reviewRefund invoked.");
        LibLog.log("id", _id);
        LibLog.log("auditInfo", _auditInfo);

        if(!t_audit.fromJson(_auditInfo)) {
            LibLog.log("auditorInfo invalid");
            Notify(errno_prefix+uint(PaymentError.JSON_INVALID), "auditorInfo invalid");
            return;
        }

        string memory _json = getPaymentById(_id);
        if(_json.equals("")) {
            LibLog.log("the payment id does not exist");
            Notify(errno_prefix+uint(PaymentError.BAD_PARAMETER), "the payment id does not exist");
            return;
        }
        t_payment.fromJson(_json);

        if(t_payment.refundApplyStatus != uint(LibPayment.RefundApplyStatus.UNREVIEWED)) {
            LibLog.log("refundApplyStatus error");
            Notify(errno_prefix+uint(PaymentError.STATUS_ERROR), "refundApplyStatus error");
            return;
        }

        uint _status;
        if(t_audit.operateCode == LibAudit.OperateCode.PASS) {
            _status = uint(LibPayment.RefundApplyStatus.UNCONFIRMED);
        } else if(t_audit.operateCode == LibAudit.OperateCode.FAIL) {
            _status = uint(LibPayment.RefundApplyStatus.REVIEW_FAILED);
        } else {
            LibLog.log("operateCode error");
            Notify(errno_prefix+uint(PaymentError.OPERATE_CODE_ERROR), "operateCode error");
            return;
        }
        t_audit.oldStatus = uint(t_payment.refundApplyStatus);
        t_audit.status = _status;
        t_audit.auditTime = now*1000;

        //写入拒绝历史记录
        if(t_audit.operateCode == LibAudit.OperateCode.FAIL) {
            t_rejectHistory.fromJson(t_audit.toJson());
            t_rejectHistory.applyTime = t_payment.refundApplyTime;
            t_rejectHistory.applyReason = t_payment.refundApplyReason;
            t_payment.rejectHistory.push(t_rejectHistory);
        }

        t_payment.audits.push(t_audit);
        t_payment.refundApplyStatus = _status;
        t_payment.updateTime = now*1000;

        uint _ret = od.update_Payment(t_payment.toJson());
        if(_ret != 0) {
            LibLog.log("call dao error", _ret.toString());
            Notify(errno_prefix + uint(PaymentError.DAO_ERROR), "call dao error");
        } else {
            Notify(0, "success");
        }
    }

    //确认已退款
    function confirmRefund(uint _id, string _auditInfo) getOrderDao {
        LibLog.log("PaymentManager confirmRefund invoked.");
        LibLog.log("id", _id);
        LibLog.log("auditInfo", _auditInfo);
        
        if(!t_audit.fromJson(_auditInfo)) {
            LibLog.log("auditorInfo invalid");
            Notify(errno_prefix+uint(PaymentError.JSON_INVALID), "auditorInfo invalid");
            return;
        }

        string memory _json = getPaymentById(_id);
        if(_json.equals("")) {
            LibLog.log("the payment id does not exist");
            Notify(errno_prefix+uint(PaymentError.BAD_PARAMETER), "the payment id does not exist");
            return;
        }
        t_payment.fromJson(_json);

        if(t_payment.refundApplyStatus != uint(LibPayment.RefundApplyStatus.UNCONFIRMED)) {
            LibLog.log("refundApplyStatus error");
            Notify(errno_prefix+uint(PaymentError.STATUS_ERROR), "refundApplyStatus error");
            return;
        }
        // t_audit.reset();
        // t_audit.operateCode = _operateCode;
        // t_audit.auditComment = _auditComment;
        // t_audit.auditorId = _auditorId;
        

        uint _status;
        if(t_audit.operateCode == LibAudit.OperateCode.PASS) {
            _status = uint(LibPayment.RefundApplyStatus.REFUNDED);
            t_payment.status = uint(LibPayment.PaymentStatus.REFUNDED);
        } else {
            LibLog.log("operateCode error");
            Notify(errno_prefix+uint(PaymentError.OPERATE_CODE_ERROR), "operateCode error");
            return;
        }
        t_audit.oldStatus = uint(t_payment.refundApplyStatus);
        t_audit.status = _status;
        t_audit.auditTime = now*1000;

        t_payment.audits.push(t_audit);
        t_payment.refundApplyStatus = _status;
        t_payment.updateTime = now*1000;

        uint _ret = od.update_Payment(t_payment.toJson());
        if(_ret != 0) {
            LibLog.log("call dao error", _ret.toString());
            Notify(errno_prefix + uint(PaymentError.DAO_ERROR), "call dao error");
        } else {
            Notify(0, "success");
        }
    }

    // LibPayment.Condition _cond;

    function findByCond(string _json) getOrderDao constant returns (string _ret) {
        uint len = od.pagePayment(_json);
        return LibStack.popex(len);
    }

    function findById(uint _id) getOrderDao constant returns (string _ret) {
        uint len = od.select_Payment_byId(_id);
        return LibStack.popex(len);
    }

    /* 以下是内部调用接口 */

    function getPaymentById(uint _id) getOrderDao constant internal returns(string _ret) {
       _ret = findById(_id);
       LibJson.push(_ret);
       _ret = _ret.jsonRead("data.items[0]");
       LibJson.pop();
    }

    event Notify(uint _errorno, string _info);
}