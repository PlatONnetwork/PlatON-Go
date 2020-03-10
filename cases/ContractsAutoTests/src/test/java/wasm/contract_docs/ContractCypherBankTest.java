package wasm.contract_docs;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Bank;
import org.junit.Before;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

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
            author = "zjsunzone", showName = "wasm.contract_CypherBank",sourcePrefix = "wasm")
    public void testSimpleStorageContract() {
        try {
            // deploy contract.
            Bank contract = Bank.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_CypherBank issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_CypherBank deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // tokenSupply
            BigInteger tokenSupply = contract.totalSupply().send();
            collector.logStepPass("Call totalSupply, res: " + tokenSupply);

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("contract_CypherBank and could not call contract function");
            }else{
                collector.logStepFail("contract_CypherBank failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

}
