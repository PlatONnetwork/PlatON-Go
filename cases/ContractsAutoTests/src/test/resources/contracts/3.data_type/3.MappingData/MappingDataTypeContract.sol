pragma solidity 0.5.13;


/**
 *
 * 
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试映射功能点
 *
 *映射(Mapping)简述：映射类型，一种键值对的映射关系存储结构
 *
 *
 *-----------------  测试点   ------------------------------
 *1、定义映射
 * mapping(keyType => keyValue) public maps;
 *  键类型--允许除映射、变长数组、合约、枚举、结构体外的几乎所有类型
 *  值类型--值类型没有任何限制，可以为任何类型包含映射类型
 *2、赋值、取值
 *   map[keyType] ← 不同数值类型;
 *   
 *
 */


contract MappingContractTest {


  
  /**
   *验证：1、定义映射
   * 键类型--允许除映射、变长数组、合约、枚举、结构体外的几乎所有类型
   * 值类型--值类型没有任何限制，可以为任何类型包含映射类型
   * ---------------------------------------------------------
   *验证结果：1)、mapping定义键类型，包含映射、变长数组、合约、枚举、结构体，会编译异常
   *         2)、mapping定义值类型，没有任何限制
   */

   enum SizeEnum {XL, XXL, XXXL}
   struct PeopleStruct {
      uint id;
      string name;
      bool sex;
   }
   mapping (int => SizeEnum) map;
   uint[] uintArr;


   mapping (uint => address) public addressMap; //正确
   mapping (int => bool) public boolMap;//正确
   mapping (bool => byte) public  byteMap;//正确
   mapping (bytes1 => string) public stringMap;//正确
   mapping (bytes => uint) public uintMap;//正确
   mapping (string => int) public intMap;//正确
   mapping (address => bytes) public bytesMap;//正确
   mapping (int => mapping (int => SizeEnum)) public sizeMap;//正确
   mapping (int => PeopleStruct) public peopleMap;//正确
   mapping (int => SizeEnum) public SizeEnumMap;//正确

   // mapping (SizeEnum => int) public enumMap;//异常
   // mapping (PeopleStruct => bytes10) public structMap;//异常
   // mapping (map => string) public stringMap1;//异常
   // mapping (uintArr => string) public stringMap1;//异常


/**
 *2、赋值、取值
 *   map[keyType] ← 不同数值类型;
 */
 
    mapping (uint => string) nameMap;
    string[] nameArr = ["Lucy","Ella","Lily"];


   function addName() public {
         
         for (uint i = 0; i < nameArr.length; i++) {
            nameMap[i] =  nameArr[i];
         }
   }

   function getName(uint index) returns (string) {
      
      return nameMap[index];
   }
 

  





    
}
