//package wasm.contract_event;
//
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.wasm.ContractEmitEvent3;
//import network.platon.contracts.wasm.ContractEmitEvent4;
//import org.junit.Test;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import wasm.beforetest.WASMContractPrepareTest;
//
//import java.util.List;
//
///**
// * @title PLATON_EVENT 测试4个主题
// * @description:
// * @author: hudenian
// * @create: 2020/02/10
// */
//public class ContractEmitEvent4Test extends WASMContractPrepareTest {
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "hudenian", showName = "wasm.contract_event合约4个主题事件",sourcePrefix = "wasm")
//    public void testFourEventContract() {
//
//        String name = "hudenian";
//        Integer value = 1;
//        String nationality = "myNationality";
//        String city = "shanghai";
//        String town = "pudongxinqu";
//        try {
//            prepare();
//            ContractEmitEvent4 contractEmitEvent4 = ContractEmitEvent4.deploy(web3j, transactionManager, provider).send();
//            String contractAddress = contractEmitEvent4.getContractAddress();
//            String transactionHash = contractEmitEvent4.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("contractEmitEvent4 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
//
//            //调用包含4个主题事件的合约
//            TransactionReceipt transactionReceipt = contractEmitEvent4.four_emit_event4(name,nationality,city,town,value).send();
//            collector.logStepPass("contractEmitEvent4 call zero_emit_event successfully hash:" + transactionReceipt.getTransactionHash());
//
//            List<ContractEmitEvent4.Platon_event4_transferEventResponse> eventList = contractEmitEvent4.getPlaton_event4_transferEvents(transactionReceipt);
//            //4个主题暂时不支持，先注释掉
////            String data = eventList.get(0).log.getData();
////            collector.assertEqual(eventList.get(0).arg1,value);
////            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0).toString());
//
//
//        } catch (Exception e) {
//            collector.logStepFail("ContractEmitEvent4Test failure,exception msg:" , e.getMessage());
//            e.printStackTrace();
//        }
//    }
//}
