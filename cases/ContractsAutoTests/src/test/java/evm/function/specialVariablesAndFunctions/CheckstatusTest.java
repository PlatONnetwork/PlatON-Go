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
 * @title 验证同一个合约部署两次后，两个合约间是否会相互影响
 * @description:
 * @author: hudenian
 * @create: 2020/02/05 13:37
 **/

public class CheckstatusTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "function.CheckstatusTest-两个合约间是否会相互影响", sourcePrefix = "evm")
    public void Selfdestructfunction() {
        try {
            //第一次部署合约
            SelfdestructFunctions selfdestructFunctions = SelfdestructFunctions.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = selfdestructFunctions.getContractAddress();
            TransactionReceipt tx = selfdestructFunctions.getTransactionReceipt().get();
            collector.logStepPass("SelfdestructFunctionsTest first deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("SelfdestructFunctionsTest first deploy gasUsed:" + selfdestructFunctions.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt increaseCount = selfdestructFunctions.increment().send();
            collector.logStepPass("交易Hash：" + increaseCount.getTransactionHash());
            BigInteger resultCount = selfdestructFunctions.getCount().send();
            collector.logStepPass("getCount函数返回值：" + resultCount);
            collector.assertEqual("5",resultCount.toString());


            //第二次部署合约
            SelfdestructFunctions selfdestructFunctionsTwo = SelfdestructFunctions.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddressTwo = selfdestructFunctionsTwo.getContractAddress();
            tx = selfdestructFunctionsTwo.getTransactionReceipt().get();
            collector.logStepPass("SelfdestructFunctionsTest second deploy successfully.contractAddress:" + selfdestructFunctionsTwo + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("SelfdestructFunctionsTest second deploy gasUsed:" + selfdestructFunctionsTwo.getTransactionReceipt().get().getGasUsed());

            increaseCount = selfdestructFunctionsTwo.increment().send();
            collector.logStepPass("交易Hash：" + increaseCount.getTransactionHash());
            resultCount = selfdestructFunctionsTwo.getCount().send();
            collector.logStepPass("getCount函数返回值：" + resultCount);
            collector.assertEqual("5",resultCount.toString());

            //第二个合约自杀后看对第一个合约是否有影响
            TransactionReceipt selfkill = selfdestructFunctionsTwo.selfKill().send();
            collector.logStepPass("自杀函数交易Hash：" + selfkill.getTransactionHash());
            resultCount = selfdestructFunctions.getCount().send();
            collector.logStepPass("第二个合约自杀后查询第一个合约getCount函数返回值：" + resultCount);
            collector.assertEqual("5",resultCount.toString());

            //调用第二个合约查询函数返回值（抛出异常）
            BigInteger count1 = selfdestructFunctionsTwo.getCount().send();
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

