pragma solidity 0.5.6;
/**
 * 1.验证函数修饰器的关键词modifier
 * 2.验证多个修饰器的用法
 * 3.验证修饰器接收参数和无参修饰器的用法
 * @author liweic
 * @dev 2019/12/25 15:10
 */
contract Modifiers {
    uint a = 10;
    
    modifier mf1 (uint b) {
        uint c = b;
        _;
        c = a;
        a = 11;
    }
    
    modifier mf2 () {
        uint c = a;
        _;
    }
    
    modifier mf3() {
        a = 12;
        return ;
        _;
        a = 13;
    }
    
    function test1() mf1(a) mf2 mf3 public {
        a = 1;
    }
    
    function test2 () public view returns (uint) {
        return a;
    }
}