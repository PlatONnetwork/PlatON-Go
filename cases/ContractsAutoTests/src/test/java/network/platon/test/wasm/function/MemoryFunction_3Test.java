package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MemoryFunction_3;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * 验证内存memory实现函数
 * 1.函数calloc
 * 2.函数free
 * 3.函数memcpy
 * @author liweic
 * @create: 2020/02/17
 */

public class MemoryFunction_3Test extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.MemoryFunction_3Test验证内存实现函数calloc",sourcePrefix = "wasm")
    public void Memoryfunction3() {

        try {
            prepare();
            MemoryFunction_3 memory3 = MemoryFunction_3.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = memory3.getContractAddress();
            String transactionHash = memory3.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MemoryFunction_3 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MemoryFunction_3 deploy gasUsed:" + memory3.getTransactionReceipt().get().getGasUsed());

            String calloc = memory3.getcalloc().send();
            collector.logStepPass("calloc函数返回值:" + calloc);
            collector.assertEqual(calloc, "WasmTest3");


        } catch (Exception e) {
            collector.logStepFail("MemoryFunction_3 failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
