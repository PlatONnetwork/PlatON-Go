package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.BasicDataIntegerTypeContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试整型基本类型
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class BasicDataIntegerTypeTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.basicDataTypeTest整型基本类型验证测试",sourcePrefix = "wasm")
    public void testBasicDataIntegerTypeTest() {

         //部署合约
        BasicDataIntegerTypeContract basicDataIntegerTypeContract = null;
        try {
            prepare();
            basicDataIntegerTypeContract = BasicDataIntegerTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = basicDataIntegerTypeContract.getContractAddress();
            TransactionReceipt tx = basicDataIntegerTypeContract.getTransactionReceipt().get();
            collector.logStepPass("basicDataIntegerTypeContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("basicDataIntegerTypeContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：整型有符号/无符号类型
            TransactionReceipt  transactionReceipt = basicDataIntegerTypeContract.setSignedInt().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整型有符号/无符号类型】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：整型uint8_t取值
            Byte expectValue = 3;
            Byte actualUint8Value = basicDataIntegerTypeContract.getStorageUint8().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数uint8取值】 执行getStorageUint8() successfully actualValue:" + actualUint8Value);
            collector.assertEqual(actualUint8Value,expectValue, "checkout  execute success.");
            //3、验证：整型int8_t取值
            Byte expectValue1 = 1;
            Byte actualInt8Value = basicDataIntegerTypeContract.getStorageInt8().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数int8_t取值】 执行getStorageInt8() successfully actualValue:" + actualInt8Value);
            //collector.assertEqual(actualInt8Value,expectValue1, "checkout  execute success.");
            //4、验证：整型int16_t取值
            Integer expectValue2 = -10;
            Integer actualInt16Value = basicDataIntegerTypeContract.getStorageInt16().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整型int16_t取值】 执行getStorageInt16() successfully actualValue:" + actualInt16Value);
            //collector.assertEqual(actualInt16Value,expectValue2, "checkout  execute success.");
            //5、验证：无符号整数位数
            TransactionReceipt  transactionReceipt1 = basicDataIntegerTypeContract.setUint().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证无符号整数位数】 successfully hash:" + transactionReceipt1.getTransactionHash());


            //6、验证：有符号整数位数
           /* TransactionReceipt  transactionReceipt2 = basicDataIntegerTypeContract.setInt().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证有符号整数位数】 successfully hash:" + transactionReceipt2.getTransactionHash());
            //7、验证：大位数整型赋值
            TransactionReceipt  transactionReceipt3 = basicDataIntegerTypeContract.setBigInt().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证验证大位数整型赋值】 successfully hash:" + transactionReceipt3.getTransactionHash());
*/
        } catch (Exception e) {
            collector.logStepFail("basicDataIntegerTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
