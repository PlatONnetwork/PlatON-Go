package network.platon.test.evm.oop.abstracttest;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AbstractContractBSubclass;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @title 测试：抽象合约被继承，且实现抽象方法，是否可正常执行
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class AbstractContractBInhertAllImpTest extends ContractPrepareTest {

    private String name,resultName;

    @Before
    public void before() {
       this.prepare();
        name = driverService.param.get("name");
        resultName = driverService.param.get("resultName");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "AbstractContract.抽象合约单继承且实现执行情况",sourcePrefix = "evm")
    public void testAbstractContract() {

        AbstractContractBSubclass abstractContractBSubclass= null;
        try {
            //合约部署
            abstractContractBSubclass = AbstractContractBSubclass.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = abstractContractBSubclass.getContractAddress();
            TransactionReceipt tx = abstractContractBSubclass.getTransactionReceipt().get();
            collector.logStepPass("abstractContract issued successfully.contractAddress:" + contractAddress
                                            + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("abstractContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            //设置用户名称setParentName()
            TransactionReceipt transactionReceipt =  abstractContractBSubclass.setParentName(name).send();
            collector.logStepPass("执行【设置用户名称合约方法setParentName()】,生成hash：" + transactionReceipt.getTransactionHash());
            //获取用户名称parentName()
            String actualValue = abstractContractBSubclass.parentName().send();
            collector.logStepPass("执行【获取用户名称 parentName()】 successful.actualValue:" + actualValue);
            collector.assertEqual(actualValue,resultName, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("abstractContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
