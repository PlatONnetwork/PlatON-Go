pragma solidity 0.5.13;

/**
 * type(C).name()：获取合约名称
 *
 * @author hudenian
 * @dev 2019/12/23 13:57
 */
contract TypeName {
  string public constant name = type(TypeName).name;
  
  function getContractName() public returns(string memory contractName){
      return name;
  }
}