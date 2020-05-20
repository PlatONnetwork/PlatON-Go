package evm.versioncompatible.v0_5_13.v2_publicOverExternal;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.PersonPublic;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title   允许public函数覆盖external函数
 * @description:
 * @author: hudenian
 * @create: 2019/12/27
 */
public class PublicCoverExternalTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "PublicCoverExternalTest-允许public函数覆盖external函数", sourcePrefix = "evm")
    public void testPublicCoverExternal() {
        try {

            PersonPublic personPublic = PersonPublic.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = personPublic.getContractAddress();
            TransactionReceipt tx = personPublic.getTransactionReceipt().get();

            collector.logStepPass("PublicCoverExternalTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + personPublic.getTransactionReceipt().get().getGasUsed());

            String birthDay = personPublic.birthDay().send().toString();

            collector.logStepPass("获取到的birthDay的值为："+birthDay);

            collector.assertEqual("2020-12-15",birthDay);


        } catch (Exception e) {
            collector.logStepFail("PublicCoverExternalTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
