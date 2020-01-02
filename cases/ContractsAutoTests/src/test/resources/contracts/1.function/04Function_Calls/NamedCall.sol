pragma solidity 0.5.13;
/**
 * 验证函数的具名调用
 * @author liweic
 * @dev 2019/12/26 16:10
 */
contract NamedCall {
    //交换传入值的顺序并返回
    function exchange(uint key, uint value) public returns (uint, uint){ 
        return (value, key);
    }

    function calltest() public returns (uint, uint){
        //任意顺序的通过变量名来指定参数值
        return exchange({value: 2, key: 1}); 
    }
}
