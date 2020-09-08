//package network.platon.test.wasm.contract_cross_call;
//
//
//import com.platon.rlp.network.platon.test.datatypes.Uint64;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.wasm.ContractDelegateCallOriginType;
//import network.platon.contracts.wasm.ContractOriginType;
//import org.junit.Test;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import network.platon.test.wasm.beforetest.WASMContractPrepareTest;
//
//public class ContractDelegateCallOriginTypeTest extends WASMContractPrepareTest {
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "xujiacan", showName = "wasm.contract_delegate_call_origin_type",sourcePrefix = "wasm")
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
//            // deploy the delegate_call_origin_type  contract second
//            ContractDelegateCallOriginType crossCall = ContractDelegateCallOriginType.deploy(web3j, transactionManager, provider, chainId).send();
//            collector.logStepPass("gas used after deploy delegate_call_origin_type contract:" + crossCall.getTransactionReceipt().get().getGasUsed());
//
//
//            String delegateCallAddr = crossCall.getContractAddress();
//            String delegateCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("delegate_call_origin_type deployed sucessfully, contractAddress:" + delegateCallAddr + ", txHash:" + delegateCallTxHash);
//
//
//            // check vec size 1st
//            Uint64 originVecSize = origin.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of  storge_origin contract:" + originVecSize);
//            collector.assertEqual(originVecSize.getValue().longValue(), 0l);
//
//            Uint64 delegateCallVecSize = crossCall.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of delegate_call_origin_type contract:" + delegateCallVecSize);
//            collector.assertEqual(delegateCallVecSize.getValue().longValue(), 0l);
//
//            // delegate call contract start
//            ContractDelegateCallOriginType.My_message myMessage = new ContractDelegateCallOriginType.My_message();
//            myMessage.baseClass = new ContractDelegateCallOriginType.Message();
//            myMessage.baseClass.head = "Gavin Head";
//            myMessage.body = "Gavin Body";
//            myMessage.end = "Gavin End";
//
//            // cross call contract start
//            TransactionReceipt receipt = crossCall.delegate_call_add_message(originAddr, myMessage, Uint64.of(60000000l)).send();
//            collector.logStepPass("delegate_call_origin_type call_add_message successfully txHash:" + receipt.getTransactionHash());
//
//
//            // check arr size 2nd
//            originVecSize = origin.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of  storge_origin contract:" + originVecSize);
//            collector.assertEqual(originVecSize.getValue().longValue(), 0l);
//
//            delegateCallVecSize = crossCall.get_vector_size().send();
//            collector.logStepPass("the msg count in arr of delegate_call_origin_type contract:" + delegateCallVecSize);
//            collector.assertEqual(delegateCallVecSize.getValue().longValue(), 1l);
//
//        } catch (Exception e) {
//            collector.logStepFail("Failed to call delegate_call_origin_type Contract,exception msg:" , e.getMessage());
//            e.printStackTrace();
//        }
//    }
//}
