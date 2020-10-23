package network.platon.test.evm.oop.Interface;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InterfaceContractStructTest;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：在5.0以后版本，接口可以声明结构体
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InterfaceContractBLimitStructTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InterfaceContractLimitStruct.验证接口定义结构体", sourcePrefix = "evm")
    public void testInterfaceContractLimitEnum() {

        InterfaceContractStructTest interfaceContractStructTest= null;
        try {
            //合约部署
            interfaceContractStructTest = InterfaceContractStructTest.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = interfaceContractStructTest.getContractAddress();
            TransactionReceipt tx = interfaceContractStructTest.getTransactionReceipt().get();

            collector.logStepPass("InterfaceContractStructTest issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());

            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("InterfaceContractStructTest deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法

        //1、执行setBook（）
        try {
            TransactionReceipt setBookTransaction = interfaceContractStructTest.setBook().send();
            collector.logStepPass("调用合约setBook()方法完毕 successful hash:" + setBookTransaction.getTransactionHash());
            //collector.assertEqual(sumBigInt, new BigInteger(sumParam), "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractStructTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }


        //2、执行getBookID()
        try {
            BigInteger expectBookResult = new BigInteger("1");
            BigInteger actualBigInteger = interfaceContractStructTest.getBookID().send();
            collector.logStepPass("调用合约getBookID()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectBookResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractStructTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}
