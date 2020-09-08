package wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractCrossCallStorageString;
import network.platon.contracts.wasm.ContractStorageString;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

public class ContractCrossCallBalanceTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call_balance",sourcePrefix = "wasm")
    public void testCrossCallContract() {

        try {
            prepare();

            // deploy the target contract which the name is `storge_str`, first
            ContractStorageString strc = ContractStorageString.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy storge_str contract:" + strc.getTransactionReceipt().get().getGasUsed());

            String strcAddr = strc.getContractAddress();
            String strcTxHash = strc.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("storge_str deployed sucessfully, contractAddress:" + strcAddr + ", txHash:" + strcTxHash);


            // deploy the cross_call_storage_str  contract second
            ContractCrossCallStorageString crossCall = ContractCrossCallStorageString.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy cross_call_storage_str contract:" + crossCall.getTransactionReceipt().get().getGasUsed());

            String crossCallAddr = crossCall.getContractAddress();
            String crossCallTxHash = crossCall.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("cross_call_storage_str deployed sucessfully, contractAddress:" + crossCallAddr + ", txHash:" + crossCallTxHash);


            // check contract balance 1st
            BigInteger strBalance = web3j.platonGetBalance(strcAddr, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("check contract balance 1st: the storage_str contract balance is:" + strBalance.toString());
            collector.assertEqual(strBalance.longValue(), 0l);

            BigInteger crosscallBalance = web3j.platonGetBalance(crossCallAddr, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("check contract balance 1st: the cross_call_balance contract balance is:" + crosscallBalance.toString());
            collector.assertEqual(crosscallBalance.longValue(), 0l);

            String transferMoneyStr = "1000";
            long transferMoneyL = Long.valueOf(transferMoneyStr).longValue();
            // transfer some balance to crosscall contract
            Transfer transfer = new Transfer(web3j, transactionManager);
            transfer.sendFunds(crossCallAddr, new BigDecimal(transferMoneyStr), Convert.Unit.VON).send();

            crosscallBalance = web3j.platonGetBalance(crossCallAddr, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("after transfer balance to crosscall contract: the cross_call_balance contract balance is:" + crosscallBalance.toString());
            collector.assertEqual(crosscallBalance.longValue(), transferMoneyL);


            long value = 100l;

            // cross call contract start
            TransactionReceipt receipt = crossCall.call_set_string(strcAddr, "", Uint64.of(value), Uint64.of(60000000l)).send();
            collector.logStepPass("cross_call_storage_str call_add_message successfully txHash:" + receipt.getTransactionHash());


            // check contract balance 2nd
            strBalance = web3j.platonGetBalance(strcAddr, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("check contract balance 2nd: the storage_str contract balance is:" + strBalance.toString());
            collector.assertEqual(strBalance.longValue(), value);

            crosscallBalance = web3j.platonGetBalance(crossCallAddr, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("check contract balance 2nd: the cross_call_balance contract balance is:" + crosscallBalance.toString());
            collector.assertEqual(crosscallBalance.longValue(), transferMoneyL-value);

        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_call_storage_str Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
