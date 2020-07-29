package wasm.data_type;

import com.platon.rlp.datatypes.Uint32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeArrayContract;
import network.platon.contracts.wasm.ReferenceDataTypeArrayFuncContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型数组(array类型)属性/函数
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeArrayFuncTest extends WASMContractPrepareTest {

    private String indexValue;
    private String expectIndexValue;
    private String expectFirstValue;
    private String fillValue;
    private String expectFill;

    @Before
    public void before() {
        indexValue = driverService.param.get("indexValue");
        expectIndexValue = driverService.param.get("expectIndexValue");
        expectFirstValue = driverService.param.get("expectFirstValue");
        fillValue = driverService.param.get("fillValue");
        expectFill = driverService.param.get("expectFill");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeArrayTest验证数组属性及函数",sourcePrefix = "wasm")
    public void testReferenceDataTypeArrayFunc() {

         //部署合约
        ReferenceDataTypeArrayFuncContract referenceDataTypeArrayFuncContract = null;
        try {
            prepare();
            referenceDataTypeArrayFuncContract = ReferenceDataTypeArrayFuncContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeArrayFuncContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeArrayFuncContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeArrayFuncContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeArrayFuncContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、数组初始化数据源
            TransactionReceipt  transactionReceipt = referenceDataTypeArrayFuncContract.setInitArrayDate().send();
            collector.logStepPass("referenceDataTypeArrayFuncContract 【数组初始化数据源】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：数组属性empty()
            Boolean actualValue = referenceDataTypeArrayFuncContract.getArrayIsEmpty().send();
            collector.logStepPass("referenceDataTypeArrayFuncContract 【验证：数组属性empty()】 执行getArrayIsEmpty() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,false, "checkout  execute success.");
            //3、验证：数组属性at()
            String actualStringValue = referenceDataTypeArrayFuncContract.getArrayValueIndex(Uint32.of(indexValue)).send();
            collector.logStepPass("referenceDataTypeArrayFuncContract 【验证：数组属性at()】 执行getArrayValueIndex() successfully actualStringValue:" + actualStringValue);
            collector.assertEqual(actualStringValue,expectIndexValue, "checkout  execute success.");
            //4、验证：数组属性front()
            String actualStringValue1 = referenceDataTypeArrayFuncContract.getArrayFirstValue().send();
            collector.logStepPass("referenceDataTypeArrayFuncContract 【验证：数组属性front()】 执行getArrayFirstValue() successfully actualStringValue1:" + actualStringValue1);
            collector.assertEqual(actualStringValue1,expectFirstValue, "checkout  execute success.");
            //5、验证：数组属性fill()
            TransactionReceipt  transactionReceipt1 = referenceDataTypeArrayFuncContract.setArrayFill(fillValue).send();
            collector.logStepPass("referenceDataTypeArrayFuncContract 【验证：数组属性fill()】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //取值验证
            String actualStringValue2 = referenceDataTypeArrayFuncContract.getArrayValueIndex(Uint32.of(expectFill)).send();
            collector.logStepPass("referenceDataTypeArrayFuncContract 【验证：数组属性at()】 执行getArrayValueIndex() successfully actualStringValue2:" + actualStringValue2);
            collector.assertEqual(actualStringValue2,fillValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeArrayFuncContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
