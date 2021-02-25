package network.platon.test.evm.controlstruct;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.WhileError;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 *  for控制结构语法
 *
 *
 * @author hudenian
 * @dev 2020/1/7 13:42
 */
public class WhileErrorTest extends ContractPrepareTest {


    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "WhileErrorTest-while控制结构测试", sourcePrefix = "evm")
    public void whileStruct() {
        try {

            WhileError whileError = WhileError.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = whileError.getContractAddress();
            TransactionReceipt tx = whileError.getTransactionReceipt().get();

            collector.logStepPass("WhileErrorTest for控制结构返回的结构体如果是指针类型需要提前初始化");
            collector.logStepPass("deploy gas used:" + whileError.getTransactionReceipt().get().getGasUsed());

            //1.while控制结构验证
            Boolean whileFlg = whileError.getWhileControlRes().send();

            collector.logStepPass( "WhileErrorTest->whileStruct 测试获取链上的结果是:" + whileFlg);
            collector.assertEqual(false,whileFlg);

        } catch (Exception e) {
            collector.logStepFail("WhileErrorTest testCase process fail",e.toString());
            e.printStackTrace();
        }
    }
}
