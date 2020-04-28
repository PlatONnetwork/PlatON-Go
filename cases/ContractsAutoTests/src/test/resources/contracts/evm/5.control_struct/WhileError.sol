pragma solidity ^0.5.13;
/**
 *  while控制结构语法
 *
 *
 * @author hudenian
 * @dev 2020/1/7 11:30
 */


contract WhileError {

    struct S { bool f; }
    S s;

    /**
    * 结构体变量返回的是指针类型变量，可以未分配的情况下返回，导致未定义的异常
    */
    /*function f() internal view returns (S storage c) {
        while(false) {
            c = s;
        }
    }*/

    function f() internal view returns (S storage c) {
        while((c = s).f) {
        }
    }

    function getWhileControlRes() public view returns(bool){
        S storage whileS = f();
        return whileS.f;
    }

}