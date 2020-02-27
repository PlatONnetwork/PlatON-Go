package wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithListParams;
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
        list.add("list1");

        //list嵌套list
        List<List<String>> listlist = new ArrayList<List<String>>();
        listlist.add(list);

        try {
            prepare();
            provider = new ContractGasProvider(BigInteger.valueOf(50000000004L), BigInteger.valueOf(90000000L));
            InitWithListParams initWithListParams = InitWithListParams.deploy(web3j, transactionManager, provider,list).send();
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

        } catch (Exception e) {
            collector.logStepFail("InitWithListParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
