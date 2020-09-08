package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MemoryFunction_2;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * 验证内存memory实现函数
 * 1.函数realloc
 * 2.函数free
 * 3.函数memset
 * 4.C库函数strcpy
 * @author liweic
 * @create: 2020/02/16
 */

public class MemoryFunction_2Test extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.MemoryFunction_2Test验证内存实现函数realloc",sourcePrefix = "wasm")
    public void Memoryfunction2() {

        try {
            prepare();
            MemoryFunction_2 memory2 = MemoryFunction_2.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = memory2.getContractAddress();
            String transactionHash = memory2.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MemoryFunction_2 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MemoryFunction_2 deploy gasUsed:" + memory2.getTransactionReceipt().get().getGasUsed());

            String realloc = memory2.getrealloc().send();
            collector.logStepPass("realloc函数返回值:" + realloc);
            collector.assertEqual(realloc, "WasmTest2");


        } catch (Exception e) {
            collector.logStepFail("MemoryFunction_2 failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
