package lib;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LibraryStaticUsing;
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
 * @title 库引用类似引用static方法测试
 * 解释：如果L作为库的名称，f()是库L的函数，则可以通过L.f()的方式调用
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class LibraryStaticUsingTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "emitEvent",
            author = "albedo", showName = "lib.LibraryStaticUsingTest-类static方式引用")
    public void testEmitEvent() {
        try {
            prepare();
            LibraryStaticUsing using = LibraryStaticUsing.deploy(web3j, transactionManager, provider).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("LibraryStaticUsing issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = using.register(new BigInteger("12")).send();
            List<LibraryStaticUsing.ResultEventResponse> eventData = using.getResultEvents(receipt);
            String data = eventData.get(0).log.getData();
            collector.assertEqual(subHexData(data), subHexData("1"), "checkout static method using library function");
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
