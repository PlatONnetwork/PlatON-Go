package network.platon.test.wasm.contract_event;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractEmitEvent;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.util.List;

/**
 * @title PLATON_EVENT 测试零个主题
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class ContractEmitEventTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event合约0个主题事件",sourcePrefix = "wasm")
    public void testZeroEventContract() {

        String name = "hudenian";
        try {
            prepare();
            ContractEmitEvent contractEmitEvent = ContractEmitEvent.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractEmitEvent.getContractAddress();
            String transactionHash = contractEmitEvent.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractEmitEvent issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractEmitEvent.getTransactionReceipt().get().getGasUsed());

            //调用包含零个主题事件的合约
            TransactionReceipt transactionReceipt = contractEmitEvent.zero_emit_event(name).send();
            collector.logStepPass("ContractEmitEvent call zero_emit_event successfully hash:" + transactionReceipt.getTransactionHash());

            List<ContractEmitEvent.TransferEventResponse> eventList = contractEmitEvent.getTransferEvents(transactionReceipt);
            String data = eventList.get(0).log.getData();
            collector.assertEqual(eventList.get(0).arg1,name);


        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
