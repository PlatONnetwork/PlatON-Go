pragma solidity ^0.5.13;
/**
 * 验证函数四种可见性external,public,internal,private
 * @author liweic
 * @dev 2019/12/27 10:10
 */

contract Visibility {

    //private：私有函数和状态变量仅在当前合约中可以访问，在继承的合约内，不可访问
    function fpri(uint a) view private returns(uint) {
        return a + 1;
    }

    //external:外部函数是合约接口的一部分，所以我们可以从其它合约或通过交易来发起调用
    function fe(uint a) view external returns(uint){
        return a + 2;
    }

    //public:公开函数是合约接口的一部分，可以通过内部，或者消息来进行调用
    function fpub(uint a) view public returns(uint) {
        return a + 3;
    }

    //internal：这样声明的函数和状态变量只能通过内部访问,也可以在继承合约里调用
    function add(uint a, uint b) view internal returns(uint) {
        return a+b;
    }
}


contract VisibilityCall {
    function readData() public payable returns(uint localA, uint localB){
        Visibility visibility = new Visibility();
        //uint local = visibilitytest.fpri(10); // error: member "fpri" is not visible,编译报错
        localA = visibility.fe(1);
        localB = visibility.fpub(1);
        //uint localB = visibilitytest.add(1, 2); // error: member "add" is not visible,编译报错
    }

}


contract Inter is Visibility {
    function g() view public returns(uint){
        //继承合约里调用internal函数
        return add(1,2);  // acces to internal member (from derivated to parent contract)
    }
}