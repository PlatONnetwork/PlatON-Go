package evm.function.contractRelatedAndFallback;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.FallBack;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 1.验证合约关键字this,表示当前合约，可以显示的转换为Address
 2.验证Fallback函数,调用了未命名的函数等方式
 * @author liweic
 * @dev 2020/01/02 18:00
 */

public class FallbackTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.FallbackTest-this和回退函数测试", sourcePrefix = "evm")
    public void fallback() {
        try {
            FallBack fallback = FallBack.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = fallback.getContractAddress();
            TransactionReceipt tx = fallback.getTransactionReceipt().get();
            collector.logStepPass("FallBack deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("FallBack deploy gasUsed:" + fallback.getTransactionReceipt().get().getGasUsed());

            //验证this和回退函数
            BigInteger a = fallback.getA().send();
            collector.logStepPass("回退函数调用前a的值：" + a);

            TransactionReceipt functionnotexist = fallback.CallFunctionNotExist().send();
            collector.logStepPass("打印fallback交易hash：" + functionnotexist.getTransactionHash());

            BigInteger falla = fallback.getA().send();
            collector.logStepPass("回退函数调用后a的值：" + falla);
            collector.assertEqual(new BigInteger("100"),falla);

        } catch (Exception e) {
            collector.logStepFail("FallbackContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}




