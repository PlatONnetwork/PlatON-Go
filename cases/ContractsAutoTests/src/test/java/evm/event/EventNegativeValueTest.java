package evm.event;


import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Eventer;
import network.platon.contracts.LoopCall;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title EventNegativeValueTest
 * @description 事件中负值的调用
 * @author qcxiao
 * @updateTime 2020/6/3 19:38
 */
public class EventNegativeValueTest extends ContractPrepareTest {

    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "event.EventNegativeValueTest-事件负值调用", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {
            Eventer eventer = Eventer.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = eventer.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + eventer.getTransactionReceipt().get().getGasUsed());


            TransactionReceipt transactionReceipt = eventer.load(contractAddress, web3j, transactionManager, provider, chainId).getEvent().send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            collector.assertEqual(eventer.getTestInt8Events(transactionReceipt).get(0).out1, BigInteger.valueOf(-2), "事件中第一个值的比较");
            collector.assertEqual(eventer.getTestInt8Events(transactionReceipt).get(0).out2, BigInteger.valueOf(-3), "事件中第二个值的比较");
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
