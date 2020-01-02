pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试引用数据类型
 *1、数组（Array）
 *2、结构体（Struct）
 */

contract ReferenceDataTypeStructContract {

  /**
   * 2、结构体 (struct)
   * 简述：使用关键字struct定义结构体，struct可以用于映射和数组中作为元素，
   * 其本身也可以包含映射和数组类型，不能声明一个struct同时将自身struct作为成员。
   * -----------------  测试点   ------------------------------
   *1)、定义不同子元素数据（包含基本类型、包含引用数据类型）
   *2)、赋值及取值
   */
  

  /**
   *验证：1)、定义不同子元素数据（包含基本类型、包含引用数据类型）
   *---------------------------------------------------------
   * 验证结果：结构体可以包含基本数据类型，引用数据类型，包含结构体（不可以包含自身结构体，否则会编译不过）
   * 
   */

  //定义结构体包含基本数据类型
    struct  Person {
        uint id;
        address addr;
        uint age;
        bool isVIP;
    }


    //定义结构体包含引用数据类型
    struct Group {
        uint idGroup;
        address addrGroup;
        mapping(uint => address) addressMap;
        string nameGroup;
        uint[] pIdArray;
    }

    //定义结构体包含结构体(但是结构体不能包含自身)
    struct LucyPerson {
        uint id;
        string name;
        Person p; //包含其他结构体
       // LucyPerson lucy; //结构体不能包含自身，否则编译不过
    }

  /**
   * 2)、验证结构体赋值及取值
   * ---------------------
   * 三种赋值方式 
   */

   //赋值方式一： 按入参顺序赋值
   function initDataStruct1() public view returns (uint,address,uint,bool) {
      
      Person memory  person = Person(2,msg.sender,25,true);
      return (person.id,person.addr,person.age,person.isVIP);
   }

  //赋值方式二： 按命名参数赋值
   function initDataStruct2() public view returns (uint,address,uint,bool) {
      
      Person memory  person = Person({id:2,addr:msg.sender,age:25,isVIP:true});
      return (person.id,person.addr,person.age,person.isVIP);
   }

  //赋值方式三：结构体中映射的初始化

 function initDataStruct3() public view returns (uint,address,uint,bool) {
      
      Person memory p1;//分配映射空间
      Person memory  person = Person({id:2,addr:msg.sender,age:25,isVIP:true});
      p1 = person;
      return (p1.id,p1.addr,p1.age,p1.isVIP);
   }
}



