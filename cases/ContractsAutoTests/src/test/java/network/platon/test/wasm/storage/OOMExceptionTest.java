package network.platon.test.wasm.storage;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.OOMException;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title OOMExceptionTest
 * @description 内存溢出测试
 * @author qcxiao
 * @updateTime 2020/4/21 15:58
 */
public class OOMExceptionTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.OOMExceptionTest",sourcePrefix = "wasm")
    public void test() {
        prepare();

        try {
            OOMException oomException = OOMException.deploy(web3j, transactionManager, provider,chainId).send();
            String contractAddress = oomException.getContractAddress();
            String transactionHash = oomException.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("OOMException issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + oomException.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = OOMException.load(contractAddress, web3j, transactionManager, provider,chainId).memory_limit().send();
            collector.logStepFail("OOMExceptionTest fail" , "内存限制未生效");

        } catch (Exception e) {
            collector.logStepPass("OOMException memory restriction effective.");
            e.printStackTrace();
        }


    }
}
