package network.platon.test.evm.data_type.BasicDataType;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.BasicDataTypeConstantContract;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple3;

import java.math.BigInteger;

/**
 * @title 测试：合约基本数据类型字面常量
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class BasicDataTypeConstantContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "BasicDataTypeConstantContractTest.合约基本数据类型字面常量",sourcePrefix = "evm")
    public void testBasicDataTypeContract() {

        BasicDataTypeConstantContract basicDataTypeConstantContract = null;
        try {
            //合约部署
            basicDataTypeConstantContract = BasicDataTypeConstantContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = basicDataTypeConstantContract.getContractAddress();
            TransactionReceipt tx =  basicDataTypeConstantContract.getTransactionReceipt().get();
            collector.logStepPass("BasicDataTypeConstantContract issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、验证：address类型，获取当前合约余额
        try {
            BigInteger expectValue = new BigInteger("0");
            //执行getCurrentBalance()
            BigInteger actualValue = basicDataTypeConstantContract.getCurrentBalance().send();
            collector.logStepPass("BasicDataTypeConstantContract 执行查询当前合约中余额 getCurrentBalance() successfully contractBalance:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、验证：字符串字面常量，进行赋值
        try {
            String expectValue = "hello";
            //执行getCurrentBalance()
            String actualValue = basicDataTypeConstantContract.getStrA().send();
            collector.logStepPass("BasicDataTypeConstantContract 字符串字面常量进行赋值 执行getStrA() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //3、验证：字符串是特殊的动态长度字节数组，获取长度
        try {
            BigInteger expectValue = new BigInteger("5");
            //执行getCurrentBalance()
            BigInteger actualValue = basicDataTypeConstantContract.getStrALength().send();
            collector.logStepPass("BasicDataTypeConstantContract 验证字符串字面常量 执行getStrALength() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //4、验证：字符串是特殊的动态长度字节数组，进行转换
        try {
            String expectValue = "aello";
            //赋值执行setStr1()
            String actualValue = basicDataTypeConstantContract.setStrA().send();
            collector.logStepPass("BasicDataTypeConstantContract 验证字符串字面常量 执行setStrA() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }


        //5、验证：十六进制字面常量，定义赋值取值
        try {
            //执行getHexLiteral()
            String expectValue = "c8";
            byte[] byteValue = basicDataTypeConstantContract.getHexLiteraA().send();
            String actualValue = DataChangeUtil.bytesToHex(byteValue);
            collector.logStepPass("BasicDataTypeConstantContract 验证十六进制字面常量 执行getHexLiteraA() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");

            //执行getHexLitera2()
            String expectValue2 = "01f4";
            byte[] byteValue2 = basicDataTypeConstantContract.getHexLiteraB().send();
            String actualValue2 = DataChangeUtil.bytesToHex(byteValue2);
            collector.logStepPass("BasicDataTypeConstantContract 验证十六进制字面常量 执行getHexLiteraB() successfully actualValue2:" + actualValue2);
            collector.assertEqual(actualValue2,expectValue2, "checkout  execute success.");

            //执行getHexLitera3()
            String expectValueA = "01f4";
            String expectValueB = "01";
            String expectValueC = "f4";
            Tuple3<byte[], byte[], byte[]> tuple3 = basicDataTypeConstantContract.getHexLiteraC().send();
            String actualValueA = DataChangeUtil.bytesToHex(tuple3.getValue1());
            String actualValueB = DataChangeUtil.bytesToHex(tuple3.getValue2());
            String actualValueC = DataChangeUtil.bytesToHex(tuple3.getValue3());
            collector.logStepPass("BasicDataTypeConstantContract 验证十六进制字面常量 执行getHexLiteraC() successfully  actualValueA:" + actualValueA +
                                  ",  actualValueB:" + actualValueB + ",  actualValueC:" + actualValueC);

            collector.assertEqual(actualValueA,expectValueA, "checkout  execute success.");
            collector.assertEqual(actualValueB,expectValueB, "checkout  execute success.");
            collector.assertEqual(actualValueC,expectValueC, "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //6、验证：枚举类型
        try {
            BigInteger expectValue1 = new BigInteger("1");
            BigInteger expectValue2 = new BigInteger("3");
            //执行getSeason1()
            BigInteger actualValue1 = basicDataTypeConstantContract.getSeasonA().send();
            collector.logStepPass("BasicDataTypeConstantContract 验证枚举类型执行 getSeasonA() successfully actualValue1:" + actualValue1);
            collector.assertEqual(actualValue1,expectValue1, "checkout  execute success.");

            //执行getSeason2()
            BigInteger actualValue2 = basicDataTypeConstantContract.getSeasonB().send();
            collector.logStepPass("BasicDataTypeConstantContract 验证枚举类型执行 getSeasonB() successfully actualValue2:" + actualValue2);
            collector.assertEqual(actualValue2,expectValue2, "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //7、验证：有理数常量
        try {
            BigInteger expectValue1 = new BigInteger("3");
            BigInteger expectValue2 = new BigInteger("20000000000");
            BigInteger expectValue3 = new BigInteger("3000000000000000");
            //执行getValue()
            Tuple3<BigInteger, BigInteger, BigInteger> tuple3 = basicDataTypeConstantContract.getValue().send();
            BigInteger actualValue1 = tuple3.getValue1();
            BigInteger actualValue2 = tuple3.getValue2();
            BigInteger actualValue3 = tuple3.getValue3();

            collector.logStepPass("BasicDataTypeConstantContract 验证有理数常量执行 getValue() successfully actualValue1:" + actualValue1 +
                                  ",  actualValue2: " + actualValue2 + ",  actualValue3:" + actualValue3) ;
            collector.assertEqual(actualValue1,expectValue1, "checkout  execute success.");
            collector.assertEqual(actualValue2,expectValue2, "checkout  execute success.");
            collector.assertEqual(actualValue3,expectValue3, "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeConstantContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
