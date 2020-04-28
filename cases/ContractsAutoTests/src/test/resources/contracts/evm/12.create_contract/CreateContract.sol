pragma solidity ^0.5.13;
/**
 * new关键字创建合约验证
 *
 *
 * @author Albedo
 * @dev 2019/12/19
 **/
//目标创建合约
contract NewTargetCreateContract {
    uint public x;  // solidity自动为public变量创建同名方法x()
    uint public amount;

    constructor(uint _a) public payable {
        x = _a;
        amount = msg.value;
    }
}

contract CreateContract {
    //因为getTargetCreateContractData使用view修饰，所以将new过程放在函数外
    NewTargetCreateContract createContract = new NewTargetCreateContract(1000);
    function getTargetCreateContractData() public view returns (uint,uint){
        //获取new创建合约实际值
        return (createContract.x(),createContract.amount());
    }
}


