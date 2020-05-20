package evm.controlstruct;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ShortCircuitError;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 *  && || 短路语法
 *
 * @author hudenian
 * @dev 2020/1/6 19:54
 */
public class ShortCircuitErrorTest extends ContractPrepareTest {


    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ShortCircuitErrorTest-短路语法", sourcePrefix = "evm")
    public void shortCircuitStruct() {
        try {

            ShortCircuitError shortCircuitError = ShortCircuitError.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = shortCircuitError.getContractAddress();
            TransactionReceipt tx = shortCircuitError.getTransactionReceipt().get();

            collector.logStepPass("ShortCircuitError deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + shortCircuitError.getTransactionReceipt().get().getGasUsed());


            Boolean fFlg = shortCircuitError.getF().send();

            collector.logStepPass( "ShortCircuitErrorTest 测试获取链上的结果是:" + fFlg);
            collector.assertEqual(false,fFlg);


            Boolean gflg = shortCircuitError.getG().send();

            collector.logStepPass( "ShortCircuitErrorTest 测试获取链上的结果是:" + gflg);
            collector.assertEqual(false,gflg);

        } catch (Exception e) {
            collector.logStepFail("ShortCircuitErrorTest testCase process fail",e.toString());
            e.printStackTrace();
        }
    }
}
