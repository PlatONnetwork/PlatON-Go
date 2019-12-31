pragma solidity 0.5.13;

/**
 * type(C).name()：提供对合约名称的访问
 *
 * @author hudenian
 * @dev 2019/12/23 13:57
 */
contract TypeName {
  string public constant name = type(typeName).name;
  
  function getContractName() public returns(string memory contractName){
      return name;
  }
}