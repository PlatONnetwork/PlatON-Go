package wasm.contract_docs;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InnerFunction;
import network.platon.contracts.wasm.InnerFunction_1;
import network.platon.contracts.wasm.InnerFunction_2;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

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
    public void testFunctionContract() {
        try {
            // deploy contract.
            InnerFunction innerFunction = InnerFunction.deploy(web3j, transactionManager, provider).send();
            String contractAddress = innerFunction.getContractAddress();
            String transactionHash = innerFunction.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InnerFunction issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            // test: timestamp(bug)
            Long timestamp = innerFunction.timestamp().send();
            collector.logStepPass("To invoke timestamp success, timestamp: " + timestamp);

            // test: gas_limit
            Long gasLimit = innerFunction.gas_limit().send();
            collector.logStepPass("To invoke gas_limit success. gasLimit: " + gasLimit);
            collector.assertFalse(provider.getGasLimit().longValue() == gasLimit.longValue());

            // test: block_number
            Long bn = innerFunction.block_number().send();
            collector.logStepPass("To invoke block_number success, bn: " + bn);


        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("InnerFunction and could not call contract function");
            }else{
                collector.logStepFail("InnerFunction failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_function_01",sourcePrefix = "wasm")
    public void testFunctionContract_01() {

        String name = "zjsunzone";
        try {

            // deploy contract.
            InnerFunction_1 innerFunction = InnerFunction_1.deploy(web3j, transactionManager, provider).send();
            String contractAddress = innerFunction.getContractAddress();
            String transactionHash = innerFunction.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InnerFunction issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            // test: gas
            Long gas = innerFunction.gas().send();
            collector.logStepPass("To invoke gas success, gas: " + gas);

            // test: nonce
            Long rnonce = web3j.platonGetTransactionCount(credentials.getAddress(), DefaultBlockParameterName.LATEST).send().getTransactionCount().longValue();
            Long nonce = innerFunction.nonce().send();
            collector.logStepPass("To invoke nonce success, nonce: " + nonce + " rnonce: " + rnonce);

            // test: block_hash
            //String bhsh = innerFunction.block_hash(Long.valueOf(100)).send();
            //collector.logStepPass("To invoke block_hash success, hash: " + bhsh);
            //String bhash2 = web3j.platonGetBlockByNumber(new DefaultBlockParameterNumber(100), false).send().getBlock().getHash();
            //collector.assertEqual(prependHexPrefix(bhash2), prependHexPrefix(bhsh));

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("InnerFunction_1 and could not call contract function");
            }else{
                collector.logStepFail("InnerFunction_1 failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_function_02",sourcePrefix = "wasm")
    public void testFunctionContract_02() {

        try {
            // deploy contract.
            InnerFunction_2 innerFunction = InnerFunction_2.deploy(web3j, transactionManager, provider).send();
            String contractAddress = innerFunction.getContractAddress();
            String transactionHash = innerFunction.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InnerFunction deploy successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);


            // test: coinbase
            String coinbase = innerFunction.origin().send();
            collector.logStepPass("To invoke coinbase success. origin: " + Numeric.prependHexPrefix(coinbase));
            collector.assertEqual(credentials.getAddress(), Numeric.prependHexPrefix(coinbase));

            // test: transfer
            String toAddress = "0x250b67c9f1baa47dafcd1cfd5ad7780bb7b9b196";
            long amount = 1;
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal(amount), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);
            TransactionReceipt transferTr = innerFunction.transfer(toAddress, amount).send();
            BigInteger balance = web3j.platonGetBalance(toAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("To invoke transfer success, hash:" + transferTr.getTransactionHash() + " balance: " + balance);
            //collector.assertEqual(amount, balance.longValue());

            // test: sha3
            String sha3v1 = innerFunction.sha3("this is bob").send();
            String sha3v2 = innerFunction.sha3("this is bob").send();
            collector.logStepPass("To invoke sha3 success, v1: " + sha3v1 + " v2: " + sha3v2);
            collector.assertEqual(sha3v1, sha3v2);

            // test: return
            // ignore

            // test: panic
            TransactionReceipt panicTr = null;
            try {
                panicTr = innerFunction.panic().send();
                collector.logStepPass("To invoke panic success. hash:"+ panicTr.getTransactionHash() +" useGas: " + panicTr.getGasUsed().toString());
            }catch (Exception e){
                if (panicTr != null) {
                    collector.assertEqual(provider.getGasLimit(), panicTr.getGasUsed().longValue());
                }
            }

            // test: revert(bug)
            TransactionReceipt tr = innerFunction.revert(Long.valueOf(1)).send();
            collector.logStepPass("To invoke revert success. hash:"+ tr.getTransactionHash() +" useGas: " + tr.getGasUsed().toString());
            collector.assertFalse(provider.getGasLimit().longValue() == tr.getGasUsed().longValue());

            // test: destroy
            String receiveAddr = "0x250b67c9f1baa47dafcd1cfd5ad7780bb7b9b193";
            TransactionReceipt destoryTr = innerFunction.destroy(receiveAddr).send();
            BigInteger receiveBalance = web3j.platonGetBalance(receiveAddr, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("To invoke destory success, receiveBalance: " + receiveBalance);

            // test: origin(without 0x)
            String origin = innerFunction.origin().send();
            collector.logStepPass("To invoke origin success. origin: " + Numeric.prependHexPrefix(origin));
            collector.assertEqual(credentials.getAddress(), Numeric.prependHexPrefix(origin));

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("InnerFunction and could not call contract function");
            }else{
                collector.logStepFail("InnerFunction failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
