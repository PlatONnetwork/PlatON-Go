pragma solidity ^0.5.13;
/**
 * 跨合约调用中的被调用者
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */
contract WithBackCallee{
//    event FallbackCalledEvent(bytes data);
//    event DoubleEvent(uint256 a, uint256 b);
//    event GetNameEvent(string);

//    function() external{
//        emit FallbackCalledEvent(msg.data);
//    }



    function getDouble(uint256 a) public returns(uint256){
        uint256 _result = a + a;
//        emit DoubleEvent(a, _result);
        return _result;
    }

    function getName(string memory option,string memory name) public returns(string memory){
//        emit GetNameEvent(name);
        return strConcat(option,name);
    }


     function strConcat(string memory _a, string memory _b) internal returns (string memory){
         bytes memory _ba = bytes(_a);
         bytes memory _bb = bytes(_b);
         string memory ret = new string(_ba.length + _bb.length);
         bytes memory bret = bytes(ret);
         uint k = 0;
         for (uint i = 0; i < _ba.length; i++)bret[k++] = _ba[i];
         for (uint i = 0; i < _bb.length; i++) bret[k++] = _bb[i];
         return string(ret);
     }

}