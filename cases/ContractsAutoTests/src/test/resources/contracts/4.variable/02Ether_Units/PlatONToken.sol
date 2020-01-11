pragma solidity 0.5.12;

/**************************************************
* PlatON的单位关键字有von, kvon, mvon, gvon，szabo,finney,lat,klat,mlat,glat换算格式如下：
* 1 glat = 1 * 10^27 von
* 1 mlat = 1 * 10^24 von
* 1 klat = 1* 10^21 von
* 1 lat = 1* 10^18 von
* 1 finney = 1* 10^15 von
* 1 szabo = 1* 10^12 von
* 1 gvon = 1* 10^9 von
* 1 mvon = 1* 10^6 von
* 1 kvon = 1* 10^3 von
* 默认缺省单位是von
*************************************************/

// 对 PlatON  币的几个单位进行测试
contract PlatONToken {
    // 定义全局变量
    uint public balance;

    function Plat() public returns(uint balance){
        balance = 1 lat; //1000000000000000000
    }

    function Pfinney() public returns(uint balance){
        balance = 1 finney; //1000000000000000
    }

    function Pszabo() public returns(uint balance){
        balance = 1 szabo; //1000000000000
    }

    function Pvon() public returns(uint balance){
        balance = 1 ; //1
    }

}