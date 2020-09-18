package network.platon.test.evm.function.functionVisibilityAndDecarations;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.Inter;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 验证internal在继承合约里的调用
 * @description:
 * @author: liweic
 * @create: 2020/01/02 16:01
 **/

public class InterTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.InterTest-函数可见性继承合约调用内部方法测试", sourcePrefix = "evm")
    public void inter() {
        try {
            Inter intercall = Inter.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = intercall.getContractAddress();
            TransactionReceipt tx = intercall.getTransactionReceipt().get();
            collector.logStepPass("Inter deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("Inter deploy gasUsed:" + intercall.getTransactionReceipt().get().getGasUsed());

            //验证继承合约可以调用父合约的内部方法
            BigInteger interdata = intercall.g().send();
            collector.logStepPass("public可见函数返回值：" + interdata);
            collector.assertEqual(new BigInteger("3"),interdata);

        } catch (Exception e) {
            collector.logStepFail("InterContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}



