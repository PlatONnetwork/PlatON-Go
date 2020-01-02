pragma solidity ^0.4.4;

/**
 * 1.验证0.4.4版本下的特殊版本下的构造函数
 * 2.验证selfdestruct析構函數
 * 3.验证==表达式
 * 4.验证uint数据类型
 * @author qcxiao
 * @dev 2019/12/18 15:09
 */
contract Person {
    uint _age;
    uint _height;
    address _owner;
    string _name;

    function Person() {
        _height = 180;
        _age = 20;
        _owner = msg.sender;
        _name = "qcxiao";
    }

    function f() {
        modify(_name);
    }

    // 引用类型string需要配合storage关键字传递参数的地址,并需要配合internal或private声明为内部函数
    function modify(string storage name) internal {
        bytes(name)[0] = "Q";
    }

    function name() constant returns (string) {
        return _name;
    }

    function owner() constant returns (address) {
        return _owner;
    }

    function setAge(uint age) {
        _age = age;
    }

    function age() constant returns (uint) {
        return _age;
    }

    function setHeight(uint height) {
        _height = height;
    }

    // constant代表只讀,不需要消耗gas
    function height() constant returns (uint) {
        return _height;
    }

    function kill() constant {
        if (_owner == msg.sender) {
            // 析構函數
            selfdestruct(_owner);
        }
    }


}