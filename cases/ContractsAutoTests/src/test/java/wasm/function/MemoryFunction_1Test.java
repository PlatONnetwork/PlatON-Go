package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MemoryFunction_1;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * 验证内存memory实现函数
 * 1.函数malloc
 * 2.函数free
 * 3.函数memset
 * 4.C库函数strcpy
 * @author liweic
 *@create: 2020/02/16
 */

public class MemoryFunction_1Test extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.MemoryFunction_1Test验证内存实现函数malloc",sourcePrefix = "wasm")
    public void Memoryfunction1() {

        try {
            prepare();
            MemoryFunction_1 memory1 = MemoryFunction_1.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = memory1.getContractAddress();
            String transactionHash = memory1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MemoryFunction_1 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MemoryFunction_1 deploy gasUsed:" + memory1.getTransactionReceipt().get().getGasUsed());

            String malloc = memory1.getmalloc().send();
            collector.logStepPass("malloc函数返回值:" + malloc);
            collector.assertEqual(malloc, "WasmTest");


        } catch (Exception e) {
            collector.logStepFail("MemoryFunctions_1 failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
