package wasm.data_type;

import com.platon.rlp.datatypes.Int32;
import com.platon.rlp.datatypes.Int64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.BasicDataTypeContract;
import org.junit.Before;
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

    private String uint8ByteValueStr;
    private String stringValueStr;
    private String stringValueStrLength;
    private String floatValueStr;


    @Before
    public void before() {
        uint8ByteValueStr = driverService.param.get("uint8ByteValueStr");
        stringValueStr = driverService.param.get("stringValueStr");
        stringValueStrLength = driverService.param.get("stringValueStrLength");
        floatValueStr = driverService.param.get("floatValueStr");

    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.basicDataTypeTest基本类型验证测试",sourcePrefix = "wasm")
    public void testBasicDataType() {
         //部署合约
        BasicDataTypeContract basicDataTypeContract = null;
        try {
            prepare();
            basicDataTypeContract = BasicDataTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = basicDataTypeContract.getContractAddress();
            TransactionReceipt tx = basicDataTypeContract.getTransactionReceipt().get();
            collector.logStepPass("basicDataTypeContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
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
            Boolean actualBoolValue = basicDataTypeContract.get_bool().send();
            collector.logStepPass("basicDataTypeContract 【验证布尔值取值】 执行getBool() successfully actualBoolValue:" + actualBoolValue);
            collector.assertEqual(actualBoolValue,boolValue, "checkout  execute success.");
            //2、验证：字节类型（byte）
            Uint8 uint8ByteValue = Uint8.of(uint8ByteValueStr);
            TransactionReceipt transactionReceipt1 = basicDataTypeContract.set_byte(uint8ByteValue).send();
            collector.logStepPass("basicDataTypeContract 【验证字节类型（byte）】 successfully hash:" + transactionReceipt1.getTransactionHash());
            Uint8 actualByteValue = basicDataTypeContract.get_byte().send();
            collector.logStepPass("basicDataTypeContract 【验证字节类型取值】 执行get_byte() successfully actualByteValue:" + actualByteValue);
            collector.assertEqual(actualByteValue,uint8ByteValue, "checkout  execute success.");
            //3、验证：字符串赋值
            TransactionReceipt transactionReceipt2 = basicDataTypeContract.set_string(stringValueStr).send();
            collector.logStepPass("basicDataTypeContract 【验证字符串赋值】 successfully hash:" + transactionReceipt2.getTransactionHash());
            String actualStringValue = basicDataTypeContract.get_string().send();
            collector.logStepPass("basicDataTypeContract 【验证字符串取值】 执行get_string() successfully actualStringValue:" + actualStringValue);
            collector.assertEqual(actualStringValue,stringValueStr, "checkout  execute success.");
            //4、验证：字符串长度
            Uint8 actualStringLength = basicDataTypeContract.get_string_length().send();
            collector.logStepPass("basicDataTypeContract 【验证字符串长度】 执行get_string_length() successfully actualStringLength:" + actualStringLength);
            collector.assertEqual(actualStringLength,Uint8.of(stringValueStrLength), "checkout  execute success.");
            //5、地址类型(Address)
            TransactionReceipt  transactionReceipt4 = basicDataTypeContract.set_address().send();
            collector.logStepPass("basicDataTypeContract 【验证地址类型(Address)】 successfully hash:" + transactionReceipt4.getTransactionHash());
            String actualAddreeValue = basicDataTypeContract.get_address().send();
            collector.logStepPass("basicDataTypeContract 【验证地址取值】 执行getString() successfully actualAddreeValue:" + actualAddreeValue);
            //collector.assertEqual(actualStringValue,expectStringValue, "checkout  execute success.")
            //6、浮点类型(float)
            // Float floatValue = Float.parseFloat(floatValueStr);//-3.4E-38f
           /*  Float floatValue = 1.5f;
            TransactionReceipt  transactionReceipt5 = basicDataTypeContract.set_float(floatValue).send();
            collector.logStepPass("basicDataTypeContract 【验证浮点类型(float)】 successfully hash:" + transactionReceipt5.getTransactionHash());
            Float actualFloatValue = basicDataTypeContract.get_float().send();
            collector.logStepPass("basicDataTypeContract 【验证浮点类型(float)取值】 执行get_float() successfully actualFloatValue:" + actualFloatValue);
            collector.assertEqual(actualFloatValue,floatValue, "checkout  execute success.");
            //7、浮点类型(double)
           // Double doubleValue = 6.577;
            Double doubleValue = 2.4791E2;
            TransactionReceipt  transactionReceipt6 = basicDataTypeContract.set_double(doubleValue).send();
            collector.logStepPass("basicDataTypeContract 【验证浮点类型(double)】 successfully hash:" + transactionReceipt6.getTransactionHash());
            Double actualDoubleValue = basicDataTypeContract.get_double().send();
            collector.logStepPass("basicDataTypeContract 【验证浮点类型(double)取值】 执行get_double() successfully actualDoubleValue:" + actualDoubleValue);
            collector.assertEqual(actualDoubleValue,doubleValue, "checkout  execute success.");*/
           //8、验证浮点型局部变量
            TransactionReceipt  transactionReceipt6 = basicDataTypeContract.set_float_type_local().send();
            collector.logStepPass("basicDataTypeContract 【验证浮点型局部变量】 successfully hash:" + transactionReceipt6.getTransactionHash());

            //9、验证long整型
            Int32 longValue = Int32.of(50);
            TransactionReceipt  transactionReceipt7 = basicDataTypeContract.set_long(longValue).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证long整型类型】 successfully hash:" + transactionReceipt7.getTransactionHash());
            Int32 actualLongValue = basicDataTypeContract.get_long().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数long取值】 执行get_long() successfully actualLongValue:" + actualLongValue);
            collector.assertEqual(actualLongValue,longValue, "checkout  execute success.");

            //10、验证long long整型
            Int64 longlongValue = Int64.of(100);
            TransactionReceipt  transactionReceipt8 = basicDataTypeContract.set_long_long(longlongValue).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证longlong整型类型】 successfully hash:" + transactionReceipt8.getTransactionHash());
            Int64 actualLongLongValue = basicDataTypeContract.get_long_long().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数longlong取值】 执行get_long_long() successfully actualLongLongValue:" + actualLongLongValue);
            collector.assertEqual(actualLongLongValue,longlongValue, "checkout  execute success.");
            //11、enum枚举赋值
            Uint8 actualEnumValue = basicDataTypeContract.set_enum_assignment().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证enum枚举赋值】 执行set_enum_assignment() successfully actualEnumValue:" + actualEnumValue);
            collector.assertEqual(actualEnumValue,Uint8.of(3), "checkout  execute success.");
            //12、enum限制作用域枚举赋值
            Uint8 actualEnumValue1 = basicDataTypeContract.set_enum_class_assignment().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证enum限制作用域枚举赋值】 执行set_enum_class_assignment() successfully actualEnumValue1:" + actualEnumValue1);
            collector.assertEqual(actualEnumValue1,Uint8.of(2), "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("basicDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
