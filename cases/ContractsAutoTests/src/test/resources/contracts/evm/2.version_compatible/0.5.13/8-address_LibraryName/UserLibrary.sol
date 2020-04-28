pragma solidity 0.5.13;
/**
 * 1. 允许public和external库函数的参数和返回变量使用mapping类型
 *
 * @author hudenian
 * @dev 2019/12/20 11:09
 */

library UserLibrary {
    struct _User {
        string name;
        uint age;
        mapping(uint=>string) user;
    }

    function set_age(uint new_age) internal returns(_User memory)  {
        _User memory user;
        user.age =new_age;
        return user;
    }

    /**
     * library库中允许包含mapping做为入参与出参
     * 
     */
    function set_user(mapping(uint=>uint) storage id_age_map,uint _id) public returns(mapping(uint=>uint) storage)  {
        _User memory user;
        id_age_map[_id]= 23;
        return id_age_map;
    }

    /**
     * library库中允许包含mapping做为入参与出参
     * 
     */
    function set_user_inter(mapping(uint=>uint) storage id_age_map,uint _id) external returns(mapping(uint=>uint) storage)  {
        _User memory user;
        id_age_map[_id]= 24;
        return id_age_map;
    }

    function get_age() public {

    }

}