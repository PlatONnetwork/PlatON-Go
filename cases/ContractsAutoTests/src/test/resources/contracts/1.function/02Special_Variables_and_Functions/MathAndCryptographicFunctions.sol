pragma solidity 0.5.13;
/**
 * 验证数学和加密函数,sha3在0.5.0版本之后被剔除,用keccak256取代
 * 计算(x + y) % k
 * 计算(x * y) % k
 * 计算输入的Keccak256Hash值
 * 计算输入的Sha256Hash值
 * 计算输入的Ripemd160Hash值
 * 从椭圆曲线签名中恢复与公钥相关的地址，或在出错时返回零
 * 函数参数对应于签名的ECDSA值: r – 签名的前32字节; s: 签名的第二个32字节; v: 签名的最后一个字节
 * @author liweic
 * @dev 2019/12/27 20:10
 */

contract MathAndCryptographicFunctions {
    
    function callAddMod() public pure returns(uint){
        return addmod(2, 3, 3);
    }
   
    function callMulMod() public pure returns(uint){
        return mulmod(2, 3, 3);
    }
   
    function callKeccak256() public pure returns(bytes32 result){
        return keccak256("ABC");
    }
   
    function callSha256() public pure returns(bytes32 result){
        return sha256("ABC");
    }
   
    function callRipemd160() public pure returns(bytes32 result){
        return ripemd160("ABC");
    }
   
    //hash: "0xe281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d"
    //V = 27
    //R = "0x55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe"
    //S = "0x2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6"
    //address: "0x8a9B36694F1eeeb500c84A19bB34137B05162EC5"
    function callEcrecover(bytes32 hash, uint8 v, bytes32 r, bytes32 s) public pure returns (address) {
        address x = ecrecover(hash, v, r, s);
        return x;
    }
}