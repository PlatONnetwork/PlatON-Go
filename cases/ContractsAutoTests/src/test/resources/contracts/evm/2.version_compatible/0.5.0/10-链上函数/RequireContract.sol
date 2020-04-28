pragma solidity ^0.5.13;
/**
 * 10-链上函数
 * 函数 require()
 *
 * @author hudenian
 * @dev 2020/1/8 09:57
 *
 */

contract RequireContract {

    uint result;
    /**
     * require退回剩下的gas
     * 验证输入参数合法性
     */
    function toSenderAmount(uint frist,uint second) public {
        require(frist > second);
        result =  frist - second;
    }


    function getResult() view public returns(uint){
        return result;
    }


}