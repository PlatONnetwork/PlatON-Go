pragma solidity 0.5.13;
import "./BaseAbstractInterface.sol";

/**
 * @author qudong
 * @dev 2019/12/23
 * 基础合约（定义抽象合约继承接口）
 */

contract AbstractContractESubclass is AbstractContractAInterface {
    function setInterAge(int v) public;
}
