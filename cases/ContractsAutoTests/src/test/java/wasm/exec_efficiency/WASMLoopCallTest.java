package wasm.exec_efficiency;


import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.LoopCall;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @title 执行效率-循环调用
 * @description:
 * @author: qcxiao
 * @create: 2019/12/26 14:38
 **/
public class WASMLoopCallTest extends WASMContractPrepareTest {

    private Long numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new Long(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.exec_efficiency.LoopCallTest-循环执行", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {
            LoopCall loopCall = LoopCall.deploy(web3j, transactionManager, provider).send();
            contractAddress = loopCall.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);

            TransactionReceipt transactionReceipt = LoopCall.load(contractAddress, web3j, transactionManager, provider)
                    .loopCallTest(numberOfCalls).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());
            Long sum = LoopCall.load(contractAddress, web3j, transactionManager, provider).get_sum().send();
            collector.logStepPass("computing result:" + sum);
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
