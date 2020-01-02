pragma solidity 0.5.13;
/**
 * 验证函数的内部调用,一个函数在同一个合约中调用另一个函数
 * @author liweic
 * @dev 2019/12/26 17:10
 */

contract IntenalCall {
    
    function getSum() public view returns(uint sum){
        //定义两个局部变量
       uint a = 1;
       uint b = 2;
       sum = a + b;
    }
   
    function getResult() public view returns(uint product){
       uint c = 3;
       product = c * getSum();  //内部调用之前定义的函数
    }
}