pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 * 测试数据类型操作符delete 
 * solidity中通过操作符delete用于释放空间，对于某个值类型赋予初始化值
 */

contract BasicDataTypeDeleteContract {

     bool boolValue;
     uint uintValue;
     address addr;
     bytes3 bytesValue;
     string str;
     int intValue;
     uint8[] arrayInt;

     enum LightEnum{RED, GREEN, YELLOW}
     LightEnum lightEnum;

     struct Person{
          uint8 id;
          string name;
     }
     Person person;
     mapping (uint8 => uint8) uintMap;
     
  /**
   * 数据赋值
   */
   function initBasicData() public {
         boolValue = true;
         uintValue = 2;
         addr = msg.sender;
         bytesValue = "123";
         str = "hello";
         intValue = 5;
         //结构体赋值
         person = Person(1,"Ella");
         //数组赋值
         arrayInt = [1,2,3];
         //mapping赋值
         uintMap[1] = 1;
         uintMap[2] = 2;
         uintMap[3] = 3;
         //枚举赋值
         lightEnum = LightEnum.YELLOW;
   }

   //  -----  delete基本数据类型 -----
   function deleteBasicData() public {
         delete boolValue;
         delete uintValue;
         delete addr;
         delete bytesValue;
         delete str;
         delete intValue; 
   }

   function getBasicData() public view returns (bool,uint,address,bytes3,string memory,int) {
         return (boolValue,uintValue,addr,bytesValue,str,intValue);
   }
   

   //   ------  delete结构体  ------
   function deleteStruct() public {
         
         delete person;
   }

   function getStruct() public view returns (uint8,string memory) {
        return (person.id,person.name);
   }

   // -----  delete数组 --------
   function deleteArray() public {
         delete arrayInt;         
   }

   function getArrayLength() public view returns (uint256) {
        return arrayInt.length;     
   }

   // -----  delete枚举  ------
   function deleteEnum() public {
          
          delete lightEnum;          
   }

   function getEnum() public view returns (LightEnum) {
         return lightEnum;
   }

   /**
    *   -------  delete映射 ----------
    * 映射的特殊性，由于映射的键并不总能有效遍历，所以不能直接删除一个映射
    * 但是，可以按照某一项删除映射 
    */
   function deleteMapping() public {
       
        //直接删除映射,编译异常
        //delete uintMap;
        //删除映射中某一项
        delete uintMap[2];
   }

   function getMapping() public view returns (uint8) {
        return uintMap[2];    
   }
}