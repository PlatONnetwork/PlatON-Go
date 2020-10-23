package network.platon.test.evm.exec_efficiency;


import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.BubbleSort;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;


/**
 * @title 冒泡排序
 * @description:
 * @author: liweic
 * @create: 2020/03/02
 **/
public class BubbleSortTest extends ContractPrepareTest {

    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "exec_efficiency.BubbleSortTest-冒泡排序", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {

            Integer numberOfCalls = Integer.valueOf(driverService.param.get("numberOfCalls"));

            BubbleSort bubblesort = BubbleSort.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = bubblesort.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + bubblesort.getTransactionReceipt().get().getGasUsed());

            List<BigInteger> array = new ArrayList<>(numberOfCalls);

            int min = -1000, max = 2000;
            for (int i = 0; i < numberOfCalls; i++) {
                BigInteger a = BigInteger.valueOf(min + (int) (Math.random() * 3001));
                array.add(a);
            }

            collector.logStepPass("before sort:" + array.toString());
            BigInteger n = BigInteger.valueOf(array.size());
            TransactionReceipt transactionReceipt = bubblesort.BubbleArrays(array, n, new BigInteger("1000")).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            List resultarr = bubblesort.get_arr().send();
            collector.logStepPass("after sort:" + resultarr);

        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

    /**
     * int 数组转 BigInteger
     * @param intArray
     * @return
     */
    private List<BigInteger> changeIntToBigInteger(Integer[] intArray) {
        List<BigInteger> bigIntegerList = new ArrayList<>();
        for (int i = 0; i <intArray.length ; i++) {
            bigIntegerList.add(new BigInteger(String.valueOf(intArray[i])));
        }
        return bigIntegerList;
    }

}

