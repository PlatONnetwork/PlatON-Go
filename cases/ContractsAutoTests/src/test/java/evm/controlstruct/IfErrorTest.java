package evm.controlstruct;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.IfError;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 *  for控制结构语法
 *
 *
 * @author hudenian
 * @dev 2020/1/6 11:09
 */
public class IfErrorTest extends ContractPrepareTest {


    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ForErrorTest-for控制结构测试", sourcePrefix = "evm")
    public void ifStruct() {
        try {

            IfError ifError = IfError.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = ifError.getContractAddress();
            TransactionReceipt tx = ifError.getTransactionReceipt().get();

            collector.logStepPass("IfErrorTest if控制结构返回的结构体如果是指针需要提前初始化");
            collector.logStepPass("deploy gas used:" + ifError.getTransactionReceipt().get().getGasUsed());

            //1.if控制结构验证 ifControlValue
            Boolean ifFlg = ifError.getIfControlRes().send();

            collector.logStepPass( "IfErrorTest 测试获取链上的结果是:" + ifFlg);
            collector.assertEqual(false,ifFlg);

            //2.if控制结构验证 ifControlValue
            Boolean ifFlg1 = ifError.getIfControlRes1().send();

            collector.logStepPass( "IfErrorTest 测试获取链上的结果是:" + ifFlg1);
            collector.assertEqual(false,ifFlg1);

        } catch (Exception e) {
            collector.logStepFail("IfErrorTest testCase process fail",e.toString());
            e.printStackTrace();
        }
    }
}
