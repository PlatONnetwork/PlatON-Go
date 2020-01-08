pragma solidity ^0.4.12;
/**
* @file LibBiz.sol
* @author yiyating
* @time 2016-12-27
* @desc 业务定义
*/


import "./LibAudit.sol";
import "./LibTradeUser.sol";
import "./LibTradeOperator.sol";
import "./LibAttachInfo.sol";

library LibBiz {
    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibAudit for *;
    using LibTradeUser for *;
    using LibTradeOperator for *;
    using LibBiz for *;
    using LibAttachInfo for *;

    /* 业务状态 */
    enum BizStatus {

        /* 以下为证券质押业务状态 */
        NONE,               //0-无
        PLEDGE_PLEDGEE_UNCONFIRMED,   //1-待质权人确认
        PLEDGE_PLEDGOR_UNPAID,      //2-待出质人支付
        PLEDGE_UNAUDITED,         //3-待审核
        PLEDGE_AUDIT_SUCCESS,       //4-审核通过
        PLEDGE_UNPROCESSED,       //5-已发送到核心系统
        PLEDGE_PLEDGEE_DENIED,      //6-质权人已拒绝
        PLEDGE_AUDIT_FAILED,      //7-投业部审核拒绝
        PLEDGE_COMPLETE,        //8-冻结成功
        PLEDGE_PROCESS_FAILED,      //9-冻结失败
        PLEDGE_TIMEOUT_FAILED,      //10-支付超时失败
        PLEDGE_PLEDGOR_UNAUTH,      //11-出质人待人脸认证
        PLEDGE_PLEDGOR_AUTH_FAILED,   //12-出质人人脸认证不通过
        PLEDGE_PLEDGEE_UNAUTH,      //13-质权人待人脸认证
        PLEDGE_PLEDGEE_AUTH_FAILED,   //14-质权人人脸认证不通过
        PLEDGE_PLEDGEE_UNPAID,      //15-待质权人支付
        PLEDGE_REFUND,          //16-已退款
        PLEDGE_PLEDGOR_UNSIGNED,  //17-待出质人签署
        PLEDGE_BIZ_TIMEOUT_FAIL,    //18-超时失败
        PLEDGE_COMPLETE_PARTIAL,    //19-部分质押申请冻结成功

        /* 以下为解除证券质押业务状态 */
        DISPLEDGE_PLEDGOR_UNCONFIRMED,  //20-待出质人确认
        DISPLEDGE_UNAUDITED,      //21-待审核
        DISPLEDGE_AUTIT_SUCCESS,    //22-审核通过
        DISPLEDGE_UNPROCESSED,      //23-已发送到核心系统，待处理
        DISPLEDGE_PLEDGOR_DENIED,   //24-出质人已拒绝（部分解除）
        DISPLEDGE_AUDIT_FALIED,     //25-审核不通过
        DISPLEDGE_COMPLETE_PARTIAL,   //26-部分解除成功，废除
        DISPLEDGE_COMPLETE_ALL,     //27-解除成功
        DISPLEDGE_PLEDGEE_UNCONFIRMED,  //28-待质权人确认（部分解除）
        DISPLEDGE_PLEDGEE_UNAUTH,   //29-质权人待人脸认证
        DISPLEDGE_PLEDGEE_AUTH_FALIED,  //30-质权人人脸认证已拒绝
        DISPLEDGE_PLEDGOR_UNAUTH,   //31-出质人待人脸认证
        DISPLEDGE_PLEDGOR_AUTH_FALIED,   //32-出质人人脸认证已拒绝
        DISPLEDGE_PLEDGEE_DENIED, //33 解除质押质权人拒绝（部分解除）
        DISPLEDGE_PROCESS_FAILED,    //34解除质押核心系统操作失败（解冻）

        NONE35,NONE36,NONE37,NONE38,NONE39,NONE40,
        NONE41,NONE42,NONE43,NONE44,NONE45,NONE46,NONE47,NONE48,


        /* 以下为中证登柜台录入质押业务状态 */
        CSDC_PLEDGE_PAYNOTIFY_UNCREATED_LATER,   //49-待生成付款通知（修改后）
        CSDC_PLEDGE_PAYMENT_CONFIRMED_LATER,     //50-已确认支付（修改后）
        CSDC_PLEDGE_STASHED,                //51-暂存中，作废
        CSDC_PLEDGE_CREATED,                //52-已创建
        CSDC_PLEDGE_PAYNOTIFY_UNCREATED,    //53-待生成付款通知,作废
        CSDC_PLEDGE_PAYNOTIFY_CREATED,      //54-已生成付款通知
        CSDC_PLEDGE_PAYMENT_CONFIRMED,      //55-已确认支付
        CSDC_PLEDGE_UNAUDITED,              //56-办理人已录入，待复核
        CSDC_PLEDGE_AUDIT_DENIED,           //57-复核驳回，需办理人重新录入
        CSDC_PLEDGE_AUDIT_SUCCESS,          //58-复核通过，待AS400处理
        CSDC_PLEDGE_UNPROCESS,              //59-已发送核心系统
        CSDC_PLEDGE_COMPLETE_ALL,           //60-已办结（全部）
        CSDC_PLEDGE_COMPLETE_PARTIAL,       //61-已办结（部分）
        CSDC_PLEDGE_FAILED,                 //62-已办结（失败）

        NONE63,NONE64,NONE65,NONE66,NONE67,NONE68,

        /* 以下为中证登柜台录入解质押业务状态 */
        CSDC_DISPLEDGE_STASHED,             //69-办理人已保存
        CSDC_DISPLEDGE_CREATED,             //70-办理人已录入
        CSDC_DISPLEDGE_REVIEW_DENIED,       //71-复核已驳回，需代理人从新录入
        CSDC_DISPLEDGE_REVIEW_SUCCESS,      //72-复核通过，待发送到核心系统
        CSDC_DISPLEDGE_UNPROCESS,           //73-已发送到核心系统，待处理     
        CSDC_DISPLEDGE_COMPLETE,            //74-已办结（成功）
        CSDC_DISPLEDGE_PROCESS_FAIL,        //75-已办结（失败）

        NONE76,NONE77,NONE78,NONE79,NONE80,
        NONE81,NONE82,NONE83,NONE84,NONE85,NONE86,NONE87,NONE88,NONE89,NONE90,
        NONE91,NONE92,NONE93,NONE94,NONE95,NONE96,NONE97,NONE98,NONE99,NONE100,

        /* 以下为券商代理质押业务状态 */
        BROKER_PLEDGE_STASHED,                //101-暂存中
        BROKER_PLEDGE_CREATED,                //102-已创建
        BROKER_PLEDGE_UNAUDITED,              //103-办理人已录入，待复核
        BROKER_PLEDGE_UNRETYPED,              //104-复核驳回，需办理人重新录入
        BROKER_PLEDGE_CSDC_UNAUDITED,         //105-券商复核通过，待投业部初审
        BROKER_PLEDGE_CSDC_UNREVIEWED,        //106-投业部初审通过，待投业部复审
        BROKER_PLEDGE_LEADER_UNREVIEWED,      //107-待领导审阅
        BROKER_PLEDGE_LEADER_SUCCESS,         //108-领导审阅通过，待发送到AS400
        BROKER_PLEDGE_UNPROCESS,              //109-已发送到AS400
        BROKER_PLEDGE_COMPLETE_ALL,           //110-已办结（全部）
        BROKER_PLEDGE_COMPLETE_PARTIAL,       //111-已办结（部分）
        BROKER_PLEDGE_FAILED,                 //112-已办结（失败）

        QUERY_BROKER_UNCHECKED,             //113-待异地代理点查看（仅用于查询）
        QUERY_CSDC_UNCHECKED,               //114-待投业部审核人员查看（仅用于查询）

        NONE115,NONE116,NONE117,NONE118,NONE119,NONE120,
        NONE121,NONE122,NONE123,NONE124,NONE125,NONE126,NONE127,NONE128,NONE129,NONE130,
        NONE131,NONE132,NONE133,NONE134,NONE135,NONE136,NONE137,NONE138,NONE139,NONE140,
        NONE141,NONE142,NONE143,NONE144,NONE145,NONE146,NONE147,NONE148,NONE149,NONE150,

        BROKER_DISPLEDGE_STASHED,                //151-暂存中
        BROKER_DISPLEDGE_CREATED,                //152-已创建
        BROKER_DISPLEDGE_UNAUDITED,              //153-办理人已录入，待复核
        BROKER_DISPLEDGE_UNRETYPED,              //154-复核驳回，需办理人重新录入
        BROKER_DISPLEDGE_CSDC_UNAUDITED,         //155-券商复核通过，待投业部初审
        BROKER_DISPLEDGE_CSDC_UNREVIEWED,        //156-投业部初审通过，待投业部复审
        BROKER_DISPLEDGE_LEADER_UNREVIEWED,      //157-待领导审阅
        BROKER_DISPLEDGE_LEADER_SUCCESS,         //158-领导审阅通过，待发送到AS400
        BROKER_DISPLEDGE_UNPROCESS,              //159-已发送到AS400
        BROKER_DISPLEDGE_COMPLETE_ALL,           //160-已办结（全部）
        BROKER_DISPLEDGE_COMPLETE_PARTIAL,       //161-已办结（部分）
        BROKER_DISPLEDGE_FAILED                  //162-已办结（失败）
    }

    enum CheckStatus {
        NONE,                       //不能看
        NEITHER_READ,               //都未看
        BROKER_READ_CSDC_UNREAD,    //券商已看，中证登未看
        BROKER_UNREAD_CSDC_READ,    //券商未看，中证登已看
        BOTH_READ                  //均已看
    }

    /* 业务类型 */
    enum BizType { 
        NONE,              //无
        PLEDGE_BIZ,        //证券质押登记业务
        DISPLEDGE_BIZ      //解除证券质押登记业务
    }

    /* 业务办理渠道 */
    enum ChannelType {
        NONE,
        ONLINE,     //1-在线
        BY_BROKER,  //2-券商办理
        BY_CSDC     //3-中证登柜台办理
    }

    enum RejectStatus {
        NO_REJECT,
        BROKER_REJECT,
        CSDC_REJECT
    }

    struct Biz{
        uint id;        //id
        
        address pledgorId;    //出质人地址
        string pledgorName;   //出质人名称
        address pledgeeId;    //质权人地址
        string pledgeeName;   //质权人姓名
        address managerId;    //经办人地址

        LibTradeUser.TradeUser[]        pledgors;       //出质人
        LibTradeUser.TradeUser          pledgee;        //质权人   
        LibTradeOperator.TradeOperator  tradeOperator;  //经办人
        LibTradeOperator.TradeOperator  csdcLeader;     //中证登领导
        
        uint startTime;     //业务发起时间（完成支付时间）
        uint endTime;     //业务结束时间

        BizStatus status;   //业务状态
        CheckStatus checkStatus;    //券商/柜面人员查看状态

        BizType bizType;    //业务类型
        LibAudit.Audit[] audits;//审核信息
        uint paymentId;     //缴费信息id
        uint relatedId;     //关联申请信息id
        
        string businessNo;   //业务流水号（同质押、解质押业务编号）
        string pledgeContractNo;  //中证登合同编号
        
        string pledgeContractFileId;  //质押合同文件编号
        string pledgeContractName;    //质押合同文件名称
    
        uint channelType;       //办理渠道类型
        string desc;            //业务描述

        uint rejectStatus;      //驳回状态
        LibAttachInfo.AttachInfo[]  backAttachments;    //后端操作普通附件

        string pledgeNotify5No;         //预留文件5
        string pledgeNotify5FileId;     //预留文件5
        string pledgeNotify5FileName;   //预留文件5

        uint createTime;  //创建时间
        uint updateTime;  //更新时间

        uint pledgeStatus;  //质物状态status

        string cxbh;    //查询编号
        string lczy;    //录入操作员
        string fhcz;    //复核操作员
    
    }

    /**
    *@desc fromJson for Biz
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function fromJson(Biz storage _self, string _json) internal returns(bool succ) {
        _self.reset();
        if(LibJson.push(_json) == 0) {
            return false;
        }

        if (!_json.isJson()) {
            LibJson.pop();
            return false;
        }

        _self.id = _json.jsonRead("id").toUint();
        _self.pledgorId = _json.jsonRead("pledgorId").toAddress();
        _self.pledgorName = _json.jsonRead("pledgorName");
        _self.pledgeeId = _json.jsonRead("pledgeeId").toAddress();
        _self.pledgeeName = _json.jsonRead("pledgeeName");
        _self.managerId = _json.jsonRead("managerId").toAddress();
        _self.pledgors.fromJsonArray(_json.jsonRead("pledgors"));
        _self.pledgee.fromJson(_json.jsonRead("pledgee"));
        _self.tradeOperator.fromJson(_json.jsonRead("tradeOperator"));
        _self.csdcLeader.fromJson(_json.jsonRead("csdcLeader"));
        _self.startTime = _json.jsonRead("startTime").toUint();
        _self.endTime = _json.jsonRead("endTime").toUint();
        _self.status = BizStatus(_json.jsonRead("status").toUint());
        _self.checkStatus = CheckStatus(_json.jsonRead("checkStatus").toUint());
        _self.bizType = BizType(_json.jsonRead("bizType").toUint());
        _self.audits.fromJsonArray(_json.jsonRead("audits"));
        _self.paymentId = _json.jsonRead("paymentId").toUint();
        _self.relatedId = _json.jsonRead("relatedId").toUint();
        _self.businessNo = _json.jsonRead("businessNo");
        _self.pledgeContractNo = _json.jsonRead("pledgeContractNo");
        _self.pledgeContractFileId = _json.jsonRead("pledgeContractFileId");
        _self.pledgeContractName = _json.jsonRead("pledgeContractName");
        _self.channelType = _json.jsonRead("channelType").toUint();
        _self.desc = _json.jsonRead("desc");
        _self.rejectStatus = _json.jsonRead("rejectStatus").toUint();
        _self.backAttachments.fromJsonArray(_json.jsonRead("backAttachments"));
        _self.pledgeNotify5No = _json.jsonRead("pledgeNotify5No");
        _self.pledgeNotify5FileId = _json.jsonRead("pledgeNotify5FileId");
        _self.pledgeNotify5FileName = _json.jsonRead("pledgeNotify5FileName");
        _self.createTime = _json.jsonRead("createTime").toUint();
        _self.updateTime = _json.jsonRead("updateTime").toUint();
        _self.pledgeStatus = _json.jsonRead("pledgeStatus").toUint();
        _self.cxbh = _json.jsonRead("cxbh");
        _self.lczy = _json.jsonRead("lczy");
        _self.fhcz = _json.jsonRead("fhcz");

        LibJson.pop();
        return true;
    }

    /**
    *@desc toJson for Biz
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function toJson(Biz storage _self) internal constant returns (string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("id", _self.id);
        len = LibStack.appendKeyValue("pledgorId", _self.pledgorId);
        len = LibStack.appendKeyValue("pledgorName", _self.pledgorName);
        len = LibStack.appendKeyValue("pledgeeId", _self.pledgeeId);
        len = LibStack.appendKeyValue("pledgeeName", _self.pledgeeName);
        len = LibStack.appendKeyValue("managerId", _self.managerId);
        len = LibStack.appendKeyValue("pledgors", _self.pledgors.toJsonArray());
        len = LibStack.appendKeyValue("pledgee", _self.pledgee.toJson());
        len = LibStack.appendKeyValue("tradeOperator", _self.tradeOperator.toJson());
        len = LibStack.appendKeyValue("csdcLeader", _self.csdcLeader.toJson());
        len = LibStack.appendKeyValue("startTime", _self.startTime);
        len = LibStack.appendKeyValue("endTime", _self.endTime);
        len = LibStack.appendKeyValue("status", uint(_self.status));
        len = LibStack.appendKeyValue("checkStatus", uint(_self.checkStatus));
        len = LibStack.appendKeyValue("bizType", uint(_self.bizType));
        len = LibStack.appendKeyValue("audits", _self.audits.toJsonArray());
        len = LibStack.appendKeyValue("paymentId", _self.paymentId);
        len = LibStack.appendKeyValue("relatedId", _self.relatedId);
        len = LibStack.appendKeyValue("businessNo", _self.businessNo);
        len = LibStack.appendKeyValue("pledgeContractNo", _self.pledgeContractNo);
        len = LibStack.appendKeyValue("pledgeContractFileId", _self.pledgeContractFileId);
        len = LibStack.appendKeyValue("pledgeContractName", _self.pledgeContractName);
        len = LibStack.appendKeyValue("channelType", _self.channelType);
        len = LibStack.appendKeyValue("desc", _self.desc);
        len = LibStack.appendKeyValue("rejectStatus", _self.rejectStatus);
        len = LibStack.appendKeyValue("backAttachments", _self.backAttachments.toJsonArray());
        len = LibStack.appendKeyValue("pledgeNotify5No", _self.pledgeNotify5No);
        len = LibStack.appendKeyValue("pledgeNotify5FileId", _self.pledgeNotify5FileId);
        len = LibStack.appendKeyValue("pledgeNotify5FileName", _self.pledgeNotify5FileName);
        len = LibStack.appendKeyValue("createTime", _self.createTime);
        len = LibStack.appendKeyValue("updateTime", _self.updateTime);
        len = LibStack.appendKeyValue("pledgeStatus", _self.pledgeStatus);
        len = LibStack.appendKeyValue("cxbh", _self.cxbh);
        len = LibStack.appendKeyValue("lczy", _self.lczy);
        len = LibStack.appendKeyValue("fhcz", _self.fhcz);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    /**
    *@desc update for Biz
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function update(Biz storage _self, string _json) internal returns(bool succ) {
        if(LibJson.push(_json) == 0) {
            return false;
        }

        if (!_json.isJson()) {
            LibJson.pop();
            return false;
        }

        if (_json.jsonKeyExists("id"))
            _self.id = _json.jsonRead("id").toUint();
        if (_json.jsonKeyExists("pledgorId"))
            _self.pledgorId = _json.jsonRead("pledgorId").toAddress();
        if (_json.jsonKeyExists("pledgorName"))
            _self.pledgorName = _json.jsonRead("pledgorName");
        if (_json.jsonKeyExists("pledgeeId"))
            _self.pledgeeId = _json.jsonRead("pledgeeId").toAddress();
        if (_json.jsonKeyExists("pledgeeName"))
            _self.pledgeeName = _json.jsonRead("pledgeeName");
        if (_json.jsonKeyExists("managerId"))
            _self.managerId = _json.jsonRead("managerId").toAddress();
        if (_json.jsonKeyExists("pledgors"))
            _self.pledgors.fromJsonArray(_json.jsonRead("pledgors"));
        if (_json.jsonKeyExists("pledgee"))
            _self.pledgee.fromJson(_json.jsonRead("pledgee"));
        if (_json.jsonKeyExists("tradeOperator"))
            _self.tradeOperator.fromJson(_json.jsonRead("tradeOperator"));
        if (_json.jsonKeyExists("csdcLeader"))
            _self.csdcLeader.fromJson(_json.jsonRead("csdcLeader"));
        if (_json.jsonKeyExists("startTime"))
            _self.startTime = _json.jsonRead("startTime").toUint();
        if (_json.jsonKeyExists("endTime"))
            _self.endTime = _json.jsonRead("endTime").toUint();
        if (_json.jsonKeyExists("status"))
            _self.status = BizStatus(_json.jsonRead("status").toUint());
        if (_json.jsonKeyExists("checkStatus"))
            _self.checkStatus = CheckStatus(_json.jsonRead("checkStatus").toUint());
        if (_json.jsonKeyExists("bizType"))
            _self.bizType = BizType(_json.jsonRead("bizType").toUint());
        if (_json.jsonKeyExists("audits"))
            _self.audits.fromJsonArray(_json.jsonRead("audits"));
        if (_json.jsonKeyExists("paymentId"))
            _self.paymentId = _json.jsonRead("paymentId").toUint();
        if (_json.jsonKeyExists("relatedId"))
            _self.relatedId = _json.jsonRead("relatedId").toUint();
        if (_json.jsonKeyExists("businessNo"))
            _self.businessNo = _json.jsonRead("businessNo");
        if (_json.jsonKeyExists("pledgeContractNo"))
            _self.pledgeContractNo = _json.jsonRead("pledgeContractNo");
        if (_json.jsonKeyExists("pledgeContractFileId"))
            _self.pledgeContractFileId = _json.jsonRead("pledgeContractFileId");
        if (_json.jsonKeyExists("pledgeContractName"))
            _self.pledgeContractName = _json.jsonRead("pledgeContractName");
        if (_json.jsonKeyExists("channelType"))
            _self.channelType = _json.jsonRead("channelType").toUint();
        if (_json.jsonKeyExists("desc"))
            _self.desc = _json.jsonRead("desc");
        if (_json.jsonKeyExists("rejectStatus"))
            _self.rejectStatus = _json.jsonRead("rejectStatus").toUint();
        if (_json.jsonKeyExists("backAttachments"))
            _self.backAttachments.fromJsonArray(_json.jsonRead("backAttachments"));
        if (_json.jsonKeyExists("pledgeNotify5No"))
            _self.pledgeNotify5No = _json.jsonRead("pledgeNotify5No");
        if (_json.jsonKeyExists("pledgeNotify5FileId"))
            _self.pledgeNotify5FileId = _json.jsonRead("pledgeNotify5FileId");
        if (_json.jsonKeyExists("pledgeNotify5FileName"))
            _self.pledgeNotify5FileName = _json.jsonRead("pledgeNotify5FileName");
        if (_json.jsonKeyExists("createTime"))
            _self.createTime = _json.jsonRead("createTime").toUint();
        if (_json.jsonKeyExists("updateTime"))
            _self.updateTime = _json.jsonRead("updateTime").toUint();
        if (_json.jsonKeyExists("pledgeStatus"))
            _self.pledgeStatus = _json.jsonRead("pledgeStatus").toUint();
        if (_json.jsonKeyExists("cxbh"))
            _self.cxbh = _json.jsonRead("cxbh");
        if (_json.jsonKeyExists("lczy"))
            _self.lczy = _json.jsonRead("lczy");
        if (_json.jsonKeyExists("fhcz"))
            _self.fhcz = _json.jsonRead("fhcz");

        LibJson.pop();
        return true;
    }

    /**
    *@desc reset for Biz
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function reset(Biz storage _self) internal {
        delete _self.id;
        delete _self.pledgorId;
        delete _self.pledgorName;
        delete _self.pledgeeId;
        delete _self.pledgeeName;
        delete _self.managerId;
        _self.pledgors.length = 0;
        _self.pledgee.reset();
        _self.tradeOperator.reset();
        _self.csdcLeader.reset();
        delete _self.startTime;
        delete _self.endTime;
        delete _self.status;
        delete _self.checkStatus;
        delete _self.bizType;
        _self.audits.length = 0;
        delete _self.paymentId;
        delete _self.relatedId;
        delete _self.businessNo;
        delete _self.pledgeContractNo;
        delete _self.pledgeContractFileId;
        delete _self.pledgeContractName;
        delete _self.channelType;
        delete _self.desc;
        delete _self.rejectStatus;
        _self.backAttachments.length = 0;
        delete _self.pledgeNotify5No;
        delete _self.pledgeNotify5FileId;
        delete _self.pledgeNotify5FileName;
        delete _self.createTime;
        delete _self.updateTime;
        delete _self.pledgeStatus;
        delete _self.cxbh;
        delete _self.lczy;
        delete _self.fhcz;
    }


    function create(Biz storage _self, address _pledgorId, string _pledgorName, address _pledgeeId, string _pledgeeName, address _managerId, uint _bizType, uint _status, uint _relatedId) internal {
        _self.reset();
        _self.pledgorId = _pledgorId;
        _self.pledgorName = _pledgorName;
        _self.pledgeeId = _pledgeeId;
        _self.pledgeeName = _pledgeeName;
        _self.managerId = _managerId;
        _self.bizType   = LibBiz.BizType(_bizType);
        _self.status    = LibBiz.BizStatus(_status);
        _self.relatedId = _relatedId;
        _self.startTime = now*1000;
    }
}