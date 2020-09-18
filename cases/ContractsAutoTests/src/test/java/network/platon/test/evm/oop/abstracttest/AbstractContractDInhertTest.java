package network.platon.test.evm.oop.abstracttest;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AbstractContractGSubclass;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：抽象合约是否可以继承接口(反之接口是否可以继承抽象合约)
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class AbstractContractDInhertTest extends ContractPrepareTest {

   private String age,resultAge;

    @Before
    public void before() {
       this.prepare();
        age = driverService.param.get("age");
        resultAge = driverService.param.get("resultAge");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "AbstractContract.抽象合约继承接口执行情况",sourcePrefix = "evm")
    public void testAbstractContract() {

        AbstractContractGSubclass abstractContractGSubclass= null;
        try {
            //合约部署
            abstractContractGSubclass = AbstractContractGSubclass.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = abstractContractGSubclass.getContractAddress();
            TransactionReceipt tx = abstractContractGSubclass.getTransactionReceipt().get();
            collector.logStepPass("abstractContract issued successfully.contractAddress:" + contractAddress
                                           + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("abstractContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            //设置用户年龄setInterAge()
            BigInteger resultValue = new BigInteger(resultAge);
            TransactionReceipt transactionReceipt =  abstractContractGSubclass.setInterAge(new BigInteger(age)).send();
            collector.logStepPass("执行【设置用户年龄合约方法setInterAge()】,生成hash：" + transactionReceipt.getTransactionHash());
            //获取用户名称getInterAge()
            BigInteger actualValue = abstractContractGSubclass.aInterAge().send();
            collector.logStepPass("执行【获取用户年龄 getInterAge()】 successful.actualValue:" + actualValue);
            collector.assertEqual(actualValue,resultValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("abstractContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
