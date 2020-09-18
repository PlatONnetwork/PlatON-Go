package network.platon.test.evm.data_type.BasicDataType;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.BasicDataTypeContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：合约基本数据类型
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class BasicDataTypeContractTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "BasicDataTypeUintOverTest.合约基本数据类型",sourcePrefix = "evm")
    public void testBasicDataTypeContract() {

        BasicDataTypeContract basicDataTypeContract = null;
        try {
            //合约部署
            basicDataTypeContract = BasicDataTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = basicDataTypeContract.getContractAddress();
            TransactionReceipt tx =  basicDataTypeContract.getTransactionReceipt().get();
            collector.logStepPass("BasicDataTypeContract issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、验证：定长字节数组
        try {
            BigInteger expectLength = new BigInteger("2");
            //赋值执行getBytes1Length()
            BigInteger actualLength = basicDataTypeContract.getBytes1Length().send();
            collector.logStepPass("BasicDataTypeContract 执行getBytes1Length() successfully actualValue:" + actualLength);
            collector.assertEqual(actualLength,expectLength, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、验证：变长字节数组
        try {
            BigInteger expectLength = new BigInteger("3");
            //赋值执行getBytesLength()
            BigInteger actualLength = basicDataTypeContract.getBytesLength().send();
            collector.logStepPass("BasicDataTypeContract 执行getBytesLength() successfully actualValue:" + actualLength);
            collector.assertEqual(actualLength,expectLength, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
