package wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractCrossCall;
import network.platon.contracts.wasm.ContractHello;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import wasm.beforetest.WASMContractPrepareTest;


public class ContractCrossCallTest extends WASMContractPrepareTest {



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call",sourcePrefix = "wasm")
    public void testCrossCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `hello`, first
            ContractHello hello = ContractHello.deploy(web3j, transactionManager, provider).send();

            String helloAddr = hello.getContractAddress();
            String helloTxHash = hello.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractHello issued sucessfully, contractAddress:" + helloAddr + ", txHash:" + helloTxHash);


            // deploy the cross_call  contract second
            ContractCrossCall crossCall = ContractCrossCall.deploy(web3j, transactionManager, provider).send();

            String crossCallAddr = crossCall.getContractAddress();
            String crossCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractCrossCall issued sucessfully, contractAddress:" + crossCallAddr + ", txHash:" + crossCallTxHash);


            // check arr size 1st
            Uint64 helloArrLen = hello.get_vector_size().send();
            System.out.println("the msg count in arr of  hello contract:" + helloArrLen.getValue().toString());
            collector.assertEqual(helloArrLen.getValue().longValue(), 0l);

            Uint64 crossCallArrLen = crossCall.get_vector_size().send();
            System.out.println("the msg count in arr of crossCall contract:" + crossCallArrLen.getValue().toString());
            collector.assertEqual(crossCallArrLen.getValue().longValue(), 0l);


            // cross call contract start
            ContractCrossCall.My_message myMessage = new ContractCrossCall.My_message();
            myMessage.baseClass = new ContractCrossCall.Message();
            myMessage.baseClass.head = "Gavin Head";
            myMessage.body = "Gavin Body";
            myMessage.end = "Gavin End";

            TransactionReceipt receipt = crossCall.call_add_message(helloAddr, myMessage, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("ContractCrossCall call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            helloArrLen = hello.get_vector_size().send();
            System.out.println("the msg count in arr of  hello contract:" + helloArrLen);
            collector.assertEqual(helloArrLen.getValue().longValue(), 1l);

            crossCallArrLen = crossCall.get_vector_size().send();
            System.out.println("the msg count in arr of crossCall contract:" + crossCallArrLen);
            collector.assertEqual(crossCallArrLen.getValue().longValue(), 0l);

        } catch (Exception e) {
            collector.logStepFail("Failed to CrossCall Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "xujiacan", showName = "wasm.contract_cross_call",sourcePrefix = "wasm")
//    public void testCrossCallContract() {
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
//            // deploy the cross_call  contract second
//            ContractCrossCall crossCall = ContractCrossCall.deploy(web3j, transactionManager, provider).send();
//
//            String crossCallAddr = crossCall.getContractAddress();
//            String crossCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("ContractCrossCall issued sucessfully, contractAddress:" + crossCallAddr + ", txHash:" + crossCallTxHash);
//
//
//            // check arr size 1st
//            String helloArrLen = hello.get_string().send();
//            System.out.println("the msg count in arr of  hello contract:" + helloArrLen);
//
//
//            String crossCallArrLen = crossCall.get_string().send();
//            System.out.println("the msg count in arr of crossCall contract:" + crossCallArrLen);
//
//
//
//            // cross call contract start
//
//            TransactionReceipt receipt = crossCall.call_add_message(helloAddr, "Gavin", 0l, 60000000l).send();
//            collector.logStepPass("ContractCrossCall call_add_message successfully txHash:" + receipt.getTransactionHash());
//
//
//            // check arr size 2nd
//            helloArrLen = hello.get_string().send();
//            System.out.println("the msg count in arr of  hello contract:" + helloArrLen);
//
//
//            crossCallArrLen = crossCall.get_string().send();
//            System.out.println("the msg count in arr of crossCall contract:" + crossCallArrLen);
//
//        } catch (Exception e) {
//            collector.logStepFail("Failed to CrossCall Contract,exception msg:" , e.getMessage());
//            e.printStackTrace();
//        }
//    }

}
