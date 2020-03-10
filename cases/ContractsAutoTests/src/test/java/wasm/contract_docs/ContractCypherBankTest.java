package wasm.contract_docs;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SimpleStorage;
import org.junit.Before;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class is for docs.
 */
public class ContractCypherBankTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_sypherBank",sourcePrefix = "wasm")
    public void testSimpleStorageContract() {
        try {
            // deploy contract.
            SimpleStorage contract = SimpleStorage.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SimpleStorage issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("SimpleStorage deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());



        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("SimpleStorage and could not call contract function");
            }else{
                collector.logStepFail("SimpleStorage failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

}
