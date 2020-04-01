pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 *
 *测试结构体嵌套delete操作
 */

contract ReferenceDataTypeStructDeleteContract {

     //内嵌结构体
     struct nestedStruct {
          uint nestedValue;
          mapping (uint => bool) nestedMapping; 
     }
     //外部结构体
     struct topStruct {
          nestedStruct nStruct;//声明内嵌结构体
          uint topValue;
          mapping (uint => uint) topMapping;
     }
     uint toDeletInt;
     topStruct tStruct;//声明内嵌结构体
     //构造函数初始化
     constructor() public{
         toDeletInt = 5;
         //外部结构体赋值
         tStruct.topValue = 10;
         tStruct.topMapping[1] = 1;
         tStruct.topMapping[2] = 2;
          //外部结构体嵌套的内部结构体赋值
         tStruct.nStruct.nestedValue = 3;
         tStruct.nStruct.nestedMapping[1] = true;
         tStruct.nStruct.nestedMapping[2] = true;
         //执行delete
         delete tStruct;
         delete toDeletInt;
     }

     function getToDeleteInt() public view returns (uint) {
          return toDeletInt;
     }
   
    function getTopValue() public view returns(uint){
          return tStruct.topValue;
    }

    function getTopMapping() public view returns(uint) {
          return tStruct.topMapping[1]; 
    }

    function getNestedValue() public view returns(uint){
         return tStruct.nStruct.nestedValue;
    }
    
    function getNestedMapping() public view returns(bool) {
        return tStruct.nStruct.nestedMapping[1];
    }
}



