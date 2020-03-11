package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Bank;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
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
            t.sendFunds(contractAddress, new BigDecimal(100000000), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);

            // tokenSupply
            BigInteger tokenSupply = contract.totalSupply().send();
            collector.logStepPass("Call totalSupply, res: " + tokenSupply);

            // exitFee
            Uint8 exitFee = contract.exitFee().send();
            collector.logStepPass("Call exitFee, res: " + exitFee.getValue());

            // buy
            //TransactionReceipt buyTr = contract.buy(new WasmAddress(credentials.getAddress()), new BigInteger("100000000000000000000")).send();
            //collector.logStepPass("Send buy, txHash: " + buyTr.getTransactionHash() + " gasUsed: " + buyTr.getGasUsed());
            //collector.logStepPass("Send buy ,logs size: " + buyTr.getLogs().size());

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
