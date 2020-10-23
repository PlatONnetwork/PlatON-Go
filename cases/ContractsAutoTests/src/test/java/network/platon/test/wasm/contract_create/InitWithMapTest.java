package network.platon.test.wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithMap;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约init包含Map测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/16
 */
public class InitWithMapTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init带Map",sourcePrefix = "wasm")
    public void testNewContract() {

        String mapKey = "name";
        String mapValue = "Lily";

        try {
            prepare();
            InitWithMap initWithMap = InitWithMap.deploy(web3j, transactionManager, provider, chainId,mapKey,mapValue).send();
            String contractAddress = initWithMap.getContractAddress();
            String transactionHash = initWithMap.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithMap issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initWithMap.getTransactionReceipt().get().getGasUsed());

            Byte idx = 0;
            String chainMapValue = initWithMap.get_map(mapKey).send();
            collector.assertEqual(chainMapValue,mapValue);

        } catch (Exception e) {
            collector.logStepFail("InitWithMapTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
