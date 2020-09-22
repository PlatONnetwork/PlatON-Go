package network.platon.test.wasm.contract_create;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithStruct;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约init带一个入参结构体测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class InitWithStructTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init带一个入参",sourcePrefix = "wasm")
    public void testNewContract() {

        String name = "myName";
        Uint64 age = Uint64.of(12L);
        try {
            prepare();
            InitWithStruct initWithStruct = InitWithStruct.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = initWithStruct.getContractAddress();
            String transactionHash = initWithStruct.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithStructTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initWithStruct.getTransactionReceipt().get().getGasUsed());

            transactionHash =initWithStruct.add_vector(name,age).send().getTransactionHash();
            collector.logStepPass("InitWithStructTest invoke successfully  hash:" + transactionHash);

            Uint8 idx = Uint8.of(0);
            String chainName = initWithStruct.get_vector(idx).send();
            collector.assertEqual(chainName,name);

        } catch (Exception e) {
            collector.logStepFail("InitWithStructTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
