package evm.oop.inherit;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.InheritContractOverload;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：合约函数重载(Overload)：合约可以有多个同名函数，可以有不同输入参数。
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InheritContractOverloadTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InheritContract.多重继承函数重载",sourcePrefix = "evm")
    public void testInheritContractMutipleTest1() {

        InheritContractOverload inheritContractOverload = null;
        try {
            //合约部署
            inheritContractOverload = InheritContractOverload.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = inheritContractOverload.getContractAddress();
            TransactionReceipt tx =  inheritContractOverload.getTransactionReceipt().get();
            collector.logStepPass("InheritContractOverload issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("InheritContractOverload deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、执行getData1()
        try {
            BigInteger expectResult = new BigInteger("3");
            BigInteger actualBigInteger = inheritContractOverload.getDataA().send();
            collector.logStepPass("调用合约getData1()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractSubclass Calling Method fail.", e.toString());
            e.printStackTrace();
        }

        //2、执行getData2()
        try {
            BigInteger expectResult = new BigInteger("6");
            BigInteger actualBigInteger = inheritContractOverload.getDataB().send();
            collector.logStepPass("调用合约getData2()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractSubclass Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
