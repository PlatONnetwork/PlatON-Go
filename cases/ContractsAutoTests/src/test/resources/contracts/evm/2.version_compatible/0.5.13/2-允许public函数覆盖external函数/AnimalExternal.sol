pragma solidity ^0.5.13;
/**
 * 验证0.5.13版本允许public函数覆盖external函数
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */

contract AnimalExternal {

    string _birthDay; // 生日
    int public _age; // 年龄
    int internal _weight; // 身高
    string private _name; // 姓名

    constructor() public{
      _age = 29;
      _weight = 170;
      _name = "Lucky dog";
      _birthDay = "2011-01-01";
    }
    //声明external函数
    function birthDay() view external returns (string memory) {
      return _birthDay;
    }

    function age() view public returns (int) {
      return _age;
    }

    function height() view internal returns (int) {
      return _weight;
    }

    function name() view private returns (string memory) {
      return _name;
    }
}
