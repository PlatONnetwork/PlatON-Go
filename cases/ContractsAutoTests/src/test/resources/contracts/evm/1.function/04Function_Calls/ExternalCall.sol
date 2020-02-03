pragma solidity 0.5.13;
/**
 * 验证函数的外部调用,一个函数调用另一个合约的函数
 * @author liweic
 * @dev 2019/12/27 10:10
 */

contract External {
    function getSum() external view returns(uint sum){
        //定义两个局部变量
        uint a = 1;
        uint b = 2;
        sum = a + b;
    }
}

contract CallExternal {
    event ExternalCValue(uint value);
    function getResult() payable public returns(uint c){
        External d = new External();
        c = d.getSum();
        emit ExternalCValue(c);
    }
}