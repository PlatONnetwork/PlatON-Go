pragma solidity ^0.5.13;

contract PlatonUnit {

    uint256 public balance; //合约金额

    /**
     * 默认函数
     *
     * 默认函数，可以向合约直接打款
     */
    function () payable external {
        balance = add(balance,msg.value);
    }


    //累加函数
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        assert(c >= a);
        return c;
    }

    //查看当前合约余额
    function getBalance() view public returns (uint256){
        return address(this).balance;
    }

}