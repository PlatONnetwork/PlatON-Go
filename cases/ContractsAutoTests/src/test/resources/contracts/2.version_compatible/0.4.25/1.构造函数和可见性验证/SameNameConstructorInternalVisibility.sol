pragma solidity ^0.4.25;

/**
 * 1. 0.4.25版本同名函数构造函数定义，声明internal可见性验证；
 * 2. 0.4.25版本接口(interface)函数支持external和public两种可见性，可见性声明非必须验证
 * （1）默认可见性（默认public）函数声明
 * （2）public可见性函数声明
 * （3）external可见性声明
 * 3. 0.4.25版本支持，但0.5.x已弃用变量验证
 * (1)0.4.25版本允许声明0长度的定长数组类型
 * (2)0.4.25版本允许声明0结构体成员的结构体类型
 * (3)0.4.25版本允许定义非编译期常量的 constant常量
 * (4)0.4.25版本允许使用空元组组件
 * (5)0.4.25版本允许声明未初始化的storage变量
 * (6)0.4.25版本允许使用var
 * @author Albedo
 * @dev 2019/12/19
 **/
interface VisibilityInterface {
    //默认可见性
    function defaultVisibility();
    //public可见性
    function publicVisibility() public;
    //external可见性
    function externalVisibility() external;
}

contract SameNameConstructorInternalVisibility {
    uint256 conParam;
    //声明0长度的定长数组类型
    int256[0] zeroArr;
    //声明0结构体成员的结构体类型
    struct ZeroStruct {}
    //允许定义非编译期常量的 constant常量
    uint constant time = 1;

    struct UnInitialized {
        int age;
        string name;
    }
    //同名函数构造函数，internal可见性
    function SameNameConstructorPublicVisibility(uint256 param) internal {
        conParam = param;
    }
    //弃用字面量及后缀整体覆盖验证
    function discardVariable() public view returns (uint, uint, uint){
        //允许使用空元组组件验证
        uint x;
        uint y;
        (x,) = (1,);
        //允许声明未初始化的storage变量
        UnInitialized storage unInitialized;
        //使用var
        var i = 1;
        return (x, y, i);
    }

}
