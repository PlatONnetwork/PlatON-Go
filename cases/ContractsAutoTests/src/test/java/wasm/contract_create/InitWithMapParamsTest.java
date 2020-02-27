package wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithMapParams;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

import java.util.HashMap;
import java.util.Map;

/**
 * @title 创建合约init入参包含map
 * @description:
 * @author: hudenian
 * @create: 2020/02/26
 */
public class InitWithMapParamsTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init入参包含map",sourcePrefix = "wasm")
    public void testMapParams() {

        Map<String,String> maps = new HashMap<String,String>();
        maps.put("key1","value1");
        maps.put("key2","value2");
        try {
            prepare();
            InitWithMapParams initWithMapParams = InitWithMapParams.deploy(web3j, transactionManager, provider,maps).send();
            String contractAddress = initWithMapParams.getContractAddress();
            String transactionHash = initWithMapParams.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithMapParamsTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //查询包含map对象
            Map chainMap = initWithMapParams.get_map().send();

            collector.assertEqual(maps.get("key1").toString(),chainMap.get("key1").toString());
            collector.assertEqual(maps.get("key2").toString(),chainMap.get("key2").toString());

        } catch (Exception e) {
            collector.logStepFail("InitWithMapParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
