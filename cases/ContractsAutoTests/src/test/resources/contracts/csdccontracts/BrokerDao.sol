pragma solidity ^0.4.12;
/**
* @file BrokerDao.sol
* @author liaoyan
* @time 2017-07-05
* @desc the definition of BrokerDao contract
*/


import "./csdc_library/LibBroker.sol";
import "./Sequence.sol";

contract BrokerDao is OwnerNamed {

    using LibInt for *;
    using LibString for *;
    using LibBroker for *;
    using LibJson for *;

    mapping(uint => LibBroker.Broker) m_brokerMap; //id => object
    uint[] m_brokerIds;

    LibBroker.Broker t_broker;

    uint MAX_PAGESIZE = 10;
    struct Condition {
        string  orgNo;
        uint pageSize;
        uint pageNo;
    }

    function fromJson_Condition(Condition storage _self, string _json) internal returns(bool succ) {
        reset_Condition(_self);
        LibJson.push(_json);

        if (!_json.isJson()) {
            LibJson.pop();
            return false;
        }

        _self.orgNo = _json.jsonRead("orgNo");
        _self.pageSize = _json.jsonRead("pageSize").toUint();
        _self.pageNo = _json.jsonRead("pageNo").toUint();

        if (_self.pageSize > MAX_PAGESIZE) {
            LibLog.log("pageSize is in excess of 10, now it's set to 10.");
            _self.pageSize = MAX_PAGESIZE;
        }

        LibJson.pop();
        return true;
    }

    function reset_Condition(Condition storage _self) internal {
        delete _self.orgNo;
        delete _self.pageSize;
        delete _self.pageNo;
    }

    Condition _cond;

    enum DaoError {
        NO_ERROR,
        BAD_PARAMETER,
        ID_NOT_EXISTS,
        ID_CONFLICTED
    }

    function BrokerDao() {
        register("CsdcModule", "0.0.1.0", "BrokerDao", "0.0.1.0");
    }

    //Dao for Broker

    function insert_Broker(string _json) returns(uint _ret) {
        if (!t_broker.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);
        
        for (uint i=0; i<m_brokerIds.length; ++i) {
            if (m_brokerIds[i] == t_broker.id)
                return uint(DaoError.ID_CONFLICTED); 
        }

        m_brokerMap[t_broker.id] = t_broker;
        m_brokerIds.push(t_broker.id);

        return uint(DaoError.NO_ERROR);
    }

    function select_Broker_byId(uint _id) constant returns(uint) {
        string memory item = "";
        if(m_brokerMap[_id].id != 0) {
            item = m_brokerMap[_id].toJson();
        }

        return itemStackPush(item, m_brokerIds.length);
    }

    function select_Broker_all() constant returns(uint _ret) {
        uint len = 0;
        len = LibStack.push("");
        for (uint i=0; i<m_brokerIds.length; ++i) {
            if (i > 0) {
                len = LibStack.append(",");
            }
            len = LibStack.append(m_brokerMap[m_brokerIds[i]].toJson());
        }

        return itemStackPush(LibStack.popex(len), m_brokerIds.length);
    }

    function delete_Broker_byId(uint _id) returns(uint _ret) {
        bool found = false;
        for (uint i=0; i<m_brokerIds.length; ++i) {
            if (!found) {
                if (m_brokerIds[i] == _id) {
                    found = true;
                }
            }

            if (found && i < m_brokerIds.length-1) {
                m_brokerIds[i] = m_brokerIds[i+1];
            }
        }

        if (!found)
            return uint(DaoError.ID_NOT_EXISTS);

        m_brokerIds.length -= 1;
        m_brokerMap[_id].reset();

        return uint(DaoError.NO_ERROR);
    }

    function update_Broker(string _json) returns(uint _ret) {
        if (!t_broker.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        if (m_brokerMap[t_broker.id].id == 0) {
            return uint(DaoError.ID_NOT_EXISTS);
        }

        if (m_brokerMap[t_broker.id].update(_json))
            return uint(DaoError.NO_ERROR);
        else
            return uint(DaoError.BAD_PARAMETER);
    }

    function pageBroker(string _json) constant returns(uint) {
        if (m_brokerIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_brokerIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = m_brokerIds.length; i >= 1 ; i--) {
            t_broker = m_brokerMap[m_brokerIds[i-1]];

            if (_json.jsonKeyExists("orgNo") && !_cond.orgNo.equals("") && !_cond.orgNo.equals(t_broker.orgNo)) {
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
                len = LibStack.append(t_broker.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
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
}
