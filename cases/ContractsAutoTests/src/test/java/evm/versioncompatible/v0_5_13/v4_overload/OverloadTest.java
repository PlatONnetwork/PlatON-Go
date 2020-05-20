package evm.versioncompatible.v0_5_13.v4_overload;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Overload;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title   带有指定参数的函数调用可以处理重载函数
 * @description:
 * @author: hudenian
 * @create: 2019/12/27
 */
public class OverloadTest extends ContractPrepareTest {

    //合约中的期望值
    private String expectValue="2";

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "OverloadTest-带有指定参数的函数调用可以处理重载函数", sourcePrefix = "evm")
    public void testStringMapping() {
        try {

            Overload overload = Overload.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = overload.getContractAddress();
            TransactionReceipt tx = overload.getTransactionReceipt().get();

            collector.logStepPass("OverloadTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + overload.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = overload.g().send();

            collector.logStepPass("StringmappingSupportTest testMapping successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String chainRe = overload.getRe().send().toString();

            collector.assertEqual(expectValue,chainRe);

        } catch (Exception e) {
            collector.logStepFail("OverloadTest  process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
