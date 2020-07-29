package evm.oop.inherit;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.InheritContractAMutipleClass;
import network.platon.contracts.InheritContractBMutipleClass;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：多重合约继承重名问题，遵循最远继承原则
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InheritContractAMutipleTest extends ContractPrepareTest {


    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InheritContract.合约继承重名问题(遵循最远继承原则)",sourcePrefix = "evm")
    public void testInheritContracAtMutipleTest() {

        InheritContractAMutipleClass inheritContractMutipleTest1 = null;
        InheritContractBMutipleClass inheritContractMutipleTest2 = null;
        try {
            //合约部署(InheritContractMutipleTest1)
            inheritContractMutipleTest1 = InheritContractAMutipleClass.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = inheritContractMutipleTest1.getContractAddress();
            TransactionReceipt tx =  inheritContractMutipleTest1.getTransactionReceipt().get();
            collector.logStepPass("InheritContractMutipleTest1 issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());


            //合约部署(InheritContractMutipleTest1)
            inheritContractMutipleTest2 = InheritContractBMutipleClass.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress2 = inheritContractMutipleTest2.getContractAddress();
            TransactionReceipt tx2 =  inheritContractMutipleTest2.getTransactionReceipt().get();
            collector.logStepPass("InheritContractMutipleTest2 issued successfully.contractAddress:" + contractAddress2
                    + ", hash:" + tx2.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx2.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("InheritContractMutipleTest deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、执行callGetDate1()
        try {
            BigInteger expectBookResult = new BigInteger("1");
            BigInteger actualBigInteger = inheritContractMutipleTest1.callGetDateA().send();
            collector.logStepPass("调用合约callGetDateA()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractMutipleTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、执行callGetDate2()
        try {
            BigInteger expectBookResult = new BigInteger("2");
            BigInteger actualBigInteger = inheritContractMutipleTest2.callGetDateB().send();
            collector.logStepPass("调用合约callGetDateB()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractMutipleTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
