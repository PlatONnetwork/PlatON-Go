pragma solidity ^0.5.0;
/**
 * 05-弃用变量
 * 1-不允许声明0长度的定长数组类型
 * 2-不允许声明0结构体成员的结构体类型
 * 3-不允许使用空元组组件
 * 4-不允许声明未初始化的storage变量
 * 5-不允许定义非编译期常量的 constant常量
 * 6-var已弃用，需要明确声明变量类型
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract DisallowVariable {
    
    /**
      *1-不允许声明0长度的定长数组类型
      */ 
    //  int[0]  arri0; //error,编译报错
     int[1] arri1; //right


      /**
       * 2-不允许声明0结构体成员的结构体类型
       */ 
    // struct Person{}//error,编译出错
           
      struct Person{
         uint age;
      }
      
      
      mapping (uint =>Person) persons;

   function tesEmpty(uint _id) payable public returns(uint time)  {
     
     /**
       * 3-不允许使用空元组组件
       */ 
    //  (uint x,uint y) = (,); //error,编译出错
     (uint x,uint y) = (2,3); //right

       /**
        * 4-不允许声明未初始化的storage变量
        */ 
        // Person storage p;//error,编译出错
        Person storage p = persons[_id];//right
   }
   
   
   /**
    * 5-不允许定义非编译期常量的 constant常量
    */ 
    // uint constant time = now; //error,编译出错
    
    /**
    * 6-var已弃用，需要明确声明变量类型
    */ 
    // var i=1; //编译出错，需明确指定变量类型
    uint j=1;

}