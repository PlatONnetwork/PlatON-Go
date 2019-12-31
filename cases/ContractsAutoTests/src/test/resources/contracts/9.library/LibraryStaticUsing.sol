pragma solidity ^0.5.13;

/**
 * 库引用类似引用static方法验证
 * 解释：如果L作为库的名称，f()是库L的函数，则可以通过L.f()的方式调用
 *
 * @author Albedo
 * @dev 2019/12/25
 **/

library BaseStaticLibrary {
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

contract LibraryStaticUsing {
    BaseStaticLibrary.Data knownValues;

    //库引用类静态方式验证。
    function register(uint value) public {
        // 不需要库的特定实例就可以调用库函数，
        // 因为当前合约就是“instance”。
        require(BaseStaticLibrary.insert(knownValues, value));
    }
}