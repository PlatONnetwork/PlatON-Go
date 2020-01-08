pragma solidity ^0.4.12;
/**
* file LibDecode.sol
* author yangzhou
* time 2017-6-19
* desc the defination of LibDecode libary
*/


import "../utillib/LibString.sol";

library LibDecode {
    using LibString for *;
    using LibDecode for *;
    function decode(string _signString, bytes32 _hash) internal returns (address addr){
        bytes memory signedString = _signString.toHex();
        string memory defaultAddr = "0000000000000000000000000000000000000000";
        if (signedString.length < 65)
        {
            return defaultAddr.toAddress();
        }
        bytes32  r = bytesToBytes32(slice(signedString, 0, 32));
        bytes32  s = bytesToBytes32(slice(signedString, 32, 32));
        bytes1  v = slice(signedString, 64, 1)[0];
        return ecrecover(_hash, uint8(v), r, s);
    }

    //将原始数据按段切割出来指定长度
    function slice(bytes memory _data, uint _start, uint _len) internal returns (bytes){
        bytes memory b = new bytes(_len);

        for(uint i = 0; i < _len; i++){
        b[i] = _data[i + _start];
        }

        return b;
    }


    //bytes转换为bytes32
    function bytesToBytes32(bytes memory source) internal returns (bytes32 result) {
        assembly {
            result := mload(add(source, 32))
        }
    }
}
