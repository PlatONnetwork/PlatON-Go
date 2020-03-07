pragma solidity ^0.5.0;
/**
 * 函数可见性必须显式声明（0.4.x可以不显式声明）
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 *
 */


contract FunctionDeclaraction{
    uint balance;

    //1.函数可见性声明为public
    function update_public(uint amount_pu) public returns (address, uint){
        update_internal(amount_pu);
        return (msg.sender, balance);
    }


    //2.函数可见性显示声明为external
    function update_external(uint amount_ex) external returns (address, uint){
        update_private(amount_ex);
        return (msg.sender, balance);
    }

    //3.函数可见性显示声明为internal
    function update_internal(uint amount_in) internal returns (address, uint){
        balance += amount_in;
        return (msg.sender, balance);
    }

    //4.函数可见性显示声明为private
    function update_private(uint amount_pr) private returns (address, uint){
        balance += amount_pr;
        return (msg.sender, balance);
    }

    //5.查询balance
    function getBalance() public view returns (uint){
        return balance;
    }
}
