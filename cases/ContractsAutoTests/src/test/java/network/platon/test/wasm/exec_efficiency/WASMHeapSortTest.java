package network.platon.test.wasm.exec_efficiency;

import com.platon.rlp.datatypes.Int32;
import com.platon.rlp.datatypes.Int64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.HeapSort;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.Arrays;

/**
 * @title WASMHeapSortTest
 * @description 执行效率 - 堆排序
 * @author liweic
 * @updateTime 2020/3/2 15:09
 */
public class WASMHeapSortTest extends WASMContractPrepareTest {
    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.exec_efficiency-堆排序", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {

            Integer numberOfCalls = Integer.valueOf(driverService.param.get("numberOfCalls"));

            HeapSort heapsort = HeapSort.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = heapsort.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + heapsort.getTransactionReceipt().get().getGasUsed());

            Int64[] arr = new Int64[numberOfCalls];

            int min = -1000, max = 2000;

            for (int i = 0; i < numberOfCalls; i++) {
                arr[i] = Int64.of(min + (int) (Math.random() * (max - min + 1)));
            }

            collector.logStepPass("before sort:" + Arrays.toString(arr));
            TransactionReceipt transactionReceipt = heapsort.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .sort(arr, Int32.of(arr.length)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            Int64[] generationArr = heapsort.load(contractAddress, web3j, transactionManager, provider, chainId).get_array().send();

            collector.logStepPass("after sort:" + Arrays.toString(generationArr));
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}