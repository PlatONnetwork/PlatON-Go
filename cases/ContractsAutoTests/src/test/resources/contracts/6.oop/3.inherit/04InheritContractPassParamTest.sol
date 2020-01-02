pragma solidity 0.5.13;

/**
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约继承功能点
 *
 *继承(is)简述：合约支持多重继承，即当一个合约从多个合约继承时，
 *在区块链上只有一个合约被创建，所有基类合约的代码被复制到创建合约中。
 *-----------------  测试点   ------------------------------
 *1、多重继承情况
 *2、多重继承(合约存在父子关系)
 *3、继承支持传参
 *4、合约函数重载(Overload)
 */

/**
 *验证：3、继承支持传参(继承中基类构造函数的传参)
 *-------------------------------------------
 * 验证结果：继承支持传参两种方式
 */
 
 contract InheritContractBase {
      
      uint a = 0;
      constructor(uint x) public {
          a = x;
      }
 }

//给基类构造函数传参方式一
contract InheritContractSub1 is InheritContractBase(2) {
    function getData1() public view returns (uint) {
        return a;
    }
}

//给基类构造函数传参方式二
contract InheritContractSub2 is InheritContractBase {
    uint b;
    constructor () InheritContractBase(3) public {
        uint y = 1;
        b = y + a;
    }

    function getData2() public view returns(uint){
        return b;
    }
}

