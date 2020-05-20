package evm.data_type.ReferenceData;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ReferenceDataTypeArrayOperatorContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tuples.generated.Tuple5;

import java.math.BigInteger;

/**
 * @title 测试：验证数组支持的运算符
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeArrayOperatorTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "ReferenceDataTypeArray.数组支持的运算符",sourcePrefix = "evm")
    public void testReferenceDataTypeArrayTest() {

        ReferenceDataTypeArrayOperatorContract referenceDataTypeArrayOperator = null;
        try {
            //合约部署
            referenceDataTypeArrayOperator = ReferenceDataTypeArrayOperatorContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeArrayOperator.getContractAddress();
            TransactionReceipt tx =  referenceDataTypeArrayOperator.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeArrayOperator issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeArrayOperator deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            //1、比较运算符 arrayCompare()
            Tuple5<Boolean, Boolean, Boolean, Boolean, Boolean> tuple5 = referenceDataTypeArrayOperator.arrayCompare().send();
            boolean actualValue1 = tuple5.getValue1();
            boolean actualValue2 = tuple5.getValue2();
            boolean actualValue3 = tuple5.getValue3();
            boolean actualValue4 = tuple5.getValue4();
            boolean actualValue5 = tuple5.getValue5();

            collector.logStepPass("执行【比较运算符 arrayCompare()】 successfully.hash:" + tuple5.toString());
            collector.assertEqual(actualValue1,true, "checkout < execute success.");
            collector.assertEqual(actualValue2,false, "checkout > execute success.");
            collector.assertEqual(actualValue3,true, "checkout == execute success.");
            collector.assertEqual(actualValue4,true, "checkout != execute success.");
            collector.assertEqual(actualValue5,true, "checkout >= execute success.");

            //2、&(按位与)
            Tuple2<byte[], BigInteger> tuple2 =referenceDataTypeArrayOperator.arrayBitAndOperators().send();
            BigInteger andOperatorsValue = tuple2.getValue2();
            collector.logStepPass("执行【&(按位与) arrayBitAndOperators()】 successfully.hash:" + tuple2.toString());
            collector.assertEqual(andOperatorsValue,new BigInteger("128"), "checkout execute success.");

            //3、|(按位或)
            Tuple2<byte[], BigInteger> tuple3 =referenceDataTypeArrayOperator.arrayBitOrOperators().send();
            BigInteger orOperatorsValue = tuple3.getValue2();
            collector.logStepPass("执行【|(按位或) arrayBitOrOperators()】 successfully.hash:" + tuple3.toString());
            collector.assertEqual(orOperatorsValue,new BigInteger("129"), "checkout execute success.");

            //4、~（按位取反）
            Tuple2<byte[], BigInteger> tuple4 =referenceDataTypeArrayOperator.arrayBitInverseOperators().send();
            BigInteger inverseOperatorsValue = tuple4.getValue2();
            collector.logStepPass("执行【~（按位取反） arrayBitInverseOperators()】 successfully.hash:" + tuple4.toString());
            collector.assertEqual(inverseOperatorsValue,new BigInteger("126"), "checkout execute success.");

            //5、^（按位异或）
            Tuple2<byte[], BigInteger> tuple6 =referenceDataTypeArrayOperator.arrayBitXOROperators().send();
            BigInteger xOROperatorsValue = tuple6.getValue2();
            collector.logStepPass("执行【^（按位异或） arrayBitXOROperators()】 successfully.hash:" + tuple6.toString());
            collector.assertEqual(xOROperatorsValue,new BigInteger("1"), "checkout execute success.");

            //6、<<（左移位）
            Tuple2<byte[], BigInteger> tuple7 =referenceDataTypeArrayOperator.arrayBitLeftShiftperators().send();
            BigInteger leftShiftperatorsValue = tuple7.getValue2();
            collector.logStepPass("执行【<<（左移位） arrayBitLeftShiftperators()】 successfully.hash:" + tuple7.toString());
            collector.assertEqual(leftShiftperatorsValue,new BigInteger("2"), "checkout execute success.");

            //7、<<（右移位）
            Tuple2<byte[], BigInteger> tuple8 =referenceDataTypeArrayOperator.arrayBitRightShiftperators().send();
            BigInteger rightShiftperatorsValue = tuple8.getValue2();
            collector.logStepPass("执行【<<（右移位） arrayBitLeftShiftperators()】 successfully.hash:" + tuple8.toString());
            collector.assertEqual(rightShiftperatorsValue,new BigInteger("64"), "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeArrayOperator Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }

}
