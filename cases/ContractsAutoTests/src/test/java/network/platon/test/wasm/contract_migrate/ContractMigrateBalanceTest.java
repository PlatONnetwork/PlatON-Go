package network.platon.test.wasm.contract_migrate;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractMigrate_v1;
import org.junit.Test;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Arrays;

/**
 * @title contract migrate
 * @description:
 * @author: yuanwenjun
 * @create: 2020/02/12
 */
public class ContractMigrateBalanceTest extends WASMContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "yuanwenjun", showName = "wasm.contract_migrate",sourcePrefix = "wasm")
    public void testMigrateContractBalance() {

        Uint64 transfer_value = Uint64.of(100000L);
        BigInteger origin_contract_value = BigInteger.valueOf(10000);

        try {
            prepare();
            ContractMigrate_v1 contractMigratev1 = ContractMigrate_v1.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractMigratev1.getContractAddress();
            String transactionHash = contractMigratev1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractMigrateV1 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractMigratev1.getTransactionReceipt().get().getGasUsed());
            
            Transfer transfer = new Transfer(web3j, transactionManager);
            transfer.sendFunds(contractAddress, new BigDecimal(origin_contract_value), Convert.Unit.VON).send();
            BigInteger originBalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("origin contract balance is: " + originBalance);

            String code = WasmFunctionEncoder.encodeConstructor(contractMigratev1.getContractBinary(), Arrays.asList());
            byte[] data = Numeric.hexStringToByteArray(code);

            TransactionReceipt transactionReceipt = contractMigratev1.migrate_contract(data,transfer_value, Uint64.of(90000000L)).send();
            collector.logStepPass("Contract Migrate V1  successfully hash:" + transactionReceipt.getTransactionHash());
            
            BigInteger originAfterMigrateBalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("After migrate, origin contract balance is: " + originAfterMigrateBalance);
            collector.assertEqual(originAfterMigrateBalance, BigInteger.valueOf(0), "checkout origin contract balance");
            
            String newContractAddress = contractMigratev1.getTransferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);
            BigInteger newMigrateBalance = web3j.platonGetBalance(newContractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("new contract balance is: " + newMigrateBalance);
            collector.assertEqual(newMigrateBalance, origin_contract_value.add(transfer_value.getValue()), "checkout new contract balance");

        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
