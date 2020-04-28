pragma solidity ^0.5.13;
/**
 * 1.revert()函数————终止运行并撤销状态更改————验证
 * 2.revert(string reason)函数————终止运行并撤销状态更改,并提供一个解释性的字符串————验证
 *
 * @author Albedo
 * @dev 2019/12/30
 **/
contract RevertHandle {
    //revert()终止运行并撤销状态更改
    function revertCheck(uint param) public {
        if (param > 10) {revert();}
    }

    //revert(string reason) 终止运行并撤销状态更改,并提供一个解释性的字符串
    function revertReasonCheck(uint param) public {
        if (param > 10) {revert("check catch exception");}
    }
}