pragma solidity ^0.5.13;
import "./BaseStaticLibrary.sol";
/**
 * 库引用类似引用static方法验证
 * 解释：如果L作为库的名称，f()是库L的函数，则可以通过L.f()的方式调用
 *
 * @author Albedo
 * @dev 2019/12/25
 **/


contract LibraryStaticUsing {
    event Result(bool result);

    //库引用类静态方式验证。
    function register(uint value) public returns (bool result) {
        result = BaseStaticLibrary.compare(123, value);
        emit Result(result);
    }
}