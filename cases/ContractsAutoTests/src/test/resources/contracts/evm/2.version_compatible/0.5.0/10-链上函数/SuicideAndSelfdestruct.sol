pragma solidity ^0.5.0;
/**
 * 10-链上函数
 * 1- suicide()已弃用, 请使用 selfdestruct()
 *  部署合约账号必须要有balance
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 *
 */

contract SuicideAndSelfdestruct {
    uint count = 0;
    address payable owner;

    constructor() public{
        owner = msg.sender;
    }

    function increment() public {
        uint step = 10;
        if (owner == msg.sender) {
            count = count + step;
        }
    }

    function getCount() view public returns (uint){
        return count;
    }

    /**
     *0.5.0版本使用新的自杀函数
     */
    function kill() public{
        if (owner == msg.sender) {
            selfdestruct(owner); //0.5.0使用的函数
            //   suicide(owner); //0.4.25使用的函数（0.5.0废弃的函数）
        }
    }
}
