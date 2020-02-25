package wasm.exec_efficiency;

import com.platon.rlp.datatypes.Int8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InsertSort;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;
import java.math.BigInteger;

/**
 * @title WASMInsertSortTest
 * @description 执行效率 - 插入排序
 * @author qcxiao
 * @updateTime 2020/2/25 11:38
 */
public class WASMInsertSortTest extends WASMContractPrepareTest {
    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.exec_efficiency-空间复杂度", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {
            InsertSort insertSort = InsertSort.deploy(web3j, transactionManager, provider).send();
            contractAddress = insertSort.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);

            Int8[] arr = new Int8[]{Int8.of((byte) 1), Int8.of((byte) -1), Int8.of((byte) 5),
                    Int8.of((byte) 8), Int8.of((byte) 10), Int8.of((byte) 11), Int8.of((byte) 20),
                    Int8.of((byte) 30), Int8.of((byte) 32), Int8.of((byte) 127)};
            TransactionReceipt transactionReceipt = InsertSort.load(contractAddress, web3j, transactionManager, provider)
                    .sort(arr, Int8.of((byte) arr.length)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            Int8[] generationArr = InsertSort.load(contractAddress, web3j, transactionManager, provider).get_array().send();

            for (Int8 ele : generationArr) {
                System.out.print(ele.value + ",");
            }
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
