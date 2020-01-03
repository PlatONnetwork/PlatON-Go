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
 * 验证：2)、当继承多个合约时，这些父合约中不允许出现相同的函数名，事件名，修改器名，或互相重名。
 *        另外，隐藏情况，默认状态变量的getter函数导致的重名。
 *--------------------------------------------
 *验证结果：当继承多个合约时，父类重名函数与修饰器、事件与修饰器、默认getter函数名同名，会编译异常
 */

contract InheritContractParentBase1 {
    address public owner;//编译器自动为所有public 状态变量创建getter函数
    modifier ownd1(){
        if (msg.sender == owner)
            _;
    }

    function getAddress1() public view returns (address) {
         return msg.sender;
    }
}

contract InheritContractParentBase2 {
    address owner2;
    modifier ownd2() {
        if (msg.sender == owner2)
            _;
    }

    function getAddress2() public view returns (address) {
        return msg.sender;
    }
   //1、同名函数名与修饰器,编译异常
   /* function ownd1() public view returns (address) {
         return msg.sender;
    }*/

    //2、同名修饰器与事件,编译异常
   // event ownd1(address addr);

    //3、默认状态变量 owner的getter函数导致的重名,编译异常
  /*  function owner() public view returns (address) {
         return msg.sender;
    }*/
}

contract InheritContractParentMutipleTest is InheritContractParentBase2,InheritContractParentBase1 {
    function sum() public view {

    }
}




