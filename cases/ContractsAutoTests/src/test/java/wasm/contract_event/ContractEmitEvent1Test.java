package wasm.contract_event;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.EventCallContract;
import network.platon.contracts.wasm.ContractEmitEvent;
import network.platon.contracts.wasm.ContractEmitEvent1;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.util.List;

/**
 * @title PLATON_EVENT 测试1个主题
 * @description:
 * @author: hudenian
 * @create: 2020/02/11
 */
public class ContractEmitEvent1Test extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_event合约1个主题事件",sourcePrefix = "wasm")
    public void testZeroEventContract() {

        String name = "hudenian";
        Integer value = 1;
        try {
            prepare();
            ContractEmitEvent1 contractEmitEvent1 = ContractEmitEvent1.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contractEmitEvent1.getContractAddress();
            String transactionHash = contractEmitEvent1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractEmitEvent1 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //调用包含1个主题事件的合约
            TransactionReceipt transactionReceipt = contractEmitEvent1.one_emit_event1(name,value).send();
            collector.logStepPass("ContractEmitEvent1 call zero_emit_event successfully hash:" + transactionReceipt.getTransactionHash());

            //对事件信息进行解析
            List<ContractEmitEvent1.TransferEventResponse> eventList = contractEmitEvent1.getTransferEvents(transactionReceipt);
            String data = eventList.get(0).log.getData();
            collector.assertEqual(eventList.get(0).arg1,value);
            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0).toString());

        } catch (Exception e) {
            collector.logStepFail("ContractEmitEvent1Test failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
