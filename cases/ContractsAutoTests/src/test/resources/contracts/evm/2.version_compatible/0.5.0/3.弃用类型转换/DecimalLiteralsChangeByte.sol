pragma solidity ^0.5.0;
/**
 * 10进制数值不能直接转换成 bytesX类型
 * 必须先转换到与 bytesX相同长度的 uintY，再转换到 bytesX类型
 *
 * @author hudenian
 * @dev 2019/12/24 09:57
 */


contract DecimalLiteralsChangeByte {

    bytes4 public b4;

    function testChange(uint a) public returns(bytes4) {
        //   bytes8 bt = bytes8(a); //编译出错,10进制数值不能直接转换成 bytesX类型，要像如下进行转换
        uint32 u1 = uint32(a); //先转成uint

        //       bytes5 b2 = bytes5(u1);//编译出错, uint32长度与byte4一致，不能转成byte5类型的长度

        b4 = bytes4(u1);//再转成byte

        return b4;
    }

    //5.查询b4的值
    function getB4() public view returns (bytes4){
        return b4;
    }



}