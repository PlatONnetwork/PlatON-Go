package evm.function.Modifier;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Modifiers;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 验证多修饰器
 * @description:
 * @author: liweic
 * @create: 2020/01/02 20:01
 **/

public class ModifiersTest extends ContractPrepareTest {

    private String modifierstest;

    @Before
    public void before() {
        this.prepare();
        modifierstest = driverService.param.get("modifierstest");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.ModifiersTest-多修饰器函数测试", sourcePrefix = "evm")
    public void modifiers() {
        try {
            Modifiers modifiers = Modifiers.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = modifiers.getContractAddress();
            TransactionReceipt tx = modifiers.getTransactionReceipt().get();
            collector.logStepPass("modifiers deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("modifiers deploy gasUsed:" + modifiers.getTransactionReceipt().get().getGasUsed());

            //验证多修饰器函数调用
            TransactionReceipt result = modifiers.test1().send();
            collector.logStepPass("交易hash：" + result.getTransactionHash());

            BigInteger a = modifiers.test2().send();
            collector.logStepPass("modifiers函数返回值：" + a);
            collector.assertEqual(modifierstest ,a.toString());


        } catch (Exception e) {
            collector.logStepFail("ModifiersContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}


