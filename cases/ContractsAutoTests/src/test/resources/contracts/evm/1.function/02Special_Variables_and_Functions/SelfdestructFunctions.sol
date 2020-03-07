pragma solidity 0.5.7;
/**
 * 验证构造函数constructor
 * 验证销毁函数selfdestruct,0.5.x之后弃用suicide
 * @author liweic
 * @dev 2019/12/28 16:00
 */

contract SelfdestructFunctions {
    //声明状态变量
	uint count = 0;    
	address payable owner;
    
    //构造函数
    constructor () public {       
    	owner = msg.sender;    
    } 
    
    function increment() public {
    	uint num = 5;
        if(owner == msg.sender) {
	    	count = count + num;       
	    }    
	}
	
	function getCount() view public returns (uint) {
    	return count; 
    }
    
 	function selfKill() public {
 	if (owner == msg.sender) {
       	selfdestruct(owner);
       }
 	}
}