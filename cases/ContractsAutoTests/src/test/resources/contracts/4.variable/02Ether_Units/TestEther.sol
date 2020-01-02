/*************************************************
* Ether的单位关键字有wei, finney, szabo, ether，换算格式如下：
* 1 ether = 1 * 10^18 wei
* 1 ether = 1 * 10^6 szabo
* 1 ether = 1* 10^3 finney
* 默认缺省单位是wei
*************************************************/

pragma solidity 0.5.13;
/**
 * 对 比特币 Ether 的几个单位进行测试
 */
contract test {
    // 定义全局变量
    uint public balance;

    function testEther() public{
        balance = 1 ether;  //1000000000000000000
    }

    function fFinney() public{
      balance = 1 finney; //1000000000000000
    }

    function fSzabo() public{
      balance = 1 szabo;  //1000000000000
    }

    function fWei() public{
      balance = 1 wei; //1
    }
}
