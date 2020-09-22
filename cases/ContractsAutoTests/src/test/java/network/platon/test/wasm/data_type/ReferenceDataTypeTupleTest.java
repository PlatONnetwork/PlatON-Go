package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeTupleContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型（Tuple）
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeTupleTest extends WASMContractPrepareTest {

    private String expectValue1;
    private String expectValue2;
    private String a;
    private String b;
    private String name;
    private String age;


    @Before
    public void before() {
        expectValue1 = driverService.param.get("expectValue1");
        expectValue2 = driverService.param.get("expectValue2");
        a = driverService.param.get("a");
        b = driverService.param.get("b");
        name = driverService.param.get("name");
        age = driverService.param.get("age");
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeTupleTest定义元组赋值&取值",sourcePrefix = "wasm")
    public void testReferenceDataTypeTuple() {

         //部署合约
        ReferenceDataTypeTupleContract referenceDataTypeTupleContract = null;
        try {
            prepare();
            referenceDataTypeTupleContract = ReferenceDataTypeTupleContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeTupleContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeTupleContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeTupleContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeTupleContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：tuple元组初始化赋值方式一
            TransactionReceipt  transactionReceipt = referenceDataTypeTupleContract.setInitTupleModeOne().send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证元组初始化赋值方式一】 successfully hash:" + transactionReceipt.getTransactionHash());
           //2、验证：tuple根据索引取值
            String actualValue1 = referenceDataTypeTupleContract.getTupleValueIndex1().send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple元组根据索引取值】 执行getTupleValueIndex1() successfully actualValue1:" + actualValue1);
            collector.assertEqual(actualValue1,expectValue1, "checkout execute success.");
            Uint8 actualValue2 = referenceDataTypeTupleContract.getTupleValueIndex2().send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple元组根据索引取值】 执行getTupleValueIndex2() successfully actualValue2:" + actualValue2);
            collector.assertEqual(actualValue2,Uint8.of(expectValue2), "checkout execute success.");

            //3、验证：tuple元组初始化赋值方式二(使用make_tuple函数)
            TransactionReceipt  transactionReceipt1 = referenceDataTypeTupleContract.setInitTupleModeTwo(a,Uint8.of(b)).send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple元组初始化赋值方式二(使用make_tuple函数)】 successfully hash:" + transactionReceipt1.getTransactionHash());
            String actualValue3 = referenceDataTypeTupleContract.getTupleValueIndex3().send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple元组根据索引取值】 执行getTupleValueIndex() successfully actualValue3:" + actualValue3);
            collector.assertEqual(actualValue3,a, "checkout execute success.");

            //4、验证:定义包含引用类型
            TransactionReceipt  transactionReceipt2 = referenceDataTypeTupleContract.setInitTupleModeThree(name,Uint64.of(age)).send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证元组定义包含引用类型】 successfully hash:" + transactionReceipt2.getTransactionHash());
            ReferenceDataTypeTupleContract.Person person = referenceDataTypeTupleContract.getTupleValueIndex4().send();
            String actualValueName = person.name;
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple元组根据索引取值】 执行getTupleValueIndex4() successfully actualValueName:" + actualValueName);
            collector.assertEqual(actualValueName,name, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeTupleContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
