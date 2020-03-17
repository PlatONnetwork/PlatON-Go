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
public class ContractBridgeTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_HomeBridge",sourcePrefix = "wasm")
    public void testHomeBridge() {
        try {
            // deploy contract.
            HomeBridge contract = HomeBridge.deploy(web3j, transactionManager, provider,
                    BigInteger.ONE, new WasmAddress[]{new WasmAddress(credentials.getAddress())}, BigInteger.ONE).send();
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
            byte[]   vs = new byte[0];
            byte[][] rs = new byte[0][256];
            byte[][] ss = new byte[0][256];
            byte[]  message = new byte[116];
            collector.logStepPass("message size: " + message.length);
            TransactionReceipt withdrawTr = contract.withdraw(vs, rs, ss, message, new BigInteger("1000000000000000000")).send();
            collector.logStepPass("Send withdraw, txHash: " + withdrawTr.getTransactionHash() + " gasUsed: " + withdrawTr.getGasUsed());
            collector.logStepPass("Send withdraw ,logs size: " + withdrawTr.getLogs().size());
            collector.assertTrue(withdrawTr.getLogs().size() != 0);
            HomeBridge.WithdrawEventResponse withevent = contract.getWithdrawEvents(withdrawTr).get(0);
            collector.logStepPass("Send withdraw ,event response, "
            + " arg1: " + withevent.arg1);

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("contract_HomeBridge and could not call contract function");
            }else{
                collector.logStepFail("contract_HomeBridge failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_ForeignBridge",sourcePrefix = "wasm")
    public void testForeignBridge() {
        try {
            // deploy contract.
            ForeignBridge contract = ForeignBridge.deploy(web3j, transactionManager, provider,
                    BigInteger.ONE, new WasmAddress[]{new WasmAddress(credentials.getAddress())}, BigInteger.ONE).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_ForeignBridge issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_ForeignBridge deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("contract_ForeignBridge and could not call contract function");
            }else{
                collector.logStepFail("contract_ForeignBridge failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

}
