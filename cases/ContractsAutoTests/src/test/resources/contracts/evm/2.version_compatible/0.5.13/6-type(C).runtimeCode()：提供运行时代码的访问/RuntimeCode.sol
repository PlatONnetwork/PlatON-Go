pragma solidity 0.5.13;
/**
 * 6. type(C).runtimeCode()：提供运行时代码的访问
 * https://github.com/ethereum/solidity/issues/5647
 * https://solidity.readthedocs.io/en/v0.5.3/units-and-global-variables.html
 * 说明:主要用于内联汇编中
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */ 

import "./RuntimeCodeType.sol";
 
contract RuntimeCode {
  
  bytes public constant runtimeCodeInfo = type(RuntimeCodeType).runtimeCode;
  
  function getContractName() public view returns(bytes memory contractName){
      return runtimeCodeInfo;
  }
}