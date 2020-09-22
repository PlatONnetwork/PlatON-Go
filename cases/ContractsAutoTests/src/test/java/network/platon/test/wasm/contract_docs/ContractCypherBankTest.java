package network.platon.test.wasm.contract_docs;

import com.platon.rlp.datatypes.Uint128;
import com.platon.rlp.datatypes.Uint8;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.test.datatypes.Xuint128;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Bank;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
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
            Bank contract = Bank.deploy(web3j, transactionManager, provider,chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_CypherBank issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_CypherBank deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // transfer
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal("1"), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);

            // tokenSupply
            Uint128 tokenSupply = contract.totalSupply().send();
            collector.logStepPass("Call totalSupply, res: " + tokenSupply);

            // exitFee
            Uint8 exitFee = contract.exitFee().send();
            collector.logStepPass("Call exitFee, res: " + exitFee.getValue());

            // buy
//            WasmAddress receiver = new WasmAddress("lax1uqug0zq7rcxddndleq4ux2ft3tv6dqljphydrl");
            WasmAddress receiver = new WasmAddress("lax19xdjrg06xz9te85c839zqtelmaj2tgt047fh5q");

            TransactionReceipt buyTr = contract.buy(receiver, new BigInteger("100000000000000000000")).send();
            collector.logStepPass("Send buy, txHash: " + buyTr.getTransactionHash() + " gasUsed: " + buyTr.getGasUsed());
            collector.logStepPass("Send buy ,logs size: " + buyTr.getLogs().size());
            /*Bank.TestDataEventResponse testData = contract.getTestDataEvents(buyTr).get(0);
            collector.logStepPass("parse logs, args1: " +
                    testData.arg1.toString() + "" +
                    " args2: " + testData.arg2 +
                    " args3: " + testData.arg3);*/

            // IDD
            TransactionReceipt iddTr = contract.IDD().send();
            collector.logStepPass("Send IDD, txHash: " + iddTr.getTransactionHash() + " gasUsed: " + iddTr.getGasUsed());
            collector.logStepPass("Send IDD ,logs size: " + iddTr.getLogs().size());

            // DivsAddon
            //TransactionReceipt divsAddonTr = contract.DivsAddon().send();
            //collector.logStepPass("Send DivsAddon, txHash: " + divsAddonTr.getTransactionHash() + " gasUsed: " + divsAddonTr.getGasUsed());
            //collector.logStepPass("Send DivsAddon ,logs size: " + divsAddonTr.getLogs().size());

            // reinvest
            TransactionReceipt reinvestTr = contract.reinvest().send();
            collector.logStepPass("Send reinvest, txHash: " + reinvestTr.getTransactionHash() + " gasUsed: " + reinvestTr.getGasUsed());
            collector.logStepPass("Send reinvest ,logs size: " + reinvestTr.getLogs().size());

            // withdraw
            TransactionReceipt withdrawTr = contract.withdraw().send();
            collector.logStepPass("Send withdraw, txHash: " + withdrawTr.getTransactionHash() +
                    " gasUsed: " + withdrawTr.getGasUsed());
            collector.logStepPass("Send withdraw ,logs size: " + withdrawTr.getLogs().size());

            // sell
            TransactionReceipt sellTr = contract.sell(new Xuint128("1000")).send();
            collector.logStepPass("Send sell, txHash: " + sellTr.getTransactionHash() +
                    " gasUsed: " + sellTr.getGasUsed());
            collector.logStepPass("Send sell ,logs size: " + sellTr.getLogs().size());

            // transfer
            TransactionReceipt transferTr = contract.transfer(new WasmAddress("lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2"), new Xuint128("300000000000000000")).send();
            collector.logStepPass("Send transfer, txHash: " + transferTr.getTransactionHash() +
                    " gasUsed: " + transferTr.getGasUsed());
            collector.logStepPass("Send transfer ,logs size: " + transferTr.getLogs().size());

            // balanceOf
            Uint128 balanceOf = contract.balanceOf(new WasmAddress("lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2")).send();
            collector.logStepPass("Call balanceOf, res: " + balanceOf + " expect: " + "300000000000000000"); //

            // myDividends
            Uint128 myDividends = contract.myDividends(false).send();
            collector.logStepPass("Call myDividends, res: " + myDividends);

            // dividendsOf
//            BigInteger dividendsOf = contract.dividendsOf(new WasmAddress(credentials.getAddress(chainId))).send();
            Uint128 dividendsOf = contract.dividendsOf(new WasmAddress("lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2")).send();
            collector.logStepPass("Call dividendsOf, res: " + dividendsOf);

            // totalEthereumBalance
            Uint128 totalEthereumBalance = contract.totalEthereumBalance().send();
            collector.logStepPass("Call totalEthereumBalance, res: " + totalEthereumBalance);

            // myTokens
            Uint128 myTokens = contract.myTokens().send();
            collector.logStepPass("Call myTokens, res: " + myTokens);

            // sellPrice
            Uint128 sellPrice = contract.sellPrice().send();
            collector.logStepPass("Call sellPrice, res: " + sellPrice);

            // buyPrice
            Uint128 buyPrice = contract.sellPrice().send();
            collector.logStepPass("Call buyPrice, res: " + buyPrice);

        } catch (Exception e) {
            collector.logStepFail("contract_CypherBank failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
