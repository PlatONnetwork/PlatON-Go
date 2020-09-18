package network.platon.test.evm.oop.inherit;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InheritContractSubclass;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：多重继承(父类合约存在父子关系)，合约继承必须遵循先父到子的继承顺序
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InheritContractBMutipleTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InheritContract.多重继承(遵循先父到子的继承顺序)",sourcePrefix = "evm")
    public void testInheritContractMutipleTest1() {

        InheritContractSubclass inheritContractSubclass = null;
        try {
            //合约部署
            inheritContractSubclass = InheritContractSubclass.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = inheritContractSubclass.getContractAddress();
            TransactionReceipt tx =  inheritContractSubclass.getTransactionReceipt().get();
            collector.logStepPass("InheritContractSubclass issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("InheritContractSubclass deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法
        //1、执行getDataThree()
        try {
            BigInteger expectResult = new BigInteger("3");
            BigInteger actualBigInteger = inheritContractSubclass.getDataThree().send();
            collector.logStepPass("调用合约getDataThree()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InheritContractSubclass Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
