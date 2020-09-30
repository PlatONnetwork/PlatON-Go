package network.platon.test.evm.exceptionhandle;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.RevertHandle;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;

import java.math.BigInteger;

/**
 * @title revert函数测试
 * 1.revert()函数————终止运行并撤销状态更改————验证
 * 2.revert(string reason)函数————终止运行并撤销状态更改,并提供一个解释性的字符串————验证
 * @description:
 * @author: albedo
 * @create: 2019/12/31
 */
public class RevertHandleTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "revertCheck",
            author = "albedo", showName = "exceptionhandle.RevertHandle-revert()函数", sourcePrefix = "evm")
    public void testRevertCheck() {
        try {
            prepare();
            RevertHandle handle = RevertHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RevertHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + handle.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt =handle.revertCheck(new BigInteger("5")).send();
            collector.logStepPass("checkout revert normal,transactionHah="+receipt.getTransactionHash());
            try {
                handle.revertCheck(new BigInteger("11")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout revert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RevertHandleTest testRevertCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "revertReasonCheck",
            author = "albedo", showName = "exceptionhandle.RevertHandle-revert(string reason)函数", sourcePrefix = "evm")
    public void testParamException() {
        try {
            prepare();
            RevertHandle handle = RevertHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RevertHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + handle.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt =handle.revertReasonCheck(new BigInteger("5")).send();
            collector.logStepPass("checkout revert normal,transactionHah="+receipt.getTransactionHash());
            try {
                handle.revertReasonCheck(new BigInteger("11")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout revert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RevertHandleTest testParamException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
