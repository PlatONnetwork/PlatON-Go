pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 *
 * 测试不同类型转换
 * -----------------  测试点   ------------------------------
 * 字节转换整型
 */

contract TypeConversionBytesToUintContract {

      //字节转换大位整型
    function bytesToBigUint() public view returns(uint64) {
        bytes4 a = "abcd";//hex：0x61626364
        uint32 b = uint32(a);//dec：1633837924
        uint64 c = uint64(b);//dec：1633837924
        return c;
    }

    //字节转换相同位数整数
    function bytesToSameUint() public view returns(uint8) {
       bytes1 a = "a";//hex：0x61
       uint8 b = uint8(a);//dec：97
       return b;
    }

    //字节转换小位整型
    function bytesToSmallUint() public view returns (uint16) {
        bytes4 a = "abcd";//hex：0x61626364，bin：0‭110 0001 0110 0010 0110 0011 0110 0100‬
        uint32 b = uint32(a);//dec：1633837924 
        uint16 c = uint16(b);//dec：25444 ，bin：0110 0011 0110 0100‬
        return c;
    }
}
