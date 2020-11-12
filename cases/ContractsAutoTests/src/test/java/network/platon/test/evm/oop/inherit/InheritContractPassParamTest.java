package network.platon.test.evm.oop.inherit;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InheritContractASub;
import network.platon.contracts.evm.InheritContractBSub;
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
public class InheritContractPassParamTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InheritContract.合约继承支持传参",sourcePrefix = "evm")
    public void testInheritContractMutipleTest1() {

        InheritContractASub inheritContractSub1 = null;
        InheritContractBSub inheritContractSub2 = null;
        try {
            //合约部署(inheritContractASub)
            inheritContractSub1 = InheritContractASub.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = inheritContractSub1.getContractAddress();
            TransactionReceipt tx =  inheritContractSub1.getTransactionReceipt().get();
            collector.logStepPass("InheritContractASub issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());


            //合约部署(inheritContractBSub)
            inheritContractSub2 = InheritContractBSub.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress2 = inheritContractSub2.getContractAddress();
            TransactionReceipt tx2 =  inheritContractSub2.getTransactionReceipt().get();
            collector.logStepPass("InheritContractBSub issued successfully.contractAddress:" + contractAddress2
                    + ", hash:" + tx2.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx2.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("inheritContractSub deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、执行getDataA()
        try {
            BigInteger expectBookResult = new BigInteger("2");
            BigInteger actualBigInteger = inheritContractSub1.getDataA().send();
            collector.logStepPass("调用合约getDataA()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractStructTest1 Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、执行getDataB()
        try {
            BigInteger expectBookResult = new BigInteger("4");
            BigInteger actualBigInteger = inheritContractSub2.getDataB().send();
            collector.logStepPass("调用合约getDataB()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractStructTest1 Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }

}
