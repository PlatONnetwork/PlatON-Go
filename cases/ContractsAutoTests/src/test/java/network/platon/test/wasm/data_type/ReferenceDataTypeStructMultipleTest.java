package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint32;
import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeStructMultipleContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试结构体包含复杂类型赋值&取值
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeStructMultipleTest extends WASMContractPrepareTest {

    private String myGroupName;
    private String myGroupId;
    private String oneValue;
    private String twoValue;
    private String arrayIndex;


    @Before
    public void before() {
        myGroupName = driverService.param.get("myGroupName");
        myGroupId = driverService.param.get("myGroupId");
        oneValue = driverService.param.get("oneValue");
        twoValue = driverService.param.get("twoValue");
        arrayIndex = driverService.param.get("arrayIndex");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeStructTest验证结构体包含复杂类型赋值&取值",sourcePrefix = "wasm")
    public void testReferenceDataTypeStructMultiple() {

         //部署合约
        ReferenceDataTypeStructMultipleContract referenceDataTypeStructMultipleContract = null;
        try {
            prepare();
            referenceDataTypeStructMultipleContract = ReferenceDataTypeStructMultipleContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeStructMultipleContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeStructMultipleContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeStructMultipleContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructMultipleContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {

            //1、包含(引用类型及基本类型)struct类型赋值
            TransactionReceipt  transactionReceipt = referenceDataTypeStructMultipleContract.setGroupValue(myGroupName,Uint64.of(myGroupId)).send();
            collector.logStepPass("referenceDataTypeStructMultipleContract 【验证struct类型包含引用类型及基本类型赋值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：验证struct结构体 groupName取值
            String actualGroupName = referenceDataTypeStructMultipleContract.getGroupName().send();
            collector.logStepPass("referenceDataTypeStructMultipleContract 【验证struct结构体 groupName取值】 执行getGroupName() successfully actualGroupName:" + actualGroupName);
            collector.assertEqual(actualGroupName,myGroupName, "checkout  execute success.");

            //3、包含(引用类型及基本类型)struct类型赋值
            TransactionReceipt  transactionReceipt1 = referenceDataTypeStructMultipleContract.setGroupArrayValue(oneValue,twoValue).send();
            collector.logStepPass("referenceDataTypeStructMultipleContract 【验证struct类型数组赋值】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //4、验证：验证struct结构体 getGroupArrayIndexValue取值
            String actualArrayValue = referenceDataTypeStructMultipleContract.getGroupArrayIndexValue(Uint32.of(arrayIndex)).send();
            collector.logStepPass("referenceDataTypeStructMultipleContract 【验证struct结构体数组取值】 执行getGroupArrayIndexValue() successfully actualArrayValue:" + actualArrayValue);
             collector.assertEqual(actualArrayValue,twoValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructMultipleContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
