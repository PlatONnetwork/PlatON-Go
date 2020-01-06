pragma solidity ^0.5.13;
/**
 *  for控制结构语法
 *
 *
 * @author hudenian
 * @dev 2020/1/6 18:02
 */


contract ForError {

    struct S { bool f; }
    S s;
    S s1;

    /**
    * 结构体变量返回的是指针类型变量，可以未分配的情况下返回，导致未定义的异常
    */
    /*function f() internal view returns (S storage c) {
        for(;; c = s) {
        }
    }*/


    function forControlFirst() internal view returns (S storage c) {
        for(c = s;;) {
        }
    }

    function getForControlRes() public view returns(bool){
        S storage forS = forControlFirst();
        return forS.f;
    }

    /**
    * 结构体变量返回的是指针类型变量，可以未分配的情况下返回，导致未定义的异常
    */
    /*function g() internal view returns (S storage c) {
        for(;;) {
            c = s;
        }
    }*/

    function forControlSecond() internal view returns (S storage c) {
        for(; (c = s1).f;) {
        }
    }

    function getForControlRes1() public view returns(bool){
        S storage forS = forControlSecond();
        return forS.f;
    }

}