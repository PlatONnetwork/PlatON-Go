package evm.controlstruct;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.DoWhileError;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 *  dowhile控制结构语法
 *
 *
 * @author hudenian
 * @dev 2020/1/6 11:09
 */
public class DoWhileErrorTest extends ContractPrepareTest {


    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ControlTest-dowhile控制结构测试", sourcePrefix = "evm")
    public void doWhileStruct() {
        try {

            DoWhileError doWhileError = DoWhileError.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = doWhileError.getContractAddress();
            TransactionReceipt tx = doWhileError.getTransactionReceipt().get();

            collector.logStepPass("DoWhileErrorTest dowhile控制结构返回的结构体如果是指针需要提前初始化");
            collector.logStepPass("deploy gas used:" + doWhileError.getTransactionReceipt().get().getGasUsed());

            //1.if控制结构验证 ifControlValue
            Boolean doWhileFlg = doWhileError.getDoWhileControlRes().send();

            collector.logStepPass( "DoWhileErrorTest 测试获取链上的结果是:" + doWhileFlg);
            collector.assertEqual(false,doWhileFlg);

        } catch (Exception e) {
            collector.logStepFail("DoWhileErrorTest process fail.", e.toString());
            e.printStackTrace();
        }
    }
}
