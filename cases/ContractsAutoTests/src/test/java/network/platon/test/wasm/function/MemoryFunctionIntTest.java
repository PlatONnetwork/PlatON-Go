package network.platon.test.wasm.function;

import com.platon.rlp.datatypes.Int32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MemoryFunctionInt;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * 验证内存memory函数-整型存储
 * 1.函数realloc
 * 2.函数free
 * 3.函数memset
 * @author liweic
 * @create 2020/02/19
 */

public class MemoryFunctionIntTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.MemoryFunctionIntTest验证memory函数-整型存储",sourcePrefix = "wasm")
    public void Memoryfunctionint() {

        try {
            prepare();
            MemoryFunctionInt memoryint = MemoryFunctionInt.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = memoryint.getContractAddress();
            String transactionHash = memoryint.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MemoryFunctionInt issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MemoryFunctionInt deploy gasUsed:" + memoryint.getTransactionReceipt().get().getGasUsed());

            Int32 mallocint = memoryint.getmallocint().send();
            collector.logStepPass("malloc函数返回值:" + mallocint.value);
            collector.assertEqual(mallocint.value, -1);

        } catch (Exception e) {
            collector.logStepFail("MemoryFunctionInt failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}

