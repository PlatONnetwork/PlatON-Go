package wasm.exec_efficiency;

import com.platon.rlp.datatypes.Int8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SpaceComplexity;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;
import java.math.BigInteger;

/**
 * @title SpaceComplexityTest
 * @description 空间复杂度场景测试
 * @author qcxiao
 * @updateTime 2019/12/28 14:39
 */
public class WASMSpaceComplexityTest extends WASMContractPrepareTest {
    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.exec_efficiency-空间复杂度", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {
            SpaceComplexity spaceComplexity = SpaceComplexity.deploy(web3j, transactionManager, provider).send();
            contractAddress = spaceComplexity.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);

            Int8[] arr = new Int8[]{Int8.of((byte) 1), Int8.of((byte) -1), Int8.of((byte) 5),
                    Int8.of((byte) 8), Int8.of((byte) 10), Int8.of((byte) 11), Int8.of((byte) 20),
                    Int8.of((byte) 30), Int8.of((byte) 32), Int8.of((byte) 127)};
            TransactionReceipt transactionReceipt = SpaceComplexity.load(contractAddress, web3j, transactionManager, provider)
                    .sort(arr, Int8.of((byte) 0), Int8.of((byte) 9)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            Int8[] generationArr = SpaceComplexity.load(contractAddress, web3j, transactionManager, provider).get_array().send();

            for (Int8 ele : generationArr) {
                System.out.print(ele.value + ",");
            }
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
