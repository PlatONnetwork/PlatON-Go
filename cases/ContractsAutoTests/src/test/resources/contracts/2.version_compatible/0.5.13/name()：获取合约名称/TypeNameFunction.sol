pragma solidity 0.5.13;

/**
 * type(C).name()：提供对合约名称的访问
 *
 * @author hudenian
 * @dev 2019/12/23 13:57
 */
contract TypeNameFunction {
    function f() public pure returns (string memory) {
        return type(TypeNameFunction).name;
    }
}