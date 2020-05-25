pragma solidity ^0.4.24;

/**
 * 0.4.25版本合约和地址成员变量/函数验证
 * 1.0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 send()成员函数验证
 * 2.0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 transfer()成员函数验证
 * 3.0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 balance成员变量验证
 * 4.0.4.25版本msg.sender类型所属验证
 * @author Albedo
 * @dev 2019/12/24
 **/
contract ContractAndAddressFunction {

    function addressCheck() public view returns (address,uint256,uint256) {
        address x = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqfreqpsmj"; //0x123、lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqfrk9nl4a
        address myAddress = this;
        //0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 balance成员变量验证
        if (x.balance < 10 && myAddress.balance >= 10)
        {
            //0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 transfer()成员函数验证
            x.transfer(10);
            //0.4.25版本contract合约类型包括 address类型的成员函数，可以直接使用 send()成员函数验证
            x.send(10);
        }

        //0.4.25版本msg.sender类型所属验证
        address sender=msg.sender;
        return (sender,myAddress.balance,x.balance);
    }
    function() payable external {}
}