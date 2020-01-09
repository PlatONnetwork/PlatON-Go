pragma solidity ^0.4.12;

import "./BrokerDao.sol";

contract BrokerManager is OwnerNamed {

    using LibInt for *;
    using LibString for *;
    using LibBroker for *;

    LibBroker.Broker t_broker;
    BrokerDao brokerDao;
    Sequence sq;

    enum BrokerError {
      NONE,
      BAD_PARAMETER,
      DAO_ERROR
    }

    function BrokerManager() {
        register("CsdcModule", "0.0.1.0", "BrokerManager", "0.0.1.0");
    }

    function insert(string _json) {

        if (!t_broker.fromJson(_json)) {
            Notify(uint(BrokerError.BAD_PARAMETER), "json invalid");
            return;
        }

        sq = Sequence(rm.getContractAddress("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0"));
        t_broker.id = sq.getSeqNo("Broker.id");

        brokerDao = BrokerDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "BrokerDao", "0.0.1.0"));
        uint _errno = brokerDao.insert_Broker(t_broker.toJson());
        if(_errno != 0) {
            LibLog.log("call dao error: ", _errno.toString());
            Notify(uint(BrokerError.DAO_ERROR), "insert broker failed");
            return;
        }
        string memory success = "";
        success = success.concat("{\"ret\":0,\"message\": \"success\", \"data\":{\"total\": 1");
        success = success.concat(",\"items\":[]");
        success = success.concat(",\"id\":", t_broker.id.toString(),  "}}");
        Notify(0, success);
    }

    function findById(uint _id) constant returns(string _ret) {
        brokerDao = BrokerDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "BrokerDao", "0.0.1.0"));
        uint len = brokerDao.select_Broker_byId(_id);
        return LibStack.popex(len); 
    }

    function pageBroker(string _json) constant returns(string _ret) {
        brokerDao = BrokerDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "BrokerDao", "0.0.1.0"));
        uint len = brokerDao.pageBroker(_json);
        return LibStack.popex(len);
    }

    function listAll() constant returns(string) {
        brokerDao = BrokerDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "BrokerDao", "0.0.1.0"));
        uint len = brokerDao.select_Broker_all();
        return LibStack.popex(len);
    }

    function deleteBrokerbyId(uint _id) {
        brokerDao = BrokerDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "BrokerDao", "0.0.1.0"));
        uint _errno = brokerDao.delete_Broker_byId(_id);
        if(_errno != 0) {
            LibLog.log("call dao error: ", _errno.toString());
            Notify(uint(BrokerError.DAO_ERROR), "delete broker failed");
        } else {
            Notify(0, "success");
        }
    }

    function updateBroker(string _json) {
        brokerDao = BrokerDao(rm.getContractAddress("CsdcModule", "0.0.1.0", "BrokerDao", "0.0.1.0"));
        uint _errno = brokerDao.update_Broker(_json);
        if(_errno != 0) {
            LibLog.log("call dao error: ", _errno.toString());
            Notify(uint(BrokerError.DAO_ERROR), "update broker failed");
        } else {
            Notify(0, "success");
        }
    }

    event Notify(uint _errorno, string _info);
}
