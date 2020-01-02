pragma solidity 0.5.13;
/**
 * 验证函数声明方式payable,用于转账
 * @author liweic
 * @dev 2019/12/27 16:10
 */

contract Payable {
    mapping (address => uint) balances;
        
    function balanceOf(address _user) public returns (uint) { 
        return balances[_user]; 
    }
    
    function deposit() public payable { 
        balances[msg.sender] += msg.value; 
    }
    
    function withdraw(uint _amount) public {
        require(balances[msg.sender] - _amount > 0);  // 确保账户有足够的余额可以提取
        msg.sender.transfer(_amount);
        balances[msg.sender] -= _amount;
    }
}