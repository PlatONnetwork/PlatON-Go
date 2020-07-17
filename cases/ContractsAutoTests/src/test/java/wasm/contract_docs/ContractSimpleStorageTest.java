package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SimpleStorage;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class is for docs.
 */
public class ContractSimpleStorageTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_simple_storage",sourcePrefix = "wasm")
    public void testSimpleStorageContract() {
        try {
            // deploy contract.
            SimpleStorage contract = SimpleStorage.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SimpleStorage issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("SimpleStorage deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            TransactionReceipt tr = contract.set(Uint64.of(10000)).send();
            collector.logStepPass("To invoke set success, txHash: " + tr.getTransactionHash());
            Uint64 result = contract.get().send();
            collector.logStepPass("To invoke get success, result: " + result.value.toString());
            collector.assertEqual(result.value.toString(), "10000");

        } catch (Exception e) {
            collector.logStepFail("SimpleStorage failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
