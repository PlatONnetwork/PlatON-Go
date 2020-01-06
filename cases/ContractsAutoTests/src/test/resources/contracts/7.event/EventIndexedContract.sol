pragma solidity ^0.5.13;

/**
 *  事件声明类型验证
 * （1）一维数组索引
 * （2）二维数组但一维确定长度 索引
 * （3）字符串索引
 * （4）枚举索引
 * （4）复杂索引
 * （5）匿名事件索引认证
 * @author Albedo
 * @dev 2020/01/06
 **/
contract EventIndexedContract {
    uint[] uintArray = [0,1,2,3,4];
    //一维数组索引测试
    event OneDimensionalArrayEvent(uint[] indexed array);
    function testOneDimensionalArray() public{
        emit OneDimensionalArrayEvent(uintArray);
    }
    uint[2][]  multiArray = [[0,0],[0,1],[0,2],[0,3],[0,4],[0,5]];

    event TwoDimensionalArrayEvent(uint[2][] indexed array);
    //二维数组索引
    function testTwoDimensionalArray() public{
        emit TwoDimensionalArrayEvent(multiArray);
    }
    string str = "1234567890097865432112345678900987654321123456789009764354666663242444444444475831546856";
    event StringEvent(string indexed str);
    //string类型索引
    function testStr() public{
        emit StringEvent(str);
    }
    enum ActionChoices { GoLeft, GoRight, GoStraight, SitStill }
    event EnumEvent(ActionChoices indexed choices);

    //测试枚举类型索引
    function testEnum() public{
        emit EnumEvent(ActionChoices.GoLeft);
    }

    event ComplexIndexedEvent(uint[] indexed array,ActionChoices indexed choice,string indexed str);
    //复杂类型索引
    function testComplex() public{
        emit ComplexIndexedEvent(uintArray,ActionChoices.GoLeft,str);
    }

    event AnonymousIndexedNum(uint indexed u1,uint indexed u2,uint indexed u3,uint indexed u4) anonymous;
    //匿名事件索引数量测试
    function testAnonymousIndexed() public{
        emit AnonymousIndexedNum(1,2,3,4);
    }
}