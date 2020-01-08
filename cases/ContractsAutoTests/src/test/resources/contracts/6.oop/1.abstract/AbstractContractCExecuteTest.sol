pragma solidity 0.5.13;
import "./BaseAbstract.sol";
import "./BaseAbstractParent.sol";

/**
 * @author qudong
 * @dev 2019/12/23
 * 测试抽象合约功能点
 * 抽象合约(Abstract )简述：
 * 1、合约函数缺少实现，或者除了包含未实现的函数还包含已经实现的函数。
 * 2、如果合约继承自抽象合约，并且没有去重写实现所有未实现的函数，那么它本身依旧是抽象合约。
 *-----------------  测试点   ------------------------------
 * 测试操作：
 * 1、抽象合约是否可以继承抽象合约
 * 2、抽象合约是否可以继承接口（但是反之不可以）
 */

/**
 * 1、抽象合约是否可以继承抽象合约
 * ----------------------------------
 * 验证结果： 抽象合约是可以继承抽象合约
 */
 contract AbstractContractFSubclass is AbstractContractDSubclass {

     function parentName() public view returns (string memory v){
           return myName;
     }

     function dSubClassName() public view returns (string memory dSubName){

     }
 }


























