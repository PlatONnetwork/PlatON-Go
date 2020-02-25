package wasm.data_type;

import com.platon.rlp.datatypes.Uint8;
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
            Boolean boolValue = true;
            TransactionReceipt  transactionReceipt = basicDataTypeContract.set_bool(boolValue).send();
            collector.logStepPass("basicDataTypeContract 【验证布尔值赋值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //Boolean expectValue = true;
            Boolean actualValue = basicDataTypeContract.get_bool().send();
            collector.logStepPass("basicDataTypeContract 【验证布尔值取值】 执行getBool() successfully actualValue:" + actualValue);
            //collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
            //2、验证：字节类型（byte）
           Uint8 uint8Value = Uint8.of("2");
            TransactionReceipt transactionReceipt1 = basicDataTypeContract.set_byte(uint8Value).send();
            collector.logStepPass("basicDataTypeContract 【验证字节类型（byte）】 successfully hash:" + transactionReceipt1.getTransactionHash());
            Uint8 actualByteValue = basicDataTypeContract.get_byte().send();
            collector.logStepPass("basicDataTypeContract 【验证字节类型取值】 执行get_byte() successfully actualByteValue:" + actualByteValue);
            //collector.assertEqual(actualValue,expectValue, "checkout  execute success.");*/
            //3、验证：字符串赋值
            String str = "ABC";
            TransactionReceipt transactionReceipt2 = basicDataTypeContract.set_string(str).send();
            collector.logStepPass("basicDataTypeContract 【验证字符串赋值】 successfully hash:" + transactionReceipt2.getTransactionHash());
            String actualStringValue = basicDataTypeContract.get_string().send();
            collector.logStepPass("basicDataTypeContract 【验证字符串取值】 执行get_string() successfully actualStringValue:" + actualStringValue);
            collector.assertEqual(actualStringValue,str, "checkout  execute success.");
              //验证：字符串长度
            Byte expectStringLength = 3;
            Uint8 actualStringLength = basicDataTypeContract.get_string_length().send();
            collector.logStepPass("basicDataTypeContract 【验证字符串长度】 执行get_string_length() successfully actualStringLength:" + actualStringLength);
            //collector.assertEqual(actualStringLength,expectStringLength, "checkout  execute success.");

            //4、地址类型(Address)
             TransactionReceipt  transactionReceipt4 = basicDataTypeContract.set_address().send();
            collector.logStepPass("basicDataTypeContract 【验证地址类型(Address)】 successfully hash:" + transactionReceipt4.getTransactionHash());
            String actualAddreeValue = basicDataTypeContract.get_address().send();
            collector.logStepPass("basicDataTypeContract 【验证地址取值】 执行getString() successfully actualAddreeValue:" + actualAddreeValue);
            //collector.assertEqual(actualStringValue,expectStringValue, "checkout  execute success.");*/

            //6、浮点类型(float、double)
           /* TransactionReceipt  transactionReceipt3 = basicDataTypeContract.setFloat().send();
            collector.logStepPass("basicDataTypeContract 【验证浮点类型(float、double)】 successfully hash:" + transactionReceipt3.getTransactionHash());*/


        } catch (Exception e) {
            collector.logStepFail("basicDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
