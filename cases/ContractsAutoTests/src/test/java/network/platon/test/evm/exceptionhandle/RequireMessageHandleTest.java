package network.platon.test.evm.exceptionhandle;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.RequireMessageHandle;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;

import java.math.BigInteger;

/**
 * @title require(bool condition, string message)函数（该函数可以自定义message信息）测试
 * @description:
 * @author: albedo
 * @create: 2019/12/31
 */
public class RequireMessageHandleTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "exceptionhandle.RequireMessageHandle-require(bool condition, string message)函数", sourcePrefix = "evm")
    public void testParamException() {
        try {
            prepare();
            RequireMessageHandle handle = RequireMessageHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireMessageHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + handle.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = handle.paramException(new BigInteger("5")).send();
            collector.logStepPass("checkout require normal,transactionHash="+receipt.getTransactionHash());
            try {
                handle.paramException(new BigInteger("11")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout require throw exception:"+e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RequireMessageHandleTest testParamException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
