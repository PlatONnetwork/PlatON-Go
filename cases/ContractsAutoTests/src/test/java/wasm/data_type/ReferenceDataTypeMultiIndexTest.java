package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeArrayContract;
import network.platon.contracts.wasm.ReferenceDataTypeMultiIndexContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型多索引
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeMultiIndexTest extends WASMContractPrepareTest {

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
            author = "qudong", showName = "wasm.referenceDataTypeMultiIndex测试引用类型多索引",sourcePrefix = "wasm")
    public void testReferenceDataTypeMultiIndex() {

         //部署合约
        ReferenceDataTypeMultiIndexContract referenceDataTypeMultiIndexContract = null;
        try {
            prepare();
            referenceDataTypeMultiIndexContract = ReferenceDataTypeMultiIndexContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = referenceDataTypeMultiIndexContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeMultiIndexContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeMultiIndexContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeMultiIndexContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：多索引插入数据
            TransactionReceipt  transactionReceipt = referenceDataTypeMultiIndexContract.addInitMultiIndex(my_name,Byte.parseByte(my_age),Byte.parseByte(my_sex)).send();
            collector.logStepPass("referenceDataTypeMultiIndexContract 【验证多索引插入数据】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：find（）多索引取值(查询年龄为10)
            Byte  actualValueAge =  referenceDataTypeMultiIndexContract.getMultiIndexFind(Byte.parseByte(my_age)).send();
            collector.logStepPass("referenceDataTypeMultiIndexContract 【验证find（）多索引取值(查询年龄为10)】 执行getMultiIndexFind() successfully actualValueAge:" + actualValueAge);
            collector.assertEqual(actualValueAge,Byte.parseByte(my_age), "checkout  execute success.");
           //3、验证：cbegin()多索引迭代器起始位置
            Boolean actualSame= referenceDataTypeMultiIndexContract.getMultiIndexCbegin().send();
            collector.logStepPass("referenceDataTypeMultiIndexContract 【验证cbegin()多索引迭代器起始位置】 执行getMultiIndexCbegin() successfully actualSame:" + actualSame);
            collector.assertEqual(actualSame,true, "checkout  execute success.");
            //4、验证：erase()多索引删除数据（删除年龄为10）
            TransactionReceipt  transactionReceipt1 = referenceDataTypeMultiIndexContract.deleteMultiIndexErase().send();
            collector.logStepPass("referenceDataTypeMultiIndexContract 【验证erase()多索引删除数据】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //查询删除数据
            Byte  actualValueAge1 =  referenceDataTypeMultiIndexContract.getMultiIndexFind(Byte.parseByte(my_age)).send();
            collector.logStepPass("referenceDataTypeMultiIndexContract 【验证find（）多索引取值(查询年龄为10)】 执行getMultiIndexFind() successfully actualValueAge1:" + actualValueAge1);
            collector.assertEqual(actualValueAge1,Byte.parseByte(delete_age), "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeMultiIndexContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
