package data_type.TypeConversion;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.TypeConversionContractTest;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;

import java.math.BigInteger;

/**
 * @title 测试：基本类型之间的转换（隐式/显示）
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class TypeConversionContract01_Test extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "TypeConversionContract.基本类型之间的转换（隐式/显示）")
    public void testTypeConversionContract() {

        TypeConversionContractTest typeConversionContractTest = null;
        try {
            //合约部署
            typeConversionContractTest = TypeConversionContractTest.deploy(web3j, transactionManager, provider).send();
            String contractAddress = typeConversionContractTest.getContractAddress();
            TransactionReceipt tx =  typeConversionContractTest.getTransactionReceipt().get();
            collector.logStepPass("TypeConversionContractTest issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、验证：基本类型隐式转换(运算符操作隐式转换)
        try {
            BigInteger expectValue = new BigInteger("102");
            //赋值执行sum()
            BigInteger actualValue = typeConversionContractTest.sum().send();
            collector.logStepPass("TypeConversionContractTest 运算符操作隐式转换执行sum() successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、验证：基本类型隐式转换(赋值操作隐式转换)
        try {
            BigInteger expectValue = new BigInteger("10");
            //赋值执行conversion()
            BigInteger actualValue = typeConversionContractTest.conversion().send();
            collector.logStepPass("TypeConversionContractTest 赋值操作隐式转换执行conversion() successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //3、验证：基本类型显示转换(无符合与有符号转换)
        try {
            BigInteger expectValue = new BigInteger("1");
            //执行displayConversion()
            BigInteger actualValue = typeConversionContractTest.displayConversion().send();
            collector.logStepPass("TypeConversionContractTest 无符合与有符号转换执行displayConversion() successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //4、验证：基本类型显示转换(转换成更小的类型，会丢失高位)
        try {
            BigInteger expectValue1 = new BigInteger("22136");
            String expectValue2 = "0x5678";
            //执行displayConversion1()
            Tuple2<BigInteger, byte[]> tuple2 = typeConversionContractTest.displayConversion1().send();
            BigInteger actualValue1 = tuple2.getValue1();
            String actualValue2 = "0x" + DataChangeUtil.bytesToHex(tuple2.getValue2());

            collector.logStepPass("TypeConversionContractTest 转换成更小的类型执行displayConversion1() successfully.actualValue1:" + actualValue1 +
                                  ",actualValue2:" + actualValue2);
            collector.assertEqual(actualValue1,expectValue1, "checkout execute success.");
            collector.assertEqual(actualValue2,expectValue2, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //5、验证：转换成更大的类型(将向左侧添加填充位)
        try {
            BigInteger expectValue1 = new BigInteger("4660");
            //执行displayConversion2()
            Tuple2<BigInteger, byte[]> tuple2 = typeConversionContractTest.displayConversion2().send();
            BigInteger actualValue1 = tuple2.getValue1();

            collector.logStepPass("TypeConversionContractTest 转换成更大的类型执行displayConversion2() successfully.actualValue1:" + actualValue1);
            collector.assertEqual(actualValue1,expectValue1, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //6、验证：转换到更小的字节类型
        try {
            String expectValue = "0x12";
            //执行displayConversion3()
            byte[] byteValue = typeConversionContractTest.displayConversion3().send();
            String actualValue = "0x" + DataChangeUtil.bytesToHex(byteValue);

            collector.logStepPass("TypeConversionContractTest 转换到更小的字节类型执行displayConversion3() successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //7、验证：转换为更大的字节类型
        try {
            String expectValue = "0x12340000";
            //执行displayConversion4()
            byte[] byteValue = typeConversionContractTest.displayConversion4().send();
            String actualValue = "0x" + DataChangeUtil.bytesToHex(byteValue);

            collector.logStepPass("TypeConversionContractTest 转换为更大的字节类型执行displayConversion4() successfully.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("TypeConversionContractTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }





    }

}
