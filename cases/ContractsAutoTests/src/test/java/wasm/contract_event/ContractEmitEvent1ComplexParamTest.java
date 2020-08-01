package wasm.contract_event;

import com.platon.rlp.datatypes.Uint32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractEmitEvent1;
import network.platon.contracts.wasm.ContractEmitEvent1ComplexParam;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.util.ArrayList;
import java.util.List;

/**
 * @title PLATON_EVENT 合约入参及事件包含List
 * @description:
 * @author: hudenian
 * @create: 2020/02/26
 */
public class ContractEmitEvent1ComplexParamTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.ContractEmitEvent1ComplexParamTest合约入参及事件包含List",sourcePrefix = "wasm")
    public void testComplexParamContract() {

        String name = "myName";
        Uint32 value = Uint32.of(1L);
        List<String> stringList = new ArrayList<String>();
        stringList.add("listOne");
        stringList.add("listTwo");

        try {
            prepare();
            ContractEmitEvent1ComplexParam contractEmitEvent1 = ContractEmitEvent1ComplexParam.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractEmitEvent1.getContractAddress();
            String transactionHash = contractEmitEvent1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractEmitEvent1ComplexParamTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractEmitEvent1.getTransactionReceipt().get().getGasUsed());

            //调用包含1个主题事件的合约传一个list
            TransactionReceipt transactionReceipt = contractEmitEvent1.one_emit_event1(name,value,stringList).send();
            collector.logStepPass("ContractEmitEvent1ComplexParamTest call one_emit_event1 successfully hash:" + transactionReceipt.getTransactionHash());

            //对事件信息进行解析
            List<ContractEmitEvent1ComplexParam.TransferEventResponse> eventList = contractEmitEvent1.getTransferEvents(transactionReceipt);
            String data = eventList.get(0).log.getData();

            //获取包含list集合的合约
            collector.assertEqual(eventList.get(0).arg1,value);
            collector.assertEqual(eventList.get(0).arg2.get(0),stringList.get(0));
            collector.assertEqual(eventList.get(0).arg2.get(1),stringList.get(1));

            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0).toString());

        } catch (Exception e) {
            collector.logStepFail("ContractEmitEvent1ComplexParamTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
