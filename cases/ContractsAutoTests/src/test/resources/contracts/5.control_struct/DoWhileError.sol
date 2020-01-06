pragma solidity ^0.5.13;
/**
 *  dowhile控制结构语法
 *
 *
 * @author hudenian
 * @dev 2020/1/6 11:09
 */


contract DoWhileError {

    struct S { bool f; }
    S s;

    /**
    * 结构体变量返回的是指针类型变量，可以未分配的情况下返回，导致未定义的异常
    */
    //    function f() internal view returns (S storage c) {
    //        do {
    //            break;
    //            c = s;
    //        } while(false);
    //    }

    /**
     * 正确
     */
    function doWhileControl() internal view returns (S storage c) {
        do {} while((c = s).f);
    }

    function getDoWhileControlRes() public view returns(bool){
        S storage doWileS = doWhileControl();
        return doWileS.f;
    }

}