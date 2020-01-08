pragma solidity ^0.4.12;

import "./Sequence.sol";
import "./OrderDao.sol";

contract NewSecPledgeManager  is OwnerNamed  {
    using LibDisSecPledgeApply for *;
    using LibString for *;
    using LibInt for *;
    using LibBiz for *;
    using LibSecPledge for *;
    using LibJson for *;

    OrderDao orderDao;
    Sequence sq;
    //inner setting member
    LibDisSecPledgeApply.DisSecPledgeApply internal tmp_DisSecPledgeApply;
    LibSecPledge.SecPledge internal tmp_SecPledge;
    
    LibBiz.Biz internal tmp_Biz;

    event Notify(uint _errorno, string _info);
    
    /** @brief errno for test case */
    enum SecPledgeError {
        NO_ERROR,
        BAD_PARAMETER,
        DAO_ERROR,
        OPERATE_NOT_ALLOWED,
        ID_EMPTY,   
        PAYERTYPE_EMPTY,
        INSERT_FAILED,
        USER_STATUS_ERROR,
        APPLIEDSECURITIES_ERROR,
        PLEDGOR_PLEDGEE_SAME,
        BUSINESSNO_ERROR //业务编号获取失败
    }

    enum OperateCode{
        NONE, 
        PASS, //通过
        REJECT, //拒绝
        WAIT //等待处理
    }
    uint errno_prefix = 10000;

    function NewSecPledgeManager() {
        register("CsdcModule", "0.0.1.0", "NewSecPledgeManager", "0.0.1.0");
        
        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
        orderDao = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
    }

    modifier getOrderDao(){ 
      orderDao = OrderDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0"));
      _;
    }

    //20.1.根据id主键查询
    function findById(uint _id) getOrderDao constant returns(string _ret) {
        uint len = orderDao.select_SecPledge_byId(_id);
        return LibStack.popex(len);
    }

    function insertSecPledge(string _json) getOrderDao {
        LibLog.log("insert: ", _json);
        if(!tmp_SecPledge.fromJson(_json)) {
            LibLog.log("json invalid");
            Notify(errno_prefix + uint(SecPledgeError.BAD_PARAMETER), "json invalid");
            return;
        }
        if (orderDao.insert_SecPledge(tmp_SecPledge.toJson()) != 0) {
          Notify(errno_prefix + uint(SecPledgeError.DAO_ERROR), "call dao error");
          return;
        }
        Notify(0, "success");
    }


    //20.2.更新质物记录信息
    function updateSecPledge(string _json) getOrderDao returns (uint) {
        if (orderDao.update_SecPledge(_json) != 0) {
          Notify(errno_prefix + uint(SecPledgeError.DAO_ERROR), "call dao error");
        }
        Notify(0, "success");
    }
    //20.3.质物记录保存
    function saveSecPledge(string _json) getOrderDao returns(bool _ret) {
        _saveSecPledge(_json);
    }
  
    //20.4.撤销质物记录
    function cancelSecPledge(uint id) getOrderDao returns(bool _ret) {
        uint daoResult = orderDao.delete_SecPledge_byId(id);
        notify(daoResult);
    }

    /* 以下是internal接口 */

    //用于保存质物记录的内部方法
    function _saveSecPledge(string _json) internal returns(bool _ret) {
        LibLog.log("_saveSecPledge: ", _json);


        //封装apply       
        tmp_SecPledge.reset();       
        tmp_SecPledge.fromJson(_json); 
       
        //分配一致的id，如id已存在，则选择更新
        uint id = 0;
        if (tmp_SecPledge.id != 0) {
            id = tmp_SecPledge.id;
        } else {
            id  = sq.getSeqNo("Biz.id");
        }
        tmp_SecPledge.id = id;

        if (__isExistSecPledge(id)) { //更新数据
            //申请入库
            uint daoResult = orderDao.update_SecPledge(_json);
            LibLog.log("_saveSecPledge: update_SecPledge: ", daoResult);

        } else { //插入新数据
            //申请入库
            orderDao.insert_SecPledge(tmp_SecPledge.toJson());
        }


        string memory success = "";
        uint _total = 1;
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\":", _total.toString() );
        success = success.concat(",\"items\":[]");
        success = success.concat(",\"id\":", id.toString(),  "}}");
        Notify(0, success);
        
        return true;

    }

    //检验某个记录是否存在
    function __isExistSecPledge(uint id) internal returns(bool _ret) {
        _ret = false;
        string memory _json = findById(id);
        LibLog.log(_json);

        LibJson.push(_json);
        if (_json.jsonKeyExists("data.items[0]")) {
            LibLog.log("key exists");
            _ret = true;
        }
        LibJson.pop();
    }
    
    function notify(uint result) internal {
        if (result == 0) {
            Notify(result, "success");
        } else {
            Notify(result, "error");
        }
    }

    function listAll() getOrderDao constant returns(string) {
        uint len = orderDao.select_SecPledge_all();
        return LibStack.popex(len);
    }

    //分页功能
    function pageByCond(string _json) getOrderDao constant returns (string) {
        uint len = orderDao.pageByCond_SecPledge(_json);
        return LibStack.popex(len);
    }
    
}
