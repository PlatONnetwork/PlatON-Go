package wasm.exec_efficiency;

import com.platon.rlp.Int64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SpaceComplexity;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.Arrays;

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
            author = "qcxiao", showName = "wasm.exec_efficiency.SpaceComplexityTest-空间复杂度", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {
            SpaceComplexity spaceComplexity = SpaceComplexity.deploy(web3j, transactionManager, provider).send();
            contractAddress = spaceComplexity.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);

            Int64[] arr = new Int64[]{Int64.of(1), Int64.of(1000)};
            TransactionReceipt transactionReceipt = SpaceComplexity.load(contractAddress, web3j, transactionManager, provider)
                    .sort(arr, Int64.of(1), Int64.of(1000)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            Int64[] generationArr = SpaceComplexity.load(contractAddress, web3j, transactionManager, provider).get_array().send();

            System.out.println(Arrays.toString(generationArr));
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
