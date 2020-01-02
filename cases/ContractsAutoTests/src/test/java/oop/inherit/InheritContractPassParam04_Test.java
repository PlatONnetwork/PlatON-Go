package oop.inherit;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.InheritContractSub1;
import network.platon.contracts.InheritContractSub2;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：继承支持传参(继承中基类构造函数的传参)
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InheritContractPassParam04_Test extends ContractPrepareTest {


    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "01InheritContractMultipleTest.多重合约继承重名问题(遵循最远继承原则)")
    public void testInheritContractMutipleTest1() {

        InheritContractSub1 inheritContractSub1 = null;
        InheritContractSub2 inheritContractSub2 = null;
        try {
            //合约部署(inheritContractSub1)
            inheritContractSub1 = InheritContractSub1.deploy(web3j, transactionManager, provider).send();
            String contractAddress = inheritContractSub1.getContractAddress();
            TransactionReceipt tx =  inheritContractSub1.getTransactionReceipt().get();
            collector.logStepPass("InheritContractSub1 issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());


            //合约部署(inheritContractSub2)
            inheritContractSub2 = InheritContractSub2.deploy(web3j, transactionManager, provider).send();
            String contractAddress2 = inheritContractSub2.getContractAddress();
            TransactionReceipt tx2 =  inheritContractSub2.getTransactionReceipt().get();
            collector.logStepPass("inheritContractSub2 issued successfully.contractAddress:" + contractAddress2
                    + ", hash:" + tx2.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx2.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("inheritContractSub deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、执行getData1()
        try {
            BigInteger expectBookResult = new BigInteger("2");
            BigInteger actualBigInteger = inheritContractSub1.getData1().send();
            collector.logStepPass("调用合约getData1()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractStructTest1 Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、执行getData2()
        try {
            BigInteger expectBookResult = new BigInteger("4");
            BigInteger actualBigInteger = inheritContractSub2.getData2().send();
            collector.logStepPass("调用合约getData1()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractStructTest1 Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }

}
