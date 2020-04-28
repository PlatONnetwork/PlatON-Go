pragma solidity ^0.5.13;
/**
 * 1. 允许public和external库函数的参数和返回变量使用mapping类型
 *
 * @author hudenian
 * @dev 2019/12/20 11:09
 */

import "./UserLib.sol";

contract UserMapping{

    using UserLib for *;


    mapping(uint=>uint) id_age_map;

    function setOutUser(uint _age,uint _id) public {

        // uint age = _age;
        // 正确的用法：等价于set_age(age);
        // age.set_age();

        id_age_map[_id] =_age;
        // mapping(uint=>uint) storage myIdAgeMap = set_user(id_age_map,_id);

        id_age_map.set_user(_id);


        // return id_age_map[1];
        // 错误用法：get_age()无参数或和第一个参数不匹配,就不能这么调用
        // age.get_age();
    }

    function getOutUser(uint _id) view public returns(uint age){
        return id_age_map[_id];
    }
}