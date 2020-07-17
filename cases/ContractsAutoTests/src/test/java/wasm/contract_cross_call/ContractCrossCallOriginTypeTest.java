//package wasm.contract_cross_call;
//
//import com.platon.rlp.datatypes.Uint64;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.wasm.*;
//import org.junit.Test;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import wasm.beforetest.WASMContractPrepareTest;
//
//public class ContractCrossCallOriginTypeTest extends WASMContractPrepareTest {
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "xujiacan", showName = "wasm.contract_cross_call_origin_type",sourcePrefix = "wasm")
//    public void testCrossCallContract() {
//
//        try {
//            prepare();
//
//            // deploy the target contract which the name is `storge_origin`, first
//            ContractOriginType origin = ContractOriginType.deploy(web3j, transactionManager, provider, chainId).send();
//            collector.logStepPass("gas used after deploy storge_origin contract:" + origin.getTransactionReceipt().get().getGasUsed());
//
//
//            String originAddr = origin.getContractAddress();
//            String originTxHash = origin.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("storge_str deployed sucessfully, contractAddress:" + originAddr + ", txHash:" + originTxHash);
//
//
//            // deploy the cross_call_origin_type  contract second
//            ContractCrossCallOriginType crossCall = ContractCrossCallOriginType.deploy(web3j, transactionManager, provider, chainId).send();
//            collector.logStepPass("gas used after deploy cross_call_origin_type contract:" + crossCall.getTransactionReceipt().get().getGasUsed());
//
//
//            String crossCallAddr = crossCall.getContractAddress();
//            String crossCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("cross_call_origin_type deployed sucessfully, contractAddress:" + crossCallAddr + ", txHash:" + crossCallTxHash);
//
//
//            // check vec size 1st
//            Uint64 originVecSize = origin.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of  storge_origin contract:" + originVecSize);
//            collector.assertEqual(originVecSize.getValue().longValue(), 0l);
//
//            Uint64 crossCallVecSize = crossCall.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of cross_call_origin_type contract:" + crossCallVecSize);
//            collector.assertEqual(crossCallVecSize.getValue().longValue(), 0l);
//
//            // delegate call contract start
//            ContractCrossCallOriginType.My_message myMessage = new ContractCrossCallOriginType.My_message();
//            myMessage.baseClass = new ContractCrossCallOriginType.Message();
//            myMessage.baseClass.head = "Gavin Head";
//            myMessage.body = "Gavin Body";
//            myMessage.end = "Gavin End";
//
//            // cross call contract start
//            TransactionReceipt receipt = crossCall.cross_call_add_message(originAddr, myMessage, Uint64.of(0), Uint64.of(60000000l)).send();
//            collector.logStepPass("cross_call_origin_type call_add_message successfully txHash:" + receipt.getTransactionHash());
//
//
//            // check arr size 2nd
//            originVecSize = origin.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of  storge_origin contract:" + originVecSize);
//            collector.assertEqual(originVecSize.getValue().longValue(), 1l);
//
//            crossCallVecSize = crossCall.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of cross_call_origin_type contract:" + crossCallVecSize);
//            collector.assertEqual(crossCallVecSize.getValue().longValue(), 0l);
//
//        } catch (Exception e) {
//            collector.logStepFail("Failed to call cross_call_origin_type Contract,exception msg:" , e.getMessage());
//            e.printStackTrace();
//        }
//    }
//}
