pragma solidity ^0.5.0;
/**
 * 16进制数值如果长度与 bytesX不相等，也不能直接转换成 bytesX类型
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */


contract HexLiteralsChangeByte {

   bytes1 public b1;

   function testChange() payable public returns(bytes1) {
      bytes16  b16 = "ab";
      //   b16 = 0xff; //编译出错，16进制数值长度与 byte16长度不相等

      b1 = 0xf1; //16进制数值长度与 byte1长度不相等，正常转换(ff占用8位=1byte)
      return b1;
   }

   //查询y的值
   function getY() public view returns (bytes1){
      return b1;
   }
}