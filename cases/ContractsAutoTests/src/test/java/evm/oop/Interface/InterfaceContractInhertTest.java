package evm.oop.Interface;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.InterfaceContractInheritMultipleTest;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：合约是否可以继承多个接口
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InterfaceContractInhertTest extends ContractPrepareTest {

   private String a,b,sumValue;
   private String c,d,reduceValue;

    @Before
    public void before() {
       this.prepare();
        a = driverService.param.get("a");
        b = driverService.param.get("b");
        sumValue = driverService.param.get("sumValue");
        c = driverService.param.get("c");
        d = driverService.param.get("d");
        reduceValue = driverService.param.get("reduceValue");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "AbstractContract.合约多继承接口执行情况",sourcePrefix = "evm")
    public void testAbstractContract() {

        InterfaceContractInheritMultipleTest interfaceInheritMultiple = null;
        try {
            //合约部署
            interfaceInheritMultiple = InterfaceContractInheritMultipleTest.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = interfaceInheritMultiple.getContractAddress();
            TransactionReceipt tx = interfaceInheritMultiple.getTransactionReceipt().get();
            collector.logStepPass("interfaceContract issued successfully.contractAddress:" + contractAddress
                                           + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("interfaceContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            BigInteger resultSum = new BigInteger(sumValue);
            BigInteger resultReduce = new BigInteger(reduceValue);

            //合约加法 sum()
            BigInteger actualSumValue =  interfaceInheritMultiple.sum(new BigInteger(a),new BigInteger(b)).send();
            collector.logStepPass("执行【合约加法 sum()】 actualSumValue：" + actualSumValue);
            collector.assertEqual(resultSum,actualSumValue, "checkout  execute success.");

            //合约减法 reduce()
            BigInteger actualReduceValue = interfaceInheritMultiple.reduce(new BigInteger(c),new BigInteger(d)).send();
            collector.logStepPass("执行【合约减法 reduce()】 actualReduceValue:" + actualReduceValue);
            collector.assertEqual(resultReduce,actualReduceValue, "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("interfaceContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
