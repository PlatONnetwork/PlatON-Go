package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MultiIndexContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 多索引容器
 * @description:
 * @author: liweic
 * @create: 2020/04/20
 */
public class MultiIndexContractTest extends WASMContractPrepareTest {

    private String my_name;
    private String my_age;
    private String my_sex;
    private String delete_age;


    @Before
    public void before() {
        my_name = driverService.param.get("my_name");
        my_age = driverService.param.get("my_age");
        my_sex = driverService.param.get("my_sex");
        delete_age = driverService.param.get("delete_age");
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.MultiIndex测试多索引容器",sourcePrefix = "wasm")
    public void test() {

        try {
            prepare();
            MultiIndexContract multiindexcontract = MultiIndexContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = multiindexcontract.getContractAddress();
            TransactionReceipt tx = multiindexcontract.getTransactionReceipt().get();
            collector.logStepPass("MultiIndexContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

            //1、验证：多索引插入数据
            TransactionReceipt  transactionReceipt = multiindexcontract.addInitMultiIndex(my_name,Uint8.of(my_age),Uint8.of(my_sex)).send();
            collector.logStepPass("MultiIndexContract验证多索引插入数据 successfully hash:" + transactionReceipt.getTransactionHash());

            //2、验证：cbegin()多索引迭代器起始位置
            Boolean actualbegin= multiindexcontract.getMultiIndexCbegin("lucy").send();
            collector.logStepPass("MultiIndexContract 验证cbegin()多索引迭代器起始位置 执行getMultiIndexCbegin() successfully actualbegin:" + actualbegin);
            collector.assertEqual(actualbegin,true, "checkout  execute success.");

            //3、验证：cend()多索引迭代器结束位置
            Boolean actualend= multiindexcontract.getMultiIndexCend(Uint8.of(1)).send();
            collector.logStepPass("MultiIndexContract 验证cend()多索引迭代器结束位置 getMultiIndexCend() successfully actualend:" + actualend);
            collector.assertEqual(actualend,true, "checkout  execute success.");

            //4、验证：count获取与索引值对应的数据的数量
            Uint8 actualcount = multiindexcontract.getMultiIndexCount(Uint8.of(10)).send();
            collector.logStepPass("MultiIndexContract 验证索引对应数据的数量 getMultiIndexCount() successfully actualcount:" + actualcount);
            collector.assertEqual(actualcount,Uint8.of(1), "checkout  execute success.");

            //5、验证：find多索引取值
            String  actualfind = multiindexcontract.getMultiIndexFind("lucy").send();
            collector.logStepPass("MultiIndexContract 验证find多索引取值 执行getMultiIndexFind() successfully actualfind:" + actualfind);
            collector.assertEqual(actualfind,"lucy", "checkout  execute success.");

            //6、验证：get_index获取非唯一索引的索引对象
            Boolean actualindex = multiindexcontract.getMultiIndexIndex(Uint8.of(10)).send();
            collector.logStepPass("MultiIndexContract 验证get_index获取非唯一索引的索引对象 执行getMultiIndexFind() successfully actualindex:" + actualindex);
            collector.assertEqual(actualindex, true, "checkout  execute success.");

//            //7、验证：modify基于迭代器修改数据(待验证)
//            TransactionReceipt actualtx = multiindexcontract.MultiIndexModify("messi").send();
//            collector.logStepPass("MultiIndexContract验证多索引修改数据 successfully hash:" + actualtx.getTransactionHash());
//            String actualfind2 = multiindexcontract.getMultiIndexFind("lucy").send();
//            collector.logStepPass("MultiIndexContract 验证数据是否修改成功 successfully actualfind2:" + actualfind2);
//            collector.assertEqual(actualfind2,"lucy", "checkout  execute success.");

            //8、验证：erase()多索引删除数据
            TransactionReceipt erase = multiindexcontract.MultiIndexErase("lucy").send();
            collector.logStepPass("MultiIndexContract验证多索引修改数据 successfully hash:" + erase.getTransactionHash());
            String actualfind3 = multiindexcontract.getMultiIndexFind("lucy").send();
            collector.logStepPass("MultiIndexContract 验证数据是否删除成功 successfully actualfind2:" + actualfind3);
            collector.assertEqual(actualfind3,"", "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeMultiIndexContract deploy fail.", e.toString());
            e.printStackTrace();
        }
    }
}

