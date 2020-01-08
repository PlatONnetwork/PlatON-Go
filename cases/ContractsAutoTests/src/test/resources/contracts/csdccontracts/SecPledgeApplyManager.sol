pragma solidity ^0.4.12;

import "./csdc_library/LibBilling.sol";
import "./csdc_base/CommonContract.sol";

import "./BizManager.sol";
import "./PaymentManager.sol";
import "./SecPledgeManager.sol";
import "./PerUserManager.sol";
import "./OrgUserManager.sol";

contract SecPledgeApplyManager  is CommonContract  {
    using LibSecPledgeApply for *;
    using LibPledgeSecurity for *;
    using LibString for *;
    using LibInt for *;
    using LibBiz for *;
    using LibPayment for *;
    using LibSecPledge for *;
    using LibTradeUser for *;
    using LibPerUser for *;
    using LibOrgUser for *;
    using LibBiz for *;
    using LibJson for *;

     //inner setting member
    LibSecPledgeApply.SecPledgeApply internal tmp_SecPledgeApply;
    LibPayment.Payment _paymentTemp;
    LibSecPledge.SecPledge _secPledge;
    LibPledgeSecurity.PledgeSecurity _tmpSec;

    LibTradeUser.TradeUser t_tradeUser;
    LibPerUser.PerUser t_perUser;
    LibOrgUser.OrgUser t_orgUser;
    LibBiz.Biz t_biz;

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
        BUSINESSNO_ERROR, //业务编号获取失败
        DAO_ERROR
    }

    enum OperateCode{
        NONE, 
        PASS, //通过
        REJECT, //拒绝
        WAIT //等待处理
    }
    uint errno_prefix = 10000;
    
    /* only for calculate payment */
     mapping(string=>uint) num0;    //每个证券代码的优先股股数
     mapping(string=>uint) num1;    //每个证券代码的其他股股数
     string[] public codes;         //证券代码list

    Sequence sq;
    BizManager bizManager;
    PaymentManager paymentManager;
    SecPledgeManager secPledgeManager;
    PerUserManager pm;
    OrgUserManager om;

    modifier getNewestContract(){
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
        bizManager = BizManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "BizManager", "0.0.1.0"));
        paymentManager = PaymentManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PaymentManager","0.0.1.0"));
        secPledgeManager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager", "0.0.1.0"));
        pm = PerUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PerUserManager", "0.0.1.0"));
        om = OrgUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrgUserManager","0.0.1.0"));
        od = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao","0.0.1.0"));
        _;
    }

    modifier getSq(){ sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0")); _;}

    modifier getBizManager(){ bizManager = BizManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "BizManager", "0.0.1.0")); _;}

    modifier getPaymentManager(){ paymentManager = PaymentManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PaymentManager","0.0.1.0")); _;}

    modifier getSecPledgeManager(){ secPledgeManager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager", "0.0.1.0")); _;}

    modifier getPm(){ pm = PerUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PerUserManager", "0.0.1.0")); _;}

    modifier getOm(){ om = OrgUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrgUserManager","0.0.1.0")); _;}

    function SecPledgeApplyManager() {
        register("CsdcModule", "0.0.1.0", "SecPledgeApplyManager", "0.0.1.0");
    }

    function updateSecPledgeApply(string _json) getOrderDao {
        LibLog.log("update_SecPledgeApply: ", _json);
        if (od.update_SecPledgeApply(_json) != 0) {
          Notify(errno_prefix + uint(SecPledgeApplyError.DAO_ERROR), "call dao error");
          return;
        }
        Notify(0, "success");
    }

    // 创建申请
    function createPledgeApplyCommon(string code, string _json, uint _status) returns(uint _id) {

        LibLog.log("createPledgeApplyCommon: json", _json);

        if(!tmp_SecPledgeApply.fromJson(_json)){
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER),"json invalid");
            return;
        }

        LibLog.log("tmp_SecPledgeApply: ", tmp_SecPledgeApply.toJson());

        if (tmp_SecPledgeApply.pledgorId == address(0) || tmp_SecPledgeApply.pledgeeId == address(0)) {
            Notify(errno_prefix+uint(SecPledgeApplyError.ID_EMPTY),"pledgor or pledgee can not be null");
            return;         
        }

        if (tmp_SecPledgeApply.payerType != LibSecPledgeApply.PayerType.PLEDGOR && tmp_SecPledgeApply.payerType != LibSecPledgeApply.PayerType.PLEGDEE) {
            Notify(errno_prefix+uint(SecPledgeApplyError.PAYERTYPE_EMPTY),"PAYER_TYPE is wrong");
            return;             
        }

        if (tmp_SecPledgeApply.pledgorId == tmp_SecPledgeApply.pledgeeId) {
            Notify(errno_prefix+uint(SecPledgeApplyError.PLEDGOR_PLEDGEE_SAME),"Pledgor and Pledgee can not be the same person");
            return;             
        }

        userToTrader(t_tradeUser, tmp_SecPledgeApply.pledgeeId);
        tmp_SecPledgeApply.pledgee = t_tradeUser;

        userToTrader(t_tradeUser, tmp_SecPledgeApply.pledgorId);
        tmp_SecPledgeApply.pledgors.push(t_tradeUser);
        //将质押证券信息中的证券账户号copy到出质人信息中
        tmp_SecPledgeApply.pledgors[0].account = tmp_SecPledgeApply.appliedSecurities[0].secAccount;


        //检查证券信息的数据
        for (uint i = 0; i < tmp_SecPledgeApply.appliedSecurities.length; i++) {
               //是否解除红利字段合法
            if (tmp_SecPledgeApply.appliedSecurities[i].isProfit != LibPledgeSecurity.IsProfit.T && tmp_SecPledgeApply.appliedSecurities[i].isProfit != LibPledgeSecurity.IsProfit.F) {
                Notify(errno_prefix+uint(SecPledgeApplyError.APPLIEDSECURITIES_ERROR),"APPLIEDSECURITIES isProfit is illegal");
                return; 
            }

            if (tmp_SecPledgeApply.appliedSecurities[i].pledgeNum <= 0) {
                Notify(errno_prefix+uint(SecPledgeApplyError.APPLIEDSECURITIES_ERROR),"APPLIEDSECURITIES pledgeNum is illegal");
                return; 
            }

            if ( tmp_SecPledgeApply.appliedSecurities[i].secAccount.equals("")
                    || tmp_SecPledgeApply.appliedSecurities[i].secCode.equals("")
                    || tmp_SecPledgeApply.appliedSecurities[i].secName.equals("")
                    //|| tmp_SecPledgeApply.appliedSecurities[i].secType  == _tmpSec.secType 
                    || tmp_SecPledgeApply.appliedSecurities[i].hostedUnit.equals("")
                    || tmp_SecPledgeApply.appliedSecurities[i].hostedUnitName.equals("")
                    || tmp_SecPledgeApply.appliedSecurities[i].secProperty.equals("") ) {
                        Notify(errno_prefix+uint(SecPledgeApplyError.APPLIEDSECURITIES_ERROR),"APPLIEDSECURITIES property can not be null");
                        return; 
                    }
            tmp_SecPledgeApply.appliedSecurities[i].shareholderIdNo = tmp_SecPledgeApply.pledgors[0].idNo;
            tmp_SecPledgeApply.appliedSecurities[i].id = i+1;
    
        }

        uint _ret = 0;
        string memory businessNo;
        if(tmp_SecPledgeApply.id == 0) {

            //json中不带id，新业务
            businessNo = sq.genBusinessNo(code, "01").recoveryToString();
        } else {

            //重新发起的业务，保留原业务流水号
            LibLog.log("id exists", tmp_SecPledgeApply.id.toString());
            string memory bizJson = getBizById(tmp_SecPledgeApply.id);

            //业务流水号获取失败，业务已被重新发起
            if(bizJson.equals("")) {
                Notify(errno_prefix+uint(SecPledgeApplyError.BUSINESSNO_ERROR),"该业务已被重新发起");
                return;
            }
            LibJson.push(bizJson);
            businessNo = bizJson.jsonRead("businessNo");
            LibJson.pop();

            //删除原业务
            _ret = od.delete_Biz_byId(tmp_SecPledgeApply.id);
            LibLog.log("delete_Biz_byId:", _ret.toString());
            _ret = od.delete_SecPledgeApply_byId(tmp_SecPledgeApply.id);
            LibLog.log("delete_SecPledgeApply_byId:", _ret.toString());
        }

        tmp_SecPledgeApply.businessNo = businessNo;

        // string memory businessNo = sq.genBusinessNo(code, "01").recoveryToString();
        tmp_SecPledgeApply.pledgeContractNo =  "CSDC-ZXZY".concat(businessNo);
        
        tmp_SecPledgeApply.pledgorName = LibStack.popex(pm.findNameById(tmp_SecPledgeApply.pledgorId));
        tmp_SecPledgeApply.pledgeeName = LibStack.popex(pm.findNameById(tmp_SecPledgeApply.pledgeeId));

        //设置支付信息
        address _payerAccount = address(0);
        if (tmp_SecPledgeApply.payerType == LibSecPledgeApply.PayerType.PLEDGOR) {
            _payerAccount = tmp_SecPledgeApply.pledgorId;
        } else {
            _payerAccount = tmp_SecPledgeApply.pledgeeId;
        }
        tmp_SecPledgeApply.payerAccount = _payerAccount;
        tmp_SecPledgeApply.payerName = LibStack.popex(pm.findNameById(_payerAccount));
        tmp_SecPledgeApply.payAmount = __calcPayment(tmp_SecPledgeApply.appliedSecurities);
        
        //创建申请记录
        _id = sq.getSeqNo("Biz.id");
        tmp_SecPledgeApply.id = _id;
        tmp_SecPledgeApply.bizId = _id;
        tmp_SecPledgeApply.applyTime = now*1000;

        //* 【warning】
        //  Biz和Apply
        //  使用同一个id
        //*/
        t_biz.reset();
        t_biz.id = t_biz.relatedId = _id;
        t_biz.pledgorId = tmp_SecPledgeApply.pledgorId;
        t_biz.pledgeeId = tmp_SecPledgeApply.pledgeeId;
        t_biz.pledgorName = tmp_SecPledgeApply.pledgorName;
        t_biz.pledgeeName = tmp_SecPledgeApply.pledgeeName;
        t_biz.bizType = LibBiz.BizType.PLEDGE_BIZ;
        t_biz.status = LibBiz.BizStatus(_status);
        t_biz.businessNo = businessNo;
        t_biz.pledgeContractNo = "CSDC-ZXZY".concat(businessNo);
        t_biz.channelType = uint(LibBiz.ChannelType.ONLINE);
        t_biz.startTime = now*1000;
        t_biz.createTime = now*1000;
        t_biz.updateTime = now*1000;
        bizJson = t_biz.toJson();
        LibLog.log("insertBiz: ", bizJson);
        _ret = od.insert_Biz(bizJson);
        LibLog.log("insert_biz result: ", _ret);
        addAudit(_id, tmp_SecPledgeApply.pledgorId, uint(LibAudit.OperateCode.PASS), "", "createPledgeApplyCommon", uint(LibBiz.BizStatus.NONE), _status);
        // bizManager.insert(_id, tmp_SecPledgeApply.pledgorId, tmp_SecPledgeApply.pledgeeId, address(0), uint(LibBiz.BizType.PLEDGE_BIZ), _status, _id);
        // bizManager.setPledgeContractNo(_id, "CSDC-ZXZY".concat(businessNo));
        // bizManager.setBusinessNo(_id, businessNo);

        _ret = od.insert_SecPledgeApply(tmp_SecPledgeApply.toJson());
        LibLog.log("insert_SecPledgeApply result: ", _ret);


        string memory success = "";
        uint _total = 1;
        //success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString(), ",\"businessNo\":", businessNo, ",\"items\":[");
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString() );
        success = success.concat(",\"businessNo\":\"",  businessNo,"\"");
        success = success.concat(",\"items\":[");
        success = success.concat("{\"id\":", _id.toString(),  "}]}}");
        Notify(0, success);
    }
    
    
    // (个人)创建申请
    function createPledgeApply(string _json) getNewestContract returns(uint _id) {
        // string memory businessNo = sq.genBusinessNo("A", "01").recoveryToString(); //手机客户端质押申请
        LibLog.log("------------------------start----------------------");
        _id = createPledgeApplyCommon("A", _json, uint(LibBiz.BizStatus.PLEDGE_PLEDGOR_UNAUTH));
        LibLog.log("------------------------end----------------------| id=",_id.toString());
    }


          // (机构)创建申请
    function createPledgeApplyInstitution(string _json) getNewestContract returns(uint _id) {
        // string memory businessNo = sq.genBusinessNo("W", "01").recoveryToString(); //web客户端质押申请
        _id = createPledgeApplyCommon("W", _json, uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNCONFIRMED));
    }
    
    
    //(机构)质权人确认
    function updatePledgeApplyByPledgeeInstitution(uint id, uint operateCode, string rejectReason) getNewestContract {
        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        uint bizId = tmp_SecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        //处理出质人申请的情况
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNCONFIRMED)) { 
            if (operateCode == uint(OperateCode.PASS) ) {
                
                if (tmp_SecPledgeApply.payAmount == 0) {    //支付金额为0时直接进入待审核流程
                    bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_UNAUDITED));
                    _status = LibBiz.BizStatus.PLEDGE_UNAUDITED;

                } else if (tmp_SecPledgeApply.payerType == LibSecPledgeApply.PayerType.PLEDGOR) {
                    bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGOR_UNPAID));//2-待出质人支付
                    _status = LibBiz.BizStatus. PLEDGE_PLEDGOR_UNPAID;
                } else {
                    bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNPAID));//15-待质权人支付
                    _status = LibBiz.BizStatus. PLEDGE_PLEDGEE_UNPAID;
                }

                address _payerAccount = tmp_SecPledgeApply.payerAccount;
                uint amount = tmp_SecPledgeApply.payAmount; //__calcPayment(tmp_SecPledgeApply.id);
                uint paymentId = tmp_SecPledgeApply.id;
                _addPayment(amount, tmp_SecPledgeApply.id, _payerAccount);
                
                bizManager.addPaymentId(bizId, paymentId);

                // secPledgeApplyMap[id].paymentId = paymentId;
                tmp_SecPledgeApply.paymentId = paymentId;
                od.update_SecPledgeApply(tmp_SecPledgeApply.toJson());

            } else if (operateCode == uint(OperateCode.REJECT)) {
                //质权人拒绝
                bizManager.endBiz(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_DENIED));  
                _status = LibBiz.BizStatus.PLEDGE_PLEDGEE_DENIED;
            }
            addAudit(bizId, tmp_SecPledgeApply.pledgeeId, operateCode, rejectReason, "updatePledgeApplyByPledgeeInstitution", _oldStatus, uint(_status));
            Notify(0, 'success');
        } else{
            Notify(errno_prefix+uint(SecPledgeApplyError.OPERATE_NOT_ALLOWED), 'Status is wrong');
        }
    }


    //合同流转-质权人通过/拒绝
    function updatePledgeApplyByPledgee(uint id, uint operateCode, string rejectReason) getBizManager {

        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        uint bizId = tmp_SecPledgeApply.bizId;
        uint _oldStatus = getBizStatus(bizId);
        uint _status;
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNCONFIRMED)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                //待质权人人脸认证
                bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNAUTH));
                _status = uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNAUTH);
                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //质权人已拒绝
                bizManager.endBiz(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_DENIED));
                _status = uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_DENIED);
                Notify(0, 'success');
            }
            addAudit(bizId, tmp_SecPledgeApply.pledgeeId, operateCode, rejectReason, "updatePledgeApplyByPledgee", _oldStatus, _status);
        } else { 

            Notify(errno_prefix+uint(SecPledgeApplyError.OPERATE_NOT_ALLOWED), 'Status is wrong');
        }
    }

     //合同流转-质权人确认人脸确认
    function updatePledgeApplyByPledgeeFaceAuth(uint id, uint operateCode, string rejectReason) getNewestContract {
        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        uint bizId = tmp_SecPledgeApply.bizId;
        // t_biz = getBizById(bizId);
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        //处理质权人通过的情况
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNAUTH)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                address _payerAccount = address(0);
                if (tmp_SecPledgeApply.payAmount == 0) {    
                    //支付金额为0时直接进入待审核流程
                     _status = LibBiz.BizStatus.PLEDGE_UNAUDITED;

                } else if (tmp_SecPledgeApply.payerType == LibSecPledgeApply.PayerType.PLEDGOR) {  //待出质人支付，也可以质权人支付
                     _status = LibBiz.BizStatus.PLEDGE_PLEDGOR_UNPAID;
                     _payerAccount = tmp_SecPledgeApply.pledgorId;
                     
                } else {
                     _status = LibBiz.BizStatus.PLEDGE_PLEDGEE_UNPAID;
                     _payerAccount = tmp_SecPledgeApply.pledgeeId;
                }
                bizManager.changeStatus(bizId, uint(_status));

                uint amount = tmp_SecPledgeApply.payAmount; //__calcPayment(tmp_SecPledgeApply.id);
                uint paymentId = tmp_SecPledgeApply.id;
                _addPayment(amount, tmp_SecPledgeApply.id, _payerAccount);
                
                bizManager.addPaymentId(bizId, paymentId);

                tmp_SecPledgeApply.paymentId = paymentId;
                od.update_SecPledgeApply(tmp_SecPledgeApply.toJson());
                
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //人脸识别失败，质权人人脸验证失败
                _status = LibBiz.BizStatus.PLEDGE_PLEDGEE_AUTH_FAILED;
                bizManager.endBiz(bizId, uint(_status));
            }
            addAudit(bizId, tmp_SecPledgeApply.pledgeeId, operateCode, rejectReason, "updatePledgeApplyByPledgeeFaceAuth", _oldStatus, uint(_status));
            Notify(0, 'success');
            return;
        } 

        //处理质权人之前拒绝的情况 PLEDGE_PLEDGEE_DENIED
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_DENIED)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                //不需要再做操作
                bizManager.endBiz(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_AUTH_FAILED));
            }
            Notify(0, 'success');
            return;
        }
        Notify(errno_prefix+uint(SecPledgeApplyError.OPERATE_NOT_ALLOWED), 'Status is wrong');
    }

    //合同流转-申请人脸确认
    function updatePledgeApplyByPledgorFaceAuth(uint id, uint operateCode, string rejectReason) getBizManager {
        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        uint bizId = tmp_SecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        //处理出质人申请的情况
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGOR_UNAUTH) || _oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGOR_UNSIGNED)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                //进入待质权人认证的状态
                bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNCONFIRMED));
                _status = LibBiz.BizStatus.PLEDGE_PLEDGEE_UNCONFIRMED;

                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //人脸识别失败
                bizManager.endBiz(bizId, uint(LibBiz.BizStatus.PLEDGE_PLEDGOR_AUTH_FAILED));
                _status = LibBiz.BizStatus.PLEDGE_PLEDGOR_AUTH_FAILED;
                Notify(0, 'success');
            }
            addAudit(bizId, tmp_SecPledgeApply.pledgorId, operateCode, rejectReason, "updatePledgeApplyByPledgorFaceAuth", _oldStatus, uint(_status));
        } 
        //Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "no operate done");
    }

    //合同流转-付款完成/超时
    function updatePledgeApplyByPayment(uint id, uint operateCode, string rejectReason) getBizManager {
        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        uint bizId = tmp_SecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGOR_UNPAID) || _oldStatus == uint(LibBiz.BizStatus.PLEDGE_PLEDGEE_UNPAID) ) {
            if (operateCode == uint(OperateCode.PASS) ) {
                //付款通过进入待审核
                bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_UNAUDITED));
                _status = LibBiz.BizStatus.PLEDGE_UNAUDITED;

                __updatePayment(id, tmp_SecPledgeApply.payAmount);
                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //付款失败
                bizManager.endBiz(bizId, uint(LibBiz.BizStatus.PLEDGE_TIMEOUT_FAILED));
                _status = LibBiz.BizStatus.PLEDGE_TIMEOUT_FAILED;
                Notify(0, 'success');
            }
            addAudit(bizId, address(0), operateCode, rejectReason, "updatePledgeApplyByPayment", _oldStatus, uint(_status));
           
        } else { 

            Notify(errno_prefix+uint(SecPledgeApplyError.OPERATE_NOT_ALLOWED), 'Status is wrong');
        }
    }

    //合同流转-管理员审核通过/拒绝
    function updatePledgeApplyByAdmin(uint id, uint operateCode, address auditorId, string rejectReason) getBizManager {
        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        uint bizId = tmp_SecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_UNAUDITED) ) {
            if (operateCode == uint(OperateCode.PASS) ) {
                //审核通过
                bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_AUDIT_SUCCESS));
                _status = LibBiz.BizStatus.PLEDGE_AUDIT_SUCCESS;
                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //审核拒绝
                bizManager.endBiz(bizId, uint(LibBiz.BizStatus.PLEDGE_AUDIT_FAILED));
                _status = LibBiz.BizStatus.PLEDGE_AUDIT_FAILED;
                Notify(0, 'success');
            }
           addAudit(bizId, auditorId, operateCode, rejectReason, "updatePledgeApplyByAdmin", _oldStatus, uint(_status));
        } else { 

            Notify(errno_prefix+uint(SecPledgeApplyError.OPERATE_NOT_ALLOWED), 'Status is wrong');
        }
    }

    //合同流转-退款
    function refund(uint id, string rejectReason) getBizManager getPaymentManager {
        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        uint bizId = tmp_SecPledgeApply.bizId;
        uint paymentId = tmp_SecPledgeApply.paymentId;
        uint _oldStatus = getBizStatus(bizId);
        //只有审核拒绝和冻结失败时可以置退款
        if (_oldStatus == uint(LibBiz.BizStatus.PLEDGE_AUDIT_FAILED) || _oldStatus == uint(LibBiz.BizStatus.PLEDGE_PROCESS_FAILED)) {
            bizManager.changeStatus(bizId, uint(LibBiz.BizStatus.PLEDGE_REFUND));
            paymentManager.refundPayment(paymentId);
            Notify(0, 'success');
           
        } else { 

            Notify(errno_prefix+uint(SecPledgeApplyError.OPERATE_NOT_ALLOWED), 'Status is wrong');
        }
    }
   
    // 根据id查询
    function findById(uint id) getOrderDao constant returns (string _ret) {
        uint len = od.select_SecPledgeApply_byId(id);
        return LibStack.popex(len);
    }
    
    //添加质押合同文件id和名称
    function updatePledgeContractFile(uint id, string pledgeContractFileId, string pledgeContractFileName) getNewestContract returns (bool) {
        string memory _json = getSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(SecPledgeApplyError.BAD_PARAMETER), "the apply id does not exist");
            return;
        }

        tmp_SecPledgeApply.fromJson(_json);

        bizManager.setPledgeContracFile(id, pledgeContractFileId, pledgeContractFileName);

        // secPledgeApplyMap[id].pledgeContractFileId = pledgeContractFileId;
        // secPledgeApplyMap[id].pledgeContractFileName = pledgeContractFileName;
        tmp_SecPledgeApply.pledgeContractFileId = pledgeContractFileId;
        tmp_SecPledgeApply.pledgeContractFileName = pledgeContractFileName;
        od.update_SecPledgeApply(tmp_SecPledgeApply.toJson());
        
        Notify(0, 'success');
        return true;
    }

    event Notify(uint _errorno, string _info);
    
    /* 以下是内部调用接口 */
    function __calcPayment(LibPledgeSecurity.PledgeSecurity[] appliedSecurities) internal returns (uint _amount) {
        for(uint i = 0; i < appliedSecurities.length; i++) {
            if(num0[appliedSecurities[i].secCode]==0 && num1[appliedSecurities[i].secCode]==0) {
                codes.push(appliedSecurities[i].secCode);
            }
            if(appliedSecurities[i].secType.equals('45')) {// 优先股
                num0[appliedSecurities[i].secCode] += appliedSecurities[i].pledgeNum;
            } else {   // 其他股
                num1[appliedSecurities[i].secCode] += appliedSecurities[i].pledgeNum;
            }
        }
        for (i = 0; i <codes.length; i++) {
            _amount += LibBilling.calc_secPledge(num0[codes[i]], 10000);    //优先股每股100元 = 10000分
            _amount += LibBilling.calc_secPledge(num1[codes[i]], 100);      //其他股每股1元 = 100分

            num0[codes[i]] = num1[codes[i]] = 0;    //清空数据，待下次调用
        }
        delete codes;
    }

    //增加支付请求
    function _addPayment(uint amount, uint id, address account) getPaymentManager internal returns (uint paymentId) {
    //  Payment{
    //  uint id;                //id
    //  string flow;            //交易流水
    //  string payChannel;      //付费渠道
    //  uint amount;            //交易金额
    //  PaymentStatus status;   //交易状态
    //  PaymentType paymentType;        //缴费方式
    //  uint account;           //缴费账号
    //  uint relatedId;         //合约申请表示
    // }
        _paymentTemp.reset();
        _paymentTemp.amount = amount;
        // _paymentTemp.receivedAmount = amount;
        _paymentTemp.relatedId = id;
        _paymentTemp.id = id;
        // _paymentTemp.account = uint(account);

        paymentManager.createPayment(_paymentTemp.toJson());
    }

    function __updatePayment(uint _id, uint _amount) getPaymentManager internal {
        string memory _json = "{";
        _json = _json.jsonCat("id", _id);
        _json = _json.jsonCat("receivedAmount", _amount);
        _json = _json.concat("}");

        paymentManager.updatePaymentStatus(_json);
    }

    function userToTrader(LibTradeUser.TradeUser storage _traderUser, address _userId) internal {
        if(pm.userExists(_userId) == 1) {
            //出质人为个人
            t_perUser.fromJson(LibStack.popex(pm.findById(_userId)));

            _traderUser.reset();
            _traderUser.traderType = uint(LibTradeUser.TraderType.PERUSER);
            _traderUser.traderId = t_perUser.id;
            _traderUser.idNo = t_perUser.idNo;
            _traderUser.idType = uint(t_perUser.idType);
            _traderUser.userType = 1;   //境内自然人
            _traderUser.name = t_perUser.name;
        } else if(om.userExists(_userId) == 1) {
            //出质人为机构
            t_orgUser.fromJson(LibStack.popex(om.findById(_userId)));

            _traderUser.reset();
            _traderUser.traderType = uint(LibTradeUser.TraderType.ORGUSER);
            _traderUser.traderId = t_orgUser.id;
            _traderUser.idNo = t_orgUser.businessLicenseNo;
            _traderUser.idType = uint(t_orgUser.organIdType);
            _traderUser.userType = 2;   //境内法人
            _traderUser.name = t_orgUser.organFullName;
        }
    }

}