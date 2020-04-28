pragma solidity ^0.5.0;
/**
 * 06-弃用语法
 * 1-不允许布尔表达式使用算术运算符
 * 2-不允许具有一个或多个返回值的函数使用空返回语句
 * 3-不允许未实现的函数使用修饰符(modifier)
 * 4-constant函数修饰符已弃用，不能用作修饰函数状态可变性，请使用view关键字
 * 5-不允许使用一元运算符 +
 * 6-不允许在if/while/for等非{}包含的块中使用单个语句声明/定义变量
 * 7-不允许定义具有命名返回值的函数类型
 * 8-不允许 msg.value用在非 payable函数里以及此函数的修饰符(modifier)里
 * 9-新的关键字alias,apply,auto,copyof,define,immutable,implements,macro,mutable,override,partial,promise,reference,sealed,sizeof,supports,typedefandunchecked，但这些关键字并未实际声明用处，会在以后的版本进行设置
 * 10-接口允许定义结构体与枚举
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract DisallowSyntax {
    
     /**
      *1-不允许布尔表达式使用算术运算符
      */ 
     bool flg = false;
     uint a = 1;
    //  uint b = a + flg;


   /**
    * 2-不允许具有一个或多个返回值的函数使用空返回语句
    */     
   function tesReturn(uint _id) payable public returns(uint res1)  {
      // return;//error,编译出错
      return 1; //right
   }
   
   /**
    * 3-不允许未实现的函数使用修饰符(modifier)
    */
//   function method() public modifier; //error,编译报错
   function method() public; //right
   
   /**
    * 4-constant函数修饰符已弃用，不能用作修饰函数状态可变性，请使用view关键字
    */
    // function getValue(address _to, uint _value) constant external returns (uint){ //error,编译出错，请使用view或者pure
    function getValue(address _to, uint _value) view external returns (uint){ //right
        return a;
    }
    
    /**
    * 5-不允许使用一元运算符 +
    */
    // uint value = +1 ether; //编译出错，可以写成 value = 1ether
    uint value1 = 1 lat; //编译出错，可以写成 value = 1ether
    
    
    /**
    * 6-不允许在if/while/for等非{}包含的块中使用单个语句声明/定义变量
    */
    function testBlock(address _to, uint _value) public{
        // if(true) bool b=false; //error,编译出错，必须用 {}包含单个语句
        if(true) {bool b=false;}
    }
    
     /**
      * 7-不允许定义具有命名返回值的函数类型
      * 函数执行时出错
      */
    function mulvalue2(uint a,uint b) public returns(uint add){
       add = a+b;
    }
    
    /**
      * 8-不允许 msg.value用在非 payable函数里以及此函数的修饰符(modifier)里
      * 
      */
    //   function show() public returns(uint) { return msg.value; } //error编译出错，需加payable修饰符 
      function getMsgValue() public payable returns(uint) {return msg.value; }//right
    
    /**
      * 9-新的关键字alias,apply,auto,copyof,define,immutable,implements,macro,mutable,override,partial,promise,reference,sealed,sizeof,supports,typedefandunchecked，
      * 但这些关键字并未实际声明用处，会在以后的版本进行设置
      * 
      */
      
}