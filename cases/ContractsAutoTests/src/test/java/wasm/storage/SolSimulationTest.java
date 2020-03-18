package wasm.storage;

import com.platon.rlp.datatypes.Uint32;
import com.platon.rlp.datatypes.Uint64;
import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Sol_simulation;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @title SolSimulationTest
 * @description 验证存储
 * @author qcxiao
 * @updateTime 2020/3/16 20:39
 */
public class SolSimulationTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.storage.SolSimulationTest", sourcePrefix = "wasm")
    public void test() {

        try {
            prepare();
            Sol_simulation solSimulation = Sol_simulation.deploy(web3j, transactionManager, provider, Uint32.of(32), Uint32.of(0)).send();
            String contractAddress = solSimulation.getContractAddress();
            String transactionHash = solSimulation.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("deploy successfully.contractAddress:" + contractAddress
                    + ", hash:" + transactionHash
                    + ", gasUsed:" + solSimulation.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceiptAction = Sol_simulation.load(contractAddress, web3j, transactionManager, provider).action().send();

            TransactionReceipt transactionReceipt = Sol_simulation.load(contractAddress, web3j, transactionManager, provider).debug().send();

            collector.logStepPass("action and debug are gas used: " + transactionReceiptAction.getGasUsed() + "," + transactionReceipt.getGasUsed());


        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("", e.getMessage());
        }
    }
}
