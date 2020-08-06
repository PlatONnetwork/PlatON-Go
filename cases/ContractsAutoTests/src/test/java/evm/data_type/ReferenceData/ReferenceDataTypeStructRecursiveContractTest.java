package evm.data_type.ReferenceData;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ReferenceDataTypeStructRecursiveContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple3;

import java.math.BigInteger;

/**
 * @title 测试：结构体嵌套递归
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeStructRecursiveContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "ReferenceDataTypeStructRecursive.结构体嵌套递归",sourcePrefix = "evm")
    public void testReferenceDataTypeStructRecursive() {

        ReferenceDataTypeStructRecursiveContract referenceDataTypeStructRecursive = null;
        try {
            //合约部署
            referenceDataTypeStructRecursive = ReferenceDataTypeStructRecursiveContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeStructRecursive.getContractAddress();
            TransactionReceipt tx =  referenceDataTypeStructRecursive.getTransactionReceipt().get();
            collector.logStepPass("ReferenceDataTypeStructRecursive issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeStructRecursive deploy fail.", e.toString());
            e.printStackTrace();
        }
        try {
            //结构体嵌套递归，获取结构体数组长度
            Tuple3<BigInteger, BigInteger, BigInteger> tuple3 = referenceDataTypeStructRecursive.getStructPersonLength().send();
            BigInteger actuaLength1 = tuple3.getValue1();
            BigInteger actuaLength2 = tuple3.getValue2();
            BigInteger actuaLength3 = tuple3.getValue3();
            collector.logStepPass("ReferenceDataTypeStructRecursive 【执行获取结构体数组长度 getStructPersonLength()】方法 successfully.tuple3：" + tuple3.toString());
            collector.assertEqual(actuaLength1,new BigInteger("2"), "checkout Id execute success.");
            collector.assertEqual(actuaLength2,new BigInteger("10"), "checkout Age execute success.");
            collector.assertEqual(actuaLength3,new BigInteger("20"), "checkout VIP execute success.");
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeStructRecursive Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
