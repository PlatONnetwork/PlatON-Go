pragma solidity 0.5.13;



/**
 *
 * 
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约继承功能点
 *
 *继承(is)简述：合约支持多重继承，即当一个合约从多个合约继承时，
 *在区块链上只有一个合约被创建，所有基类合约的代码被复制到创建合约中。
 *
 *
 *-----------------  测试点   ------------------------------
 *1、多重继承情况
 *2、多重继承(合约存在父子关系)
 *3、继承支持传参
 *4、合约函数重载(Overload)
 */
 




/**
 *
 *验证：2、多重继承(父类合约存在父子关系)，如果继承的合约之间有父子关系，是否必须遵循先父到子的继承顺序
 *-------------------------------
 *验证结果：多重继承(父类合约存在父子关系)，合约继承必须遵循先父到子的继承顺序，否则会编译异常。
 *
 */

import "./01InheritContractMultipleTest.sol";


 contract InheritContractParentThreeClass is InheritContractParentTwoClass {

 function getDataThree() public view returns (uint) {
        return 3;
    }

}




contract InheritContractSubclass is InheritContractParentThreeClass,InheritContractParentTwoClass {
  
  function getSubData() public view returns (uint) {
        return 3;
    }

}


