package network.platon.test.wasm.contract_clone;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.test.datatypes.Xuint128;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SimpleContract;
import network.platon.contracts.wasm.CreateContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;
import com.platon.rlp.datatypes.Int32;

/**
 * @author wanghengtao
 * <p>
 * This class is for docs.
 */
public class CreateContractTest extends WASMContractPrepareTest {

    @Before
    public void before() {
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "wanghengtao", showName = "wasm.CreateContract", sourcePrefix = "wasm")
    public void testCloneAndCreate() {
        try {
            // deploy contract.
            SimpleContract simpleContract = SimpleContract.deploy(web3j, transactionManager, provider, chainId).send();
            CreateContract createContract = CreateContract.deploy(web3j, transactionManager, provider, chainId).send();

            // test simple contract
            WasmAddress wasmCreateAddress = new WasmAddress(createContract.getContractAddress());
            simpleContract.set_address(wasmCreateAddress).send();
            WasmAddress simpleCreateAddress = simpleContract.get_address().send();
            collector.assertEqual(wasmCreateAddress, simpleCreateAddress);
            collector.logStepPass("The set_simple_address and get_simple_address methods have been verified to work.");

            // test create contract
            WasmAddress wasmSimpleAddress = new WasmAddress(simpleContract.getContractAddress());
            createContract.set_simple_address(wasmSimpleAddress).send();
            WasmAddress createSimpleAddress = createContract.get_simple_address().send();
            collector.assertEqual(wasmSimpleAddress, createSimpleAddress);

            // deploy simple contract
            TransactionReceipt receipt = createContract.deploy_contract(wasmSimpleAddress).send();
            collector.assertTrue(receipt.isStatusOK());
            WasmAddress deploySimple = createContract.get_deploy_address().send();
            collector.logStepPass("deploy contract successfully, contract address:" + deploySimple.toString());

            // get deploy simple contract code length
            Int32 deployLength = createContract.get_contract_length(deploySimple).send();
            Int32 simpleLength = createContract.get_contract_length(wasmSimpleAddress).send();
            collector.assertEqual(deployLength, simpleLength);
            collector.logStepPass("The deployment contract length is verified correctly");

            // test deploy simple contract
            SimpleContract deploySimpleContract = SimpleContract.load(deploySimple.toString(), web3j, transactionManager, provider, chainId);
            deploySimpleContract.set(Uint64.of(7)).send();
            Uint64 deploySimpleGet = deploySimpleContract.get().send();
            collector.assertEqual(deploySimpleGet, Uint64.of(7));
            collector.logStepPass("The deployment contract method is called normally");

            // clone simple contract
            receipt = createContract.clone_contract(wasmSimpleAddress).send();
            collector.assertTrue(receipt.isStatusOK());
            WasmAddress cloneSimple = createContract.get_clone_address().send();
            collector.logStepPass("clone contract successfully, contract address:" + cloneSimple.toString());

            // get clone simple contract code length
            Int32 cloneLength = createContract.get_contract_length(cloneSimple).send();
            collector.assertEqual(cloneLength, simpleLength);
            collector.logStepPass("The cloned contract length is verified correctly");

            // test clone simple contract
            SimpleContract cloneSimpleContract = SimpleContract.load(cloneSimple.toString(), web3j, transactionManager, provider, chainId);
            cloneSimpleContract.set(Uint64.of(7)).send();
            Uint64 cloneSimpleGet = cloneSimpleContract.get().send();
            collector.assertEqual(cloneSimpleGet, Uint64.of(7));
            collector.logStepPass("The cloned contract method is called normally");

            // migrate test
            CreateContract migrateCreateContract = CreateContract.deploy(web3j, transactionManager, provider, chainId).send();
            migrateCreateContract.set_simple_address(deploySimple).send();
            receipt = migrateCreateContract.migrate(deploySimple).send();
            collector.assertTrue(receipt.isStatusOK());
            WasmAddress migrateSimpleAddress = deploySimpleContract.get_address().send();
            collector.logStepPass("migrate contract successfully, contract address:" + migrateSimpleAddress.toString());

            // test migrate simple contract
            SimpleContract migrateSimpleContract = SimpleContract.load(migrateSimpleAddress.toString(), web3j, transactionManager, provider, chainId);
            migrateSimpleContract.set(Uint64.of(7)).send();
            Uint64 migrateSimpleGet = migrateSimpleContract.get().send();
            collector.assertEqual(migrateSimpleGet, Uint64.of(7));
            collector.logStepPass("The migrate contract method is called normally");

            // migrate clone test
            CreateContract migrateCloneCreateContract = CreateContract.deploy(web3j, transactionManager, provider, chainId).send();
            migrateCloneCreateContract.set_simple_address(cloneSimple).send();
            receipt = migrateCloneCreateContract.migrate_clone(cloneSimple).send();
            collector.assertTrue(receipt.isStatusOK());
            WasmAddress migrateCloneSimpleAddress = cloneSimpleContract.get_address().send();
            collector.logStepPass("migrate clone contract successfully, contract address:" + migrateCloneSimpleAddress.toString());

            // test migrate clone simple contract
            SimpleContract migrateCloneSimpleContract = SimpleContract.load(migrateCloneSimpleAddress.toString(), web3j, transactionManager, provider, chainId);
            migrateCloneSimpleContract.set(Uint64.of(7)).send();
            Uint64 migrateCloneSimpleGet = migrateCloneSimpleContract.get().send();
            collector.assertEqual(migrateCloneSimpleGet, Uint64.of(7));
            collector.logStepPass("The migrate clone contract method is called normally");

            // Failed to migrate contract test
            CreateContract failMigrateCreateContract = CreateContract.deploy(web3j, transactionManager, provider, chainId).send();
            WasmAddress failMigrateCreateContractAddress = new WasmAddress(createContract.getContractAddress());
            WasmAddress failMigrateDeployAddress;
            CreateContract failMigrateDeployContract;
            for(int i = 0; i < 20; i++){
                // create a new create contract
                receipt = failMigrateCreateContract.deploy_contract(wasmCreateAddress).send();
                collector.assertTrue(receipt.isStatusOK());
                failMigrateDeployAddress = failMigrateCreateContract.get_deploy_address().send();
                collector.logStepPass("create contract create a new contratc successfully, contract address:" + failMigrateDeployAddress.toString());

                // new create contract clone migrate
                failMigrateDeployContract = CreateContract.load(failMigrateDeployAddress.toString(), web3j, transactionManager, provider, chainId);
                receipt = failMigrateDeployContract.migrate_clone(failMigrateCreateContractAddress).send();
                collector.assertTrue(receipt.isStatusOK());
                collector.logStepPass("deploy create contract clone successfully");
            }

        } catch (Exception e) {
            collector.logStepFail("contract_VIDToken failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }


}
