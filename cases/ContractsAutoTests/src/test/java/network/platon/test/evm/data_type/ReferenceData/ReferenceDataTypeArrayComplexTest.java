package network.platon.test.evm.data_type.ReferenceData;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ReferenceDataTypeArrayComplexContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

/**
 * @title 测试：验证含数组（Array）运算逻辑合约
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeArrayComplexTest extends ContractPrepareTest {

    private String a,b,c,d,e;
    private String sum;


    @Before
    public void before() {
       this.prepare();
        a = driverService.param.get("a");
        b = driverService.param.get("b");
        c = driverService.param.get("c");
        d = driverService.param.get("d");
        e = driverService.param.get("e");
        sum = driverService.param.get("sum");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "ReferenceDataTypeArrayComplex.含数组（Array）运算逻辑合约",sourcePrefix = "evm")
    public void testReferenceDataTypeArrayTest() {
        try{
            ReferenceDataTypeArrayComplexContract referenceDataTypeArrayComplex = null;
            referenceDataTypeArrayComplex = ReferenceDataTypeArrayComplexContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeArrayComplex.getContractAddress();
            TransactionReceipt tx =  referenceDataTypeArrayComplex.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeArrayComplex issued successfully.contractAddress:" + contractAddress
                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

            if(getIntParam("seq") == 3){
                List<BigInteger> array = new ArrayList<BigInteger>();
                BigInteger actualValue = referenceDataTypeArrayComplex.sumComplexArray(array).send();
                collector.logStepPass("referenceDataTypeArrayComplex 执行sumComplexArray() successfully.hash:" + actualValue);
                collector.assertEqual(actualValue,BigInteger.ZERO, "checkout execute success.");
                return;
            }
            BigInteger sumBig = new BigInteger(sum);
            List<BigInteger> array = new ArrayList<BigInteger>();
            array.add(new BigInteger(a));
            array.add(new BigInteger(b));
            array.add(new BigInteger(c));
            array.add(new BigInteger(d));
            array.add(new BigInteger(e));

            //赋值执行sumComplexArray()
            BigInteger actualValue = referenceDataTypeArrayComplex.sumComplexArray(array).send();
            collector.logStepPass("referenceDataTypeArrayComplex 执行sumComplexArray() successfully.hash:" + actualValue);
            collector.assertEqual(actualValue,sumBig, "checkout execute success.");
        }catch (Exception e){
            collector.logStepFail(Thread.currentThread().getStackTrace()[1].getMethodName() + " fail.", e.getMessage());
            e.printStackTrace();
        }

    }

}
