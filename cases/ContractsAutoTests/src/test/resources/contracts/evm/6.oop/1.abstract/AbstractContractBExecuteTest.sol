pragma solidity 0.5.13;
import "./BaseAbstract.sol";

/**
 * @author qudong
 * @dev 2019/12/23
 * 测试抽象合约功能点
 * 抽象合约(Abstract )简述：
 * 1、合约函数缺少实现，或者除了包含未实现的函数还包含已经实现的函数。
 * 2、如果合约继承自抽象合约，并且没有去重写实现所有未实现的函数，那么它本身依旧是抽象合约。
 *-----------------  测试点   ------------------------------
 * 测试操作：
 * 1、抽象合约被继承，且被实现抽象方法，是否可正常执行
 * 2、普通合约是否可以继承多个抽象合约
 */

/**
 * 1、抽象合约被继承，且实现抽象方法，是否可正常执行
 * -----------------------------------------------
 * 验证结果：抽象合约被继承并实现了抽象方法后，是可以部署执行的
 */
contract AbstractContractBSubclass is AbstractContractParentClass {

    //实现父类抽象函数
    function parentName() public view returns (string memory v){
        return myName;
    }

    function bSubName() public view returns (string memory v) {
        string memory name = "bSubName";
        return name;
    }
}

/**
 * 2、普通合约是否可以继承多个抽象合约,且实现抽象方法，是否可以正常编译部署执行
 * -----------------------------------------------
 *  验证结果：普通合约继承多个抽象合约，可以正常编译执行
 */
 contract AbstractContractCSubclass is AbstractContractASubclass,AbstractContractParentClass {

     //实现ParentClass父类函数
     function parentName() public view returns (string memory v){
          string memory name = "parentName";
          return name;
     }

     //实现ASubclass父类函数
     function aSubName() public view returns (string memory v){
         return subName;
     }

     function setASubName(string memory v) public{
        subName = v;
     }

     function cSubName() public view returns (string memory v) {
        string memory name = "cSubName";
        return name;
    }
 }




  






