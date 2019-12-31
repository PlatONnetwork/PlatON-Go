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
 *1、连续继承情况 
 *2、多重继承情况
 *3、多重继承(合约存在父子关系)
 *4、继承中基类构造函数的传参
 *5、合约函数重载(Overload)
 *6、继承中涉及函数及状态变量访问范围（可见性测试）
 */




/**
*1、连续继承
* 测试：连续继承重名问题(状态变量的名字、函数的名字都是相同情况下，C合约状态变量及函数以继承的B合约中的状态变量和函数为准)
*
*/


contract A{


}

contract B is A{


}


contract C is B{


}

/**
 *
 *2、多重继承:合约可以继承多个合约，也可以被多个合约继承
 *测试：
 *1)、多重合约继承重名问题，继承顺序很重要，以最后继承的为主(Subclass函数以ParentTwoClass为准)。
 *2)、当继承多个合约时，这些父合约中不允许出现相同的函数名，事件名，修改器名，或互相重名。
 *    另外，隐藏情况，默认状态变量的getter函数导致的重名。
 */

contract ParentOneClass{


}
contract ParentTwoClass{



}
contract Subclass is ParentOneClass,ParentTwoClass{


}


/**
 *3、 多重继承(合约存在父子关系)：如果继承的合约之间有父子关系，那么合约要按照先父到子的顺序排序
 *测试：合约中继承顺序
 *
 */

contract ParentClass1{

}

contract ASubclass is ParentClass1{

}

contract BSubclass is ParentClass1,ASubclass{


}

/**
 *4、继承中基类构造函数的传参
 *
 */


//方式一
contract Base {
    uint x;
    function Base(uint _x) public { x = _x; }
}
//方式二
contract Derived is Base(7) {
    function Derived(uint _y) Base(_y * _y) public {
    }
}


/**
 *5、合约函数重载(Overload)：合约可以有多个同名函数，可以有不同输入参数。
 *   重载解析和参数匹配
 *
 *
 **/

contract A {

    function f(uint _in) public pure returns (uint out) {
        out = 1;
    }

    function f(uint _in, bytes32 _key) public pure returns (uint out) {
        out = 2;
    }
}


/**
 *6、继承中涉及函数及状态变量访问范围（可见性测试）
 *   external、public、internal、private
 */

