package function.functionVisibilityAndDecarations;


import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Selector;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.utils.DataChangeUtil;


/**
 * 验证public (或 external) 函数有一个特殊的成员selector, 它对应一个ABI 函数选择器.
 * @author liweic
 * @dev 2020/01/11 20:30
 */

public class SelectorTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.SelectorTest-Selector测试")
    public void selector() {
        try {
            Selector selector = Selector.deploy(web3j, transactionManager, provider).send();

            String contractAddress = selector.getContractAddress();
            TransactionReceipt tx = selector.getTransactionReceipt().get();
            collector.logStepPass("Selector deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //验证payable声明
            byte[] result = selector.f().send();
            String f = DataChangeUtil.bytesToHex(result);

            collector.logStepPass("selector：" + f);
            collector.assertEqual("b8c9d365",f);

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
