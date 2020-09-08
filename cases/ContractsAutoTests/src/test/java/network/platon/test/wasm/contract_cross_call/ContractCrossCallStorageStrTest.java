package network.platon.test.wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractCrossCallStorageString;
import network.platon.contracts.wasm.ContractStorageString;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

public class ContractCrossCallStorageStrTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call_storage_str",sourcePrefix = "wasm")
    public void testCrossCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `storge_str`, first
            ContractStorageString strc = ContractStorageString.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy storge_str contract:" + strc.getTransactionReceipt().get().getGasUsed());

            String strcAddr = strc.getContractAddress();
            String strcTxHash = strc.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("storge_str deployed sucessfully, contractAddress:" + strcAddr + ", txHash:" + strcTxHash);


            // deploy the cross_call_storage_str  contract second
            ContractCrossCallStorageString crossCall = ContractCrossCallStorageString.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy cross_call_storage_str contract:" + crossCall.getTransactionReceipt().get().getGasUsed());

            String crossCallAddr = crossCall.getContractAddress();
            String crossCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("cross_call_storage_str deployed sucessfully, contractAddress:" + crossCallAddr + ", txHash:" + crossCallTxHash);


            // check str value 1st
            String strcStr = strc.get_string().send();
            collector.logStepPass("the msg count in arr of  storge_str contract:" + strcStr);
            collector.assertEqual(strcStr, "");

            String crossCallStr = crossCall.get_string().send();
            collector.logStepPass("the msg count in arr of cross_call_storage_str contract:" + crossCallStr);
            collector.assertEqual(crossCallStr, "");

            String msg = "Gavin";

            // cross call contract start
            TransactionReceipt receipt = crossCall.call_set_string(strcAddr, msg, Uint64.of(0), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_storage_str call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check str value 2nd
            strcStr = strc.get_string().send();
            collector.logStepPass("the msg count in arr of  storge_str contract:" + strcStr);
            collector.assertEqual(strcStr, msg);

            crossCallStr = crossCall.get_string().send();
            collector.logStepPass("the msg count in arr of cross_call_storage_str contract:" + crossCallStr);
            collector.assertEqual(crossCallStr, "");

        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_call_storage_str Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
