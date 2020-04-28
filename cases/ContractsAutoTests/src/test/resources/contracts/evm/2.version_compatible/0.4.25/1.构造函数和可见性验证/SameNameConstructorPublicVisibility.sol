pragma solidity ^0.4.20;
/**
 * 1. 0.4.25版本同名函数构造函数定义，声明public可见性验证（0.5.x版本弃用同名函数定义构造函数）；
 * 2. 0.4.25版本fallback函数可见性非强制声明（默认public可见性）验证
 * 3. 0.4.25版本支持字面量及后缀（0.5.x版本已弃用）验证
 * （1）0.4.25版本支持year时间单位
 * （2）0.4.25版本允许小数点后不跟数字的数值写法
 * （3）0.4.25版本十六进制数字支持带“0X”和“0x”等2种前缀表示
 * （4）0.4.25版本支持十六进制数与以太币单位组合
 *
 * @author Albedo
 * @dev 2019/12/19
 **/
contract SameNameConstructorPublicVisibility {
    uint256 conParam;
    //同名函数构造函数，public可见性
    function SameNameConstructorPublicVisibility(uint256 param) public   {
        conParam=param;
    }
    //弃用字面量及后缀整体覆盖验证
    function discardLiteralsAndSuffixes() public view returns (uint,uint,uint,uint,uint,uint){
        //0.4.25版本years时间单位支持验证
        uint time = 1 years;
        //0.4.25版本允许小数点后不跟数字的数值写法验证
        uint ethPoint=255. lat;
        //0.4.25版本十六进制数字支持带“0X”和“0x”等2种前缀表示验证
        uint hexX=0Xff;
        uint hexx=0xff;
        //0.4.25版本支持十六进制数与以太币单位组合验证
        uint hexEther=0xff lat;
        return (conParam,time,ethPoint,hexX,hexx,hexEther);
    }


    //fallback函数可见性非强制声明验证
    function (){

    }
}
