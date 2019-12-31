pragma solidity ^0.5.13;
/**
 * 引用using for 方式验证
 * 解释：
 * （1）指令using A for B 可用于附加库函数（从库 A）到任何类型（B）。 这些函数将接收到调用它们的对象作为它们的第一个参数。
 * （2）using A for * 的效果是，库 A 中的函数被附加在任意的类型上。
 * @author Albedo
 * @dev 2019/12/25
 **/
library BaseLibrary {
    // 我们定义了一个新的结构体数据类型，用于在调用合约中保存数据。
    struct Data {mapping(uint => bool) flags;}

    // 注意第一个参数是“storage reference”类型，因此在调用中参数传递的只是它的存储地址而不是内容。
    // 这是库函数的一个特性。如果该函数可以被视为对象的方法，则习惯称第一个参数为 `self` 。
    function insert(Data storage self, uint value) public returns (bool)
    {
        if (self.flags[value])
            return false;
        // 已经存在
        self.flags[value] = true;
        return true;
    }
}


library SearchLibrary {
    function indexOf(uint[] storage self, uint value)
    public
    view
    returns (uint)
    {
        for (uint i = 0; i < self.length; i++)
            if (self[i] == value) return i;
        return uint(- 1);
    }
}

contract LibraryUsingFor {
    //using A for B
    using BaseLibrary for BaseLibrary.Data;
    BaseLibrary.Data knownValues;
    //using A for *
    using SearchLibrary for *;
    uint[] data;

    function register(uint value) public {
        // 这里， BaseLib.Data 类型的所有变量都有与之相对应的成员函数。
        // 下面的函数调用和 `BaseLib.insert(knownValues, value)` 的效果完全相同。
        require(knownValues.insert(value));
    }

    function replace(uint _old, uint _new) public {
        // 执行库函数调用
        uint index = data.indexOf(_old);
        if (index == uint(- 1))
            data.push(_new);
        else
            data[index] = _new;
    }
}