pragma solidity ^0.5.0;
/**
 * constructor必须强制使用constructor声明
 * 且去除同名函数定义,0.4.25版本则通过同名函数定义，且constructor非必须
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */


contract ConstructorTest{

    uint public count = 0;

    //constructor必须强制使用constructor声明
    constructor(uint _count) public {
        count = _count;
    }

    function update(uint amount) public returns (address, uint){
        count += amount;
        return (msg.sender, count);
    }

    //5.查询count
    function getCount() public view returns (uint){
        return count;
    }

    //0.5.0弃用同名构造函数
    // function constructorTest() external  { 
    //     a = 1;     
    // }
}