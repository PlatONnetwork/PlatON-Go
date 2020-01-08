pragma solidity ^0.4.12;
/**
*@file      CommonContract.sol
*@author    yiyating
*@time      2017-08-15
*@desc      the defination of CommonContract contract
*/


import "../OrderDao.sol";

contract CommonContract is OwnerNamed {

    using LibJson for *;
    using LibAudit for *;

    OrderDao od;

    LibAudit.Audit t_audit;

    modifier getOrderDao(){ 
        od = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
        _;
    }

    function getBizById(uint id) getOrderDao constant returns (string _ret) {
        uint len = od.select_Biz_byId(id);
        _ret = LibStack.popex(len);
        
        LibJson.push(_ret);
        _ret = _ret.jsonRead("data.items[0]");
        LibJson.pop();
    }

    function getSecPledgeApplyById(uint id) getOrderDao constant returns (string _ret) {
        uint len = od.select_SecPledgeApply_byId(id);
        _ret = LibStack.popex(len);
        
        LibJson.push(_ret);
        _ret = _ret.jsonRead("data.items[0]");
        LibJson.pop();
    }

    function getDisSecPledgeApplyById(uint id) getOrderDao constant returns (string _ret) {
        uint len = od.select_DisSecPledgeApply_byId(id);
        _ret = LibStack.popex(len);
        
        LibJson.push(_ret);
        _ret = _ret.jsonRead("data.items[0]");
        LibJson.pop();
    }

    function getSecPledgeById(uint id) getOrderDao constant returns (string _ret) {
        uint len = od.select_SecPledge_byId(id);
        _ret = LibStack.popex(len);
        
        LibJson.push(_ret);
        _ret = _ret.jsonRead("data.items[0]");
        LibJson.pop();
    }

    function getBizStatus(uint _bizId) getOrderDao constant returns (uint) {
      return od.select_Biz_Status_ById(_bizId);
    }

    // 记录审计结果
    function addAudit(
      uint _bizId,
      address _auditorId,
      uint _operateCode, 
      string _auditComment,
      string _methodName,
      uint _oldStatus,
      uint _status
    ) getOrderDao returns (bool) {
        if (/* 以下是【解质押】流程信息 */
            _methodName.equals("createDisPledgeApply") ||
            _methodName.equals("createDisPledgeApplyInstitution") ||
            _methodName.equals("updateDisPledgeApplyByPledgeeFaceAuth") ||
            _methodName.equals("updateDisPledgeApplyByPledgorFaceAuth") ||
            _methodName.equals("updateDisPledgeApplyByAdmin") ||
            _methodName.equals("pledgorSubmit") ||
            _methodName.equals("pledgorSubmitInstitution") ||
            _methodName.equals("pledgeeSubmit") ||
            _methodName.equals("pledgeeSubmitInstitution") ||

            /* 以下是【质押】流程信息 */
            _methodName.equals("createPledgeApplyCommon") || 
            _methodName.equals("updatePledgeApplyByPledgorFaceAuth") || 
            _methodName.equals("updatePledgeApplyByPledgeeFaceAuth") ||
            _methodName.equals("updatePledgeApplyByPayment") ||
            _methodName.equals("updatePledgeApplyByAdmin") || 
            _methodName.equals("updatePledgeApplyByPledgee") ||
            _methodName.equals("updatePledgeApplyByPledgeeInstitution")
        ) {
            t_audit.reset();
            t_audit.auditorId = _auditorId;
            t_audit.auditComment = _auditComment;
            t_audit.operateCode = LibAudit.OperateCode(_operateCode);
            t_audit.auditTime = now*1000;
            t_audit.oldStatus = _oldStatus;
            t_audit.status = _status;
            od.add_Biz_Audit(_bizId, t_audit.toJson());
        }
        return true;
    }

}