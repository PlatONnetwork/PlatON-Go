pragma solidity 0.5.13;


/**
 *
 * @author qudong
 * @dev 2019/12/23
 *
 *
 *测试抽象合约功能点
 *
 *抽象合约(Abstract )简述：
 *1、合约函数缺少实现，或者除了包含未实现的函数还包含已经实现的函数。
 *2、如果合约继承自抽象合约，并且没有去重写实现所有未实现的函数，那么它本身依旧是抽象合约。
 *
 *-----------------  测试点   ------------------------------
 *情况分类：
 *1、不含任何实现的抽象合约
 *2、包含部分实现的抽象合约
 *测试操作：
 *1、抽象合约是否可编译、部署、执行
 *2、抽象合约被继承，但未被实现抽象方法，是否可正常执行
 *
 *
 */




/**
  *
  * 1、测试：不含任何实现的抽象合约，是否可编译部署执行
  * ------------------------------------------------
  * 验证结果：抽象合约可以编译、部署，但是不可以执行调用方法
  */
contract AbstractContractGrandpa {

    function name() public view returns (string memory v);

}


/**
 *  1.1、测试：包含部分实现的抽象合约，可编译部署执行
 * ------------------------------------------------
 *  验证结果：抽象合约可以编译、部署，但是不可以执行调用方法
 */
contract AbstractContractFather {

    function fatherName() public view returns (string memory v);

    function fatherAge() public view returns (int v) {
        int age = 20;
        return age;
    }
}


/**
 * 2、测试：抽象合约被继承，但未被实现抽象方法，是否可正常执行
 * ------------------------------------------------
 *   验证结果：抽象合约被继承，但未被实现抽象方法,可以编译、部署，但是不可以执行调用方法
 *
 */

contract AbstractContractSon is AbstractContractFather {

    function sonName() public view returns (string memory v) {
        string memory name = "sonName";
        return name;
    }

}






