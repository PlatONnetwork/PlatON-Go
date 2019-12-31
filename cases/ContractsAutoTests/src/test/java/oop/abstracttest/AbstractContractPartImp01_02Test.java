package oop.abstracttest;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.AbstractContractFather;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 1、抽象合约实现部分方法，验证是否可编译、部署、执行
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class AbstractContractPartImp01_02Test extends ContractPrepareTest {


    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "01AbstractContractExecuteTest_02.验证抽象合约(实现部分抽象方法)是否可以编译部署执行")
    public void testAbstractContract() {

        AbstractContractFather fatherAbstractContract = null;
        try {
            //合约部署
            fatherAbstractContract = AbstractContractFather.deploy(web3j, transactionManager, provider).send();
            String contractAddress = fatherAbstractContract.getContractAddress();
            TransactionReceipt tx = fatherAbstractContract.getTransactionReceipt().get();

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
            BigInteger age = fatherAbstractContract.fatherAge().send();
            collector.logStepPass("调用合约方法完毕 successful.age:" + age.toString());

        } catch (Exception e) {
            collector.logStepFail("grandpaAbstractContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }



    }

}
