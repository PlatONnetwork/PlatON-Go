package network.platon.test.wasm.exec_efficiency;

import com.platon.rlp.datatypes.Int32;
import com.platon.rlp.datatypes.Int64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.BubbleSort;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.Arrays;

/**
 * @title WASMBubbleSortTest
 * @description 执行效率 - 冒泡排序
 * @author liweic
 * @updateTime 2020/2/28 19:38
 */
public class WASMBubbleSortTest extends WASMContractPrepareTest {
    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.exec_efficiency-冒泡排序", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {

            Integer numberOfCalls = Integer.valueOf(driverService.param.get("numberOfCalls"));

            BubbleSort bubblesort = BubbleSort.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = bubblesort.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + bubblesort.getTransactionReceipt().get().getGasUsed());

            Int64[] arr = new Int64[numberOfCalls];

            int min = -1000, max = 2000;

            for (int i = 0; i < numberOfCalls; i++) {
                arr[i] = Int64.of(min + (int) (Math.random() * (max - min + 1)));
            }

            collector.logStepPass("before sort:" + Arrays.toString(arr));
            TransactionReceipt transactionReceipt = bubblesort.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .sort(arr, Int32.of(arr.length)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            Int64[] generationArr = bubblesort.load(contractAddress, web3j, transactionManager, provider, chainId).get_array().send();

            collector.logStepPass("after sort:" + Arrays.toString(generationArr));
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
