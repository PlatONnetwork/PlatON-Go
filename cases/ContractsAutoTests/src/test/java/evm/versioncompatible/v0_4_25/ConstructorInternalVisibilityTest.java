package evm.versioncompatible.v0_4_25;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ConstructorInternalVisibility;
import org.junit.Test;

import java.math.BigInteger;
/**
 * @title 构造函数和可见性测试
 * 1. 0.4.25版本验证使用constructor关键字定义构造函数，使用internal声明可见性
 * 2. 0.4.25版本验证子合约直接声明父合约构造函数，但是构造函数参数与父合约不一致
 *  如：父合约：constructor(uint _x) 子合约：is Base()
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class ConstructorInternalVisibilityTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "evm.ConstructorInternalVisibilityTest-构造函数和可见性", sourcePrefix = "evm")
    public void testGetOutI() {
        try {
            prepare();
            BigInteger constructorValue = new BigInteger("10000000000000");
            ConstructorInternalVisibility visibility = ConstructorInternalVisibility.deploy(web3j, transactionManager, provider, chainId, constructorValue).send();

            String contractAddress = visibility.getContractAddress();
            String transactionHash = visibility.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ConstructorInternalVisibility issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + visibility.getTransactionReceipt().get().getGasUsed());
            BigInteger outI = visibility.getOutI().send();
            collector.assertEqual(outI, constructorValue, "checkout constructor initial param");
        } catch (Exception e) {
            collector.logStepFail("ConstructorInternalVisibilityTest testGetOutI failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
