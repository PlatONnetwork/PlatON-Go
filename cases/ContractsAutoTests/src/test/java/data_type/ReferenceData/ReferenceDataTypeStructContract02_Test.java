package data_type.ReferenceData;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ReferenceDataTypeStructContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple4;

import java.math.BigInteger;

/**
 * @title 测试：结构体定义、赋值及取值
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeStructContract02_Test extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "02ReferenceDataTypeStructContract.结构体定义、赋值及取值")
    public void testReferenceDataTypeStructTest() {

        ReferenceDataTypeStructContract referenceDataTypeStructContract = null;
        try {
            //合约部署
            referenceDataTypeStructContract = ReferenceDataTypeStructContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = referenceDataTypeStructContract.getContractAddress();
            TransactionReceipt tx =  referenceDataTypeStructContract.getTransactionReceipt().get();
            collector.logStepPass("ReferenceDataTypeStructContract issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeStructContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        BigInteger expectId = new BigInteger("2");
        BigInteger expectAge = new BigInteger("25");
        boolean expectVIP = true;

        //1、赋值方式一： 按入参顺序赋值
        try {
            //赋值执行initDataStruct1()
            Tuple4<BigInteger, String, BigInteger, Boolean> tuple4 = referenceDataTypeStructContract.initDataStruct1().send();
            BigInteger actualId = tuple4.getValue1();
            BigInteger actualAge = tuple4.getValue3();
            boolean actualVIP = tuple4.getValue4();

            collector.logStepPass("ReferenceDataTypeStructContract 执行initDataStruct1()方法 successfully.");
            collector.assertEqual(actualId,expectId, "checkout Id execute success.");
            collector.assertEqual(actualAge,expectAge, "checkout Age execute success.");
            collector.assertEqual(actualVIP,expectVIP, "checkout VIP execute success.");
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeStructContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }


        //2、赋值方式二： 按命名参数赋值
        try {
            //赋值执行initDataStruct2()
            Tuple4<BigInteger, String, BigInteger, Boolean> tuple4 = referenceDataTypeStructContract.initDataStruct2().send();
            BigInteger actualId = tuple4.getValue1();
            BigInteger actualAge = tuple4.getValue3();
            boolean actualVIP = tuple4.getValue4();

            collector.logStepPass("ReferenceDataTypeStructContract 执行initDataStruct2()方法 successfully.");
            collector.assertEqual(actualId,expectId, "checkout Id execute success.");
            collector.assertEqual(actualAge,expectAge, "checkout Age execute success.");
            collector.assertEqual(actualVIP,expectVIP, "checkout VIP execute success.");
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeStructContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //3、赋值方式三：结构体中映射的初始化
        try {
            //赋值执行initDataStruct3()
            Tuple4<BigInteger, String, BigInteger, Boolean> tuple4 = referenceDataTypeStructContract.initDataStruct3().send();
            BigInteger actualId = tuple4.getValue1();
            BigInteger actualAge = tuple4.getValue3();
            boolean actualVIP = tuple4.getValue4();

            collector.logStepPass("ReferenceDataTypeStructContract 执行initDataStruct3()方法 successfully.");
            collector.assertEqual(actualId,expectId, "checkout Id execute success.");
            collector.assertEqual(actualAge,expectAge, "checkout Age execute success.");
            collector.assertEqual(actualVIP,expectVIP, "checkout VIP execute success.");
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeStructContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }

}
