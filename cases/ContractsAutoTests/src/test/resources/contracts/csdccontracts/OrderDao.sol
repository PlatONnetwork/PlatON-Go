pragma solidity ^0.4.12;
/**
* @file OrderDao.sol
* @author liaoyan, huanggaofeng
* @time 2017-07-01
* @desc the definition of OrderDao contract
*/


import "./utillib/LibLog.sol";
import "./sysbase/OwnerNamed.sol";

import "./csdc_library/LibAudit.sol";
import "./csdc_library/LibAttachInfo.sol";
import "./csdc_library/LibBrokerUser.sol";
import "./csdc_library/LibSecPledgeApply.sol";
import "./csdc_library/LibDisSecPledgeApply.sol";
import "./csdc_library/LibSecPledge.sol";
import "./csdc_library/LibBiz.sol";
import "./csdc_library/LibPayment.sol";

contract OrderDao is OwnerNamed {

    using LibInt for *;
    using LibString for *;
    using LibSecPledgeApply for *;
    using LibDisSecPledgeApply for *;
    using LibSecPledge for *;
    using LibBiz for *;
    using LibAudit for *;
    using LibJson for *;
    using LibAttachInfo for *;
    using LibPayment for *;

    mapping(uint => LibSecPledgeApply.SecPledgeApply) m_secPledgeApplyMap; //id => object
    uint[] m_secPledgeApplyIds;

    mapping(uint => LibDisSecPledgeApply.DisSecPledgeApply) m_disSecPledgeApplyMap; //id => object
    uint[] m_disSecPledgeApplyIds;

    mapping(uint => LibSecPledge.SecPledge) m_secPledgeMap; //id => object
    uint[] m_secPledgeIds;

    mapping(uint => LibBiz.Biz) m_bizMap; //id => object
    uint[] m_bizIds;

    mapping(uint => LibPayment.Payment) m_paymentMap; //id => object
    uint[] m_paymentIds;

    mapping(uint=>mapping(uint=>address)) m_operators;

    LibSecPledgeApply.SecPledgeApply t_secPledgeApply;
    LibDisSecPledgeApply.DisSecPledgeApply t_disSecPledgeApply;
    LibSecPledge.SecPledge t_secPledge;
    LibBiz.Biz t_biz;
    LibAudit.Audit t_audit;
    LibAttachInfo.AttachInfo t_attchInfo;
    LibPayment.Payment t_payment;

    uint MAX_PAGESIZE = 500;

    enum QueryScope {
        NONE,
        ALL,        //1-所有
        MY_TODO,    //2-我待办
        MY_PASSBY   //3-我经办
    }

    struct Condition {
        string  businessNo;
        address userId;
        address managerId;
        uint brokerId;
        uint role;
        uint minStartTime;
        uint maxStartTime;
        uint minUpdateTime;
        uint maxUpdateTime;
        uint scope;
        uint[] statuses;
        uint[] checkStatuses;
        uint bizType;
        string desc;
        uint channelType;

        uint bizId;
        string pledgeeName;
        string pledgorName;

        uint secPledgeId;
        string pledgeRegisterNo;
        string pledgeRegisterFileId;
        uint minPledgeTime;
        uint maxPledgeTime;
        uint isOnline;
        uint isEvidenceMailed;
        uint isEvidenceApplied;

        uint[] paymentStatuses;     //付款状态
        uint[] refundApplyStatuses; //退款申请状态
        string payerName;           //付款人名称
        uint paymentType;           //付款方式
        uint minPaymentTime;        //付款时间起
        uint maxPaymentTime;        //付款时间止
        uint minRefundApplyTime;    //退款申请时间起
        uint maxRefundApplyTime;    //退款申请时间止

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

        _self.businessNo = _json.jsonRead("businessNo");
        _self.userId = _json.jsonRead("userId").toAddress();
        _self.managerId = _json.jsonRead("managerId").toAddress();
        _self.brokerId = _json.jsonRead("brokerId").toUint();
        _self.role = _json.jsonRead("role").toUint();
        _self.minStartTime = _json.jsonRead("minStartTime").toUint();
        _self.maxStartTime = _json.jsonRead("maxStartTime").toUint();
        _self.minUpdateTime = _json.jsonRead("minUpdateTime").toUint();
        _self.maxUpdateTime = _json.jsonRead("maxUpdateTime").toUint();
        _self.scope = _json.jsonRead("scope").toUint();
        _self.statuses.fromJsonArray(_json.jsonRead("statuses"));
        _self.checkStatuses.fromJsonArray(_json.jsonRead("checkStatuses"));
        _self.bizType = _json.jsonRead("bizType").toUint();
        _self.desc = _json.jsonRead("desc");
        _self.channelType = _json.jsonRead("channelType").toUint();
        _self.bizId = _json.jsonRead("bizId").toUint();
        _self.pledgeeName = _json.jsonRead("pledgeeName");
        _self.pledgorName = _json.jsonRead("pledgorName");
        _self.secPledgeId = _json.jsonRead("secPledgeId").toUint();
        _self.pledgeRegisterNo = _json.jsonRead("pledgeRegisterNo");
        _self.pledgeRegisterFileId = _json.jsonRead("pledgeRegisterFileId");
        _self.minPledgeTime = _json.jsonRead("minPledgeTime").toUint();
        _self.maxPledgeTime = _json.jsonRead("maxPledgeTime").toUint();
        _self.isOnline = _json.jsonRead("isOnline").toUint();
        _self.isEvidenceMailed = _json.jsonRead("isEvidenceMailed").toUint();
        _self.isEvidenceApplied = _json.jsonRead("isEvidenceApplied").toUint();
        _self.paymentStatuses.fromJsonArray(_json.jsonRead("paymentStatuses"));
        _self.refundApplyStatuses.fromJsonArray(_json.jsonRead("refundApplyStatuses"));
        _self.payerName = _json.jsonRead("payerName");
        _self.paymentType = _json.jsonRead("paymentType").toUint();
        _self.minPaymentTime = _json.jsonRead("minPaymentTime").toUint();
        _self.maxPaymentTime = _json.jsonRead("maxPaymentTime").toUint();
        _self.minRefundApplyTime = _json.jsonRead("minRefundApplyTime").toUint();
        _self.maxRefundApplyTime = _json.jsonRead("maxRefundApplyTime").toUint();
        _self.pageSize = _json.jsonRead("pageSize").toUint();
        _self.pageNo = _json.jsonRead("pageNo").toUint();

        if (_self.pageSize > MAX_PAGESIZE) {
            LibLog.log("pageSize is in excess of 500, now it's set to 500");
            _self.pageSize = MAX_PAGESIZE;
        }

        LibJson.pop();
        return true;
    }

    function reset_Condition(Condition storage _self) internal {
        delete _self.businessNo;
        delete _self.userId;
        delete _self.managerId;
        delete _self.brokerId;
        delete _self.role;
        delete _self.minStartTime;
        delete _self.maxStartTime;
        delete _self.minUpdateTime;
        delete _self.maxUpdateTime;
        delete _self.scope;
        _self.statuses.length = 0;
        _self.checkStatuses.length = 0;
        delete _self.bizType;
        delete _self.desc;
        delete _self.channelType;
        delete _self.bizId;
        delete _self.pledgeeName;
        delete _self.pledgorName;
        delete _self.secPledgeId;
        delete _self.pledgeRegisterNo;
        delete _self.pledgeRegisterFileId;
        delete _self.minPledgeTime;
        delete _self.maxPledgeTime;
        delete _self.isOnline;
        delete _self.isEvidenceMailed;
        delete _self.isEvidenceApplied;
        _self.paymentStatuses.length = 0;
        _self.refundApplyStatuses.length = 0;
        delete _self.payerName;
        delete _self.paymentType;
        delete _self.minPaymentTime;
        delete _self.maxPaymentTime;
        delete _self.minRefundApplyTime;
        delete _self.maxRefundApplyTime;
        delete _self.pageSize;
        delete _self.pageNo;
    }

    Condition _cond;

    enum DaoError {
        NO_ERROR,
        BAD_PARAMETER,
        ID_NOT_EXISTS,
        ID_CONFLICTED,
        ITEM_NOT_EXIST
    }

    function OrderDao() {
        register("CsdcModule", "0.0.1.0", "OrderDao", "0.0.1.0");
    }

    //Dao for SecPledgeApply

    function insert_SecPledgeApply(string _json) returns(uint _ret) {
        if (!t_secPledgeApply.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        for (uint i=0; i<m_secPledgeApplyIds.length; ++i) {
            if (m_secPledgeApplyIds[i] == t_secPledgeApply.id)
                return uint(DaoError.ID_CONFLICTED); 
        }

        m_secPledgeApplyMap[t_secPledgeApply.id] = t_secPledgeApply;
        m_secPledgeApplyIds.push(t_secPledgeApply.id);

        return uint(DaoError.NO_ERROR);
    }

    function select_SecPledgeApply_byId(uint _id) constant returns(uint) {

        string memory item = "";
        if(m_secPledgeApplyMap[_id].id != 0) {
            item = m_secPledgeApplyMap[_id].toJson();
        }

        return itemStackPush(item, m_secPledgeApplyIds.length);
    }

    function select_SecPledgeApply_all() constant returns(uint) {
        uint len = 0;
        len = LibStack.push("");
        for (uint i=0; i<m_secPledgeApplyIds.length; ++i) {
            if (i > 0) {
                len = LibStack.append(",");
            }
            len = LibStack.append(m_secPledgeApplyMap[m_secPledgeApplyIds[i]].toJson());
        }

        return itemStackPush(LibStack.popex(len), m_secPledgeApplyIds.length);
    }

    function delete_SecPledgeApply_byId(uint _id) returns(uint _ret) {
        bool found = false;
        for (uint i=0; i<m_secPledgeApplyIds.length; ++i) {
            if (!found) {
                if (m_secPledgeApplyIds[i] == _id) {
                    found = true;
                }
            }

            if (found && i < m_secPledgeApplyIds.length-1) {
                m_secPledgeApplyIds[i] = m_secPledgeApplyIds[i+1];
            }
        }

        if (!found)
            return uint(DaoError.ID_NOT_EXISTS);

        m_secPledgeApplyIds.length -= 1;
        m_secPledgeApplyMap[_id].reset();

        return uint(DaoError.NO_ERROR);
    }

    function update_SecPledgeApply(string _json) returns(uint _ret) {
        if (!t_secPledgeApply.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        if (m_secPledgeApplyMap[t_secPledgeApply.id].id == 0) {
            return uint(DaoError.ID_NOT_EXISTS);
        }

        if (m_secPledgeApplyMap[t_secPledgeApply.id].update(_json))
            return uint(DaoError.NO_ERROR);
        else
            return uint(DaoError.BAD_PARAMETER);
        
    }

    //Dao for DisSecPledgeApply

    function insert_DisSecPledgeApply(string _json) returns(uint _ret) {
        if (!t_disSecPledgeApply.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        for (uint i=0; i<m_disSecPledgeApplyIds.length; ++i) {
            if (m_disSecPledgeApplyIds[i] == t_disSecPledgeApply.id)
                return uint(DaoError.ID_CONFLICTED); 
        }

        m_disSecPledgeApplyMap[t_disSecPledgeApply.id] = t_disSecPledgeApply;
        m_disSecPledgeApplyIds.push(t_disSecPledgeApply.id);

        return uint(DaoError.NO_ERROR);
    }

    function select_DisSecPledgeApply_byId(uint _id) constant returns(uint) {
        string memory item = "";
        if(m_disSecPledgeApplyMap[_id].id != 0) {
            item = m_disSecPledgeApplyMap[_id].toJson();
        }

        return itemStackPush(item, m_disSecPledgeApplyIds.length);
    }

    function select_DisSecPledgeApply_all() constant returns(uint) {
        uint len = 0;
        len = LibStack.push("");
        for (uint i=0; i<m_disSecPledgeApplyIds.length; ++i) {
            if (i > 0) {
                len = LibStack.append(",");
            }
            len = LibStack.append(m_disSecPledgeApplyMap[m_disSecPledgeApplyIds[i]].toJson());
        }

        return itemStackPush(LibStack.popex(len), m_disSecPledgeApplyIds.length);
    }

    function delete_DisSecPledgeApply_byId(uint _id) returns(uint _ret) {
        bool found = false;
        for (uint i=0; i<m_disSecPledgeApplyIds.length; ++i) {
            if (!found) {
                if (m_disSecPledgeApplyIds[i] == _id) {
                    found = true;
                }
            }

            if (found && i < m_disSecPledgeApplyIds.length-1) {
                m_disSecPledgeApplyIds[i] = m_disSecPledgeApplyIds[i+1];
            }
        }

        if (!found)
            return uint(DaoError.ID_NOT_EXISTS);

        m_disSecPledgeApplyIds.length -= 1;
        m_disSecPledgeApplyMap[_id].reset();

        return uint(DaoError.NO_ERROR);
    }

    function update_DisSecPledgeApply(string _json) returns(uint _ret) {
        if (!t_disSecPledgeApply.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        if (m_disSecPledgeApplyMap[t_disSecPledgeApply.id].id == 0) {
            return uint(DaoError.ID_NOT_EXISTS);
        }

        if (m_disSecPledgeApplyMap[t_disSecPledgeApply.id].update(_json))
            return uint(DaoError.NO_ERROR);
        else
            return uint(DaoError.BAD_PARAMETER);
    }

    //Dao for Payment

    function insert_Payment(string _json) returns(uint _ret) {
        if (!t_payment.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        for (uint i=0; i<m_paymentIds.length; ++i) {
            if (m_paymentIds[i] == t_payment.id)
                return uint(DaoError.ID_CONFLICTED); 
        }

        m_paymentMap[t_payment.id] = t_payment;
        m_paymentIds.push(t_payment.id);

        return uint(DaoError.NO_ERROR);
    }

    function select_Payment_byId(uint _id) constant returns(uint) {
        string memory item = "";
        if(m_paymentMap[_id].id != 0) {
            item = m_paymentMap[_id].toJson();
        }

        return itemStackPush(item, m_paymentIds.length);
    }

    function delete_Payment_byId(uint _id) returns(uint _ret) {
        bool found = false;
        for (uint i=0; i<m_paymentIds.length; ++i) {
            if (!found) {
                if (m_paymentIds[i] == _id) {
                    found = true;
                }
            }

            if (found && i < m_paymentIds.length-1) {
                m_paymentIds[i] = m_paymentIds[i+1];
            }
        }

        if (!found)
            return uint(DaoError.ID_NOT_EXISTS);

        m_paymentIds.length -= 1;
        m_paymentMap[_id].reset();

        return uint(DaoError.NO_ERROR);
    }

    function update_Payment(string _json) returns(uint _ret) {
        if (!t_payment.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        if(m_paymentMap[t_payment.id].id == 0)
            return uint(DaoError.ID_NOT_EXISTS);

        if (m_paymentMap[t_payment.id].update(_json)){
            return uint(DaoError.NO_ERROR);
        } else {
            return uint(DaoError.BAD_PARAMETER);
        }

    }

    //Dao for SecPledge

    function insert_SecPledge(string _json) returns(uint _ret) {
        if (!t_secPledge.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        for (uint i=0; i<m_secPledgeIds.length; ++i) {
            if (m_secPledgeIds[i] == t_secPledge.id)
                return uint(DaoError.ID_CONFLICTED); 
        }

        m_secPledgeMap[t_secPledge.id] = t_secPledge;
        m_secPledgeIds.push(t_secPledge.id);

        return uint(DaoError.NO_ERROR);
    }

    function select_SecPledge_byId(uint _id) constant returns(uint) {
        string memory item = "";
        if(m_secPledgeMap[_id].id != 0) {
            item = m_secPledgeMap[_id].toJson();
        }

        return itemStackPush(item, m_secPledgeIds.length);
    }

    function select_SecPledge_all() constant returns(uint) {
        uint len = 0;
        len = LibStack.push("");
        for (uint i=0; i<m_secPledgeIds.length; ++i) {
            if (i > 0) {
                len = LibStack.append(",");
            }
            len = LibStack.append(m_secPledgeMap[m_secPledgeIds[i]].toJson());
        }

        return itemStackPush(LibStack.popex(len), m_secPledgeIds.length);
    }

    function select_SecPledge_status_ById(uint _id) constant returns(uint) {
        return uint(m_secPledgeMap[_id].status);
    }

    function add_SecPledge_disId(uint _id, uint _disId) returns(uint _ret) {
        if(m_secPledgeMap[_id].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        m_secPledgeMap[_id].disSecPedgeApplyIds.push(_disId);
        return uint(DaoError.NO_ERROR);
    }

    function extractPledgeByPledgeRegisterNo(string _no) constant returns(uint) {
        if(m_secPledgeIds.length == 0) {
            return itemStackPush("", 0);
        }

        for(uint i=0; i<m_secPledgeIds.length; i++) {
            t_secPledge = m_secPledgeMap[m_secPledgeIds[i]];
            if(t_secPledge.pledgeRegisterNo.equals(_no)) {
                return itemStackPush(t_secPledge.toJson(), 1);
            }
        }
        return itemStackPush("", 0);
    }

    function delete_SecPledge_byId(uint _id) returns(uint _ret) {
        bool found = false;
        for (uint i=0; i<m_secPledgeIds.length; ++i) {
            if (!found) {
                if (m_secPledgeIds[i] == _id) {
                    found = true;
                }
            }

            if (found && i < m_secPledgeIds.length-1) {
                m_secPledgeIds[i] = m_secPledgeIds[i+1];
            }
        }

        if (!found)
            return uint(DaoError.ID_NOT_EXISTS);

        m_secPledgeIds.length -= 1;
        m_secPledgeMap[_id].reset();

        return uint(DaoError.NO_ERROR);
    }

    function update_SecPledge(string _json) returns(uint _ret) {
        if (!t_secPledge.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        if (m_secPledgeMap[t_secPledge.id].id == 0) {
            return uint(DaoError.ID_NOT_EXISTS);
        }

        if (m_secPledgeMap[t_secPledge.id].update(_json))
            return uint(DaoError.NO_ERROR);
        else
            return uint(DaoError.BAD_PARAMETER);
    }

    function update_SecPledge_status(uint _id, uint _status) returns(uint _ret) {
        if(m_secPledgeMap[_id].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        m_secPledgeMap[_id].status = LibSecPledge.PledgeStatus(_status);
        return uint(DaoError.NO_ERROR);
    }

    function undo_SecPledgeStatus_byDis(uint _id) returns(uint _ret) {
        if(m_disSecPledgeApplyMap[_id].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        uint _secPledgeId = m_disSecPledgeApplyMap[_id].secPledgeId;
        m_secPledgeMap[_secPledgeId].status = m_secPledgeMap[_secPledgeId].statusShow;
        removeFromArray(_id, m_secPledgeMap[_secPledgeId].disSecPedgeApplyIds);
        return uint(DaoError.NO_ERROR);
    }

    //Dao for Biz

    function insert_Biz(string _json) returns(uint _ret) {
        if (!t_biz.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        //t_biz.id = sq.getSeqNo("Biz.id");

        for (uint i=0; i<m_bizIds.length; ++i) {
            if (m_bizIds[i] == t_biz.id)
                return uint(DaoError.ID_CONFLICTED); 
        }

        m_bizMap[t_biz.id] = t_biz;
        m_bizIds.push(t_biz.id);

        return uint(DaoError.NO_ERROR);
    }

    function select_Biz_byId(uint _id) constant returns(uint) {
        string memory item = "";
        if(m_bizMap[_id].id != 0) {
            item = bizToJson(m_bizMap[_id]);
        }

        return itemStackPush(item, m_bizIds.length);
    }

    function select_Biz_all() constant returns(uint) {
        uint len = 0;
        len = LibStack.push("");
        for (uint i=0; i<m_bizIds.length; ++i) {
            if (i > 0) {
                len = LibStack.append(",");
            }
            len = LibStack.append(bizToJson(m_bizMap[m_bizIds[i]]));
        }

        return itemStackPush(LibStack.popex(len), m_bizIds.length);
    }

    function pageBiz(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];
        
            if (_json.jsonKeyExists("channelType") && _cond.channelType != 0 && _cond.channelType != uint(t_biz.channelType)) {
                continue;
            }

            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function pageBizForCsdc(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_CSDC)) {
                continue;
            }

            if (_json.jsonKeyExists("managerId") && _cond.managerId != address(0) && _cond.managerId != t_biz.tradeOperator.id) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function pageBizForProxy(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function pageBizForBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId") || !_json.jsonKeyExists("brokerId") || !_json.jsonKeyExists("role")) {
            LibLog.log("userId, brokerId and role cannot be empty.");
            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__savedByMe(_cond.userId, t_biz) && !__needHandledByBroker(_cond.userId, _cond.brokerId, _cond.role, t_biz) && !__passedByMe(_cond.userId, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }
        
            //质押/解质押业务菜单只包括我经手的，不包括我保存未提交的(需求取消)
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
                continue;
            }

            if (_json.jsonKeyExists("checkStatuses") && _cond.checkStatuses.length > 0 && !uint(t_biz.checkStatus).inArray(_cond.checkStatuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findToDoForBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId") || !_json.jsonKeyExists("brokerId") || !_json.jsonKeyExists("role")) {
            LibLog.log("userId, brokerId and role cannot be empty.");
            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__savedByMe(_cond.userId, t_biz) && !__needHandledByBroker(_cond.userId, _cond.brokerId, _cond.role, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
                continue;
            }
    
            if (_json.jsonKeyExists("checkStatuses") && _cond.checkStatuses.length > 0 && !uint(t_biz.checkStatus).inArray(_cond.checkStatuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findHandledForBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId") || !_json.jsonKeyExists("brokerId") || !_json.jsonKeyExists("role")) {
            LibLog.log("userId, brokerId and role cannot be empty.");

            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__savedByMe(_cond.userId, t_biz) && !__passedByMe(_cond.userId, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
                continue;
            }

            if (_json.jsonKeyExists("checkStatuses") && _cond.checkStatuses.length > 0 && !uint(t_biz.checkStatus).inArray(_cond.checkStatuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findAllBizForCsdcByBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId")) {
            LibLog.log("userId cannot be empty.");

            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (_json.jsonKeyExists("scope")) {
                if (_cond.scope == uint(QueryScope.MY_TODO) && !__needHandledByCsdcFromBroker(_cond.userId, t_biz)) {
                    continue;
                } else if (_cond.scope == uint(QueryScope.MY_PASSBY) && !__passedByMe(_cond.userId, t_biz)) {
                    continue;
                } else if (_cond.scope == uint(QueryScope.ALL) && !__needHandledByCsdcFromBroker(_cond.userId, t_biz) && !__passedByMe(_cond.userId, t_biz)) {
                    continue;
                }
            } else {
                if (!__needHandledByCsdcFromBroker(_cond.userId, t_biz) && !__passedByMe(_cond.userId, t_biz)) {
                    continue;
                }
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }

            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
                continue;
            }

            if (_json.jsonKeyExists("checkStatuses") && _cond.checkStatuses.length > 0 && !uint(t_biz.checkStatus).inArray(_cond.checkStatuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findUnAuditedBizForCsdcByBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId")) {
            LibLog.log("userId cannot be empty.");

            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__needAuditByCsdcFromBroker(_cond.userId, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }

            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findUnReviewedBizForCsdcLeaderByBroker(string _json) constant returns(uint) {
        LibJson.push(_json);

        if (!_json.jsonKeyExists("userId")) {
            LibLog.log("userId cannot be empty.");

            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__needReviewByCsdcLeaderFromBroker(_cond.userId, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }

            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findUnReviewedBizForCsdcByBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId")) {
            LibLog.log("userId cannot be empty.");

            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__needReviewByCsdcFromBroker(_cond.userId, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }

            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findUnCheckedBizForCsdcByBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId")) {
            LibLog.log("userId cannot be empty.");

            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__needCheckByCsdcFromBroker(_cond.userId, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }

            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function findClosedBizForCsdcByBroker(string _json) constant returns(uint) {
        LibJson.push(_json);
        if (!_json.jsonKeyExists("userId")) {
            LibLog.log("userId cannot be empty.");

            LibJson.pop();
            return itemStackPush("", 0);
        }

        if (m_bizIds.length <= 0) {
            LibJson.pop();
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            LibJson.pop();
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__closedByCsdcFromBroker(_cond.userId, t_biz)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }
        
            if (_json.jsonKeyExists("desc") && !_cond.desc.equals("") && !_cond.desc.equals(t_biz.desc)) {
                continue;
            } 
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") && _cond.maxStartTime != 0
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }

            if (_json.jsonKeyExists("minUpdateTime") && _json.jsonKeyExists("maxUpdateTime") && _cond.maxUpdateTime != 0
                && ( t_biz.updateTime < _cond.minUpdateTime || t_biz.updateTime > _cond.maxUpdateTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length > 0 && !__isInStatuses(t_biz, _cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    function getToDoIdListByBroker(address _userId, uint _brokerId, uint _role) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        uint _count = 0;

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__savedByMe(_userId, t_biz) && !__needHandledByBroker(_userId, _brokerId, _role, t_biz)) {
                continue;
            }
        
            if (_count++ > 0) {
              len = LibStack.append(",");
            }
            len = LibStack.append(m_bizIds[i-1].toString());
        }
        return itemStackPush(LibStack.popex(len), _count);
    }

    function getToDoIdListByCsdcFromBroker(address _userId) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        uint _count = 0;

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.BY_BROKER)) {
                continue;
            }

            if (!__needHandledByCsdcFromBroker(_userId, t_biz)) {
                continue;
            }
        
            if (_count++ > 0) {
              len = LibStack.append(",");
            }
            len = LibStack.append(m_bizIds[i-1].toString());
        }
        return itemStackPush(LibStack.popex(len), _count);
    }

    function select_Biz_Status_ById(uint _bizId) constant returns(uint) {
        if(m_bizMap[_bizId].id == 0) {
          return;
        }
        return uint(m_bizMap[_bizId].status);
    }

    function delete_Biz_byId(uint _id) returns(uint _ret) {
        bool found = false;
        for (uint i=0; i<m_bizIds.length; ++i) {
            if (!found) {
                if (m_bizIds[i] == _id) {
                    found = true;
                }
            }

            if (found && i < m_bizIds.length-1) {
                m_bizIds[i] = m_bizIds[i+1];
            }
        }

        if (!found)
            return uint(DaoError.ID_NOT_EXISTS);

        m_bizIds.length -= 1;
        m_bizMap[_id].reset();

        return uint(DaoError.NO_ERROR);
    }

    function update_Biz(string _json) returns(uint _ret) {
        if (!t_biz.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        if (m_bizMap[t_biz.id].id == 0) {
            return uint(DaoError.ID_NOT_EXISTS);
        }

        if (m_bizMap[t_biz.id].update(_json))
            return uint(DaoError.NO_ERROR);
        else
            return uint(DaoError.BAD_PARAMETER);
    }

    function add_Biz_Audit(uint _bizId, string _json) returns(uint _ret) {
        if(!t_audit.fromJson(_json)) {
            return uint(DaoError.BAD_PARAMETER);
        }
        if(m_bizMap[_bizId].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        m_bizMap[_bizId].audits.push(t_audit);
        m_bizMap[_bizId].updateTime = now*1000;
        set_audit_of_status(_bizId, t_audit.oldStatus, t_audit.auditorId);
        return uint(DaoError.NO_ERROR);
    }

    function update_Biz_rejectStatus(uint _bizId, uint _rejectStatus) returns(uint _ret) {
        if(m_bizMap[_bizId].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        if(m_bizMap[_bizId].rejectStatus == uint(LibBiz.RejectStatus.NO_REJECT)) {
            m_bizMap[_bizId].rejectStatus = _rejectStatus;
        }
    }

    function set_audit_of_status(uint _bizId, uint _status, address _operatorId) {
        if(m_operators[_bizId][_status] == address(0)) {
            m_operators[_bizId][_status] = _operatorId;
        }
    }

    function add_Biz_Attachment(uint _bizId, string _json) returns(uint _ret) {
        if(!t_attchInfo.fromJson(_json)) {
            return uint(DaoError.BAD_PARAMETER);
        }
        if(m_bizMap[_bizId].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        m_bizMap[_bizId].backAttachments.push(t_attchInfo);
        return uint(DaoError.NO_ERROR);
    }

    function add_SecpledgeApply_Attachment(uint _id, string _json) returns(uint _ret) {
        if(!t_attchInfo.fromJson(_json)) {
            return uint(DaoError.BAD_PARAMETER);
        }
        if(m_secPledgeApplyMap[_id].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        m_secPledgeApplyMap[_id].frontAttachments.push(t_attchInfo);
        return uint(DaoError.NO_ERROR);
    }

    function update_Biz_Status(uint _bizId, uint _status) returns (uint) {
        if(m_bizMap[_bizId].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        m_bizMap[_bizId].status = LibBiz.BizStatus(_status);
        return uint(DaoError.NO_ERROR);
    }

    function update_Biz_CheckStatus(uint _bizId, uint _checkStatus) returns (uint) {
        if(m_bizMap[_bizId].id == 0) {
            return uint(DaoError.ITEM_NOT_EXIST);
        }
        if (_checkStatus == uint(LibBiz.CheckStatus.BROKER_READ_CSDC_UNREAD) && (m_bizMap[_bizId].checkStatus == LibBiz.CheckStatus.BROKER_UNREAD_CSDC_READ)) {
            _checkStatus = uint(LibBiz.CheckStatus.BOTH_READ);
        }

        if (_checkStatus == uint(LibBiz.CheckStatus.BROKER_UNREAD_CSDC_READ) && (m_bizMap[_bizId].checkStatus == LibBiz.CheckStatus.BROKER_READ_CSDC_UNREAD)) {
            _checkStatus = uint(LibBiz.CheckStatus.BOTH_READ);
        }

        m_bizMap[_bizId].checkStatus = LibBiz.CheckStatus(_checkStatus);
        return uint(DaoError.NO_ERROR);
    }

    //直接通过id更新
    function update_SecPledgeApply_ById(uint id, string _json) returns(uint _ret) {
        if (!t_secPledgeApply.fromJson(_json))
            return uint(DaoError.BAD_PARAMETER);

        if (m_secPledgeApplyMap[id].id == 0) {
            return uint(DaoError.ID_NOT_EXISTS);
        }

        if (m_secPledgeApplyMap[id].update(_json))
            return uint(DaoError.NO_ERROR);
        else
            return uint(DaoError.BAD_PARAMETER);
    }


    //在线业务查询方法 (channelType = 1)

    //查询待办事项
    function findToDo_Online(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];
        
            /* 查询条件 - start */
            if (!__inToDo(_cond.userId, t_biz)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        return itemStackPush(LibStack.popex(len), _count);
    }

    //查询是否有待办事项
    function hasTodo_Online(address _userId) constant returns(bool _ret) {
        for (uint i = 0; i < m_bizIds.length; i++) {
            t_biz = m_bizMap[m_bizIds[i]];
            if (__inToDo(_userId, t_biz)) {
                return true;
            }
        }
    }

    //在线业务分页查询
    function pageBiz_Online(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = m_bizIds.length; i >= 1 ; i--) {
            t_biz = m_bizMap[m_bizIds[i-1]];

            if (t_biz.channelType != uint(LibBiz.ChannelType.ONLINE)) {
                continue;
            }

            if (_json.jsonKeyExists("bizId") && _cond.bizId != 0 && _cond.bizId != t_biz.id) {
                continue;
            }

            if (_json.jsonKeyExists("bizType") && _cond.bizType != 0 && _cond.bizType != uint(t_biz.bizType)) {
                continue;
            }
        
            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }

            if (_json.jsonKeyExists("pledgeeName") && !_cond.pledgeeName.equals("") && !_cond.pledgeeName.equals(t_biz.pledgeeName)) {
                continue;
            }

            if (_json.jsonKeyExists("pledgorName") && !_cond.pledgorName.equals("") && !_cond.pledgorName.equals(t_biz.pledgorName)) {
                continue;
            }
        
            if (_json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != t_biz.pledgorId && _cond.userId != t_biz.pledgeeId && _cond.userId != t_biz.managerId) {
                continue;
            }
        
            if (_json.jsonKeyExists("minStartTime") && _json.jsonKeyExists("maxStartTime") &&_cond. maxStartTime != 0 
                && ( t_biz.startTime < _cond.minStartTime || t_biz.startTime > _cond.maxStartTime)) {
                continue;
            }
        
            if (_json.jsonKeyExists("statuses") && _cond.statuses.length != 0 && !uint(t_biz.status).inArray(_cond.statuses)) {
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
                len = LibStack.append(bizToJson(t_biz));
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    //分页查询解质押记录
    function pageByCond_DisSecPledgeApply(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = 0; i < m_disSecPledgeApplyIds.length; i++) {
            t_disSecPledgeApply = m_disSecPledgeApplyMap[m_disSecPledgeApplyIds[i]];
        
            /* 查询条件 - start */

            if (_json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != t_disSecPledgeApply.pledgorId && _cond.userId != t_disSecPledgeApply.pledgeeId) {
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
                len = LibStack.append(t_disSecPledgeApply.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    //分页查询质物记录
    function pageByCond_SecPledge(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = 0; i < m_secPledgeIds.length; i++) {

            t_secPledge = m_secPledgeMap[m_secPledgeIds[i]];
        
            /* 查询条件 - start */

            if (_json.jsonKeyExists("secPledgeId") && _cond.secPledgeId != 0 && _cond.secPledgeId != t_secPledge.id) {
                continue;
            }

            if (_json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != t_secPledge.pledgorId && _cond.userId != t_secPledge.pledgeeId && _cond.userId != t_secPledge.managerId) {
                continue;
            }

            if (_json.jsonKeyExists("pledgeRegisterNo") && !_cond.pledgeRegisterNo.equals("") && !_cond.pledgeRegisterNo.equals(t_secPledge.pledgeRegisterNo)) {
                continue;
            }

            if (_json.jsonKeyExists("pledgeRegisterFileId") && !_cond.pledgeRegisterFileId.equals("") && !_cond.pledgeRegisterFileId.equals(t_secPledge.pledgeRegisterFileId)) {
                continue;
            }

            if (_json.jsonKeyExists("minPledgeTime") && _json.jsonKeyExists("maxPledgeTime") && _cond.maxPledgeTime != 0 
                && ( t_secPledge.pledgeTime < _cond.minPledgeTime || t_secPledge.pledgeTime > _cond.maxPledgeTime)) {
                continue;
            }

            if (_json.jsonKeyExists("isOnline") && _cond.isOnline != 0 && _cond.isOnline != uint(t_secPledge.isOnline)) {
                continue;
            }

            if (_json.jsonKeyExists("isEvidenceMailed") && _cond.isEvidenceMailed != 0 && _cond.isEvidenceMailed != uint(t_secPledge.isEvidenceMailed)) {
                continue;
            }

            if (_json.jsonKeyExists("isEvidenceApplied") && _cond.isEvidenceApplied != 0 && _cond.isEvidenceApplied != uint(t_secPledge.isEvidenceApplied)) {
                continue;
            }

            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_secPledge.businessNo) ) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length != 0 && !uint(t_secPledge.status).inArray(_cond.statuses)) {
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
                len = LibStack.append(t_secPledge.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    //查找可以解除質押的成交記錄
    function findToDis_SecPledge_Online(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = 0; i < m_secPledgeIds.length; i++) {
            t_secPledge = m_secPledgeMap[m_secPledgeIds[i]];
        
            /* 查询条件 - start */
            /* 在线只能解除在线办理的质物状态 */
            if (t_secPledge.isOnline != uint(LibSecPledge.IsOnline.YES) || 
                m_bizMap[t_secPledge.secPledgeApplyId].channelType != uint(LibBiz.ChannelType.ONLINE) )
            {
                continue;
            }

            if (_json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != t_secPledge.pledgeeId) {
                continue;
            }

            if (t_secPledge.status == LibSecPledge.PledgeStatus.NONE || t_secPledge.status == LibSecPledge.PledgeStatus.DISPLEGING) {
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
                len = LibStack.append(t_secPledge.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }

    /* 付款信息查询 
        channelType、userId、pageNo、pageSize
        businessNo、refundApplyStatuses、pledgorName、pledgeeName、
        payerName、paymentType、time
    */
    function pagePayment(string _json) constant returns(uint) {
        if (m_bizIds.length <= 0) {
            return itemStackPush("", 0);
        }

        fromJson_Condition(_cond, _json);
        
        uint _startIndex = _cond.pageSize * _cond.pageNo;
        
        if (_startIndex >= m_bizIds.length) {
            return itemStackPush("", 0);
        }
        
        uint _count = 0; //满足条件的数据条数
        uint _total = 0; //满足条件的指定页数的数据条数

        uint len = 0;
        len = LibStack.push("");
        LibJson.push(_json);
        for (uint i = 0; i < m_paymentIds.length; i++) {
            t_payment = m_paymentMap[m_paymentIds[i]];
            t_biz = m_bizMap[t_payment.id];
            t_secPledgeApply = m_secPledgeApplyMap[t_payment.id];
        
            /* 查询条件 - start */

            /* 忽略付款金额为0的数据 */
            if (t_secPledgeApply.payAmount == 0) {
                continue;
            }

            /* 在线只能解除在线办理的质物状态 */
            if (_json.jsonKeyExists("channelType") && _cond.channelType != 0 && _cond.channelType != t_biz.channelType) {
                continue;
            }

            if (_json.jsonKeyExists("userId") && _cond.userId != address(0) && _cond.userId != t_secPledgeApply.payerAccount) {
                continue;
            }

            if (_json.jsonKeyExists("statuses") && _cond.statuses.length != 0 && !uint(t_biz.status).inArray(_cond.statuses)) {
                continue;
            }

            if (_json.jsonKeyExists("paymentStatuses") && _cond.paymentStatuses.length != 0 && !uint(t_payment.status).inArray(_cond.paymentStatuses)) {
                continue;
            }

            if (_json.jsonKeyExists("refundApplyStatuses") && _cond.refundApplyStatuses.length != 0 && !uint(t_payment.refundApplyStatus).inArray(_cond.refundApplyStatuses)) {
                continue;
            }

            if (_json.jsonKeyExists("businessNo") && !_cond.businessNo.equals("") && !_cond.businessNo.equals(t_biz.businessNo)) {
                continue;
            }

            if (_json.jsonKeyExists("pledgorName") && !_cond.pledgorName.equals("") && !__isNameInUsers(_cond.pledgorName, t_secPledgeApply.pledgors)) {
                continue;
            }

            if (_json.jsonKeyExists("pledgeeName") && !_cond.pledgeeName.equals("") && !_cond.pledgeeName.equals(t_secPledgeApply.pledgee.name)) {
                continue;
            }

            if (_json.jsonKeyExists("payerName") && !_cond.payerName.equals("") && !_cond.payerName.equals(t_secPledgeApply.payerName)) {
                continue;
            }

            if (_json.jsonKeyExists("paymentType") && _cond.paymentType != 0 && _cond.paymentType != uint(t_payment.paymentType)) {
                continue;
            }

            if (_json.jsonKeyExists("minPaymentTime") && _json.jsonKeyExists("maxPaymentTime") && _cond.maxPaymentTime != 0 
                && ( t_payment.time < _cond.minPaymentTime || t_payment.time > _cond.maxPaymentTime)) {
                continue;
            }

            if (_json.jsonKeyExists("minRefundApplyTime") && _json.jsonKeyExists("maxRefundApplyTime") && _cond.maxRefundApplyTime != 0 
                && ( t_payment.refundApplyTime < _cond.minRefundApplyTime || t_payment.refundApplyTime > _cond.maxRefundApplyTime)) {
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
                len = LibStack.append(t_payment.toJson());
            }
        }
        LibJson.pop();
        return itemStackPush(LibStack.popex(len), _count);
    }


    /* 以下是内部调用方法 */

    function __isNameInUsers(string _name, LibTradeUser.TradeUser[] storage _users) internal returns(bool) {
        for(uint i=0; i<_users.length; i++) {
            if(_name.equals(_users[i].name)) {
                return true;
            }
        }
    }
    
    function __isInStatuses(LibBiz.Biz _biz, uint[] storage _statuses) internal returns(bool) {
        uint _status;
        for(uint i=0; i<_statuses.length; i++) {
            _status = _statuses[i];
            if(_status == uint(_biz.status)) {
                return true;
            } else {
                if( _status == uint(LibBiz.BizStatus.QUERY_BROKER_UNCHECKED) &&
                   (_biz.checkStatus == LibBiz.CheckStatus.NEITHER_READ || _biz.checkStatus == LibBiz.CheckStatus.BROKER_UNREAD_CSDC_READ)
                ) {
                    return true;
                } else if(_status == uint(LibBiz.BizStatus.QUERY_CSDC_UNCHECKED) &&
                   (_biz.checkStatus == LibBiz.CheckStatus.NEITHER_READ || _biz.checkStatus == LibBiz.CheckStatus.BROKER_READ_CSDC_UNREAD)
                ) {
                    return true;
                }
            }
        }
    }

    //判断该业务是否在某用户的待办事项中
    function __inToDo(address _userId, LibBiz.Biz _biz) internal returns (bool) {
        if(_biz.bizType == LibBiz.BizType.PLEDGE_BIZ && __inToDo_pledge(_userId, _biz)) {
          return true;
        }
        if(_biz.bizType == LibBiz.BizType.DISPLEDGE_BIZ && __inToDo_disPledge(_userId, _biz)) {
          return true;
        }
    }

  //1.证券质押业务
    function __inToDo_pledge(address _userId, LibBiz.Biz _biz) internal returns (bool) {
        // 当前用户为质权人
        if(_biz.pledgeeId == _userId) {
          if( _biz.status == LibBiz.BizStatus.PLEDGE_PLEDGEE_UNCONFIRMED  ||  //待质权人确认
              _biz.status == LibBiz.BizStatus.PLEDGE_PLEDGEE_UNPAID       ||  //待质权人付款
              _biz.status == LibBiz.BizStatus.PLEDGE_PLEDGEE_UNAUTH           //待质权人人脸认证
            )
            return true;
        }
        
        // 当前用户为出质人
        if(_biz.pledgorId == _userId) {
          if( _biz.status == LibBiz.BizStatus.PLEDGE_PLEDGOR_UNPAID  ||  //待出质人付款
              _biz.status == LibBiz.BizStatus.PLEDGE_PLEDGOR_UNAUTH      //待出质人人脸认证
            )
            return true;
        }
    }

    //2.解除证券质押业务
    function __inToDo_disPledge(address _userId, LibBiz.Biz _biz) internal returns (bool) {
        // 当前用户为质权人
        if(_biz.pledgeeId == _userId) {
          if( _biz.status == LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNCONFIRMED  ||  //待质权人确认
              _biz.status == LibBiz.BizStatus.DISPLEDGE_PLEDGEE_UNAUTH           //待质权人人脸认证
            )
            return true;
        }
        
        // 当前用户为出质人
        if(_biz.pledgorId == _userId) {
          if( _biz.status == LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNCONFIRMED  ||  //待出质人确认
              _biz.status == LibBiz.BizStatus.DISPLEDGE_PLEDGOR_UNAUTH           //待出质人人脸认证
            )
            return true;
        }
    }    

    /* 由我保存的 */
    function __savedByMe(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        return _biz.tradeOperator.id == _userId && 
        (_biz.status == LibBiz.BizStatus.BROKER_PLEDGE_STASHED ||
         _biz.status == LibBiz.BizStatus.BROKER_PLEDGE_CREATED ||
         _biz.status == LibBiz.BizStatus.BROKER_DISPLEDGE_STASHED ||
         _biz.status == LibBiz.BizStatus.BROKER_DISPLEDGE_CREATED );
    }

    /* 需要我办理的（券商用户）*/
    /* 
        1.同一券商用户经办
        2.当前状态为待券商复核
        3.我的角色为复核人
    */
    function __needHandledByBroker(address _userId, uint _brokerId, uint _role, LibBiz.Biz _biz) internal constant returns(bool) {
        //同一券商业务
        if (_brokerId == _biz.tradeOperator.brokerId) {
            //我的角色为复核人，业务状态为待券商复核
            //1.该笔业务未被复核过 2.该笔业务由我复核
            if(_role == uint(LibBrokerUser.Role.PLEDGE_REVIEWER)) {
                return (_biz.status == LibBiz.BizStatus.BROKER_PLEDGE_UNAUDITED || _biz.status == LibBiz.BizStatus.BROKER_DISPLEDGE_UNAUDITED) &&
                       (m_operators[_biz.id][uint(_biz.status)] == address(0) || m_operators[_biz.id][uint(_biz.status)] == _userId);
            }
            //我的角色为经办人，
            // 1.当前业务待重新录入
            // 2.业务状态为待经办人查看结果
            if(_role == uint(LibBrokerUser.Role.PLEDGE_OPERATOR) && _biz.tradeOperator.id == _userId) {
                return  _biz.status == LibBiz.BizStatus.BROKER_PLEDGE_UNRETYPED || 
                        _biz.status == LibBiz.BizStatus.BROKER_DISPLEDGE_UNRETYPED ||
                        _biz.checkStatus == LibBiz.CheckStatus.NEITHER_READ || 
                        _biz.checkStatus == LibBiz.CheckStatus.BROKER_UNREAD_CSDC_READ;
            }
        }
    }

    /* 待中证登初审（券商提交） */
    function __needAuditByCsdcFromBroker(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        return  _biz.channelType == uint(LibBiz.ChannelType.BY_BROKER) &&
                (_biz.status == LibBiz.BizStatus.BROKER_PLEDGE_CSDC_UNAUDITED || _biz.status == LibBiz.BizStatus.BROKER_DISPLEDGE_CSDC_UNAUDITED) &&
                (m_operators[_biz.id][uint(_biz.status)] == address(0) || m_operators[_biz.id][uint(_biz.status)] == _userId);
    }

    /* 待中证登人员复核（券商提交） */
    function __needReviewByCsdcFromBroker(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        return _biz.channelType == uint(LibBiz.ChannelType.BY_BROKER) &&
              //质押业务，当前状态为带中证登复核，且初审人不是当前用户
            ( (_biz.status == LibBiz.BizStatus.BROKER_PLEDGE_CSDC_UNREVIEWED && m_operators[_biz.id][uint(LibBiz.BizStatus.BROKER_PLEDGE_CSDC_UNAUDITED)] != _userId) ||
              //解除质押业务，当前状态为带中证登复核，且初审人不是当前用户
              (_biz.status == LibBiz.BizStatus.BROKER_DISPLEDGE_CSDC_UNREVIEWED && m_operators[_biz.id][uint(LibBiz.BizStatus.BROKER_DISPLEDGE_CSDC_UNAUDITED)] != _userId) ) && 
              (m_operators[_biz.id][uint(_biz.status)] == address(0) || m_operators[_biz.id][uint(_biz.status)] == _userId);
    }

    function __needReviewByCsdcLeaderFromBroker(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        return  _biz.channelType == uint(LibBiz.ChannelType.BY_BROKER) &&
                _biz.csdcLeader.id == _userId &&
               (_biz.status == LibBiz.BizStatus.BROKER_PLEDGE_LEADER_UNREVIEWED || 
                _biz.status == LibBiz.BizStatus.BROKER_DISPLEDGE_LEADER_UNREVIEWED);
    }

    /* 待中证登人员（初审人）查看（券商提交） */
    function __needCheckByCsdcFromBroker(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        return  _biz.channelType == uint(LibBiz.ChannelType.BY_BROKER) &&
                ( _biz.checkStatus == LibBiz.CheckStatus.NEITHER_READ || 
                  _biz.checkStatus == LibBiz.CheckStatus.BROKER_READ_CSDC_UNREAD
                ) && 
                ( m_operators[_biz.id][uint(LibBiz.BizStatus.BROKER_PLEDGE_CSDC_UNAUDITED)] == _userId ||
                  m_operators[_biz.id][uint(LibBiz.BizStatus.BROKER_DISPLEDGE_CSDC_UNAUDITED)] == _userId
                );
    }

    /* 需要我办理的（中证登用户）*/
    /* 
        1.待审核
        2.待复核
        3.待查看
    */
    function __needHandledByCsdcFromBroker(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        return __needAuditByCsdcFromBroker(_userId, _biz) || 
               __needReviewByCsdcFromBroker(_userId, _biz) ||
               __needCheckByCsdcFromBroker(_userId, _biz) ||
               __needReviewByCsdcLeaderFromBroker(_userId, _biz);
    }

    /* 已办结且已查看 */
    function __closedByCsdcFromBroker(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        return  _biz.checkStatus == LibBiz.CheckStatus.BOTH_READ &&
                __passedByMe(_userId, _biz);
    }

    /* 我经手的业务(包括我经办的和我审核的) */
    function __passedByMe(address _userId, LibBiz.Biz _biz) internal constant returns(bool) {
        //不包括还未提交的业务
        if(__savedByMe(_userId, _biz)) {
            return false;
        }
        //我经办的业务
        if(_biz.tradeOperator.id == _userId) {
            return true;
        }
        //我审核的业务
        for(uint i = 0; i < _biz.audits.length; i++) {
            t_audit = _biz.audits[i];
            if(t_audit.auditorId == _userId) {
                return true;
            }
        }
    }

    function getOperator(uint _bizId, uint _statuses) constant returns(address) {
        return m_operators[_bizId][_statuses];
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

    function bizToJson(LibBiz.Biz storage _biz) constant internal returns(string) {
        if(_biz.bizType == LibBiz.BizType.PLEDGE_BIZ) {
            _biz.pledgeStatus = uint(m_secPledgeMap[_biz.id].statusShow);
        } else if(_biz.bizType == LibBiz.BizType.DISPLEDGE_BIZ) {
            _biz.pledgeStatus = uint(m_secPledgeMap[m_disSecPledgeApplyMap[_biz.id].secPledgeId].statusShow);
        }
        return _biz.toJson();
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