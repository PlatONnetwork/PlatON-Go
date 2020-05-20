package evm.function.assembly;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.SumAssembly;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 1.验证内联汇编在库中的使用
 2.验证汇编的操作码add，mul等
 * @author liweic
 * @dev 2020/01/08 14:30
 */

public class SumAssemblyTest extends ContractPrepareTest {
    private String sum;


    @Before
    public void before() {
        this.prepare();
        sum = driverService.param.get("sum");

    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.SumAssemblyTest-汇编操作指令测试", sourcePrefix = "evm")
    public void Sumassembly() {
        try {
            SumAssembly sumassembly = SumAssembly.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = sumassembly.getContractAddress();
            TransactionReceipt tx = sumassembly.getTransactionReceipt().get();
            collector.logStepPass("SumAssembly deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("SumAssembly deploy gasUsed:" + sumassembly.getTransactionReceipt().get().getGasUsed());

            //验证内联汇编操作指令
            BigInteger result = sumassembly.sum().send();

            collector.logStepPass("SumAssembly返回值：" + result);
            collector.assertEqual(sum ,result.toString());


        } catch (Exception e) {
            collector.logStepFail("SumAssemblyContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}


