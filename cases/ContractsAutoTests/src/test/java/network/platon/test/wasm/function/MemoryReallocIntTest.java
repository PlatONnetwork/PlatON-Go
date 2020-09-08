package network.platon.test.wasm.function;

import com.platon.rlp.datatypes.Int32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MemoryReallocInt;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * 验证内存memory函数-整型存储
 * 1.函数realloc
 * @author liweic
 * @create 2020/02/24
 */

public class MemoryReallocIntTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.MemoryReallocIntTest验证realloc-int",sourcePrefix = "wasm")
    public void Memoryreallocint() {

        try {
            prepare();
            MemoryReallocInt mri = MemoryReallocInt.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = mri.getContractAddress();
            String transactionHash = mri.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MemoryReallocInt issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MemoryReallocInt deploy gasUsed:" + mri.getTransactionReceipt().get().getGasUsed());

            Int32 reallocint = mri.getrealloc().send();
            collector.logStepPass("realloc函数返回值:" + reallocint.value);
            collector.assertEqual(reallocint.value, 100);

        } catch (Exception e) {
            collector.logStepFail("MemoryReallocInt failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}

