package evm.data_type.ReferenceData;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ReferenceDataTypeArrayContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;

import java.math.BigInteger;

/**
 * @title 测试：多维数组声明及初始化及取值
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeMultiArrayContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "ReferenceDataTypeMultiArray.多维数组声明及初始化及取值",sourcePrefix = "evm")
    public void testReferenceDataTypeArrayTest() {

        ReferenceDataTypeArrayContract referenceDataTypeArrayContract = null;
        try {
            //合约部署
            referenceDataTypeArrayContract = ReferenceDataTypeArrayContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeArrayContract.getContractAddress();
            TransactionReceipt tx =  referenceDataTypeArrayContract.getTransactionReceipt().get();
            collector.logStepPass("ReferenceDataTypeArrayContract issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeArrayContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法

        //验证：多维数组声明及初始化及取值
        try {
            BigInteger expectValue = new BigInteger("100");
            BigInteger expectLength = new BigInteger("2");
            //赋值执行setMultiArray()
            TransactionReceipt transactionReceipt = referenceDataTypeArrayContract.setMultiArray().send();
            collector.logStepPass("ReferenceDataTypeArrayContract 执行setMultiArray() successfully.hash:" + transactionReceipt.getTransactionHash());
            //获取值getMultiArray()
            Tuple2<BigInteger, BigInteger> tuple2 = referenceDataTypeArrayContract.getMultiArray().send();
            BigInteger actualValue = tuple2.getValue1();
            BigInteger actualLength = tuple2.getValue2();

            collector.logStepPass("调用合约setMultiArray() 方法完毕 successful actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout value execute success.");
            collector.assertEqual(actualLength,expectLength, "checkout length execute success.");
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeArrayContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
