package exec_efficiency;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LoopCall;
import network.platon.contracts.LoopCallOfView;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

public class LoopCallOfViewTest extends ContractPrepareTest {

    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "exec_efficiency.LoopCallOfView-循环执行")
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