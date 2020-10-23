package network.platon.test.evm.exec_efficiency;


import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.RecursionCall;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 执行效率-递归调用
 * @description:
 * @author: qcxiao
 * @create: 2019/12/26 14:38
 **/
public class RecursionCallTest extends ContractPrepareTest {


    private BigInteger numberOfCalls;
    private String contractAddress;

    @Before
    public void before() {
        numberOfCalls = new BigInteger(driverService.param.get("numberOfCalls"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "exec_efficiency.RecursionCallTest-递归执行", sourcePrefix = "evm")
    public void test() {
        prepare();
        try {
            RecursionCall recursionCall = RecursionCall.deploy(web3j, transactionManager, provider, chainId).send();
            contractAddress = recursionCall.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + recursionCall.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = RecursionCall.load(contractAddress, web3j, transactionManager, provider, chainId)
                    .recursionCallTest(numberOfCalls, new BigInteger("0")).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());
            collector.logStepPass("get result value:" + RecursionCall.load(contractAddress, web3j, transactionManager, provider, chainId).get_total().send());
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}
