package evm.oop.abstracttest;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.AbstractContractSon;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @title 测试：抽象合约被继承，但未被实现抽象方法，是否可正常执行
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class AbstractContractAInhertNoImpTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "AbstractContract.抽象合约被继承未实现执行情况",sourcePrefix = "evm")
    public void testAbstractContract() {

        AbstractContractSon sonAbstractContract = null;
        try {
            //合约部署
            sonAbstractContract = AbstractContractSon.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = sonAbstractContract.getContractAddress();
            TransactionReceipt tx = sonAbstractContract.getTransactionReceipt().get();
            collector.logStepPass("abstractContract issued successfully.contractAddress:" + contractAddress
                                           + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("abstractContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            String expectResult = "sonName";
            String actualName = sonAbstractContract.sonName().send();
            collector.logStepFail("abstractContract Calling Method Fail.","抽象合约被继承未实现方法无法执行");
        } catch (Exception e) {
            collector.logStepPass("执行【抽象合约被继承未实现方法，调用函数sonName()】,抽象合约部分实现是无法执行方法");
            collector.assertEqual(e.getMessage(),"Empty value (0x) returned from contract","checkout  execute success.");
            //e.printStackTrace();
        }
    }
}
