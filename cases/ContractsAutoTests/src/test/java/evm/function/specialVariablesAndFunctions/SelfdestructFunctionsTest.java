package evm.function.specialVariablesAndFunctions;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.SelfdestructFunctions;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 验证合约自毁函数
 * @description:
 * @author: liweic
 * @create: 2019/12/30 19:01
 **/

public class SelfdestructFunctionsTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.SelfdestructFunctionsTest-合约相关函数测试", sourcePrefix = "evm")
    public void Selfdestructfunction() {
        try {
            SelfdestructFunctions selfdestructFunctions = SelfdestructFunctions.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = selfdestructFunctions.getContractAddress();
            TransactionReceipt tx = selfdestructFunctions.getTransactionReceipt().get();
            collector.logStepPass("SelfdestructFunctionsTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("SelfdestructFunctionsTest deploy gasUsed:" + selfdestructFunctions.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt increaseCount = selfdestructFunctions.increment().send();
            collector.logStepPass("交易Hash：" + increaseCount.getTransactionHash());
            BigInteger resultCount = selfdestructFunctions.getCount().send();
            collector.logStepPass("getCount函数返回值：" + resultCount);
            collector.assertEqual("5",resultCount.toString());


            //调用自杀函数
            TransactionReceipt selfkill = selfdestructFunctions.selfKill().send();
            collector.logStepPass("自杀函数交易Hash：" + selfkill.getTransactionHash());

            collector.logStepPass("Selfdestruct successful.transactionHash:" + selfkill.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + selfkill.getBlockNumber());


            BigInteger count1 = selfdestructFunctions.getCount().send();

            collector.logStepPass("调用自杀函数后链上的count值为："+count1);


        } catch (Exception e) {
            if(e.getMessage().startsWith("Empty")){
                collector.logStepPass("调用自杀函数后链上的count值为:Empty");
            }
            collector.assertContains(e.toString(), "ContractCallException");
            e.printStackTrace();
        }
    }
}

