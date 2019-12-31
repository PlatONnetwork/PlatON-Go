pragma solidity ^0.5.0;
/**
 * 因为填充(padding)的时候， bytesX是右填充(低比特位补0)，
 * 而 uintY是左填充(高比特位补0)，二者直接转换可能会导致不符合预期的结果，
 * 所以现在当 bytesX和uintY长度大小不一致(即X*8 != Y)时,不能直接转换,
 * 必须先转换到相同长度,再转换到相同类型
 *
 * @author hudenian
 * @dev 2019/12/24 09:57
 */


contract DisallowTypeChange {

    uint32 public y;

    function testChange() payable public{
        //由于uint32（4字节）小于bytes8（8字节）
        // bytes8 x1; // 8*8 length
        // uint32 y1 = uint32(x1); //error, 64 length can not change to 32 length

        //  bytes1 a = 255;
        //     bytes2 b = "aA";

        bytes8 x = hex"aaaa" ;
        bytes4 x4 = bytes4(x); //right,转将长度转成一致
        y = uint32(x4); //right,再进行类型转换
    }

    //查询y的值
    function getY() public view returns (uint32){
        return y;
    }

}