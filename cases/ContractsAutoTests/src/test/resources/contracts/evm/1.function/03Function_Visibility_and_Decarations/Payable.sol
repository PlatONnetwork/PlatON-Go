pragma solidity 0.5.13;
/**
 * 验证函数声明方式payable,用于转账
 * @author liweic
 * @dev 2019/12/27 16:10
 */

contract Payable {

    //获取地址的余额
    function getBalances(address addr) view public returns (uint){
        return addr.balance;
    }

    function transfer(address payable addr) public payable{
        addr.transfer(msg.value);
    }
}