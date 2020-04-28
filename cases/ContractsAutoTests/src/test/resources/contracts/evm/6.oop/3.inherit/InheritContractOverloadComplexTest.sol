pragma solidity 0.5.13;

/**

 * @author qudong
 * @dev 2019/12/23
 * 
 *测试合约继承功能点
 *继承(is)简述：合约支持多重继承，即当一个合约从多个合约继承时，
 *在区块链上只有一个合约被创建，所有基类合约的代码被复制到创建合约中。
 *-----------------  测试点   ------------------------------
 *1、多重继承情况
 *2、多重继承(合约存在父子关系)
 *3、继承支持传参
 *4、合约函数重载(Overload)
 */


/**
 *4、合约函数重载(Overload)复杂情况
 **/

contract InheritContractOverloadBaseBase {
	uint public x;
	uint public y;
	function init(uint a, uint b) public {
		x = b;
		y = a;
	}
	function init(uint a) public {
		x = a + 1;
	}
    function getX() public view returns (uint) {
        return x;
    } 
     function getY() public view returns (uint) {
        return y;
    } 
}

contract InheritContractOverloadBase is InheritContractOverloadBaseBase {
	function init(uint a, uint b) public {
		x = a;
		y = b;
	}
	function init(uint a) public {
		x = a;
	}
}

contract InheritContractOverloadChild is InheritContractOverloadBase {
	function initBase(uint c) public {
		InheritContractOverloadBase.init(c);
	}
	function initBase(uint c, uint d) public {
		InheritContractOverloadBase.init(c, d);
	}
	function initBaseBase(uint c) public {
		InheritContractOverloadBaseBase.init(c);
	}
	function initBaseBase(uint c, uint d) public {
		InheritContractOverloadBaseBase.init(c, d);
	}
}




