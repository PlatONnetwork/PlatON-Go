package wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractCrossCallStorageVector;
import network.platon.contracts.wasm.ContractStorageVector;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import wasm.beforetest.WASMContractPrepareTest;


public class ContractCrossCallStorageVecTest extends WASMContractPrepareTest {



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call_storage_vector",sourcePrefix = "wasm")
    public void testCrossCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `storage_vec`, first
            ContractStorageVector storage_vec = ContractStorageVector.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy storage_vec contract:" + storage_vec.getTransactionReceipt().get().getGasUsed());

            String storage_vec_Addr = storage_vec.getContractAddress();
            String storage_vec_TxHash = storage_vec.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractHello deployed sucessfully, contractAddress:" + storage_vec_Addr + ", txHash:" + storage_vec_TxHash);


            // deploy the cross_call  contract second
            ContractCrossCallStorageVector crossCall = ContractCrossCallStorageVector.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy cross_call_storage_vec contract:" + crossCall.getTransactionReceipt().get().getGasUsed());


            String crossCallAddr = crossCall.getContractAddress();
            String crossCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractCrossCall deployed sucessfully, contractAddress:" + crossCallAddr + ", txHash:" + crossCallTxHash);


            // check arr size 1st
            Uint64 storageVecLen = storage_vec.get_vector_size().send();
            collector.logStepPass("the msg count in arr of  storage_vec contract:" + storageVecLen);
            collector.assertEqual(storageVecLen.getValue().longValue(), 0l);

            Uint64 crossCallVecLen = crossCall.get_vector_size().send();
            collector.logStepPass("the msg count in arr of crossCall contract:" + crossCallVecLen);
            collector.assertEqual(crossCallVecLen.getValue().longValue(), 0l);


            // cross call contract start
            ContractCrossCallStorageVector.My_message myMessage = new ContractCrossCallStorageVector.My_message();
            myMessage.baseClass = new ContractCrossCallStorageVector.Message();
            myMessage.baseClass.head = "Gavin Head";
            myMessage.body = "Gavin Body";
            myMessage.end = "Gavin End";

            TransactionReceipt receipt = crossCall.call_add_message(storage_vec_Addr, myMessage, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("ContractCrossCall call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            storageVecLen = storage_vec.get_vector_size().send();
            collector.logStepPass("the msg count in arr of  storage_vec contract:" + storageVecLen);
            collector.assertEqual(storageVecLen.getValue().longValue(), 1l);

            crossCallVecLen = crossCall.get_vector_size().send();
            collector.logStepPass("the msg count in arr of crossCall contract:" + crossCallVecLen);
            collector.assertEqual(crossCallVecLen.getValue().longValue(), 0l);

        } catch (Exception e) {
            collector.logStepFail("Failed to CrossCall Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
