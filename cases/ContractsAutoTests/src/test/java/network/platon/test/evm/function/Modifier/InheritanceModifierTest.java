package network.platon.test.evm.function.Modifier;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.InheritanceModifier;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 *1.验证单一修饰器
 *2.验证特殊_的用法，符合函数修饰器定义的条件，才可以执行函数体内容
 *3.验证修饰器可以接收参数
 *4.验证合约继承情况下的修饰器的使用
 * @author liweic
 * @dev 2020/01/02 20:50
 */

public class InheritanceModifierTest extends ContractPrepareTest {
    private String modifiertest;

    @Before
    public void before() {
        this.prepare();
        modifiertest = driverService.param.get("modifiertest");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.InheritanceModifierTest-单修饰器函数测试", sourcePrefix = "evm")
    public void inheritancemodifier() {
        try {
            InheritanceModifier inheritanceModifier = InheritanceModifier.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = inheritanceModifier.getContractAddress();
            TransactionReceipt tx = inheritanceModifier.getTransactionReceipt().get();
            collector.logStepPass("InheritanceModifier deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("InheritanceModifier deploy gasUsed:" + inheritanceModifier.getTransactionReceipt().get().getGasUsed());

            //验证单修饰器函数调用
            BigInteger result = inheritanceModifier.getA().send();
            collector.logStepPass("InheritanceModifier函数返回值：" + result);
            collector.assertEqual(modifiertest ,result.toString());


        } catch (Exception e) {
            collector.logStepFail("InheritanceModifierContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}



