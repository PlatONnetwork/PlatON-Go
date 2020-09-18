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
 * @title 测试：数组（Array）赋值取值及方法
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class ReferenceDataTypeSetArrayContractTest extends ContractPrepareTest {

  private String insertNo;
  private String insertValue;

    @Before
    public void before() {
       this.prepare();
        insertNo = driverService.param.get("insertNo");
        insertValue = driverService.param.get("insertValue");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "ReferenceDataTypeArray.数组的声明及初始化及取值",sourcePrefix = "evm")
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
        //1、验证：数组的声明及初始化及取值(定长数组、可变数组)
        try {
            BigInteger index =  new BigInteger(insertNo);
            BigInteger value = new BigInteger(insertValue);
            //赋值执行setArray()
            TransactionReceipt transactionReceipt = referenceDataTypeArrayContract.setArray(index,value).send();
            collector.logStepPass("ReferenceDataTypeArrayContract 执行setArray() successfully.hash:" + transactionReceipt.getTransactionHash());
            //获取值getArray()
            BigInteger actualValue = referenceDataTypeArrayContract.getArray(index).send();
            collector.logStepPass("调用合约getArray()方法完毕 successful actualValue:" + actualValue);
            collector.assertEqual(actualValue,value, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractSubclass Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
