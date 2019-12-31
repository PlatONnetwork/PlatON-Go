pragma solidity 0.5.13;


/**
 *
 * 
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约接口功能点
 *
 *接口(interface)简述：接口类似于抽象合约，但是其不能实现任何函数。
 *
 *
 *-----------------  测试点   ------------------------------
 *1、接口限制点测试(5.0以后版本)
 *1)、 可以声明结构体
 *2)、 可以声明枚举
 *
 *2、接口被继承情况测试
 *1)、 普通合约继承单/多个接口
 *2)、 接口无法继承其他合约或接口
 */



interface InterfaceContractEnum {

   /**
    * 验证：1、接口可以声明枚举
    *-----------------------------
    *验证结果：*/

   enum FreshJuiceSize {
        SMALL,
        MEDIUM, 
        LARGE 
   }
}


contract  InterfaceContractTestEnum is InterfaceContractEnum {

    FreshJuiceSize choice;
     //默认饮料类型
    FreshJuiceSize defaultChoice =  FreshJuiceSize.MEDIUM;

   function setLarge() public {
      choice = FreshJuiceSize.LARGE;
   }

   function getChoice() public view returns (FreshJuiceSize) {
      return choice;
   }
   
   function getDefaultChoice() public view returns (uint) {
      return uint(defaultChoice);
   }
}


/*interface InterfaceContractStruct {

   *//**
    * 验证：2、接口可以声明结构体
    *-----------------------------
    *验证结果：
    *//*
   struct Book {
        int bookID;
        string title;
        string describe;
   }
}*/

/*contract  InterfaceContractTestStruct is InterfaceContractStruct {

    Book book;

    function setBook() public {

         book = Book(1,"PlatON","PlatON Describe");
    }

    function getBookID() public view  returns (int) {
         
          return book.bookID;
    }
}*/







