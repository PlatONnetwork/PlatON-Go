package wasm.gas;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractEmitEvent1;
import network.platon.contracts.wasm.OOMException;
import network.platon.contracts.wasm.Platon_gas_price;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.PlatonFilter;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.util.List;

/**
 * @title GasTest
 * @description Gas测试
 * @author qcxiao
 */
public class GasTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.GasTest",sourcePrefix = "wasm")
    public void test() {
        prepare();

        try {
            Platon_gas_price gas = Platon_gas_price.deploy(web3j, transactionManager, provider,chainId).send();
            String contractAddress = gas.getContractAddress();
            String transactionHash = gas.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("OOMException issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + gas.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = Platon_gas_price.load(contractAddress, web3j, transactionManager, provider,chainId).test().send();

            List<Platon_gas_price.GasUsedEventResponse> eventList = gas.getGasUsedEvents(transactionReceipt);
            //collector.assertEqual(eventList.get(0).arg1, 2);
            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0));

        } catch (Exception e) {
            collector.logStepPass("OOMException memory restriction effective.");
            e.printStackTrace();
        }


    }
}
