pragma solidity 0.5.13;

/**
 * @author liweic
 * @dev 2020/01/11
 * -----------------  测试点   ------------------------------
 * 验证合约数组
 */
contract ContractArray {
    ContractArray[] y = new ContractArray[](3);
    ContractArray[3] x;
    function f() public {
        ContractArray[3] memory z;
        y.push(this);
        x[0] = this;
        z[0] = this;
    }

    function gety() view public returns(ContractArray[] memory){
        return y;
    }

    function getx() view public returns(ContractArray){
        return x[0];
    }
}