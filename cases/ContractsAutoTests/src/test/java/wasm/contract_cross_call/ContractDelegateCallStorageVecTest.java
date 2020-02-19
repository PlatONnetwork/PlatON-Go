package wasm.contract_cross_call;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDelegateCallStorageVector;
import network.platon.contracts.wasm.ContractStorageVector;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

public class ContractDelegateCallStorageVecTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_storage_vector",sourcePrefix = "wasm")
    public void testDelegateCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `storage_vec`, first
            ContractStorageVector target = ContractStorageVector.deploy(web3j, transactionManager, provider).send();

            String storage_Addr = target.getContractAddress();
            String sotrage_vec_TxHash = target.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractTarget issued sucessfully, contractAddress:" + storage_Addr + ", txHash:" + sotrage_vec_TxHash);


            // deploy the delegate_call  contract second
            ContractDelegateCallStorageVector delegateCall = ContractDelegateCallStorageVector.deploy(web3j, transactionManager, provider).send();

            String delegateCallAddr = delegateCall.getContractAddress();
            String delegateCallTxHash = delegateCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDelegateCall issued sucessfully, contractAddress:" + delegateCallAddr + ", txHash:" + delegateCallTxHash);


            // check arr size 1st
            Long vecLen = target.get_vector_size().send();
            System.out.println("the msg count in arr of  storage_vec contract:" + vecLen);
            collector.assertEqual(vecLen.longValue(), 0l);

            Long delegateCallVecLen = delegateCall.get_vector_size().send();
            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallVecLen);
            collector.assertEqual(delegateCallVecLen.longValue(), 0l);

            // delegate call contract start
            ContractDelegateCallStorageVector.My_message myMessage = new ContractDelegateCallStorageVector.My_message();
            myMessage.baseClass = new ContractDelegateCallStorageVector.Message();
            myMessage.baseClass.head = "Gavin Head";
            myMessage.body = "Gavin Body";
            myMessage.end = "Gavin End";

            TransactionReceipt receipt = delegateCall.delegate_call_add_message(storage_Addr, myMessage, 60000000l).send();
            collector.logStepPass("ContractDelegateCall call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            vecLen = target.get_vector_size().send();
            System.out.println("the msg count in arr of  storage_vec contract:" + vecLen);
            collector.assertEqual(vecLen.longValue(), 0l);

            delegateCallVecLen = delegateCall.get_vector_size().send();
            System.out.println("the msg count in arr of delegateCall contract:" + delegateCallVecLen);
            collector.assertEqual(delegateCallVecLen.longValue(), 1l);

        } catch (Exception e) {
            collector.logStepFail("Failed to DelegateCall Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
