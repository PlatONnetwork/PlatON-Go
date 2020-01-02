pragma solidity 0.5.13;

/**
 * *
 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约继承功能点
 *
 *继承(is)简述：合约支持多重继承，即当一个合约从多个合约继承时，
 *在区块链上只有一个合约被创建，所有基类合约的代码被复制到创建合约中。
 *
 *
 *-----------------  测试点   ------------------------------
 *1、多重继承情况
 *2、多重继承(合约存在父子关系)
 *3、继承支持传参
 *4、合约函数重载(Overload)
 */
 





/**
 *4、合约函数重载(Overload)：合约可以有多个同名函数，可以有不同输入参数。
 *   重载解析和参数匹配
 *
 *
 **/

contract InheritContractOverload {

     uint sumInt;

    function sum(uint a,uint b) public pure returns(uint sumInt) {
        sumInt = a + b;
    }

    function sum(uint a,uint b, uint c) public pure returns(uint sumInt) {
        sumInt = a + b + c;
    }

    function getData1() public view returns(uint) {
          
          return sum(1,2);
    }

     function getData2() public view returns(uint) {
          
          return sum(1,2,3);
    }


}


