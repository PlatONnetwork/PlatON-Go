pragma solidity ^0.4.25;

import "./strings.sol";

contract NewSecPledgeApplyManager{

    using strings for *;

    struct SecPledgeApply {
        string id;
        string businessNo;
        string bizId;
        string pledgorId;
        string pledgorName;
        string pledgeeId;
        string pledgeeName;
        string managerId;
        string financingAmount;
        string financingDateStart;
        string financingDateEnd;
        string financingRate;
        string financingTarget;
        string paymentId;
        string payerType;
        string warnLine;
        string closeLine;
        string pledgeContractNo;
        string pledgeContractFileId;
        string pledgeContractFileName;
        string secPledgeId;
        // string applyTime;
        // string mainContractCode;
        // string mainContractName;
        // string disputeMethodInfo;  //   争议解决处理内容       24
        // string payerAccount; //支付者地址 25
        // string payerName; //付款人姓名 26
        // string payAmount; //支付金额 27
        // string receivedAmount;    //到账金额 28
        // string receivedTime;      //到账时间 29
        // string payType;           //付款方式 30
        // string payFlow;           //银行流水号 31
        // string invoiceId;         //邮寄发票id 32
        // string djfs;              //冻结方式 33
        // string pledgeRegisterNo; //质押业务登记编号 34
        // string bizAccount;      //券商备付金账户 35
        // string ywbh;    //业务编号 36
        // string sfje;      //收费金额 37
        // string enter400time;  //进入400时间 38
        // string pledgeTime;    //质押成功时间 39

        TradeUser        tradeUser;       //出质人 40
        TradeOperator  tradeOperator;  //经办人 41
        // PledgeSecurity pledgeSecurity; //质押证券 42
    }


    struct TradeOperator {
        string id;
        string brokerId;         //券商id
        string name;            //姓名
        //     string department;      //部门名称
        //     string phone;           //电话
        //     string fax;             //传真
        //     string mobile;          //手机
        //     string email;           //邮箱
    }

    struct TradeUser {
        string      traderType;     //交易者类型
        string      traderId;       //交易者address
        string      account;        //证券账户号
        string      accountType;    //证券账户类型
        string      idType;         //身份证件类型
        //     string      idNo;           //身份证件号/注册号码
        //     string      userType;       //用户类型
        //     string      name;           //姓名
        //     string      isShareholder;  //是否国有股东
        //     string      isSecAccount;   //是否定向资管专用证券账户
        //     string      isFundAccount;  //是否基金账户
        //     string      isAgency;       //是否代办
        //     string      agentName;      //代办人姓名
        //     string      agentMobile;    //代办人移动电话
        //     string      email;          //Email
        //     string      receiverName;   //收件人姓名
        //     string      receiverMobile; //收件人联系方式
        //     string      postCode;       //邮政编码
        //     string      postAddr;       //邮政地址
        //     string      receiverCompany;//收件人单位名称
        //     string      postType;       //邮寄方式
        //     string      expressCompany; //快递公司名称
        //     string      isPost;         //是否邮寄
    }

    struct PledgeSecurity {
        string id;            //id
        string secAccount;  //证券账户
        string secCode;     //证券代码
        string secName;     //证券简称
        //     string secType;   //证券类别（股份性质）
        //     string hostedUnit;  //托管单元
        //     string hostedUnitName;//证券单元名称
        //     string secProperty;//股份性质
        //     string pledgeNum; //质押数量
        //     string remainPledgeNum; //剩余质押数量
        //     string isProfit;      //是否解除红利
        //     string isProfitRemain;        //是否解除红利修改后
        //     string profitAmount;  //现金红利金额
        //     string bonusShareAmount; //剩余红股
        //     string uniAcctNbr; //一码通号
        //     string shareholderIdNo;     //股东证件号码
        //     string freezeNo;            //冻结序号
        //     string subFreezeNo;         //冻结子序号
        //     string totalMarketAmount;     //市场总股本
        //     string shareHoldingRatio;     //持股比例
        //     string pledgeRatio;           //质押比例
        //     string frozenFlag;            //冻结是否成功
        //     string frozenNum;             //已冻结股数
        //     string preFrozenNum;          //预冻结股数
        //     string originalPledgeNum;     //原质押数量
        //     string judiciaryFreezeNum;    //司法冻结股数
        //     string isJudiciaryFrozen;     //是否已司法冻结
        //     string doContinueWithJudiciaryFreeze;  //司法冻结情况下是否继续接触质押
        //     string remainPositionNum;     //剩余持仓股数
    }

    mapping(string => SecPledgeApply) secPledgeApplyMap; //id => object
    string[] secPledgeApplyIds;//件组ids
    PledgeSecurity[] pledgeSecuritys;
    mapping(string =>PledgeSecurity[]) pledgeSecurityArrMap;

    //新增质押申请信息
    function createPledgeApplyCommon(string memory secPledgeApplyJson) public returns(string memory) {
        strings.slice memory s = secPledgeApplyJson.toSlice();
        strings.slice memory delim = "-".toSlice();
        string[] memory parts = new string[](s.count(delim));
        for (uint i = 0; i < parts.length; i++) {
            parts[i] = s.split(delim).toString();
        }

        secPledgeApplyIds.push(parts[0]);

        SecPledgeApply memory secPledgeApply;
        secPledgeApply.id = parts[0];
        secPledgeApply.businessNo = parts[1];
        secPledgeApply.bizId = parts[2];
        secPledgeApply.pledgorId = parts[3];
        secPledgeApply.pledgorName = parts[4];
        secPledgeApply.pledgeeId = parts[5];
        secPledgeApply.pledgeeName = parts[6];
        secPledgeApply.managerId = parts[7];
        secPledgeApply.financingAmount = parts[8];
        secPledgeApply.financingDateStart = parts[9];
        secPledgeApply.financingDateEnd = parts[10];
        secPledgeApply.financingRate = parts[11];
        secPledgeApply.financingTarget = parts[12];
        secPledgeApply.paymentId = parts[13];
        secPledgeApply.payerType = parts[14];
        secPledgeApply.warnLine = parts[15];
        secPledgeApply.closeLine = parts[16];
        secPledgeApply.pledgeContractNo = parts[17];
        secPledgeApply.pledgeContractFileId = parts[18];
        secPledgeApply.pledgeContractFileName = parts[19];
        secPledgeApply.secPledgeId = parts[20];
        // secPledgeApply.applyTime = parts[21];
        // secPledgeApply.mainContractCode = parts[22];
        // secPledgeApply.mainContractName = parts[23];
        // secPledgeApply.disputeMethodInfo = parts[24];
        // secPledgeApply.payerAccount = parts[25];
        // secPledgeApply.payerName = parts[26];
        // secPledgeApply.payAmount = parts[27];
        // secPledgeApply.receivedAmount = parts[28];
        // secPledgeApply.receivedTime = parts[29];
        // secPledgeApply.payType = parts[30];
        // secPledgeApply.payFlow = parts[31];
        // secPledgeApply.invoiceId = parts[32];
        // secPledgeApply.djfs = parts[33];
        // secPledgeApply.pledgeRegisterNo = parts[34];
        // secPledgeApply.bizAccount = parts[35];
        // secPledgeApply.ywbh = parts[36];
        // secPledgeApply.sfje = parts[37];
        // secPledgeApply.enter400time = parts[38];
        // secPledgeApply.pledgeTime = parts[39];

        //新增证券信息
        // strings.slice memory s1 = pledgeSecurityJson.toSlice();
        // strings.slice memory delim1 = "-".toSlice();
        // string[] memory parts1 = new string[](s1.count(delim1));
        // for (uint i = 0; i < parts1.length; i++) {
        //     parts1[i] = s1.split(delim1).toString();
        // }

        PledgeSecurity memory pledgeSecurity;
        pledgeSecurity.id = parts[8];
        pledgeSecurity.secAccount = parts[11];
        pledgeSecurity.secCode = parts[12];
        pledgeSecurity.secName = parts[13];

        pledgeSecuritys.push(pledgeSecurity);
        pledgeSecurityArrMap[parts[0]] = pledgeSecuritys;

        // secPledgeApply.pledgeSecuritys = pledgeSecurityArr;
        // secPledgeApply.pledgeSecurity.secType = parts1[4];
        // secPledgeApply.pledgeSecurity.hostedUnit = parts1[5];
        // secPledgeApply.pledgeSecurity.hostedUnitName = parts1[6];
        // secPledgeApply.pledgeSecurity.secProperty = parts1[7];
        // secPledgeApply.pledgeSecurity.pledgeNum = parts1[8];
        // secPledgeApply.pledgeSecurity.remainPledgeNum = parts1[9];
        // secPledgeApply.pledgeSecurity.isProfit = parts1[10];
        // secPledgeApply.pledgeSecurity.isProfitRemain = parts1[11];
        // secPledgeApply.pledgeSecurity.profitAmount = parts1[12];
        // secPledgeApply.pledgeSecurity.bonusShareAmount = parts1[13];
        // secPledgeApply.pledgeSecurity.uniAcctNbr = parts1[14];
        // secPledgeApply.pledgeSecurity.shareholderIdNo = parts1[15];
        // secPledgeApply.pledgeSecurity.freezeNo = parts1[16];
        // secPledgeApply.pledgeSecurity.subFreezeNo = parts1[17];
        // secPledgeApply.pledgeSecurity.totalMarketAmount = parts1[18];
        // secPledgeApply.pledgeSecurity.shareHoldingRatio = parts1[19];
        // secPledgeApply.pledgeSecurity.pledgeRatio = parts1[20];
        // secPledgeApply.pledgeSecurity.frozenFlag = parts1[21];
        // secPledgeApply.pledgeSecurity.frozenNum = parts1[22];
        // secPledgeApply.pledgeSecurity.preFrozenNum = parts1[23];
        // secPledgeApply.pledgeSecurity.originalPledgeNum = parts1[24];
        // secPledgeApply.pledgeSecurity.judiciaryFreezeNum = parts1[25];
        // secPledgeApply.pledgeSecurity.isJudiciaryFrozen = parts1[26];
        // secPledgeApply.pledgeSecurity.doContinueWithJudiciaryFreeze = parts1[27];
        // secPledgeApply.pledgeSecurity.remainPositionNum = parts1[28];



        //  //新增交易用户信息
        // strings.slice memory s2 = tradeUserJson.toSlice();
        // strings.slice memory delim2 = "-".toSlice();
        // string[] memory parts2 = new string[](s2.count(delim2));
        // for (uint i = 0; i < parts2.length; i++) {
        //     parts2[i] = s2.split(delim2).toString();
        // }

        secPledgeApply.tradeUser.traderType = parts[0];
        secPledgeApply.tradeUser.traderId = parts[1];
        secPledgeApply.tradeUser.account = parts[2];
        secPledgeApply.tradeUser.accountType = parts[3];
        secPledgeApply.tradeUser.idType = parts[4];
        // secPledgeApply.tradeUser.idNo = parts2[5];
        // secPledgeApply.tradeUser.userType = parts2[6];
        // secPledgeApply.tradeUser.name = parts2[7];
        // secPledgeApply.tradeUser.isShareholder = parts2[8];
        // secPledgeApply.tradeUser.isSecAccount = parts2[9];
        // secPledgeApply.tradeUser.isFundAccount = parts2[10];
        // secPledgeApply.tradeUser.isAgency = parts2[11];
        // secPledgeApply.tradeUser.agentName = parts2[12];
        // secPledgeApply.tradeUser.agentMobile = parts2[13];
        // secPledgeApply.tradeUser.email = parts2[14];
        // secPledgeApply.tradeUser.receiverName = parts2[15];
        // secPledgeApply.tradeUser.receiverMobile = parts2[16];
        // secPledgeApply.tradeUser.postCode = parts2[17];
        // secPledgeApply.tradeUser.postAddr = parts2[18];
        // secPledgeApply.tradeUser.receiverCompany = parts2[19];
        // secPledgeApply.tradeUser.postType = parts2[20];
        // secPledgeApply.tradeUser.expressCompany = parts2[21];
        // secPledgeApply.tradeUser.isPost = parts2[22];


        // //新增操作用户
        // strings.slice memory s3 = tradeOperatorJson.toSlice();
        // strings.slice memory delim3 = "-".toSlice();
        // string[] memory parts3 = new string[](s3.count(delim3));
        // for (uint i = 0; i < parts3.length; i++) {
        //     parts3[i] = s3.split(delim3).toString();
        // }

        secPledgeApply.tradeOperator.id = parts[4];
        secPledgeApply.tradeOperator.brokerId = parts[5];
        secPledgeApply.tradeOperator.name = parts[6];
        // secPledgeApply.tradeOperator.department = parts3[3];
        // secPledgeApply.tradeOperator.phone = parts3[4];
        // secPledgeApply.tradeOperator.fax = parts3[5];
        // secPledgeApply.tradeOperator.mobile = parts3[6];
        // secPledgeApply.tradeOperator.email = parts3[7];

        secPledgeApplyMap[secPledgeApply.id] = secPledgeApply;
        // secPledgeApplyIds.push(parts[0]);

        return secPledgeApply.id;

    }

    //查询质押申请信息（将所有的返回参数拼接成字符串）
    function select_SecPledgeApply_byId(string memory _id) public constant returns(string memory) {
        return secPledgeApplyMap[_id].businessNo;
    }

    //查询交易用户
    function select_tradeUser_byId(string memory _id) public constant returns(string memory) {
        return secPledgeApplyMap[_id].tradeUser.account;
    }

    //查询操作用户
    function select_tradeOperator_bytId(string memory _id) public constant returns(string memory) {
        return secPledgeApplyMap[_id].tradeOperator.name;
    }

    //查询证券信息用户
    function select_pledgeSecurity_bytId(string memory _id) public constant returns(string memory) {
        return pledgeSecurityArrMap[_id][0].secName;
    }
}