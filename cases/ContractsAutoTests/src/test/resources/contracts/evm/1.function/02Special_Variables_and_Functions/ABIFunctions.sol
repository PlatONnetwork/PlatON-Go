pragma solidity 0.5.11;
/**
 * 验证ABI编解码相关的函数
 * @author liweic
 * @dev 2019/12/27 20:10
 */
contract ABIFunctions {
    function getEncodeWithSignature() public view returns (bytes memory) {
        return abi.encodeWithSignature("set(uint256)", 1); //计算函数set(uint256) 及参数1 的ABI 编码
    }
    
    function getEncode() public view returns (bytes memory) {
        return abi.encode(1); //计算参数 1 的ABI 编码
    }
    
    function getEncodePacked() public view returns (bytes memory) {
        return abi.encodePacked("1"); //计算参数 1 的紧密打包编码
    }
}