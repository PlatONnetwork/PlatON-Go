pragma solidity ^0.5.13;
/**
 * 跨合约调用注意事项：方法签名必须用全名，而不是别名。返回的数据是16进制，所以数据要进行转换。
 *
 * @author hudenian
 * @dev 2020/04/14 11:09
 */
contract WithBackCaller{

    uint256 callerUintResult;

    string callerStringResult;


    /**
     * 参数类型要为uint256不能为uint
     *
     */
    function callAddlTest(address other) public {
        // bytes4 messageId = bytes4(keccak256("add(uint, uint)"));
        // other.call(messageId, 5, 60);
        other.call(abi.encodeWithSignature("add(uint256,uint256)", 85, 60));
    }

    function callDoublelTest(address other,uint256 a) public{
        // bytes4 messageId = bytes4(keccak256("add(uint, uint)"));
        // other.call(messageId, 5, 60);
        (bool success, bytes memory data) = other.call(abi.encodeWithSignature("getDouble(uint256)", a));
        if(!success){
            revert();
        }
        callerUintResult = abi.decode(data,(uint256));
    }


    function getuintResult() view public returns(uint256){
        return callerUintResult;
    }


    function callgetNameTest(address other,string memory name) public {
        (bool success, bytes memory data) = other.call(abi.encodeWithSignature("getName(string,string)","hello",name));
        if(!success){
            revert();
        }
        callerStringResult = abi.decode(data,(string));
    }

    function callgetNameTestWithGas(address other,string memory name,uint256 gasValue) public {
        (bool success, bytes memory data) = other.call.gas(gasValue)(abi.encodeWithSignature("getName(string,string)","hellogas",name));
        if(!success){
            revert();
        }
        callerStringResult = abi.decode(data,(string));
    }

    function getStringResult() view public returns(string memory){
        return callerStringResult;
    }


}
