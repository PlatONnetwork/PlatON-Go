pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 * 基础合约
 */
contract AbstractContractParentClass {

    string myName =  "";
    function parentName() public view returns (string memory v);

    function setParentName(string memory name) public {
        myName = name;
    }
}

contract AbstractContractASubclass {

    string subName =  "";
    function aSubName() public view returns (string memory v);
    function aSubAge() public view returns (int v) {
        int age = 20;
        return age;
    }
}


