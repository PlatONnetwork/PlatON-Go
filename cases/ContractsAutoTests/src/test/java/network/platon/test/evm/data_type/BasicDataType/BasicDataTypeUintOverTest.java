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
 * @title 测试：无符号8位整数数据溢出(8位无符号整数取值范围0~255)
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class BasicDataTypeUintOverTest extends ContractPrepareTest {

    private String uintParam;
    private String resultParam;


    @Before
    public void before() {
       this.prepare();
        uintParam = driverService.param.get("uintParam");
        resultParam = driverService.param.get("resultParam");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "BasicDataTypeUintOverTest.无符号8位整数数据溢出",sourcePrefix = "evm")
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
        //1、验证：无符号8位整数数据溢出(8位无符号整数取值范围0~255)
        try {
            BigInteger resultValue = new BigInteger(resultParam);
            //赋值执行addUintOverflow()
            BigInteger actualValue = basicDataTypeContract.addUintOverflow(new BigInteger(uintParam)).send();
            collector.logStepPass("BasicDataTypeContract 执行addUintOverflow() successfully actualValue:" + actualValue +",resultValue:"+ resultValue);
            collector.assertEqual(actualValue,resultValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("BasicDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }



    }
}
