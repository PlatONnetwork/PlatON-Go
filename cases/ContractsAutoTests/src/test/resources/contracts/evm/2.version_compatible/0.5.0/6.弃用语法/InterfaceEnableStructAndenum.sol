pragma solidity ^0.5.0;
/**
 * 09-弃用语法
 * 10-接口允许定义结构体与枚举
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

interface InterfaceEnableStructAndenum {
    
     //结构体
     struct Voter {
         string name;
         uint age;
     }
     
     //枚举
     enum gender{
         MALE,
         FEMALE
     }
      
}