pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 *测试合约接口功能点
 *接口(interface)简述：接口类似于抽象合约，但是其不能实现任何函数。
 *-----------------  测试点   ------------------------------
 *1、接口被继承情况测试
 *1)、 普通合约继承多个接口
 *2)、 接口无法继承接口
 *3)、 接口无法继承其他合约(在抽象函数已经验证过此问题，估此处不再验证)
 */
interface InterfaceContractInheritOne {
      function sum(uint a, uint b) external view returns (uint);
}

interface InterfaceContractInheritTwo {
      function reduce(uint c, uint d) external view returns (uint);
}

/**
 *验证：1)、普通合约是否可以继承多个接口
 *-----------------------------------
 *验证结果：普通合约可以继承多个接口
 */
contract InterfaceContractInheritMultipleTest is InterfaceContractInheritOne,
                                                   InterfaceContractInheritTwo {

    function sum(uint a, uint b) external view returns (uint) {
         return a + b;
    }

    function reduce(uint c, uint d) external view returns (uint) {
         return c - d;
    }


     
}

/**
 * 验证：2)、 接口是否可以继承接口
 * -------------------------------
 * 验证结果：接口是不可以继承接口的，无法编译通过
 */

/* interface InterfaceContractInheritTest is InterfaceContractInheritOne {

    function multiply(uint e, uint f) external returns (uint); 
 }*/







