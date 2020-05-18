package wasm.contract_create;

import com.platon.rlp.datatypes.Uint;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithListParams;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.gas.ContractGasProvider;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.*;

/**
 * @title 创建合约init入参包含list
 * @description:
 * @author: hudenian
 * @create: 2020/02/27
 */
public class InitWithListParamsTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init入参包含list及嵌套测试",sourcePrefix = "wasm")
    public void testMapParams() {

        //单个list
        List<String> list = new ArrayList<String>();
        String list1 = "list1";
        list.add(list1);

        //list中要添加的元素
        String list2 = "list2";
        String list3 = "list3";
        String list4 = "list4";
        String list5 = "list5";
        String list6 = "list6";

        List<String> listMerge = Arrays.asList(new String[]{"list7","list8"});

                //list嵌套list
        List<List<String>> listlist = new ArrayList<List<String>>();
        listlist.add(list);

        try {
            prepare();
            InitWithListParams initWithListParams = InitWithListParams.deploy(web3j, transactionManager, provider, chainId,list).send();
            String contractAddress = initWithListParams.getContractAddress();
            String transactionHash = initWithListParams.getTransactionReceipt().get().getTransactionHash();
            String gasUsed =  initWithListParams.getTransactionReceipt().get().getGasUsed().toString();
            collector.logStepPass("InitWithListParamsTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash+",deploy gas used:"+gasUsed);
            collector.logStepPass("deploy gas used:" + initWithListParams.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt tx = initWithListParams.set_list(list).send();
            collector.logStepPass("InitWithListParamsTest call set_list successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查询包含list对象
            List<String> chainList = initWithListParams.get_list().send();

            collector.assertEqual(list.get(0).toString(),chainList.get(0).toString());


            //list
            tx = initWithListParams.set_list_list(listlist).send();
            collector.logStepPass("InitWithListParamsTest call add_map_map successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //调用list嵌套list
            List<List<String>> chainlistlist = initWithListParams.get_list_list().send();
            collector.logStepPass("InitWithMapParamsTest call add_map_list successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.assertEqual(listlist.get(0).get(0).toString(), chainlistlist.get(0).get(0).toString());

            //list中添加一个元素
            tx = initWithListParams.add_list_element(list2).send();
            collector.logStepPass("InitWithListParamsTest call add_list_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //取出list中最后一个元素
            String chainLastElement = initWithListParams.get_list_last_element().send();
            collector.assertEqual(list2,chainLastElement);

            //取出list中第一个元素
            String chainFirstElement = initWithListParams.get_list_first_element().send();
            collector.assertEqual(list1,chainFirstElement);

            //list中添加一个元素
            tx = initWithListParams.add_list_element(list3).send();
            collector.logStepPass("InitWithListParamsTest call add_list_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //list中删除一个元素list2(还有list1 与list3)
            tx = initWithListParams.list_remove_element(list2).send();
            collector.logStepPass("InitWithListParamsTest call list_remove_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //list中删除第一个元素(还剩下list3)
            tx = initWithListParams.list_pop_front().send();
            collector.logStepPass("InitWithListParamsTest call list_remove_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查询包含list对象
            chainList = initWithListParams.get_list().send();
            collector.assertEqual(list3,chainList.get(0).toString());

            //list 中再添加list3
            tx = initWithListParams.add_list_element(list3).send();
            collector.logStepPass("InitWithListParamsTest call add_list_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看list的大小
            Uint listSize =initWithListParams.list_size().send();
            collector.assertEqual("2",listSize.value.toString());

            //去掉list中的重复元素
            tx = initWithListParams.list_unique().send();
            collector.logStepPass("InitWithListParamsTest call list_unique successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看list的大小(理论上只有一个元素)
            listSize =initWithListParams.list_size().send();
            collector.assertEqual("1",listSize.value.toString());

            //list中添加另一个list
            tx = initWithListParams.list_merge(listMerge).send();
            collector.logStepPass("InitWithListParamsTest call list_merge successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看list的大小(理论上只有一个元素)
            listSize =initWithListParams.list_size().send();
            collector.assertEqual("3",listSize.value.toString());

        } catch (Exception e) {
            collector.logStepFail("InitWithListParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
