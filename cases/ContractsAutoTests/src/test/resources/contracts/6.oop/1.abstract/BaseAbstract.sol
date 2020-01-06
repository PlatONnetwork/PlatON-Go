pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 * 基础合约
 */
contract AbstractContractParentClass {

    function parentName() public view returns (string memory v);
}

contract AbstractContractASubclass {

    function aSubName() public view returns (string memory v);
    function aSubAge() public view returns (int v) {
        int age = 20;
        return age;
    }
}
