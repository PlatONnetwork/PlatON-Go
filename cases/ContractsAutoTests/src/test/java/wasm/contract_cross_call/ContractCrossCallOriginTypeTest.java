package wasm.contract_cross_call;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.*;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

public class ContractCrossCallOriginTypeTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call",sourcePrefix = "wasm")
    public void testCrossCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `storge_origin`, first
            ContractOriginType origin = ContractOriginType.deploy(web3j, transactionManager, provider).send();

            String originAddr = origin.getContractAddress();
            String originTxHash = origin.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("storge_str issued sucessfully, contractAddress:" + originAddr + ", txHash:" + originTxHash);


            // deploy the cross_call_storage_str  contract second
            ContractCrossCallOriginType crossCall = ContractCrossCallOriginType.deploy(web3j, transactionManager, provider).send();

            String crossCallAddr = crossCall.getContractAddress();
            String crossCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("cross_call_origin_type issued sucessfully, contractAddress:" + crossCallAddr + ", txHash:" + crossCallTxHash);


            // check vec size 1st
            Long originVecSize = origin.get_vector_size().send();
            System.out.println("the msg count in arr of  storge_origin contract:" + originVecSize);
            collector.assertEqual(originVecSize, 0l);

            Long crossCallVecSize = crossCall.get_vector_size().send();
            System.out.println("the msg count in arr of cross_call_origin_type contract:" + crossCallVecSize);
            collector.assertEqual(crossCallVecSize, 0l);

            // delegate call contract start
            ContractCrossCallOriginType.My_message myMessage = new ContractCrossCallOriginType.My_message();
            myMessage.baseClass = new ContractCrossCallOriginType.Message();
            myMessage.baseClass.head = "Gavin Head";
            myMessage.body = "Gavin Body";
            myMessage.end = "Gavin End";

            // cross call contract start
            TransactionReceipt receipt = crossCall.cross_call_add_message(originAddr, myMessage, 0l, 60000000l).send();
            collector.logStepPass("cross_call_origin_type call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            originVecSize = origin.get_vector_size().send();
            System.out.println("the msg count in arr of  storge_origin contract:" + originVecSize);
            collector.assertEqual(originVecSize, 1l);

            crossCallVecSize = crossCall.get_vector_size().send();
            System.out.println("the msg count in arr of cross_call_origin_type contract:" + crossCallVecSize);
            collector.assertEqual(crossCallVecSize, 0l);

        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_call_origin_type Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
