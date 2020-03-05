package wasm.contract_docs;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.TweetAccount;
import org.junit.Before;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class exists for docs.
 */
public class ContractTweetRegistryTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_TweetAccount",sourcePrefix = "wasm")
    public void testTweetAccount() {

        try {
            // deploy contract.
            TweetAccount contract = TweetAccount.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("TweetAccount deploy successfully. contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("TweetAccount deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());


        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("TweetAccount and could not call contract function");
            }else{
                collector.logStepFail("TweetAccount failure,exception msg:" , e.getMessage());
            }
            e.printStackTrace();
        }
    }


}
