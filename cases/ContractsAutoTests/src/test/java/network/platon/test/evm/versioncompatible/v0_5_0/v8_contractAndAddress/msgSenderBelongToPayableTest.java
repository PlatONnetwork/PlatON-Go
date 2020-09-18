package network.platon.test.evm.versioncompatible.v0_5_0.v8_contractAndAddress;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.MsgSenderBelongToPayable;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title  08-合约和地址
 *         验证msg.sender属于address payable类型
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class msgSenderBelongToPayableTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "DataLocationTest-存储区域", sourcePrefix = "evm")
    public void dataLocationTest() {
        try {

            MsgSenderBelongToPayable msgSenderBelongToPayable = MsgSenderBelongToPayable.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = msgSenderBelongToPayable.getContractAddress();
            TransactionReceipt tx = msgSenderBelongToPayable.getTransactionReceipt().get();

            collector.logStepPass("FunctionDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + msgSenderBelongToPayable.getTransactionReceipt().get().getGasUsed());

            String contractAddr = msgSenderBelongToPayable.getContractAddr().send();

            collector.assertEqual(contractAddress,contractAddr);

            String msgSenderAddr = msgSenderBelongToPayable.getMsgSenderAddr().send();

        } catch (Exception e) {
            collector.logStepFail("msgSenderBelongToPayableTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
