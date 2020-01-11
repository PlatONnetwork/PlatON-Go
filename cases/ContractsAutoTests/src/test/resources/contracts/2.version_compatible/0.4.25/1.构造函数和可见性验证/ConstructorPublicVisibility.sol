pragma solidity ^0.4.17;

/**
 * 1. constructor声明构造函数，public可见性验证；
 * 2. 允许调用没有括号的基类构造函数验证
 * 3. ，0.4.25支持语法(0.5.X版本已弃用)覆盖验证
 * （1）0.4.25版本允许未实现的函数使用修饰符(modifier)验证()
 * （2）0.4.25版本允许布尔变量使用算术运算验证
 * （3）0.4.25版本允许使用一元运算符"+"验证
 * （4）0.4.25版本允许在if包含的块中使用单个语句声明/定义变量验证
 * （5）0.4.25版本允许在while包含的块中使用单个语句声明/定义变量验证使用算术原酸运算符
 * （6）0.4.25版本允许在for包含的块中使用单个语句声明/定义变量验证
 * （7）0.4.25版本允许具有一个或多个返回值的函数使用空返回语句验证
 * （8）0.4.25版本允许具有一个或多个返回值的函数使用空返回语句验证
 *  4. 0.4.25版本允许constant用作修饰函数状态可变性验证
 *  5. 0.4.25版本允许定义具有命名返回值的函数类型验证
 *  6. 0.4.25版本允许 msg.value用在非 payable函数里以及此函数的修饰符(modifier)里验证
 * @author Albedo
 * @dev 2019/12/23
 **/

//0.4.25版本允许未实现的函数使用修饰符(modifier)验证
contract BasePublic {
    uint x;
    modifier preMethod {
        x = 123;
        _;
    }
    //允许未实现的函数使用修饰符(modifier)
    function abstractFunction() public view preMethod returns (uint);

}

contract ConstructorPublicVisibility is BasePublic {
    uint outI;
    constructor(uint _y) public {outI = _y;}
    function abstractFunction() public view preMethod returns (uint){
        outI = outI + x;
        return outI;
    }

    function getOutI() public view returns (uint) {
        return outI;
    }
    //弃用语法验证
    function grammarCheck() public view returns (uint, uint){
        //0.4.25版本允许布尔变量使用算术运算验证
        bool evenFlag = true;
        bool evenFlag1 = false;
        //实际验证不可行
        // bool result = evenFlag + evenFlag1;
        //0.4.25版本允许使用一元运算符"+"验证
        uint value = + 1;
        //0.4.25版本允许在if包含的块中使用单个语句声明/定义变量验证
        if (evenFlag) {uint ifVar = 1;}
        //0.4.25版本允许在while包含的块中使用单个语句声明/定义变量验证使用算术原酸运算符
        while (evenFlag1) {uint whileVal = 1;}
        //0.4.25版本允许在for包含的块中使用单个语句声明/定义变量验证
        for (uint i = 0; i < 2; i++) {uint forVal = 1;}
        //0.4.25版本允许具有一个或多个返回值的函数使用空返回语句验证
        return;
    }

    //0.4.25版本允许constant用作修饰函数状态可变性验证
    function constantCheck() public constant returns (uint) {
        return outI;
    }


    //0.4.25版本允许定义具有命名返回值的函数类型验证
    function namedReturn(uint a, uint b) public view returns (uint add){
        add = a + b;
    }

    //0.4.25版本允许 msg.value用在非 payable函数里以及此函数的修饰符(modifier)里验证
    function show() public view returns (uint) {return msg.value;}
}
