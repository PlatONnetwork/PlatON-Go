pragma solidity ^0.5.13;
/**
 * 递归调用验证
 *
 *
 * @author Albedo
 * @dev 2019/12/19
 **/
contract RecursionCall {
    //实际调用次数
    uint total;
    /**
     * 递归验证：分别测试递归调用执行效率
     * 0~100
     * 100~1000
     * 1000~10000
     * 10000~100000
     * 100000+
     **/
    function recursionCallTest(uint n) public {
        if (total < n) {
            //业务逻辑（转账）

            ++total;
            recursionCallTest(n);
        }
    }
}