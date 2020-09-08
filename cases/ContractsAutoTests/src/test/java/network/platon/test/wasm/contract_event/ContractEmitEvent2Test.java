package wasm.contract_event;

import com.platon.rlp.datatypes.Uint32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractEmitEvent2;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.util.List;

/**
 * @title PLATON_EVENT 测试2个主题
 * @description:
 * @author: hudenian
 * @create: 2020/02/10
 */
public class ContractEmitEvent2Test extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event合约2个主题事件",sourcePrefix = "wasm")
    public void testTwoEventContract() {

        String name = "hudenian";
        Uint32 value = Uint32.of(1L);
        String nationality = "myNationality";
        try {
            prepare();
            ContractEmitEvent2 contractEmitEvent2 = ContractEmitEvent2.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractEmitEvent2.getContractAddress();
            String transactionHash = contractEmitEvent2.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contractEmitEvent2 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractEmitEvent2.getTransactionReceipt().get().getGasUsed());

            //调用包含2个主题事件的合约
            TransactionReceipt transactionReceipt = contractEmitEvent2.two_emit_event2(name,nationality,value).send();
            collector.logStepPass("contractEmitEvent2 call zero_emit_event successfully hash:" + transactionReceipt.getTransactionHash());

            List<ContractEmitEvent2.TransferEventResponse> eventList = contractEmitEvent2.getTransferEvents(transactionReceipt);
            String data = eventList.get(0).log.getData();
            collector.assertEqual(eventList.get(0).arg1,value);
            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0).toString());


        } catch (Exception e) {
            collector.logStepFail("ContractEmitEvent2Test failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
