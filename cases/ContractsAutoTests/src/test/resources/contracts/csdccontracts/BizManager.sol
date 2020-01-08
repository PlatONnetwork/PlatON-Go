pragma solidity ^0.4.12;

import "./csdc_base/CsdcBaseInterface.sol";
import "./csdc_base/CommonContract.sol";
import "./Sequence.sol";

contract BizManager is CommonContract, CsdcBaseInterface {
    using LibInt for *;
    using LibString for *;
    using LibBiz for *;
    
    CsdcBaseInterface bi;
    
    LibBiz.Biz _tmpBiz;

    enum BizError {
      NONE,
      DAO_ERROR
    }

    uint errno_prefix = 18000;

    function BizManager() {
        register("CsdcModule", "0.0.1.0", "BizManager", '0.0.1.0');
    }

    function insertBiz(string _json) getOrderDao {
        LibLog.log("BizManager insertBiz invoked.");
        if (od.insert_Biz(_json) != 0) {
          Notify(errno_prefix + uint(BizError.DAO_ERROR), "call dao error");
          return;
        }
        Notify(0, "success");
    }

    function insert(
        uint _id,
        address _pledgorId,
        address _pledgeeId,
        address _managerId,   
        uint _bizType, 
        uint _status, 
        uint _relatedId 
      ) getOrderDao returns (uint) {

        bi = CsdcBaseInterface(rm.getContractAddress("CsdcModule", "0.0.1.0", "PerUserManager", "0.0.1.0"));
        
        string memory _pledgorName = LibStack.popex(bi.findNameById(_pledgorId));
        string memory _pledgeeName = LibStack.popex(bi.findNameById(_pledgeeId));

        _tmpBiz.create(_pledgorId, _pledgorName, _pledgeeId, _pledgeeName, _managerId, _bizType, _status, _relatedId);
        _tmpBiz.id = _id;
        _tmpBiz.channelType = uint(LibBiz.ChannelType.ONLINE);
        _tmpBiz.createTime = now*1000;
        _tmpBiz.updateTime = now*1000;

        return od.insert_Biz(_tmpBiz.toJson());
    }

    function changeStatus(uint _bizId, uint _status) getOrderDao returns (bool) {
        string memory _json = getBizById(_bizId);
        if(_json.equals("")) {
            LibLog.log("bizId not exists");
            return;
        }
            
        _tmpBiz.fromJson(_json);
        _tmpBiz.status = LibBiz.BizStatus(_status);
        _tmpBiz.updateTime = now*1000;
        od.update_Biz(_tmpBiz.toJson());
        return true;
    }
    
    function endBiz(uint _bizId, uint _status) getOrderDao returns (bool) {
        string memory _json = getBizById(_bizId);
        if(_json.equals("")) {
            LibLog.log("bizId not exists");
            return;
        }
        
        _tmpBiz.fromJson(_json);
        _tmpBiz.status = LibBiz.BizStatus(_status);
        _tmpBiz.updateTime = now*1000;
        _tmpBiz.endTime = now*1000;
        od.update_Biz(_tmpBiz.toJson());
        return true;
    }

    function addPaymentId(uint _bizId, uint _paymentId) getOrderDao returns (bool) {
        string memory _json = getBizById(_bizId);
        if(_json.equals("")) {
            LibLog.log("bizId not exists");
            return;
        }

        _tmpBiz.fromJson(_json);
        _tmpBiz.paymentId = _paymentId;
        od.update_Biz(_tmpBiz.toJson());
        return true;
    }

    function setBusinessNo(uint _bizId, string _businessNo) getOrderDao returns (bool) {
        string memory _json = getBizById(_bizId);
        if(_json.equals("")) {
            LibLog.log("bizId not exists");
            return;
        }

        _tmpBiz.fromJson(_json);
        _tmpBiz.businessNo = _businessNo;
        od.update_Biz(_tmpBiz.toJson());
        return true;
    }

    function setPledgeContractNo(uint _bizId, string _contractNo) getOrderDao returns (bool) {
        string memory _json = getBizById(_bizId);
        if(_json.equals("")) {
            LibLog.log("bizId not exists");
            return;
        }

        _tmpBiz.fromJson(_json);
        _tmpBiz.pledgeContractNo = _contractNo;
        od.update_Biz(_tmpBiz.toJson());
        return true;
    }

    function setPledgeContracFile(uint _bizId, string _fileId, string _fileName) getOrderDao returns (bool) {
        string memory _json = getBizById(_bizId);
        if(_json.equals("")) {
            LibLog.log("bizId not exists");
            return;
        }

        _tmpBiz.fromJson(_json);
        _tmpBiz.pledgeContractFileId = _fileId;
        _tmpBiz.pledgeContractName = _fileName;
        od.update_Biz(_tmpBiz.toJson());
        return true;
    }
    
    function updateBiz(string _json) getOrderDao {
        if (od.update_Biz(_json) != 0) {
          Notify(errno_prefix + uint(BizError.DAO_ERROR), "call dao error");
          return;
        }
        Notify(0, "success");
    }

    function isStatus(uint _bizId, uint _status) constant returns (bool) {
        string memory _json = getBizById(_bizId);
        if(_json.equals("")) {
            LibLog.log("bizId not exists");
            return;
        }
        _tmpBiz.fromJson(_json);
        return _tmpBiz.status==LibBiz.BizStatus(_status);
    }

    function findToDo(string _json) getOrderDao constant returns (string) {
        uint len = od.findToDo_Online(_json);
        return LibStack.popex(len);
    }

  /**
   * @dev 分页查询业务
   * @param _json 查询条件 
   *      bizId-业务单号，userId-用户地址，
   *      minStartTime-开始时间起，maxStartTime-开始时间止，statuses-账户状态list，bizType-业务状态
   *      pageSize-页面大小, pageNo-页面号
   * @return _ret 
   */
    function pageBiz(string _json) getOrderDao constant returns (string) {
        LibLog.log("=======================start====================");
        uint len = od.pageBiz_Online(_json);
        LibLog.log("=======================end====================");
        return LibStack.popex(len);
    }

    function findById(uint _id) getOrderDao constant returns (string) {
        uint len = od.select_Biz_byId(_id);
        return LibStack.popex(len);
    }

    function hasTodo(address _userId) getOrderDao constant returns (bool) {
        return od.hasTodo_Online(_userId);
    }
    
    event Notify(uint _errorno, string _info);

    /* for CsdcBaseInterface */
    function findNameById(address _id) constant returns (uint) {}
}