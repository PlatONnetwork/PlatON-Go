package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.BasicDataTypeContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试基本类型
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class BasicDataTypeTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.basicDataTypeTest基本类型验证测试",sourcePrefix = "wasm")
    public void testBasicDataType() {

         //部署合约
        BasicDataTypeContract basicDataTypeContract = null;
        try {
            prepare();
            basicDataTypeContract = BasicDataTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = basicDataTypeContract.getContractAddress();
            TransactionReceipt tx = basicDataTypeContract.getTransactionReceipt().get();
            collector.logStepPass("basicDataTypeContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("basicDataTypeContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证:布尔值赋值
            TransactionReceipt  transactionReceipt = basicDataTypeContract.setBool().send();
            collector.logStepPass("basicDataTypeContract 【验证布尔值赋值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：布尔值取值
            Boolean expectValue = true;
            Boolean actualValue = basicDataTypeContract.getBool().send();
            collector.logStepPass("basicDataTypeContract 【验证布尔值取值】 执行getBool() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
            //3、验证：字节类型（byte）
            TransactionReceipt transactionReceipt1 = basicDataTypeContract.setByte().send();
            collector.logStepPass("basicDataTypeContract 【验证字节类型（byte）】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //4、验证：字符串赋值
            String str = "ABC";
            TransactionReceipt transactionReceipt2 = basicDataTypeContract.setString(str).send();
            collector.logStepPass("basicDataTypeContract 【验证字符串赋值】 successfully hash:" + transactionReceipt2.getTransactionHash());
            //5、验证：字符串取值
            String actualStringValue = basicDataTypeContract.getString().send();
            collector.logStepPass("basicDataTypeContract 【验证字符串取值】 执行getString() successfully actualByteValue:" + actualStringValue);
            collector.assertEqual(actualStringValue,str, "checkout  execute success.");
            //5、验证：字符串长度
            Byte expectStringLength = 3;
            Byte actualStringLength = basicDataTypeContract.getStringLength().send();
            collector.logStepPass("basicDataTypeContract 【验证字符串长度】 执行getStringLength() successfully actualStringLength:" + actualStringLength);
            collector.assertEqual(actualStringLength,expectStringLength, "checkout  execute success.");
            //6、浮点类型(float、double)
           /* TransactionReceipt  transactionReceipt3 = basicDataTypeContract.setFloat().send();
            collector.logStepPass("basicDataTypeContract 【验证浮点类型(float、double)】 successfully hash:" + transactionReceipt3.getTransactionHash());*/
            //7、地址类型(Address)
            TransactionReceipt  transactionReceipt4 = basicDataTypeContract.setContractCallAddress().send();
            collector.logStepPass("basicDataTypeContract 【验证地址类型(Address)】 successfully hash:" + transactionReceipt4.getTransactionHash());
            //8、地址取值
            String actualAddreeValue = basicDataTypeContract.getContractCallAddress().send();
            collector.logStepPass("basicDataTypeContract 【验证地址取值】 执行getString() successfully actualAddreeValue:" + actualAddreeValue);
            //collector.assertEqual(actualStringValue,expectStringValue, "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("basicDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
