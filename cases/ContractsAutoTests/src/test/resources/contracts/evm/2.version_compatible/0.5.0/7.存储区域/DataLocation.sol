pragma solidity ^0.5.0;
/**
 * 06- 存储区域
 * 1-结构体(struct)，数组(array)，
 *   映射(mapping)类型的变量必须显式声明存储区域( storage， memeory， calldata)，
 *   包括函数参数和返回值变量都必须显式声明
 * 2-external 的函数的引用参数和映射参数需显式声明为 calldata
 *
 * @author hudenian
 * @dev 2019/12/24 09:57
 */

contract DataLocation {

    struct Person{
        string name;
        uint age;
    }

    mapping(uint => Person) persons;

    bytes res_data;


    /**
     * 1.结构体(struct)类型的变量入参及返回值必须显示声明存储区域
     */
    function set_struct_person(Person memory _person,uint _id) internal{
        persons[_id] = _person;
    }


    /**
     * 2. 数组(array) 类型的变量入参及返回值必须显示声明存储区域
     */
    // function transfer(address _to,uint _value,bytes  _data) external returns (bool success){ //数组数据类型存储区域必须声明为calldata
    function testBytes(bytes calldata _data) external returns (bytes memory){
        res_data = _data;
        return res_data;
    }

    /**
     * 3.mapping 类型的变量入参及返回值必须声明存储位置
     * 允许internal函数的参数及返回值为mapping指针
     */
    function set_person(mapping(uint=>uint) storage id_age_map,uint _id) internal returns(mapping(uint=>uint) storage)  {
        Person memory person;
        id_age_map[_id]= 23;
        return id_age_map;
    }


    /**
     * public 函数入参必须显示声明为memory
     */
    function getPerson(uint _id) view external returns (string memory,uint){
        return (persons[_id].name,persons[_id].age);
    }

    /**
     * public 函数入参必须显示声明为memory
     */
    function getBytes() view external returns (bytes memory){
        return res_data;
    }

    /**
     * public 函数入参必须显示声明为memory
     */
    function savePerson(uint _id,string memory _name,uint _age)  public returns (bool){
        Person memory _person = Person(_name,_age);
        set_struct_person(_person,_id);
        return true;
    }

}