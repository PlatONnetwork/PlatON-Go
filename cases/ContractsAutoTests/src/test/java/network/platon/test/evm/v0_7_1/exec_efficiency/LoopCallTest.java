package network.platon.test.evm.v0_7_1.exec_efficiency;


import com.alaya.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import org.junit.Before;
import org.junit.Test;
import java.math.BigInteger;

/**
 * @title 执行效率-循环调用
 * @description:
 * @author: qcxiao
 * @create: 2019/12/26 14:38
 **/
public class LoopCallTest extends ContractPrepareTest {

    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "exec_efficiency.LoopCallTest-循环执行", sourcePrefix = "evm/0.7.1")
    public void test() {
        prepare();
        try {
            LoopCall loopCall = LoopCall.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = loopCall.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + loopCall.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = LoopCall.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .loopCallTest(numberOfCalls).send();

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
