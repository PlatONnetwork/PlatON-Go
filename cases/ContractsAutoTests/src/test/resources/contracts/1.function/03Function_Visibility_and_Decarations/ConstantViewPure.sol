pragma solidity ^0.4.25;
/**
 * 验证三个函数修饰词constant, view, pure
 * @author liweic
 * @dev 2019/12/27 14:10
 */

contract ConstantViewPure{
    string name;
    uint public age;
    
    function constantViewPure() public{
        name = "fanxian";
        age = 19;
    }
    
    function getAgeByConstant() public constant returns(uint){
        age += 1;  //声明为constant，在函数体中又试图去改变状态变量的值，编译会报warning, 但是可以通过
        return age;  // return 20, 状态变量age的值不会改变，仍然为19！
    } 
    
    function getAgeByView() public view returns(uint){
        age += 1; //view和constant效果一致，编译会报warning，但是可以通过
        return age; // return 20，状态变量age的值不会改变，仍然为19！
    }
    
    function getAgeByPure() public pure returns(uint){
        //return age; //编译报错！pure比constant和view都要严格，pure完全禁止读写状态变量！
        return 1;
    }
}