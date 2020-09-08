package network.platon.test.wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDelegateCallStorageString;
import network.platon.contracts.wasm.ContractStorageString;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

public class ContractDelegateCallStorageStrTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_delegate_call_storage_str",sourcePrefix = "wasm")
    public void testDelegateCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `storge_str`, first
            ContractStorageString strc = ContractStorageString.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy storge_str contract:" + strc.getTransactionReceipt().get().getGasUsed());


            String strcAddr = strc.getContractAddress();
            String helloTxHash = strc.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("storge_str deployed sucessfully, contractAddress:" + strcAddr + ", txHash:" + helloTxHash);


            // deploy the delegate_call  contract second
            ContractDelegateCallStorageString delegateCall = ContractDelegateCallStorageString.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy delegate_call_storge_str contract:" + delegateCall.getTransactionReceipt().get().getGasUsed());

            String delegateCallAddr = delegateCall.getContractAddress();
            String delegateCallTxHash = delegateCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("cross_delegate_call_storage_str deployed sucessfully, contractAddress:" + delegateCallAddr + ", txHash:" + delegateCallTxHash);


            // check arr size 1st
            String strcStr = strc.get_string().send();
            collector.logStepPass("the msg count in arr of  storge_str contract:" + strcStr);
            collector.assertEqual(strcStr, "");

            String delegateCallStr = delegateCall.get_string().send();
            collector.logStepPass("the msg count in arr of cross_delegate_call_storage_str contract:" + delegateCallStr);
            collector.assertEqual(delegateCallStr, "");

            String msg = "Gavin";

            TransactionReceipt receipt = delegateCall.delegate_call_set_string(strcAddr, msg, Uint64.of(60000000l)).send();
            collector.logStepPass("cross_delegate_call_storage_str call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check arr size 2nd
            strcStr = strc.get_string().send();
            collector.logStepPass("the msg count in arr of  storge_str contract:" + strcStr);
            collector.assertEqual(strcStr, "");

            delegateCallStr = delegateCall.get_string().send();
            collector.logStepPass("the msg count in arr of cross_delegate_call_storage_str contract:" + delegateCallStr);
            collector.assertEqual(delegateCallStr, msg);

        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_delegate_call_storage_str Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
