package evm.data_type.ReferenceData;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ReferenceDataTypeStructDeleteContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：结构体嵌套delete操作
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeStructDeleteContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "ReferenceDataTypeStructDelete.结构体嵌套delete操作",sourcePrefix = "evm")
    public void testReferenceDataTypeStructRecursive() {

        ReferenceDataTypeStructDeleteContract referenceDataTypeStructDelete = null;
        try {
            //合约部署
            referenceDataTypeStructDelete = ReferenceDataTypeStructDeleteContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeStructDelete.getContractAddress();
            TransactionReceipt tx =  referenceDataTypeStructDelete.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeStructDelete issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructDelete deploy fail.", e.toString());
            e.printStackTrace();
        }
        try {
            //1、执行delete基本类型之后值
            BigInteger deleteInt = referenceDataTypeStructDelete.getToDeleteInt().send();
            collector.logStepPass("ReferenceDataTypeStructRecursive 【执行delete基本类型之后值 getToDeleteInt()】方法 successfully.deleteInt：" + deleteInt);
            collector.assertEqual(deleteInt,new BigInteger("0"), "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        try {
            //2、执行delete外部结构体uint类型
            BigInteger actualValue = referenceDataTypeStructDelete.getTopValue().send();
            collector.logStepPass("ReferenceDataTypeStructRecursive 【执行delete外部结构体包含uint类型 getTopValue()】方法 successfully.actualValue：" + actualValue);
            collector.assertEqual(actualValue,new BigInteger("0"), "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        try {
            //3、执行delete外部结构体包含mapping类型
            BigInteger actualValue = referenceDataTypeStructDelete.getTopMapping().send();
            collector.logStepPass("ReferenceDataTypeStructRecursive 【执行delete外部结构体包含mapping类型 getTopMapping()】方法 successfully.actualValue：" + actualValue);
            collector.assertEqual(actualValue,new BigInteger("1"), "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        try {
            //4、执行delete内部部结构体包含uint类型
            BigInteger actualValue = referenceDataTypeStructDelete.getNestedValue().send();
            collector.logStepPass("ReferenceDataTypeStructRecursive 【执行delete内部部结构体包含uint类型 getNestedValue()】方法 successfully.actualValue：" + actualValue);
            collector.assertEqual(actualValue,new BigInteger("0"), "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        try {
            //5、执行delete内部部结构体包含mapping类型
            Boolean actualValue = referenceDataTypeStructDelete.getNestedMapping().send();
            collector.logStepPass("ReferenceDataTypeStructRecursive 【执行delete内部部结构体包含mapping类型 getNestedMapping()】方法 successfully.actualValue：" + actualValue);
            collector.assertEqual(actualValue,true, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructDelete Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
