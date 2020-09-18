package network.platon.test.evm.data_type.ReferenceData;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ReferenceDataTypeArrayContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：数组（Array）方法
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeArrayMethodTest extends ContractPrepareTest {

    private String insertValue;
    private String arrayLength;

    @Before
    public void before() {
       this.prepare();
        insertValue = driverService.param.get("insertValue");
        arrayLength = driverService.param.get("arrayLength");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "ReferenceDataTypeMultiArrayContractTest.数组（Array）方法",sourcePrefix = "evm")
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
        //验证：数组的属性及方法
        try {
            BigInteger expectLength = new BigInteger(arrayLength);
            //赋值执行setArrayPush()
            TransactionReceipt transactionReceipt = referenceDataTypeArrayContract.setArrayPush(insertValue).send();
            collector.logStepPass("ReferenceDataTypeArrayContract 执行setArrayPush() successfully.hash:" + transactionReceipt.getTransactionHash());
           //获取数组长度getArrayLength()
            BigInteger actualLength = referenceDataTypeArrayContract.getArrayLength().send();
            collector.logStepPass("调用合约getArrayLength()方法完毕 successful actualLength:" + actualLength);
            collector.assertEqual(actualLength,expectLength, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeArrayContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }

}
