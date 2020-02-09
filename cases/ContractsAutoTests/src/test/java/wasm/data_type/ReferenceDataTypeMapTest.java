package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeArrayContract;
import network.platon.contracts.wasm.ReferenceDataTypeMapContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型Map集合
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeMapTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeMapTest",sourcePrefix = "wasm")
    public void testReferenceDataTypeMap() {

         //部署合约
        ReferenceDataTypeMapContract referenceDataTypeMapContract = null;
        try {
            prepare();
            referenceDataTypeMapContract = ReferenceDataTypeMapContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = referenceDataTypeMapContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeMapContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeMapContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeMapContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：map中的key与value可以是任意类型
            TransactionReceipt  transactionReceipt = referenceDataTypeMapContract.setMap().send();
            collector.logStepPass("referenceDataTypeMapContract 【验证map中的key与value可以是任意类型】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：map属性方法新增
            TransactionReceipt  transactionReceipt1 = referenceDataTypeMapContract.addMap().send();
            collector.logStepPass("referenceDataTypeMapContract 【验证map属性方法新增】 执行addMap() successfully hash:" + transactionReceipt1.getTransactionHash());
            //3、验证：获取数组容器大小
            Byte expectMapLength = 3;
            Byte actualMapLength = referenceDataTypeMapContract.getMapSize().send();
            collector.logStepPass("referenceDataTypeMapContract 【验证获取Map容器大小】 执行getMapSize() successfully actualMapLength:" + actualMapLength);
            collector.assertEqual(actualMapLength,expectMapLength, "checkout  execute success.");
            //4、验证：key关键字只能在map出现一次
            TransactionReceipt transactionReceipt2 = referenceDataTypeMapContract.setSameKeyMap().send();
            collector.logStepPass("referenceDataTypeMapContract 【验证Map中key值唯一性】 successfully hash:" + transactionReceipt2.getTransactionHash());


           /* Byte expectByteValue = 1;
            Byte actualByteValue = referenceDataTypeArrayContract.getBytesArrayIndex().send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证字节数组取值】 执行getBytesArrayIndex() successfully actualByteValue:" + actualByteValue);
            collector.assertEqual(actualByteValue,expectByteValue, "checkout  execute success.");*/
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeArrayContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
