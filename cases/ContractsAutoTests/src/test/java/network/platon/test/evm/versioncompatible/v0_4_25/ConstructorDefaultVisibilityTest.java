package network.platon.test.evm.versioncompatible.v0_4_25;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ConstructorDefaultVisibility;
import org.junit.Test;

import java.math.BigInteger;
/**
 * @title ERC200412Token测试
 * 1. 0.4.25版本验证使用constrictor关键字定义构造函数，但是不强制声明可见性(默认为public可见性）
 * 2. 0.4.25版本同一继承层次结构中允许多次指定基类构造函数参数验证:
 * (1) 允许合约直接声明构造函数 ———— is Base(7)
 *（2）子合约构造函数继承父合约构造函数———— constructor(uint _y) Base(_y * _y)
 * 两种引用构造函数方式共存时，合约优先选择（2）方式
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class ConstructorDefaultVisibilityTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "network.platon.test.evm.ConstructorDefaultVisibilityTest-构造函数和可见性", sourcePrefix = "evm")
    public void testGetOutI() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000000");
            ConstructorDefaultVisibility visibility = ConstructorDefaultVisibility.deploy(web3j, transactionManager, provider, chainId, constructorValue).send();

            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ConstructorDefaultVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            BigInteger outI = visibility.getOutI().send();
            collector.assertEqual(outI, constructorValue, "checkout constructor initial param");
        } catch (Exception e) {
            collector.logStepFail("ConstructorDefaultVisibilityTest testGetOutI failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
