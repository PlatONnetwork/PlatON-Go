package wasm.contract_cross_call;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractCrossCall;
import network.platon.contracts.wasm.ContractDelegateCall;
import network.platon.contracts.wasm.ContractHello;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

public class ContractDelegateCallTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call",sourcePrefix = "wasm")
    public void testDelegateCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `hello`, first
            ContractHello hello = ContractHello.deploy(web3j, transactionManager, provider).send();

            String helloAddr = hello.getContractAddress();
            String helloTxHash = hello.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractHello issued sucessfully, contractAddress:" + helloAddr + ", txHash:" + helloTxHash);


            // deploy the delegate_call  contract second
            ContractDelegateCall delegateCall = ContractDelegateCall.deploy(web3j, transactionManager, provider).send();

            String delegateCallAddr = delegateCall.getContractAddress();
            String delegateCallTxHash = delegateCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDelegateCall issued sucessfully, contractAddress:" + delegateCallAddr + ", txHash:" + delegateCallTxHash);


            // check arr size 1st
            Long helloArrLen = hello.get_vector_size().send();
            System.out.println("the msg count in arr of  hello contract:" + helloArrLen);
            collector.assertEqual(helloArrLen.longValue(), 0l);

            Long delegateCallArrLen = delegateCall.get_vector_size().send();
            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallArrLen);
            collector.assertEqual(delegateCallArrLen.longValue(), 0l);

            // delegate call contract start
            ContractDelegateCall.My_message myMessage = new ContractDelegateCall.My_message();
            myMessage.baseClass = new ContractDelegateCall.Message();
            myMessage.baseClass.head = "Gavin Head";
            myMessage.body = "Gavin Body";
            myMessage.end = "Gavin End";

            TransactionReceipt receipt = delegateCall.delegate_call_add_message(helloAddr, myMessage, 60000000l).send();
            collector.logStepPass("ContractDelegateCall call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            helloArrLen = hello.get_vector_size().send();
            System.out.println("the msg count in arr of  hello contract:" + helloArrLen);
            collector.assertEqual(helloArrLen.longValue(), 0l);

            delegateCallArrLen = delegateCall.get_vector_size().send();
            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallArrLen);
            collector.assertEqual(delegateCallArrLen.longValue(), 1l);

        } catch (Exception e) {
            collector.logStepFail("Failed to DelegateCall Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
