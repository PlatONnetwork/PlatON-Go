package network.platon.test.evm.versioncompatible.v0_4_25;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ConstructorPublicVisibility;
import org.junit.Test;
import org.web3j.tuples.generated.Tuple2;

import java.math.BigInteger;
/**
 * @title 构造函数和可见性测试
 * 1. constructor声明构造函数，public可见性验证；
 * 2. 允许调用没有括号的基类构造函数验证
 * 3. ，0.4.25支持语法(0.5.X版本已弃用)覆盖验证
 * （1）0.4.25版本允许未实现的函数使用修饰符(modifier)验证()
 * （2）0.4.25版本允许布尔变量使用算术运算验证
 * （3）0.4.25版本允许使用一元运算符"+"验证
 * （4）0.4.25版本允许在if包含的块中使用单个语句声明/定义变量验证
 * （5）0.4.25版本允许在while包含的块中使用单个语句声明/定义变量验证使用算术原酸运算符
 * （6）0.4.25版本允许在for包含的块中使用单个语句声明/定义变量验证
 * （7）0.4.25版本允许具有一个或多个返回值的函数使用空返回语句验证
 * （8）0.4.25版本允许具有一个或多个返回值的函数使用空返回语句验证
 *  4. 0.4.25版本允许constant用作修饰函数状态可变性验证
 *  5. 0.4.25版本允许定义具有命名返回值的函数类型验证
 *  6. 0.4.25版本允许 msg.value用在非 payable函数里以及此函数的修饰符(modifier)里验证
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class ConstructorPublicVisibilityTest extends ContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testConstantCheck",
            author = "albedo", showName = "network.platon.test.evm.ConstructorPublicVisibilityTest-constant修饰函数状态可变性", sourcePrefix = "evm")
    public void testConstantCheck() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000000");
            ConstructorPublicVisibility visibility = ConstructorPublicVisibility.deploy(web3j, transactionManager, provider, chainId,constructorValue).send();
            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ConstructorPublicVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            BigInteger constantCheck = visibility.constantCheck().send();
            collector.assertEqual(constantCheck, constructorValue, "checkout allows constant to be used as a modifier for state variability");
        } catch (Exception e) {
            collector.logStepFail("ConstructorPublicVisibilityTest testConstantCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testNamedReturn",
            author = "albedo", showName = "ConstructorPublicVisibilityTest-定义具有命名返回值的函数类型", sourcePrefix = "evm")
    public void testNamedReturn() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000");
            ConstructorPublicVisibility visibility = ConstructorPublicVisibility.deploy(web3j, transactionManager, provider, chainId,constructorValue).send();
            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            BigInteger result = visibility.namedReturn(new BigInteger("12"),new BigInteger("23")).send();
            collector.logStepPass("ConstructorPublicVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            collector.assertEqual(result, new BigInteger("35"), "checkout allows you to define function types with named return values");
        } catch (Exception e) {
            collector.logStepFail("ConstructorPublicVisibilityTest testNamedReturn failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testGrammarCheck",
            author = "albedo", showName = "ConstructorPublicVisibilityTest-支持语法覆盖", sourcePrefix = "evm")
    public void testGrammarCheck() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000000");
            ConstructorPublicVisibility visibility = ConstructorPublicVisibility.deploy(web3j, transactionManager, provider, chainId,constructorValue).send();
            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            Tuple2<BigInteger,BigInteger> result = visibility.grammarCheck().send();
            Tuple2<BigInteger,BigInteger> checkResult=new Tuple2<>(new BigInteger("0"),new BigInteger("0"));
            collector.logStepPass("ConstructorPublicVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            collector.assertEqual(result, checkResult, "checkout discard syntax validation ");
        } catch (Exception e) {
            collector.logStepFail("ConstructorPublicVisibilityTest testGrammarCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testAbstractFunction",
            author = "albedo", showName = "ConstructorPublicVisibilityTest-未实现的函数使用修饰符", sourcePrefix = "evm")
    public void testAbstractFunction() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000000");
            ConstructorPublicVisibility visibility = ConstructorPublicVisibility.deploy(web3j, transactionManager, provider, chainId,constructorValue).send();
            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            BigInteger result = visibility.abstractFunction().send();
            collector.logStepPass("ConstructorPublicVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            collector.assertEqual(result, constructorValue.add(new BigInteger("123")), "checkout allows unimplemented functions are validated with a modifier");
        } catch (Exception e) {
            collector.logStepFail("ConstructorPublicVisibilityTest testAbstractFunction failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testShow",
            author = "albedo", showName = "ConstructorPublicVisibilityTest-msg.value用在非 payable函数里以及此函数的修饰符", sourcePrefix = "evm")
    public void testShow() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000000");
            ConstructorPublicVisibility visibility = ConstructorPublicVisibility.deploy(web3j, transactionManager, provider, chainId,constructorValue).send();
            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            BigInteger result = visibility.show().send();
            collector.logStepPass("ConstructorPublicVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            collector.assertEqual(result, new BigInteger("0"), "checkout msg.value is allowed to be used in non-payable functions and in the modifier for this function");
        } catch (Exception e) {
            collector.logStepFail("ConstructorPublicVisibilityTest testShow failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
