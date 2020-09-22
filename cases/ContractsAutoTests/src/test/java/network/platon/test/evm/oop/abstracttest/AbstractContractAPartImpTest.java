package network.platon.test.evm.oop.abstracttest;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.AbstractContractFather;
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
public class AbstractContractAPartImpTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "AbstractContract.抽象合约实现部分执行情况",sourcePrefix = "evm")
    public void testAbstractContract() {

        AbstractContractFather fatherAbstractContract = null;
        try {
            //合约部署
            fatherAbstractContract = AbstractContractFather.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = fatherAbstractContract.getContractAddress();
            TransactionReceipt tx = fatherAbstractContract.getTransactionReceipt().get();
            collector.logStepPass("abstractContract issued successfully.contractAddress:" + contractAddress
                                           + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("abstractContract deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        try {
            BigInteger age = fatherAbstractContract.fatherAge().send();
            collector.logStepFail("abstractContract Calling Method Fail.","抽象合约部分实现是无法执行方法");
        } catch (Exception e) {
            collector.logStepPass("执行【抽象合约部分实现调用函数fatherAge()】,抽象合约部分实现是无法执行方法");
            collector.assertEqual(e.getMessage(),"Empty value (0x) returned from contract","checkout  execute success.");
            //e.printStackTrace();
        }
    }
}
