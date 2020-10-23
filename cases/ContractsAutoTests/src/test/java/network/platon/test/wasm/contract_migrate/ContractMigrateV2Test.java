package network.platon.test.wasm.contract_migrate;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractMigrate_v1;
import network.platon.contracts.wasm.ContractMigrate_v2;
import org.junit.Test;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.utils.Numeric;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.util.Arrays;

/**
 * @title 合约升级增加升级权限校验
 * @description:
 * @author: hudenian
 * @create: 2020/02/16
 */
public class ContractMigrateV2Test extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_migrateV2合约升级带权限验证",sourcePrefix = "wasm")
    public void testMigrateContract() {

        String name = "hello";

        try {
            prepare();
            ContractMigrate_v2 contractMigratev2 = ContractMigrate_v2.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractMigratev2.getContractAddress();
            String transactionHash = contractMigratev2.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contractMigratev2 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractMigratev2.getTransactionReceipt().get().getGasUsed());

            //设置值
            transactionHash = contractMigratev2.set_string(name).send().getTransactionHash();
            collector.logStepPass("contractMigratev2 set_string successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //查询结果
            String chainName = contractMigratev2.get_string().send();
            collector.assertEqual(chainName,name);


            String code = WasmFunctionEncoder.encodeConstructor(contractMigratev2.getContractBinary(), Arrays.asList());
            byte[] data = Numeric.hexStringToByteArray(code);

            //合约升级
            TransactionReceipt transactionReceipt = contractMigratev2.migrate_contract(data, Uint64.of(0L), Uint64.of(90000000L)).send();
            collector.logStepPass("Contract Migrate V1  successfully hash:" + transactionReceipt.getTransactionHash());

            //获取升级后的合约地址(需要通过事件获取)
            String newContractAddress = contractMigratev2.getTransferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);

            //调用升级后的合约
            ContractMigrate_v1 new_contractMigrate_v1 = ContractMigrate_v1.load(newContractAddress,web3j,credentials,provider, chainId);
            String newContractChainName = new_contractMigrate_v1.get_string().send();
            collector.assertContains(newContractChainName,name);

            //调用旧的合约(旧合约已被销毁不能再调用)
            try{
                String chainName2 =  contractMigratev2.get_string().send();
            }catch (Exception e){
                collector.logStepPass("old contract can not migrate");
            }

        } catch (Exception e) {
            collector.logStepFail("ContractMigrateV2Test failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
