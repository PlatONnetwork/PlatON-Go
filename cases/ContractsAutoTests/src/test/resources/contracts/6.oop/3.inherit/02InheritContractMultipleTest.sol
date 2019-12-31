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
 *
 *-----------------  测试点   ------------------------------
 *1、连续继承情况 
 *2、多重继承情况
 *3、多重继承(合约存在父子关系)
 *4、继承支持传参
 *5、合约函数重载(Overload)
 */



/**
 *
 *1、多重继承:合约可以继承多个合约，也可以被多个合约继承
 *
 * 验证：2)、当继承多个合约时，这些父合约中不允许出现相同的函数名，事件名，修改器名，或互相重名。
 *        另外，隐藏情况，默认状态变量的getter函数导致的重名。
 *--------------------------------------------
 *
 */

contract Base1 {
    address owner;

    modifier test(){
        if (msg.sender == owner)
            _;
    }


    function test1() public view {

    }
}

contract Base2  {
    address owner;

    modifier ownd() {
        if (msg.sender == owner )
            _;
    }


    function test() public view {

    }
}


contract InheritContractMutipleTest1 is Base2,Base1 {

    function sum() public view {

    }

}




