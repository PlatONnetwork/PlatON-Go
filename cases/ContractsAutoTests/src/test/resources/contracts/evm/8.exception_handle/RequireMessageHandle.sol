pragma solidity ^0.5.13;
/**
 *  require(bool condition, string message)函数（该函数可以自定义message信息）验证
 * 如果条件不满足则撤销状态更改，用于检查由输入或者外部组件引起的错误
 * 此时为require异常返回异常信息，非运行时异常检测
 * @author Albedo
 * @dev 2019/12/19
 **/
contract RequireMessageHandle {
    //如果调用require的参数为false
    function paramException(uint param) public {
        require(param<10,"整型大小比较异常");
    }
}