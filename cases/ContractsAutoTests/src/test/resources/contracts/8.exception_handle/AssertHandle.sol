pragma solidity ^0.5.13;
/**
 * assert(bool condition)函数验证：
 * 如果条件不满足，则使当前交易没有效果 ，gas正常消耗，用于检查内部错误
 * 1.数组越界访问产生生异常验证，如i >= x.length 或 i < 0时访问x[i]
 * 2.定长bytesN数组越界访问产生异常验证
 * 3.被除数为0或取模运算产生异常验证
 * 4.对一个二进制移动一个负的值产生异常验证
 * 5.整数可以显式转换为枚举时，如果将过大值，负值转为枚举类型则抛出异常被除数为0异常验证
 * 6.调用内部函数类型的零初始化变量验证
 * 7.用assert的参数为false产生异常验证
 * @author Albedo
 * @dev 2019/12/19
 **/

library ArrayUtils {
    // 编译异常：6.如果函数类型调用未初始化internal函数时，将会产生异常
    function map(uint[] memory self, function (uint) pure returns (uint) f)
    internal
    pure
    returns (uint[] memory r)
    {
        r = new uint[](self.length);
        for (uint i = 0; i < self.length; i++) {
            r[i] = f(self[i]);
        }
    }

    function reduce(
        uint[] memory self,
        function (uint, uint) pure returns (uint) f
    )
    internal
    pure
    returns (uint r)
    {
        r = self[0];
        for (uint i = 1; i < self.length; i++) {
            r = f(r, self[i]);
        }
    }

    function range(uint length) internal pure returns (uint[] memory r) {
        r = new uint[](length);
        for (uint i = 0; i < r.length; i++) {
            r[i] = i;
        }
    }
}

contract AssertHandle {

    using ArrayUtils for *;
    //1.如果越界，或负的序号值访问数组，如i >= x.length 或 i < 0时访问x[i]
    function outOfBoundsException() public {
        //编译异常：数组越界访问
        //uint8[3] memory balance = [1, 2, 3];
        //balance[4]=12;
        //balance[-1]=12;
    }

    //2.如果序号越界，或负的序号值时访问一个定长的bytesN
    function noOutOfBoundsException() public {
        //编译异常：定长bytesN越界访问
        // bytes4 b4=0x12345678;
        // b4[4];
        // b4[-1];
    }

    //3.被除数为0或取模运算
    function dividendZeroException() public {
        //编译异常：被除数为0和取模运算
        // uint result = 12/0;
        // uint result = 12%    0;
    }

    //4.移位负数位
    function binaryMoveMinusException() public {
        //编译异常：移位负数位
        // uint8 uu=2;
        // uu<<-2;
    }
    //7.如果调用assert的参数为false
    function paramException() public {
        assert(1 < 0);
    }

}