pragma solidity 0.5.13;
/**
 * 验证address的查余额方法banlance和转账方法send,transfer
 * @author liweic
 * @dev 2019/12/28 10:10
 */

contract AddressFunctions {
    
    //获取地址的余额
    function getBalance(address addr) public returns (uint){
        return addr.balance;
    }
    
    //当前合约的余额  
    function getBalanceOf() public returns (uint){
        return address(this).balance;
    }

    function transfer(address payable addr) public payable{
        addr.transfer(msg.value);
    }

    function transfer4(address payable addr) public payable {
        addr.send(10 ether);
    }

}