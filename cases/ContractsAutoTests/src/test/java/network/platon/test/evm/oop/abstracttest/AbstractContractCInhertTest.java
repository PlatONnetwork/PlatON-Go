package network.platon.test.evm.oop.abstracttest;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AbstractContractFSubclass;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @title 测试：抽象合约是否可以继承抽象合约，是否可以正常编译部署执行
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class AbstractContractCInhertTest extends ContractPrepareTest {

   private String name,resultName;

    @Before
    public void before() {
       this.prepare();
        name = driverService.param.get("name");
        resultName = driverService.param.get("resultName");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "AbstractContract.抽象合约继承抽象合约执行情况",sourcePrefix = "evm")
    public void testAbstractContract() {

        AbstractContractFSubclass abstractContractFSubclass= null;
        try {
            //合约部署
            abstractContractFSubclass = AbstractContractFSubclass.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = abstractContractFSubclass.getContractAddress();
            TransactionReceipt tx = abstractContractFSubclass.getTransactionReceipt().get();
            collector.logStepPass("abstractContract issued successfully.contractAddress:" + contractAddress
                                           + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("abstractContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            //设置用户名称setASubName()
            TransactionReceipt transactionReceipt =  abstractContractFSubclass.setParentNameD(name).send();
            collector.logStepPass("执行【设置用户名称合约方法setParentNameD()】,生成hash：" + transactionReceipt.getTransactionHash());
            //获取用户名称aSubName()
            String actualValue = abstractContractFSubclass.parentName().send();
            collector.logStepPass("执行【获取用户名称 parentName()】 successful.actualValue:" + actualValue);
            collector.assertEqual(actualValue,resultName, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("abstractContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
