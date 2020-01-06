pragma solidity ^0.5.13;
/**
 *  控制结构
 *  1. if...else
 *  2. do...while
 *  3. for循环
 *  4. for循环包含break
 *  5. for循环包含continue
 *  6. for循环包含return
 *  7. 三目运算符
 *
 *
 * @author hudenian
 * @dev 2020/1/6 11:09
 */


contract ControlError {

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

}