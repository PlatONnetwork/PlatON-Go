package wasm.function;

import com.platon.rlp.datatypes.Int32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MemoryCallocInt;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * 验证内存memory函数-整型存储
 * 1.函数calloc
 * @author liweic
 * @create 2020/02/24
 */

public class MemoryCallocIntTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.MemoryCallocIntTest验证calloc-int",sourcePrefix = "wasm")
    public void Memoryreallocint() {

        try {
            prepare();
            MemoryCallocInt mci = MemoryCallocInt.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = mci.getContractAddress();
            String transactionHash = mci.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MemoryCallocInt issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MemoryCallocInt deploy gasUsed:" + mci.getTransactionReceipt().get().getGasUsed());

            Int32 callocint = mci.getcalloc().send();
            collector.logStepPass("calloc函数返回值:" + callocint.value);
            collector.assertEqual(callocint.value, 2500);

        } catch (Exception e) {
            collector.logStepFail("MemoryCallocInt failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}

