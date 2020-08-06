package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint8;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Bank;
import network.platon.contracts.wasm.ForeignBridge;
import network.platon.contracts.wasm.HomeBridge;
import org.junit.Before;
import org.junit.Test;
import org.web3j.crypto.Wallet;
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
 * <p>
 * This class is for docs.
 */
public class ContractBridgeTest extends WASMContractPrepareTest {

    @Before
    public void before() {
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_HomeBridge", sourcePrefix = "wasm")
    public void testHomeBridge() {
        try {
            // deploy contract.
            WasmAddress wasmAddress1 = new WasmAddress("lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2");
            WasmAddress wasmAddress2 = new WasmAddress("lax1fyeszufxwxk62p46djncj86rd553skpptsj8v6");
            HomeBridge contract = HomeBridge.deploy(web3j, transactionManager, provider, chainId,
                    BigInteger.ONE, new WasmAddress[]{wasmAddress1, wasmAddress2}, BigInteger.ONE).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_HomeBridge issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_HomeBridge deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // transfer to contract
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal(100), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);

            //
            TransactionReceipt setTr = contract.setGasLimitWithdrawRelay(new BigInteger("100000")).send();
            collector.logStepPass("Send setGasLimitWithdrawRelay, txHash: " + setTr.getTransactionHash()
                    + " gasUsed: " + setTr.getGasUsed());
            collector.logStepPass("Send setGasLimitWithdrawRelay ,logs size: " + setTr.getLogs().size());
            HomeBridge.GasConsumptionLimitsUpdatedEventResponse setEvent = contract.getGasConsumptionLimitsUpdatedEvents(setTr).get(0);
            collector.logStepPass("Send setGasLimitWithdrawRelay, event args1:" + setEvent.arg1);

            // withdraw.
            byte[] vs = new byte[0];
            byte[][] rs = new byte[0][256];
            byte[][] ss = new byte[0][256];
            byte[] message = new byte[116];
            collector.logStepPass("message size: " + message.length);
            TransactionReceipt withdrawTr = contract.withdraw(vs, rs, ss, message, new BigInteger("1000000000000000000")).send();
            collector.logStepPass("Send withdraw, txHash: " + withdrawTr.getTransactionHash() + " gasUsed: " + withdrawTr.getGasUsed());
            collector.logStepPass("Send withdraw ,logs size: " + withdrawTr.getLogs().size());
            collector.assertTrue(withdrawTr.getLogs().size() != 0);
            HomeBridge.WithdrawEventResponse withevent = contract.getWithdrawEvents(withdrawTr).get(0);
            collector.logStepPass("Send withdraw ,event response, "
                    + " arg1: " + withevent.arg1);

        } catch (Exception e) {
            collector.logStepFail("contract_HomeBridge failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_ForeignBridge", sourcePrefix = "wasm")
    public void testForeignBridge() {
        try {
            // deploy contract.
            ForeignBridge contract = ForeignBridge.deploy(web3j, transactionManager, provider, chainId,
                    BigInteger.ONE, new WasmAddress[]{new WasmAddress(credentials.getAddress(chainId))}, BigInteger.ONE).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_ForeignBridge issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_ForeignBridge deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // setTokenAddress
            TransactionReceipt tokenTr = contract.setTokenAddress(new WasmAddress(contractAddress)).send();
            collector.logStepPass("Send setTokenAddress, txHash: " + tokenTr.getTransactionHash()
                    + " gasUsed: " + tokenTr.getGasUsed());
            collector.logStepPass("Send setTokenAddress ,logs size: " + tokenTr.getLogs().size());

            // deposit
            TransactionReceipt depositTr = contract.deposit(new WasmAddress(credentials.getAddress(chainId)), BigInteger.valueOf(1000000), new byte[]{}).send();
            collector.logStepPass("Send deposit, txHash: " + depositTr.getTransactionHash()
                    + " gasUsed: " + depositTr.getGasUsed());
            collector.logStepPass("Send deposit ,logs size: " + depositTr.getLogs().size());
            collector.assertTrue(depositTr.getLogs().size() != 0);

            // submitSignature
            TransactionReceipt submitSignatureTr = contract.submitSignature(new byte[]{}, new byte[116]).send();
            collector.logStepPass("Send submitSignature, txHash: " + submitSignatureTr.getTransactionHash()
                    + " gasUsed: " + submitSignatureTr.getGasUsed());
            collector.logStepPass("Send submitSignature ,logs size: " + submitSignatureTr.getLogs().size());
            collector.assertTrue(submitSignatureTr.getLogs().size() != 0);

            // signature
            byte[] signature = contract.signature(new byte[]{}, BigInteger.TEN).send();
            collector.logStepPass("Call signature, response: " + Numeric.toHexString(signature));


        } catch (Exception e) {
            collector.logStepFail("contract_ForeignBridge failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

}
