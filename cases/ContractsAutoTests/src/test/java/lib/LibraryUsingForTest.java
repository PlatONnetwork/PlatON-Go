package lib;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LibraryUsingFor;
import org.apache.commons.lang.StringUtils;
import org.junit.Test;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.PlatonGetTransactionCount;
import org.web3j.protocol.core.methods.response.PlatonSendTransaction;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.tx.response.PollingTransactionReceiptProcessor;
import org.web3j.utils.Numeric;

import java.io.IOException;
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
            author = "albedo", showName = "lib.LibraryUsingForTest-using A for B")
    public void testRegister() {
        try {
            prepare();
            LibraryUsingFor using = LibraryUsingFor.deploy(web3j, transactionManager, provider).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("LibraryUsingFor issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = using.register(new BigInteger("12")).send();
            List<LibraryUsingFor.ResultEventResponse> eventData = using.getResultEvents(receipt);
            String data = eventData.get(0).log.getData();
            collector.assertEqual(subHexData(data), subHexData("1"), "checkout using A for B");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
    private String subHexData(String hexStr) {
        if (StringUtils.isBlank(hexStr)) {
            throw new IllegalArgumentException("string is blank");
        }
        if (StringUtils.startsWith(hexStr, "0x")) {
            hexStr = StringUtils.substringAfter(hexStr, "0x");
        }
        byte[] addi = hexStr.getBytes();
        for (int i = 0; i < addi.length; i++) {
            if (addi[i] != 0) {
                hexStr = StringUtils.substring(hexStr, i - 1);
                break;
            }
        }
        return hexStr;
    }

}
