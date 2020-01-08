pragma solidity ^0.4.12;

import "./csdc_base/CsdcBaseInterface.sol";
import "./csdc_base/CommonContract.sol";

import "./Sequence.sol";

contract SecPledgeManager is CommonContract {
    using LibSecPledge for *;
    using LibString for *;
    using LibInt for *;
    using LibBiz for *;
    using LibPledgeSecurity for *;

     //inner setting member
    LibSecPledge.SecPledge internal tmp_SecPledge;
    LibPledgeSecurity.PledgeSecurity _tmpSec;

    Sequence sq;
    CsdcBaseInterface bi;

    modifier getSq() { 
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0")); 
        _;
    }

    function SecPledgeManager() {
        register("CsdcModule", "0.0.1.0", 'SecPledgeManager', '0.0.1.0');
    }

        /** @brief errno for test case */
    enum SecPledgeError {
        NO_ERROR,
        BAD_PARAMETER,
        OPERATE_NOT_ALLOWED,
        ID_EMPTY,
        DAO_ERROR   
    }

    uint errno_prefix = 10000;

    function updateSecPledge(string _json) getOrderDao returns (uint) {
        LibLog.log("updateSecPledge: ", _json);
        if (od.update_SecPledge(_json) != 0) {
          Notify(errno_prefix + uint(SecPledgeError.DAO_ERROR), "call dao error");
          return;
        }
        Notify(0, "success");
    }

    //分页功能
    function pageByCond(string _json) getOrderDao constant returns (string) {
        uint len = od.pageByCond_SecPledge(_json);
        return LibStack.popex(len);
    }

    //查找可以解除質押的成交記錄
    function findToDis(string _json) getOrderDao constant returns (string) {
        uint len = od.findToDis_SecPledge_Online(_json);
        return LibStack.popex(len);
    }

    //安卓调用版本
    function findById(uint id) getOrderDao constant returns (string _ret) {
        uint len = od.select_SecPledge_byId(id);
        return LibStack.popex(len);
    }

    //确定该质押记录及质押人质权人是否匹配
    function isMatching(uint id, address pledgorId, address pledgeeId) constant returns (bool _ret) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        return tmp_SecPledge.pledgorId == pledgorId && tmp_SecPledge.pledgeeId == pledgeeId;
    }

    //请求解押的证券是否合规, 输入的json为待解除的证券及其数量
    function isMatchingAppliedSecurities(uint id, string json) constant returns (bool) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        _tmpSec.fromJson(json);

        for (uint i = 0; i < tmp_SecPledge.appliedSecurities.length; i++) {
            if (_tmpSec.id == tmp_SecPledge.appliedSecurities[i].id) {
                if ( tmp_SecPledge.appliedSecurities[i].secAccount.equals(_tmpSec.secAccount)
                && tmp_SecPledge.appliedSecurities[i].secCode.equals(_tmpSec.secCode)
                && tmp_SecPledge.appliedSecurities[i].secName.equals(_tmpSec.secName)
                && tmp_SecPledge.appliedSecurities[i].secType.equals(_tmpSec.secType) 
                && tmp_SecPledge.appliedSecurities[i].hostedUnit.equals(_tmpSec.hostedUnit)
                && tmp_SecPledge.appliedSecurities[i].hostedUnitName.equals(_tmpSec.hostedUnitName)
                && tmp_SecPledge.appliedSecurities[i].secProperty.equals(_tmpSec.secProperty)
                ) {
                //计算红股+剩余股数是否大等于申请的解冻股数量
                    LibPledgeSecurity.PledgeSecurity p = tmp_SecPledge.appliedSecurities[i];
                    return (p.bonusShareAmount + p.remainPledgeNum)>=_tmpSec.pledgeNum;
                }
            }
         }
    }

     //获得某个证券的红利数量
    function getAppliedSecuritiesProfitAmount(uint id, string json) constant returns (uint _ret) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        _tmpSec.fromJson(json);

        for (uint i = 0; i < tmp_SecPledge.appliedSecurities.length; i++) {
            if (_tmpSec.id == tmp_SecPledge.appliedSecurities[i].id) {
                //获取红利数量
                return tmp_SecPledge.appliedSecurities[i].profitAmount;
            }
        }
    }

    function getAppliedSecuritiesFreezeNo(uint id, string json) constant returns (uint _ret) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        _tmpSec.fromJson(json);

        for (uint i = 0; i < tmp_SecPledge.appliedSecurities.length; i++) {
            if (_tmpSec.id == tmp_SecPledge.appliedSecurities[i].id) {
                //获取红利数量
                return tmp_SecPledge.appliedSecurities[i].freezeNo.storageToUint();
            }
        }
    }

    function getAppliedSecuritiesSubFreezeNo(uint id, string json) constant returns (uint _ret) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        _tmpSec.fromJson(json);

        for (uint i = 0; i < tmp_SecPledge.appliedSecurities.length; i++) {
            if (_tmpSec.id == tmp_SecPledge.appliedSecurities[i].id) {
                //获取红利数量
                return tmp_SecPledge.appliedSecurities[i].subFreezeNo.storageToUint();
            }
        }
    }

    //用于检查状态
    function isStatus(uint id, uint _status) returns (bool) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        return tmp_SecPledge.status==LibSecPledge.PledgeStatus(_status);
    }

    //用于检查证明文件是否已申请
    function isEvidenceApplied(uint id, uint _status) returns(bool){
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        return tmp_SecPledge.isEvidenceApplied == uint(LibSecPledge.IsEvidenceApplied(_status));
    }

    //用于检查证明文件是否已邮寄
    function isEvidenceMailed(uint id, uint _status) returns(bool){
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        return tmp_SecPledge.isEvidenceMailed == uint(LibSecPledge.IsEvidenceMailed(_status));
    }

    //新增解押申请时调用，会修改成功记录的状态为处理中，并记录解押申请单
    function addNewDisSecPledgeAplly(uint id, uint disSecPedgeApplyId) getOrderDao returns (bool _ret) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        tmp_SecPledge.status = LibSecPledge.PledgeStatus.DISPLEGING;
        tmp_SecPledge.disSecPedgeApplyIds.push(disSecPedgeApplyId);

        od.update_SecPledge(tmp_SecPledge.toJson());
    }

    //回退质押，即将状态置回
    function undoDisSecPledgeApply(uint id) getOrderDao returns (bool _ret) {
        string memory _json = getSecPledgeById(id);
        if(_json.equals("")) {
            LibLog.log("secPledge not exists");
            return;
        }
        tmp_SecPledge.fromJson(_json);
        tmp_SecPledge.status = tmp_SecPledge.statusShow;

        od.update_SecPledge(tmp_SecPledge.toJson());
    }

    event Notify(uint _errorno, string _info);    
}