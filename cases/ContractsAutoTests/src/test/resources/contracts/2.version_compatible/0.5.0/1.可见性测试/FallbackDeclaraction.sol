pragma solidity ^0.5.0;
/**
 * fallback函数必须声明为external
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract FallbackDeclaraction{

    uint public a;

    //回退事件，会把调用的数据打印出来
    event FallbackCalled(bytes data);
    //fallback函数，注意是没有名字的，没有参数，没有返回值的
    function() external{
        a =111;
        emit FallbackCalled(msg.data);
    }

    //调用已存在函数的事件，会把调用的原始数据，请求参数打印出来
    event ExistFuncCalled(bytes data, uint256 para);
    //一个存在的函数
    function existFunc(uint256 para) public{
        emit ExistFuncCalled(msg.data, para);
    }


    //模拟从外部对一个不存在的函数发起一个调用，由于匹配不到函数，将调用回退函数
    function callNonExistFunc() public returns(uint){
        address(this).delegatecall("functionNotExist()");
        return 1;
    }

    function getA() public view returns(uint){
        return a;
    }
}
