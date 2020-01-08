pragma solidity ^0.4.12;
/**
* @file LibSequence.sol
* @author yiyating
* @time 2016-12-27
* @desc id管理合約
*/


import "./utillib/DateTime.sol";
import "./sysbase/OwnerNamed.sol";

contract Sequence is OwnerNamed{
  
    using LibInt for *;

    mapping(string=>uint) sequences;
    
    function Sequence(){
        register("CsdcModule", "0.0.1.0", "Sequence", "0.0.1.0");
    }

    function getSeqNo(string _key) public returns (uint) {
      // return (++sequences[_key]).toString();
      return ++sequences[_key];
    }

    function genBusinessNo(string _char, string _bizType) returns(uint){
        _char = _char.concat(_bizType);
        uint year = uint(DateTime.getYear(now));
        string memory yearStr = (year - year/100*100).toString();
        uint month = uint(DateTime.getMonth(now));
        string memory monthStr;
        if (month<10) {
            monthStr = "0".concat(month.toString());
        } else {
            monthStr = month.toString();
        }
        uint day = uint(DateTime.getDay(now));
        string memory dayStr;
        if (day<10) {
            dayStr = "0".concat(day.toString());
        } else {
            dayStr = day.toString();
        }
        string memory noKey = _char.concat(yearStr, monthStr, dayStr);
        string memory no = getZeroNo(getSeqNo(noKey)).recoveryToString();
        return noKey.concat(no).storageToUint();
    }

    function refreshBusinessNo(string _char, string _bizType) {
        uint _uNo = genBusinessNo(_char, _bizType);
        string memory _sNo = _uNo.recoveryToString();

        string memory _ret = "{\"ret\":0, \"message\": \"success\", \"data\":{\"total\":1,";
        _ret = _ret.concat(_sNo.toKeyValue("businessNo"), ", \"items\":[]}}");
        Notify(0, _ret);
    }

  /* 以下为内部调用方法 */

    function getZeroNo(uint _no) internal constant returns(uint) {
        string memory _zeros;
        if(_no < 10) {
          _zeros = "000";
        } else if (_no < 100) {
          _zeros = "00";
        } else if (_no < 1000) {
          _zeros = "0";
        } else if (_no < 10000){
          _zeros = "";
        }
        _zeros = _zeros.concat(_no.toString());
        return _zeros.storageToUint();
    }

    event Notify(uint _errorno, string _info);

}