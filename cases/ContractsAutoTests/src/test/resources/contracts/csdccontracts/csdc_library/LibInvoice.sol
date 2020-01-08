pragma solidity ^0.4.12;
/**
* file LibInvoice.sol
* author Xiaofeng Liu
* time 2017-03-20
* desc the defination of Post Information
*/
import "../utillib/LibInt.sol";
import "../utillib/LibString.sol";
import "../utillib/LibJson.sol";
import "../utillib/LibStack.sol";

library LibInvoice{
    using LibInt for *;
    using LibString for *;
    using LibJson for *;
    using LibInvoice for *;

    //发票状态
    // enum InvoiceStatus {
    //     NONE,
    //     Invoice_INIT,     //未申请
    //     Invoice_APPLIED,  //已申请
    //     MAILED,            //已邮寄
    //     WAIT_MAIL         //待邮寄
    // }
    enum InvoiceStatus {
        NONE,
        Invoice_INIT,     //待申请
        WAIT_MAIL,        //待邮寄
        MAILED,           //已邮寄
        WAIT_INIT_AUDIT, //待初审
        WAIT_FINAL_AUDIT  //待复核
    }
    //发票类型
    enum InvoiceType {
        NONE,
        Invoice_NORMAL,   //普通发票
        Invoice_VAT       //增值税专用发票
    }
    //客户类型
    enum CustomerType {
        NONE,
        personal,         //个人
        organization      //机构
    }
    //发票领取方式
    enum ReceiveWay {
        NONE,
        COUNTER,          //柜台
        MAIL              //邮寄
    }
    //发票邮寄方式
    enum DeliveryWay{
        NONE,
        sendPay,          //到付
        receivePay        //寄付
    }

    //快递公司
    enum CompanyNo{
        NONE,
        SHUNFENG
    }

    //业务类型
    enum BizType{
        NONE,
        PLEDGE_BIZ       //质押登记业务
    }

    // 客户类型
    enum VatCustomerType{
        NONE,
        enterprise,         //企业
        PUBLICINSTITUTION,  //事业单位
        GOVERNMENTAGENCIES, //政府机关
        PERSONAL,           //个人
        INDIVIDUALBUSINESS  //个体工商户
    }

    //审核结果
    enum AuditResult{
        NONE,
        PASS,   //通过
        REJECT, //拒绝
        WAIT    //等待处理
    }
    //业务渠道
    enum HandleChannel{
        NONE,
        ONLINE //在线办理
    }

    // 纳税人类型            1:小规模纳税人 2:一般纳税人
    enum TaxpayerType { NONE,SMALLTAXPAYER,GENERALTAXPAYER }

    struct Invoice{
        //普通发票信息
        uint id;                    // id
        address userId;             //用户id 
        CustomerType customerType;  //客户类型
        VatCustomerType vatCustomerType;  //客户类型
        uint bizId;                 //业务id 审核用
        string pledgeeName;         //质权人全称 机构
        string pledgorName;         //出质人全称 机构
        string businessNo;          //业务ID
        string pledgeContractNo;    //质押合同编号
        string invoiceTitle;        //发票抬头
        InvoiceType invoiceType;    //发票类型
        BizType bizType;            //业务类型
        uint invoiceAmount;         //发票金额
        InvoiceStatus status;       //发票状态
        HandleChannel handleChannel;//业务渠道
        uint initDate;              //初始化日期
        uint invoiceDate;           //开票日期
        //增值税专用发票信息
        TaxpayerType taxpayerType;    // 纳税人类型
        string taxpayerIdentifyNo;    // 纳税人识别号
        string customerName;          // 客户名称
        string certCode;              // 税务登记证号码或统一社会信用代码
        string depositBank;           // 开户行
        string bankAccount;           // 开户行账号
        string phoneNumber;           // 电话
        string companyAddress;        // 地址
        string contactIDcard;         // 专票联系人（领取人）身份证号码
        string email;                 // 电子邮箱
        //共同属性
        string receiver;            //收件人 专票联系人（领取人）姓名
        string mobile;              //收件人手机 专票联系人（领取人）手机号码
        string phone;               //收件人联系电话 专票联系人（领取人）联系电话
        DeliveryWay deliveryWay;    //发票邮寄方式 发票领取方式
        CompanyNo companyNo;        //快递公司名称
        string deliveryNo;          //物流单号
        string postCode;            //邮政编码 发票寄送地邮编
        string detailAddress;       //详细地址（所在地区+详细地址拼接）
        string receiverUnit;        //收件人单位 专票联系人（领取人）单位名称
        // 附件
        string bizLicenseFileId;        // 营业执照副本id
        string bizLicenseFileName;      // 营业执照副本name
        string contactIDcardFileId;     // 专票联系人（领取人）身份证件id
        string contactIDcardFileName;   // 专票联系人（领取人）身份证件name
        string taxRegistFileId;         // 税务登记证副本id
        string taxRegistFileName;       // 税务登记证副本name
        string depositOpenFileId;       // 基本存款账户开户许可证id
        string depositOpenFileName;     // 基本存款账户开户许可证name
        string qualificationFileId;     // 增值税一般纳税人资格证明
        string qualificationFileName;   // 增值税一般纳税人资格证明
        // 审核
        address initAuditClaimer;            // 初审领取人员
        address finalAuditClaimer;           // 复审领取人员
        string initAuditOpinion;             // 初审意见
        string finalAuditOpinion;            // 复核意见
    }
    struct Condition{
        address userId;
        string businessNo;
        uint initDate;
        uint minInitDate;
        uint maxInitDate;
        uint minInvoiceDate;
        uint maxInvoiceDate;
        uint invoiceDate;
        string invoiceTitle;
        string receiver;
        uint containInit; //0 包含 1 不包含
        uint containWaitInitAudit; // 0 包含 1 不包含
        uint containWaitFinalAudit;   // 0 包含 1 不包含
        InvoiceType invoiceType;    //发票类型
        InvoiceStatus status;       //发票状态
        uint[] statuses;

        uint pageSize;
        uint pageNo;
    }

    /**
    *@desc fromJson for Condition
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function fromJson(Condition storage _self, string _json) internal returns(bool succ) {
        _self.reset();
        if(LibJson.push(_json) == 0) {
            LibLog.log("json empty");
            return false;
        }

        if (!_json.isJson()) {
            LibJson.pop();
            return false;
        }

        _self.userId = _json.jsonRead("userId").toAddress();
        _self.businessNo = _json.jsonRead("businessNo");
        _self.initDate = _json.jsonRead("initDate").toUint();
        _self.minInitDate = _json.jsonRead("minInitDate").toUint();
        _self.maxInitDate = _json.jsonRead("maxInitDate").toUint();
        _self.minInvoiceDate = _json.jsonRead("minInvoiceDate").toUint();
        _self.maxInvoiceDate = _json.jsonRead("maxInvoiceDate").toUint();
        _self.invoiceDate = _json.jsonRead("invoiceDate").toUint();
        _self.invoiceTitle = _json.jsonRead("invoiceTitle");
        _self.receiver = _json.jsonRead("receiver");
        _self.containInit = _json.jsonRead("containInit").toUint();
        _self.containWaitInitAudit = _json.jsonRead("containWaitInitAudit").toUint();
        _self.containWaitFinalAudit = _json.jsonRead("containWaitFinalAudit").toUint();
        _self.invoiceType = InvoiceType(_json.jsonRead("invoiceType").toUint());
        _self.status = InvoiceStatus(_json.jsonRead("status").toUint());
        _self.statuses.fromJsonArray(_json.jsonRead("statuses"));
        _self.pageSize = _json.jsonRead("pageSize").toUint();
        _self.pageNo = _json.jsonRead("pageNo").toUint();
        
        LibJson.pop();
        return true;
    }

    function reset(Condition storage _self) internal {
        delete _self.userId;
        delete _self.businessNo;
        delete _self.initDate;
        delete _self.minInitDate;
        delete _self.maxInitDate;
        delete _self.minInvoiceDate;
        delete _self.maxInvoiceDate;
        delete _self.invoiceDate;
        delete _self.invoiceTitle;
        delete _self.receiver;
        delete _self.containInit;
        delete _self.containWaitInitAudit;
        delete _self.containWaitFinalAudit;
        delete _self.invoiceType;
        delete _self.status;
        _self.statuses.length = 0;
        delete _self.pageSize;
        delete _self.pageNo;
    }

    /**
    *@desc fromJson for Invoice
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function fromJson(Invoice storage _self, string _json) internal returns(bool succ) {
        _self.reset();
        if(LibJson.push(_json) == 0) {
            LibLog.log("json empty");
            return false;
        }

        if (!_json.isJson()) {
            LibJson.pop();
            return false;
        }

        _self.id = _json.jsonRead("id").toUint();
        _self.userId = _json.jsonRead("userId").toAddress();
        _self.customerType = CustomerType(_json.jsonRead("customerType").toUint());
        _self.vatCustomerType = VatCustomerType(_json.jsonRead("vatCustomerType").toUint());
        _self.bizId = _json.jsonRead("bizId").toUint();
        _self.pledgeeName = _json.jsonRead("pledgeeName");
        _self.pledgorName = _json.jsonRead("pledgorName");
        _self.businessNo = _json.jsonRead("businessNo");
        _self.pledgeContractNo = _json.jsonRead("pledgeContractNo");
        _self.invoiceTitle = _json.jsonRead("invoiceTitle");
        _self.invoiceType = InvoiceType(_json.jsonRead("invoiceType").toUint());
        _self.bizType = BizType(_json.jsonRead("bizType").toUint());
        _self.invoiceAmount = _json.jsonRead("invoiceAmount").toUint();
        _self.status = InvoiceStatus(_json.jsonRead("status").toUint());
        _self.handleChannel = HandleChannel(_json.jsonRead("handleChannel").toUint());
        _self.initDate = _json.jsonRead("initDate").toUint();
        _self.invoiceDate = _json.jsonRead("invoiceDate").toUint();
        _self.taxpayerType = TaxpayerType(_json.jsonRead("taxpayerType").toUint());
        _self.taxpayerIdentifyNo = _json.jsonRead("taxpayerIdentifyNo");
        _self.customerName = _json.jsonRead("customerName");
        _self.certCode = _json.jsonRead("certCode");
        _self.depositBank = _json.jsonRead("depositBank");
        _self.bankAccount = _json.jsonRead("bankAccount");
        _self.phoneNumber = _json.jsonRead("phoneNumber");
        _self.companyAddress = _json.jsonRead("companyAddress");
        _self.contactIDcard = _json.jsonRead("contactIDcard");
        _self.email = _json.jsonRead("email");
        _self.receiver = _json.jsonRead("receiver");
        _self.mobile = _json.jsonRead("mobile");
        _self.phone = _json.jsonRead("phone");
        _self.deliveryWay = DeliveryWay(_json.jsonRead("deliveryWay").toUint());
        _self.companyNo = CompanyNo(_json.jsonRead("companyNo").toUint());
        _self.deliveryNo = _json.jsonRead("deliveryNo");
        _self.postCode = _json.jsonRead("postCode");
        _self.detailAddress = _json.jsonRead("detailAddress");
        _self.receiverUnit = _json.jsonRead("receiverUnit");
        _self.bizLicenseFileId = _json.jsonRead("bizLicenseFileId");
        _self.bizLicenseFileName = _json.jsonRead("bizLicenseFileName");
        _self.contactIDcardFileId = _json.jsonRead("contactIDcardFileId");
        _self.contactIDcardFileName = _json.jsonRead("contactIDcardFileName");
        _self.taxRegistFileId = _json.jsonRead("taxRegistFileId");
        _self.taxRegistFileName = _json.jsonRead("taxRegistFileName");
        _self.depositOpenFileId = _json.jsonRead("depositOpenFileId");
        _self.depositOpenFileName = _json.jsonRead("depositOpenFileName");
        _self.qualificationFileId = _json.jsonRead("qualificationFileId");
        _self.qualificationFileName = _json.jsonRead("qualificationFileName");
        _self.initAuditClaimer = _json.jsonRead("initAuditClaimer").toAddress();
        _self.finalAuditClaimer = _json.jsonRead("finalAuditClaimer").toAddress();
        _self.initAuditOpinion = _json.jsonRead("initAuditOpinion");
        _self.finalAuditOpinion = _json.jsonRead("finalAuditOpinion");
        
        LibJson.pop();
        return true;
    }

    /**
    *@desc toJson for Invoice
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function toJson(Invoice storage _self) internal constant returns (string _json) {
        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("id", _self.id);
        len = LibStack.appendKeyValue("userId", _self.userId);
        len = LibStack.appendKeyValue("customerType", uint(_self.customerType));
        len = LibStack.appendKeyValue("vatCustomerType", uint(_self.vatCustomerType));
        len = LibStack.appendKeyValue("bizId", _self.bizId);
        len = LibStack.appendKeyValue("pledgeeName", _self.pledgeeName);
        len = LibStack.appendKeyValue("pledgorName", _self.pledgorName);
        len = LibStack.appendKeyValue("businessNo", _self.businessNo);
        len = LibStack.appendKeyValue("pledgeContractNo", _self.pledgeContractNo);
        len = LibStack.appendKeyValue("invoiceTitle", _self.invoiceTitle);
        len = LibStack.appendKeyValue("invoiceType", uint(_self.invoiceType));
        len = LibStack.appendKeyValue("bizType", uint(_self.bizType));
        len = LibStack.appendKeyValue("invoiceAmount", _self.invoiceAmount);
        len = LibStack.appendKeyValue("status", uint(_self.status));
        len = LibStack.appendKeyValue("handleChannel", uint(_self.handleChannel));
        len = LibStack.appendKeyValue("initDate", _self.initDate);
        len = LibStack.appendKeyValue("invoiceDate", _self.invoiceDate);
        len = LibStack.appendKeyValue("taxpayerType", uint(_self.taxpayerType));
        len = LibStack.appendKeyValue("taxpayerIdentifyNo", _self.taxpayerIdentifyNo);
        len = LibStack.appendKeyValue("customerName", _self.customerName);
        len = LibStack.appendKeyValue("certCode", _self.certCode);
        len = LibStack.appendKeyValue("depositBank", _self.depositBank);
        len = LibStack.appendKeyValue("bankAccount", _self.bankAccount);
        len = LibStack.appendKeyValue("phoneNumber", _self.phoneNumber);
        len = LibStack.appendKeyValue("companyAddress", _self.companyAddress);
        len = LibStack.appendKeyValue("contactIDcard", _self.contactIDcard);
        len = LibStack.appendKeyValue("email", _self.email);
        len = LibStack.appendKeyValue("receiver", _self.receiver);
        len = LibStack.appendKeyValue("mobile", _self.mobile);
        len = LibStack.appendKeyValue("phone", _self.phone);
        len = LibStack.appendKeyValue("deliveryWay", uint(_self.deliveryWay));
        len = LibStack.appendKeyValue("companyNo", uint(_self.companyNo));
        len = LibStack.appendKeyValue("deliveryNo", _self.deliveryNo);
        len = LibStack.appendKeyValue("postCode", _self.postCode);
        len = LibStack.appendKeyValue("detailAddress", _self.detailAddress);
        len = LibStack.appendKeyValue("receiverUnit", _self.receiverUnit);
        len = LibStack.appendKeyValue("bizLicenseFileId", _self.bizLicenseFileId);
        len = LibStack.appendKeyValue("bizLicenseFileName", _self.bizLicenseFileName);
        len = LibStack.appendKeyValue("contactIDcardFileId", _self.contactIDcardFileId);
        len = LibStack.appendKeyValue("contactIDcardFileName", _self.contactIDcardFileName);
        len = LibStack.appendKeyValue("taxRegistFileId", _self.taxRegistFileId);
        len = LibStack.appendKeyValue("taxRegistFileName", _self.taxRegistFileName);
        len = LibStack.appendKeyValue("depositOpenFileId", _self.depositOpenFileId);
        len = LibStack.appendKeyValue("depositOpenFileName", _self.depositOpenFileName);
        len = LibStack.appendKeyValue("qualificationFileId", _self.qualificationFileId);
        len = LibStack.appendKeyValue("qualificationFileName", _self.qualificationFileName);
        len = LibStack.appendKeyValue("initAuditClaimer", _self.initAuditClaimer);
        len = LibStack.appendKeyValue("finalAuditClaimer", _self.finalAuditClaimer);
        len = LibStack.appendKeyValue("initAuditOpinion", _self.initAuditOpinion);
        len = LibStack.appendKeyValue("finalAuditOpinion", _self.finalAuditOpinion);
        len = LibStack.append("}");
        _json = LibStack.popex(len);
    }

    /**
    *@desc update for Invoice
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function update(Invoice storage _self, string _json) internal returns(bool succ) {
        if(LibJson.push(_json) == 0) {
            return false;
        }

        if (!_json.isJson()) {
            LibJson.pop();
            return false;
        }

        if (_json.jsonKeyExists("id"))
            _self.id = _json.jsonRead("id").toUint();
        if (_json.jsonKeyExists("userId"))
            _self.userId = _json.jsonRead("userId").toAddress();
        if (_json.jsonKeyExists("customerType"))
            _self.customerType = CustomerType(_json.jsonRead("customerType").toUint());
        if (_json.jsonKeyExists("vatCustomerType"))
            _self.vatCustomerType = VatCustomerType(_json.jsonRead("vatCustomerType").toUint());
        if (_json.jsonKeyExists("bizId"))
            _self.bizId = _json.jsonRead("bizId").toUint();
        if (_json.jsonKeyExists("pledgeeName"))
            _self.pledgeeName = _json.jsonRead("pledgeeName");
        if (_json.jsonKeyExists("pledgorName"))
            _self.pledgorName = _json.jsonRead("pledgorName");
        if (_json.jsonKeyExists("businessNo"))
            _self.businessNo = _json.jsonRead("businessNo");
        if (_json.jsonKeyExists("pledgeContractNo"))
            _self.pledgeContractNo = _json.jsonRead("pledgeContractNo");
        if (_json.jsonKeyExists("invoiceTitle"))
            _self.invoiceTitle = _json.jsonRead("invoiceTitle");
        if (_json.jsonKeyExists("invoiceType"))
            _self.invoiceType = InvoiceType(_json.jsonRead("invoiceType").toUint());
        if (_json.jsonKeyExists("bizType"))
            _self.bizType = BizType(_json.jsonRead("bizType").toUint());
        if (_json.jsonKeyExists("invoiceAmount"))
            _self.invoiceAmount = _json.jsonRead("invoiceAmount").toUint();
        if (_json.jsonKeyExists("status"))
            _self.status = InvoiceStatus(_json.jsonRead("status").toUint());
        if (_json.jsonKeyExists("handleChannel"))
            _self.handleChannel = HandleChannel(_json.jsonRead("handleChannel").toUint());
        if (_json.jsonKeyExists("initDate"))
            _self.initDate = _json.jsonRead("initDate").toUint();
        if (_json.jsonKeyExists("invoiceDate"))
            _self.invoiceDate = _json.jsonRead("invoiceDate").toUint();
        if (_json.jsonKeyExists("taxpayerType"))
            _self.taxpayerType = TaxpayerType(_json.jsonRead("taxpayerType").toUint());
        if (_json.jsonKeyExists("taxpayerIdentifyNo"))
            _self.taxpayerIdentifyNo = _json.jsonRead("taxpayerIdentifyNo");
        if (_json.jsonKeyExists("customerName"))
            _self.customerName = _json.jsonRead("customerName");
        if (_json.jsonKeyExists("certCode"))
            _self.certCode = _json.jsonRead("certCode");
        if (_json.jsonKeyExists("depositBank"))
            _self.depositBank = _json.jsonRead("depositBank");
        if (_json.jsonKeyExists("bankAccount"))
            _self.bankAccount = _json.jsonRead("bankAccount");
        if (_json.jsonKeyExists("phoneNumber"))
            _self.phoneNumber = _json.jsonRead("phoneNumber");
        if (_json.jsonKeyExists("companyAddress"))
            _self.companyAddress = _json.jsonRead("companyAddress");
        if (_json.jsonKeyExists("contactIDcard"))
            _self.contactIDcard = _json.jsonRead("contactIDcard");
        if (_json.jsonKeyExists("email"))
            _self.email = _json.jsonRead("email");
        if (_json.jsonKeyExists("receiver"))
            _self.receiver = _json.jsonRead("receiver");
        if (_json.jsonKeyExists("mobile"))
            _self.mobile = _json.jsonRead("mobile");
        if (_json.jsonKeyExists("phone"))
            _self.phone = _json.jsonRead("phone");
        if (_json.jsonKeyExists("deliveryWay"))
            _self.deliveryWay = DeliveryWay(_json.jsonRead("deliveryWay").toUint());
        if (_json.jsonKeyExists("companyNo"))
            _self.companyNo = CompanyNo(_json.jsonRead("companyNo").toUint());
        if (_json.jsonKeyExists("deliveryNo"))
            _self.deliveryNo = _json.jsonRead("deliveryNo");
        if (_json.jsonKeyExists("postCode"))
            _self.postCode = _json.jsonRead("postCode");
        if (_json.jsonKeyExists("detailAddress"))
            _self.detailAddress = _json.jsonRead("detailAddress");
        if (_json.jsonKeyExists("receiverUnit"))
            _self.receiverUnit = _json.jsonRead("receiverUnit");
        if (_json.jsonKeyExists("bizLicenseFileId"))
            _self.bizLicenseFileId = _json.jsonRead("bizLicenseFileId");
        if (_json.jsonKeyExists("bizLicenseFileName"))
            _self.bizLicenseFileName = _json.jsonRead("bizLicenseFileName");
        if (_json.jsonKeyExists("contactIDcardFileId"))
            _self.contactIDcardFileId = _json.jsonRead("contactIDcardFileId");
        if (_json.jsonKeyExists("contactIDcardFileName"))
            _self.contactIDcardFileName = _json.jsonRead("contactIDcardFileName");
        if (_json.jsonKeyExists("taxRegistFileId"))
            _self.taxRegistFileId = _json.jsonRead("taxRegistFileId");
        if (_json.jsonKeyExists("taxRegistFileName"))
            _self.taxRegistFileName = _json.jsonRead("taxRegistFileName");
        if (_json.jsonKeyExists("depositOpenFileId"))
            _self.depositOpenFileId = _json.jsonRead("depositOpenFileId");
        if (_json.jsonKeyExists("depositOpenFileName"))
            _self.depositOpenFileName = _json.jsonRead("depositOpenFileName");
        if (_json.jsonKeyExists("qualificationFileId"))
            _self.qualificationFileId = _json.jsonRead("qualificationFileId");
        if (_json.jsonKeyExists("qualificationFileName"))
            _self.qualificationFileName = _json.jsonRead("qualificationFileName");
        if (_json.jsonKeyExists("initAuditClaimer"))
            _self.initAuditClaimer = _json.jsonRead("initAuditClaimer").toAddress();
        if (_json.jsonKeyExists("finalAuditClaimer"))
            _self.finalAuditClaimer = _json.jsonRead("finalAuditClaimer").toAddress();
        if (_json.jsonKeyExists("initAuditOpinion"))
            _self.initAuditOpinion = _json.jsonRead("initAuditOpinion");
        if (_json.jsonKeyExists("finalAuditOpinion"))
            _self.finalAuditOpinion = _json.jsonRead("finalAuditOpinion");
        
        LibJson.pop();
        return true;
    }

    /**
    *@desc reset for Invoice
    *      Generated by juzhen SolidityStructTool automatically.
    *      Not to edit this code manually.
    */
    function reset(Invoice storage _self) internal {
        delete _self.id;
        delete _self.userId;
        delete _self.customerType;
        delete _self.vatCustomerType;
        delete _self.bizId;
        delete _self.pledgeeName;
        delete _self.pledgorName;
        delete _self.businessNo;
        delete _self.pledgeContractNo;
        delete _self.invoiceTitle;
        delete _self.invoiceType;
        delete _self.bizType;
        delete _self.invoiceAmount;
        delete _self.status;
        delete _self.handleChannel;
        delete _self.initDate;
        delete _self.invoiceDate;
        delete _self.taxpayerType;
        delete _self.taxpayerIdentifyNo;
        delete _self.customerName;
        delete _self.certCode;
        delete _self.depositBank;
        delete _self.bankAccount;
        delete _self.phoneNumber;
        delete _self.companyAddress;
        delete _self.contactIDcard;
        delete _self.email;
        delete _self.receiver;
        delete _self.mobile;
        delete _self.phone;
        delete _self.deliveryWay;
        delete _self.companyNo;
        delete _self.deliveryNo;
        delete _self.postCode;
        delete _self.detailAddress;
        delete _self.receiverUnit;
        delete _self.bizLicenseFileId;
        delete _self.bizLicenseFileName;
        delete _self.contactIDcardFileId;
        delete _self.contactIDcardFileName;
        delete _self.taxRegistFileId;
        delete _self.taxRegistFileName;
        delete _self.depositOpenFileId;
        delete _self.depositOpenFileName;
        delete _self.qualificationFileId;
        delete _self.qualificationFileName;
        delete _self.initAuditClaimer;
        delete _self.finalAuditClaimer;
        delete _self.initAuditOpinion;
        delete _self.finalAuditOpinion;
    }
}
