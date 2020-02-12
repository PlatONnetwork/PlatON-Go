package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeTupleContract;
import network.platon.contracts.wasm.ReferenceDataTypeVectorContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型（Tuple）
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeTupleTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeTupleTest定义元组赋值&取值",sourcePrefix = "wasm")
    public void testReferenceDataTypeTuple() {

         //部署合约
        ReferenceDataTypeTupleContract referenceDataTypeTupleContract = null;
        try {
            prepare();
            referenceDataTypeTupleContract = ReferenceDataTypeTupleContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = referenceDataTypeTupleContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeTupleContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeTupleContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeTupleContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：tuple类型初始化赋值
            TransactionReceipt  transactionReceipt = referenceDataTypeTupleContract.setInitTuple().send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple类型初始化赋值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：tuple类型通过make_tuple()生成tuple对象
            TransactionReceipt  transactionReceipt1 = referenceDataTypeTupleContract.setTupleObject().send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple类型通过make_tuple()生成tuple对象】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //3、验证：tuple根据索引取值
            String expectValue = "test";
            String actualValue = referenceDataTypeTupleContract.getTupleValueIndex().send();
            collector.logStepPass("referenceDataTypeTupleContract 【验证tuple根据索引取值】 执行getTupleValueIndex() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeTupleContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
