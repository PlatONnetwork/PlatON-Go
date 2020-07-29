package evm.oop.Interface;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.InterfaceContractEnumTest;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 测试：在5.0以后版本，接口中可以正常定义枚举
 * @description:
 * @author: qudong
 * @create: 2019/12/25 15:09
 **/
public class InterfaceContractBLimitEnumTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qudong", showName = "InterfaceContractLimitEnum.验证接口中定义枚举",sourcePrefix = "evm")
    public void testInterfaceContractLimitEnum() {

        InterfaceContractEnumTest interfaceContractEnumTest= null;
        try {
            //合约部署
            interfaceContractEnumTest = InterfaceContractEnumTest.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = interfaceContractEnumTest.getContractAddress();
            TransactionReceipt tx = interfaceContractEnumTest.getTransactionReceipt().get();

            collector.logStepPass("InterfaceContractEnumTest issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());

            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("InterfaceContractEnumTest deploy fail.", e.toString());
            e.printStackTrace();
        }

        //调用合约方法

        //1、执行setLarge（）
        try {
            TransactionReceipt setLargeTransaction = interfaceContractEnumTest.setLarge().send();
            collector.logStepPass("调用合约setLarge()方法完毕 successful hash:" + setLargeTransaction.getTransactionHash());
            //collector.assertEqual(sumBigInt, new BigInteger(sumParam), "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractEnumTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }


        //2、执行getChoice（）
        try {
            BigInteger expectChoiceResult = new BigInteger("2");
            BigInteger actualBigInteger = interfaceContractEnumTest.getChoice().send();
            collector.logStepPass("调用合约getChoice()方法完毕 successful actualValue:" + actualBigInteger);
            collector.assertEqual(actualBigInteger,expectChoiceResult, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("InterfaceContractEnumTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }





    }

}
