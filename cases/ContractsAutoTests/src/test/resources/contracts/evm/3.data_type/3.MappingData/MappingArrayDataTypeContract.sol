pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 *
 *测试映射功能点
 *映射(Mapping)简述：映射类型，一种键值对的映射关系存储结构
 *-----------------  测试点   ------------------------------
 * 验证mapping数组类型
 */

contract MappingArrayDataTypeContract {

   //定义mapping定长数组
   mapping (uint8 => uint8)[2] a;
   mapping (uint8 => uint8)[2] b;
   //定义赋值内部函数
   function setMappingInternal(mapping (uint8 => uint8)[2] storage mapArray,uint8 key,uint8 value1,
                                                     uint8 value2) internal returns (uint8,uint8) {
            //获取mapping数组中值                                             
            uint8 oldValue1 = mapArray[0][key];
            uint8 oldValue2 = mapArray[1][key];
            //给mapping数组赋新值
            mapArray[0][key] = value1;
            mapArray[1][key] = value2;
            return (oldValue1,oldValue2);
   }
   //调用内部函数进行赋值
   function set(uint8 key,uint8 value_a1,uint8 value_a2,uint8 value_b1,uint8 value_b2) public {
               
            setMappingInternal(a,key,value_a1,value_a2);
            setMappingInternal(b,key,value_b1,value_b2);
   }
   //查询mapping定长数组值
   function getValueByKey(uint8 key) public view returns (uint8,uint8,uint8,uint8) {
            return(a[0][key],a[1][key],b[0][key],b[1][key]);
   }
}
