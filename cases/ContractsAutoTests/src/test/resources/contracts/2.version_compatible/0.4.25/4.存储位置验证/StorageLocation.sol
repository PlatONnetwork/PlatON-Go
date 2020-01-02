pragma solidity ^0.4.25;

/**
 * 存储位置验证；
 * （1）0.4.25版本结构体(struct)、数组(array)、映射(mapping)等类型的变量不必显式声明存储位置验证
 * （2）0.4.25版本函数参数变量为数组(array)类型不须显式声明验证
 * （3）0.4.25版本external 的函数的数组(array)类型参数不需显式声明为 calldata验证
 * @author Albedo
 * @dev 2019/12/24
 **/
contract StorageLocation {
    bytes data;

    //0.4.25版本结构体(struct)类型的变量不必显式声明存储位置验证
    struct Person{
        string name;
        int8 age;
        int16 high;
    }
    //0.4.25版本映射(mapping)类型的变量不必显式声明存储位置验证
    mapping(address=>bytes) addValue;


    //0.4.25版本函数参数变量为数组(array)类型不须显式声明验证
    function storageLocaltionCheck(bytes _data) public view returns (bytes){
        addValue[msg.sender] =_data;
        return addValue[msg.sender];
    }


    //0.4.25版本external 的函数的数组(array)类型参数不需显式声明为 calldata验证
    function transfer(bytes _data) external view returns (bytes){
        data=_data;
        return data;
    }
}