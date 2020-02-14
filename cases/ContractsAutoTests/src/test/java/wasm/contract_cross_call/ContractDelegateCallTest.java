package wasm.contract_cross_call;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDelegateCall;
import network.platon.contracts.wasm.ContractHello;
import network.platon.contracts.wasm.ContractTarget;
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
            ContractTarget target = ContractTarget.deploy(web3j, transactionManager, provider).send();

            String targetAddr = target.getContractAddress();
            String targetTxHash = target.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractTarget issued sucessfully, contractAddress:" + targetAddr + ", txHash:" + targetTxHash);


            // deploy the delegate_call  contract second
            ContractDelegateCall delegateCall = ContractDelegateCall.deploy(web3j, transactionManager, provider).send();

            String delegateCallAddr = delegateCall.getContractAddress();
            String delegateCallTxHash = delegateCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDelegateCall issued sucessfully, contractAddress:" + delegateCallAddr + ", txHash:" + delegateCallTxHash);


            // check arr size 1st
            Long targetArrLen = target.get_vector_size().send();
            System.out.println("the msg count in arr of  target contract:" + targetArrLen);
            collector.assertEqual(targetArrLen.longValue(), 0l);

            Long delegateCallArrLen = delegateCall.get_vector_size().send();
            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallArrLen);
            collector.assertEqual(delegateCallArrLen.longValue(), 0l);

            // delegate call contract start
            ContractDelegateCall.My_message myMessage = new ContractDelegateCall.My_message();
            myMessage.baseClass = new ContractDelegateCall.Message();
            myMessage.baseClass.head = "Gavin Head";
            myMessage.body = "Gavin Body";
            myMessage.end = "Gavin End";

            TransactionReceipt receipt = delegateCall.delegate_call_add_message(targetAddr, myMessage, 60000000l).send();
            collector.logStepPass("ContractDelegateCall call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            targetArrLen = target.get_vector_size().send();
            System.out.println("the msg count in arr of  target contract:" + targetArrLen);
            collector.assertEqual(targetArrLen.longValue(), 0l);

            delegateCallArrLen = delegateCall.get_vector_size().send();
            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallArrLen);
            collector.assertEqual(delegateCallArrLen.longValue(), 1l);

        } catch (Exception e) {
            collector.logStepFail("Failed to DelegateCall Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "xujiacan", showName = "wasm.contract_delegate_call",sourcePrefix = "wasm")
//    public void testDelegateCallContract() {
//
//        try {
//            prepare();
//
//            // deploy the target contract which the name is `hello`, first
//            ContractHello hello = ContractHello.deploy(web3j, transactionManager, provider).send();
//
//            String helloAddr = hello.getContractAddress();
//            String helloTxHash = hello.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("ContractHello issued sucessfully, contractAddress:" + helloAddr + ", txHash:" + helloTxHash);
//
//
//            // deploy the delegate_call  contract second
//            ContractDelegateCall delegateCall = ContractDelegateCall.deploy(web3j, transactionManager, provider).send();
//
//            String delegateCallAddr = delegateCall.getContractAddress();
//            String delegateCallTxHash = delegateCall.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("ContractDelegateCall issued sucessfully, contractAddress:" + delegateCallAddr + ", txHash:" + delegateCallTxHash);
//
//
//            // check arr size 1st
//            String helloStr = hello.get_string().send();
//            System.out.println("the msg count in arr of  hello contract:" + helloStr);
//
//
//            String delegateCallStr = delegateCall.get_string().send();
//            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallStr);
//
//
//
//
//            TransactionReceipt receipt = delegateCall.delegate_call_add_message(helloAddr, "Gavin", 60000000l).send();
//            collector.logStepPass("ContractDelegateCall call_add_message successfully txHash:" + receipt.getTransactionHash());
//
//
//            // check arr size 2nd
//            helloStr = hello.get_string().send();
//            System.out.println("the msg count in arr of  hello contract:" + helloStr);
//
//
//            delegateCallStr = delegateCall.get_string().send();
//            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallStr);
//
//        } catch (Exception e) {
//            collector.logStepFail("Failed to DelegateCall Contract,exception msg:" , e.getMessage());
//            e.printStackTrace();
//        }
//    }

}
