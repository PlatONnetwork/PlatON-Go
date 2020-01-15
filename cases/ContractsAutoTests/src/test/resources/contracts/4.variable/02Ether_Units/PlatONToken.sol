pragma solidity 0.5.13;

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
    uint public platontoken;

    function Token() public{
        platontoken = 1 von;
    }

    function Plat() public view returns(uint platontoken){
        //1lat = 1000000000000000000
        return platontoken + 1 lat;
    }

    function Pfinney() public view returns(uint platontoken){
        //1finney = 1000000000000000
        return platontoken + 1 finney;
    }

    function Pszabo() public view  returns(uint platontoken){
        //1 szabo = 1000000000000
        return platontoken + 1 szabo;
    }

    function Pvon() public view returns(uint platontoken){
        //默认缺省单位是von
        return platontoken + 1;
    }

}