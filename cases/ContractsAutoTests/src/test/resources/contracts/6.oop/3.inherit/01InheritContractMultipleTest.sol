pragma solidity 0.5.13;


/**
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约继承功能点
 *
 *继承(is)简述：合约支持多重继承，即当一个合约从多个合约继承时，
 *在区块链上只有一个合约被创建，所有基类合约的代码被复制到创建合约中。
 *-----------------  测试点   ------------------------------
 *1、多重继承情况
 *2、多重继承(合约存在父子关系)
 *3、继承支持传参
 *4、合约函数重载(Overload)
 */


/**
 *
 *1、多重继承:合约可以继承多个合约，也可以被多个合约继承
 *
 *验证： 1)、多重合约继承重名问题，继承顺序很重要，是否遵循最远继承原则
 * ---------------------------------------------------------------
 *验证结果：多重合约继承，如果父类合约有同名函数，则遵循最远继承原则。
 *
 */

contract InheritContractParentOneClass {

    function getDate() public view returns (uint) {
        return 1;
    }
}


contract InheritContractParentTwoClass {

 function getDate() public view returns (uint) {
        return 2;
    }

}


contract InheritContractMutipleTest1 is InheritContractParentTwoClass,InheritContractParentOneClass {

    function callGetDate1() public view returns (uint) {

        return getDate();
    }
}

contract InheritContractMutipleTest2 is InheritContractParentOneClass,InheritContractParentTwoClass {

    function callGetDate2() public view returns (uint) {

        return getDate();
    }
}




