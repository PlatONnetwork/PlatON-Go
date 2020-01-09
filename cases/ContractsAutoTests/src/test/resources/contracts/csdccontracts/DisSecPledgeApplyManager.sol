pragma solidity ^0.4.12;
/**
* @file LibSecPledge.sol
* @author zhangyu
* @time 2017-04-24
* @desc
*/


import "./csdc_base/CommonContract.sol";
import "./Sequence.sol";
import "./BizManager.sol";
import "./SecPledgeManager.sol";
import "./PerUserManager.sol";

contract DisSecPledgeApplyManager is CommonContract {
    using LibDisSecPledgeApply for *;
    using LibPledgeSecurity for *;
    using LibString for *;
    using LibInt for *;
    using LibBiz for *;
    using LibJson for *;

    //inner setting member
    LibDisSecPledgeApply.DisSecPledgeApply internal tmp_DisSecPledgeApply;

    LibSecPledge.SecPledge _secPledge;
    LibPledgeSecurity.PledgeSecurity _tmpSec;
    LibBiz.Biz t_biz;
    LibBiz.Biz t_biz_secpledgeApply;

    Sequence sq;
    BizManager bizManager;
    SecPledgeManager secPledgeManager;
    PerUserManager pm;
    uint _tmpbizId;
    
    /** @brief errno for test case */
    enum DisSecPledgeApplyError {
        NO_ERROR,
        BAD_PARAMETER,
        PledgeStatus_NOT_ALLOWED, //质押的状态处于处理中或已全部解除
        OPERATE_NOT_ALLOWED,
        PLEDGORID_EMPTY,   
        PLEDGEEID_EMPTY,
        SECPLEDGEID_ERROR, //不存在或不属于该pledgee
        APPLY_TYPE_NOT_ALLOWED, //解除类型不存在
        IS_FREEZE_INVALID, //司法冻结是否解除值无效
        APPLIEDSECURITIES_ERROR,  //解压证券错误，如不存在或数量已超出
        STATUS_NOT_ALLOWED_SUCH_METHOD, //该状态下不能走此流程
        BUSINESSNO_ERROR, //业务编号获取失败
        IS_FinancingAmountRemain,
        EVIDENCE_ALREADY_APPLIED, //证明文件已申请
        EVIDENCE_ALREADY_MAILED,   //证明文件已邮寄
        DAO_ERROR

    }

    enum OperateCode{
        NONE, 
        PASS, //通过
        REJECT, //拒绝
        WAIT //等待处理
    }
    uint errno_prefix = 10000;

    modifier getNewestContract(){
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
        bizManager = BizManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "BizManager", "0.0.1.0"));
        secPledgeManager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager", "0.0.1.0"));
        pm = PerUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PerUserManager", "0.0.1.0"));
        od = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao","0.0.1.0"));
        _;
    }

    modifier getSq(){ sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0")); _;}

    modifier getBizManager(){ bizManager = BizManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "BizManager", "0.0.1.0")); _;}

    modifier getSecPledgeManager(){ secPledgeManager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager", "0.0.1.0")); _;}

    modifier getPm(){ pm = PerUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PerUserManager", "0.0.1.0")); _;}

    function DisSecPledgeApplyManager() {
        register("CsdcModule", "0.0.1.0", "DisSecPledgeApplyManager", "0.0.1.0");
    }

    event Notify(uint _errorno, string _info);


    function updateDisSecPledgeApply(string _json) getOrderDao returns (uint) {
        LibLog.log("updateDisSecPledgeApply: ", _json);
        if (od.update_DisSecPledgeApply(_json) != 0) {
          Notify(errno_prefix + uint(DisSecPledgeApplyError.DAO_ERROR), "call dao error");
          return;
        }
        Notify(0, "success");
    }

    //新增解除质押申请业务
    function createDisPledgeApplyCommon(string code,string _json) getNewestContract internal returns(uint _id) {
        LibLog.log("createDisPledgeApplyCommon: ", _json);
        if(!tmp_DisSecPledgeApply.fromJson(_json)){
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER),"json invalid");
            return;
        }

        //uint secPledgeId = tmp_DisSecPledgeApply.secPledgeId;
        //成功质押记录存在且各个用户id匹配
        if (!secPledgeManager.isMatching(tmp_DisSecPledgeApply.secPledgeId, tmp_DisSecPledgeApply.pledgorId, tmp_DisSecPledgeApply.pledgeeId)) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.SECPLEDGEID_ERROR),"SecPledge does not exist or pledgor/pledgee id is wrong");
            return; 
        }

        //检查质押成交记录状态,只有质押中和部分质押中的可以申请解除
        if (!secPledgeManager.isStatus(tmp_DisSecPledgeApply.secPledgeId, uint(LibSecPledge.PledgeStatus.PLEDGING))
             && !secPledgeManager.isStatus(tmp_DisSecPledgeApply.secPledgeId, uint(LibSecPledge.PledgeStatus.PARTIAL_PLEDGED))
        ) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.PledgeStatus_NOT_ALLOWED),"SecPledge status is not for dis");
            return; 
        }

        // 校验解冻证券是否已质押且未超出
        for (uint i = 0; i < tmp_DisSecPledgeApply.appliedSecurities.length; i++) {
            //Notify(0, tmp_DisSecPledgeApply.appliedSecurities[i].toJson());
            //必须携带质押申请后携带的编号
            if (tmp_DisSecPledgeApply.appliedSecurities[i].id == 0) {
                Notify(errno_prefix+uint(DisSecPledgeApplyError.APPLIEDSECURITIES_ERROR),"APPLIEDSECURITIES id is 0");
                 return; 
            }
            //是否解除红利字段合法
            if (tmp_DisSecPledgeApply.appliedSecurities[i].isProfit != LibPledgeSecurity.IsProfit.T && tmp_DisSecPledgeApply.appliedSecurities[i].isProfit != LibPledgeSecurity.IsProfit.F) {
                Notify(errno_prefix+uint(DisSecPledgeApplyError.APPLIEDSECURITIES_ERROR),"APPLIEDSECURITIES isProfit is illegal");
                 return; 
            }
            if (!secPledgeManager.isMatchingAppliedSecurities(tmp_DisSecPledgeApply.secPledgeId, tmp_DisSecPledgeApply.appliedSecurities[i].toJson())) {
                 Notify(errno_prefix+uint(DisSecPledgeApplyError.APPLIEDSECURITIES_ERROR),"APPLIEDSECURITIES is not match");
                 return; 
            }
            if (tmp_DisSecPledgeApply.appliedSecurities[i].isProfit == LibPledgeSecurity.IsProfit.T ) {
                tmp_DisSecPledgeApply.appliedSecurities[i].profitAmount = 0;
            } else { //为F表示解除红利
                tmp_DisSecPledgeApply.appliedSecurities[i].profitAmount = 
                    secPledgeManager.getAppliedSecuritiesProfitAmount(tmp_DisSecPledgeApply.secPledgeId, tmp_DisSecPledgeApply.appliedSecurities[i].toJson());
            }
            
            tmp_DisSecPledgeApply.appliedSecurities[i].freezeNo = 
                    secPledgeManager.getAppliedSecuritiesFreezeNo(tmp_DisSecPledgeApply.secPledgeId, tmp_DisSecPledgeApply.appliedSecurities[i].toJson()).recoveryToString();
            tmp_DisSecPledgeApply.appliedSecurities[i].subFreezeNo = 
                    secPledgeManager.getAppliedSecuritiesSubFreezeNo(tmp_DisSecPledgeApply.secPledgeId, tmp_DisSecPledgeApply.appliedSecurities[i].toJson()).recoveryToString();
            


        }

        //校验司法冻结情况是否解除质押
        if (tmp_DisSecPledgeApply.isFreeze != uint(LibDisSecPledgeApply.IsFreeze.YES) 
            && tmp_DisSecPledgeApply.isFreeze != uint(LibDisSecPledgeApply.IsFreeze.NOT) 
        ) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.IS_FREEZE_INVALID),"IS_FREEZE is INVALID");
            return; 
        }

        //检查剩余融资金额是不是小于当前融资金额
        //作废，改成只要大等于0即可
        //if (!secPledgeManager.isMatchingFinancingAmountRemain(secPledgeId, tmp_DisSecPledgeApply.financingAmountRemain)) {
        //作废
        // if (tmp_DisSecPledgeApply.financingAmountRemain < 0) {
        //     Notify(errno_prefix+uint(DisSecPledgeApplyError.IS_FinancingAmountRemain),"FinancingAmountRemain is less than 0");
        //     return; 
        // }
        

        //检查解押类型
        //LibBiz.BizStatus _status;
        if (tmp_DisSecPledgeApply.applyType != LibDisSecPledgeApply.DispledgeType.ALL && tmp_DisSecPledgeApply.applyType != LibDisSecPledgeApply.DispledgeType.PARTIALLY)  {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.APPLY_TYPE_NOT_ALLOWED),"DisSecPledgeApply APPLY_TYPE_NOT_ALLOWED");
            return;             
        }

        //检查证明文件是否已经申请或者已经邮寄
        if(secPledgeManager.isEvidenceApplied(tmp_DisSecPledgeApply.secPledgeId,1)){
            Notify(errno_prefix+uint(DisSecPledgeApplyError.EVIDENCE_ALREADY_APPLIED),"DisSecPledgeApply EVIDENCE_ALREADY_APPLIED");
            return;
        }
        if(secPledgeManager.isEvidenceMailed(tmp_DisSecPledgeApply.secPledgeId,1)){
            Notify(errno_prefix+uint(DisSecPledgeApplyError.EVIDENCE_ALREADY_MAILED),"DisSecPledgeApply EVIDENCE_ALREADY_MAILED");
            return;
        }

        string memory businessNo;
        if(tmp_DisSecPledgeApply.id == 0) {
            //json中不带id，新业务
            businessNo = sq.genBusinessNo(code, "02").recoveryToString();
        } else {
            //重新发起的业务，保留原业务流水号
            LibLog.log("id exists", tmp_DisSecPledgeApply.id.toString());
            string memory bizJson = getBizById(tmp_DisSecPledgeApply.id);

            //业务流水号获取失败，业务已被重新发起
            if(bizJson.equals("")) {
                Notify(errno_prefix+uint(DisSecPledgeApplyError.BUSINESSNO_ERROR),"该业务已被重新发起");
                return;
            }
            LibJson.push(bizJson);
            businessNo = bizJson.jsonRead("businessNo");
            LibJson.pop();

            //将解除id从原质押记录中去掉
            od.undo_SecPledgeStatus_byDis(tmp_DisSecPledgeApply.id);

            //删除原业务
            uint _ret = od.delete_Biz_byId(tmp_DisSecPledgeApply.id);
            LibLog.log("delete_Biz_byId:", _ret.toString());
            _ret = od.delete_DisSecPledgeApply_byId(tmp_DisSecPledgeApply.id);
            LibLog.log("delete_DisSecPledgeApply_byId:", _ret.toString());

        }
        tmp_DisSecPledgeApply.businessNo = businessNo;
        // string memory businessNo = sq.genBusinessNo(code, "02").recoveryToString();

        tmp_DisSecPledgeApply.pledgorName = LibStack.popex(pm.findNameById(tmp_DisSecPledgeApply.pledgorId));
        tmp_DisSecPledgeApply.pledgeeName = LibStack.popex(pm.findNameById(tmp_DisSecPledgeApply.pledgeeId));

        //在本合约增加记录，再往Biz增加记录，并修改SecPledge中status为处理中及增加disSecPedgeApplyIds
        _id = sq.getSeqNo("Biz.id");
        tmp_DisSecPledgeApply.id = _id;
        tmp_DisSecPledgeApply.bizId = _id;
        // tmp_DisSecPledgeApply.businessNo = businessNo;    
        tmp_DisSecPledgeApply.applyTime = now*1000;


        t_biz.reset();
        t_biz.id = t_biz.relatedId = _id;
        t_biz.pledgorId = tmp_DisSecPledgeApply.pledgorId;
        t_biz.pledgeeId = tmp_DisSecPledgeApply.pledgeeId;
        t_biz.pledgorName = tmp_DisSecPledgeApply.pledgorName;
        t_biz.pledgeeName = tmp_DisSecPledgeApply.pledgeeName;
        t_biz.bizType = LibBiz.BizType.DISPLEDGE_BIZ;
        t_biz.businessNo = businessNo;
        t_biz.channelType = uint(LibBiz.ChannelType.ONLINE);
        t_biz.startTime = now*1000;
        t_biz.createTime = now*1000;
        t_biz.updateTime = now*1000;

        if (t_biz_secpledgeApply.fromJson(getBizById(tmp_DisSecPledgeApply.secPledgeId))) {
            for( i = t_biz_secpledgeApply.backAttachments.length; i >= 1; i--) {
                if((t_biz_secpledgeApply.backAttachments[i-1].fileType == 103 || t_biz_secpledgeApply.backAttachments[i-1].fileType == 104) &&
                   (t_biz_secpledgeApply.backAttachments[i-1].ext2.equals("jpg") || t_biz_secpledgeApply.backAttachments[i-1].ext2.equals("pdf"))
                ) {
                    t_biz.backAttachments.push(t_biz_secpledgeApply.backAttachments[i-1]);
                    if(t_biz.backAttachments.length == 2) {
                        break;
                    }
                }
            }
        } else {
            LibLog.log("biz not exists");
        }

        bizJson = t_biz.toJson();
        LibLog.log("insertBiz: ", bizJson);
        _ret = od.insert_Biz(bizJson);
        LibLog.log("insert_biz result: ", _ret);

        // bizManager.insert(_id, tmp_DisSecPledgeApply.pledgorId, tmp_DisSecPledgeApply.pledgeeId, address(0), uint(LibBiz.BizType.DISPLEDGE_BIZ), uint(LibBiz.BizStatus.NONE), _id);
        // bizManager.setBusinessNo(_id, businessNo);

        // tmp_DisSecPledgeApplyMap[_id] = tmp_DisSecPledgeApply;
        // tmp_DisSecPledgeApplyIds.push(_id);
        LibLog.log("insert_DisSecPledgeApply: ", tmp_DisSecPledgeApply.toJson());
        _ret = od.insert_DisSecPledgeApply(tmp_DisSecPledgeApply.toJson());
        LibLog.log("ret: ", _ret);

        secPledgeManager.addNewDisSecPledgeAplly(tmp_DisSecPledgeApply.secPledgeId, _id);

 
        string memory success = "";
        uint _total = 1;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString() );
        success = success.concat(",\"id\":",  _id.toString());
       // success = success.concat(",\"businessNo\":",  businessNo);
        success = success.concat(",\"businessNo\":\"",  businessNo,"\"");
        success = success.concat(",\"items\":[");
        success = success.concat("{\"id\":", _id.toString(),  "}]}}");
        Notify(0, success);

        return;
    }
    
    //6.1 新增解除质押申请业务
    function createDisPledgeApply(string _json) getSq getBizManager returns(bool _ret) {   
        // string memory businessNo = sq.genBusinessNo("A", "02").recoveryToString(); //手机客户端质押申请
        uint _id = createDisPledgeApplyCommon("A",_json);

        if(_id > 0) {
            //检查解押类型
            LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
            tmp_DisSecPledgeApply.fromJson(getDisSecPledgeApplyById(_id));   
            if (tmp_DisSecPledgeApply.applyType == LibDisSecPledgeApply.DispledgeType.ALL) {
                //全部解除进入待质权人人脸识别
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNAUTH;
            } else if (tmp_DisSecPledgeApply.applyType == LibDisSecPledgeApply.DispledgeType.PARTIALLY) {
                //部分解除进入待出质人确认的状态
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNCONFIRMED;
            } else {
                Notify(errno_prefix+uint(DisSecPledgeApplyError.APPLY_TYPE_NOT_ALLOWED),"DisSecPledgeApply APPLY_TYPE_NOT_ALLOWED");
                return; 
            }
            bizManager.changeStatus(_id, uint(_status));
            addAudit(_id, tmp_DisSecPledgeApply.pledgeeId, 1, "", "createDisPledgeApply", uint(LibBiz.BizStatus.NONE), uint(_status));
            return ; 
        }
    }  


    //6.2 根据质押登记ID查询解除质押详情
    function findById(uint id) constant getOrderDao returns (string) {
        uint len = od.select_DisSecPledgeApply_byId(id);
        return LibStack.popex(len);
    }

    //分页功能
    function pageByCond(string _json) constant getOrderDao returns (string) {
        uint len = od.pageByCond_DisSecPledgeApply(_json);
        return LibStack.popex(len);
    }

    //6.8.合同流转-质权人确认人脸确认
    function updateDisPledgeApplyByPledgeeFaceAuth(uint id, uint operateCode, string rejectReason) getSecPledgeManager getBizManager  {
        string memory _json = getDisSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER), "the dis apply id does not exist");
            return;
        }

        tmp_DisSecPledgeApply.fromJson(_json);

        uint bizId = tmp_DisSecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNAUTH)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                //根据解押是全部还是部分确定下一步的状态
                if (tmp_DisSecPledgeApply.applyType == LibDisSecPledgeApply.DispledgeType.ALL) {
                    //全部解除进入待审核装
                    _status = LibBiz.BizStatus.DISPLEDGE_UNAUDITED;
                } else if (tmp_DisSecPledgeApply.applyType == LibDisSecPledgeApply.DispledgeType.PARTIALLY) {
                    //部分解除确认后进入待审核状态
                    _status = LibBiz.BizStatus.DISPLEDGE_UNAUDITED;
                } else {
                    Notify(errno_prefix+uint(DisSecPledgeApplyError.APPLY_TYPE_NOT_ALLOWED),"DisSecPledgeApply APPLY_TYPE_NOT_ALLOWED");
                    return; 
                }
                bizManager.changeStatus(bizId, uint(_status));
                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //人脸识别失败，质权人人脸验证失败
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGEE_AUTH_FALIED;
                bizManager.endBiz(bizId, uint(_status));
                secPledgeManager.undoDisSecPledgeApply(tmp_DisSecPledgeApply.secPledgeId);
                
                Notify(0, 'success');
            }
            addAudit(bizId, tmp_DisSecPledgeApply.pledgeeId, operateCode, rejectReason, "updateDisPledgeApplyByPledgeeFaceAuth", _oldStatus, uint(_status));
            return;
        } 
        Notify(errno_prefix+uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED_SUCH_METHOD),"Only apply to DISPLEDGE_PLEDGEE_UNAUTH");
    }

    //6.3.合同流转 - 出质人确认（部分解除）
    function pledgorSubmit(uint id, uint operateCode, string rejectReason) getSecPledgeManager getBizManager {
        string memory _json = getDisSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER), "the dis apply id does not exist");
            return;
        }

        tmp_DisSecPledgeApply.fromJson(_json);

        uint bizId = tmp_DisSecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNCONFIRMED)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNAUTH;
                bizManager.changeStatus(bizId, uint(_status));

                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //出质人拒绝
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGOR_DENIED;
                bizManager.endBiz(bizId, uint(_status));
                secPledgeManager.undoDisSecPledgeApply(tmp_DisSecPledgeApply.secPledgeId);
                Notify(0, 'success');
            }
            addAudit(bizId, tmp_DisSecPledgeApply.pledgorId, operateCode, rejectReason, "pledgorSubmit", _oldStatus, uint(_status));
            return;
        }
        Notify(errno_prefix+uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED_SUCH_METHOD),"Only apply to DISPLEDGE_PLEDGOR_UNCONFIRMED");
    }

    //6.7.合同流转-出质人人脸确认
    function updateDisPledgeApplyByPledgorFaceAuth(uint id, uint operateCode, string rejectReason) getSecPledgeManager getBizManager {
        string memory _json = getDisSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER), "the dis apply id does not exist");
            return;
        }

        tmp_DisSecPledgeApply.fromJson(_json);

        uint bizId = tmp_DisSecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNAUTH)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNCONFIRMED;
                bizManager.changeStatus(bizId, uint(_status));
                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //人脸识别失败，出质人人脸验证失败
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGOR_AUTH_FALIED;
                bizManager.endBiz(bizId, uint(_status));
                secPledgeManager.undoDisSecPledgeApply(tmp_DisSecPledgeApply.secPledgeId);
                Notify(0, 'success');
            }
            addAudit(bizId, tmp_DisSecPledgeApply.pledgorId, operateCode, rejectReason, "updateDisPledgeApplyByPledgorFaceAuth", _oldStatus, uint(_status));
            return;
        }
        Notify(errno_prefix+uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED_SUCH_METHOD),"Only apply to DISPLEDGE_PLEDGOR_UNAUTH");
    }
    
    //6.4.合同流转 - 质权人确认（部分解除）
    function pledgeeSubmit(uint id, uint operateCode, string rejectReason) getSecPledgeManager getBizManager {
        string memory _json = getDisSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER), "the dis apply id does not exist");
            return;
        }

        tmp_DisSecPledgeApply.fromJson(_json);

        uint bizId = tmp_DisSecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNCONFIRMED)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNAUTH;
                bizManager.changeStatus(bizId, uint(_status));
                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //质权人拒绝
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGEE_DENIED;
                bizManager.endBiz(bizId, uint(_status));
                secPledgeManager.undoDisSecPledgeApply(tmp_DisSecPledgeApply.secPledgeId);
                Notify(0, 'success');
            }
            addAudit(bizId, tmp_DisSecPledgeApply.pledgeeId, operateCode, rejectReason, "pledgeeSubmit", _oldStatus, uint(_status));
            return;
        }
        Notify(errno_prefix+uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED_SUCH_METHOD),"Only apply to DISPLEDGE_PLEDGEE_UNCONFIRMED");
    
    }

    //6.5 合同流转-管理员审核通过/拒绝
    function updateDisPledgeApplyByAdmin(uint id, uint operateCode, address auditorId, string rejectReason) getSecPledgeManager getBizManager {
        string memory _json = getDisSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER), "the dis apply id does not exist");
            return;
        }

        tmp_DisSecPledgeApply.fromJson(_json);

        uint bizId = tmp_DisSecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.DISPLEDGE_UNAUDITED)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                _status = LibBiz.BizStatus.DISPLEDGE_AUTIT_SUCCESS;
                bizManager.changeStatus(bizId, uint(_status));
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //管理员拒绝
                _status = LibBiz.BizStatus.DISPLEDGE_AUDIT_FALIED;
                bizManager.endBiz(bizId, uint(_status));
                secPledgeManager.undoDisSecPledgeApply(tmp_DisSecPledgeApply.secPledgeId);

            }
            addAudit(bizId, auditorId, operateCode, rejectReason, "updateDisPledgeApplyByAdmin", _oldStatus, uint(_status));
            Notify(0, 'success');
            return;
        }
        Notify(errno_prefix+uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED_SUCH_METHOD),"Only apply to DISPLEDGE_UNAUDITED");

    }
     
    //6.10.（机构）新增解除质押申请业务
      function createDisPledgeApplyInstitution(string _json) getSq getBizManager returns(bool _ret) {
        // string memory businessNo = sq.genBusinessNo("W", "02").recoveryToString(); //WEB客户端质押申请
        uint _id = createDisPledgeApplyCommon("W", _json);

        if(_id > 0) {
            tmp_DisSecPledgeApply.fromJson(getDisSecPledgeApplyById(_id));
            uint _status;
            if(tmp_DisSecPledgeApply.applyType == LibDisSecPledgeApply.DispledgeType.ALL){
                _status = uint(LibBiz.BizStatus.DISPLEDGE_UNAUDITED);   //变为DISPLEDGE_UNAUDITED，待审核
            }else if(tmp_DisSecPledgeApply.applyType == LibDisSecPledgeApply.DispledgeType.PARTIALLY){
                _status = uint(LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNCONFIRMED); // 变为DISPLEDGE_PLEDGOR_UNCONFIRMED,待出质人确认
            }
            bizManager.changeStatus(_id, _status);
            addAudit(_id, tmp_DisSecPledgeApply.pledgeeId, 1, "", "createDisPledgeApplyInstitution", uint(LibBiz.BizStatus.NONE), _status);
            return;
        }
      }
     
    //6.11.合同流转 - （机构）出质人确认（部分解除）
    function pledgorSubmitInstitution(uint id,uint operateCode,string rejectReason) getSecPledgeManager getBizManager {
        string memory _json = getDisSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER), "the dis apply id does not exist");
            return;
        }

        tmp_DisSecPledgeApply.fromJson(_json);

        uint bizId=tmp_DisSecPledgeApply.bizId;
        LibBiz.BizStatus _status=LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if(_oldStatus == uint(LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNCONFIRMED)) {	  
       
            if(operateCode==uint(OperateCode.PASS)){
                bizManager.changeStatus(bizId,uint(LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNCONFIRMED));  //28-待质权人确认（部分解除）
                _status=LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNCONFIRMED;
                Notify(0,"success");
            } else if(operateCode==uint(OperateCode.REJECT)) {
                bizManager.endBiz(bizId,uint(LibBiz.BizStatus.DISPLEDGE_PLEDGOR_DENIED)); //出质人拒绝
                secPledgeManager.undoDisSecPledgeApply(tmp_DisSecPledgeApply.secPledgeId);
                _status=LibBiz.BizStatus.DISPLEDGE_PLEDGOR_DENIED;	
                Notify(0,"success");		   
            }
            addAudit(bizId, tmp_DisSecPledgeApply.pledgorId, operateCode, rejectReason, "pledgorSubmitInstitution", _oldStatus, uint(_status));
            return;
       }  
       Notify(errno_prefix+uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED_SUCH_METHOD),"Only apply to DISPLEDGE_PLEDGOR_UNCONFIRMED");
    } 
     
     //6.12.合同流转 - (机构）质权人确认（部分解除）
    function pledgeeSubmitInstitution(uint id,uint operateCode,string rejectReason) getSecPledgeManager getBizManager {
        string memory _json = getDisSecPledgeApplyById(id);
        if(_json.equals("")) {
            Notify(errno_prefix+uint(DisSecPledgeApplyError.BAD_PARAMETER), "the dis apply id does not exist");
            return;
        }

        tmp_DisSecPledgeApply.fromJson(_json);

        uint bizId = tmp_DisSecPledgeApply.bizId;
        LibBiz.BizStatus _status = LibBiz.BizStatus.NONE;
        uint _oldStatus = getBizStatus(bizId);
        if (_oldStatus == uint(LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNCONFIRMED)) {
            if (operateCode == uint(OperateCode.PASS) ) {
                _status = LibBiz.BizStatus.DISPLEDGE_UNAUDITED;  //待审核
                bizManager.changeStatus(bizId, uint(_status));
                Notify(0, 'success');
            } else if (operateCode == uint(OperateCode.REJECT)) {
                //质权人拒绝
                _status = LibBiz.BizStatus.DISPLEDGE_PLEDGEE_DENIED;  //出质人已拒绝（部分解除）
                bizManager.endBiz(bizId, uint(_status));
                secPledgeManager.undoDisSecPledgeApply(tmp_DisSecPledgeApply.secPledgeId);
                
                Notify(0, 'success');
            }
            addAudit(bizId, tmp_DisSecPledgeApply.pledgeeId, operateCode, rejectReason, "pledgeeSubmitInstitution", _oldStatus, uint(_status));
            return;
        }
        Notify(errno_prefix+uint(DisSecPledgeApplyError.STATUS_NOT_ALLOWED_SUCH_METHOD),"Only apply to DISPLEDGE_PLEDGEE_UNCONFIRMED");
    }

    /* for CsdcBaseInterface */
    function hasTodo(address _userId) constant returns (bool) {}
    function findNameById(address _id) constant returns (uint) {}    
}	
     

