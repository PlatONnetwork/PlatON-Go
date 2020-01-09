pragma solidity ^0.4.12;
/**
* file EvidenceManager.sol
* author xuhui
* time 2017-05-11
* desc the defination of Post Information
*/

import "./csdc_library/LibEvidence.sol";
import "./Sequence.sol";
import "./SecPledgeManager.sol";

contract EvidenceManager is OwnerNamed{
    using LibInt for *;
    using LibString for *;
    using LibEvidence for *;
    using LibSecPledge for *;
    using LibJson for *;
    using LibStack for *;
    
    //错误码
    enum Error{
        NO_ERROR,
        INVOICEID_ERROR,
        ID_NOT_EXIST,
        EVIDENCETYPE_EMPTY,
        BAD_PARAMETER,
        BIZID_EMPTY,
        BIZID_NOTEXIST,
        BIZSTATUS_ERROE,
        INVOICESTATUS_ERROR,
        AUDITORID_EMPTY,
        USER_NOT_EXISTS,
        STATUS_ERROR,
        PLEDGEREGISTERFILE_EMPTY,
        PLEDGEREGISTERFILEID_EMPTY,
        PLEDGEREGISTERFILENAME_EMPTY,
        PLEDGEEFULLNAME_EMPTY,
        PLEDGORFULLNAME_EMPTY,
        RECEIVER_EMPTY,
        MOBILE_EMPTY,
        RECEIVERUNIT_EMPTY,
        DELIVERYNO_EMPTY,
        DETAILADDRESS_EMPTY,
        POSTCODE_EMPTY,
        DELIVERYWAY_INCORRECT,
        COMPANYNO_INCORRECT,
        HANDLECHANNEL_EMPTY,
        BUSINESSNO_REPEAT
    }
    uint errno_prefix = 16500;
    
    //修改数据结构
    mapping(address=>uint[]) user2evidenceIds;
    mapping(uint => LibEvidence.Evidence) evidenceMap;
    uint[] evidenceIds;
    uint[] tempList;
    
    LibEvidence.Evidence _evidenceTemp;
    Sequence sq;
    SecPledgeManager secpledgemanager;

    modifier getSq(){ sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0")); _;}

    function EvidenceManager(){
        register("CsdcModule", "0.0.1.0", "EvidenceManager", '0.0.1.0');
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
        secpledgemanager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager","0.0.1.0"));
    }

    function initEvidence(string _invoiceJson) getSq public returns(uint _id){
        _evidenceTemp.reset();
        if(!_evidenceTemp.fromJson(_invoiceJson)){
            LibLog.log("initEvidence json invalid");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"json invalid");
            return; 
        }
        if(_evidenceTemp.evidenceType == LibEvidence.EvidenceType.NONE){
            LibLog.log("initEvidence evidenceType can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.EVIDENCETYPE_EMPTY),"evidenceType can not empty");
            return;
        }
        if(_evidenceTemp.handleChannel == LibEvidence.HandleChannel.NONE){
            LibLog.log("initEvidence handleChannel can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.HANDLECHANNEL_EMPTY),"handleChannel can not empty");
            return;
        }
        if(_evidenceTemp.businessNo.equals("")){
            LibLog.log("initEvidence businessNo can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.BIZID_NOTEXIST),"businessNo can not empty");
            return;
        }
        //业务流水号不能重复
        for(uint i = 0; i< evidenceIds.length;i++){
            if(evidenceMap[evidenceIds[i]].businessNo.equals(_evidenceTemp.businessNo)){
                LibLog.log("initEvidence","initEvidence: businessNo already exist");
                Notify(errno_prefix+uint(Error.BUSINESSNO_REPEAT),"initInvoice: businessNo already exist");
                return;
            }
        }
        if(_evidenceTemp.pledgeeName.equals("")){
            LibLog.log("initEvidence pledgeeName can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.PLEDGEEFULLNAME_EMPTY),"pledgeeName can not empty");
            return;
        }
        if(_evidenceTemp.pledgorName.equals("")){
            LibLog.log("initEvidence pledgorName can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.PLEDGORFULLNAME_EMPTY),"pledgorName can not empty");
            return;
        }
        _evidenceTemp.id = sq.getSeqNo("Evidence.id");
        _evidenceTemp.initDate = now*1000;
        _evidenceTemp.deliveryWay = LibEvidence.DeliveryWay.sendPay;
        _evidenceTemp.companyNo = LibEvidence.CompanyNo.SHUNFENG;
        _evidenceTemp.status = LibEvidence.Status.INIT;
        user2evidenceIds[_evidenceTemp.userId].push(_evidenceTemp.id);
        evidenceIds.push(_evidenceTemp.id);
        evidenceMap[_evidenceTemp.id] = _evidenceTemp;
        evidenceMap[_evidenceTemp.id].status = LibEvidence.Status.INIT;
        _id = _evidenceTemp.id;
        string memory success = "";
        uint _total = 1;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString(), ",\"items\":[");
        success = success.concat("{\"id\":", _id.toString(),  "}]}}");
        Notify(0, success);
    }

    function applyEvidence(string _invoiceJson) {
        _evidenceTemp.reset();
        if(!_evidenceTemp.fromJson(_invoiceJson)){
            LibLog.log("applyEvidence json invalid");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"json invalid");
            return;
        }
        //通过获取到的businessNo获取证明文件id
        _evidenceTemp.id = findIdByBusinessNo(_evidenceTemp.businessNo);
        if(evidenceMap[_evidenceTemp.id].id == 0){
            LibLog.log("applyEvidence evidence id not exist");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.ID_NOT_EXIST),"evidence id not exist");
            return;
        }
        if(evidenceMap[_evidenceTemp.id].status!=LibEvidence.Status.INIT){
            LibLog.log("applyEvidence evidence status error");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.STATUS_ERROR),"evidence status error");
            return;
        }
        if(_evidenceTemp.receiver.equals("")){
            LibLog.log("applyEvidence receiver can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.RECEIVERUNIT_EMPTY),"receiver can not empty");
            return;
        }
        if(_evidenceTemp.mobile.equals("")){
            LibLog.log("applyEvidence mobile can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.MOBILE_EMPTY),"mobile can not empty");
            return;
        }
        // if(_evidenceTemp.receiverUnit.equals("")){
        //     Notify(errno_prefix+uint(Error.RECEIVERUNIT_EMPTY),"receiverUnit can not empty");
        //     return;
        // }
        if(_evidenceTemp.detailAddress.equals("")){
            LibLog.log("applyEvidence detailAddress can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.DETAILADDRESS_EMPTY),"detailAddress can not empty");
            return;
        }
        if(_evidenceTemp.deliveryWay == LibEvidence.DeliveryWay.NONE||uint(_evidenceTemp.deliveryWay)>uint(LibEvidence.DeliveryWay.receivePay)){
            LibLog.log("applyEvidence deliveryWay incorrect");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.RECEIVERUNIT_EMPTY),"deliveryWay incorrect");
            return;
        }
        if(_evidenceTemp.companyNo == LibEvidence.CompanyNo.NONE||uint(_evidenceTemp.companyNo)>uint(LibEvidence.CompanyNo.SHUNFENG)){
            LibLog.log("applyEvidence companyNo incorrect");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.COMPANYNO_INCORRECT),"companyNo incorrect");
            return;
        }
        //如果传业务流水号
        if(!_evidenceTemp.businessNo.equals("")){
            evidenceMap[_evidenceTemp.id].businessNo = _evidenceTemp.businessNo;
        }
        evidenceMap[_evidenceTemp.id].evidenceDate = now*1000;
        evidenceMap[_evidenceTemp.id].receiver = _evidenceTemp.receiver;
        evidenceMap[_evidenceTemp.id].mobile = _evidenceTemp.mobile;
        evidenceMap[_evidenceTemp.id].receiverUnit = _evidenceTemp.receiverUnit;
        evidenceMap[_evidenceTemp.id].detailAddress = _evidenceTemp.detailAddress;
        evidenceMap[_evidenceTemp.id].postCode = _evidenceTemp.postCode;
        evidenceMap[_evidenceTemp.id].deliveryWay = _evidenceTemp.deliveryWay;
        evidenceMap[_evidenceTemp.id].companyNo = _evidenceTemp.companyNo;
        evidenceMap[_evidenceTemp.id].status = LibEvidence.Status.APPLIED;
        string memory _json = "{";
        _json = _json.jsonCat("id", evidenceMap[_evidenceTemp.id].secPledgeId);
        _json = _json.jsonCat("isEvidenceApplied", uint(LibSecPledge.IsEvidenceApplied.YES));
        _json = _json.concat("}");
        secpledgemanager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager","0.0.1.0"));
        secpledgemanager.updateSecPledge(_json);
        Notify(0,"success");
    }

    // 专门为券商做的申请证明文件
    function applyEvidenceOfBroker(string _invoiceJson) {
        _evidenceTemp.reset();
        if(!_evidenceTemp.fromJson(_invoiceJson)){
            LibLog.log("applyEvidence json invalid");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.BAD_PARAMETER),"json invalid");
            return;
        }
        //通过获取到的businessNo获取证明文件id
        _evidenceTemp.id = findIdByBusinessNo(_evidenceTemp.businessNo);
        if(evidenceMap[_evidenceTemp.id].id == 0){
            LibLog.log("applyEvidence evidence id not exist");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.ID_NOT_EXIST),"evidence id not exist");
            return;
        }
        if(evidenceMap[_evidenceTemp.id].status!=LibEvidence.Status.INIT){
            LibLog.log("applyEvidence evidence status error");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.STATUS_ERROR),"evidence status error");
            return;
        }
        if(_evidenceTemp.receiver.equals("")){
            LibLog.log("applyEvidence receiver can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.RECEIVERUNIT_EMPTY),"receiver can not empty");
            return;
        }
        if(_evidenceTemp.mobile.equals("")){
            LibLog.log("applyEvidence mobile can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.MOBILE_EMPTY),"mobile can not empty");
            return;
        }
        // if(_evidenceTemp.receiverUnit.equals("")){
        //     Notify(errno_prefix+uint(Error.RECEIVERUNIT_EMPTY),"receiverUnit can not empty");
        //     return;
        // }
        if(_evidenceTemp.detailAddress.equals("")){
            LibLog.log("applyEvidence detailAddress can not empty");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.DETAILADDRESS_EMPTY),"detailAddress can not empty");
            return;
        }
        if(_evidenceTemp.deliveryWay == LibEvidence.DeliveryWay.NONE||uint(_evidenceTemp.deliveryWay)>uint(LibEvidence.DeliveryWay.receivePay)){
            LibLog.log("applyEvidence deliveryWay incorrect");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.RECEIVERUNIT_EMPTY),"deliveryWay incorrect");
            return;
        }
        if(_evidenceTemp.companyNo == LibEvidence.CompanyNo.NONE||uint(_evidenceTemp.companyNo)>uint(LibEvidence.CompanyNo.SHUNFENG)){
            LibLog.log("applyEvidence companyNo incorrect");
            LibLog.log("json: ",_invoiceJson);
            Notify(errno_prefix+uint(Error.COMPANYNO_INCORRECT),"companyNo incorrect");
            return;
        }
        //如果传业务流水号
        if(!_evidenceTemp.businessNo.equals("")){
            evidenceMap[_evidenceTemp.id].businessNo = _evidenceTemp.businessNo;
        }
        evidenceMap[_evidenceTemp.id].evidenceDate = now*1000;
        evidenceMap[_evidenceTemp.id].receiver = _evidenceTemp.receiver;
        evidenceMap[_evidenceTemp.id].mobile = _evidenceTemp.mobile;
        evidenceMap[_evidenceTemp.id].receiverUnit = _evidenceTemp.receiverUnit;
        evidenceMap[_evidenceTemp.id].detailAddress = _evidenceTemp.detailAddress;
        evidenceMap[_evidenceTemp.id].postCode = _evidenceTemp.postCode;
        evidenceMap[_evidenceTemp.id].deliveryWay = _evidenceTemp.deliveryWay;
        evidenceMap[_evidenceTemp.id].companyNo = _evidenceTemp.companyNo;
        evidenceMap[_evidenceTemp.id].status = LibEvidence.Status.APPLIED;
        // string memory _json = "{";
        // _json = _json.jsonCat("id", evidenceMap[_evidenceTemp.id].secPledgeId);
        // _json = _json.jsonCat("isEvidenceApplied", uint(LibSecPledge.IsEvidenceApplied.YES));
        // _json = _json.concat("}");
        // secpledgemanager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager","0.0.1.0"));
        // secpledgemanager.updateSecPledge(_json);
        Notify(0,"success");
    }

    //质物状态里通过查询businessNo获取证明文件id;
    function findIdByBusinessNo(string _businessNo) public returns(uint _id){
        for(uint i = 0; i< evidenceIds.length;i++){
            if(evidenceMap[evidenceIds[i]].businessNo.equals(_businessNo)){
                _id = evidenceMap[evidenceIds[i]].id;
                return;
            }
        }
    }
    // 通过BusinessNo查找证明文件详情 已经确认证明文件的BusinessNo唯一
    // 2017.09.22 java端之前写了这个方法合约没有
    function findByBusinessNo(string _businessNo) constant public returns(string _ret){
        LibLog.log(_businessNo);
        uint _id = findIdByBusinessNo(_businessNo);
        _ret = findById(_id);
    }

    //根据发票id添加运单号
    function addDeliveryNo(uint _id, string _deliveryNo) {
        if(evidenceMap[_id].id==0){
            LibLog.log("id not exist");
            Notify(errno_prefix+uint(Error.ID_NOT_EXIST),"id not exist");
            return;
        }
        if(evidenceMap[_id].status!=LibEvidence.Status.APPLIED){
            LibLog.log("status must be applied");
            Notify(errno_prefix+uint(Error.STATUS_ERROR),"status must be applied");
            return;
        }
        if(_deliveryNo.equals("")){
            LibLog.log("deliveryNo can not empty");
            Notify(errno_prefix+uint(Error.DELIVERYNO_EMPTY),"deliveryNo can not empty");
            return;
        }
        evidenceMap[_id].deliveryNo = _deliveryNo;
        evidenceMap[_id].status = LibEvidence.Status.MAILED;

        //更新质物状态isEvidenceMailed
        string memory _json = "{";
        _json = _json.jsonCat("id", evidenceMap[_id].secPledgeId);
        _json = _json.jsonCat("isEvidenceMailed", uint(LibSecPledge.IsEvidenceMailed.YES));
        _json = _json.concat("}");
        secpledgemanager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager","0.0.1.0"));
        secpledgemanager.updateSecPledge(_json);
        Notify(0,"success");
    }

    //专门为券商做的添加运单号
    function addDeliveryNoOfBroker(uint _id, string _deliveryNo) {
        if(evidenceMap[_id].id==0){
            LibLog.log("id not exist");
            Notify(errno_prefix+uint(Error.ID_NOT_EXIST),"id not exist");
            return;
        }
        if(evidenceMap[_id].status!=LibEvidence.Status.APPLIED){
            LibLog.log("status must be applied");
            Notify(errno_prefix+uint(Error.STATUS_ERROR),"status must be applied");
            return;
        }
        if(_deliveryNo.equals("")){
            LibLog.log("deliveryNo can not empty");
            Notify(errno_prefix+uint(Error.DELIVERYNO_EMPTY),"deliveryNo can not empty");
            return;
        }
        evidenceMap[_id].deliveryNo = _deliveryNo;
        evidenceMap[_id].status = LibEvidence.Status.MAILED;

        //更新质物状态isEvidenceMailed
        // string memory _json = "{";
        // _json = _json.jsonCat("id", evidenceMap[_evidenceTemp.id].secPledgeId);
        // _json = _json.jsonCat("isEvidenceMailed", uint(LibSecPledge.IsEvidenceMailed.YES));
        // _json = _json.concat("}");
        // secpledgemanager = SecPledgeManager(rm.getContractAddress("CsdcModule", "0.0.1.0", "SecPledgeManager","0.0.1.0"));
        // secpledgemanager.updateSecPledge(_json);
        Notify(0,"success");
    }

    //根据业务流水号添加运单号
    function addDeliveryNoByBusinessNo(string _businessNo, string _deliveryNo) {
        uint _id = findIdByBusinessNo(_businessNo);
        addDeliveryNo(_id,_deliveryNo);
    }

    function findById(uint _id) constant public returns(string _ret) {
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": []}}";
        if(evidenceMap[_id].id==0){
            return;
        }
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": [";
        _ret = _ret.concat(evidenceMap[_id].toJson());
        _ret = _ret.concat("]}}");
    }

    LibEvidence.Cond _cond;

    //个人用户根据状态分页查询 日期 业务单号 //发票类型 发票抬头 联系人
    function pageByUserId(string _json) constant public returns(string _ret){
        _cond.fromJson(_json);

        if (user2evidenceIds[_cond.userId].length <= 0) {
            return LibStack.popex(itemStackPush("", 0));
        }
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= user2evidenceIds[_cond.userId].length) {
            return LibStack.popex(itemStackPush("", 0));
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for(uint i = 0; i< user2evidenceIds[_cond.userId].length; i++){
            _evidenceTemp = evidenceMap[user2evidenceIds[_cond.userId][i]];

            //选择不包含并且该发票的状态为init;  
            if(_cond.containInit == 1 && _evidenceTemp.status == LibEvidence.Status.INIT) {
                continue;
            }

            if (_json.jsonKeyExists("initDate") && _cond.initDate != 0 && _cond.initDate != _evidenceTemp.initDate){
                continue;
            }
            if (_json.jsonKeyExists("minInitDate") && _json.jsonKeyExists("maxInitDate") && _cond.maxInitDate !=0 
                &&(_evidenceTemp.initDate < _cond.minInitDate || _evidenceTemp.initDate > _cond.maxInitDate)){
                continue;
            }
            if (_json.jsonKeyExists("evidenceDate") && _cond.evidenceDate!=0 && _cond.evidenceDate!=_evidenceTemp.evidenceDate) {
                continue;
            }
            if (_json.jsonKeyExists("evidenceType") && _cond.evidenceType != LibEvidence.EvidenceType.NONE&&uint(_cond.evidenceType) != uint(_evidenceTemp.evidenceType)){
                continue;
            }
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(_evidenceTemp.businessNo)) {
                continue;
            }  
            if (_json.jsonKeyExists("receiver") && !_cond.receiver.equals("") && !_cond.receiver.equals(_evidenceTemp.receiver)) {
                continue;
            }
            if (_json.jsonKeyExists("status") && _cond.status != LibEvidence.Status.NONE&&uint(_cond.status) != uint(_evidenceTemp.status)){
                continue;
            }
            if (_json.jsonKeyExists("handleChannel") && _cond.handleChannel !=0 && _cond.handleChannel != uint(_evidenceTemp.handleChannel)){
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
                len = LibStack.append(_evidenceTemp.toJson());
            }
        }
        LibJson.pop();
        return LibStack.popex(itemStackPush(LibStack.popex(len), _count));
    }

    //中证登用户发票查询 日流水 业务单号 收件人姓名
    function pageByCond(string _json) constant public returns(string _ret){
        if (evidenceIds.length <= 0) {
            return LibStack.popex(itemStackPush("", 0));
        }

        _cond.fromJson(_json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= evidenceIds.length) {
            return LibStack.popex(itemStackPush("", 0));
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for(uint i = 0; i< evidenceIds.length; i++){
            _evidenceTemp = evidenceMap[evidenceIds[i]];

            if (_json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != _evidenceTemp.userId) {
                continue;
            }
            //选择不包含并且该发票的状态为init;  
            if ( _cond.containInit == 1 && _evidenceTemp.status == LibEvidence.Status.INIT) {
                continue;
            }

            if (_json.jsonKeyExists("initDate") && _cond.initDate != 0 && _cond.initDate != _evidenceTemp.initDate){
                continue;
            }
            if (_json.jsonKeyExists("minInitDate") && _json.jsonKeyExists("maxInitDate") && _cond.maxInitDate !=0 
                &&(_evidenceTemp.initDate < _cond.minInitDate || _evidenceTemp.initDate > _cond.maxInitDate)){
                continue;
            }
            if (_json.jsonKeyExists("evidenceDate") && _cond.evidenceDate!=0 && _cond.evidenceDate!=_evidenceTemp.evidenceDate) {
                continue;
            }
            if (_json.jsonKeyExists("evidenceType") && _cond.evidenceType != LibEvidence.EvidenceType.NONE&&uint(_cond.evidenceType) != uint(_evidenceTemp.evidenceType)){
                continue;
            }
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(_evidenceTemp.businessNo)) {
                continue;
            }  
            if (_json.jsonKeyExists("receiver") && !_cond.receiver.equals("") && !_cond.receiver.equals(_evidenceTemp.receiver)) {
                continue;
            }
            if (_json.jsonKeyExists("status") && _cond.status != LibEvidence.Status.NONE&&uint(_cond.status) != uint(_evidenceTemp.status)){
                continue;
            }
            if (_json.jsonKeyExists("handleChannel") && _cond.handleChannel !=0 && _cond.handleChannel != uint(_evidenceTemp.handleChannel)){
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
                len = LibStack.append(_evidenceTemp.toJson());
            }
        }
        LibJson.pop();
        return LibStack.popex(itemStackPush(LibStack.popex(len), _count));
    }

    function listAll() constant public returns(string _ret){
        _ret = "{ \"ret\": 0, \"message\": \"success\", \"data\": { \"total\": 0, \"items\": [] }}";
        if(evidenceIds.length<=0){
            return;
        }
        _ret = "{\"ret\":0, \"message\": \"success\", \"data\":{";
        _ret = _ret.concat(uint(evidenceIds.length).toKeyValue("total"), ",\"items\":[");
        for(uint i = 0; i < evidenceIds.length; i++){
            if(i==evidenceIds.length-1){
                _ret = _ret.concat(evidenceMap[evidenceIds[i]].toJson());
            }else{
                _ret = _ret.concat(evidenceMap[evidenceIds[i]].toJson(),",");
            }
        }
        _ret = _ret.concat("]}}");
        Notify(0,_ret);
    }

    function deleteById(uint id) {
        delete tempList;
        for (uint i = 0; i < evidenceIds.length; i++)  {
            if(id == evidenceIds[i]) {
                continue;
            }
            else {
                tempList.push(evidenceIds[i]);
            }
        }
        delete evidenceIds;
        for (uint j = 0; j < tempList.length; ++j) {
            evidenceIds.push(tempList[j]);
        }
        delete evidenceMap[id];
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

    // function updateIsEvidenceApplied(uint id, uint isEvidenceApplied){}
    // function updateIsEvidenceMailed(uint id, uint isEvidenceMailed){}

}