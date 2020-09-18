package network.platon.test.evm.exceptionhandle;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AssertHandle;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;

import java.math.BigInteger;

/**
 * @title assert(bool condition)函数测试
 * 如果条件不满足，则使当前交易没有效果 ，gas正常消耗，用于检查内部错误
 * 1.数组越界访问产生生异常验证，如i >= x.length 或 i < 0时访问x[i]————(编译异常)
 * 2.定长bytesN数组越界访问产生异常验证————(编译异常)
 * 3.被除数为0或取模运算产生异常验证————(编译异常)
 * 4.对一个二进制移动一个负的值产生异常验证————(编译异常)
 * 5.整数显式转换为枚举，将过大值，负值转为枚举类型则抛出异常
 * 6.调用内部函数类型的零初始化变量验证————(编译异常)
 * 7.用assert的参数为false产生异常验证
 * @description:
 * @author: albedo
 * @create: 2019/12/31
 */
public class AssertHandleTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "intChangeException",
            author = "albedo", showName = "exceptionhandle.AssertHandle-整数显式转换为枚举", sourcePrefix = "evm")
    public void testIntChangeException() {
        try {
            prepare();
            AssertHandle handle = AssertHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("AssertHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + handle.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = handle.intChangeException(new BigInteger("2")).send();
            collector.logStepPass("checkout integer contains enums,transactionHah="+receipt.getTransactionHash());
            try {
                handle.intChangeException(new BigInteger("5")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout integer large then enums:" + e.getMessage());
            }
            try {
                handle.intChangeException(new BigInteger("-1")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout integer less then enums:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("AssertHandleTest testIntChangeException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "paramException",
            author = "albedo", showName = "exceptionhandle.AssertHandle-调用assert的参数为false", sourcePrefix = "evm")
    public void testParamException() {
        try {
            prepare();
            AssertHandle handle = AssertHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("AssertHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + handle.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = handle.paramException(new BigInteger("5")).send();
            collector.logStepPass("checkout normal,transactionHash="+receipt.getTransactionHash());
            try {
                handle.paramException(new BigInteger("11")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout assert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("AssertHandleTest testParamException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
