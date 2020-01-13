pragma solidity ^0.5.13;
/**
 * 引用using for方式验证
 * 解释：指令using A for B 可用于附加库函数（从库 A）到任何类型（B）。 这些函数将接收到调用它们的对象作为它们的第一个参数。
 * @author Albedo
 * @dev 2019/12/25
 **/
library BaseLibrary {
    // 我们定义了一个新的结构体数据类型，用于在调用合约中保存数据。
    struct Data {mapping(uint => bool) flags;}

    // 注意第一个参数是“storage reference”类型，因此在调用中参数传递的只是它的存储地址而不是内容。
    // 这是库函数的一个特性。如果该函数可以被视为对象的方法，则习惯称第一个参数为 `self` 。
    function insert(Data storage self, uint value) internal returns (bool)
    {
        if (self.flags[value])
            return false;
        // 已经存在
        self.flags[value] = true;
        return true;
    }
}

contract LibraryUsingFor {
    event Result(bool result);


    //using A for B
    using BaseLibrary for BaseLibrary.Data;
    BaseLibrary.Data knownValues;

    function register(uint value) public returns (bool result){
        result=knownValues.insert(value);
        emit Result(result);
    }

}