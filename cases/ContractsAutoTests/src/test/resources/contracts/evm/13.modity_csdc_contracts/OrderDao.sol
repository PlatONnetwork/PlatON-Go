pragma solidity ^0.4.25;

import "./strings.sol";

contract OrderDao{
    using strings for *;

    struct SecPledgeApply {
        string id;                    //id 0
        string businessNo;          //交易申请流水，业务单号 1
        string bizId;             //业务基本信息id 2
        string bizId2;             //业务基本信息id 2
        string pledgorId;    //出质人地址 3
        string pledgorName;   //出质人姓名 4
        string pledgeeId;    //质权人地址 5
        string pledgeeName;    //质权人姓名 6
        string managerId;          //经办人id 7

        // LibTradeUser.TradeUser[]        pledgors;       //出质人
        // LibTradeUser.TradeUser          pledgee;        //质权人
        // LibTradeOperator.TradeOperator  tradeOperator;  //经办人

        string financingAmount;       //融资金额 8
        string financingDateStart;    //融资期限起 9
        string financingDateEnd;      //融资期限止 10
        string financingRate;     //融资利率 11
        string financingTarget;       //融资投向 12
        string paymentId;             //缴费信息id 13
        // PayerType payerType;        //付费方
        string warnLine;          //预警线 14
        string closeLine;         //平仓线 15
        string pledgeContractNo;    //质押合同编号 16
        string pledgeContractFileId;        //质押合同文件id 17
        string pledgeContractFileName; //18
        string secPledgeId;               //成功单号 19
        string applyTime;   //20
        string mainContractCode;  //    主合同编号       21
        string mainContractName;  //    主合同名称       22
        // DisputeMethod disputeMethod;    //  争议解决方式      s   “1”：方式一 “2”：方式二 23
        string disputeMethodInfo;  //   争议解决处理内容        s  24

        string payerAccount; //支付者地址 25
        string payerName; //付款人姓名 26
        string payAmount; //支付金额 27
        string receivedAmount;    //到账金额 28
        string receivedTime;      //到账时间 29
        string payType;           //付款方式 30
        string payFlow;           //银行流水号 31

        string invoiceId;         //邮寄发票id 32
        string djfs;              //冻结方式 33
        string pledgeRegisterNo; //质押业务登记编号 34

        // LibPledgeSecurity.PledgeSecurity[] appliedSecurities; //质押证券

        string bizAccount;      //券商备付金账户 35
        // LibAttachInfo.AttachInfo[]  frontAttachments;    //前端操作普通附件
        string ywbh;    //业务编号 36
        string sfje;      //收费金额 37

        string enter400time;  //进入400时间 38
        string pledgeTime;    //质押成功时间 39
    }

    mapping(string => SecPledgeApply) m_secPledgeApplyMap; //id => object
    string[] m_secPledgeApplyIds;//件组ids

    string public part1;
    string public part2;

    function insert_SecPledgeApply(string memory param) public returns(uint) {
        strings.slice memory s = param.toSlice();
        strings.slice memory delim = "-".toSlice();
        string[] memory parts = new string[](s.count(delim));
        for (uint i = 0; i < parts.length; i++) {
            parts[i] = s.split(delim).toString();
        }


        m_secPledgeApplyIds.push(parts[0]);
        SecPledgeApply memory secPledgeApply;
        secPledgeApply.id = parts[0];
        secPledgeApply.businessNo = parts[1];
        secPledgeApply.bizId = parts[2];
        secPledgeApply.bizId2 = parts[3];
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
        secPledgeApply.warnLine = parts[14];
        secPledgeApply.closeLine = parts[15];
        secPledgeApply.pledgeContractNo = parts[16];
        secPledgeApply.pledgeContractFileId = parts[17];
        secPledgeApply.pledgeContractFileName = parts[18];
        secPledgeApply.secPledgeId = parts[19];
        secPledgeApply.applyTime = parts[20];
        secPledgeApply.mainContractCode = parts[21];
        secPledgeApply.mainContractName = parts[22];
        secPledgeApply.disputeMethodInfo = parts[23];
        secPledgeApply.payerAccount = parts[24];
        secPledgeApply.payerName = parts[25];
        secPledgeApply.payAmount = parts[26];
        secPledgeApply.receivedAmount = parts[27];
        secPledgeApply.receivedTime = parts[28];
        secPledgeApply.payType = parts[29];
        secPledgeApply.payFlow = parts[30];
        secPledgeApply.invoiceId = parts[31];
        secPledgeApply.djfs = parts[32];
        secPledgeApply.pledgeRegisterNo = parts[33];
        secPledgeApply.bizAccount = parts[34];
        secPledgeApply.ywbh = parts[35];
        secPledgeApply.sfje = parts[36];
        secPledgeApply.enter400time = parts[37];
        secPledgeApply.pledgeTime = parts[38];


        m_secPledgeApplyMap[secPledgeApply.id] = secPledgeApply;

    }

    function select_SecPledgeApply_byId(string memory _id) public constant returns(string memory) {
        return m_secPledgeApplyMap[_id].businessNo;
    }
}