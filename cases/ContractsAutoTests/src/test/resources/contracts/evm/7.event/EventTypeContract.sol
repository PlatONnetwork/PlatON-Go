pragma solidity ^0.5.13;

/**
 *  事件声明类型验证
 * （1）一维数组
 * （2）二维数组但是二维确定长度
 * （3）字符串类型
 * （3）枚举类型
 * @author Albedo
 * @dev 2020/01/06
 **/
contract EventTypeContract {

    uint[] uintArray = [0,1,2,3,4];
    event OneDimensionalArrayEvent(uint[] array);
    //测试一维数组类型声明事件
    function testOneDimensionalArray() public{
        emit OneDimensionalArrayEvent(uintArray);
    }
    uint[2][]  multiArray = [[0,0],[0,1],[0,2],[0,3],[0,4],[0,5]];

    event TwoDimensionalArrayEvent(uint[2][] array);
    //测试二维数组声明事件
    function testTwoDimensionalArray() public{
        emit TwoDimensionalArrayEvent(multiArray);
    }
    string str = "1234567890097865432112345678900987654321123456789009764354666663242444444444475831546856";
    event StringEvent(string str);
    //测试string类型声明事件
    function testStr() public{
        emit StringEvent(str);
    }

    enum ActionChoices { GoLeft, GoRight, GoStraight, SitStill }
    event EnumEvent(ActionChoices choices);

    //测试枚举类型声明事件
    function testEnum() public{
        emit EnumEvent(ActionChoices.GoLeft);
    }

}