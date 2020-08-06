package evm.controlstruct;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ForError;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 *  for控制结构语法
 *
 *
 * @author hudenian
 * @dev 2020/1/6 18:09
 */
public class ForErrorTest extends ContractPrepareTest {


    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ForErrorTest-for控制结构测试", sourcePrefix = "evm")
    public void forStruct() {
        try {

            ForError forError = ForError.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = forError.getContractAddress();
            TransactionReceipt tx = forError.getTransactionReceipt().get();

            collector.logStepPass("ForErrorTest for控制结构返回的结构体如果是指针需要提前初始化");
            collector.logStepPass("deploy gas used:" + forError.getTransactionReceipt().get().getGasUsed());

            //1.for控制结构验证
            Boolean forFlg = forError.getForControlRes().send();

            collector.logStepPass( "ForErrorTest->getForControlRes 测试获取链上的结果是:" + forFlg);
            collector.assertEqual(false,forFlg);

            //2.for控制结构验证
            Boolean forFlg1 = forError.getForControlRes1().send();

            collector.logStepPass( "ForErrorTest->getForControlRes1测试获取链上的结果是:" + forFlg1);
            collector.assertEqual(false,forFlg1);

        } catch (Exception e) {
            collector.logStepFail("ForErrorTest testCase process fail",e.toString());
            e.printStackTrace();
        }
    }
}
