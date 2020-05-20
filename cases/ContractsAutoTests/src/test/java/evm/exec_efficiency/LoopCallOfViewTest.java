package evm.exec_efficiency;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LoopCallOfView;
import org.junit.Before;
import org.junit.Test;

import java.math.BigInteger;
/**
 * @title LoopCallOfViewTest
 * @description 循环Call调用
 * @author qcxiao
 * @updateTime 2020/1/9 20:19
 */
public class LoopCallOfViewTest extends ContractPrepareTest {

    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "exec_efficiency.LoopCallOfView-循环执行", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {
            LoopCallOfView loopCall = LoopCallOfView.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = loopCall.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + loopCall.getTransactionReceipt().get().getGasUsed());

            BigInteger sum = LoopCallOfView.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .loopCallTest(numberOfCalls).send();

            collector.logStepPass("contract load successful. after the sum:" + sum);
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}