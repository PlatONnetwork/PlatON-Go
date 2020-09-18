package network.platon.test.evm.versioncompatible.v0_5_0.v1_isibility;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.FallbackDeclaraction;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title fallback函数必须声明为external
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class FallbackDeclaractionTest extends ContractPrepareTest {

    //链上函数初始值
    private String initValue;

    @Before
    public void before() {
        this.prepare();
        initValue = driverService.param.get("initValue");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "FallbackDeclaractionTest-可见性测试", sourcePrefix = "evm")
    public void update() {
        try {

            FallbackDeclaraction fallbackDeclaraction = FallbackDeclaraction.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = fallbackDeclaraction.getContractAddress();
            TransactionReceipt tx = fallbackDeclaraction.getTransactionReceipt().get();

            initValue = fallbackDeclaraction.getA().send().toString();

            collector.logStepPass("链上函数的初始值为："+initValue);

            collector.logStepPass("FallbackDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + fallbackDeclaraction.getTransactionReceipt().get().getGasUsed());

            //调用不存在函数，将触发回退函数，导致A修改成111
            TransactionReceipt transactionReceipt = fallbackDeclaraction.callNonExistFunc().send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String afterValue = fallbackDeclaraction.getA().send().toString();
            collector.logStepPass("链上函数的执行update后的值为："+afterValue);

            collector.assertEqual("111",afterValue);
        } catch (Exception e) {
            collector.logStepFail("FallbackDeclaractionTest process fail.", e.toString());
            e.printStackTrace();
        }
    }


}
