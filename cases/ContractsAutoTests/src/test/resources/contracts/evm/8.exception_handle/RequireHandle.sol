pragma solidity ^0.5.13;
/**
 * require(bool condition)函数验证：
 * 如果条件不满足则撤销状态更改，用于检查由输入或者外部组件引起的错误
 *（1）消息调用一个函数，但在调用的过程中，并没有正确结束异常验证
 *（2）使用new创建一个新合约时因为（1）的原因没有正常完成异常验证
 *（3）调用外部函数时，被调用的对象不包含代码异常验证
 *（4）合约没有payable修饰符的public的函数在接收主币时（包括构造函数，和回退函数）异常验证
 *（5）合约通过一个public的getter函数（public getter funciton）接收主币异常验证
 *（6）.transfer()函数执行失败异常验证
 * @author Albedo
 * @dev 2019/12/19
 **/

contract InfoFeed {
    constructor () public {
        revert();
    }
    function info() public payable returns (uint ret) {return 42;}

    function nonCode() public {}
}

contract RequireHandle {
    InfoFeed feed;
    address payable get_pay = address(uint160(0x4B0897b0513fdC7C541B6d9D7E929C4e5364D2dB));
    constructor() public{
        //部署异常
        // address payable to_pay =address(uint160(0x14723A09ACff6D2A60DcdF7aA4AFf308FDDC160C));
        //  to_pay.transfer(10);
    }

    //如果调用require的参数为false
    function paramException(uint param) public {
        require(param < 10);
    }
    /**
     * 如果你通过消息调用一个函数，但在调用的过程中，并没有正确结束
     * (gas不足，没有匹配到对应的函数，或被调用的函数出现异常)。
     * 底层操作如call,send,delegatecall或callcode除外，它们不会抛出异常，但它们会通过返回false来表示失败
     **/
    function functionCallException(uint param) public {
        feed.info.value(10).gas(param)();
    }
    //如果在使用new创建一个新合约时出现第1条的原因没有正常完成
    function newContractException() public {
        feed = new InfoFeed();

    }
    //如果调用外部函数时，被调用的对象不包含代码
    function outFunctionCallException(uint count) public {
        feed.nonCode.gas(count)();
    }
    //如果合约没有payable修饰符的public的函数在接收主币时（包括构造函数和回退函数）
    function nonPayableReceiveEthException(uint count) public {
        address payable to_pay =address(uint160(0x14723A09ACff6D2A60DcdF7aA4AFf308FDDC160C));
        to_pay.transfer(count);
    }
    //如果合约通过一个public的getter函数（public getter funciton）接收主币
    function publicGetterReceiveEthException(uint count) public {
        get_pay.transfer(count);
    }
    //如果.transfer()执行失败
    function transferCallException(uint count) public payable {
        address payable to_pay =address(uint160(0x14723A09ACff6D2A60DcdF7aA4AFf308FDDC160C));
        to_pay.transfer(count);
    }

    function() external {
        address payable to_pay =address(uint160(0x14723A09ACff6D2A60DcdF7aA4AFf308FDDC160C));
        to_pay.transfer(10);
    }
}