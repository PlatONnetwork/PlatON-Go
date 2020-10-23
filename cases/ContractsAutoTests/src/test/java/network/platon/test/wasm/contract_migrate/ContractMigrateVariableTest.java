package network.platon.test.wasm.contract_migrate;

import com.platon.rlp.datatypes.Uint16;
import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractMigrate_new;
import network.platon.contracts.wasm.ContractMigrate_old;
import org.junit.Test;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.utils.Numeric;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.util.Arrays;

/**
 * @title contract migrate
 * @description:
 * @author: yuanwenjun
 * @create: 2020/02/12
 */
public class ContractMigrateVariableTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "yuanwenjun", showName = "wasm.contract_migrate",sourcePrefix = "wasm")
    public void testMigrateContractBalance() {

        Uint8 oldval = Uint8.of(12);

        try {
            prepare();
            ContractMigrate_old contractMigrateOld = ContractMigrate_old.deploy(web3j, transactionManager, provider, chainId, oldval).send();
            String contractAddress = contractMigrateOld.getContractAddress();
            String transactionHash = contractMigrateOld.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractMigrateVariableTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractMigrateOld.getTransactionReceipt().get().getGasUsed());

            Uint8 varval = contractMigrateOld.getUint8().send();
            collector.logStepPass("ContractMigrateVariableTest old contract variable value:" + varval);

            Byte newval = 23;
            short newvar = 26;
            String code = WasmFunctionEncoder.encodeConstructor(ContractMigrate_new.BINARY, Arrays.asList(newval, newvar));
            byte[] data = Numeric.hexStringToByteArray(code);
            TransactionReceipt transactionReceipt = contractMigrateOld.migrate_contract(data, Uint64.of(0L), Uint64.of(90000000L)).send();
            collector.logStepPass("ContractMigrateVariableTest migrate successfully hash:" + transactionReceipt.getTransactionHash());

            String newContractAddress = contractMigrateOld.getTransferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);

            ContractMigrate_new new_contractMigrate = ContractMigrate_new.load(newContractAddress,web3j,credentials,provider, chainId);
            Uint8 newContractval = new_contractMigrate.getUint8New().send();
            collector.logStepPass("new Contract origin variable value is:" + newContractval);
            collector.assertEqual(newContractval.value.intValue(), Integer.valueOf(newval).intValue(), "checkout old variable of new contract value");
            Uint16 newVar = new_contractMigrate.getUint16().send();
            collector.logStepPass("new Contract new variable value is:" + newVar);
            collector.assertEqual(newVar.value.intValue(), Integer.valueOf(newvar).intValue(), "checkout new variable of new contract value");
        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
