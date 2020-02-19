package wasm.contract_migrate;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.utils.FileUtil;
import network.platon.contracts.wasm.ContractMigrate_new;
import network.platon.contracts.wasm.ContractMigrate_old;
import network.platon.contracts.wasm.ContractMigrate_v1;
import network.platon.utils.RlpUtil;
import org.junit.Test;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.tx.gas.ContractGasProvider;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import wasm.beforetest.WASMContractPrepareTest;

import java.io.File;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.file.Paths;
import java.util.ArrayList;
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

        Byte oldval = 12;

        try {
            prepare();
            provider = new ContractGasProvider(BigInteger.valueOf(50000000004L), BigInteger.valueOf(90000000L));
            ContractMigrate_old contractMigrateOld = ContractMigrate_old.deploy(web3j, transactionManager, provider, oldval).send();
            String contractAddress = contractMigrateOld.getContractAddress();
            String transactionHash = contractMigrateOld.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractMigrateVariableTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            //contractMigrateOld.setUint8(new Byte((byte)13)).send();

            Byte varval = contractMigrateOld.getUint8().send();
            collector.logStepPass("ContractMigrateVariableTest old contract variable value:" + varval);

            Byte newval = 23;
            Integer newvar = 26;
            String code = WasmFunctionEncoder.encodeConstructor(ContractMigrate_new.BINARY, Arrays.asList(newval, newvar));
            byte[] data = Numeric.hexStringToByteArray(code);
            TransactionReceipt transactionReceipt = contractMigrateOld.migrate_contract(data,0L, 90000000L).send();
            collector.logStepPass("ContractMigrateVariableTest migrate successfully hash:" + transactionReceipt.getTransactionHash());

            String newContractAddress = contractMigrateOld.getTransferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);

            ContractMigrate_new new_contractMigrate = ContractMigrate_new.load(newContractAddress,web3j,credentials,provider);
            Byte newContractval = new_contractMigrate.getUint8New().send();
            collector.logStepPass("new Contract origin variable value is:" + newContractval);
            collector.assertEqual(newContractval, newval, "checkout old variable of new contract value");
            short newVar = new_contractMigrate.getUint16().send();
            collector.logStepPass("new Contract new variable value is:" + newVar);
            collector.assertEqual(newVar, newvar, "checkout new variable of new contract value");
        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
