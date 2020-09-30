package network.platon.test.wasm.exec_efficiency;

import com.platon.rlp.datatypes.Int64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.QuickSort;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;
import java.math.BigInteger;
import java.util.Arrays;

/**
 * @title WASMQuickSortTest
 * @description 快排
 * @author qcxiao
 * @updateTime 2019/12/28 14:39
 */
public class WASMQuickSortTest extends WASMContractPrepareTest {
    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.exec_efficiency-快速排序", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {
            Integer numberOfCalls = Integer.valueOf(driverService.param.get("numberOfCalls"));
            QuickSort quickSort = QuickSort.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = quickSort.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + quickSort.getTransactionReceipt().get().getGasUsed());

            Int64[] arr = new Int64[numberOfCalls];

            // 定义一个数组中的最小值与最大值
            int min = -1000, max = 2000;

            // numberOfCalls表示数组的长度
            for (int i = 0; i < numberOfCalls; i++) {
                arr[i] = Int64.of(min + (int) (Math.random() * (max - min + 1)));
            }

            collector.logStepPass("before sort:" + Arrays.toString(arr));
            TransactionReceipt transactionReceipt = QuickSort.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .sort(arr, Int64.of(0), Int64.of(arr.length)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            Int64[] generationArr = QuickSort.load(contractAddress, web3j, transactionManager, provider, chainId).get_array().send();

            collector.logStepPass("after sort:" + Arrays.toString(generationArr));
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
