package wasm.contract_event;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.EventCallContract;
import network.platon.contracts.wasm.ContractEmitEvent;
import network.platon.contracts.wasm.ContractEmitEvent4;
import network.platon.contracts.wasm.Contract_panic;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import wasm.beforetest.WASMContractPrepareTest;

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
            ContractEmitEvent contractEmitEvent = ContractEmitEvent.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contractEmitEvent.getContractAddress();
            String transactionHash = contractEmitEvent.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractEmitEvent issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //调用包含零个主题事件的合约
            TransactionReceipt transactionReceipt = contractEmitEvent.zero_emit_event(name).send();
            collector.logStepPass("ContractEmitEvent call zero_emit_event successfully hash:" + transactionReceipt.getTransactionHash());

            List<ContractEmitEvent.Platon_event0_transferEventResponse> eventList = contractEmitEvent.getPlaton_event0_transferEvents(transactionReceipt);
            String data = eventList.get(0).log.getData();
            collector.assertEqual(eventList.get(0).arg1,name);


        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
