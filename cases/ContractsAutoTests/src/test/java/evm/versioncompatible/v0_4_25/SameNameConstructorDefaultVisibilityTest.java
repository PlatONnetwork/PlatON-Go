package evm.versioncompatible.v0_4_25;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.SameNameConstructorDefaultVisibility;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
/**
 * @title 构造函数和可见性测试
 * 1.0.4.25版本验证同名函数构造函数定义，可见性未声明（默认public可见性）验证；
 * 2.函数可见性非强制声明验证。5种函数
 * (1) 函数默认可见性声明
 * (2) 函数public可见性声明
 * (3) 函数external可见性声明
 * (4) 函数internal可见性声明
 * (5) 函数private可见性声明
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class SameNameConstructorDefaultVisibilityTest extends ContractPrepareTest {
    SameNameConstructorDefaultVisibility visibility;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testSameNameConstructorVisibility",
            author = "albedo", showName = "evm.SameNameConstructorDefaultVisibilityTest-同名函数未声明可见性构造函数", sourcePrefix = "evm")
    public void testSameNameConstructorVisibility() {
        try {
            TransactionReceipt sameNameConstructorTx = deployContract();
            collector.logStepPass("SameNameConstructorDefaultVisibility issued successfully sameNameConstructorTx contractAddress:" + sameNameConstructorTx.getContractAddress() + ", hash:" + sameNameConstructorTx.getTransactionHash());
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorDefaultVisibilityTest testDoWhileCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testDefaultVisibility",
            author = "albedo", showName = "SameNameConstructorDefaultVisibilityTest-函数默认可见性声明", sourcePrefix = "evm")
    public void testDefaultVisibility() {
        try {
            deployContract();
            BigInteger result = visibility.defaultVisibility(new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("12"), "checkout visibility assignment result");
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorDefaultVisibilityTest testDefaultVisibility failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testPublicVisibility",
            author = "albedo", showName = "SameNameConstructorDefaultVisibilityTest-函数public可见性声明", sourcePrefix = "evm")
    public void testPublicVisibility() {
        try {
            deployContract();
            BigInteger result = visibility.publicVisibility(new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("12"), "checkout visibility assignment result");
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorDefaultVisibilityTest testPublicVisibility failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testExternalVisibility",
            author = "albedo", showName = "SameNameConstructorDefaultVisibilityTest-函数external可见性声明", sourcePrefix = "evm")
    public void testExternalVisibility() {
        try {
            deployContract();
            BigInteger result = visibility.externalVisibility(new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("12"), "checkout visibility assignment result");
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorDefaultVisibilityTest testExternalVisibility failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testPrivateVisibilityCheck",
            author = "albedo", showName = "SameNameConstructorDefaultVisibilityTest-函数private可见性声明", sourcePrefix = "evm")
    public void testPrivateVisibilityCheck() {
        try {
            deployContract();
            BigInteger result = visibility.privateVisibilityCheck(new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("12"), "checkout visibility assignment result");
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorDefaultVisibilityTest testPrivateVisibilityCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testInternalVisibilityCheck",
            author = "albedo", showName = "SameNameConstructorDefaultVisibilityTest-函数internal可见性声明", sourcePrefix = "evm")
    public void testInternalVisibilityCheck() {
        try {
            deployContract();
            BigInteger result = visibility.internalVisibilityCheck(new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("12"), "checkout visibility assignment result");
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorDefaultVisibilityTest testInternalVisibilityCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    private TransactionReceipt deployContract() throws Exception {
        prepare();
        BigInteger constructorValue = new BigInteger("10000000000");
        visibility = SameNameConstructorDefaultVisibility.deploy(web3j, transactionManager, provider, chainId).send();
        String contractAddress = visibility.getContractAddress();
        String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
        collector.logStepPass("SameNameConstructorDefaultVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
        collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
        return visibility.SameNameConstructorVisibility(constructorValue).send();
    }
}
