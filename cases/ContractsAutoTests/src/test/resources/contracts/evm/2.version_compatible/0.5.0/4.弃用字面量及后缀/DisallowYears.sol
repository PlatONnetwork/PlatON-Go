pragma solidity ^0.5.0;
/**
 * 弃用字面量及后缀
 * 1-years时间单位已弃用。因为闰年计算容易导致各种问题
 * 2-不允许小数点后不跟数字的数值写法
 * 3-十六进制数字不允许带“0X”前缀，只能用“0x”
 * 4-不允许将十六进制数与单位组合
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract DisallowYears {

    uint256 time1;
    uint256 etherValue;
    uint256 hexValue;
    uint256 hexComValue;


    function testyear(uint a) payable public returns(uint time)  {
        /**
         *1-years时间单位已弃用。因为闰年计算容易导致各种问题
         */
        //   uint _time = 1 years;   //error,编译出错，可以days代替 uint _time = 365 days
        time1 = 365 days; //right,


        /**
         * 2-不允许小数点后不跟数字的数值写法
         */
        //uint  value=255. ether; //error,编译出错
        etherValue=255.0 lat;//right

        /**
         * 3-十六进制数字不允许带“0X”前缀，只能用“0x”
         */
        // uint value=0Xff; //error,编译出错
        hexValue=0xff; //right

        /**
        * 4-不允许将十六进制数与单位组合
        */
        // uint value=0xff ether; //error,编译出错
        hexComValue=0xff*1 lat;//right
    }

    //查询time1的值
    function getTime1() public view returns (uint){
        return time1;
    }

    //查询ether的值
    function getEtherValue() public view returns (uint){
        return etherValue;
    }

    //查询hexValue的值
    function getHexValue() public view returns (uint){
        return hexValue;
    }

    //查询time1的值
    function getHexComValue() public view returns (uint){
        return hexComValue;
    }

}