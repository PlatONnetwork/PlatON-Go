package wasm.contract_migrate;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractMigrate_types;
import network.platon.contracts.wasm.ContractMigrate_v1;
import org.junit.Test;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.gas.ContractGasProvider;
import org.web3j.utils.Numeric;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.Arrays;


/**
 * @title 合约升级
 * @description:
 * @author: yuanwenjun
 * @create: 2020/02/15
 */
public class ContractMigrateTypesTest extends WASMContractPrepareTest {

	@Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "yuanwenjun", showName = "wasm.contract_migrate",sourcePrefix = "wasm")
    public void testMigrateContract() {

        String name = "hello";

        try {
            prepare();
            provider = new ContractGasProvider(BigInteger.valueOf(50000000004L), BigInteger.valueOf(90000000L));
            ContractMigrate_types contractMigrateTypes = ContractMigrate_types.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contractMigrateTypes.getContractAddress();
            String transactionHash = contractMigrateTypes.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contractMigrateTypes issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            String structvalue = "testvalue";
            ContractMigrate_types.Message msg = new ContractMigrate_types.Message();
            msg.head = structvalue;
            contractMigrateTypes.setMessage(msg).send();
            
            Integer vecEle1 = 12, vecEle2 = 13;
            contractMigrateTypes.pushVector(vecEle1).send();
            contractMigrateTypes.pushVector(vecEle2).send();
            
            String mapKey1 = "key1", mapValue1= "value1", mapKey2 = "key2", mapValue2 = "value2";
            contractMigrateTypes.setMap(mapKey1, mapValue1).send();
            contractMigrateTypes.setMap(mapKey2, mapValue2).send();

            String code = WasmFunctionEncoder.encodeConstructor(contractMigrateTypes.getContractBinary(), Arrays.asList());
            byte[] data = Numeric.hexStringToByteArray(code);
            TransactionReceipt transactionReceipt = contractMigrateTypes.migrate_contract(data,0L, 90000000L).send();
            collector.logStepPass("contractMigrateTypes migrate successfully hash:" + transactionReceipt.getTransactionHash());

            String newContractAddress = contractMigrateTypes.getTransferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);

            ContractMigrate_types new_contractMigrate = ContractMigrate_types.load(newContractAddress,web3j,credentials,provider);
            ContractMigrate_types.Message newMsg = new_contractMigrate.getMessage().send();
            collector.logStepPass("new Contract message variable is:" + newMsg.head);
            collector.assertEqual(newMsg.head, structvalue, "check migrate struct value");
            
            Integer newVecEle1 = new_contractMigrate.getVectorElement(Long.valueOf(0)).send();
            Integer newVecEle2 = new_contractMigrate.getVectorElement(Long.valueOf(1)).send();
            collector.logStepPass("new Contract vector variable 0 is:" + newVecEle1);
            collector.logStepPass("new Contract vector variable 1 is:" + newVecEle2);
            collector.assertEqual(newVecEle1, vecEle1, "check vector variable 0");
            collector.assertEqual(newVecEle2, vecEle2, "check vector variable 1");
            
            String newMapValue1 = new_contractMigrate.getMapElement(mapKey1).send();
            String newMapValue2 = new_contractMigrate.getMapElement(mapKey2).send();
            collector.logStepPass("new Contract map value of key1 is:" + newMapValue1);
            collector.logStepPass("new Contract map value of key2 is:" + newMapValue2);
            collector.assertEqual(newMapValue1, mapValue1, "check map value of key1");
            collector.assertEqual(newMapValue2, mapValue2, "check map value of key2");
            
        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
