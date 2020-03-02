package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.TweetAccount;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class exists for docs.
 */
public class ContractTweetAccountTest extends WASMContractPrepareTest {

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

            //  get owner
            WasmAddress owner = contract.getOwnerAddress().send();
            collector.logStepPass("Call getOwnerAddress, owner: " + owner.getAddress());
            collector.assertEqual(owner.getAddress(), credentials.getAddress());

            // isAdmin
            Boolean isAdmin = contract.isAdmin().send();
            collector.logStepPass("Call isAdmin, res: " + isAdmin.toString());
            collector.assertTrue(isAdmin.booleanValue());

            // create tweet
            TransactionReceipt tweetTR = contract.tweet("Hello bob is bob.").send();
            collector.logStepPass("Send tweet, hash: " + tweetTR.getTransactionHash() + " gasUsed: " + tweetTR.getGasUsed().toString());

            // call getTweet
            String tweet = contract.getTweet(Uint64.of(0)).send();
            collector.logStepPass("Call getTweet, res: " + tweet);

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("Fibonacci and could not call contract function");
            }else{
                collector.logStepFail("Fibonacci failure,exception msg:" , e.getMessage());
            }
            e.printStackTrace();
        }
    }


}
