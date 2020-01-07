pragma solidity 0.5.13;
import "./BaseAbstract.sol";

/**
 * @author qudong
 * @dev 2019/12/23
 *测试抽象合约功能点
 *抽象合约(Abstract )简述：
 *1、合约函数缺少实现，或者除了包含未实现的函数还包含已经实现的函数。
 *2、如果合约继承自抽象合约，并且没有去重写实现所有未实现的函数，那么它本身依旧是抽象合约。
 *-----------------  测试点   ------------------------------
 *测试操作：
 *1、抽象合约是否可以继承抽象合约
 *2、抽象合约是否可以继承接口（但是反之不可以）
 */

/**
 * 1、抽象合约是否可以继承抽象合约
 * ----------------------------------
 * 验证结果： 抽象合约是可以继承抽象合约
 */
 contract AbstractContractDSubclass is AbstractContractParentClass {
     function parentName() public view returns (string memory v){
         string memory parentName = "parentName";
         return parentName;
     }

     function dSubName() public view returns (string memory v) {
         string memory name = "dSubName";
         return name;
     }
 }

/**
 * 2、抽象合约是否可以继承接口(反之接口是否可以继承抽象合约)
 * -----------------------------
 * 测试结果：抽象合约是可以继承接口（但是反过来接口是不可以继承抽象合约）
 */

interface AbstractContractAInterface{
     function aInterAge() external view returns (int age);
}

//1)、接口是不可以继承抽象合约（编译报错）
/* interface AbstractContractBInterface is AbstractContractParentClass{
      function bInterAge() external returns (int  age);
 }*/

//2)、抽象合约继承接口
contract AbstractContractESubclass is AbstractContractAInterface {

    function aInterAge() external view returns (int age){
        int age = 15;
        return age;
    }

    function eSubName() public view returns (string memory v){
        string memory name = "eSubName";
        return name;
    }
}






















