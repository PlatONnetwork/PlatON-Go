pragma solidity ^0.6.12;

/**
 * @author qudong
 * @dev 2019/12/23
 * 基础合约
 */
abstract contract AbstractContractParentClass {

    string myName =  "";
    function parentName() public virtual returns (string memory v);

    function setParentName(string memory name) public {
        myName = name;
    }
}

abstract contract AbstractContractASubclass {

    string subName =  "";
    function aSubName() public virtual returns (string memory v);
    function aSubAge() public view returns (int v) {
        int age = 20;
        return age;
    }
}


