pragma solidity ^0.5.0;
/**
 * 10-链上函数
 * 1- 0.5.0版本函数 delegatecall() 代替 0.4.25版本函数 callcode()
 * 2- 0.5.0版本函数 selfdestruct() 代替 0.4.25版本函数 suicide()
 *
 * 4- 0.5.0版本函数 revert()， require()，assert() 代替 0.4.25版本函数 throw
 * 5- 0.5.0版本函数 call()族只接受一个参数 bytes，返回成功是否的bool及函数执行的返回值
 *    0.4.25版本函数 call()族函数接收多个参数方式，只返回成功是否的 bool
 * 跨合约调用见11.cross_contract_call
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 *
 */

contract ChainFunction {

    address owner;
    uint fortune;
    bool isDeceased;

    constructor() public payable{
        owner = msg.sender;
        fortune = msg.value;
        isDeceased = false;
    }

    modifier onlyOwner{
        require(msg.sender == owner);
        _;
    }

    modifier mustBeDeceased{
        require(isDeceased == true);
        _;
    }

    /**
     * 0.5.0 版本使用revert()， require()，assert()关键字
     * _isDeceased为false，或者less9小于9则抛出异常
     */
    function deceased(bool isDeceased,uint less9) view public returns(address){
        assert(isDeceased == true);
        if(less9 < 9){
            revert();
        }
        return msg.sender;
    }

    /**
     * _isDeceased为false则抛出异常
     */
    function deceasedWithModify(bool _isDeceased) view public returns(address){
        require(_isDeceased == true);
        return msg.sender;
    }
}
