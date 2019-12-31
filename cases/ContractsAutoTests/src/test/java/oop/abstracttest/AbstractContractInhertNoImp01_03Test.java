package oop.abstracttest;

import beforetest.ContractPrepareTest;
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
public class AbstractContractInhertNoImp01_03Test extends ContractPrepareTest {


    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "01AbstractContractExecuteTest_03.抽象合约被继承，但未被实现抽象方法，是否可正常执行")
    public void testAbstractContract() {

        AbstractContractSon sonAbstractContract = null;
        try {
            //合约部署
            sonAbstractContract = AbstractContractSon.deploy(web3j, transactionManager, provider).send();
            String contractAddress = sonAbstractContract.getContractAddress();
            TransactionReceipt tx = sonAbstractContract.getTransactionReceipt().get();

            collector.logStepPass("AbstractContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());

            //collector.assertEqual(tokenName, token.name().send(), "checkout tokenName");
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("grandpaAbstractContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            String sonName = sonAbstractContract.sonName().send();
            collector.logStepPass("调用合约方法完毕 successful.age:" + sonName);

        } catch (Exception e) {
            collector.logStepFail("grandpaAbstractContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }



    }

}
