package wasm.contract_event;

import com.platon.rlp.datatypes.Uint32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractEmitEvent2;
import network.platon.contracts.wasm.ContractEmitEvent3;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.util.List;

/**
 * @title PLATON_EVENT 测试3个主题
 * @description:
 * @author: hudenian
 * @create: 2020/02/10
 */
public class ContractEmitEvent3Test extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event合约3个主题事件",sourcePrefix = "wasm")
    public void testThreeEventContract() {

        String name = "hudenian";
        Uint32 value = Uint32.of(1L);
        String nationality = "myNationality";
        String city = "shanghai";
        try {
            prepare();
            ContractEmitEvent3 contractEmitEvent3 = ContractEmitEvent3.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractEmitEvent3.getContractAddress();
            String transactionHash = contractEmitEvent3.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contractEmitEvent3 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractEmitEvent3.getTransactionReceipt().get().getGasUsed());

            //调用包含3个主题事件的合约
            TransactionReceipt transactionReceipt = contractEmitEvent3.three_emit_event3(name,nationality,city,value).send();
            collector.logStepPass("contractEmitEvent3 call zero_emit_event successfully hash:" + transactionReceipt.getTransactionHash());

            List<ContractEmitEvent3.TransferEventResponse> eventList = contractEmitEvent3.getTransferEvents(transactionReceipt);
            String data = eventList.get(0).log.getData();
            collector.assertEqual(eventList.get(0).arg1,value);
            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0).toString());


        } catch (Exception e) {
            collector.logStepFail("ContractEmitEvent3Test failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
