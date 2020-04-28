pragma solidity ^0.5.13;
/**
 *  if控制结构语法
 *
 * @author hudenian
 * @dev 2020/1/6 18:02
 */


contract ifError {

    struct S { bool f; }
    S s;

    /**
    * 结构体变量返回的是指针类型变量，可以未分配的情况下返回，导致未定义的异常
    */
    /*function f(bool flag) internal view returns (S storage c) {
        if (flag) c = s;
    }*/


    function ifControl(bool flag) internal view returns (S storage c) {
        if (flag) c = s;
        else c = s;
    }

    function getIfControlRes() public view returns(bool){
        S storage ifS = ifControl(true);
        return ifS.f;
    }

    /**
    * 结构体变量返回的是指针类型变量，可以未分配的情况下返回，导致未定义的异常
    */
    /*function g(bool flag) internal returns (S storage c) {
        if (flag) c = s;
        else
        {
            if (!flag) c = s;
            else s.f = true;
        }
    }*/

    function ifControlSecond(bool flag) internal view returns (S storage c) {
        if (flag) c = s;
        else
        {
            if (!flag) c = s;
            else c = s;
        }
    }

    function getIfControlRes1() public view returns(bool){
        S storage ifS = ifControlSecond(true);
        return ifS.f;
    }

}