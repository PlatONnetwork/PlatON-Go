package wasm.contract_migrate;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.utils.FileUtil;
import network.platon.contracts.wasm.ContractMigrate_v1;
import network.platon.utils.RlpUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.io.File;
import java.nio.file.Paths;
import java.util.ArrayList;

/**
 * @title 合约升级
 * @description:
 * @author: hudenian
 * @create: 2020/02/10
 */
public class ContractMigrateV1Test extends WASMContractPrepareTest {

    //需要升级的合约
    private String wasmFile = "ContractMigrate_v1.bin";

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_migrate合约升级",sourcePrefix = "wasm")
    public void testMigrateContract() {

        Byte[] init_arg = null;
        Long transfer_value = 100000L;
        Long gas_value = 200000L;
        String name = "hello";

        try {
            prepare();
            ContractMigrate_v1 contractMigratev1 = ContractMigrate_v1.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contractMigratev1.getContractAddress();
            String transactionHash = contractMigratev1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractMigrateV1 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //设置值
            transactionHash = contractMigratev1.set_string(name).send().getTransactionHash();
            collector.logStepPass("ContractMigrateV1 set_string successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //查询结果
            String chainName = contractMigratev1.get_string().send();
            collector.assertEqual(chainName,name);


            /**
             * 加载需要升级的合约
             * init_arg 参数为  magic number +  RLP(code, RLP("init", init_paras...))
             * transfer_value 为转到新合约地址的金额，gas_value 为预估消耗的 gas
             */
            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "wasm", "contract_migrate").toUri().getPath());
            init_arg = RlpUtil.loadInitArg(filePath+File.separator+wasmFile,new ArrayList<String>());

            //合约升级
            TransactionReceipt transactionReceipt = contractMigratev1.migrate_contract(init_arg,transfer_value,gas_value).send();
            collector.logStepPass("Contract Migrate V1  successfully hash:" + transactionReceipt.getTransactionHash());

            //获取升级后的合约地址(需要通过事件获取)
            String newContractAddress = contractMigratev1.getPlaton_event1_transferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);

            //调用升级后的合约
            //FIXME 等bug修复后放开
//            ContractMigrate_v1 new_contractMigrate_v1 = ContractMigrate_v1.load(newContractAddress,web3j,credentials,provider);
//            String newContractChainName = new_contractMigrate_v1.get_string().send();
//            collector.assertContains(newContractChainName,name);

        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
