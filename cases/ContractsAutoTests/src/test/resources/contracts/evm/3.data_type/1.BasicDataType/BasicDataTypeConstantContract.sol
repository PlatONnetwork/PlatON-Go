pragma solidity 0.5.3;

/**
 * @author qudong
 * @dev 2019/12/23
 *
 *测试合约基本数据类型字面常量
 *
 *1、地址类型(Address)
 *2、地址字面常量
 *3、有理数常量
 *4、整数字面常量
 *5、字符串字面常量
 *6、十六进制字面常量
 *7、枚举类型
 */

contract BasicDataTypeConstantContract {

    /**
     *1、地址类型(Address)
     *包含两种形式：
     *address：保存一个20字节的值（以太坊地址的大小），160位。
     *address payable ：可支付地址，与 address 相同，不过有成员函数 transfer 和 send 。
     *------------- 测试点 ----------
     *1、地址类型赋值使用
     *2、地址成员变量：
     *address.balance()：查询一个地址余额
     *address.transfer()：发送货币，[常情况转账、异常情况转账]
     *address.send() :transfer 的低级版本。如果执行失败，当前的合约不会因为异常而终止，但 send 会返回 false。
     */

    //查询地址余额
    function getBalance(address addr) public view returns (uint) {
        return addr.balance;
    }

     //当前合约的余额  
    function getCurrentBalance() public view returns (uint) {
        return address(this).balance;
    }

    //转账transfer
    function goTransfer(address payable addr) public payable {
        addr.transfer(msg.value);
    }

    //转账send，向当前合约发送货币
    function goSend(address payable addr) public payable returns(uint amount, bool success) {
        //msg.sender 全局变量，调用合约的发起方
        //msg.value 全局变量，调用合约的发起方转发的货币量，以wei为单位。
        //send() 执行的结果
        return (msg.value, addr.send(msg.value));
    }

    /**
     *2、地址字面常量
     *   通过了地址校验和测试的十六进制字面常量,长度在 39 到 41 个数字的会作为地址字面常量
     */
    function getAddress() public returns (address v) {

        address b = address("lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2");
        return b;
    }

    /**
     *3、有理数常量
     *   整数、分数(有限小数、无限循环小数)
     *
     */
    function getValue() public view returns (uint128,uint,uint) {

        uint128 o = 2.5 + 0.5;
        uint v = 2e10;
        uint l = 3e15;
        return (o,v,l);
    }

    /**
    *4、整数字面常量
    *   整数字面常量由范围在0-9的一串数字组成，表现成十进制。
    *
    */

    function getInt() public returns (int  v) {
        int  a = 6;
        return a;
    }

    /**
    *5、字符串字面常量:
    *   字符串字面常量是指由双引号或单引号引起来的字符串
    *-------  测试点----
    * 1)、赋值操作
    * 2)、字符串是特殊的动态长度字节数组
    *      转换：字符串字面常数的类型———>转换成 bytes1，……，bytes32 ,或者bytes
    * 3)、字符串不能够字节的修改长度和内容，需要转换为bytes动态字节数组
    */

   //1)、赋值
     string strA = 'hello';
     string strB = "world";

    function getStrA() public view returns (string memory) {
        return strA;
    }

    //2)、转换(字符串是特殊的动态字节数组),字符串不能直接的获取长度和内容，可以通过转换字节数组进行获取
    function getStrALength() public view returns (uint) {
        return bytes(strA).length;
    }

    function setStrA() public view  returns (string memory) {
        bytes memory b = bytes(strA);
        b[0] = 'a';
        return string(b);
    }

    /**
     *6、十六进制字面常量:
     *   十六进制字面量，以关键字hex打头，后面紧跟用单或双引号包裹的字符串。
     *   由于一个字节是8位，所以一个hex是由两个[0-9a-z]字符组成的
     *---------  测试点 ----------
     * 十六进制字面常量定义赋值取值操作
     *
     */

    function getHexLiteraA() public view  returns(bytes1){

         bytes1 b = hex"c8";//十进制数字200 <====> 十六进制c8 <===> 二进制11001000
         return  b;
    }

    function getHexLiteraB() public view  returns(bytes2){

         bytes2 b = hex"01f4"; //十进制数字256 <===> 十六进制1f4 <===> 二进制111110100
         return  b;
    }

    function getHexLiteraC() public view returns (bytes2, bytes1, bytes1){

        bytes2 b = hex"01f4";//十进制数字256 <===> 十六进制1f4 <===> 二进制111110100
        return (b, b[0], b[1]);
    }


    /**
     *7、枚举类型
     *   枚举是在Solidity中创建用户定义类型的一种方法，其可以从整型显示转换成枚举。
     *   选项从0开始无符号整数值表示。
     *------------   测试点  --------------------
     *1、定义
     *2、取值
     */

    enum Season{Spring, Summer, Autumn, Winter}

    function getSeasonA() public view returns(Season){
        return printSeason(Season.Summer);
    }

    function getSeasonB() public view returns(Season){
        //Season s = Season(5);//越界
        Season s = Season(3);
        return s;
    }

    function getSeasonIndex() public view returns(uint){
        uint s = uint(Season.Spring);
        return s;
    }

    function printSeason(Season s) public view returns(Season) {
        return s;
    }
}