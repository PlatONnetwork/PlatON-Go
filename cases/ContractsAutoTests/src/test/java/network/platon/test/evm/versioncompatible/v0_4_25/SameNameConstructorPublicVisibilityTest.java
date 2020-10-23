package network.platon.test.evm.versioncompatible.v0_4_25;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.SameNameConstructorPublicVisibility;
import org.junit.Test;
import org.web3j.tuples.generated.Tuple6;

import java.math.BigInteger;
/**
 * @title 构造函数和可见性测试
 * 1. 0.4.25版本同名函数构造函数定义，声明public可见性验证（0.5.x版本弃用同名函数定义构造函数）；
 * 2. 0.4.25版本fallback函数可见性非强制声明（默认public可见性）验证
 * 3. 0.4.25版本支持字面量及后缀（0.5.x版本已弃用）验证
 * （1）0.4.25版本支持year时间单位
 * （2）0.4.25版本允许小数点后不跟数字的数值写法
 * （3）0.4.25版本十六进制数字支持带“0X”和“0x”等2种前缀表示
 * （4）0.4.25版本支持十六进制数与以太币单位组合
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class SameNameConstructorPublicVisibilityTest extends ContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testDiscardLiteralsAndSuffixes",
            author = "albedo", showName = "network.platon.test.evm.SameNameConstructorPublicVisibilityTest-构造函数和可见性", sourcePrefix = "evm")
    public void testDiscardLiteralsAndSuffixes() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000");
            SameNameConstructorPublicVisibility visibility = SameNameConstructorPublicVisibility.deploy(web3j, transactionManager, provider, chainId, constructorValue).send();
            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SameNameConstructorPublicVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger> result = visibility.discardLiteralsAndSuffixes().send();
            Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger> expect =
                    new Tuple6(new BigInteger("10000000000"), new BigInteger("31536000"), new BigInteger("255000000000000000000"),
                            new BigInteger("255"), new BigInteger("255"), new BigInteger("255000000000000000000"));
            collector.assertEqual(result, expect, "checkout visibility assignment result");
        } catch (Exception e) {
            collector.logStepFail("SameNameConstructorPublicVisibilityTest testDiscardLiteralsAndSuffixes failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
