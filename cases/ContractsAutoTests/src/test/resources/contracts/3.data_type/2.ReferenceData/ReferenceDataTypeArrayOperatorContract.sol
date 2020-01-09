pragma solidity 0.5.13;

contract ReferenceDataTypeArrayOperatorContract {

    /**
     * 验证数组支持的运算符
     * 1)、比较运算符
     * 2)、位运算符
     */

     //1)、比较运算符:支持的比较运算符有<=，<，==，!=，>=，>，返回的结果是一个bool
    function arrayCompare() public view returns (bool,bool,bool,bool,bool) {
        
        bytes1 a = "a";
        bytes1 b = "b";
        bytes1 c = bytes1(uint8(97));

        bool r1 = a < b;//返回  ture (ASCII码 Dec a= 97，b = 98)
        bool r2 = a > b;//返回  false
        bool r3 = a == c;//返回 ture
        bool r4 = b != c;//返回 ture 
        bool r5 = a >= c;//返回 ture
        return (r1,r2,r3,r4,r5);
    }

    //2)、定长字节数组支持位运算符 &(按位与)，|(按位或)，^（按位异或），~（按位取反），以及<<（左移位），<<（右移位）


    //&(按位与)--两个操作数中位都为1，结果才为1，否则结果为0
   function arrayBitAndOperators() public view returns(bytes1,uint8) {

        bytes1 a = bytes1(uint8(129));//bin:1000 0001
        bytes1 b = bytes1(uint8(128));//bin:1000 0000
        bytes1 c = a&b;//hex:0x80,bin:1000 0000,dec:128
        uint8 d = uint8(c);
        return (c,d);
   }


  //|(按位或) -- 两个位只要有一个为1，那么结果就是1，否则就为0
   function arrayBitOrOperators() public view returns(bytes1,uint8) {

        bytes1 a = bytes1(uint8(129));//bin:1000 0001
        bytes1 b = bytes1(uint8(128));//bin:1000 0000
        bytes1 c = a|b;//hex:0x81,bin:1000 0001,dec:129
        uint8 d = uint8(c);
        return (c,d);
   }

   //~（按位取反） -- 如果位为0，结果是1，如果位为1，结果是0
   function arrayBitInverseOperators() public view returns(bytes1,uint8) {

        bytes1 a = bytes1(uint8(129));//bin:1000 0001
        bytes1 c = ~a;//hex:0x7e,bin:0111 1110,dec:126
        uint8 d = uint8(c);
        return (c,d);
   }

  //^（按位异或） -- 两个操作数的位中，相同则结果为0，不同则结果为1
   function arrayBitXOROperators() public view returns(bytes1,uint8) {

        bytes1 a = bytes1(uint8(129));//bin:1000 0001
        bytes1 b = bytes1(uint8(128));//bin:1000 0000
        bytes1 c = a^b;//hex:0x01,bin:0000 0001,dec:1
        uint8 d = uint8(c);
        return (c,d);
   }

   //<<（左移位） -- 按二进制形式把所有的数字向左移动对应的位数，高位移出（舍弃），低位的空位补零
   function arrayBitLeftShiftperators() public view returns(bytes1,uint8) {

        bytes1 a = bytes1(uint8(129));//bin:1000 0001
        bytes1 c = a<<1;//hex:0x02,bin:0000 0010,dec:2
        uint8 d = uint8(c);
        return (c,d);
   }

   //<<（右移位） -- 按二进制形式把所有的数字向右移动对应的位数，低位移出（舍弃），高位的空位补零
    function arrayBitRightShiftperators() public view returns(bytes1,uint8) {

        bytes1 a = bytes1(uint8(129));//bin:1000 0001
        bytes1 c = a>>1;//hex:0x40,bin:0100 0000,dec:64
        uint8 d = uint8(c);
        return (c,d);
   }

}
