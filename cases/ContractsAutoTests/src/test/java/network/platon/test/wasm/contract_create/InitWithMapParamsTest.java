package network.platon.test.wasm.contract_create;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithMapParams;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

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

        //map 需要添加的key与value
        String key2 = "key2";
        String value2 = "value2";

        String key3 = "key3";
        String value3 = "value3";

        String key4 = "key4";
        String value4 = "value4";

        //需要被添加进map的map对象
        String key5 = "key5";
        String value5 = "value5";

        String key6 = "key6";
        String value6 = "value6";

        Map<String,String> maps2 = new HashMap<String,String>();
        maps2.put(key5,value5);
        maps2.put(key6,value6);


        try {
            prepare();
            InitWithMapParams initWithMapParams = InitWithMapParams.deploy(web3j, transactionManager, provider, chainId,maps).send();
            String contractAddress = initWithMapParams.getContractAddress();
            String transactionHash = initWithMapParams.getTransactionReceipt().get().getTransactionHash();
            String gasUsed =  initWithMapParams.getTransactionReceipt().get().getGasUsed().toString();
            collector.logStepPass("InitWithMapParamsTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash+",deploy gas used:"+gasUsed);
            collector.logStepPass("deploy gas used:" + initWithMapParams.getTransactionReceipt().get().getGasUsed());

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

            //调用map嵌套list
            tx = initWithMapParams.add_map_list(inMaplist).send();
            collector.logStepPass("InitWithMapParamsTest call add_map_list successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            Map<String, List<String>> maplist = initWithMapParams.get_map_list().send();
            collector.assertEqual(list.get(0), maplist.get("keyList").get(0));

            //map中添加键值对
            tx = initWithMapParams.add_map_element(key3,value3).send();
            collector.logStepPass("InitWithMapParamsTest call add_map_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            chainMap = initWithMapParams.get_map().send();

            collector.assertEqual(value3,chainMap.get("key3").toString());

            //map中删除指定的key值
            tx = initWithMapParams.delete_map_element(key3).send();
            collector.logStepPass("InitWithMapParamsTest call delete_map_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            chainMap = initWithMapParams.get_map().send();

            collector.assertEqual(2,chainMap.size());

            //map中查找不存在的key
            String chainValue = initWithMapParams.find_element_bykey(key4).send();
            collector.assertEqual("",chainValue);

            //map中查找存在的key
            chainValue = initWithMapParams.find_element_bykey(key2).send();
            collector.assertEqual(value2,chainValue);

            //map 查看大小
            Uint8 mapSize = initWithMapParams.get_map_size().send();
            collector.logStepPass("当前map中元素个数为："+mapSize.value);

            //map中添加另一个map对象
            tx = initWithMapParams.add_map(maps2).send();
            collector.logStepPass("InitWithMapParamsTest call add_map successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看新的map是否添加成功
            String chainValue6 = initWithMapParams.find_element_bykey(key6).send();
            collector.assertEqual(value6,chainValue6);

        } catch (Exception e) {
            collector.logStepFail("InitWithMapParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
