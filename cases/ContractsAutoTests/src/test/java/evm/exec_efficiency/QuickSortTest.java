package evm.exec_efficiency;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.QuickSort;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Random;

/**
 * @title QuickSortTest
 * @description 快速排序
 * @author zjsunzone
 * @updateTime 2020-02-26 14:25:27
 */
public class QuickSortTest extends ContractPrepareTest {

    private BigInteger numberOfCalls;   // length of array.
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "exec_efficiency.QuickSort-快速排序", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {
            // prepare array
            List<BigInteger> array = new ArrayList<>();
            Random r = new Random();
            for (int i = 0; i < numberOfCalls.intValue(); i++) {
                int val = r.nextInt(2000);
                if(i % 2 == 0){
                    String nfval = "-" + val;
                    int nval = Integer.parseInt(nfval);
                    array.add(BigInteger.valueOf(nval));
                } else {
                    array.add(BigInteger.valueOf(val));
                }
            }

            QuickSort contract = QuickSort.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = contract.getContractAddress();
            collector.logStepPass("QuickSort contract deploy successful. contractAddress:" + contractAddress + " hash:" + contract.getTransactionReceipt().get().getTransactionHash());
            collector.logStepPass("QuickSort contract deploy successful. gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            TransactionReceipt transactionReceipt = contract.sort(array, BigInteger.ZERO, BigInteger.valueOf(array.size() - 1)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("QuickSort sort successful, gasUsed:" + gasUsed);
            collector.logStepPass("QuickSort sort successful. hash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("QuickSort currentBlockNumber:" + transactionReceipt.getBlockNumber());

            List<BigInteger> afterArray = contract.show().send();
            collector.logStepPass("QuickSort sort before:" + Arrays.toString(array.toArray()));
            collector.logStepPass("QuickSort sort after :" + Arrays.toString(afterArray.toArray()));
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
