package wasm.exec_efficiency;


import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.RecursionCall;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @title 执行效率-递归调用
 * @description:
 * @author: qcxiao
 * @create: 2019/12/26 14:38
 **/
public class WASMRecursionCallTest extends WASMContractPrepareTest {


    private Uint64 numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = Uint64.of(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.exec_efficiency.RecursionCallTest-递归执行", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {
            RecursionCall recursionCall = RecursionCall.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = recursionCall.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + recursionCall.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = RecursionCall.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .call(numberOfCalls).send();

            Uint64 sum = RecursionCall.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .get_sum().send();

            collector.logStepPass("sum:" + sum);
            collector.assertEqual(sum, numberOfCalls, "assert recursion call result");
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
