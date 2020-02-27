package evm.exec_efficiency;


import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.InsertSort;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;


/**
 * @title 插入排序
 * @description:
 * @author: liweic
 * @create: 2020/02/26
 **/
public class InsertSortTest extends ContractPrepareTest {

    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "exec_efficiency.InsertSortTest-插入排序", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {

            Integer numberOfCalls = Integer.valueOf(driverService.param.get("numberOfCalls"));

            InsertSort insertSort = InsertSort.deploy(web3j, transactionManager, provider).send();
            contractAddress = insertSort.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);

            List<BigInteger> array = new ArrayList<>(numberOfCalls);

            int min = -1000, max = 2000;

            for (int i = 0; i < numberOfCalls; i++) {
                BigInteger a = BigInteger.valueOf(min + (int) (Math.random() * 3001));
                array.add(a);
            }

            BigInteger n = BigInteger.valueOf(array.size());
            TransactionReceipt transactionReceipt = insertSort.OuputArrays(array, n, new BigInteger("1000")).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());


        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}

