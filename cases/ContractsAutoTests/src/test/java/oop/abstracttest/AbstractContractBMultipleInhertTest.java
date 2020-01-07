package oop.abstracttest;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.AbstractContractCSubclass;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @title 测试：普通合约是否可以继承多个抽象合约,且实现抽象方法，是否可以正常编译部署执行
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class AbstractContractBMultipleInhertTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "AbstractContract.合约多继承执行情况")
    public void testAbstractContract() {

        AbstractContractCSubclass abstractContractCSubclass= null;
        try {
            //合约部署
            abstractContractCSubclass = AbstractContractCSubclass.deploy(web3j, transactionManager, provider).send();
            String contractAddress = abstractContractCSubclass.getContractAddress();
            TransactionReceipt tx = abstractContractCSubclass.getTransactionReceipt().get();

            collector.logStepPass("AbstractContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());

            //collector.assertEqual(tokenName, token.name().send(), "checkout tokenName");
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("abstractContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            String expectValue = "cSubName";
            String actualValue = abstractContractCSubclass.cSubName().send();
            collector.logStepPass("调用合约方法完毕 successful.actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectValue, "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("abstractContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }



    }

}
