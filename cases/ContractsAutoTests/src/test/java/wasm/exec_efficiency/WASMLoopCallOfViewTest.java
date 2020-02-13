package wasm.exec_efficiency;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LoopCallOfView;
import org.junit.Before;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @title LoopCallOfViewTest
 * @description 循环Call调用
 * @author qcxiao
 * @updateTime 2020/1/9 20:19
 */
public class WASMLoopCallOfViewTest extends WASMContractPrepareTest {

    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.exec_efficiency.LoopCallOfView-循环执行", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {
            LoopCallOfView loopCall = LoopCallOfView.deploy(web3j, transactionManager, provider).send();
            contractAddress = loopCall.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);

            BigInteger sum = LoopCallOfView.load(contractAddress, web3j, transactionManager, provider)
                    .loopCallTest(numberOfCalls).send();

            collector.logStepPass("contract load successful. after the sum:" + sum);
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}