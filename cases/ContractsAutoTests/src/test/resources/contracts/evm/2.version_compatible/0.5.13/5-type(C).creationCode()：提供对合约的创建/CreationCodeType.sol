pragma solidity 0.5.13;
/**
 * 5. type(C).creationCode()：提供对合约的创建
 * https://solidity.readthedocs.io/en/v0.5.3/units-and-global-variables.html
 * 说明:主要用于内联汇编中
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */


contract CreationCodeType {
    function getValue() public returns (string memory){
        return "hello";
    }
}