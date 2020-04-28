pragma solidity 0.5.13;

/**
 * 验证单一输入、输出参数类型uint
 * 验证函数无返回值
 * 验证未使用参数可以省略参数名
 * 验证入参和出参均为数组
 * 验证返回值为字符串
 * 验证多个返回值且返回类型是数组
 * @author liweic
 * @dev 2019/12/28 19:09
 */

contract PramaAndReturns {
    string greeting = "What's up man";

    uint public s ;
    
    //入参出参均为uint
    function InputParam(uint a) public view returns(uint b){
        b = a;
        return b;
    }
    
    //无返回值
    function NoOutput(uint a, uint b) public {
        s = a;
    }

    //验证函数NoOutput是否调用成功
    function getS() public view returns(uint){
        return s;
    }
    
    //未使用参数可以省略参数名
    function OmitParam(uint y, uint) public pure returns(uint) {
        return y;
    }
    
    //入参和出参均为数组
    function IuputArray(uint[3] memory y) public view returns (uint[3] memory){
        y[2] = 3;
        return y;
    }

    //返回值为字符串
    function OuputString() public view returns (string memory) {
        return greeting;
    }
    
    //多个返回值且返回类型是数组
    function OuputArrays() public pure returns(uint[] memory, uint[] memory) {

        uint[] memory localMemoryArray1 = new uint[](3);  
        localMemoryArray1[0] = 1;  
        localMemoryArray1[1] = 2;  
        localMemoryArray1[2] = 3;

        uint[] memory localMemoryArray2 = localMemoryArray1;  
        localMemoryArray1[0] = 10;

        return (localMemoryArray1, localMemoryArray2); 
       //returns 1,2,3 | 10,2,3 
    }
    
}