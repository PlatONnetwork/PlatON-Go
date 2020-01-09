pragma solidity ^0.4.12;
/**
* @file LibSequence.sol
* @author yiyating
* @time 2017-01-17
* @desc 计费library
*/


import "../utillib/LibInt.sol";

library LibBilling {

    // enum BillingType { NONE, SECPLEDGE }

    function calc_secPledge(uint _shareNum, uint _sharePrice) internal returns (uint _ret) {
        if(_shareNum <= 5000000) { //500万股（含）以内按面值的1‰计算
            _ret = LibInt.divide(_shareNum*_sharePrice, 1000);
        } else { //超出500万股的部分按面试的0.1‰计算
            _ret = 5000*_sharePrice + LibInt.divide((_shareNum-5000000)*_sharePrice, 10000);
        }
    }

}