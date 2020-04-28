pragma solidity ^0.4.21;

/**
 * 0.4.25版本支持类型转换（0.5.x版本已弃用）验证
 * （1）0.4.25版本bytesX和uintY支持直接转换
 * （2）0.4.25版本10进制数值可以直接转换成 bytesX类型
 * （3）0.4.25版本16进制数值如果长度与 bytesX不相等，能直接转换成 bytesX类型
 * @author Albedo
 * @dev 2019/12/23
 **/
contract TypeCast {
    function typeCast() public view returns (uint16, bytes4, bytes4){
        //0.4.25版本bytesX和uintY支持直接转换验证
        bytes2 a = 0x12;
        uint16 b = uint16(a); //16进制转换成10进制，b的值为18
        //0.4.25版本10进制数值可以直接转换成 bytesX类型验证
        bytes4 c = bytes4(1234); //转换成16进制 0x000004d2
        //0.4.25版本16进制数值如果长度与 bytesX不相等，能直接转换成 bytesX类型验证
        bytes4 d = bytes4(0x1234); //类型转换 0x00001234
        return (b, c, d);
    }
}