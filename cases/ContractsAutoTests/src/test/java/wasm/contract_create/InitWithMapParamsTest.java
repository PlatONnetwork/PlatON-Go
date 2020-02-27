package wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithMapParams;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.gas.ContractGasProvider;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
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
            author = "hudenian", showName = "wasm.contract_create创建合约init入参包含map及嵌套测试",sourcePrefix = "wasm")
    public void testMapParams() {

        //单个map测试
        Map<String,String> maps = new HashMap<String,String>();
        maps.put("key1","value1");
        maps.put("key2","value2");

        //单个list
        List<String> list = new ArrayList<>();
        list.add("list1");

        //map嵌套map测试
        Map<String, Map<String, String>> inMapmap = new HashMap<String, Map<String, String>>();
        inMapmap.put("map1",maps);

        //map嵌套list测试
        Map<String, List<String>> inMaplist = new HashMap<String, List<String>>();
        inMaplist.put("keyList",list);

        try {
            prepare();
            provider = new ContractGasProvider(BigInteger.valueOf(50000000004L), BigInteger.valueOf(90000000L));
            InitWithMapParams initWithMapParams = InitWithMapParams.deploy(web3j, transactionManager, provider,maps).send();
            String contractAddress = initWithMapParams.getContractAddress();
            String transactionHash = initWithMapParams.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithMapParamsTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //查询包含map对象
            Map chainMap = initWithMapParams.get_map().send();

            collector.assertEqual(maps.get("key1").toString(),chainMap.get("key1").toString());
            collector.assertEqual(maps.get("key2").toString(),chainMap.get("key2").toString());

            //调用map嵌套map
            TransactionReceipt tx = initWithMapParams.add_map_map(inMapmap).send();
            collector.logStepPass("InitWithMapParamsTest call add_map_map successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            Map<String, Map<String, String>> mapmap = initWithMapParams.get_map_map().send();
            collector.assertEqual(maps.get("key1").toString(), mapmap.get("map1").get("key1").toString());
            collector.assertEqual(maps.get("key2").toString(), mapmap.get("map1").get("key2").toString());

            //调用map嵌套map
            tx = initWithMapParams.add_map_list(inMaplist).send();
            collector.logStepPass("InitWithMapParamsTest call add_map_list successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            Map<String, List<String>> maplist = initWithMapParams.get_map_list().send();
            collector.assertEqual(list.get(0), maplist.get("keyList").get(0));


        } catch (Exception e) {
            collector.logStepFail("InitWithMapParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
