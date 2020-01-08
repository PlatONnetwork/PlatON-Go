pragma solidity ^0.4.12;
/**
* @file SuppleDemandManager.sol
* @author yiyating
* @time 2017-04-24
* @desc 供需信息管理合约
*/


import "./csdc_library/LibFinancierInfo.sol";
import "./csdc_library/LibInvestorInfo.sol";
import "./csdc_library/LibAudit.sol";
import "./csdc_library/LibComplaint.sol";
import "./csdc_library/LibAppealing.sol";
import "./Sequence.sol";
import "./PerUserManager.sol";

contract SupplyDemandManager is OwnerNamed{
  
    using LibInt for *;
    using LibString for *;
    using LibFinancierInfo for *;
    using LibInvestorInfo for *;
    using LibAudit for *;
    using LibComplaint for *;
    using LibAppealing for *;
    using LibJson for *;

    mapping(uint=>LibInvestorInfo.InvestorInfo) investorInfoMap;
    mapping(uint=>LibFinancierInfo.FinancierInfo) financierInfoMap;
    uint[] invIdList;
    uint[] finIdList;

    /* 供需收藏 */
    mapping(address=>uint[]) collectedFinancierInfo;        //已收藏融资方信息
    mapping(address=>uint[]) collectedInvestorInfo;         //已收藏出资方信息

    /* 投诉&申诉 */

    enum UserStatus {
        OUT_OF_THE_LIST,    //0-不在黑名单中
        IN_THE_LIST,        //1-在黑名单中（未申诉）
        UNAUDITED,          //2-申诉待审核（已申诉）
        TIMES_EXCEEDED      //3-审核次数超出限制（连续两次）
    }

    struct AppealingInfo {
        LibAppealing.Appealing currentAppealing;            //用户当前申诉信息
        LibAppealing.Appealing[] historyAppealings;         //用户历史申诉信息
        UserStatus status;                                  //用户申诉状态map
    }

    mapping(address=>AppealingInfo) appealingInfos;         //所有用户的申诉信息
    mapping(address=>uint[]) complainedInfoIdMap;            //用户投诉过的信息id列表
    address[] blackList;                                    //黑名单用户列表

    Sequence sq;
    PerUserManager pm;

    uint errno_prefix = 16500;
    
    enum InfoStatus {
        NONE,
        CLOSED,     //1-已关闭
        VALID,      //2-有效
        IN_THE_LIST,//3-用户已列入黑名单
        INVALID     //4-已删除
    }

    enum InfoType {
        NONE,
        FINANCIERINFO,  //1-融资信息
        INVESTORINFO    //2-出资信息
    }
    
    enum SuppleDemandError {
        NONE,
        JSON_INVALID,               //1-json格式错误
        INFO_NOT_EXISTS,            //2-供需消息不存在
        PARAM_EMPTY,                //3-传入参数为空
        PARAM_ERROR,                //4-传入参数有误
        INFO_ALREADY_COLLECTED,     //5-消息已被收藏
        INFO_NOT_COLLECTED,         //6-消息未被收藏
        INFO_ALREADY_COMPLAINED,    //7-消息已投诉过
        USER_IN_THE_LIST,           //8-用户已在黑名单(无法收藏)
        NOT_IN_THE_LIST,            //9-不在黑名单
        ALREADY_APPEAL,             //10-已申诉待审核
        APPEAL_TIMES_EXCEEDED,      //11-申诉次数超过上限
        USER_NOT_APPEALED           //12-用户未申诉
    }
    
    enum QueryType { 
        NONE,
        QUERY,          //1-查询
        MATCH           //2-匹配
    }

    uint MAX_PAGESIZE = 10;
    struct Condition {  //查询条件
        string secCode;         //证券代码
        string secName;         //证券简称
        uint[] industry;      //证券所属行业
        uint[] purpose;       //融资用途
        uint maxAmountScale;    //资金规模
        uint minAmountScale;
        uint maxAmountLastDate; //资金期限
        uint minAmountLastDate;
        uint pageNo;
        uint pageSize;
        QueryType queryType;         //查询类型 1-查询，2-匹配
        address userId;         //用户地址 匹配时排除该地址

        uint[] area;          //投资人地区

        uint id;
        uint status;
    }

    Condition _cond;

    function fromJson_Condition(Condition storage _self, string _json) internal returns(bool succ) {
        reset_Condition(_self);

        LibJson.push(_json);
        if (!_json.isJson()) {
            LibJson.pop();
            return false;
        }


        _self.secCode = _json.jsonRead("secCode");
        _self.secName = _json.jsonRead("secName");
        _self.industry.fromJsonArray(_json.jsonRead("industry"));
        _self.purpose.fromJsonArray(_json.jsonRead("purpose"));
        _self.maxAmountScale = _json.jsonRead("maxAmountScale").toUint();
        _self.minAmountScale = _json.jsonRead("minAmountScale").toUint();
        _self.maxAmountLastDate = _json.jsonRead("maxAmountLastDate").toUint();
        _self.minAmountLastDate = _json.jsonRead("minAmountLastDate").toUint();
        _self.pageNo = _json.jsonRead("pageNo").toUint();
        _self.pageSize = _json.jsonRead("pageSize").toUint();
        _self.queryType = QueryType(_json.jsonRead("queryType").toUint());
        _self.userId = _json.jsonRead("userId").toAddress();
        _self.area.fromJsonArray(_json.jsonRead("area"));
        _self.id = _json.jsonRead("id").toUint();
        _self.status = _json.jsonRead("status").toUint();


        if (_self.pageSize > MAX_PAGESIZE) {
            LibLog.log("pageSize is in excess of 10, now it's set to 10.");
            _self.pageSize = MAX_PAGESIZE;
        }

        LibJson.pop();
        return true;
    }

    function reset_Condition(Condition storage _self) internal {
        delete _self.secCode;
        delete _self.secName;
        _self.industry.length = 0;
        _self.purpose.length = 0;
        delete _self.maxAmountScale;
        delete _self.minAmountScale;
        delete _self.maxAmountLastDate;
        delete _self.minAmountLastDate;
        delete _self.pageNo;
        delete _self.pageSize;
        delete _self.queryType;
        delete _self.userId;
        _self.area.length = 0;
        delete _self.id;
        delete _self.status;
    }
    
    LibFinancierInfo.FinancierInfo _tmpFinInfo;
    LibInvestorInfo.InvestorInfo _tmpInvInfo;
    LibComplaint.Complaint _tmpComplaint;
    LibAppealing.Appealing _tmpAppealing;
    LibAudit.Audit _tmpAudit;

    function SupplyDemandManager(){
        register("CsdcModule", "0.0.1.0", "SupplyDemandManager", "0.0.1.0");
    }

    modifier getSq(){ sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0")); _;}

    modifier getPm(){ pm = PerUserManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "PerUserManager", "0.0.1.0")); _;}


    function getInfoById(uint _id, uint _type) constant returns (string _ret) {
        string memory item = "";
        uint _count = 0;
        if(_type == 1 && financierInfoMap[_id].id != 0) {
            item = financierInfoMap[_id].toJson();
            _count = 1;
        } else if(_type == 2 && investorInfoMap[_id].id != 0) {
            item = investorInfoMap[_id].toJson();
            _count = 1;
        }

        return itemStackPush(item, _count);
    }

    /* 融资方信息接口 - start */
    function getFinancierInfoList(string _json) constant returns (string _ret) {

        if (finIdList.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= finIdList.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = finIdList.length; i >= 1 ; i--) {
            _tmpFinInfo = financierInfoMap[finIdList[i-1]];

            /* 查询条件 - start */

            if (_tmpFinInfo.status == uint(InfoStatus.INVALID)) {
                continue;
            }

            /* 获取所有用户供需信息时剔除黑名单用户信息 */
            if (_cond.userId == address(0) && appealingInfos[_tmpFinInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
                continue;
            }

            if (_json.jsonKeyExists("id") &&  _cond.id != 0 && _cond.id != _tmpFinInfo.id) {
                continue;
            }

            if (_json.jsonKeyExists("userId") &&  _cond.userId != address(0) && _cond.userId != _tmpFinInfo.userId) {
                continue;
            }

            if (_json.jsonKeyExists("status") &&  _cond.status != 0 && _cond.status != _tmpFinInfo.status) {
                continue;
            }

            /* 查询条件 - end */

            if (_count++ < _startIndex) {
                continue;
            }

            if (_total < _cond.pageSize) {
                if (_total > 0) {
                    len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpFinInfo.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function pubFinancierInfo(string _json) getSq getPm {
        if(!_tmpFinInfo.fromJson(_json)) {
            Notify(errno_prefix+uint(SuppleDemandError.JSON_INVALID), "json invalid");
            return ;
        }
        if (appealingInfos[_tmpFinInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
            Notify(errno_prefix+uint(SuppleDemandError.USER_IN_THE_LIST), "the user is in the blacklist");
            return;
        }
        uint _id = sq.getSeqNo("SupplyDemangInfo.id");
        _tmpFinInfo.id = _id;
        _tmpFinInfo.businessNo = sq.genBusinessNo("F", "02").recoveryToString();
        _tmpFinInfo.status = uint(InfoStatus.VALID);
        _tmpFinInfo.createTimestamp = now*1000;
        _tmpFinInfo.userName = LibStack.popex(pm.findNameById(_tmpFinInfo.userId));
        if(pm.userExists(_tmpFinInfo.userId) == 1) {    //判断用户类型
            _tmpFinInfo.userType = 1;   //个人用户
        } else {
            _tmpFinInfo.userType = 2;   //机构用户
        }
        financierInfoMap[_id] = _tmpFinInfo;
        finIdList.push(_id);
        string memory success;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":1 ,\"items\":[");
        success = success.concat("{\"businessNo\":\"", _tmpFinInfo.businessNo,  "\"}]}}");
        Notify(0, success);
    }

    function updateClosePrice(uint _id, uint _price) {
        if(financierInfoMap[_id].id == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
            return;
        }
        financierInfoMap[_id].closePrice = _price;
        Notify(0, "success");
    }

    function closeFinancingInfoById(uint _id) {
        if(financierInfoMap[_id].id == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
            return;
        }
        financierInfoMap[_id].status = uint(InfoStatus.CLOSED);
        Notify(0, "success");
    }

    function deleteFinancingInfoById(uint _id) {
        if(financierInfoMap[_id].id == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
            return;
        }
        financierInfoMap[_id].status = uint(InfoStatus.INVALID);
        Notify(0, "success");
    }

    function queryFinancierInfo(string _json) constant returns (string _ret) {
        if (finIdList.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= finIdList.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = 0; i < finIdList.length; i++) {
            _tmpFinInfo = financierInfoMap[finIdList[i]];
        
            /* 查询条件 - start */

            if ( _tmpFinInfo.status != uint(InfoStatus.VALID) ) {
                continue;
            }

            if (appealingInfos[_tmpFinInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
                continue;
            }
            
            if ( _json.jsonKeyExists("secCode") && !_cond.secCode.equals("") && !_cond.secCode.equals(_tmpFinInfo.secCode) ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("secName") && !_cond.secName.equals("") && !_cond.secName.equals(_tmpFinInfo.secName) ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("industry") && _cond.industry.length != 0 && !isInArray(_tmpFinInfo.industry, _cond.industry) ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("purpose") && _cond.purpose.length != 0 && !isInArray(_tmpFinInfo.purpose, _cond.purpose) ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("maxAmountScale") && _json.jsonKeyExists("minAmountScale") && _cond.maxAmountScale != 0 
                && (_tmpFinInfo.minAmountScale > _cond.maxAmountScale || _tmpFinInfo.maxAmountScale <= _cond.minAmountScale) 
            ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("maxAmountLastDate") && _json.jsonKeyExists("minAmountLastDate") && _cond.maxAmountLastDate != 0 
                && (_tmpFinInfo.minAmountLastDate > _cond.maxAmountLastDate || _tmpFinInfo.maxAmountLastDate <= _cond.minAmountLastDate) 
            ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("userId") && _cond.userId == _tmpFinInfo.userId ) {
              continue;
            }
            
            /* 查询条件 - end */
        
            if (_count++ < _startIndex) {
                continue;
            }
            
            if (_total < _cond.pageSize) {
                if (_total > 0) {
                  len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpFinInfo.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    /* 融资方信息接口 - end */

    /* 出资方信息接口 - start */
    function getInvestorInfoList(string _json) constant returns (string _ret) {
        if (invIdList.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= invIdList.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = invIdList.length; i >= 1 ; i--) {
            _tmpInvInfo = investorInfoMap[invIdList[i-1]];
            
            if (_tmpInvInfo.status == uint(InfoStatus.INVALID)) {
                continue;
            } 

            /* 获取所有用户供需信息时剔除黑名单用户信息 */
            if (_cond.userId == address(0) && appealingInfos[_tmpInvInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
                continue;
            }

            /* 查询条件 - start */
            
            if ( _json.jsonKeyExists("id") && _cond.id != 0 && _cond.id != _tmpInvInfo.id) {
                continue;
            }
            
            if ( _json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != _tmpInvInfo.userId) {
                continue;
            }
            
            if ( _json.jsonKeyExists("status") && _cond.status != 0 && _cond.status != _tmpInvInfo.status) {
                continue;
            }
            
            /* 查询条件 - end */
            
            if (_count++ < _startIndex) {
                continue;
            }
            
            if (_total < _cond.pageSize) {
                if (_total > 0) {
                  len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpInvInfo.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function pubInvestorInfo(string _json) getSq getPm {
        if(!_tmpInvInfo.fromJson(_json)) {
            Notify(errno_prefix+uint(SuppleDemandError.JSON_INVALID), "json invalid");
            return ;
        }
        if (appealingInfos[_tmpInvInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
            Notify(errno_prefix+uint(SuppleDemandError.USER_IN_THE_LIST), "the user is in the blacklist");
            return;
        }
        uint _id = sq.getSeqNo("SupplyDemangInfo.id");
        _tmpInvInfo.id = _id;
        _tmpInvInfo.businessNo = sq.genBusinessNo("F", "01").recoveryToString();
        _tmpInvInfo.status = uint(InfoStatus.VALID);
        _tmpInvInfo.createTimestamp = now*1000;
        _tmpInvInfo.userName = LibStack.popex(pm.findNameById(_tmpInvInfo.userId));
        investorInfoMap[_id] = _tmpInvInfo;
        invIdList.push(_id);
        string memory success;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":1 ,\"items\":[");
        success = success.concat("{\"businessNo\":\"", _tmpInvInfo.businessNo,  "\"}]}}");
        Notify(0, success);
    }

    function closeInvestorInfoById(uint _id) {
        if(investorInfoMap[_id].id == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
            return;
        }
        investorInfoMap[_id].status = uint(InfoStatus.CLOSED);
        Notify(0, "success");
    }

    function deleteInvestorInfoById(uint _id) {
        if(investorInfoMap[_id].id == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
            return;
        }
        investorInfoMap[_id].status = uint(InfoStatus.INVALID);
        Notify(0, "success");
    }

    function queryInvestorInfo(string _json) constant returns (string _ret) {
        if (invIdList.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= invIdList.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = 0; i < invIdList.length; i++) {
            _tmpInvInfo = investorInfoMap[invIdList[i]];
            
            /* 查询条件 - start */
            
            if ( _tmpInvInfo.status != uint(InfoStatus.VALID) ) {
                continue;
            }

            if (appealingInfos[_tmpInvInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
                continue;
            }
            
            if ( _json.jsonKeyExists("area") && _cond.area.length != 0 && !isInArray(_tmpInvInfo.area, _cond.area) ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("industry") && _cond.industry.length != 0 && !hasInArray(_tmpInvInfo.industry, _cond.industry) ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("minAmountScale") && _json.jsonKeyExists("maxAmountScale") && _cond.maxAmountScale != 0 
                && (_tmpInvInfo.minAmountScale > _cond.maxAmountScale || _tmpInvInfo.maxAmountScale < _cond.minAmountScale) 
            ) {
                continue;
            }
            
            if ( _json.jsonKeyExists("minAmountLastDate") && _json.jsonKeyExists("maxAmountLastDate") && _cond.maxAmountLastDate != 0 
                && (_tmpInvInfo.minAmountLastDate > _cond.maxAmountLastDate || _tmpInvInfo.maxAmountLastDate < _cond.minAmountLastDate) 
            ) {
                continue;
            }
        
            if ( _json.jsonKeyExists("userId") && _cond.userId == _tmpInvInfo.userId ) {
              continue;
            }
            
            /* 查询条件 - end */
            
            if (_count++ < _startIndex) {
              continue;
            }
            
            if (_total < _cond.pageSize) {
                if (_total > 0) {
                    len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpInvInfo.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }
    /* 出资方信息接口 - end */

    /* 收藏接口 - start */
    function collectDemandInfo(string _json) {
        address _userId = _json.getStringValueByKey("userId").toAddress();
        uint _id = uint(_json.getIntValueByKey("id"));
        uint _infoType = uint(_json.getIntValueByKey("infoType"));
        if (_userId == address(0) || _id == 0 || _infoType == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_EMPTY), "parameters cannot be empty");
            return;
        }
        if (_infoType != uint(InfoType.FINANCIERINFO) && _infoType != uint(InfoType.INVESTORINFO)) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_ERROR), "infoType must be 1 or 2");
            return;
        }
        if (appealingInfos[_userId].status != UserStatus.OUT_OF_THE_LIST) {
            Notify(errno_prefix+uint(SuppleDemandError.USER_IN_THE_LIST), "the user is in the blacklist");
            return;
        }
        if(_infoType == uint(InfoType.FINANCIERINFO)) {   //收藏融资方信息
            if(financierInfoMap[_id].id == 0) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
                return;
            }
            if(isInArray(_id, collectedFinancierInfo[_userId])) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_ALREADY_COLLECTED), "the financierInfo has already been collected.");
                return;
            }
            collectedFinancierInfo[_userId].push(_id);
        }
        if(_infoType == uint(InfoType.INVESTORINFO)) {   //收藏出资方信息
           if(investorInfoMap[_id].id == 0) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the investorInfo does not exist");
                return;
            }
            if(isInArray(_id, collectedInvestorInfo[_userId])) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_ALREADY_COLLECTED), "the investorInfo has already been collected.");
                return;
            }
            collectedInvestorInfo[_userId].push(_id);
        }
        Notify(0, "success");
    }

    function cancelCollection(string _json) {
        address _userId = _json.getStringValueByKey("userId").toAddress();
        uint _id = uint(_json.getIntValueByKey("id"));
        uint _infoType = uint(_json.getIntValueByKey("infoType"));
        if (_userId == address(0) || _id == 0 || _infoType == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_EMPTY), "parameters cannot be empty");
            return;
        }
        if (_infoType != uint(InfoType.FINANCIERINFO) && _infoType != uint(InfoType.INVESTORINFO)) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_ERROR), "infoType must be 1 or 2");
            return;
        }
        if(_infoType == uint(InfoType.FINANCIERINFO)) {   //取消收藏融资方信息
            if(financierInfoMap[_id].id == 0) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
                return;
            }
            if(!isInArray(_id, collectedFinancierInfo[_userId])) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_COLLECTED), "the financierInfo has not been collected.");
                return;
            }
            removeFromArray(_id, collectedFinancierInfo[_userId]);
        }
        if(_infoType == uint(InfoType.INVESTORINFO)) {   //取消收藏出资方信息
           if(investorInfoMap[_id].id == 0) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the investorInfo does not exist");
                return;
            }
            if(!isInArray(_id, collectedInvestorInfo[_userId])) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_COLLECTED), "the investorInfo has not been collected.");
                return;
            }
            removeFromArray(_id, collectedInvestorInfo[_userId]);
        }
        Notify(0, "success");
    }
    
    function showCollectFinancierInfo(string _json) constant returns (string _ret) {

        fromJson_Condition(_cond, _json);

        if (collectedFinancierInfo[_cond.userId].length <= 0) {
            return itemStackPush("", 0);
        }
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= collectedFinancierInfo[_cond.userId].length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = collectedFinancierInfo[_cond.userId].length; i >= 1 ; i--) {
            _tmpFinInfo = financierInfoMap[collectedFinancierInfo[_cond.userId][i-1]];

            //去除已删除和已关闭的信息
            if (_tmpFinInfo.status == uint(InfoStatus.INVALID) || _tmpFinInfo.status == uint(InfoStatus.CLOSED)) {
                continue;
            }

            if (_count++ < _startIndex) {
                continue;
            }

            if (_total < _cond.pageSize) {
                if (appealingInfos[_tmpFinInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
                    //该信息发布用户已列入黑名单
                    _tmpFinInfo.status = uint(InfoStatus.IN_THE_LIST);
                }
                if (_total > 0) {
                    len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpFinInfo.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }
    
    function showCollectInvestorInfo(string _json) constant returns (string _ret) {
        fromJson_Condition(_cond, _json);

        if (collectedInvestorInfo[_cond.userId].length <= 0) {
            return itemStackPush("", 0);
        }
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= collectedInvestorInfo[_cond.userId].length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = collectedInvestorInfo[_cond.userId].length; i >= 1 ; i--) {
            _tmpInvInfo = investorInfoMap[collectedInvestorInfo[_cond.userId][i-1]];

            //去除已删除和已关闭的信息
            if (_tmpInvInfo.status == uint(InfoStatus.INVALID) || _tmpInvInfo.status == uint(InfoStatus.CLOSED)) {
                continue;
            }

            if (_count++ < _startIndex) {
                continue;
            }

            if (_total < _cond.pageSize) {
                if (appealingInfos[_tmpInvInfo.userId].status != UserStatus.OUT_OF_THE_LIST) {
                    //该信息发布用户已列入黑名单
                    _tmpInvInfo.status = uint(InfoStatus.IN_THE_LIST);
                }
                if (_total > 0) {
                    len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpInvInfo.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }
    /* 收藏接口 - end */

    /* 关闭过期需求  */

    /* 查询所有过期融资信息 */
    function getAllOverdueFinancierInfo(uint _overdueTime) constant returns (string) {
        uint _total = 0;
        uint len = 0;
        len = LibStack.push("");
        for(uint i = 0; i < finIdList.length; i++) {
            _tmpFinInfo = financierInfoMap[finIdList[i]];
            if(_tmpFinInfo.status == uint(InfoStatus.VALID) && _tmpFinInfo.endTimestamp < _overdueTime) {
                if (_total > 0) {
                    len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpFinInfo.toJson());
            }
        }
        return itemStackPush(LibStack.popex(len), _total);
    }

    /* 查询所有过期出资信息 */
    function getAllOverdueInvestorInfo(uint _overdueTime) constant returns (string) {
        uint _total = 0;
        uint len = 0;
        len = LibStack.push("");
        for(uint i = 0; i < invIdList.length; i++) {
            _tmpInvInfo = investorInfoMap[invIdList[i]];
            if(_tmpInvInfo.status == uint(InfoStatus.VALID) && _tmpInvInfo.endTimestamp < _overdueTime) {
                if (_total > 0) {
                    len = LibStack.append(",");
                }
                _total ++;
                len = LibStack.append(_tmpInvInfo.toJson());
            }
        }
        return itemStackPush(LibStack.popex(len), _total);        
    }


    /* overdueTime-过期时间 */
    function closeOverdueInfo(uint _overdueTime) {
        closeOverdueFinancierInfo(_overdueTime);
        closeOverdueInvestorInfo(_overdueTime);
        Notify(0, "success");
    }

    function closeOverdueFinancierInfo(uint _overdueTime) {
        uint len = finIdList.length;
        for(uint i = 0; i < len; i++) {
            _tmpFinInfo = financierInfoMap[finIdList[i]];
            if(_tmpFinInfo.status == uint(InfoStatus.VALID) && _tmpFinInfo.endTimestamp < _overdueTime) {
                financierInfoMap[finIdList[i]].status = uint(InfoStatus.CLOSED);
            }
        }
    }

    function closeOverdueInvestorInfo(uint _overdueTime) {
        uint len = invIdList.length;
        for(uint i = 0; i < len; i++) {
            _tmpInvInfo = investorInfoMap[invIdList[i]];
            if(_tmpInvInfo.status == uint(InfoStatus.VALID) && _tmpInvInfo.endTimestamp < _overdueTime) {
                investorInfoMap[invIdList[i]].status = uint(InfoStatus.CLOSED);
            }
        }
    }

    /* 投诉&申诉接口 - start */

    /* 查询用户是否在黑名单 */
    function isInBlacklist(address _userId) constant returns (uint) {
        return uint(appealingInfos[_userId].status);
    }

    /* 查询已投诉 */
    function getComplainedIdList(address _userId) constant returns (string _ret) {
        _ret = "{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":0,\"items\":[]}}";
        uint len = complainedInfoIdMap[_userId].length;
        if( len == 0) {
            return;
        }
        _ret = "{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":";
        _ret = _ret.concat(len.toString(), ",\"items\":[");
        for(uint i = 0; i < len; i++) {
            if(i > 0) {
                _ret = _ret.concat(",");
            }
            _ret = _ret.concat(complainedInfoIdMap[_userId][i].toString());
        }
        _ret = _ret.concat("]}}");
    }

    /* 投诉接口 */
    function complainInfo(string _json) getPm {
        if(!_tmpComplaint.fromJson(_json)) {
            Notify(errno_prefix+uint(SuppleDemandError.JSON_INVALID), "json invalid");
            return; 
        }

        if (_tmpComplaint.infoId == 0 || _tmpComplaint.infoType == 0 || _tmpComplaint.userId == address(0) || _tmpComplaint.reason.equals("")) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_EMPTY), "parameters cannot be empty");
            return;
        }

        if (_tmpComplaint.infoType != uint(InfoType.FINANCIERINFO) && _tmpComplaint.infoType != uint(InfoType.INVESTORINFO)) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_ERROR), "infoType must be 1 or 2");
            return;
        }

        /* 用户已投诉过该信息 */
        if (isInArray(_tmpComplaint.infoId, complainedInfoIdMap[_tmpComplaint.userId])) {   //用户已投诉过该条信息
            Notify(errno_prefix+uint(SuppleDemandError.INFO_ALREADY_COMPLAINED), "you have already complained the info");
            return;
        }

        address _userId;    //被投诉用户
        if(_tmpComplaint.infoType == uint(InfoType.FINANCIERINFO)) {   //投诉融资方信息
            if(financierInfoMap[_tmpComplaint.infoId].id == 0) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the financierInfo does not exist");
                return;
            }
            _userId = financierInfoMap[_tmpComplaint.infoId].userId;
        }
        if(_tmpComplaint.infoType == uint(InfoType.INVESTORINFO)) {   //投诉出资方信息
            if(investorInfoMap[_tmpComplaint.infoId].id == 0) {
                Notify(errno_prefix+uint(SuppleDemandError.INFO_NOT_EXISTS), "the investorInfo does not exist");
                return;
            }

            _userId = investorInfoMap[_tmpComplaint.infoId].userId;
        }

        /* 被投诉用户已被列入黑名单 */
        if (appealingInfos[_userId].status != UserStatus.OUT_OF_THE_LIST) {
            Notify(errno_prefix+uint(SuppleDemandError.USER_IN_THE_LIST), "the user is in the blacklist");
            return;
        }

        _tmpComplaint.userName = LibStack.popex(pm.findNameById(_tmpComplaint.userId));

        uint complainCount = __complainInfo(_tmpComplaint, InfoType(_tmpComplaint.infoType), _userId);
        complainedInfoIdMap[_tmpComplaint.userId].push(_tmpComplaint.infoId);
        string memory success;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":1 ,\"items\":[");
        success = success.concat("{\"complainCount\":\"", complainCount.toString(),  "\"}]}}");
        Notify(0, success);
    }

    function __complainInfo(LibComplaint.Complaint _complaint, InfoType _type, address _userId) getPm internal returns(uint count) {
        if(appealingInfos[_userId].currentAppealing.userId == address(0)) {  //Map中没有对应用户的申诉信息时，初始化信息
            appealingInfos[_userId].currentAppealing.userId = _userId;
            appealingInfos[_userId].currentAppealing.userName = LibStack.popex(pm.findNameById(_userId));
            appealingInfos[_userId].status = UserStatus.OUT_OF_THE_LIST;
        }
        if(appealingInfos[_userId].status == UserStatus.OUT_OF_THE_LIST) {      //用户还未在黑名单内

            if(_type == InfoType.FINANCIERINFO) {
                appealingInfos[_userId].currentAppealing.financierInfos.push(_complaint);
            } else if (_type == InfoType.INVESTORINFO) {
                appealingInfos[_userId].currentAppealing.investorInfos.push(_complaint);
            }

            count = appealingInfos[_userId].currentAppealing.financierInfos.length + appealingInfos[_userId].currentAppealing.investorInfos.length;
            if (count == 5) {    //该用户收到投诉数量达到5条，进入黑名单
                appealingInfos[_userId].status = UserStatus.IN_THE_LIST;
                blackList.push(_userId);
            }
        }
    }

    /* 申诉接口 */
    function appeal(string _json) {
        if(!_tmpAppealing.fromJson(_json)) {
            Notify(errno_prefix+uint(SuppleDemandError.JSON_INVALID), "json invalid");
            return; 
        }
        if(_tmpAppealing.matter.equals("") || _tmpAppealing.reason.equals("")) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_EMPTY), "parameters cannot be empty");
            return;
        }
        UserStatus _status = appealingInfos[_tmpAppealing.userId].status;
        if(_status == UserStatus.OUT_OF_THE_LIST) {
            Notify(errno_prefix+uint(SuppleDemandError.NOT_IN_THE_LIST), "you are not in the blacklist");
            return;
        }
        if(_status == UserStatus.UNAUDITED) {
            Notify(errno_prefix+uint(SuppleDemandError.ALREADY_APPEAL), "you have already appealed, please wait");
            return;
        }
        if(_status == UserStatus.TIMES_EXCEEDED) {
            Notify(errno_prefix+uint(SuppleDemandError.APPEAL_TIMES_EXCEEDED), "you have already appealed more than twice, you cannot appeal again");
            return;
        }
        appealingInfos[_tmpAppealing.userId].currentAppealing.appeal(_tmpAppealing.matter, _tmpAppealing.reason, _tmpAppealing.contact);
        appealingInfos[_tmpAppealing.userId].status = UserStatus.UNAUDITED;
        Notify(0, "success");
    }

    /* 审核申诉信息 */
    function auditAppealing(string _json) {
        address _userId = _json.getStringValueByKey("userId").toAddress();
        address _auditorId = _json.getStringValueByKey("auditorId").toAddress();
        uint _operateCode = uint(_json.getIntValueByKey("operateCode"));
        string memory _auditComment = _json.getStringValueByKey("auditComment");
        if(_userId == address(0) || _auditorId == address(0) || _operateCode == 0) {
            Notify(errno_prefix+uint(SuppleDemandError.PARAM_EMPTY), "parameters cannot be empty");
            return;
        }
        if(appealingInfos[_userId].status != UserStatus.UNAUDITED) {
            Notify(errno_prefix+uint(SuppleDemandError.USER_NOT_APPEALED), "the user does not appeal");
            return;
        }
        _tmpAudit.create(_auditorId, _operateCode, _auditComment, "SupplyDemand.auditAppealing");
        appealingInfos[_userId].currentAppealing.audit = _tmpAudit;
        appealingInfos[_userId].historyAppealings.push(appealingInfos[_userId].currentAppealing);       //存入历史申诉信息
        if(_tmpAudit.operateCode == LibAudit.OperateCode.PASS) {          
            //审核通过
            appealingInfos[_userId].status = UserStatus.OUT_OF_THE_LIST;

            delete appealingInfos[_userId].currentAppealing.financierInfos; //当前投诉信息清空-1
            delete appealingInfos[_userId].currentAppealing.investorInfos;  //当前投诉信息清空-2
            removeFromArray(_userId, blackList);                            //将该用户移出黑名单
            appealingInfos[_userId].currentAppealing.count = 0;             //申诉次数清零
            
        } else if(_tmpAudit.operateCode == LibAudit.OperateCode.FAIL) {   
            //审核拒绝
            if(++appealingInfos[_userId].currentAppealing.count == 2) {
                appealingInfos[_userId].status = UserStatus.TIMES_EXCEEDED;
            } else {
                appealingInfos[_userId].status = UserStatus.IN_THE_LIST;
            }
        }
        appealingInfos[_userId].currentAppealing.clear();               //当前申诉信息清空
        
        Notify(0, "success");
    }

    /* 获取待审核申诉列表 */
    function getAppealingList(string _json) constant returns (string _ret) {
        if (blackList.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= blackList.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = blackList.length; i >= 1 ; i--) {
            _tmpAppealing = appealingInfos[blackList[i-1]].currentAppealing;

            if(appealingInfos[_tmpAppealing.userId].status != UserStatus.UNAUDITED) {
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
                len = LibStack.append(appealingToJson(_tmpAppealing));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    /* 获取某用户当前申诉信息 */
    function getCurrentAppealing(address _userId) constant returns (string _ret) {
        _ret = "{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":0,\"items\":[]}}";
        if(appealingInfos[_userId].currentAppealing.userId == address(0)) {
            return;
        }
        _ret = "";
        _ret = _ret.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\": 1", ",\"items\":[");
        _ret = _ret.concat(appealingToJson(appealingInfos[_userId].currentAppealing), "]}}");
    }

    /* 获取某用户历史申诉信息 */
    function getHistoryAppealing(address _userId)  constant returns (string _ret) {
        _ret = "{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":0,\"items\":[]}}";

        uint len = appealingInfos[_userId].historyAppealings.length;
        if (len == 0) {
            return;
        }
        _ret = "{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":";
        _ret = _ret.concat(len.toString(), ",\"items\":[");
        for (uint i = 0; i < len; i++) {
            _tmpAppealing = appealingInfos[_userId].historyAppealings[i];
            if (i > 0) {
                _ret = _ret.concat(",");
            }
            _ret = _ret.concat(appealingToJson(_tmpAppealing));
        }
        _ret = _ret.concat("]}}");
    }

    /* 投诉&申诉接口 - end */
    
    /* 以下为内部调用方法 */

    function isInArray(uint _value, uint[] storage _array) internal returns (bool) {
        for(uint i=0; i<_array.length; i++) {
            if(_array[i] == _value) {
                return true;
            }
        }
    }
    
    function hasInArray(uint[] _arraySource, uint[] storage _arrayTarget) internal returns (bool) {
        for(uint i=0; i<_arraySource.length; i++) {
            if(isInArray(_arraySource[i], _arrayTarget)) {
                return true;
            }
        }
    }
    
    function removeFromArray(uint _value, uint[] storage _array) internal {
        bool hasFind = false;
        for(uint i=0; i<_array.length; i++) {
            if(_value == _array[i]) {
                hasFind = true;
            }
            if(hasFind && i < _array.length-1) {
                _array[i] = _array[i+1];
            }
        }
        if(hasFind) {
            _array.length = _array.length -1;   //自动删除末尾元素
        }
    }

    function removeFromArray(address _value, address[] storage _array) internal {
        bool hasFind = false;
        for(uint i=0; i<_array.length; i++) {
            if(_value == _array[i]) {
                hasFind = true;
            }
            if(hasFind && i < _array.length-1) {
                _array[i] = _array[i+1];
            }
        }
        if(hasFind) {
            _array.length = _array.length -1;   //自动删除末尾元素
        }
    }

    function appealingToJson(LibAppealing.Appealing storage _self) internal constant returns (string _json) {

        uint len = 0;
        len = LibStack.push("{");
        len = LibStack.appendKeyValue("userId", uint(_self.userId).toAddrString());
        len = LibStack.appendKeyValue("userName", _self.userName);
        len = LibStack.appendKeyValue("matter", _self.matter);
        len = LibStack.appendKeyValue("count", _self.count);
        len = LibStack.appendKeyValue("reason", _self.reason);
        len = LibStack.appendKeyValue("contact", _self.contact);
        len = LibStack.appendKeyValue("createTime", _self.createTime);
        len = LibStack.append(",");

        len = LibStack.append("\"financierInfos\":[");
        uint _infoId;
        for(uint i = 0; i < _self.financierInfos.length; i++) {
            _infoId = _self.financierInfos[i].infoId;
            len = LibStack.append(_self.financierInfos[i].toJson(financierInfoMap[_infoId].toJson()));
            if(i < _self.financierInfos.length-1) {
                len = LibStack.append(",");
          } 
        }
        len = LibStack.append("],");

        len = LibStack.append("\"investorInfos\":[");
        for(i = 0; i < _self.investorInfos.length; i++) {
            _infoId = _self.investorInfos[i].infoId;
            len = LibStack.append(_self.investorInfos[i].toJson(investorInfoMap[_infoId].toJson()));
            if(i < _self.investorInfos.length-1) {
                len = LibStack.append(",");
            }
        }
        len = LibStack.append("]");

        len = LibStack.appendKeyValue("audit", _self.audit.toJson());
        len = LibStack.append("}");

        return LibStack.popex(len);
    }

    //items入栈
    function itemStackPush(string _items, uint _total) constant internal returns (string _ret){
        uint len = 0;
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
        return LibStack.popex(len);
    }

    event Notify(uint _errorno, string _info);

}