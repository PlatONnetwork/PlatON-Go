pragma solidity ^0.4.12;
/**
* file InvoiceManager.sol
* author Xiaofeng Liu
* time 2017-04-06
* desc the defination of Post Information
*/

import "./csdc_library/LibInvoice.sol";
import "./Sequence.sol";
import "./BizManager.sol";

contract InvoiceManager is OwnerNamed{
    using LibInt for *; 
    using LibString for *;
    using LibInvoice for *;
    using LibBiz for *;
    using LibLog for *;
    using LibStack for *;
    using LibJson for *;
    
    //错误码
    //1.无错误
    //2.json解析出错
    //3.信息不存在
    //4.參數错误
    //5.范围出错
    //NO_ERROR,
    //JSON_INVALID,
    //INFO_NOT_EXIST,
    //PARAM_ERROR
    //RANGE_ERROR
    enum Error{
        NO_ERROR,
        INVOICEID_ERROR,
        ID_NOT_EXIST,
        BAD_PARAMETER,
        BIZID_EMPTY,
        CUSTOMERTYPE_EMPTY,
        BIZID_NOTEXIST,
        INVOICECONTENT_EMPTY,
        INVOICEAMOUNT_EMPTY,
        BIZSTATUS_ERROE,
        INVOICESTATUS_ERROR,
        AUDITORID_EMPTY,
        USER_NOT_EXISTS,
        RECEIVER_EMPTY,
        MOBILE_EMPTY,
        RECEIVERUNIT_EMPTY,
        DETAILADDRESS_EMPTY,
        POSTCODE_EMPTY,
        DELIVERYWAY_EMPTY,
        COMPANYNO_EMPTY,
        INVOICETITLE_EMPTY,
        DELIVERYNO_EMPTY,
        OPERATION_DENIED,
        ZERO_INVOICEAMOUNT,
        BUSINESSNO_REPEAT
    }
    uint errno_prefix = 16500;
    
    // 业务审核结果
    enum OperateCode{
        NONE, 
        PASS,   //通过
        REJECT, //拒绝
        WAIT    //等待处理
    }
    
    
    uint[] IdsofUser; //零时存储个人发票id
    uint[] tempList;

    //修改数据结构
    mapping(address=>uint[]) user2InvoiceIds;
    mapping(uint => LibInvoice.Invoice) invoiceMap;
    uint[] invoiceIds;

    mapping(address => LibInvoice.Invoice) user2Invoice;    //发票信息
    
    LibInvoice.Invoice _invoiceTemp;
    Sequence sq;
    BizManager bizManager;

    function InvoiceManager(){
        register("CsdcModule", "0.0.1.0", "InvoiceManager", '0.0.1.0');
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
        bizManager = BizManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "BizManager", "0.0.1.0"));
    }

    function initInvoice(string _invoiceJson) public returns(uint _id){
        _invoiceTemp.reset();
        LibLog.log("initInvoice",_invoiceJson);
        if(!_invoiceTemp.fromJson(_invoiceJson)){
            LibLog.log("initInvoice","initInvoice: json invalid");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"initInvoice: json invalid");
            return;
        }
        if(_invoiceTemp.pledgeeName.equals("")){
            LibLog.log("initInvoice","initInvoice: pledgeeName can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"initInvoice: pledgeeName can not empty");
            return;
        }
        if(_invoiceTemp.pledgorName.equals("")){
            LibLog.log("initInvoice","initInvoice: pledgorName can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"initInvoice: pledgorName can not empty");
            return;
        }
        if(_invoiceTemp.businessNo.equals("")){
            LibLog.log("initInvoice","initInvoice: businessNo can not empty");
            Notify(errno_prefix+uint(Error.BIZID_NOTEXIST),"initInvoice: businessNo can not empty");
            return;
        }
        //业务流水号不能重复
        for(uint i = 0; i< invoiceIds.length;i++){
            if(invoiceMap[invoiceIds[i]].businessNo.equals(_invoiceTemp.businessNo)){
                LibLog.log("initInvoice","initInvoice: businessNo already exist");
                Notify(errno_prefix+uint(Error.BUSINESSNO_REPEAT),"initInvoice: businessNo already exist");
                return;
            }
        }
        if(_invoiceTemp.customerType==LibInvoice.CustomerType.NONE||uint(_invoiceTemp.customerType)>uint(LibInvoice.CustomerType.organization)){
            LibLog.log("initInvoice","initInvoice: customerType empty or incorrect");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"initInvoice: customerType empty or incorrect");
            return;
        }
        if(_invoiceTemp.invoiceType==LibInvoice.InvoiceType.NONE||uint(_invoiceTemp.invoiceType)>uint(LibInvoice.InvoiceType.Invoice_VAT)){
            LibLog.log("initInvoice","initInvoice: invoiceType can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"initInvoice: invoiceType can not empty");
            return;
        }
        if(_invoiceTemp.invoiceAmount == 0){
            LibLog.log("initInvoice","initInvoice: invoiceAmount=0 do not generate invoice");
            Notify(errno_prefix+uint(Error.ZERO_INVOICEAMOUNT),"initInvoice: invoiceAmount=0 do not generate invoice");
            return;
        }
        _invoiceTemp.id = sq.getSeqNo("Invoice.id");
        _invoiceTemp.initDate = now*1000;
        _invoiceTemp.deliveryWay = LibInvoice.DeliveryWay.sendPay;
        _invoiceTemp.companyNo = LibInvoice.CompanyNo.SHUNFENG;
        _invoiceTemp.status = LibInvoice.InvoiceStatus.Invoice_INIT;
        _invoiceTemp.handleChannel = LibInvoice.HandleChannel.ONLINE;
        user2InvoiceIds[_invoiceTemp.userId].push(_invoiceTemp.id);
        invoiceMap[_invoiceTemp.id] = _invoiceTemp;
        invoiceIds.push(_invoiceTemp.id);
        _id = _invoiceTemp.id;
        string memory success = "";
        uint _total = 1;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString(), ",\"items\":[");
        success = success.concat("{\"id\":", _id.toString(),  "}]}}");
        LibLog.log("initInvoice","success");
        Notify(0, success);
    }

    function applyInvoice(string _invoiceJson) {
        _invoiceTemp.reset();
        LibLog.log("applyInvoice",_invoiceJson);
        if(!_invoiceTemp.fromJson(_invoiceJson)){
            LibLog.log("applyInvoice","applyInvoice: json invalid");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: json invalid");
            return;
        }
        
        invoiceMap[_invoiceTemp.id].invoiceType = _invoiceTemp.invoiceType;

        if(invoiceMap[_invoiceTemp.id].id==0){
            LibLog.log("applyInvoice","applyInvoice: id not exist");
            Notify(errno_prefix+uint(Error.ID_NOT_EXIST),"applyInvoice: id not exist");
            return;
        }

        //公共判断
        if(_invoiceTemp.invoiceTitle.equals("")){
            LibLog.log("applyInvoice","applyInvoice: invoiceTitle can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: invoiceTitle can not empty");
            return;
        }
        if(_invoiceTemp.invoiceAmount==0){
            LibLog.log("applyInvoice","applyInvoice: invoiceAmount can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: invoiceAmount can not empty");
            return;
        }
        if(_invoiceTemp.companyNo == LibInvoice.CompanyNo.NONE||uint(_invoiceTemp.companyNo)>uint(LibInvoice.CompanyNo.SHUNFENG)){
            LibLog.log("applyInvoice","applyInvoice: companyNo empty or incorrect");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: companyNo empty or incorrect");
            return;
        }
        if(_invoiceTemp.detailAddress.equals("")){
            LibLog.log("applyInvoice","applyInvoice: detailAddress can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: detailAddress can not empty");
            return;
        }
        if(_invoiceTemp.postCode.equals("")){
            LibLog.log("applyInvoice","applyInvoice: postCode can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: postCode can not empty");
            return;
        }
        if(_invoiceTemp.deliveryWay==LibInvoice.DeliveryWay.NONE||uint(_invoiceTemp.deliveryWay)>uint(LibInvoice.DeliveryWay.receivePay)){
            LibLog.log("applyInvoice","applyInvoice: deliveryWay empty or incorrect");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: deliveryWay empty or incorrect");
            return;
        }
        if(_invoiceTemp.receiver.equals("")){
            LibLog.log("applyInvoice","applyInvoice: receiver can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: receiver can not empty");
            return;
        }
        if(_invoiceTemp.mobile.equals("")){
            LibLog.log("applyInvoice","applyInvoice: mobile can not empty");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: mobile can not empty");
            return;
        }
        // if(_invoiceTemp.receiverUnit.equals("")){
        //     Notify(errno_prefix+uint(Error.BAD_PARAMETER),"applyInvoice: receiverUnit can not empty");
        //     return;
        // }
        if(invoiceMap[_invoiceTemp.id].status!=LibInvoice.InvoiceStatus.Invoice_INIT){
            LibLog.log("applyInvoice","invoice status error");
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"invoice status error");
            return;
        }

        if(invoiceMap[_invoiceTemp.id].invoiceType == LibInvoice.InvoiceType.Invoice_NORMAL){
            
            invoiceMap[_invoiceTemp.id].invoiceTitle = _invoiceTemp.invoiceTitle;
            invoiceMap[_invoiceTemp.id].invoiceDate = now*1000;
            invoiceMap[_invoiceTemp.id].invoiceAmount = _invoiceTemp.invoiceAmount;
            invoiceMap[_invoiceTemp.id].companyNo = _invoiceTemp.companyNo;
            invoiceMap[_invoiceTemp.id].receiver = _invoiceTemp.receiver;
            invoiceMap[_invoiceTemp.id].deliveryWay = _invoiceTemp.deliveryWay;
            invoiceMap[_invoiceTemp.id].mobile = _invoiceTemp.mobile;
            invoiceMap[_invoiceTemp.id].receiverUnit = _invoiceTemp.receiverUnit;
            invoiceMap[_invoiceTemp.id].detailAddress = _invoiceTemp.detailAddress;
            invoiceMap[_invoiceTemp.id].postCode = _invoiceTemp.postCode;
            invoiceMap[_invoiceTemp.id].email = _invoiceTemp.email;
            invoiceMap[_invoiceTemp.id].status = LibInvoice.InvoiceStatus.WAIT_MAIL;
            LibLog.log("applyInvoice","success");
            Notify(0,"success");
            return;
        }else{
            
            invoiceMap[_invoiceTemp.id].invoiceTitle = _invoiceTemp.invoiceTitle;
            invoiceMap[_invoiceTemp.id].taxpayerType = _invoiceTemp.taxpayerType;
            invoiceMap[_invoiceTemp.id].invoiceDate = now*1000;
            invoiceMap[_invoiceTemp.id].invoiceAmount = _invoiceTemp.invoiceAmount;
            invoiceMap[_invoiceTemp.id].companyNo = _invoiceTemp.companyNo;
            invoiceMap[_invoiceTemp.id].receiver = _invoiceTemp.receiver;
            invoiceMap[_invoiceTemp.id].deliveryWay = _invoiceTemp.deliveryWay;
            invoiceMap[_invoiceTemp.id].mobile = _invoiceTemp.mobile;
            invoiceMap[_invoiceTemp.id].phone = _invoiceTemp.phone;
            invoiceMap[_invoiceTemp.id].receiverUnit = _invoiceTemp.receiverUnit;
            invoiceMap[_invoiceTemp.id].detailAddress = _invoiceTemp.detailAddress;
            invoiceMap[_invoiceTemp.id].postCode = _invoiceTemp.postCode;
            
            invoiceMap[_invoiceTemp.id].vatCustomerType = _invoiceTemp.vatCustomerType;
            invoiceMap[_invoiceTemp.id].taxpayerIdentifyNo = _invoiceTemp.taxpayerIdentifyNo;
            invoiceMap[_invoiceTemp.id].customerName = _invoiceTemp.customerName;
            invoiceMap[_invoiceTemp.id].certCode = _invoiceTemp.certCode;
            invoiceMap[_invoiceTemp.id].depositBank = _invoiceTemp.depositBank;
            invoiceMap[_invoiceTemp.id].bankAccount = _invoiceTemp.bankAccount;
            invoiceMap[_invoiceTemp.id].phoneNumber = _invoiceTemp.phoneNumber;
            invoiceMap[_invoiceTemp.id].companyAddress = _invoiceTemp.companyAddress;
            invoiceMap[_invoiceTemp.id].contactIDcard = _invoiceTemp.contactIDcard;
            invoiceMap[_invoiceTemp.id].email = _invoiceTemp.email;

            invoiceMap[_invoiceTemp.id].bizLicenseFileId = _invoiceTemp.bizLicenseFileId;
            invoiceMap[_invoiceTemp.id].bizLicenseFileName = _invoiceTemp.bizLicenseFileName;
            invoiceMap[_invoiceTemp.id].contactIDcardFileId = _invoiceTemp.contactIDcardFileId;
            invoiceMap[_invoiceTemp.id].contactIDcardFileName = _invoiceTemp.taxRegistFileName;
            invoiceMap[_invoiceTemp.id].taxRegistFileId = _invoiceTemp.taxRegistFileId;
            invoiceMap[_invoiceTemp.id].taxRegistFileName = _invoiceTemp.taxRegistFileName;
            invoiceMap[_invoiceTemp.id].depositOpenFileId = _invoiceTemp.depositOpenFileId;
            invoiceMap[_invoiceTemp.id].depositOpenFileName = _invoiceTemp.depositOpenFileName;
            invoiceMap[_invoiceTemp.id].qualificationFileId = _invoiceTemp.qualificationFileId;
            invoiceMap[_invoiceTemp.id].qualificationFileName = _invoiceTemp.qualificationFileName;


            invoiceMap[_invoiceTemp.id].status = LibInvoice.InvoiceStatus.WAIT_INIT_AUDIT;
            LibLog.log("applyInvoice", "success");
            Notify(0,"success");
        }
    }

    function addDeliveryNo(uint _id, string _deliveryNo) {
        if(invoiceMap[_id].id==0){
            Notify(errno_prefix+uint(Error.ID_NOT_EXIST),"addDeliveryNo: id not exist");
            return;
        }
        if(invoiceMap[_id].status!=LibInvoice.InvoiceStatus.WAIT_MAIL){
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"addDeliveryNo: status must be applied");
            return;
        }
        
        if(_deliveryNo.equals("")){
            Notify(errno_prefix+uint(Error.DELIVERYNO_EMPTY),"addDeliveryNo: deliveryNo can not empty");
            return;
        }
        invoiceMap[_id].deliveryNo = _deliveryNo;
        invoiceMap[_id].status = LibInvoice.InvoiceStatus.MAILED;
        Notify(0,"success");
    }

    // 根据businessNo查询发票id
    function findIdByBusinessNo(string _businessNo) public returns(uint _id){
        for(uint i = 0; i< invoiceIds.length;i++){
            if(invoiceMap[invoiceIds[i]].businessNo.equals(_businessNo)){
                _id = invoiceMap[invoiceIds[i]].id;
                return;
            }
        }
    }
    // 根据businessNo查询发票
    function findByBusinessNo(string _businessNo) constant public returns(string _ret){
        uint _id = findIdByBusinessNo(_businessNo);
        _ret = findById(_id);
    }

    // 根据businessNo添加运单号
    function addDeliveryNoByBusinessNo(string _businessNo, string _deliveryNo) {
        uint _id = findIdByBusinessNo(_businessNo);
        addDeliveryNo(_id,_deliveryNo);
    }
    
    function findById(uint _id) constant public returns(string _ret) {
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": []}}";
        if(invoiceMap[_id].id==0){
            return;
        }
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 1, \"items\": [";
        _ret = _ret.concat(invoiceMap[_id].toJson());
        _ret = _ret.concat("]}}");
    }
    
    LibInvoice.Condition _cond;

    //个人用户根据状态分页查询 日期 业务单号 //发票类型 发票抬头 联系人
    function pageByUserId(string _json) constant public returns(string _ret) {
        _cond.fromJson(_json);

        if (user2InvoiceIds[_cond.userId].length <= 0) {
            return LibStack.popex(itemStackPush("", 0));
        }
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= user2InvoiceIds[_cond.userId].length) {
            return LibStack.popex(itemStackPush("", 0));
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for(uint i = 0; i< user2InvoiceIds[_cond.userId].length; i++){
            _invoiceTemp = invoiceMap[user2InvoiceIds[_cond.userId][i]];

            //选择不包含并且该发票的状态为init;  
            if (_cond.containInit == 1 && _invoiceTemp.status == LibInvoice.InvoiceStatus.Invoice_INIT) {
                continue;
            }
            //选择不包含待初审并且该发票的状态为待初审;
            if (_cond.containWaitInitAudit == 1 && _invoiceTemp.status == LibInvoice.InvoiceStatus.WAIT_INIT_AUDIT) {
                continue;
            }
            //选择不包含待复核并且该发票的状态为待复核;
            if (_cond.containWaitFinalAudit == 1 && _invoiceTemp.status == LibInvoice.InvoiceStatus.WAIT_FINAL_AUDIT) {
                continue;
            }     

            if (_json.jsonKeyExists("initDate") && _cond.initDate!=0 && _cond.initDate!=_invoiceTemp.initDate) {
                continue;
            }
            if (_json.jsonKeyExists("minInitDate") && _json.jsonKeyExists("maxInitDate") && _cond.maxInitDate != 0 
                && ( _invoiceTemp.initDate < _cond.minInitDate || _invoiceTemp.initDate > _cond.maxInitDate)) {
                continue;
            }
            if (_json.jsonKeyExists("minInvoiceDate") && _json.jsonKeyExists("maxInvoiceDate") && _cond.maxInvoiceDate != 0 
                && ( _invoiceTemp.invoiceDate < _cond.minInvoiceDate || _invoiceTemp.invoiceDate > _cond.maxInvoiceDate)) {
                continue;
            }
            if (_json.jsonKeyExists("invoiceDate") && _cond.invoiceDate!=0 && _cond.invoiceDate!=_invoiceTemp.invoiceDate) {
                continue;
            }
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(_invoiceTemp.businessNo)) {
                continue;
            }  
            if (_json.jsonKeyExists("invoiceType") && _cond.invoiceType != LibInvoice.InvoiceType.NONE && uint(_cond.invoiceType) != uint(_invoiceTemp.invoiceType)) {
                continue;
            }    
            if (_json.jsonKeyExists("status") && _cond.status != LibInvoice.InvoiceStatus.NONE && uint(_cond.status) != uint(_invoiceTemp.status)) {
                continue;
            }   
            if (_json.jsonKeyExists("invoiceTitle") && !_cond.invoiceTitle.equals("") && !_cond.invoiceTitle.equals(_invoiceTemp.invoiceTitle)) {
                continue;
            }  
            if (_json.jsonKeyExists("receiver") && !_cond.receiver.equals("") && !_cond.receiver.equals(_invoiceTemp.receiver)) {
                continue;
            }
            if (_json.jsonKeyExists("statuses") && _cond.statuses.length != 0 && !uint(_invoiceTemp.status).inArray(_cond.statuses)) {
                continue;
            }

            if (_count++ < _startIndex) {
                continue;
            }
            
            if (_total < _cond.pageSize) {
              if (_total > 0) {
                len = LibStack.append(",");
              }
              _total ++;
              len = LibStack.append(_invoiceTemp.toJson());
            }
        }
        LibJson.pop();
        return LibStack.popex(itemStackPush(LibStack.popex(len), _count));
    }

    //中证登用户查询发票 1 业务单号 发票类型 发票抬头 收件人姓名
    //--------------------------------------------------------------todo查询增值税待邮寄发票-------------------------------------
    function pageByCond(string _json) constant public returns(string _ret){
        if (invoiceIds.length <= 0) {
            return LibStack.popex(itemStackPush("", 0));
        }        
        _cond.fromJson(_json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= invoiceIds.length) {
            return LibStack.popex(itemStackPush("", 0));
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for(uint i = 0; i< invoiceIds.length; i++){
            _invoiceTemp = invoiceMap[invoiceIds[i]];
            //选择不包含并且该发票的状态为init;  
            if (_cond.containInit == 1 && _invoiceTemp.status == LibInvoice.InvoiceStatus.Invoice_INIT) {
                continue;
            }
            //选择不包含待初审并且该发票的状态为待初审;
            if (_cond.containWaitInitAudit == 1 && _invoiceTemp.status == LibInvoice.InvoiceStatus.WAIT_INIT_AUDIT) {
                continue;
            }
            //选择不包含待复核并且该发票的状态为待复核;
            if (_cond.containWaitFinalAudit == 1 && _invoiceTemp.status == LibInvoice.InvoiceStatus.WAIT_FINAL_AUDIT) {
                continue;
            }     

            if (_json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != _invoiceTemp.userId) {
                continue;
            }

            if (_json.jsonKeyExists("initDate") && _cond.initDate!=0 && _cond.initDate!=_invoiceTemp.initDate) {
                continue;
            }
            if (_json.jsonKeyExists("minInitDate") && _json.jsonKeyExists("maxInitDate") && _cond.maxInitDate != 0 
                && ( _invoiceTemp.initDate < _cond.minInitDate || _invoiceTemp.initDate > _cond.maxInitDate)) {
                continue;
            }
            if (_json.jsonKeyExists("minInvoiceDate") && _json.jsonKeyExists("maxInvoiceDate") && _cond.maxInvoiceDate != 0 
                && ( _invoiceTemp.invoiceDate < _cond.minInvoiceDate || _invoiceTemp.invoiceDate > _cond.maxInvoiceDate)) {
                continue;
            }
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(_invoiceTemp.businessNo)) {
                continue;
            }  
            if (_json.jsonKeyExists("invoiceType") && _cond.invoiceType != LibInvoice.InvoiceType.NONE && uint(_cond.invoiceType) != uint(_invoiceTemp.invoiceType)) {
                continue;
            }    
            if (_json.jsonKeyExists("status") && _cond.status != LibInvoice.InvoiceStatus.NONE && uint(_cond.status) != uint(_invoiceTemp.status)) {
                continue;
            }   
            if (_json.jsonKeyExists("invoiceTitle") && !_cond.invoiceTitle.equals("") && !_cond.invoiceTitle.equals(_invoiceTemp.invoiceTitle)) {
                continue;
            }  
            if (_json.jsonKeyExists("receiver") && !_cond.receiver.equals("") && !_cond.receiver.equals(_invoiceTemp.receiver)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length != 0 && !uint(_invoiceTemp.status).inArray(_cond.statuses)) {
                continue;
            }

            if (_count++ < _startIndex) {
                continue;
            }
            
            if (_total < _cond.pageSize) {
                if (_total > 0) {
                    len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_invoiceTemp.toJson());
            }
        }
        LibJson.pop();
        return LibStack.popex(itemStackPush(LibStack.popex(len), _count));
    
    }

    function listAll() constant public returns(string _ret){
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": [] }}";
        if(invoiceIds.length<=0){
            return;
        }
        _ret = "{\"ret\":0, \"message\": \"success\", \"data\":{";
        _ret = _ret.concat(uint(invoiceIds.length).toKeyValue("total"), ",\"items\":[");
        for(uint i = 0; i < invoiceIds.length; i++){    
            if(i==invoiceIds.length-1){
                _ret = _ret.concat(invoiceMap[invoiceIds[i]].toJson());
            }else{
                _ret = _ret.concat(invoiceMap[invoiceIds[i]].toJson(),",");
            }
        }
        _ret = _ret.concat("]}}");
    }

    function initialAudit(uint _id, uint _opcode,address _auditorId,string _rejectReason) {
        if (invoiceMap[_id].id == 0) {
            Notify(errno_prefix+uint(Error.BAD_PARAMETER), "InitialAudit: id does not exist");
            return;
        }
        //普通发票不用审核
        if(invoiceMap[_id].invoiceType == LibInvoice.InvoiceType.Invoice_NORMAL){
            Notify(errno_prefix+uint(Error.BAD_PARAMETER), "initialAudit:normal invoice");
            return;
        }
        //待初审的发票才能初审 todo;
        if(invoiceMap[_id].status != LibInvoice.InvoiceStatus.WAIT_INIT_AUDIT){
            Notify(errno_prefix+uint(Error.BAD_PARAMETER), "initialAudit:invoice status is not applied");
            return;
        }
        //添加初审意见
        invoiceMap[_id].initAuditOpinion = _rejectReason;
        //初审通过
        if(_opcode == uint(OperateCode.PASS)){
            //变更为待复核
            invoiceMap[_id].status = LibInvoice.InvoiceStatus.WAIT_FINAL_AUDIT;
            Notify(0,"success");
            return;
        }
        //初审未通过
        else{
            //变更为初始化
            invoiceMap[_id].status = LibInvoice.InvoiceStatus.Invoice_INIT;
            invoiceMap[_id].initAuditOpinion = "";
            invoiceMap[_id].initAuditClaimer = address(0);
            Notify(0,"success");
        }
    }

    function finalAudit(uint _id,uint _opcode,address _auditorId,string _rejectReason){
        if(invoiceMap[_id].id == 0) {
            Notify(errno_prefix+uint(Error.BAD_PARAMETER), "finalAudit: id does not exist");
            return;
        }
        //普通发票不用审核 todo;
        if(invoiceMap[_id].invoiceType == LibInvoice.InvoiceType.Invoice_NORMAL){
            Notify(errno_prefix+uint(Error.BAD_PARAMETER), "finalAudit:normal invoice");
            return;
        }
        //待复核的发票才能复核
        if(invoiceMap[_id].status != LibInvoice.InvoiceStatus.WAIT_FINAL_AUDIT){
            Notify(errno_prefix+uint(Error.BAD_PARAMETER), "finalAudit:invoice status is not applied");
            return;
        }
        //添加复核意见
        invoiceMap[_id].finalAuditOpinion = _rejectReason;
        //复核通过
        if(_opcode == uint(OperateCode.PASS)){
            invoiceMap[_id].status = LibInvoice.InvoiceStatus.WAIT_MAIL;
            Notify(0,"success");
            return;
        }
        //复核未通过
        else{
            invoiceMap[_id].status = LibInvoice.InvoiceStatus.Invoice_INIT;
            invoiceMap[_id].initAuditClaimer = address(0);
            invoiceMap[_id].initAuditOpinion = "";
            invoiceMap[_id].finalAuditClaimer = address(0);
            invoiceMap[_id].finalAuditOpinion = "";
            Notify(0,"success");
        }
        
    }

    function updateUserInvoice(string _json) {
        if(!_invoiceTemp.fromJson(_json)) {
            Notify(errno_prefix + uint(Error.BAD_PARAMETER), "json invalid.");
            return;
        }

        if(_invoiceTemp.userId == address(0)) {
            Notify(errno_prefix + uint(Error.BAD_PARAMETER), "userId empty.");
            return;  
        }

        user2Invoice[_invoiceTemp.userId].update(_json);
        Notify(0, "success");
    }

    function findInvoiceByUserId(address _userId) constant returns(string _ret) {
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": []}}";
        if(user2Invoice[_userId].userId == address(0)){
            return;
        }
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 1, \"items\": [";
        _ret = _ret.concat(user2Invoice[_userId].toJson());
        _ret = _ret.concat("]}}");
    }

    function deleteById(uint id) {
        delete tempList;
        for (uint i = 0; i < invoiceIds.length; i++)  {
            if(id == invoiceIds[i]) {
                continue;
            }
            else {
                tempList.push(invoiceIds[i]);
            }
        }
        delete invoiceIds;
        for (uint j = 0; j < tempList.length; ++j) {
            invoiceIds.push(tempList[j]);
        }
        delete invoiceMap[id];
    }

    //items入栈
    function itemStackPush(string _items, uint _total) constant internal returns (uint len){
        len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("ret", uint(0));
        len = LibStack.appendKeyValue("message", "success");
        len = LibStack.append(",");
        len = LibStack.append("\"data\":{");
        len = LibStack.appendKeyValue("total", _total);
        len = LibStack.append(",");
        len = LibStack.append("\"items\":[");
        len = LibStack.append(_items);
        len = LibStack.append("]}}");
        return len;
    }

    event Notify(uint _errorno, string _info);

}