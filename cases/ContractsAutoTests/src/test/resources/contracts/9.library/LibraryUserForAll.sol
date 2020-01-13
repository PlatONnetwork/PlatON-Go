pragma solidity ^0.5.13;
/**
 * 引用using for方式验证
 * 解释：using A for * 的效果是，库 A 中的函数被附加在任意的类型上。
 * @author Albedo
 * @dev 2019/12/25
 **/
library SearchLibrary {
    function indexOf(uint[] storage self, uint value)
    internal
    returns (uint)
    {
        for (uint i = 0; i < self.length; i++)
            if (self[i] == value) return i;
        return uint(- 1);
    }
}

contract LibraryUsingForAll {

    //using A for *
    using SearchLibrary for *;
    uint[] data;

    function replace(uint _old, uint _new) public {
        // 执行库函数调用
        uint index = data.indexOf(_old);
        if (index == uint(- 1))
            data.push(_new);
        else
            data[index] = _new;
    }
}