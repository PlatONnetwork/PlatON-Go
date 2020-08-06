package evm.lib;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LibraryUsingFor;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.List;

/**
 * @title 引用using for方式验证
 * 解释：指令using A for B 可用于附加库函数（从库 A）到任何类型（B）。 这些函数将接收到调用它们的对象作为它们的第一个参数。
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class LibraryUsingForTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "lib.LibraryUsingForTest-using A for B", sourcePrefix = "evm")
    public void testRegister() {
        try {
            prepare();
            LibraryUsingFor using = LibraryUsingFor.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("LibraryUsingFor issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + using.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = using.register(new BigInteger("12")).send();
            List<LibraryUsingFor.ResultEventResponse> eventData = using.getResultEvents(receipt);
            String data = eventData.get(0).log.getData();
            collector.assertEqual(DataChangeUtil.subHexData(data), DataChangeUtil.subHexData("1"), "checkout using A for B");
        } catch (Exception e) {
            collector.logStepFail("LibraryUsingForTest testRegister method failure:",e.getMessage());
            e.printStackTrace();
        }
    }


}
