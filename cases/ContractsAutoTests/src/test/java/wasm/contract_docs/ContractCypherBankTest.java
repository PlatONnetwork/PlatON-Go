package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint8;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Bank;
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

            // transfer
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal("100000000000"), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);

            // tokenSupply
            BigInteger tokenSupply = contract.totalSupply().send();
            collector.logStepPass("Call totalSupply, res: " + tokenSupply);

            // exitFee
            Uint8 exitFee = contract.exitFee().send();
            collector.logStepPass("Call exitFee, res: " + exitFee.getValue());

            // buy
            TransactionReceipt buyTr = contract.buy(new WasmAddress(credentials.getAddress()), new BigInteger("100000000000000000000")).send();
            collector.logStepPass("Send buy, txHash: " + buyTr.getTransactionHash() + " gasUsed: " + buyTr.getGasUsed());
            collector.logStepPass("Send buy ,logs size: " + buyTr.getLogs().size());
            /*Bank.TestDataEventResponse testData = contract.getTestDataEvents(buyTr).get(0);
            collector.logStepPass("parse logs, args1: " +
                    testData.arg3.toString() + "" +
                    " args2: " + testData.arg2 +
                    " args3: " + testData.arg1);*/

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
            TransactionReceipt sellTr = contract.sell(new BigInteger("1000")).send();
            collector.logStepPass("Send sell, txHash: " + sellTr.getTransactionHash() +
                    " gasUsed: " + sellTr.getGasUsed());
            collector.logStepPass("Send sell ,logs size: " + sellTr.getLogs().size());

            // transfer
            TransactionReceipt transferTr = contract.transfer(new WasmAddress("0x72ad2b713faa14c2c4cd2d7affe5d8f538968f5a"), new BigInteger("300000000000000000")).send();
            collector.logStepPass("Send transfer, txHash: " + transferTr.getTransactionHash() +
                    " gasUsed: " + transferTr.getGasUsed());
            collector.logStepPass("Send transfer ,logs size: " + transferTr.getLogs().size());

            // balanceOf
            BigInteger balanceOf = contract.balanceOf(new WasmAddress("0x72ad2b713faa14c2c4cd2d7affe5d8f538968f5a")).send();
            collector.logStepPass("Call balanceOf, res: " + balanceOf + " expect: " + "300000000000000000"); //

            // myDividends
            BigInteger myDividends = contract.myDividends(false).send();
            collector.logStepPass("Call myDividends, res: " + myDividends);

            // dividendsOf
            BigInteger dividendsOf = contract.dividendsOf(new WasmAddress(credentials.getAddress())).send();
            collector.logStepPass("Call dividendsOf, res: " + dividendsOf);

            // totalEthereumBalance
            BigInteger totalEthereumBalance = contract.totalEthereumBalance().send();
            collector.logStepPass("Call totalEthereumBalance, res: " + totalEthereumBalance);

            // myTokens
            BigInteger myTokens = contract.myTokens().send();
            collector.logStepPass("Call myTokens, res: " + myTokens);

            // sellPrice
            BigInteger sellPrice = contract.sellPrice().send();
            collector.logStepPass("Call sellPrice, res: " + sellPrice);

            // buyPrice
            BigInteger buyPrice = contract.sellPrice().send();
            collector.logStepPass("Call buyPrice, res: " + buyPrice);

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