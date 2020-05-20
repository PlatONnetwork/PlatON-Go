package evm.lib;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LibraryStaticUsing;
import network.platon.utils.DataChangeUtil;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.util.List;

/**
 * @title 库引用类似引用static方法测试
 * 解释：如果L作为库的名称，f()是库L的函数，则可以通过L.f()的方式调用
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class LibraryStaticUsingTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "emitEvent",
            author = "albedo", showName = "lib.LibraryStaticUsingTest-类static方式引用", sourcePrefix = "evm")
    public void testEmitEvent() {
        try {
            prepare();
            LibraryStaticUsing using = LibraryStaticUsing.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("LibraryStaticUsing issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + using.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = using.register(new BigInteger("12")).send();
            List<LibraryStaticUsing.ResultEventResponse> eventData = using.getResultEvents(receipt);
            String data = eventData.get(0).log.getData();
            collector.assertEqual(DataChangeUtil.subHexData(data), DataChangeUtil.subHexData("1"), "checkout static method using library function");
        } catch (Exception e) {
            collector.logStepFail("LibraryStaticUsingTest testEmitEvent failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
