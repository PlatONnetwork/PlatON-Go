package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.TweetAccount;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

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
            TweetAccount contract = TweetAccount.deploy(web3j, transactionManager, provider,chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("TweetAccount deploy successfully. contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("TweetAccount deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            //  get owner
            WasmAddress owner = contract.getOwnerAddress().send();
            collector.logStepPass("Call getOwnerAddress, owner: " + owner.getAddress());
            collector.assertEqual(owner.getAddress(), credentials.getAddress(chainId));

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

            TransactionReceipt tweet1TR = contract.tweet("Hello alice is alice.").send();
            collector.logStepPass("Send tweet1, hash: " + tweet1TR.getTransactionHash() + " gasUsed: " + tweet1TR.getGasUsed().toString());

            // call getTweet
            String tweet1 = contract.getTweet(Uint64.of(1)).send();
            collector.logStepPass("Call getTweet 1, res: " + tweet1);

            // call last tweet
            String lastTweet = contract.getLatestTweet().send();
            collector.logStepPass("Call getLatestTweet, res " + lastTweet);

            // call getNumberOfTweets
            Uint64 numberOfTweets = contract.getNumberOfTweets().send();
            collector.logStepPass("Call getNumberOfTweets, res: " + numberOfTweets.getValue().toString());

            // call contract addr
            WasmAddress caddr = contract.caddr().send();
            collector.logStepPass("Call caddr, res: " + caddr.getAddress());

            // transfer
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal(10), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);

            String caddrBalance = contract.caddrBalance(caddr).send();
            collector.logStepPass("Call caddrBalance, res: " + caddrBalance);

            // adminRetri
            WasmAddress receiver = new WasmAddress("lax1q0cwpg3x7zq6tkhvlk3z9jhujk0d0wqp3wg4ue");
            BigInteger receiveBalanceBefore = web3j.platonGetBalance(receiver.getAddress(), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Call balance, before res: " + receiveBalanceBefore);

//            WasmAddress adminAddr = new WasmAddress(credentials.getAddress(chainId));
            WasmAddress adminAddr = new WasmAddress("lax1fyeszufxwxk62p46djncj86rd553skpptsj8v6");

            TransactionReceipt adminTr = contract.adminRetrieveDonations(adminAddr).send();
            collector.logStepPass("Send adminRetrieveDonations, hash: " + adminTr.getTransactionHash() + " gasUsed: " + adminTr.getGasUsed().toString());
            BigInteger receiveBalanceAfter = web3j.platonGetBalance(receiver.getAddress(), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Call balance, after res: " + receiveBalanceAfter);

            // adminDelete...
            BigInteger ownerBalance = web3j.platonGetBalance(credentials.getAddress(chainId), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Owner balance, before res: " + ownerBalance);
            TransactionReceipt adminDeleteTr = contract.adminDeleteAccount().send();
            collector.logStepPass("Send adminDeleteAccount, hash: " + adminDeleteTr.getTransactionHash() + " gasUsed: " + adminDeleteTr.getGasUsed());
            ownerBalance = web3j.platonGetBalance(credentials.getAddress(chainId), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Owner balance, after res: " + ownerBalance);
        } catch (Exception e) {
            collector.logStepFail("TweetAccount failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
