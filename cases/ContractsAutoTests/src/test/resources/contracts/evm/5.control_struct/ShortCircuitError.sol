pragma solidity ^0.5.13;
/**
 *  && || 短路语法
 *
 * @author hudenian
 * @dev 2020/1/6 17:38
 */

contract ShortCircuitError {

    struct S { bool f; }
    S s;
    S s1;

    /**
    * 第一个条件为false，导致c未初始化，会导致返回空指针错误
    */
/*    function f() internal view returns (S storage c) {
        false && (c = s).f;
    }*/

    //right
    function f() internal view returns (S storage c) {
        (c = s).f && false;
    }


    function getF() public view returns(bool){
        S storage c = f();
        return c.f;
    }

    /**
    * 第一个条件为true，导致c未初始化，会导致返回空指针错误
    */
    /*function g() internal view returns (S storage c) {
        true || (c = s).f;
    }*/

    /*function h() internal view returns (S storage c) {
        // expect error, although this is always fine
        true && (false || (c = s).f);
    }*/

    //right
    function g() internal view returns (S storage c) {
        (c = s1).f || true;
    }

    function getG() public view returns(bool){
        S storage c = g();
        return c.f;
    }

}